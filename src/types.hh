#pragma once

#include <cstdint>

namespace fontcatalog {

typedef unsigned char byte;
typedef uint32_t unicode_t;

enum class glyph_identifier_type { GLYPH_INDEX, UNICODE_CODEPOINT };

} // namespace fontcatalog