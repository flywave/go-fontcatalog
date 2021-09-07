#pragma once

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

#ifdef _WIN32
#define FC_LIB_EXPORT __declspec(dllexport)
#else
#define FC_LIB_EXPORT
#endif

typedef unsigned fc_glyph_index_t;
typedef uint32_t fc_unicode_t;

typedef struct _fc_font_metrics_t {
  double emSize;
  double ascenderY, descenderY;
  double lineHeight;
  double underlineY, underlineThickness;
} fc_font_metrics_t;

typedef struct _fc_glyph_box_t {
  int index;
  double advance;
  struct {
    double l, b, r, t;
  } bounds;
  struct {
    int x, y, w, h;
  } rect;
} fc_glyph_box_t;

typedef struct _fc_kerning_t {
  int first, second;
  double kerning;
} fc_kerning_t;

typedef struct _fc_font_info_t {
  int ascent;
  int descent;
  int unitsPerEm;
  int baseLine;
  int lineHeight;
  int flags;
  int *characterSet;
  int charSize;
} fc_font_info_t;

typedef struct _fc_font_holder_t fc_font_holder_t;
typedef struct _fc_font_geometry_t fc_font_geometry_t;
typedef struct _fc_font_geometry_list_t fc_font_geometry_list_t;
typedef struct _fc_glyph_geometry_t fc_glyph_geometry_t;
typedef struct _fc_glyph_geometry_list_t fc_glyph_geometry_list_t;
typedef struct _fc_glyph_range_t fc_glyph_range_t;
typedef struct _fc_generator_attributes_t fc_generator_attributes_t;
typedef struct _fc_bitmap_t fc_bitmap_t;
typedef struct _fc_bitmap_ref_t fc_bitmap_ref_t;
typedef struct _fc_charset_t fc_charset_t;
typedef struct _fc_kerning_map_t fc_kerning_map_t;

FC_LIB_EXPORT fc_font_holder_t *
fc_font_holder_load_font_memory(const unsigned char *data, long size);
FC_LIB_EXPORT void fc_font_holder_free(fc_font_holder_t *handle);
FC_LIB_EXPORT struct _fc_font_info_t fc_font_holder_get_font_info(fc_font_holder_t *handle);

FC_LIB_EXPORT fc_glyph_geometry_t *fc_new_glyph_geometry_from_glyph_index(
    fc_font_holder_t *handle, double geometryScale, fc_glyph_index_t index);
FC_LIB_EXPORT fc_glyph_geometry_t *fc_new_glyph_geometry_from_unicode(
    fc_font_holder_t *handle, double geometryScale, fc_unicode_t codepoint);
FC_LIB_EXPORT void fc_glyph_geometry_free(fc_glyph_geometry_t *geom);

FC_LIB_EXPORT void fc_glyph_geometry_edge_coloring(fc_glyph_geometry_t *geom,
                                                   uint32_t type,
                                                   double angleThreshold,
                                                   unsigned long long seed);
FC_LIB_EXPORT void fc_glyph_geometry_wrap_box(fc_glyph_geometry_t *geom,
                                              double scale, double range,
                                              double miterLimit);
FC_LIB_EXPORT void fc_glyph_geometry_place_box(fc_glyph_geometry_t *geom, int x,
                                               int y);
FC_LIB_EXPORT int fc_glyph_geometry_get_index(fc_glyph_geometry_t *geom);
FC_LIB_EXPORT fc_glyph_index_t
fc_glyph_geometry_get_glyph_index(fc_glyph_geometry_t *geom);
FC_LIB_EXPORT fc_unicode_t
fc_glyph_geometry_get_codepoint(fc_glyph_geometry_t *geom);
FC_LIB_EXPORT int fc_glyph_geometry_get_identifier(fc_glyph_geometry_t *geom,
                                                   uint32_t type);
FC_LIB_EXPORT double fc_glyph_geometry_get_advance(fc_glyph_geometry_t *geom);
FC_LIB_EXPORT void fc_glyph_geometry_get_box_rect(fc_glyph_geometry_t *geom,
                                                  int *x, int *y, int *w,
                                                  int *h);
FC_LIB_EXPORT void fc_glyph_geometry_get_box_size(fc_glyph_geometry_t *geom,
                                                  int *w, int *h);
FC_LIB_EXPORT double fc_glyph_geometry_get_box_range(fc_glyph_geometry_t *geom);
FC_LIB_EXPORT double fc_glyph_geometry_get_box_scale(fc_glyph_geometry_t *geom);
FC_LIB_EXPORT void
fc_glyph_geometry_get_box_translate(fc_glyph_geometry_t *geom, int *tx,
                                    int *ty);
FC_LIB_EXPORT fc_glyph_box_t
fc_glyph_geometry_get_glyph_box(fc_glyph_geometry_t *geom);
FC_LIB_EXPORT _Bool fc_glyph_geometry_is_whitespace(fc_glyph_geometry_t *geom);

