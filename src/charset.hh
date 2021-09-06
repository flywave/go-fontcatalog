#pragma once

#include "types.hh"

#include <cstdlib>
#include <set>

namespace fontcatalog {

class charset {
public:
  static const charset ASCII;

  void add(unicode_t cp);
  void remove(unicode_t cp);

  size_t size() const;
  bool empty() const;
  std::set<unicode_t>::const_iterator begin() const;
  std::set<unicode_t>::const_iterator end() const;

private:
  std::set<unicode_t> _codepoints;
};

} // namespace fontcatalog
