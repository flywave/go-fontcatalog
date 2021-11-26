#pragma once

#include "charset.hh"
#include "glyph_geometry.hh"
#include "types.hh"

#include <map>
#include <memory>
#include <msdfgen-ext.h>
#include <msdfgen.h>
#include <string>
#include <utility>
#include <vector>

#define MSDF_ATLAS_DEFAULT_EM_SIZE 32.0

namespace fontcatalog {

class font_geometry {
public:
  class glyph_range {
  public:
    glyph_range();
    glyph_range(const std::vector<std::shared_ptr<glyph_geometry>> *glyphs,
                size_t rangeStart, size_t rangeEnd);
    size_t size() const;
    bool empty() const;
    const std::shared_ptr<glyph_geometry> *begin() const;
    const std::shared_ptr<glyph_geometry> *end() const;

  private:
    const std::vector<std::shared_ptr<glyph_geometry>> *glyphs;
    size_t rangeStart, rangeEnd;
  };

  font_geometry();
  explicit font_geometry(
      std::vector<std::shared_ptr<glyph_geometry>> *glyphStorage);

  int load_glyphset(msdfgen::FontHandle *font, double fontScale,
                    const charset &glyphset, bool enableKerning = true);
  int load_charset(msdfgen::FontHandle *font, double fontScale,
                   const charset &charset, bool enableKerning = true);

  bool load_metrics(msdfgen::FontHandle *font, double fontScale);

  bool add_glyph(std::shared_ptr<glyph_geometry> glyph);

  int load_kerning(msdfgen::FontHandle *font);

  void set_name(const char *name);

  double get_geometry_scale() const;
  const msdfgen::FontMetrics &get_metrics() const;
  glyph_identifier_type get_preferred_identifier_type() const;
  glyph_range get_glyphs() const;

  std::shared_ptr<glyph_geometry> get_glyph(msdfgen::GlyphIndex index) const;
  std::shared_ptr<glyph_geometry> get_glyph(unicode_t codepoint) const;

  bool get_advance(double &advance, msdfgen::GlyphIndex index1,
                   msdfgen::GlyphIndex index2) const;
  bool get_advance(double &advance, unicode_t codepoint1,
                   unicode_t codepoint2) const;

  const std::map<std::pair<int, int>, double> &get_kerning() const;
  const char *get_name() const;

private:
  double geometryScale;
  msdfgen::FontMetrics metrics;
  glyph_identifier_type preferredIdentifierType;
  std::vector<std::shared_ptr<glyph_geometry>> *glyphs;
  size_t rangeStart, rangeEnd;
  std::map<int, size_t> glyphsByIndex;
  std::map<unicode_t, size_t> glyphsByCodepoint;
  std::map<std::pair<int, int>, double> kerning;
  std::vector<std::shared_ptr<glyph_geometry>> ownGlyphs;
  std::string name;
};

} // namespace fontcatalog
