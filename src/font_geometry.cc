#include "font_geometry.hh"

namespace fontcatalog {

font_geometry::glyph_range::glyph_range()
    : glyphs(), rangeStart(0), rangeEnd(0) {}

font_geometry::glyph_range::glyph_range(
    const std::vector<std::shared_ptr<glyph_geometry>> *glyphs,
    size_t rangeStart, size_t rangeEnd)
    : glyphs(glyphs), rangeStart(rangeStart), rangeEnd(rangeEnd) {}

size_t font_geometry::glyph_range::size() const { return glyphs->size(); }

bool font_geometry::glyph_range::empty() const { return glyphs->empty(); }

const std::shared_ptr<glyph_geometry> *
font_geometry::glyph_range::begin() const {
  return glyphs->data() + rangeStart;
}

const std::shared_ptr<glyph_geometry> *font_geometry::glyph_range::end() const {
  return glyphs->data() + rangeEnd;
}

font_geometry::font_geometry()
    : geometryScale(1), metrics(),
      preferredIdentifierType(glyph_identifier_type::UNICODE_CODEPOINT),
      glyphs(&ownGlyphs), rangeStart(glyphs->size()), rangeEnd(glyphs->size()) {
}

font_geometry::font_geometry(
    std::vector<std::shared_ptr<glyph_geometry>> *glyphStorage)
    : geometryScale(1), metrics(),
      preferredIdentifierType(glyph_identifier_type::UNICODE_CODEPOINT),
      glyphs(glyphStorage), rangeStart(glyphs->size()),
      rangeEnd(glyphs->size()) {}

int font_geometry::load_glyphset(msdfgen::FontHandle *font, double fontScale,
                                 const charset &glyphset, bool enableKerning) {
  if (!(glyphs->size() == rangeEnd && load_metrics(font, fontScale)))
    return -1;
  glyphs->reserve(glyphs->size() + glyphset.size());
  int loaded = 0;
  for (unicode_t index : glyphset) {
    std::shared_ptr<glyph_geometry> glyph = std::make_shared<glyph_geometry>();
    if (glyph->load(font, geometryScale, msdfgen::GlyphIndex(index))) {
      add_glyph(glyph);
      ++loaded;
    }
  }
  if (enableKerning)
    load_kerning(font);
  preferredIdentifierType = glyph_identifier_type::GLYPH_INDEX;
  return loaded;
}

int font_geometry::load_charset(msdfgen::FontHandle *font, double fontScale,
                                const charset &charset, bool enableKerning) {
  if (!(glyphs->size() == rangeEnd && load_metrics(font, fontScale)))
    return -1;
  glyphs->reserve(glyphs->size() + charset.size());
  int loaded = 0;
  if (charset.empty()) {
    return -2;
  }
  for (unicode_t cp : charset) {
    std::shared_ptr<glyph_geometry> glyph = std::make_shared<glyph_geometry>();
    if (glyph->load(font, geometryScale, cp)) {
      add_glyph(glyph);
      ++loaded;
    }
  }
  if (!charset.empty() && loaded == 0) {
    return -charset.size();
  }
  if (enableKerning)
    load_kerning(font);
  preferredIdentifierType = glyph_identifier_type::UNICODE_CODEPOINT;
  return loaded;
}

bool font_geometry::load_metrics(msdfgen::FontHandle *font, double fontScale) {
  if (!msdfgen::getFontMetrics(metrics, font))
    return false;
  if (metrics.emSize <= 0)
    metrics.emSize = MSDF_ATLAS_DEFAULT_EM_SIZE;
  geometryScale = fontScale / metrics.emSize;
  metrics.emSize *= geometryScale;
  metrics.ascenderY *= geometryScale;
  metrics.descenderY *= geometryScale;
  metrics.lineHeight *= geometryScale;
  metrics.underlineY *= geometryScale;
  metrics.underlineThickness *= geometryScale;
  return true;
}

bool font_geometry::add_glyph(std::shared_ptr<glyph_geometry> glyph) {
  if (glyphs->size() != rangeEnd)
    return false;
  glyphsByIndex.insert(std::make_pair(glyph->get_index(), rangeEnd));
  if (glyph->get_codepoint())
    glyphsByCodepoint.insert(std::make_pair(glyph->get_codepoint(), rangeEnd));
  glyphs->push_back(glyph);
  ++rangeEnd;
  return true;
}

int font_geometry::load_kerning(msdfgen::FontHandle *font) {
  int loaded = 0;
  for (size_t i = rangeStart; i < rangeEnd; ++i)
    for (size_t j = rangeStart; j < rangeEnd; ++j) {
      double advance;
      if (msdfgen::getKerning(advance, font, (*glyphs)[i]->get_glyph_index(),
                              (*glyphs)[j]->get_glyph_index()) &&
          advance) {
        kerning[std::make_pair<int, int>((*glyphs)[i]->get_index(),
                                         (*glyphs)[j]->get_index())] =
            geometryScale * advance;
        ++loaded;
      }
    }
  return loaded;
}

void font_geometry::set_name(const char *name) {
  if (name)
    this->name = name;
  else
    this->name.clear();
}

double font_geometry::get_geometry_scale() const { return geometryScale; }

const msdfgen::FontMetrics &font_geometry::get_metrics() const {
  return metrics;
}

glyph_identifier_type font_geometry::get_preferred_identifier_type() const {
  return preferredIdentifierType;
}

font_geometry::glyph_range font_geometry::get_glyphs() const {
  return glyph_range(glyphs, rangeStart, rangeEnd);
}

std::shared_ptr<glyph_geometry>
font_geometry::get_glyph(msdfgen::GlyphIndex index) const {
  std::map<int, size_t>::const_iterator it =
      glyphsByIndex.find(index.getIndex());
  if (it != glyphsByIndex.end())
    return (*glyphs)[it->second];
  return nullptr;
}

std::shared_ptr<glyph_geometry>
font_geometry::get_glyph(unicode_t codepoint) const {
  std::map<unicode_t, size_t>::const_iterator it =
      glyphsByCodepoint.find(codepoint);
  if (it != glyphsByCodepoint.end())
    return (*glyphs)[it->second];
  return nullptr;
}

bool font_geometry::get_advance(double &advance, msdfgen::GlyphIndex index1,
                                msdfgen::GlyphIndex index2) const {
  std::shared_ptr<glyph_geometry> glyph1 = get_glyph(index1);
  if (!glyph1)
    return false;
  advance = glyph1->get_advance();
  std::map<std::pair<int, int>, double>::const_iterator it = kerning.find(
      std::make_pair<int, int>(index1.getIndex(), index2.getIndex()));
  if (it != kerning.end())
    advance += it->second;
  return true;
}

bool font_geometry::get_advance(double &advance, unicode_t codepoint1,
                                unicode_t codepoint2) const {
  std::shared_ptr<glyph_geometry> glyph1;
  std::shared_ptr<glyph_geometry> glyph2;
  if (!((glyph1 = get_glyph(codepoint1)) && (glyph2 = get_glyph(codepoint2))))
    return false;
  advance = glyph1->get_advance();
  std::map<std::pair<int, int>, double>::const_iterator it = kerning.find(
      std::make_pair<int, int>(glyph1->get_index(), glyph2->get_index()));
  if (it != kerning.end())
    advance += it->second;
  return true;
}

const std::map<std::pair<int, int>, double> &
font_geometry::get_kerning() const {
  return kerning;
}

const char *font_geometry::get_name() const {
  if (name.empty())
    return nullptr;
  return name.c_str();
}

} // namespace fontcatalog
