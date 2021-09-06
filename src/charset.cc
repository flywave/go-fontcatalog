#include "charset.hh"

namespace fontcatalog {

static charset create_ascii_charset() {
  charset ascii;
  for (unicode_t cp = 0x20; cp < 0x7f; ++cp)
    ascii.add(cp);
  return ascii;
}

const charset charset::ASCII = create_ascii_charset();

void charset::add(unicode_t cp) { _codepoints.insert(cp); }

void charset::remove(unicode_t cp) { _codepoints.erase(cp); }

size_t charset::size() const { return _codepoints.size(); }

bool charset::empty() const { return _codepoints.empty(); }

std::set<unicode_t>::const_iterator charset::begin() const {
  return _codepoints.begin();
}

std::set<unicode_t>::const_iterator charset::end() const {
  return _codepoints.end();
}

} // namespace fontcatalog
