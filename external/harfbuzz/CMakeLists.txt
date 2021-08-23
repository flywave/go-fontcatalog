

option(HARFBUZZ_BUILD_ICU "Enable building of ICU" OFF)
option(HARFBUZZ_BUILD_UCDN "Enable building of UCDN" ON)

if (HARFBUZZ_BUILD_ICU)
    set(HARFBUZZ_UNICODE_SRC ${CMAKE_CURRENT_SOURCE_DIR}/src/hb-ucdn.cc)
    set(HARFBUZZ_UNICODE_INCLUDE ${CMAKE_CURRENT_SOURCE_DIR}/../icu/source/common)
    set(HARFBUZZ_UNICODE_LIB icucommon)
    set(HARFBUZZ_UNICODE_DEFINE -DHAVE_ICU -DHAVE_ICU_BUILTIN)
endif()

if (HARFBUZZ_BUILD_UCDN)
    set(HARFBUZZ_UNICODE_SRC ${CMAKE_CURRENT_SOURCE_DIR}/src/hb-ft.cc)
    set(HARFBUZZ_UNICODE_INCLUDE ${CMAKE_CURRENT_SOURCE_DIR}/src/hb-ucdn/)
    set(HARFBUZZ_UNICODE_DEFINE -DHAVE_UCDN)
endif()

set(HARFBUZZ_SRC_FILES
    ${HARFBUZZ_UNICODE_SRC}
    ${CMAKE_CURRENT_SOURCE_DIR}/src/hb-ucdn/ucdn.c
    ${CMAKE_CURRENT_SOURCE_DIR}/src/hb-ucdn.cc
    ${CMAKE_CURRENT_SOURCE_DIR}/src/hb-blob.cc
    ${CMAKE_CURRENT_SOURCE_DIR}/src/hb-buffer.cc
    ${CMAKE_CURRENT_SOURCE_DIR}/src/hb-common.cc
    ${CMAKE_CURRENT_SOURCE_DIR}/src/hb-face.cc
    ${CMAKE_CURRENT_SOURCE_DIR}/src/hb-font.cc
    ${CMAKE_CURRENT_SOURCE_DIR}/src/hb-set.cc
    ${CMAKE_CURRENT_SOURCE_DIR}/src/hb-shape.cc
    ${CMAKE_CURRENT_SOURCE_DIR}/src/hb-shape-plan.cc
    ${CMAKE_CURRENT_SOURCE_DIR}/src/hb-shaper.cc
    ${CMAKE_CURRENT_SOURCE_DIR}/src/hb-unicode.cc
    ${CMAKE_CURRENT_SOURCE_DIR}/src/hb-warning.cc
    ${CMAKE_CURRENT_SOURCE_DIR}/src/hb-ot-tag.cc
    ${CMAKE_CURRENT_SOURCE_DIR}/src/hb-ot-font.cc
    ${CMAKE_CURRENT_SOURCE_DIR}/src/hb-ot-layout.cc
    ${CMAKE_CURRENT_SOURCE_DIR}/src/hb-ot-map.cc
    ${CMAKE_CURRENT_SOURCE_DIR}/src/hb-ot-shape.cc
    ${CMAKE_CURRENT_SOURCE_DIR}/src/hb-ot-var.cc
    ${CMAKE_CURRENT_SOURCE_DIR}/src/hb-fallback-shape.cc
    ${CMAKE_CURRENT_SOURCE_DIR}/src/hb-ot-shape-complex-arabic.cc
    ${CMAKE_CURRENT_SOURCE_DIR}/src/hb-ot-shape-complex-default.cc
    ${CMAKE_CURRENT_SOURCE_DIR}/src/hb-ot-shape-complex-hangul.cc
    ${CMAKE_CURRENT_SOURCE_DIR}/src/hb-ot-shape-complex-hebrew.cc
    ${CMAKE_CURRENT_SOURCE_DIR}/src/hb-ot-shape-complex-indic.cc
    ${CMAKE_CURRENT_SOURCE_DIR}/src/hb-ot-shape-complex-indic-table.cc
    ${CMAKE_CURRENT_SOURCE_DIR}/src/hb-ot-shape-complex-myanmar.cc
    ${CMAKE_CURRENT_SOURCE_DIR}/src/hb-ot-shape-complex-thai.cc
    ${CMAKE_CURRENT_SOURCE_DIR}/src/hb-ot-shape-complex-tibetan.cc
    ${CMAKE_CURRENT_SOURCE_DIR}/src/hb-ot-shape-complex-use.cc
    ${CMAKE_CURRENT_SOURCE_DIR}/src/hb-ot-shape-complex-use-table.cc
    ${CMAKE_CURRENT_SOURCE_DIR}/src/hb-ot-shape-normalize.cc
    ${CMAKE_CURRENT_SOURCE_DIR}/src/hb-ot-shape-fallback.cc
    )

INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/generated" )

list(APPEND HARFBUZZ_SRC_FILES "${CMAKE_CURRENT_SOURCE_DIR}/generated/config.h")
add_library(harfbuzz STATIC ${HARFBUZZ_SRC_FILES})

if(CMAKE_CXX_COMPILER_ID STREQUAL "MSVC")
  set_source_files_properties(${HARFBUZZ_SRC_FILES}
     PROPERTIES
     COMPILE_FLAGS  "-utf-8 -wd4996 -wd4267 -wd4244")
endif()

set(HARFBUZZ_INCLUDE_DIRS
    ${CMAKE_CURRENT_SOURCE_DIR}/src
    ${PROJECT_BINARY_DIR}/include/freetype/config
    ${CMAKE_CURRENT_SOURCE_DIR}/../freetype2/include/
    ${CMAKE_CURRENT_SOURCE_DIR}/../freetype2/src/include/
    ${HARFBUZZ_UNICODE_INCLUDE}
    ${CMAKE_CURRENT_SOURCE_DIR}/
    CACHE INTERNAL "Harfbuzz include dirs" FORCE)

target_include_directories(harfbuzz
    PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}/generated
    ${CMAKE_CURRENT_SOURCE_DIR}/src
    "${CMAKE_CURRENT_SOURCE_DIR}/../freetype2/include/"
    "${CMAKE_CURRENT_SOURCE_DIR}/../freetype2/src/include/"
    ${HARFBUZZ_UNICODE_INCLUDE}
    PRIVATE
    ${CMAKE_CURRENT_SOURCE_DIR}/)

option(USE_INTEL_ATOMIC_PRIMITIVES "Platform needs to define HAVE_INTEL_ATOMIC_PRIMITIVES" FALSE)

if(USE_INTEL_ATOMIC_PRIMITIVES)
    add_definitions(-DHAVE_INTEL_ATOMIC_PRIMITIVES)
endif()

set(HARFBUZZ_LIBRARIES
    freetype
    ${HARFBUZZ_UNICODE_LIB}
    CACHE INTERNAL "Harfbuzz libraries" FORCE)

if (CMAKE_CXX_COMPILER_ID MATCHES "Clang")
  target_compile_options(harfbuzz
    PRIVATE
    -Wno-vla-extension -Wno-unknown-warning-option
    -Wno-extended-offsetof)
endif()

target_compile_definitions(harfbuzz PUBLIC
    PRIVATE
    -DHAVE_OT
    # Provide own config file
    -DHAVE_CONFIG_H
    ${HARFBUZZ_UNICODE_DEFINE}
    -DHB_NO_MT)

    SET_TARGET_PROPERTIES(harfbuzz PROPERTIES
ARCHIVE_OUTPUT_DIRECTORY_DEBUG ${CMAKE_CURRENT_BINARY_DIR}
ARCHIVE_OUTPUT_DIRECTORY_RELEASE ${CMAKE_CURRENT_BINARY_DIR})
SET_TARGET_PROPERTIES(harfbuzz PROPERTIES 
LIBRARY_OUTPUT_DIRECTORY_DEBUG ${CMAKE_CURRENT_BINARY_DIR}
LIBRARY_OUTPUT_DIRECTORY_RELEASE ${CMAKE_CURRENT_BINARY_DIR})
SET_TARGET_PROPERTIES(harfbuzz PROPERTIES
RUNTIME_OUTPUT_DIRECTORY_DEBUG ${CMAKE_CURRENT_BINARY_DIR}
RUNTIME_OUTPUT_DIRECTORY_RELEASE ${CMAKE_CURRENT_BINARY_DIR})

IF(FLYWAVE_ENABLE_SOLUTION_FOLDERS)
  SET_TARGET_PROPERTIES(harfbuzz PROPERTIES FOLDER external)
ENDIF(FLYWAVE_ENABLE_SOLUTION_FOLDERS)
