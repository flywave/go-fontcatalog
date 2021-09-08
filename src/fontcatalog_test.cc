#include "fontcatalog_lib.h"

#include "font_holder.hh"
#include "font_geometry.hh"

#include <algorithm>
#include <fstream>
#include <iostream>
#include <string>
#include <vector>

int main(int argc, char **argv) {
  std::ifstream is("./fonts/FiraGO_Map.ttf", std::ifstream::binary);
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

  fc_font_holder_t *fh = fc_font_holder_load_font_memory(buffer, length);

  fc_glyph_geometry_list_t *glyphs = fc_new_glyph_geometry_list();

  fc_font_geometry_t *geom = fc_new_font_geometry_with_glyphs(glyphs);

  fc_charset_t *ascii = fc_new_charset();

  fc_charset_add(ascii, 160);

  int n = fc_font_geometry_load_from_charset(geom, fh, 41, ascii);

  fc_glyph_geometry_t *ggeom = fc_font_geometry_get_glyph_from_unicode(geom, 160);

  int id = fc_glyph_geometry_get_index(ggeom);

  bool ret = fc_glyph_geometry_is_whitespace(ggeom);

  fc_glyph_geometry_wrap_box(ggeom, -1, 2.0, 1.0);

  fc_generator_attributes_t *arrr = fc_new_generator_attributes();
  

  fc_bitmap_t *bmap = fc_new_bitmap_alloc(4, 62, 62);

  fc_mtsdf_generator(bmap, ggeom, arrr);

  unsigned char * data = fc_bitmap_blit_data(bmap);

  free(data);
  fc_bitmap_free(bmap);
  fc_generator_attributes_free(arrr);
  fc_glyph_geometry_free(ggeom);
  fc_font_geometry_free(geom);
  fc_glyph_geometry_list_free(glyphs);
  fc_font_holder_free(fh);

  if (buffer) {
    delete[] buffer;
  }
}