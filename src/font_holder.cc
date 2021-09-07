#include "font_holder.hh"
#include <ft2build.h>
#include FT_FREETYPE_H
#include FT_OUTLINE_H
#include FT_TRUETYPE_TABLES_H
#include FT_SFNT_NAMES_H
#include FT_BITMAP_H
#include FT_IMAGE_H
#include "msdfgen-ext.h"

namespace fontcatalog {

font_holder::font_holder()
    : ft(msdfgen::initializeFreetype()), font(nullptr), fontFilename(nullptr) {}

font_holder::~font_holder() {
  if (ft) {
    if (font)
      msdfgen::destroyFont(font);
    msdfgen::deinitializeFreetype(ft);
  }
}

bool font_holder::load(const char *fontFilename) {
  if (ft && fontFilename) {
    if (this->fontFilename && !strcmp(this->fontFilename, fontFilename))
      return true;
    if (font)
      msdfgen::destroyFont(font);
    if ((font = msdfgen::loadFont(ft, fontFilename))) {
      this->fontFilename = fontFilename;
      return true;
    }
    this->fontFilename = nullptr;
  }
  return false;
}

bool font_holder::load(const unsigned char *data, long size) {
  if (ft && data) {
    if (font)
      msdfgen::destroyFont(font);
    if ((font = msdfgen::loadFontData(ft, data, size))) {
      return true;
    }
  }
  return false;
}

} // namespace fontcatalog