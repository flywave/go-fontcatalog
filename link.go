package fontcatalog

// #include <stdlib.h>
// #include "fontcatalog_lib.h"
// #cgo CFLAGS: -I ./lib
// #cgo linux CXXFLAGS: -I ./lib -std=c++14
// #cgo darwin CXXFLAGS: -I ./lib  -std=gnu++14
// #cgo darwin,arm CXXFLAGS: -I ./lib  -std=gnu++14
// #cgo darwin LDFLAGS: -L ./lib/darwin -lpng -lzlib -lharfbuzz -lfreetype -lmsdfgen -lmsdfgen_ext -lfontcatalog
// #cgo darwin,arm LDFLAGS: -L ./lib/darwin_arm -lpng -lzlib -lharfbuzz -lfreetype -lmsdfgen -lmsdfgen_ext -lfontcatalog
// #cgo linux LDFLAGS: -L ./lib/linux -Wl,--start-group -lpthread -ldl -lstdc++ -lm -lpng -lzlib -lharfbuzz -lfreetype -lmsdfgen -lmsdfgen_ext -lfontcatalog -Wl,--end-group
import "C"