FC_LIB_EXPORT fc_glyph_geometry_list_t *fc_new_glyph_geometry_list();
FC_LIB_EXPORT void fc_glyph_geometry_list_free(fc_glyph_geometry_list_t *list);
FC_LIB_EXPORT void
fc_glyph_geometry_list_push_geometry(fc_glyph_geometry_list_t *list,
                                     fc_glyph_geometry_t *geom);
FC_LIB_EXPORT _Bool
fc_glyph_geometry_list_empty(fc_glyph_geometry_list_t *list);
FC_LIB_EXPORT size_t
fc_glyph_geometry_list_size(fc_glyph_geometry_list_t *list);

FC_LIB_EXPORT fc_font_geometry_t *
fc_new_font_geometry_with_glyphs(fc_glyph_geometry_list_t *glyphs);
FC_LIB_EXPORT void fc_font_geometry_free(fc_font_geometry_t *geom);
FC_LIB_EXPORT int fc_font_geometry_load_from_glyphset(fc_font_geometry_t *fonts,
                                                      fc_font_holder_t *handle,
                                                      double fontScale,
                                                      fc_charset_t *charsets);
FC_LIB_EXPORT int fc_font_geometry_load_from_charset(fc_font_geometry_t *fonts,
                                                     fc_font_holder_t *handle,
                                                     double fontScale,
                                                     fc_charset_t *charsets);
FC_LIB_EXPORT _Bool fc_font_geometry_load_metrics(fc_font_geometry_t *fonts,
                                                  fc_font_holder_t *handle,
                                                  double fontScale);
FC_LIB_EXPORT _Bool fc_font_geometry_add_glyph(fc_font_geometry_t *fonts,
                                               fc_glyph_geometry_t *geom);
FC_LIB_EXPORT int fc_font_geometry_load_kerning(fc_font_geometry_t *fonts,
                                                fc_font_holder_t *handle);
FC_LIB_EXPORT void fc_font_geometry_set_name(fc_font_geometry_t *geom,
                                             const char *name);
FC_LIB_EXPORT const char *fc_font_geometry_get_name(fc_font_geometry_t *geom);
FC_LIB_EXPORT double fc_font_geometry_geometry_scale(fc_font_geometry_t *fonts);
FC_LIB_EXPORT struct _fc_font_metrics_t
fc_font_geometry_get_metrics(fc_font_geometry_t *fonts);
FC_LIB_EXPORT uint32_t
fc_font_geometry_get_preferred_identifier_type(fc_font_geometry_t *fonts);
FC_LIB_EXPORT fc_glyph_range_t *
fc_font_geometry_get_glyphs(fc_font_geometry_t *fonts);
FC_LIB_EXPORT fc_glyph_geometry_t *
fc_font_geometry_get_glyph_from_index(fc_font_geometry_t *fonts,
                                      fc_glyph_index_t index);
FC_LIB_EXPORT fc_glyph_geometry_t *
fc_font_geometry_get_glyph_from_unicode(fc_font_geometry_t *fonts,
                                        fc_unicode_t codePoint);
FC_LIB_EXPORT _Bool fc_font_geometry_get_advance_from_index(
    fc_font_geometry_t *fonts, double *advance, fc_glyph_index_t index1,
    fc_glyph_index_t index2);
FC_LIB_EXPORT _Bool fc_font_geometry_get_advance_from_unicode(
    fc_font_geometry_t *fonts, double *advance, fc_unicode_t codePoint1,
    fc_unicode_t codePoint2);
FC_LIB_EXPORT fc_kerning_map_t *
fc_font_geometry_get_kerning(fc_font_geometry_t *fonts);

FC_LIB_EXPORT void fc_kerning_map_free(fc_kerning_map_t *kmap);
FC_LIB_EXPORT fc_kerning_t *fc_kerning_map_get_kernings(fc_kerning_map_t *kmap,
                                                        size_t *si);

FC_LIB_EXPORT fc_font_geometry_list_t *fc_new_font_geometry_list();
FC_LIB_EXPORT void fc_font_geometry_list_free(fc_font_geometry_list_t *list);
FC_LIB_EXPORT void
fc_font_geometry_list_push_geometry(fc_font_geometry_list_t *list,
                                    fc_font_geometry_t *geom);
FC_LIB_EXPORT _Bool fc_font_geometry_list_empty(fc_font_geometry_list_t *list);
FC_LIB_EXPORT size_t fc_font_geometry_list_size(fc_font_geometry_list_t *list);

FC_LIB_EXPORT void fc_glyph_range_free(fc_glyph_range_t *gr);
FC_LIB_EXPORT size_t fc_glyph_range_size(fc_glyph_range_t *gr);
FC_LIB_EXPORT _Bool fc_glyph_range_empty(fc_glyph_range_t *gr);
FC_LIB_EXPORT fc_glyph_geometry_t *fc_glyph_range_get(fc_glyph_range_t *gr,
                                                      size_t index);

