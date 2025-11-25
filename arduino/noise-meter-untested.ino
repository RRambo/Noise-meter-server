#include <WiFiNINA.h>
#include <ArduinoHttpClient.h>
#include <ArduinoJson.h>

// waiting for testing

// Equation for converting the raw ADC value to dB with calibrated regression values
// adc = analogRead(soundSensorPin); // Read the ADC value
// dB = (adc + b) / a // Convert the raw ADC value to dB using regression values
// TO convert threshold to ADC (adc = a * dB - b)
// a and b will be calculated in the calibration (excel file i created)
// Excel file: https://centriafi-my.sharepoint.com/:x:/g/personal/rami_kontio_centria_fi/IQCu7gqctI9VSJnb9UOq4Mi4AS5kDfW7-6MEfM4wVLJJMRg?e=3kGB1b
// Source: https://circuitdigest.com/microcontroller-projects/arduino-sound-level-measurement

// ===== Wi-Fi & Server =====
const char* ssid = "_";   // Used wifi username      
const char* password = "_";     // Used wifi password
char server[] = "_._._._";   // Server IP Address    
int port = 8080;

// ===== Device & Room Info =====
const char* deviceID = "arduino_001";     

// ===== Pins & Sensor =====
const int ledPin = 13;
const int soundSensorPin = A0;
const int threshold = 0;   // Sound threshold
const unsigned long checkInterval = 500; // 0.5 sec

WiFiClient wifi;
HttpClient client = HttpClient(wifi, server, port);

// ===== Authentication =====
const char* username = "kids_noisemeter_admin"; 
const char* passwordAuth = "passwordkids";
/*
// ===== Variables for 1-hour average =====
float runningAverage = 0;
unsigned long readingCount = 0;
unsigned long lastHourTime = 0;
const unsigned long oneHourMillis = 3600000; // 1 hour in ms or 60000 if needed to check hourlyAverage faster
*/
// ===== LED spike logic =====
unsigned long ledOnUntil = 0;

// ===== Block/hour accumulators =====
unsigned long blockStartMillis = 0;
unsigned long lastThresholdFetch = 0;
const unsigned long thresholdFetchInterval = 60UL * 1000UL; // <-- threshold gets updated every 1 minute

const unsigned long blockSeconds = 600UL; // --> 10 minutes
const int blocksPerHour = 6;
const double calibrationFactor = 1.0;
const double P0 = 20e-6;
const double P0_SQ = P0 * P0;

unsigned long samplesPerBlockEstimate = 0;
double BlockLeq[blocksPerHour];
int currentBlockIndex = 0;

// ===== regression values =====
// Yet to be calculated
const double a = 0;
const double b = 0;

void setup() {
  Serial.begin(9600);
  pinMode(ledPin, OUTPUT);
  pinMode(buzzerPin, OUTPUT);
  pinMode(soundSensorPin, INPUT);

  Serial.println("=== Kindergarten Sound Warning System ===");
  Serial.print("Connecting to WiFi");

  WiFi.begin(ssid, password);
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }

  Serial.println("\nConnected to WiFi!");
  Serial.print("IP Address: ");
  Serial.println(WiFi.localIP());
  lastHourTime = millis(); // initialize hourly timer

  // initialize timers
  blockStartMillis = millis();
  lastThresholdFetch = 0;
}

