#include "msdfgen.h"
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

#include "msdfgen-ext.h"
#include "msdfgen_lib.h"

#ifdef __cplusplus
extern "C" {
#endif

using namespace msdfgen;

struct _font_handle_t {
  FontHandle *handle;
  double scale;
};

FreetypeHandle *ft = NULL;
bool enforceR8 = false;

MSDF_LIB_EXPORT font_handle_t *
msdfgen_load_font_memory(const unsigned char *data, long size, int fontSize,
                         struct _font_metrics_t *metrics) {
  FontHandle *handle = msdfgen::loadFontData(ft, data, size);
  FT_Face ft = msdfgen::getFreetypeFont(handle);

  double scale = (double)fontSize / (double)ft->units_per_EM * 64.;

  if (metrics != nullptr) {
    metrics->ascent = ft->ascender;
    metrics->descent = ft->descender;
    metrics->unitsPerEm = ft->units_per_EM;
    metrics->baseLine = (ft->size->metrics.ascender + 32) >> 6;
    metrics->lineHeight = (ft->size->metrics.height + 32) >> 6;
    TT_Header *header = (TT_Header *)FT_Get_Sfnt_Table(ft, FT_SFNT_HEAD);
    metrics->flags = (int)header->Mac_Style | header->Flags << 16;
  }
  return new _font_handle_t{handle, scale};
}

MSDF_LIB_EXPORT void msdfgen_free(font_handle_t *handle) {}

MSDF_LIB_EXPORT char *msdfgen_get_font_name(font_handle_t *font, long *size) {
  FT_Face face = msdfgen::getFreetypeFont(font->handle);
  FT_SfntName name;
  int count = FT_Get_Sfnt_Name_Count(face);
  for (int i = 0; i < count; i++) {
    FT_Get_Sfnt_Name(face, i, &name);
    if (name.name_id == 4 && (name.platform_id == 3 || name.platform_id == 0) &&
        name.language_id == 0x409) {
      char *data = (char *)malloc(name.string_len);
      *size = name.string_len;
      memcpy(data, name.string, name.string_len);
      return data;
    }
  }
  *size = 0;
  return nullptr;
}

MSDF_LIB_EXPORT double msdfgen_get_scale(font_handle_t *font) {
  return font->scale;
}

MSDF_LIB_EXPORT bool
msdfgen_get_glyph_metrics(font_handle_t *font, int charcode,
                          struct _glyph_metrics_t *metrics) {
  FT_Face face = msdfgen::getFreetypeFont(font->handle);
  FT_UInt index = FT_Get_Char_Index(face, charcode);
  if (index == 0 && charcode != 0)
    return false;
  FT_Error err = FT_Load_Glyph(face, index, FT_LOAD_DEFAULT);
  if (err)
    return false;
  if (metrics != nullptr) {
    FT_GlyphSlot slot = face->glyph;
    metrics->width = slot->bitmap.width;
    metrics->height = slot->bitmap.rows;
    metrics->offsetX = slot->bitmap_left;
    metrics->offsetY = slot->bitmap_top;
    metrics->advanceX = (slot->advance.x + 32) >> 6;
    metrics->descent = (slot->metrics.horiBearingY - slot->metrics.height) >> 6;
    metrics->ccw = FT_Outline_Get_Orientation(&slot->outline);
  }
  return true;
}

MSDF_LIB_EXPORT int msdfgen_get_kerning(font_handle_t *font, int left,
                                        int right) {
  FT_Vector vec;
  FT_Face face = msdfgen::getFreetypeFont(font->handle);
  FT_Error err =
      FT_Get_Kerning(face, FT_Get_Char_Index(face, left),
                     FT_Get_Char_Index(face, right), FT_KERNING_DEFAULT, &vec);
  if (err)
    return err;
  return (vec.x + 32) >> 6;
}

MSDF_LIB_EXPORT bool wrap_initialize_freetype() {
  ft = initializeFreetype();
  return ft != NULL;
}

MSDF_LIB_EXPORT void wrap_deinitialize_freetype() {
  if (ft != NULL) {
    deinitializeFreetype(ft);
  }
}

void normalizeShape(Shape &shape, bool normalizeShapes) {
  if (normalizeShapes) {
    shape.normalize();
  } else {
    for (std::vector<msdfgen::Contour>::iterator contour =
             shape.contours.begin();
         contour != shape.contours.end(); ++contour) {
      if (contour->edges.size() == 1) {
        contour->edges.clear();
      }
    }
  }
}

MSDF_LIB_EXPORT _Bool msdfgen_generate_sdf_glyph(
    font_handle_t *font, int charcode, int width, int height, uint8_t *output,
    double tx, double ty, double range, bool normalizeShapes, _Bool ccw) {
  if (width == 0 || height == 0)
    return true;

  Shape glyph;
  if (loadGlyph(glyph, font->handle, charcode)) {
    normalizeShape(glyph, normalizeShapes);
    Bitmap<float, 1> sdf(width, height);
    double scale = font->scale;
    generateSDF(sdf, glyph, range / scale, scale,
                Vector2(tx / scale, ty / scale));
    if (ccw) {
      for (int y = height - 1; y >= 0; y--) {
        uint8_t *it = &output[(height - y) * width * 4];
        for (int x = 0; x < width; x++) {
          uint8_t px = uint8_t(pixelFloatToByte(1.f - *sdf(x, y)));
          *it++ = px;
          *it++ = px;
          *it++ = px;
          *it++ = 0xff;
        }
      }
    } else {
      for (int y = height - 1; y >= 0; y--) {
        uint8_t *it = &output[(height - y) * width * 4];
        for (int x = 0; x < width; x++) {
          uint8_t px = uint8_t(pixelFloatToByte(*sdf(x, y)));
          *it++ = px;
          *it++ = px;
          *it++ = px;
          *it++ = 0xff;
        }
      }
    }
    return true;
  }
  return false;
}

MSDF_LIB_EXPORT _Bool msdfgen_generate_psdf_glyph(
    font_handle_t *font, int charcode, int width, int height, uint8_t *output,
    double tx, double ty, double range, bool normalizeShapes, _Bool ccw) {
  if (width == 0 || height == 0)
    return true;

  Shape glyph;
  if (loadGlyph(glyph, font->handle, charcode)) {
    normalizeShape(glyph, normalizeShapes);
    Bitmap<float, 1> sdf(width, height);
    double scale = font->scale;
    generatePseudoSDF(sdf, glyph, range / scale, scale,
                      Vector2(tx / scale, ty / scale));
    if (ccw) {
      for (int y = height - 1; y >= 0; y--) {
        uint8_t *it = &output[(height - y) * width * 4];
        for (int x = 0; x < width; x++) {
          uint8_t px = uint8_t(pixelFloatToByte(1.f - *sdf(x, y)));
          *it++ = px;
          *it++ = px;
          *it++ = px;
          *it++ = 0xff;
        }
      }
    } else {
      for (int y = height - 1; y >= 0; y--) {
        uint8_t *it = &output[(height - y) * width * 4];
        for (int x = 0; x < width; x++) {
          uint8_t px = uint8_t(pixelFloatToByte(*sdf(x, y)));
          *it++ = px;
          *it++ = px;
          *it++ = px;
          *it++ = 0xff;
        }
      }
    }
    return true;
  }
  return false;
}

MSDF_LIB_EXPORT _Bool msdfgen_generate_msdf_glyph(
    font_handle_t *font, int charcode, int width, int height, uint8_t *output,
    double tx, double ty, double range, bool normalizeShapes, _Bool ccw) {
  if (width == 0 || height == 0)
    return true;
  Shape glyph;
  if (loadGlyph(glyph, font->handle, charcode)) {
    normalizeShape(glyph, normalizeShapes);
    edgeColoringSimple(glyph, 3, 0);
    Bitmap<float, 3> msdf(width, height);
    double scale = font->scale;
    generateMSDF(msdf, glyph, range / scale, scale,
                 Vector2(tx / scale, ty / scale));
    if (ccw) {
      for (int y = height - 1; y >= 0; y--) {
        uint8_t *it = &output[(height - y) * width * 4];
        for (int x = 0; x < width; x++) {
          *it++ = uint8_t(pixelFloatToByte(1.f - msdf(x, y)[0]));
          *it++ = uint8_t(pixelFloatToByte(1.f - msdf(x, y)[1]));
          *it++ = uint8_t(pixelFloatToByte(1.f - msdf(x, y)[2]));
          *it++ = 0xff;
        }
      }
    } else {
      for (int y = height - 1; y >= 0; y--) {
        uint8_t *it = &output[(height - y) * width * 4];
        for (int x = 0; x < width; x++) {
          *it++ = uint8_t(pixelFloatToByte(msdf(x, y)[0]));
          *it++ = uint8_t(pixelFloatToByte(msdf(x, y)[1]));
          *it++ = uint8_t(pixelFloatToByte(msdf(x, y)[2]));
          *it++ = 0xff;
        }
      }
    }
    return true;
  }
  return false;
}

MSDF_LIB_EXPORT bool msdfgen_rasterize_glyph(font_handle_t *font, int charcode,
                                             int width, int height,
                                             uint8_t *output, int ox, int oy) {
  if (width == 0 || height == 0)
    return true;

  FT_Face face = msdfgen::getFreetypeFont(font->handle);
  FT_Error err = FT_Load_Char(face, charcode, FT_LOAD_RENDER);
  if (err)
    return false;
  FT_Bitmap *bitmap = &face->glyph->bitmap;
  FT_Library ft_lib = getFreetypeLibrary(ft);
  int multiplier = 1;
  switch (bitmap->pixel_mode) {
  case ft_pixel_mode_mono:
    FT_Bitmap grayBtm;
    FT_Bitmap_Init(&grayBtm);
    FT_Bitmap_Convert(ft_lib, bitmap, &grayBtm, 1);
    bitmap = &grayBtm;
    multiplier = 0xff;
  case ft_pixel_mode_grays:
    if (enforceR8) {
      for (int y = 0; y < height; y++) {
        uint8_t *it = &output[(oy + y) * 4 + ox];
        for (int x = 0; x < width; x++) {
          unsigned char px = bitmap->buffer[(y)*bitmap->width + x] * multiplier;
          *it++ = px;
          *it++ = px;
          *it++ = px;
          *it++ = 0xff;
        }
      }
    } else {
      for (int y = 0; y < height; y++) {
        uint8_t *it = &output[(oy + y) * 4 + ox];
        for (int x = 0; x < width; x++) {
          unsigned char px = bitmap->buffer[(y)*bitmap->width + x] * multiplier;
          *it++ = 0xff;
          *it++ = 0xff;
          *it++ = 0xff;
          *it++ = px;
        }
      }
    }

    if (bitmap->pixel_mode == FT_PIXEL_MODE_MONO) {
      FT_Bitmap_Done(ft_lib, &grayBtm);
    }
    return true;
    break;
  }
  return false;
}

#ifdef __cplusplus
}
#endif
