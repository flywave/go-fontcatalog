#include "glyph_generators.hh"

namespace fontcatalog {

void scanline_generator(const msdfgen::BitmapRef<float, 1> &output,
                        const glyph_geometry &glyph,
                        const generator_attributes &attribs) {
  msdfgen::rasterize(output, glyph.get_shape(), glyph.get_box_scale(),
                     glyph.get_box_translate(), MSDF_ATLAS_GLYPH_FILL_RULE);
}

void sdf_generator(const msdfgen::BitmapRef<float, 1> &output,
                   const glyph_geometry &glyph,
                   const generator_attributes &attribs) {
  msdfgen::generateSDF(output, glyph.get_shape(), glyph.get_box_projection(),
                       glyph.get_box_range(), attribs.config);
  if (attribs.scanlinePass)
    msdfgen::distanceSignCorrection(output, glyph.get_shape(),
                                    glyph.get_box_projection(),
                                    MSDF_ATLAS_GLYPH_FILL_RULE);
}

void psdf_generator(const msdfgen::BitmapRef<float, 1> &output,
                    const glyph_geometry &glyph,
                    const generator_attributes &attribs) {
  msdfgen::generatePseudoSDF(output, glyph.get_shape(),
                             glyph.get_box_projection(), glyph.get_box_range(),
                             attribs.config);
  if (attribs.scanlinePass)
    msdfgen::distanceSignCorrection(output, glyph.get_shape(),
                                    glyph.get_box_projection(),
                                    MSDF_ATLAS_GLYPH_FILL_RULE);
}

void msdf_generator(const msdfgen::BitmapRef<float, 3> &output,
                    const glyph_geometry &glyph,
                    const generator_attributes &attribs) {
  msdfgen::MSDFGeneratorConfig config = attribs.config;
  if (attribs.scanlinePass)
    config.errorCorrection.mode = msdfgen::ErrorCorrectionConfig::DISABLED;
  msdfgen::generateMSDF(output, glyph.get_shape(), glyph.get_box_projection(),
                        glyph.get_box_range(), config);
  if (attribs.scanlinePass) {
    msdfgen::distanceSignCorrection(output, glyph.get_shape(),
                                    glyph.get_box_projection(),
                                    MSDF_ATLAS_GLYPH_FILL_RULE);
    if (attribs.config.errorCorrection.mode !=
        msdfgen::ErrorCorrectionConfig::DISABLED) {
      config.errorCorrection.mode = attribs.config.errorCorrection.mode;
      config.errorCorrection.distanceCheckMode =
          msdfgen::ErrorCorrectionConfig::DO_NOT_CHECK_DISTANCE;
      msdfgen::msdfErrorCorrection(output, glyph.get_shape(),
                                   glyph.get_box_projection(),
                                   glyph.get_box_range(), config);
    }
  }
}

void mtsdf_generator(const msdfgen::BitmapRef<float, 4> &output,
                     const glyph_geometry &glyph,
                     const generator_attributes &attribs) {
  msdfgen::MSDFGeneratorConfig config = attribs.config;
  if (attribs.scanlinePass)
    config.errorCorrection.mode = msdfgen::ErrorCorrectionConfig::DISABLED;
  msdfgen::generateMTSDF(output, glyph.get_shape(), glyph.get_box_projection(),
                         glyph.get_box_range(), config);
  if (attribs.scanlinePass) {
    msdfgen::distanceSignCorrection(output, glyph.get_shape(),
                                    glyph.get_box_projection(),
                                    MSDF_ATLAS_GLYPH_FILL_RULE);
    if (attribs.config.errorCorrection.mode !=
        msdfgen::ErrorCorrectionConfig::DISABLED) {
      config.errorCorrection.mode = attribs.config.errorCorrection.mode;
      config.errorCorrection.distanceCheckMode =
          msdfgen::ErrorCorrectionConfig::DO_NOT_CHECK_DISTANCE;
      msdfgen::msdfErrorCorrection(output, glyph.get_shape(),
                                   glyph.get_box_projection(),
                                   glyph.get_box_range(), config);
    }
  }
}

} // namespace fontcatalog
