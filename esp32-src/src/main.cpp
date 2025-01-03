#include <Arduino.h>
#include <WiFi.h>
#include <HTTPClient.h>

const char ssid[] = "home";
const char password[] = "7andreanoS";
const char endpoint[] = "https://10.0.0.174/get";
// const char endpoint[] = "https://www.google.com";

const uint64_t sleep_time = 10 * 1000000;

void wifi_connect() {
  WiFi.begin(ssid, password);

  while(WiFi.status() != WL_CONNECTED) {
    delay(1000);
    Serial.print(".");
  }
  Serial.println("Wi-Fi connected");
}

void query_endpoint() {
  if (WiFi.status() == WL_CONNECTED) {
    HTTPClient http;
    Serial.print("Connecting to endpoint: ");
    Serial.println(endpoint);

    http.begin(endpoint);
    http.addHeader("Accept", "image/svg+xml");
    int httpResponseCode = http.GET();

    if (httpResponseCode == HTTP_CODE_OK) {
      String svgData = http.getString();
      Serial.println(svgData);
    } else {
      Serial.print("Failed to get SVG");
      Serial.println(httpResponseCode);
    }
    http.end();
  } else {
      Serial.println("WiFi not connected");
  }
}

void setup() {
  Serial.begin(9600);
  wifi_connect();
  query_endpoint();

  WiFi.disconnect(true);
  WiFi.mode(WIFI_OFF);

  esp_sleep_enable_timer_wakeup(sleep_time);
  esp_deep_sleep_start();
}

void loop() { }