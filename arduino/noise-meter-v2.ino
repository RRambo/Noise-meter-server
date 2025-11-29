#include <WiFiNINA.h>
#include <ArduinoHttpClient.h>
#include <ArduinoJson.h>
#include <math.h>

// waiting for testing

// Equation for converting the raw ADC value to dB with calibrated regression values
// adc = analogRead(soundSensorPin); // Read the ADC value
// dB = (adc + b) / a // Convert the raw ADC value to dB using regression values
// TO convert threshold to ADC (adc = a * dB - b)
// a and b will be calculated in the calibration (excel file i created)
// Excel file: https://centriafi-my.sharepoint.com/:x:/g/personal/rami_kontio_centria_fi/IQCu7gqctI9VSJnb9UOq4Mi4AS5kDfW7-6MEfM4wVLJJMRg?e=3kGB1b
// Source: https://circuitdigest.com/microcontroller-projects/arduino-sound-level-measurement

// ===== Wi-Fi & Server =====
// Prototypes
void fetchThreshold();
void sendData(const char* idValue, double soundLevelValue, double thresholdValue, const String& descriptionValue = "", bool isPeriodicValue = false);
String base64Encode(String input);
// ---------------------------------

// ===== Wi-Fi & Server =====
const char* ssid = "_";   // Used wifi username      
const char* password = "_";     // Used wifi password
char server[] = "_._._._";   // Server IP Address    
int port = 8080;

// ===== Device & Room Info =====
const char* deviceID = "arduino_001";

// ===== Pins & Sensor =====
const int ledPin = 2;
const int soundSensorPin = A0;
int thresholdADC = 0;             // Sound threshold as a raw ADC value (now mutable)
double threshold = 70;           // Sound threshold (now mutable)

WiFiClient wifi;
HttpClient client = HttpClient(wifi, server, port);

// ===== Authentication =====
const char* username = "kids_noisemeter_admin";
const char* passwordAuth = "passwordkids";

// ===== LED spike logic =====
unsigned long ledOnUntil = 0;

// ===== Block/hour accumulators =====
unsigned long blockStartMillis = 0;
unsigned long lastThresholdFetch = 0;
const unsigned long thresholdFetchInterval = 60UL * 1000UL; // threshold updated every 1 minute

unsigned long lastConstantSent = 0;
unsigned long checkConstantInterval = 3000; // 3 sec

const unsigned long blockSeconds = 600UL; // --> 10 minutes
const int blocksPerHour = 6;
double calibrationFactor = 1.0;
int currentBlockIndex = 0;
unsigned long samplesPerBlockEstimate = 0;
double BlockLeqReg[blocksPerHour];

// ===== sustained block peak state =====
int peakADC = 0;
double peak_dB = 0.0;
unsigned long candidateStartMs = 0;   // when candidate above peakADC started
int candidateMaxADC = 0;              // highest ADC seen while candidate active
const unsigned long sustainMs = 500;  // require 500 ms sustain to accept new peak

// Candidate must exceed current accepted peak by this margin to start.
// While candidate is active, dips down to (peakADC - peakMargin) are tolerated.
const int peakMargin = 6; // needs to be tuned while testing

// ===== regression values =====
double a = 6.82; // slope (ADC -> dB)
double b = -30.0; // intercept

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

  // initialize timers
  blockStartMillis = millis();
  lastThresholdFetch = 0;
}

