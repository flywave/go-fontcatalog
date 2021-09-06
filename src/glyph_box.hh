#pragma once

namespace fontcatalog {

struct glyph_box {
  int index;
  double advance;
  struct {
    double l, b, r, t;
  } bounds;
  struct {
    int x, y, w, h;
  } rect;
};

} // namespace fontcatalog
