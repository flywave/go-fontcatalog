cmake_minimum_required(VERSION 3.10)

project(msdfgen VERSION 1.8 LANGUAGES CXX)
option(MSDFGEN_BUILD_MSDFGEN_STANDALONE "Build the msdfgen standalone executable" ON)
option(MSDFGEN_USE_OPENMP "Build with OpenMP support for multithreaded code" OFF)
option(MSDFGEN_USE_CPP11 "Build with C++11 enabled" ON)
option(MSDFGEN_USE_SKIA "Build with the Skia library" OFF)
option(MSDFGEN_INSTALL "Generate installation target" ON)
option(FREETYPE_WITH_PNG "Link libpng and zlib because FreeType is configured to require it" OFF)
option(FREETYPE_WITH_HARFBUZZ "Link HarfBuzz because FreeType is configured to require it" OFF)

list(APPEND CMAKE_MODULE_PATH "${CMAKE_CURRENT_SOURCE_DIR}/cmake")

INCLUDE_DIRECTORIES(
    "${CMAKE_CURRENT_SOURCE_DIR}"
    "${CMAKE_CURRENT_SOURCE_DIR}/include"
    "${CMAKE_CURRENT_SOURCE_DIR}/../libpng"
    "${CMAKE_CURRENT_SOURCE_DIR}/../zlib"
    "${CMAKE_CURRENT_BINARY_DIR}/../zlib"
    "${CMAKE_CURRENT_SOURCE_DIR}/../freetype2/include/"
    "${CMAKE_CURRENT_SOURCE_DIR}/../harfbuzz/include/"
)

#----------------------------------------------------------------
# Gathering File
#----------------------------------------------------------------

file(GLOB_RECURSE msdfgen_HEADERS RELATIVE ${CMAKE_CURRENT_SOURCE_DIR}
	"core/*.h"
	"core/*.hpp"
)

file(GLOB_RECURSE msdfgen_SOURCES RELATIVE ${CMAKE_CURRENT_SOURCE_DIR}
	"core/*.cpp"
)

file(GLOB_RECURSE msdfgen-ext_PUBLIC_HEADERS RELATIVE ${CMAKE_CURRENT_SOURCE_DIR}
	"ext/*.h"
	"ext/*.hpp"
)

file(GLOB_RECURSE msdfgen-ext_PRIVATE_HEADERS RELATIVE ${CMAKE_CURRENT_SOURCE_DIR}
	"include/*.h"
)

file(GLOB_RECURSE msdfgen-ext_SOURCES RELATIVE ${CMAKE_CURRENT_SOURCE_DIR}
	"ext/*.cpp"
	"lib/*.cpp"
)

# Build the library (aliased name because it's the same target name the exe)
include(folderize)
folderize_sources(msdfgen_HEADERS ${CMAKE_SOURCE_DIR})
folderize_sources(msdfgen_SOURCES ${CMAKE_SOURCE_DIR})
folderize_sources(msdfgen-ext_PUBLIC_HEADERS ${CMAKE_SOURCE_DIR})
folderize_sources(msdfgen-ext_PRIVATE_HEADERS ${CMAKE_SOURCE_DIR})
folderize_sources(msdfgen-ext_SOURCES ${CMAKE_SOURCE_DIR})

#----------------------------------------------------------------
# Target configuration
#----------------------------------------------------------------

add_library(msdfgen STATIC ${msdfgen_SOURCES} ${msdfgen_HEADERS} "./msdfgen.h")
set_target_properties(msdfgen PROPERTIES PUBLIC_HEADER "${msdfgen_HEADERS}")
target_include_directories(msdfgen INTERFACE
	$<INSTALL_INTERFACE:include>
	$<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/>
)

if(MSDFGEN_USE_CPP11)
	target_compile_features(msdfgen PUBLIC cxx_std_11)
	target_compile_definitions(msdfgen PUBLIC MSDFGEN_USE_CPP11)
endif()

if(MSDFGEN_USE_OPENMP)
	# Note: Clang doesn't support OpenMP by default...
	find_package(OpenMP REQUIRED COMPONENTS CXX)
	target_link_libraries(msdfgen PRIVATE OpenMP::OpenMP_CXX)
	target_compile_definitions(msdfgen PRIVATE MSDFGEN_USE_OPENMP)
endif()

add_library(msdfgen_ext STATIC ${msdfgen-ext_SOURCES} ${msdfgen-ext_PUBLIC_HEADERS} ${msdfgen-ext_PRIVATE_HEADERS} "./msdfgen-ext.h")
target_include_directories(msdfgen_ext INTERFACE
	$<INSTALL_INTERFACE:include>
	$<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/>
)

if(MSDFGEN_USE_CPP11)
	target_compile_features(msdfgen_ext PUBLIC cxx_std_11)
	target_compile_definitions(msdfgen_ext PUBLIC MSDFGEN_USE_CPP11)
endif()

SET_TARGET_PROPERTIES(msdfgen PROPERTIES
ARCHIVE_OUTPUT_DIRECTORY_DEBUG ${CMAKE_CURRENT_BINARY_DIR}
ARCHIVE_OUTPUT_DIRECTORY_RELEASE ${CMAKE_CURRENT_BINARY_DIR})
SET_TARGET_PROPERTIES(msdfgen PROPERTIES 
LIBRARY_OUTPUT_DIRECTORY_DEBUG ${CMAKE_CURRENT_BINARY_DIR}
LIBRARY_OUTPUT_DIRECTORY_RELEASE ${CMAKE_CURRENT_BINARY_DIR})
SET_TARGET_PROPERTIES(msdfgen PROPERTIES
RUNTIME_OUTPUT_DIRECTORY_DEBUG ${CMAKE_CURRENT_BINARY_DIR}
RUNTIME_OUTPUT_DIRECTORY_RELEASE ${CMAKE_CURRENT_BINARY_DIR})

SET_TARGET_PROPERTIES(msdfgen_ext PROPERTIES
ARCHIVE_OUTPUT_DIRECTORY_DEBUG ${CMAKE_CURRENT_BINARY_DIR}
ARCHIVE_OUTPUT_DIRECTORY_RELEASE ${CMAKE_CURRENT_BINARY_DIR})
SET_TARGET_PROPERTIES(msdfgen_ext PROPERTIES 
LIBRARY_OUTPUT_DIRECTORY_DEBUG ${CMAKE_CURRENT_BINARY_DIR}
LIBRARY_OUTPUT_DIRECTORY_RELEASE ${CMAKE_CURRENT_BINARY_DIR})
SET_TARGET_PROPERTIES(msdfgen_ext PROPERTIES
RUNTIME_OUTPUT_DIRECTORY_DEBUG ${CMAKE_CURRENT_BINARY_DIR}
RUNTIME_OUTPUT_DIRECTORY_RELEASE ${CMAKE_CURRENT_BINARY_DIR})
