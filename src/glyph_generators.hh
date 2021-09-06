
#pragma once

#include "generator_attributes.hh"
#include "glyph_geometry.hh"

#include <msdfgen.h>

#define MSDF_ATLAS_GLYPH_FILL_RULE msdfgen::FILL_NONZERO

namespace fontcatalog {

void scanline_generator(const msdfgen::BitmapRef<float, 1> &output,
                        const glyph_geometry &glyph,
                        const generator_attributes &attribs);
void sdf_generator(const msdfgen::BitmapRef<float, 1> &output,
                   const glyph_geometry &glyph,
                   const generator_attributes &attribs);
void psdf_generator(const msdfgen::BitmapRef<float, 1> &output,
                    const glyph_geometry &glyph,
                    const generator_attributes &attribs);
void msdf_generator(const msdfgen::BitmapRef<float, 3> &output,
                    const glyph_geometry &glyph,
                    const generator_attributes &attribs);
void mtsdf_generator(const msdfgen::BitmapRef<float, 4> &output,
                     const glyph_geometry &glyph,
                     const generator_attributes &attribs);

} // namespace fontcatalog
