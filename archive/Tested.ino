#include <WiFiNINA.h>
#include <ArduinoHttpClient.h>
#include <ArduinoJson.h>
#include <math.h>

// Equation for converting the raw ADC value to dB with calibrated regression values
// adc = analogRead(soundSensorPin); // Read the ADC value
// dB = (adc + b) / a // Convert the raw ADC value to dB using regression values
// TO convert threshold to ADC (adc = a * dB - b)
// a and b will be calculated in the calibration (excel file i created)
// Excel file: https://centriafi-my.sharepoint.com/:x:/g/personal/rami_kontio_centria_fi/IQCu7gqctI9VSJnb9UOq4Mi4AS5kDfW7-6MEfM4wVLJJMRg?e=3kGB1b
// Source: https://circuitdigest.com/microcontroller-projects/arduino-sound-level-measurement

// ===== Wi-Fi & Server =====
void fetchThreshold();
void sendData(const char* idValue, double soundLevelValue, double thresholdValue,
              const String& descriptionValue, bool isPeriodicValue, double fiveMinAvgValue);
String base64Encode(String input);

const char* ssid = "-";
const char* password = "-";
char server[] = "-.-.-.-";
int port = 8080;

const char* deviceID = "arduino_001";

const int ledPin = 2;
const int soundSensorPin = A0;
int thresholdADC = 0;
double threshold = 80.0;   // Threshold in dB

WiFiClient wifi;
HttpClient client = HttpClient(wifi, server, port);

const char* username = "kids_noisemeter_admin";
const char* passwordAuth = "passwordkids";

unsigned long lastThresholdFetch = 0;
const unsigned long thresholdFetchInterval = 10UL * 1000UL;

unsigned long lastConstantSent = 0;
unsigned long checkConstantInterval = 3000;  // 3 sec

// ===== Block averaging =====
const unsigned long blockSeconds = 60UL;   // 1 minute blocks
const int blocksPerDemo = 5;               // 5 blocks = 5 minutes
double BlockLeqReg[blocksPerDemo];
int currentBlockIndex = 0;

// ===== Calibration constants =====
// dB = a * RMS + b
double a = 0.091;
double b = 74.1;

// EMA smoothing
double emaLinear = 0.0;
bool emaInitialized = false;
unsigned long lastSampleMicros = 0;
double sumEma = 0.0;
unsigned long emaCount = 0;
const double tau = 0.5;

// ===== Persistent 5-min average =====
double lastFiveMinAvg = 0.0;

// ===== RMS measurement function =====
double measureRMS(unsigned int samples = 400, unsigned int usDelay = 200) {
  long sum = 0;
  for (unsigned int i = 0; i < samples; i++) {
    sum += analogRead(soundSensorPin);
    delayMicroseconds(usDelay);
  }
  double mean = (double)sum / samples;

  double sqsum = 0.0;
  for (unsigned int i = 0; i < samples; i++) {
    double ac = analogRead(soundSensorPin) - mean;
    sqsum += ac * ac;
    delayMicroseconds(usDelay);
  }
  return sqrt(sqsum / samples);
}

void setup() {
  Serial.begin(9600);
  pinMode(ledPin, OUTPUT);
  pinMode(soundSensorPin, INPUT);

  Serial.println("=== Kindergarten Sound Warning System ===");
  Serial.print("Connecting to WiFi");

  WiFi.begin(ssid, password);
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }
  Serial.println("\nConnected to WiFi!");
}

