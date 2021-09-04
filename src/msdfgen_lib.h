#pragma once

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

#ifdef _WIN32
#define MSDF_LIB_EXPORT __declspec(dllexport)
#else
#define MSDF_LIB_EXPORT
#endif

struct _font_metrics_t {
  int ascent;
  int descent;
  int unitsPerEm;
  int baseLine;
  int lineHeight;
  int flags;
  int *characterSet;
  int charSize;
} font_metrics_t;

struct _glyph_metrics_t {
  int index;
  int width;
  int height;
  int offsetX;
  int offsetY;
  int advanceX;
  int descent;
  _Bool ccw;
} glyph_metrics_t;

typedef struct _font_handle_t font_handle_t;

MSDF_LIB_EXPORT bool wrap_initialize_freetype();
MSDF_LIB_EXPORT void wrap_deinitialize_freetype();

MSDF_LIB_EXPORT font_handle_t *
msdfgen_load_font_memory(const unsigned char *data, long size, int fontSize,
                         struct _font_metrics_t *metrics);

MSDF_LIB_EXPORT void msdfgen_free(font_handle_t *handle);

MSDF_LIB_EXPORT char *msdfgen_get_font_name(font_handle_t *font, long *size);

MSDF_LIB_EXPORT double msdfgen_get_scale(font_handle_t *font);

MSDF_LIB_EXPORT bool msdfgen_get_glyph_metrics(font_handle_t *font,
                                               int charcode,
                                               struct _glyph_metrics_t *output);

MSDF_LIB_EXPORT int msdfgen_get_kerning(font_handle_t *font, int left,
                                        int right);

MSDF_LIB_EXPORT _Bool msdfgen_generate_sdf_glyph(
    font_handle_t *font, int charcode, int width, int height, uint8_t *output,
    double tx, double ty, double range, bool normalizeShapes, _Bool ccw);
MSDF_LIB_EXPORT _Bool msdfgen_generate_psdf_glyph(
    font_handle_t *font, int charcode, int width, int height, uint8_t *output,
    double tx, double ty, double range, bool normalizeShapes, _Bool ccw);
MSDF_LIB_EXPORT _Bool msdfgen_generate_msdf_glyph(
    font_handle_t *font, int charcode, int width, int height, uint8_t *output,
    double tx, double ty, double range, bool normalizeShapes, _Bool ccw);
MSDF_LIB_EXPORT bool msdfgen_rasterize_glyph(font_handle_t *font, int charcode,
                                             int width, int height,
                                             uint8_t *output, int ox, int oy);

MSDF_LIB_EXPORT struct _font_metrics_t msdfgen_get_font_info(char *filename);

#ifdef __cplusplus
}
#endif