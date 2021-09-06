#pragma once

#include <msdfgen.h>

namespace fontcatalog {

struct generator_attributes {
  msdfgen::MSDFGeneratorConfig config;
  bool scanlinePass = false;
};

} // namespace fontcatalog
