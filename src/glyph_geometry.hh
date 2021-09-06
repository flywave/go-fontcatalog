
#pragma once

#include "glyph_box.hh"
#include "types.hh"
#include <msdfgen-ext.h>
#include <msdfgen.h>

namespace fontcatalog {

class glyph_geometry {
public:
  glyph_geometry();

  bool load(msdfgen::FontHandle *font, double geometryScale,
            msdfgen::GlyphIndex index, bool preprocessGeometry = true);
  bool load(msdfgen::FontHandle *font, double geometryScale,
            unicode_t codepoint, bool preprocessGeometry = true);

  void edge_coloring(void (*fn)(msdfgen::Shape &, double, unsigned long long),
                     double angleThreshold, unsigned long long seed);

  void wrap_box(double scale, double range, double miterLimit);

  void place_box(int x, int y);

  int get_index() const;

  msdfgen::GlyphIndex get_glyph_index() const;

  unicode_t get_codepoint() const;

  int get_identifier(glyph_identifier_type type) const;

  const msdfgen::Shape &get_shape() const;

  double get_advance() const;

  void get_box_rect(int &x, int &y, int &w, int &h) const;

  void get_box_size(int &w, int &h) const;

  double get_box_range() const;

  msdfgen::Projection get_box_projection() const;

  double get_box_scale() const;

  msdfgen::Vector2 get_box_translate() const;

  void get_quad_plane_bounds(double &l, double &b, double &r, double &t) const;

  void get_quad_atlas_bounds(double &l, double &b, double &r, double &t) const;

  bool is_whitespace() const;

  operator glyph_box() const;

private:
  int index;
  unicode_t codepoint;
  double geometryScale;
  msdfgen::Shape shape;
  msdfgen::Shape::Bounds bounds;
  double advance;
  struct {
    struct {
      int x, y, w, h;
    } rect;
    double range;
    double scale;
    msdfgen::Vector2 translate;
  } box;
};

} // namespace fontcatalog
