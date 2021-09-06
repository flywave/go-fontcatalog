#include "msdfgen_lib.h"

#include <algorithm>
#include <fstream>
#include <iostream>
#include <string>
#include <vector>

int main(int argc, char **argv) {
  wrap_initialize_freetype();

  std::ifstream is("./NotoSans-Regular.ttf", std::ifstream::binary);
  unsigned char *buffer = nullptr;
  int length = 0;

  if (is) {
    is.seekg(0, is.end);
    length = is.tellg();
    is.seekg(0, is.beg);

    buffer = new unsigned char[length];

    is.read((char *)buffer, length);

    is.close();
  }

  struct _font_metrics_t fontinfo = msdfgen_get_font_info(buffer, length);

  font_handle_t *font = msdfgen_load_font_memory(buffer, length, 42, nullptr);

  struct _glyph_metrics_t glyphinfo;

  _Bool ok = msdfgen_get_glyph_metrics(font, 'A', 42, &glyphinfo);

  if (!ok) {
  }

  if (buffer) {
    delete[] buffer;
  }
  wrap_deinitialize_freetype();
}