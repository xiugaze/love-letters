#ifndef MAIN_H
#define MAIN_H
#include <stdint.h>

void display_stream_HTTP(const char* url, uint16_t width, uint16_t height);
void text_centered(char text[]);
void wifi_init();
void display_init();

#endif