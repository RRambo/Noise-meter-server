#include <WiFiNINA.h>
#include <ArduinoHttpClient.h>

const char* ssid = "Darshan 1";  // Used WiFi network name       
const char* password = "Badal402";   // Used WiFi password 

char server[] = "192.168.1.105"; // IP address of the machine running API server       
int port = 8080; // Port API server listens on

const int ledPin = 13;
const int buzzerPin = 8;
const int soundSensorPin = A0;

const int threshold = 100;   // Threshold in dB (0–150)
const unsigned long checkInterval = 500;

WiFiClient wifi;
HttpClient client = HttpClient(wifi, server, port);

const char* username = "kids_noisemeter_admin"; 
const char* passwordAuth = "passwordkids";

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
}

void loop() {
  int rawValue = analogRead(soundSensorPin);
  int soundLevel = map(rawValue, 0, 1023, 0, 150);
  bool isAlert = soundLevel > threshold;

  digitalWrite(ledPin, isAlert ? HIGH : LOW);
  digitalWrite(buzzerPin, isAlert ? HIGH : LOW);

  String jsonData = "{";
jsonData += "\"device_id\":\"arduino_001\",";
jsonData += "\"room_name\":\"PlayRoom_A\",";
jsonData += "\"sound_level\":" + String(soundLevel) + ",";
jsonData += "\"threshold\":" + String(threshold) + ",";
jsonData += "\"is_alert\":" + String(isAlert ? "true" : "false") + ",";
jsonData += "\"description\":\"Kindergarten noise measurement\"";
jsonData += "}";


  if (WiFi.status() == WL_CONNECTED) {
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

    Serial.print("Status code: ");
    Serial.println(statusCode);
    Serial.print("Response: ");
    Serial.println(response);
  } else {
    Serial.println("WiFi disconnected. Reconnecting...");
    WiFi.disconnect();
    WiFi.begin(ssid, password);
  }

  delay(checkInterval);
}

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
