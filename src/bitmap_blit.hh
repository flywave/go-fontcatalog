#pragma once

#include "types.hh"
#include <msdfgen.h>

namespace fontcatalog {

void blit(const msdfgen::BitmapRef<byte, 1> &dst,
          const msdfgen::BitmapConstRef<byte, 1> &src, int dx, int dy, int sx,
          int sy, int w, int h);
void blit(const msdfgen::BitmapRef<byte, 3> &dst,
          const msdfgen::BitmapConstRef<byte, 3> &src, int dx, int dy, int sx,
          int sy, int w, int h);
void blit(const msdfgen::BitmapRef<byte, 4> &dst,
          const msdfgen::BitmapConstRef<byte, 4> &src, int dx, int dy, int sx,
          int sy, int w, int h);

void blit(const msdfgen::BitmapRef<float, 1> &dst,
          const msdfgen::BitmapConstRef<float, 1> &src, int dx, int dy, int sx,
          int sy, int w, int h);
void blit(const msdfgen::BitmapRef<float, 3> &dst,
          const msdfgen::BitmapConstRef<float, 3> &src, int dx, int dy, int sx,
          int sy, int w, int h);
void blit(const msdfgen::BitmapRef<float, 4> &dst,
          const msdfgen::BitmapConstRef<float, 4> &src, int dx, int dy, int sx,
          int sy, int w, int h);

void blit(const msdfgen::BitmapRef<byte, 1> &dst,
          const msdfgen::BitmapConstRef<float, 1> &src, int dx, int dy, int sx,
          int sy, int w, int h);
void blit(const msdfgen::BitmapRef<byte, 3> &dst,
          const msdfgen::BitmapConstRef<float, 3> &src, int dx, int dy, int sx,
          int sy, int w, int h);
void blit(const msdfgen::BitmapRef<byte, 4> &dst,
          const msdfgen::BitmapConstRef<float, 4> &src, int dx, int dy, int sx,
          int sy, int w, int h);

} // namespace fontcatalog