void loop() {
  unsigned long blockDurationMicros = blockSeconds * 1000000UL;
  unsigned long start = micros();

  unsigned long N = 0;
  double sum = 0.0;
  double sumSq = 0.0;

  // 10 minute block
  while ((unsigned long)(micros() - start) < blockDurationMicros) {
    // Fetch threshold periodically
    if (millis() - lastThresholdFetch >= thresholdFetchInterval || lastThresholdFetch == 0) {
      fetchThreshold();
      lastThresholdFetch = millis();
    }
    // === Read & smooth sound sensor ===
    const int numSamples = 20;
    int total = 0;
    for (int i = 0; i < numSamples; i++) {
      total += analogRead(soundSensorPin);
      delay(10);
    }
    // measured and rounded raw ADC sample
    int adc = total / numSamples; 

    double v = (double)adc;
    sum += v;
    sumSq += v * v;
    N++; // for counting the amount of samples in a 10 minute block

    // === LED logic ===
    if (adc > threshold) {
      if (millis() > ledOnUntil) {
        ledOnUntil = millis() + 300UL; // Minimum ON for short spike
      }
    }
    digitalWrite(ledPin, (millis() < ledOnUntil) ? HIGH : LOW);
    delayMicroSeconds(10); // for reducing CPU load
  }

  double mean = sum / (double)N;
  double meanSq = sumSq / (double)N;
  double ms_signal_adc2 = meanSq - mean * mean;
  double ms_pa2 = ms_signal_adc2 * (calibrationFactor * calibrationFactor);
  if (ms_pa2 <= 0.0) ms_pa2 = 1e-20;
  double blockLeq_dB = 10.0 * log10(ms_pa2 / P0_SQ);

  blockLeq[currentBlockIndex++] = blockLeq_dB;
  Serial.print("Block: "); Serial.println(currentBlockIndex);
  Serial.print(" Samples: "); Serial.println(samplesPerBlockEstimate);
  Serial.print("Leq dB: "); Serial.println(blockLeq_dB, 2);

  // For logging the average / hour
  // Doesn't get sent yet
  if (currentBlockIndex >= blocksPerHour) {
    double sumLinear = 0.0;
    for (int i = 0; i < blocksPerHour; ++i) sumLinear += pow(10.0, BlockLeq[i] / 10.0);
    double avgLinear = sumLinear / (double)blocksPerHour;
    double hourLeq_dB = 10.0 * log10(avgLinear);
    Serial.print("Leq 1h (dB): "); Serial.println(hourLeq_dB, 2);
    currentBlockIndex = 0;
  }

  // Converting the threshold back to dB
  threshold = (threshold + b) / a;

  delay(100);
  
  /*// === Update 1-hour running average ===
  readingCount++;
  runningAverage = ((runningAverage * (readingCount - 1)) + soundLevel) / readingCount;*/

  // === Send block data (10 minutes) to server ===
  if (WiFi.status() == WL_CONNECTED) {
    String jsonData = "{";
    jsonData += "\"device_id\":\"" + String(deviceID) + "\",";
    jsonData += "\"sound_level\":" + String(round(blockLeq_dB)) + ",";
    jsonData += "\"threshold\":" + String(threshold) + ",";
    jsonData += "\"description\":\"An average of measurements in a 10 minute block\"";
    jsonData += "}";

    client.beginRequest();
    client.post("/api/data"); 
    String auth = String(username) + ":" + String(passwordAuth);
    client.sendHeader("Authorization", "Basic " + base64Encode(auth));
    client.sendHeader("Content-Type", "application/json");
    client.sendHeader("Content-Length", jsonData.length());
    client.beginBody();
    client.print(jsonData);
    client.endRequest();

    int statusCode = client.responseStatusCode();
    String response = client.responseBody();
    Serial.print("Data post status: ");
    Serial.println(statusCode);
    Serial.print("Response: ");
    Serial.println(response);
  }

  /*// === Send 1-hour average if time passed ===
  if (millis() - lastHourTime >= oneHourMillis) {
    if (WiFi.status() == WL_CONNECTED) {
      String jsonHourly = "{";
      jsonHourly += "\"device_id\":\"" + String(deviceID) + "\",";
      jsonHourly += "\"room_name\":\"" + String(roomName) + "\",";
      jsonHourly += "\"hourly_average_sound_level\":" + String(runningAverage, 2) + ",";
      jsonHourly += "\"description\":\"1-hour average sound reading\"";
      jsonHourly += "}";

      client.beginRequest();
      client.post("/data/hourly"); // Hourly endpoint
      String auth = String(username) + ":" + String(passwordAuth);
      client.sendHeader("Authorization", "Basic " + base64Encode(auth));
      client.sendHeader("Content-Type", "application/json");
      client.sendHeader("Content-Length", jsonHourly.length());
      client.beginBody();
      client.print(jsonHourly);
      client.endRequest();

      int statusCode = client.responseStatusCode();
      String response = client.responseBody();
      Serial.print("Hourly Status code: ");
      Serial.println(statusCode);
      Serial.print("Hourly Response: ");
      Serial.println(response);
    }

    // Reset 1-hour counters
    runningAverage = 0;
    readingCount = 0;
    lastHourTime = millis();
  }*/

  //delay(checkInterval);
}

void getThreshold() {
  if (WiFi.status() !== WL_CONNECTED) return;
  client.beginRequest();
  client.get("/locations/chosen");
  String auth = String (username) + ":" + String(passwordAuth);
  client.sendHeader("Authorization", "Basic " + base64Encode(auth));
  client.endRequest();
  
  int status = client.responseStatusCode();
  String body = client.responseBody();
  Serial.print("Threshold fetch status: ");
  Serial.println(status);
  Serial.print("Body: ");
  Serial.println(body);
  if (status == 200 && body.length() > 0) {
    // json object capacity estimated to be small so 200 bytes should be enough
    StaticJsonDocument<200> doc;
    DeserializationError err = deserializeJson(doc, body);
    if (!err) {
      float t = doc["threshold"].as<float>();
      Serial.print("Current threshold: ");
      Serial.println(t);
      threshold = round(a * t - b); // convert the threshold into ADC for comparison
    } else {
      Serial.print("JSON parse error: ");
      Serial.println(err.c_str());
    }
  }
}

// === Base64 encoding function ===
String base64Encode(String input) {
  const char base64Chars[] = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";
  String encoded = "";
  int val = 0, valb = -6;
  for (unsigned char c : input) {
    val = (val << 8) + c;
    valb += 8;
    while (valb >= 0) {
      encoded += base64Chars[(val >> valb) & 0x3F];
      valb -= 6;
    }
  }
  if (valb > -6) encoded += base64Chars[((val << 8) >> (valb + 8)) & 0x3F];
  while (encoded.length() % 4) encoded += '=';
  return encoded;
}