void loop() {
  unsigned long blockDurationMicros = blockSeconds * 1000000UL;
  unsigned long start = micros();

  // regression-based dB (Leq) accumulators
  unsigned long sampleCount = 0;
  double sumEnergy = 0;

  // reset block peak state
  peakADC = 0;
  peak_dB = 0.0;
  candidateStartMs = 0;
  candidateMaxADC = 0;

  // 10 minute block
  while ((unsigned long)(micros() - start) < blockDurationMicros) {
    // Fetch threshold periodically
    if (millis() - lastThresholdFetch >= thresholdFetchInterval || lastThresholdFetch == 0) {
      fetchThreshold();
      lastThresholdFetch = millis();
    }

    // === Read & smooth sound sensor ===
    const int numSamples = 10;
    int total = 0;
    for (int i = 0; i < numSamples; i++) {
      total += analogRead(soundSensorPin);
      //Serial.print("Raw ADC: "); Serial.println(analogRead(soundSensorPin));
    }
    int adc = total / numSamples;

    // regression mapping (ensure a != 0)
    // ^ don't know why this needs to be done, since a and b are hardcoded constants
    double L_inst_dB = ((double)adc + b) / a;

    //Serial.print("values,   ADC: "); Serial.print(adc); Serial.print(" dB: "); Serial.println(L_inst_dB);

    // Accumulate energy for calculating the average (Leq)
    sumEnergy += pow(10.0, L_inst_dB / 10.0);
    sampleCount++;

    // === Sustained peak detection ===
    // For getting the highest peak within the 10 minute block
    // Start only if adc > peakADC + peakMargin to avoid unnatural noise
    if (adc > peakADC + peakMargin) {
      if (candidateStartMs == 0) {
        candidateStartMs = millis();
        candidateMaxADC = adc;
      } else {
        if (adc > candidateMaxADC) candidateMaxADC = adc;
      }
      // Promote candidate if sustained long enough (0.5 seconds)
      if ((unsigned long)(millis() - candidateStartMs) >= sustainMs) {
        peakADC = candidateMaxADC;
        peak_dB = ((double)peakADC + b) / a;
        // reset candidates
        candidateStartMs = 0;
        candidateMaxADC = 0;
      }
    } else if (candidateStartMs != 0) {
      // Candidate active: allow small dips down to (peakADC - peakMargin)
      // We keep the candidate running if ADC remains >= (peakADC - peakMargin)
      if (adc >= (peakADC - peakMargin)) {
        // Update candidateMaxADC if higher
        if (adc > candidateMaxADC) candidateMaxADC = adc;
      } else {
        // Candidate invalidated by a larger dip in noise
        candidateStartMs = 0;
        candidateMaxADC = 0;
      }
    }

    // === LED logic ===
    if (candidateMaxADC > thresholdADC) {
      if (millis() > ledOnUntil) {
        ledOnUntil = millis() + 300UL; // Minimum ON for short spike
      }
    }
    digitalWrite(ledPin, (millis() < ledOnUntil) ? HIGH : LOW);
    delay(10); // for reducing CPU load
    // Sending data every 3 seconds for current sound level card
    // No need to explicitlty send data to server whenever threshold is reached
    if (millis() - lastConstantSent >= checkConstantInterval || lastConstantSent == 0) {
      sendData(deviceID, L_inst_dB, threshold, "latest measurement", false);
      lastConstantSent = millis();
    }
  } // end 10 minute block sampling

  samplesPerBlockEstimate = sampleCount;

  // === Regression-based block Leq ===
  double blockLeqReg_dB = 0;
  if (sampleCount > 0) {
    double avgLinear = sumEnergy / (double)sampleCount;
    Serial.println();
    Serial.print("Block average (linear): "); Serial.println(avgLinear, 2);
    if (avgLinear <= 0.0) avgLinear = 1e-20;
    blockLeqReg_dB = 10.0 * log10(avgLinear);
  }
  Serial.println();
  Serial.print("Block average (regression): "); Serial.println(blockLeqReg_dB, 2);

  // === Peak ===
  Serial.print("Block peak (sustained >= "); Serial.print(sustainMs);
  Serial.print(" ms, margin "); Serial.print(peakMargin);
  Serial.print(" ADC) dB: "); Serial.println(peak_dB, 2);

  if (currentBlockIndex < blocksPerHour) BlockLeqReg[currentBlockIndex] = blockLeqReg_dB;

  // For logging the average / hour using BlockLeq
  // Doesn't get sent yet
  if (currentBlockIndex >= blocksPerHour) {
    // Hourly average using regression-based block values
    double sumLinearReg = 0.0;
    int validCountReg = 0;

    for (int i = 0; i < blocksPerHour; ++i) {
      sumLinearReg += pow(10.0, BlockLeqReg[i] / 10.0);
      validCountReg++;
    }
    if (validCountReg > 0) {
      double avgLinearReg = sumLinearReg / (double)validCountReg;
      double hourLeqReq_dB = 10.0 * log10(avgLinearReg);
      Serial.print("Leq 1h (regression, dB): "); Serial.println(hourLeqReq_dB, 2);
    }

    currentBlockIndex = 0;
  }

  delay(10);

  // === Send block data (10 minutes) to server ===
  // regression data for now. The data should be compared
  sendData(deviceID, blockLeqReg_dB, threshold, "An average of measurements in a " + String(blockSeconds / 60) + " minute block", true);
  sendData(deviceID, peak_dB, threshold, "A peak of measurements in a " + String(blockSeconds / 60) + " minute block", true);
}

// === function for sending data ===
void sendData(const char* idValue, double soundLevelValue, double thresholdValue, const String& descriptionValue, bool isPeriodicValue) {
  if (WiFi.status() == WL_CONNECTED) {
    String jsonData = "{";
    jsonData += "\"device_id\":\"" + String(idValue) + "\",";
    jsonData += "\"sound_level\":" + String(round(soundLevelValue)) + ",";
    jsonData += "\"threshold\":" + String(thresholdValue) + ",";
    jsonData += "\"description\":\"" + descriptionValue + "\",";
    jsonData += "\"is_periodic\":" + String(isPeriodicValue ? "true" : "false");
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
    /*
    int statusCode = client.responseStatusCode();
    String response = client.responseBody();
    Serial.print("Data post status: ");
    Serial.println(statusCode);
    Serial.print("Response: ");
    Serial.println(response);
    */
  }
}

// === function for getting the current threshold ===
void fetchThreshold() {
  if (WiFi.status() != WL_CONNECTED) return;
  client.beginRequest();
  client.get("/api/locations/chosen");
  String auth = String(username) + ":" + String(passwordAuth);
  client.sendHeader("Authorization", "Basic " + base64Encode(auth));
  client.endRequest();

  int status = client.responseStatusCode();
  String body = client.responseBody();
  //Serial.print("Threshold fetch status: ");
  //Serial.println(status);
  //Serial.print("Body: ");
  //Serial.println(body);
  if (status == 200 && body.length() > 0) {
    // json object capacity estimated to be small so 512 bytes should be enough
    StaticJsonDocument<512> doc;
    DeserializationError err = deserializeJson(doc, body);
    if (!err) {
      JsonVariant v = doc["data"]["threshold"];
      float t = v.as<float>();
      Serial.print("Current threshold: ");
      Serial.println(t);
      thresholdADC = round(a * t - b); // convert the threshold into ADC for comparison
      threshold = t;
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
