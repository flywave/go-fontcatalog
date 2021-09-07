#include <cstdlib>
#include <ft2build.h>
#include <iostream>
#include <lodepng.h>
#include <memory>
#include <stdio.h>
#include FT_FREETYPE_H
#include FT_OUTLINE_H
#include FT_TRUETYPE_TABLES_H
#include FT_SFNT_NAMES_H
#include FT_BITMAP_H
#include FT_IMAGE_H

#include "font_geometry.hh"
#include "font_holder.hh"
#include "fontcatalog_lib.h"
#include "generator_attributes.hh"
#include "glyph_generators.hh"
#include "glyph_geometry.hh"

#include "bitmap_blit.hh"
#include "msdfgen-ext.h"
#include "msdfgen.h"

#ifdef __cplusplus
extern "C" {
#endif

using namespace msdfgen;

struct _fc_font_holder_t {
  fontcatalog::font_holder h;
};

struct _fc_glyph_geometry_t {
  std::shared_ptr<fontcatalog::glyph_geometry> g;
};

struct _fc_glyph_geometry_list_t {
  std::vector<std::shared_ptr<fontcatalog::glyph_geometry>> gs;
};

struct _fc_font_geometry_t {
  fontcatalog::font_geometry g;
};

struct _fc_font_geometry_list_t {
  std::vector<fontcatalog::font_geometry> gs;
};

struct _fc_glyph_range_t {
  fontcatalog::font_geometry::glyph_range it;
};

struct _fc_generator_attributes_t {
  fontcatalog::generator_attributes attr;
};

struct _fc_bitmap_t {
  msdfgen::Bitmap<float, 1> bitmap1;
  msdfgen::Bitmap<float, 3> bitmap3;
  msdfgen::Bitmap<float, 4> bitmap4;
  int N;
};

struct _fc_bitmap_ref_t {
  msdfgen::BitmapRef<float, 1> bitmap1;
  msdfgen::BitmapRef<float, 3> bitmap3;
  msdfgen::BitmapRef<float, 4> bitmap4;
  int N;
};

struct _fc_charset_t {
  fontcatalog::charset c;
};

struct _fc_kerning_map_t {
  std::map<std::pair<int, int>, double> ks;
};

FC_LIB_EXPORT fc_font_holder_t *
fc_font_holder_load_font_memory(const unsigned char *data, long size) {
  fc_font_holder_t *holder = new fc_font_holder_t{};
  holder->h.load(data, size);
  return holder;
}

FC_LIB_EXPORT void fc_font_holder_free(fc_font_holder_t *handle) {
  delete handle;
}

FC_LIB_EXPORT struct _fc_font_info_t
fc_font_holder_get_font_info(fc_font_holder_t *handle) {
  FT_Face ft = msdfgen::getFreetypeFont(handle->h);

  struct _fc_font_info_t metrics;
  metrics.ascent = ft->ascender;
  metrics.descent = ft->descender;
  metrics.unitsPerEm = ft->units_per_EM;
  metrics.baseLine = (ft->size->metrics.ascender + 32) >> 6;
  metrics.lineHeight = (ft->size->metrics.height + 32) >> 6;
  TT_Header *header = (TT_Header *)FT_Get_Sfnt_Table(ft, FT_SFNT_HEAD);
  metrics.flags = (int)header->Mac_Style | header->Flags << 16;

  FT_ULong charcode;
  FT_UInt gindex;
  std::vector<int> charmap;
  charcode = FT_Get_First_Char(ft, &gindex);
  charmap.emplace_back(charcode);
  while (gindex != 0) {
    charcode = FT_Get_Next_Char(ft, charcode, &gindex);
    charmap.emplace_back(charcode);
  }
  int *data = (int *)malloc(sizeof(int) * charmap.size());
  metrics.charSize = charmap.size();
  memcpy(data, charmap.data(), sizeof(int) * charmap.size());
  metrics.characterSet = data;
  return metrics;
}

FC_LIB_EXPORT fc_glyph_geometry_t *fc_new_glyph_geometry_from_glyph_index(
    fc_font_holder_t *handle, double geometryScale, fc_glyph_index_t index) {
  fc_glyph_geometry_t *holder =
      new fc_glyph_geometry_t{std::make_shared<fontcatalog::glyph_geometry>()};
  if (holder->g->load(handle->h, geometryScale, index, false)) {
    return holder;
  }
  delete holder;
  return nullptr;
}

FC_LIB_EXPORT fc_glyph_geometry_t *fc_new_glyph_geometry_from_unicode(
    fc_font_holder_t *handle, double geometryScale, fc_unicode_t codepoint) {
  fc_glyph_geometry_t *holder =
      new fc_glyph_geometry_t{std::make_shared<fontcatalog::glyph_geometry>()};
  if (holder->g->load(handle->h, geometryScale, codepoint, false)) {
    return holder;
  }
  delete holder;
  return nullptr;
}

FC_LIB_EXPORT void fc_glyph_geometry_free(fc_glyph_geometry_t *geom) {
  delete geom;
}

enum EdgeColoring {
  EdgeColoringSimple = 0,
  EdgeColoringInkTrap = 1,
  EdgeColoringByDistance = 2,
};

FC_LIB_EXPORT void fc_glyph_geometry_edge_coloring(fc_glyph_geometry_t *geom,
                                                   uint32_t type,
                                                   double angleThreshold,
                                                   unsigned long long seed) {
  switch (type) {
  case EdgeColoringSimple:
    geom->g->edge_coloring(msdfgen::edgeColoringSimple, angleThreshold, seed);
    break;
  case EdgeColoringInkTrap:
    geom->g->edge_coloring(msdfgen::edgeColoringInkTrap, angleThreshold, seed);
    break;
  case EdgeColoringByDistance:
    geom->g->edge_coloring(msdfgen::edgeColoringByDistance, angleThreshold,
                           seed);
    break;
  default:
    break;
  }
}

FC_LIB_EXPORT void fc_glyph_geometry_wrap_box(fc_glyph_geometry_t *geom,
                                              double scale, double range,
                                              double miterLimit) {
  geom->g->wrap_box(scale, range, miterLimit);
}

FC_LIB_EXPORT void fc_glyph_geometry_place_box(fc_glyph_geometry_t *geom, int x,
                                               int y) {
  geom->g->place_box(x, y);
}

FC_LIB_EXPORT int fc_glyph_geometry_get_index(fc_glyph_geometry_t *geom) {
  return geom->g->get_index();
}

FC_LIB_EXPORT fc_glyph_index_t
fc_glyph_geometry_get_glyph_index(fc_glyph_geometry_t *geom) {
  return geom->g->get_glyph_index().getIndex();
}

FC_LIB_EXPORT fc_unicode_t
fc_glyph_geometry_get_codepoint(fc_glyph_geometry_t *geom) {
  return geom->g->get_codepoint();
}

FC_LIB_EXPORT int fc_glyph_geometry_get_identifier(fc_glyph_geometry_t *geom,
                                                   uint32_t type) {
  return geom->g->get_identifier(
      static_cast<fontcatalog::glyph_identifier_type>(type));
}

FC_LIB_EXPORT double fc_glyph_geometry_get_advance(fc_glyph_geometry_t *geom) {
  return geom->g->get_advance();
}

FC_LIB_EXPORT void fc_glyph_geometry_get_box_rect(fc_glyph_geometry_t *geom,
                                                  int *x, int *y, int *w,
                                                  int *h) {
  return geom->g->get_box_rect(*x, *y, *w, *h);
}

FC_LIB_EXPORT void fc_glyph_geometry_get_box_size(fc_glyph_geometry_t *geom,
                                                  int *w, int *h) {
  return geom->g->get_box_size(*w, *h);
}

FC_LIB_EXPORT double
fc_glyph_geometry_get_box_range(fc_glyph_geometry_t *geom) {
  return geom->g->get_box_range();
}

FC_LIB_EXPORT double
fc_glyph_geometry_get_box_scale(fc_glyph_geometry_t *geom) {
  return geom->g->get_box_scale();
}

FC_LIB_EXPORT void
fc_glyph_geometry_get_box_translate(fc_glyph_geometry_t *geom, int *tx,
                                    int *ty) {
  auto tran = geom->g->get_box_translate();
  *tx = tran.x;
  *ty = tran.y;
}

FC_LIB_EXPORT fc_glyph_box_t
fc_glyph_geometry_get_glyph_box(fc_glyph_geometry_t *geom) {
  fontcatalog::glyph_box gbox = (fontcatalog::glyph_box)(*geom->g);
  return *reinterpret_cast<fc_glyph_box_t *>(&gbox);
}

FC_LIB_EXPORT _Bool fc_glyph_geometry_is_whitespace(fc_glyph_geometry_t *geom) {
  return geom->g->is_whitespace();
}

FC_LIB_EXPORT fc_glyph_geometry_list_t *fc_new_glyph_geometry_list() {
  return new fc_glyph_geometry_list_t{};
}

FC_LIB_EXPORT void fc_glyph_geometry_list_free(fc_glyph_geometry_list_t *list) {
  delete list;
}

FC_LIB_EXPORT void
fc_glyph_geometry_list_push_geometry(fc_glyph_geometry_list_t *list,
                                     fc_glyph_geometry_t *geom) {
  list->gs.emplace_back(geom->g);
}

FC_LIB_EXPORT _Bool
fc_glyph_geometry_list_empty(fc_glyph_geometry_list_t *list) {
  return list->gs.empty();
}

FC_LIB_EXPORT size_t
fc_glyph_geometry_list_size(fc_glyph_geometry_list_t *list) {
  return list->gs.size();
}

enum ErrorCorrection {
  DISABLED = 0,
  INDISCRIMINATE = 1,
  EDGE_PRIORITY = 2,
  EDGE_ONLY = 3,
};

enum DistanceCheckMode {
  DO_NOT_CHECK_DISTANCE = 0,
  CHECK_DISTANCE_AT_EDGE = 1,
  ALWAYS_CHECK_DISTANCE = 2,
};

FC_LIB_EXPORT fc_font_geometry_t *
fc_new_font_geometry_with_glyphs(fc_glyph_geometry_list_t *glyphs) {
  return new fc_font_geometry_t{fontcatalog::font_geometry{&glyphs->gs}};
}

FC_LIB_EXPORT void fc_font_geometry_free(fc_font_geometry_t *geom) {
  delete geom;
}

FC_LIB_EXPORT int fc_font_geometry_load_from_glyphset(fc_font_geometry_t *fonts,
                                                      fc_font_holder_t *handle,
                                                      double fontScale,
                                                      fc_charset_t *charsets) {
  return fonts->g.load_glyphset(handle->h, fontScale, charsets->c, false);
}

FC_LIB_EXPORT int fc_font_geometry_load_from_charset(fc_font_geometry_t *fonts,
                                                     fc_font_holder_t *handle,
                                                     double fontScale,
                                                     fc_charset_t *charsets) {
  return fonts->g.load_charset(handle->h, fontScale, charsets->c, false);
}

FC_LIB_EXPORT _Bool fc_font_geometry_load_metrics(fc_font_geometry_t *fonts,
                                                  fc_font_holder_t *handle,
                                                  double fontScale) {
  return fonts->g.load_metrics(handle->h, fontScale);
}

FC_LIB_EXPORT _Bool fc_font_geometry_add_glyph(fc_font_geometry_t *fonts,
                                               fc_glyph_geometry_t *geom) {
  return fonts->g.add_glyph(geom->g);
}

FC_LIB_EXPORT int fc_font_geometry_load_kerning(fc_font_geometry_t *fonts,
                                                fc_font_holder_t *handle) {
  return fonts->g.load_kerning(handle->h);
}

FC_LIB_EXPORT void fc_font_geometry_set_name(fc_font_geometry_t *geom,
                                             const char *name) {
  geom->g.set_name(name);
}

FC_LIB_EXPORT const char *fc_font_geometry_get_name(fc_font_geometry_t *geom) {
  return geom->g.get_name();
}

FC_LIB_EXPORT double
fc_font_geometry_geometry_scale(fc_font_geometry_t *fonts) {
  return fonts->g.get_geometry_scale();
}

FC_LIB_EXPORT struct _fc_font_metrics_t
fc_font_geometry_get_metrics(fc_font_geometry_t *fonts) {
  auto metrics = fonts->g.get_metrics();
  return *reinterpret_cast<struct _fc_font_metrics_t *>(&metrics);
}

FC_LIB_EXPORT uint32_t
fc_font_geometry_get_preferred_identifier_type(fc_font_geometry_t *fonts) {
  return static_cast<uint32_t>(fonts->g.get_preferred_identifier_type());
}

FC_LIB_EXPORT fc_glyph_range_t *
fc_font_geometry_get_glyphs(fc_font_geometry_t *fonts) {
  return new fc_glyph_range_t{fonts->g.get_glyphs()};
}

FC_LIB_EXPORT fc_glyph_geometry_t *
fc_font_geometry_get_glyph_from_index(fc_font_geometry_t *fonts,
                                      fc_glyph_index_t index) {
  return new fc_glyph_geometry_t{
      fonts->g.get_glyph(msdfgen::GlyphIndex(index))};
}

FC_LIB_EXPORT fc_glyph_geometry_t *
fc_font_geometry_get_glyph_from_unicode(fc_font_geometry_t *fonts,
                                        fc_unicode_t codePoint) {
  return new fc_glyph_geometry_t{fonts->g.get_glyph(codePoint)};
}

FC_LIB_EXPORT _Bool fc_font_geometry_get_advance_from_index(
    fc_font_geometry_t *fonts, double *advance, fc_glyph_index_t index1,
    fc_glyph_index_t index2) {
  return fonts->g.get_advance(*advance, index1, index2);
}

FC_LIB_EXPORT _Bool fc_font_geometry_get_advance_from_unicode(
    fc_font_geometry_t *fonts, double *advance, fc_unicode_t codePoint1,
    fc_unicode_t codePoint2) {
  return fonts->g.get_advance(*advance, codePoint1, codePoint2);
}

FC_LIB_EXPORT fc_kerning_map_t *
fc_font_geometry_get_kerning(fc_font_geometry_t *fonts) {
  return new fc_kerning_map_t{fonts->g.get_kerning()};
}

FC_LIB_EXPORT void fc_kerning_map_free(fc_kerning_map_t *kmap) { delete kmap; }

FC_LIB_EXPORT fc_kerning_t *fc_kerning_map_get_kernings(fc_kerning_map_t *kmap,
                                                        size_t *si) {
  fc_kerning_t *ret =
      (fc_kerning_t *)malloc(sizeof(fc_kerning_t) * kmap->ks.size());
  int i = 0;
  for (auto kp : kmap->ks) {
    ret[i++] = fc_kerning_t{
      first : kp.first.first,
      second : kp.first.second,
      kerning : kp.second
    };
  }
  return ret;
}

FC_LIB_EXPORT fc_font_geometry_list_t *fc_new_font_geometry_list() {
  return new fc_font_geometry_list_t{};
}

FC_LIB_EXPORT void fc_font_geometry_list_free(fc_font_geometry_list_t *list) {
  delete list;
}

FC_LIB_EXPORT void
fc_font_geometry_list_push_geometry(fc_font_geometry_list_t *list,
                                    fc_font_geometry_t *geom) {
  list->gs.emplace_back(geom->g);
}

FC_LIB_EXPORT _Bool fc_font_geometry_list_empty(fc_font_geometry_list_t *list) {
  return list->gs.empty();
}

FC_LIB_EXPORT size_t fc_font_geometry_list_size(fc_font_geometry_list_t *list) {
  return list->gs.size();
}

FC_LIB_EXPORT void fc_glyph_range_free(fc_glyph_range_t *gr) { delete gr; }

FC_LIB_EXPORT size_t fc_glyph_range_size(fc_glyph_range_t *gr) {
  return gr->it.size();
}

FC_LIB_EXPORT _Bool fc_glyph_range_empty(fc_glyph_range_t *gr) {
  return gr->it.empty();
}

FC_LIB_EXPORT fc_glyph_geometry_t *fc_glyph_range_get(fc_glyph_range_t *gr,
                                                      size_t index) {
  return new fc_glyph_geometry_t{*(gr->it.begin() + index)};
}

FC_LIB_EXPORT fc_generator_attributes_t *fc_new_generator_attributes() {
  return new fc_generator_attributes_t{};
}

FC_LIB_EXPORT void fc_generator_attributes_free(fc_generator_attributes_t *ga) {
  delete ga;
}

FC_LIB_EXPORT void
fc_generator_attributes_set_min_deviation_ratio(fc_generator_attributes_t *ga,
                                                double ratio) {
  ga->attr.config.errorCorrection.minDeviationRatio = ratio;
}

FC_LIB_EXPORT void
fc_generator_attributes_set_min_improve_ratio(fc_generator_attributes_t *ga,
                                              double ratio) {
  ga->attr.config.errorCorrection.minImproveRatio = ratio;
}

FC_LIB_EXPORT void
fc_generator_attributes_set_mode(fc_generator_attributes_t *ga, uint32_t mode) {
  ga->attr.config.errorCorrection.mode =
      static_cast<ErrorCorrectionConfig::Mode>(mode);
}

FC_LIB_EXPORT void
fc_generator_attributes_set_distance_check_mode(fc_generator_attributes_t *ga,
                                                uint32_t mode) {
  ga->attr.config.errorCorrection.distanceCheckMode =
      static_cast<ErrorCorrectionConfig::DistanceCheckMode>(mode);
}

FC_LIB_EXPORT void
fc_generator_attributes_set_buffer(fc_generator_attributes_t *ga,
                                   unsigned char *buffer) {
  ga->attr.config.errorCorrection.buffer = buffer;
}

FC_LIB_EXPORT void
fc_generator_attributes_set_overlap_support(fc_generator_attributes_t *ga,
                                            _Bool overlapSupport) {
  ga->attr.config.overlapSupport = overlapSupport;
}

FC_LIB_EXPORT void
fc_generator_attributes_set_scanline_pass(fc_generator_attributes_t *ga,
                                          _Bool scanlinePass) {
  ga->attr.scanlinePass = scanlinePass;
}

FC_LIB_EXPORT fc_bitmap_t *fc_new_bitmap(int channel) {
  switch (channel) {
  case 1:
    return new fc_bitmap_t{N : 1};
  case 3:
    return new fc_bitmap_t{N : 3};
  case 4:
    return new fc_bitmap_t{N : 4};
  default:
    break;
  }
  return nullptr;
}

FC_LIB_EXPORT fc_bitmap_t *fc_new_bitmap_alloc(int channel, int width,
                                               int height) {
  switch (channel) {
  case 1:
    return new
    fc_bitmap_t{bitmap1 : msdfgen::Bitmap<float, 1>(width, height), N : 1};
  case 3:
    return new
    fc_bitmap_t{bitmap3 : msdfgen::Bitmap<float, 3>(width, height), N : 3};
  case 4:
    return new
    fc_bitmap_t{bitmap4 : msdfgen::Bitmap<float, 4>(width, height), N : 4};
  default:
    break;
  }
  return nullptr;
}

FC_LIB_EXPORT void fc_bitmap_free(fc_bitmap_t *gr) { delete gr; }

FC_LIB_EXPORT int fc_bitmap_width(fc_bitmap_t *gr) {
  switch (gr->N) {
  case 1:
    return gr->bitmap1.width();
  case 3:
    return gr->bitmap3.width();
  case 4:
    return gr->bitmap4.width();
  default:
    break;
  }
  return -1;
}

FC_LIB_EXPORT int fc_bitmap_height(fc_bitmap_t *gr) {
  switch (gr->N) {
  case 1:
    return gr->bitmap1.height();
  case 3:
    return gr->bitmap3.height();
  case 4:
    return gr->bitmap4.height();
  default:
    break;
  }
  return -1;
}

FC_LIB_EXPORT int fc_bitmap_channels(fc_bitmap_t *gr) { return gr->N; }

FC_LIB_EXPORT float *fc_bitmap_data(fc_bitmap_t *gr) {
  switch (gr->N) {
  case 1:
    return gr->bitmap1;
  case 3:
    return gr->bitmap3;
  case 4:
    return gr->bitmap4;
  default:
    break;
  }
  return nullptr;
}

FC_LIB_EXPORT unsigned char *fc_bitmap_blit_data(fc_bitmap_t *gr) {
  if (gr->N == 1) {
    msdfgen::Bitmap<byte, 1> dst(gr->bitmap1.width(), gr->bitmap1.height());
    msdfgen::BitmapRef<byte, 1> refdst = dst;
    fontcatalog::blit(refdst, gr->bitmap1, 0, 0, 0, 0, gr->bitmap1.width(),
                      gr->bitmap1.height());
    size_t bytesize = sizeof(byte) * gr->bitmap1.width() * gr->bitmap1.height();
    unsigned char *ret = (unsigned char *)malloc(bytesize);
    memcpy(ret, refdst.pixels, bytesize);
    return ret;
  } else if (gr->N == 3) {
    msdfgen::Bitmap<byte, 3> dst(gr->bitmap3.width(), gr->bitmap3.height());
    msdfgen::BitmapRef<byte, 3> refdst = dst;
    fontcatalog::blit(refdst, gr->bitmap3, 0, 0, 0, 0, gr->bitmap3.width(),
                      gr->bitmap3.height());
    size_t bytesize =
        sizeof(byte) * gr->bitmap3.width() * gr->bitmap3.height() * 3;
    unsigned char *ret = (unsigned char *)malloc(bytesize);
    memcpy(ret, refdst.pixels, bytesize);
    return ret;
  } else if (gr->N == 4) {
    msdfgen::Bitmap<byte, 4> dst(gr->bitmap4.width(), gr->bitmap4.height());
    msdfgen::BitmapRef<byte, 4> refdst = dst;
    fontcatalog::blit(refdst, gr->bitmap4, 0, 0, 0, 0, gr->bitmap4.width(),
                      gr->bitmap4.height());
    size_t bytesize =
        sizeof(byte) * gr->bitmap4.width() * gr->bitmap4.height() * 4;
    unsigned char *ret = (unsigned char *)malloc(bytesize);
    memcpy(ret, refdst.pixels, bytesize);
    return ret;
  }
  return nullptr;
}

FC_LIB_EXPORT fc_bitmap_ref_t *fc_new_bitmap_ref(float *pixels, int channel,
                                                 int width, int height) {
  switch (channel) {
  case 1:
    return new fc_bitmap_ref_t{
      bitmap1 : msdfgen::BitmapRef<float, 1>(pixels, width, height),
      N : 1
    };
  case 3:
    return new fc_bitmap_ref_t{
      bitmap3 : msdfgen::BitmapRef<float, 3>(pixels, width, height),
      N : 3
    };
  case 4:
    return new fc_bitmap_ref_t{
      bitmap4 : msdfgen::BitmapRef<float, 4>(pixels, width, height),
      N : 4
    };
  default:
    break;
  }
  return nullptr;
}

FC_LIB_EXPORT void fc_bitmap_ref_free(fc_bitmap_ref_t *gr) { delete gr; }

FC_LIB_EXPORT int fc_bitmap_ref_width(fc_bitmap_ref_t *gr) {
  switch (gr->N) {
  case 1:
    return gr->bitmap1.width;
  case 3:
    return gr->bitmap3.width;
  case 4:
    return gr->bitmap4.width;
  default:
    break;
  }
  return -1;
}

FC_LIB_EXPORT int fc_bitmap_ref_height(fc_bitmap_ref_t *gr) {
  switch (gr->N) {
  case 1:
    return gr->bitmap1.height;
  case 3:
    return gr->bitmap3.height;
  case 4:
    return gr->bitmap4.height;
  default:
    break;
  }
  return -1;
}

FC_LIB_EXPORT int fc_bitmap_ref_channels(fc_bitmap_ref_t *gr) { return gr->N; }

FC_LIB_EXPORT float *fc_bitmap_ref_data(fc_bitmap_ref_t *gr) {
  switch (gr->N) {
  case 1:
    return gr->bitmap1.pixels;
  case 3:
    return gr->bitmap3.pixels;
  case 4:
    return gr->bitmap4.pixels;
  default:
    break;
  }
  return nullptr;
}

FC_LIB_EXPORT unsigned char *fc_bitmap_ref_blit_data(fc_bitmap_ref_t *gr) {
  if (gr->N == 1) {
    msdfgen::Bitmap<byte, 1> dst(gr->bitmap1.width, gr->bitmap1.height);
    msdfgen::BitmapRef<byte, 1> refdst = dst;
    fontcatalog::blit(refdst, gr->bitmap1, 0, 0, 0, 0, gr->bitmap1.width,
                      gr->bitmap1.height);
    size_t bytesize = sizeof(byte) * gr->bitmap1.width * gr->bitmap1.height;
    unsigned char *ret = (unsigned char *)malloc(bytesize);
    memcpy(ret, refdst.pixels, bytesize);
    return ret;
  } else if (gr->N == 3) {
    msdfgen::Bitmap<byte, 3> dst(gr->bitmap3.width, gr->bitmap3.height);
    msdfgen::BitmapRef<byte, 3> refdst = dst;
    fontcatalog::blit(refdst, gr->bitmap3, 0, 0, 0, 0, gr->bitmap3.width,
                      gr->bitmap3.height);
    size_t bytesize = sizeof(byte) * gr->bitmap3.width * gr->bitmap3.height * 3;
    unsigned char *ret = (unsigned char *)malloc(bytesize);
    memcpy(ret, refdst.pixels, bytesize);
    return ret;
  } else if (gr->N == 4) {
    msdfgen::Bitmap<byte, 4> dst(gr->bitmap4.width, gr->bitmap4.height);
    msdfgen::BitmapRef<byte, 4> refdst = dst;
    fontcatalog::blit(refdst, gr->bitmap4, 0, 0, 0, 0, gr->bitmap4.width,
                      gr->bitmap4.height);
    size_t bytesize = sizeof(byte) * gr->bitmap4.width * gr->bitmap4.height * 4;
    unsigned char *ret = (unsigned char *)malloc(bytesize);
    memcpy(ret, refdst.pixels, bytesize);
    return ret;
  }
  return nullptr;
}

FC_LIB_EXPORT fc_charset_t *fc_new_charset() { return new fc_charset_t{}; }

FC_LIB_EXPORT fc_charset_t *fc_new_charset_ascii() {
  return new fc_charset_t{fontcatalog::charset::ASCII};
}

FC_LIB_EXPORT void fc_charset_free(fc_charset_t *cs) { delete cs; }

FC_LIB_EXPORT size_t fc_charset_size(fc_charset_t *cs) { return cs->c.size(); }

FC_LIB_EXPORT _Bool fc_charset_empty(fc_charset_t *cs) { return cs->c.empty(); }

FC_LIB_EXPORT void fc_charset_add(fc_charset_t *cs, fc_unicode_t code) {
  cs->c.add(code);
}

FC_LIB_EXPORT void fc_charset_remove(fc_charset_t *cs, fc_unicode_t code) {
  cs->c.remove(code);
}

FC_LIB_EXPORT fc_unicode_t *fc_charset_data(fc_charset_t *cs, size_t *si) {
  fc_unicode_t *ret =
      (fc_unicode_t *)malloc(sizeof(fc_unicode_t) * cs->c.size());
  int i = 0;
  for (auto c : cs->c) {
    ret[i++] = c;
  }
  return ret;
}

FC_LIB_EXPORT void fc_scanline_generator(fc_bitmap_t *output,
                                         fc_glyph_geometry_t *glyph,
                                         fc_generator_attributes_t *attribs) {
  msdfgen::BitmapRef<float, 1> ref = output->bitmap1;
  fontcatalog::scanline_generator(ref, *glyph->g, attribs->attr);
}

FC_LIB_EXPORT void fc_sdf_generator(fc_bitmap_t *output,
                                    fc_glyph_geometry_t *glyph,
                                    fc_generator_attributes_t *attribs) {
  msdfgen::BitmapRef<float, 1> ref = output->bitmap1;
  fontcatalog::sdf_generator(ref, *glyph->g, attribs->attr);
}

FC_LIB_EXPORT void fc_psdf_generator(fc_bitmap_t *output,
                                     fc_glyph_geometry_t *glyph,
                                     fc_generator_attributes_t *attribs) {
  msdfgen::BitmapRef<float, 1> ref = output->bitmap1;
  fontcatalog::psdf_generator(ref, *glyph->g, attribs->attr);
}

FC_LIB_EXPORT void fc_msdf_generator(fc_bitmap_t *output,
                                     fc_glyph_geometry_t *glyph,
                                     fc_generator_attributes_t *attribs) {
  msdfgen::BitmapRef<float, 3> ref = output->bitmap3;
  fontcatalog::msdf_generator(ref, *glyph->g, attribs->attr);
}

FC_LIB_EXPORT void fc_mtsdf_generator(fc_bitmap_t *output,
                                      fc_glyph_geometry_t *glyph,
                                      fc_generator_attributes_t *attribs) {
  msdfgen::BitmapRef<float, 4> ref = output->bitmap4;
  fontcatalog::mtsdf_generator(ref, *glyph->g, attribs->attr);
}

#ifdef __cplusplus
}
#endif
