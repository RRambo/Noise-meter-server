#include <WiFiNINA.h>
#include <ArduinoHttpClient.h>

// ===== Wi-Fi & Server =====
const char* ssid = "_";   // Used wifi username      
const char* password = "_";     // Used wifi password
char server[] = "_._._._";   // Server IP Address    
int port = 8080;

// ===== Device & Room Info =====
const char* deviceID = "arduino_001";     
const char* roomName = "Class_A";         

// ===== Pins & Sensor =====
const int ledPin = 13;
const int buzzerPin = 8;
const int soundSensorPin = A0;
const int threshold = 100;   // Sound threshold
const unsigned long checkInterval = 500; // 0.5 sec

WiFiClient wifi;
HttpClient client = HttpClient(wifi, server, port);

// ===== Authentication =====
const char* username = "kids_noisemeter_admin"; 
const char* passwordAuth = "passwordkids";

// ===== Variables for 1-hour average =====
float runningAverage = 0;
unsigned long readingCount = 0;
unsigned long lastHourTime = 0;
const unsigned long oneHourMillis = 3600000; // 1 hour in ms or 60000 if needed to check hourlyAverage faster

// ===== LED/Buzzer spike logic =====
unsigned long ledOnUntil = 0;

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
}

void loop() {
  // === Read & smooth sound sensor ===
  const int numSamples = 20;
  int total = 0;
  for (int i = 0; i < numSamples; i++) {
    total += analogRead(soundSensorPin);
    delay(10);
  }

  int avgRaw = total / numSamples;
  int soundLevel = map(avgRaw, 0, 1023, 0, 150);
  bool isAlert = soundLevel > threshold;

  // === LED/Buzzer logic ===
  if (isAlert) {
    if (millis() > ledOnUntil) {
      ledOnUntil = millis() + 300; // Minimum ON for short spike
    }
  }
  digitalWrite(ledPin, (millis() < ledOnUntil) ? HIGH : LOW);
  digitalWrite(buzzerPin, (millis() < ledOnUntil) ? HIGH : LOW);

  // === Update 1-hour running average ===
  readingCount++;
  runningAverage = ((runningAverage * (readingCount - 1)) + soundLevel) / readingCount;

  // === Send real-time data to server ===
  if (WiFi.status() == WL_CONNECTED) {
    String jsonData = "{";
    jsonData += "\"device_id\":\"" + String(deviceID) + "\",";
    jsonData += "\"room_name\":\"" + String(roomName) + "\",";
    jsonData += "\"current_sound_level\":" + String(soundLevel) + ",";
    jsonData += "\"average_sound_level\":" + String(runningAverage, 2) + ",";
    jsonData += "\"threshold\":" + String(threshold) + ",";
    jsonData += "\"is_alert\":" + String(isAlert ? "true" : "false") + ",";
    jsonData += "\"description\":\"Real-time + 1-hour average sound reading\"";
    jsonData += "}";

    client.beginRequest();
    client.post("/api/data"); // Real-time endpoint
    String auth = String(username) + ":" + String(passwordAuth);
    client.sendHeader("Authorization", "Basic " + base64Encode(auth));
    client.sendHeader("Content-Type", "application/json");
    client.sendHeader("Content-Length", jsonData.length());
    client.beginBody();
    client.print(jsonData);
    client.endRequest();

    int statusCode = client.responseStatusCode();
    String response = client.responseBody();
    Serial.print("Real-time Status code: ");
    Serial.println(statusCode);
    Serial.print("Response: ");
    Serial.println(response);
  }

  // === Send 1-hour average if time passed ===
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
  }

  delay(checkInterval);
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

