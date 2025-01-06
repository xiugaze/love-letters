#include <Arduino.h>
#include <GxEPD2_BW.h>
#include <Fonts/FreeMonoBold12pt7b.h>
#include <WiFi.h>
#include <HTTPClient.h>

#include <main.hpp>


#define CS 15
#define DC 1
#define RST 26
#define BUSY 4
#define CLK 13
#define MISO 12
#define MOSI 14
#define CS 15

const char* ssid     = "home";
const char* password = "7andreanoS";
const char* HOST = "https://10.0.0.174/get-bin";

/* From GxEPD2 example code*/
#define MAX_DISPLAY_BUFFER_SIZE 65536ul
#define MAX_HEIGHT(EPD) (EPD::HEIGHT <= MAX_DISPLAY_BUFFER_SIZE / (EPD::WIDTH / 8) ? EPD::HEIGHT : MAX_DISPLAY_BUFFER_SIZE / (EPD::WIDTH / 8))
GxEPD2_BW<GxEPD2_750_T7, MAX_HEIGHT(GxEPD2_750_T7)> display(GxEPD2_750_T7(/* CS */ CS, /* DC */ DC, /* RST */ RST, /* BUSY */ BUSY));
SPIClass hspi(HSPI);

void setup() {
  Serial.begin(9600);
  while(!Serial) {};
  delay(3000);

  wifi_init();
  display_init();

  display_stream_HTTP(HOST, 480, 800);
}

void loop() { }

void text_centered(char text[]) {
  //Serial.println("helloWorld");
  display.setRotation(1);
  display.setFont(&FreeMonoBold12pt7b);
  display.setTextColor(GxEPD_BLACK);
  int16_t tbx, tby; uint16_t tbw, tbh;
  display.getTextBounds(text, 0, 0, &tbx, &tby, &tbw, &tbh);
  // center bounding box by transposition of origin:
  uint16_t x = ((display.width() - tbw) / 2) - tbx;
  uint16_t y = ((display.height() - tbh) / 2) - tby;
  display.setFullWindow();
  display.firstPage();
  do {
    display.fillScreen(GxEPD_WHITE);
    display.setCursor(x, y);
    display.print(text);
  }
  while (display.nextPage());
}

void wifi_init() {
  WiFi.begin(ssid, password);
  while (WiFi.status() != WL_CONNECTED) {
      Serial.print(".");
      delay(500);
  }
  Serial.println("wifi Connected");
}

void display_init() {
  hspi.begin(CLK, MISO, MOSI, CS); 
  display.epd2.selectSPI(hspi, SPISettings(4000000, MSBFIRST, SPI_MODE0));
  display.init(115200);
  display.setRotation(1); // portrait
  Serial.println("display initialized");
}

void display_stream_HTTP(const char* url, uint16_t width, uint16_t height)
{
  HTTPClient http;
  if (!http.begin(url)) {
    Serial.println("http.begin() failed, check url (did you forget https://?)");
    return;
  }

  int httpCode = http.GET();
  if (httpCode != HTTP_CODE_OK) {
    Serial.printf("HTTP GET failed, code: %d\n", httpCode);
    http.end();
    return;
  }

  WiFiClient& stream = http.getStream();
  Serial.printf("HTTP size: %d (expected 348,000)\n", http.getSize());

  display.firstPage();
  do {
    static uint8_t lineBuf[480*2];
    for (uint16_t row = 0; row < height; row++) {
      int bytesRead = stream.readBytes((char*)lineBuf, width);
      if (bytesRead < width) {
        Serial.println("stream ended unexpectedly or error reading line.");
        break;
      }
      for (uint16_t col = 0; col < width; col++) {
        display.drawPixel(col, row, lineBuf[col]);
      }
    }
  } while (display.nextPage());

  http.end();
}