FC_LIB_EXPORT fc_generator_attributes_t *fc_new_generator_attributes();
FC_LIB_EXPORT void fc_generator_attributes_free(fc_generator_attributes_t *ga);
FC_LIB_EXPORT void
fc_generator_attributes_set_min_deviation_ratio(fc_generator_attributes_t *ga,
                                                double ratio);
FC_LIB_EXPORT void
fc_generator_attributes_set_min_improve_ratio(fc_generator_attributes_t *ga,
                                              double ratio);
FC_LIB_EXPORT void
fc_generator_attributes_set_mode(fc_generator_attributes_t *ga, uint32_t mode);
FC_LIB_EXPORT void
fc_generator_attributes_set_distance_check_mode(fc_generator_attributes_t *ga,
                                                uint32_t mode);
FC_LIB_EXPORT void
fc_generator_attributes_set_buffer(fc_generator_attributes_t *ga,
                                   unsigned char *buffer);
FC_LIB_EXPORT void
fc_generator_attributes_set_overlap_support(fc_generator_attributes_t *ga,
                                            _Bool overlapSupport);
FC_LIB_EXPORT void
fc_generator_attributes_set_scanline_pass(fc_generator_attributes_t *ga,
                                          _Bool scanlinePass);

FC_LIB_EXPORT fc_bitmap_t *fc_new_bitmap(int channel);
FC_LIB_EXPORT fc_bitmap_t *fc_new_bitmap_alloc(int channel, int width,
                                               int height);
FC_LIB_EXPORT void fc_bitmap_free(fc_bitmap_t *gr);
FC_LIB_EXPORT int fc_bitmap_width(fc_bitmap_t *gr);
FC_LIB_EXPORT int fc_bitmap_height(fc_bitmap_t *gr);
FC_LIB_EXPORT int fc_bitmap_channels(fc_bitmap_t *gr);
FC_LIB_EXPORT float *fc_bitmap_data(fc_bitmap_t *gr);
FC_LIB_EXPORT unsigned char *fc_bitmap_blit_data(fc_bitmap_t *gr);

FC_LIB_EXPORT fc_charset_t *fc_new_charset();
FC_LIB_EXPORT fc_charset_t *fc_new_charset_ascii();
FC_LIB_EXPORT void fc_charset_free(fc_charset_t *cs);
FC_LIB_EXPORT size_t fc_charset_size(fc_charset_t *cs);
FC_LIB_EXPORT _Bool fc_charset_empty(fc_charset_t *cs);
FC_LIB_EXPORT void fc_charset_add(fc_charset_t *cs, fc_unicode_t code);
FC_LIB_EXPORT void fc_charset_remove(fc_charset_t *cs, fc_unicode_t code);
FC_LIB_EXPORT fc_unicode_t *fc_charset_data(fc_charset_t *cs, size_t *si);

FC_LIB_EXPORT fc_bitmap_ref_t *fc_new_bitmap_ref(float *pixels, int channel,
                                                 int width, int height);
FC_LIB_EXPORT void fc_bitmap_ref_free(fc_bitmap_ref_t *gr);
FC_LIB_EXPORT int fc_bitmap_ref_width(fc_bitmap_ref_t *gr);
FC_LIB_EXPORT int fc_bitmap_ref_height(fc_bitmap_ref_t *gr);
FC_LIB_EXPORT int fc_bitmap_ref_channels(fc_bitmap_ref_t *gr);
FC_LIB_EXPORT float *fc_bitmap_ref_data(fc_bitmap_ref_t *gr);
FC_LIB_EXPORT unsigned char *fc_bitmap_ref_blit_data(fc_bitmap_ref_t *gr);

FC_LIB_EXPORT void fc_scanline_generator(fc_bitmap_t *output,
                                         fc_glyph_geometry_t *glyph,
                                         fc_generator_attributes_t *attribs);
FC_LIB_EXPORT void fc_sdf_generator(fc_bitmap_t *output,
                                    fc_glyph_geometry_t *glyph,
                                    fc_generator_attributes_t *attribs);
FC_LIB_EXPORT void fc_psdf_generator(fc_bitmap_t *output,
                                     fc_glyph_geometry_t *glyph,
                                     fc_generator_attributes_t *attribs);
FC_LIB_EXPORT void fc_msdf_generator(fc_bitmap_t *output,
                                     fc_glyph_geometry_t *glyph,
                                     fc_generator_attributes_t *attribs);
FC_LIB_EXPORT void fc_mtsdf_generator(fc_bitmap_t *output,
                                      fc_glyph_geometry_t *glyph,
                                      fc_generator_attributes_t *attribs);

#ifdef __cplusplus
}
#endif