void loop() {
  unsigned long blockDurationMicros = blockSeconds * 1000000UL;
  unsigned long start = micros();

  emaLinear = 0.0;
  emaInitialized = false;
  sumEma = 0.0;
  emaCount = 0;
  lastSampleMicros = 0;

  // 1 minute block
  while ((unsigned long)(micros() - start) < blockDurationMicros) {
    // Fetch threshold periodically
    if (millis() - lastThresholdFetch >= thresholdFetchInterval || lastThresholdFetch == 0) {
      fetchThreshold();
      lastThresholdFetch = millis();
    }

    // === Read RMS sound sensor ===
    double rms = measureRMS();
    double L_inst_dB = a * rms + b;

    // === LED logic with delay ===
    static unsigned long ledOnStart = 0;
    const unsigned long ledDelay = 2000;  // require 2 seconds above threshold

    if (L_inst_dB > threshold) {
      if (ledOnStart == 0) {
        ledOnStart = millis();
      }
      if (millis() - ledOnStart >= ledDelay) {
        digitalWrite(ledPin, HIGH);
      }
    } else {
      ledOnStart = 0;
      digitalWrite(ledPin, LOW);
    }

    // === Serial logging ===
    Serial.println("------ Measurement ------");
    Serial.print("RMS Value       : "); Serial.println(rms, 2);
    Serial.print("Sound Level     : "); Serial.print(L_inst_dB, 1); Serial.println(" dB");
    Serial.print("Threshold       : "); Serial.print(threshold, 1); Serial.println(" dB");
    Serial.print("LED Status      : "); Serial.println(digitalRead(ledPin) == HIGH ? "ON" : "OFF");
    Serial.print("Last 5-min Avg  : "); Serial.println(lastFiveMinAvg, 2);
    Serial.println("-------------------------");

    // EMA smoothing
    unsigned long nowMicros = micros();
    double dt = 0.02;
    if (lastSampleMicros != 0) {
      dt = (nowMicros - lastSampleMicros) / 1e6;
    }
    lastSampleMicros = nowMicros;

    double instLinear = pow(10.0, L_inst_dB / 10.0);
    double alpha = exp(-dt / tau);

    if (!emaInitialized) {
      emaLinear = instLinear;
      emaInitialized = true;
    } else {
      emaLinear = alpha * emaLinear + (1.0 - alpha) * instLinear;
    }

    sumEma += emaLinear;
    emaCount++;

    // Send latest measurement every 3 sec, always include lastFiveMinAvg
    if (millis() - lastConstantSent >= checkConstantInterval || lastConstantSent == 0) {
      sendData(deviceID, L_inst_dB, threshold, "latest measurement", false, lastFiveMinAvg);
      lastConstantSent = millis();
    }
  }

  // === Block average (1 minute) ===
  double blockLeqReg_dB = 0.0;
  if (emaCount > 0) {
    double avgLinear = sumEma / (double)emaCount;
    if (avgLinear <= 0.0) avgLinear = 1e-20;
    blockLeqReg_dB = 10.0 * log10(avgLinear);
  }

  Serial.print("Block average (dB): ");
  Serial.println(blockLeqReg_dB, 2);

  if (currentBlockIndex < blocksPerDemo) {
    BlockLeqReg[currentBlockIndex] = blockLeqReg_dB;
    currentBlockIndex++;
  }

  // === Five-minute average after 5 blocks ===
  if (currentBlockIndex >= blocksPerDemo) {
    double sumLinearReg = 0.0;
    for (int i = 0; i < blocksPerDemo; ++i) {
      sumLinearReg += pow(10.0, BlockLeqReg[i] / 10.0);
    }
    double avgLinearReg = sumLinearReg / blocksPerDemo;
    double fiveMinLeq_dB = 10.0 * log10(avgLinearReg);

    Serial.print("Five-minute average (dB): ");
    Serial.println(fiveMinLeq_dB, 2);

    // Update lastFiveMinAvg
    lastFiveMinAvg = fiveMinLeq_dB;

    // Send five-minute average to server
    sendData(deviceID, fiveMinLeq_dB, threshold, "Five-minute average", true, lastFiveMinAvg);

    currentBlockIndex = 0;
  }

  // Send block data to server, always include lastFiveMinAvg
  sendData(deviceID, blockLeqReg_dB, threshold, "Block average over 1 minute", true, lastFiveMinAvg);
}

// ===== sendData function =====
void sendData(const char* idValue,
              double soundLevelValue,
              double thresholdValue,
              const String& descriptionValue,
              bool isPeriodicValue,
              double fiveMinAvgValue) {
  if (WiFi.status() == WL_CONNECTED) {
    String jsonData = "{";
    jsonData += "\"device_id\":\"" + String(idValue) + "\",";
    jsonData += "\"sound_level\":" + String(soundLevelValue, 2) + ",";
    jsonData += "\"threshold\":" + String(thresholdValue, 1) + ",";
    jsonData += "\"description\":\"" + descriptionValue + "\",";
    jsonData += "\"is_periodic\":" + String(isPeriodicValue ? "true" : "false") + ",";
    jsonData += "\"five_minute_average_sound_level\":" + String(fiveMinAvgValue, 2) + ",";
    jsonData += "\"led_status\":\"" + String(digitalRead(ledPin) == HIGH ? "ON" : "OFF") + "\"";
    jsonData += "}";

        client.beginRequest();
    client.post("/api/data/hourly");
    String auth = String(username) + ":" + String(passwordAuth);
    client.sendHeader("Authorization", "Basic " + base64Encode(auth));
    client.sendHeader("Content-Type", "application/json");
    client.sendHeader("Content-Length", jsonData.length());
    client.beginBody();
    client.print(jsonData);
    client.endRequest();

    // Optional: read server response for debugging
    int statusCode = client.responseStatusCode();
    String response = client.responseBody();
    Serial.print("Server response code: ");
    Serial.println(statusCode);
    Serial.print("Server response body: ");
    Serial.println(response);
  }
}
// === For getting the current threshold ===
void fetchThreshold() {
  if (WiFi.status() != WL_CONNECTED) return;

  client.beginRequest();
  client.get("/api/locations/chosen");
  String auth = String(username) + ":" + String(passwordAuth);
  client.sendHeader("Authorization", "Basic " + base64Encode(auth));
  client.endRequest();

  int status = client.responseStatusCode();
  String body = client.responseBody();

  if (status == 200 && body.length() > 0) {
    StaticJsonDocument<512> doc;
    DeserializationError err = deserializeJson(doc, body);
    if (!err) {
      JsonVariant v = doc["data"]["threshold"];
      float t = v.as<float>();
      Serial.print("Current threshold: ");
      Serial.println(t);
      thresholdADC = round((t - b) / a);  // converts dB threshold into RMS threshold
      threshold = t;
    } else {
      Serial.print("JSON parse error: ");
      Serial.println(err.c_str());
    }
  }
}
/// === Base64 encoding function ===
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
  if (valb > -6) {
    encoded += base64Chars[((val << 8) >> (valb + 8)) & 0x3F];
  }
  while (encoded.length() % 4) {
    encoded += '=';
  }
  return encoded;
}


