#pragma once

#include <msdfgen.h>

#include <string>

namespace fontcatalog {

class font_holder {
  msdfgen::FreetypeHandle *ft;
  msdfgen::FontHandle *font;
  const char *fontFilename;

public:
  font_holder();
  ~font_holder();

  bool load(const char *fontFilename);
  bool load(const unsigned char *data, long size);

  std::string get_font_name() const;

  operator msdfgen::FontHandle *() const { return font; }
};

} // namespace fontcatalog