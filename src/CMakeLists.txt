
include_directories(${FLYWAVE_INCLUDE_DIRS})

FILE( GLOB fontcatalog_SOURCE_FILES ${CMAKE_CURRENT_SOURCE_DIR}/*.cc )
FILE( GLOB fontcatalog_HEADER_FILES ${CMAKE_CURRENT_SOURCE_DIR}/*.hh ${CMAKE_CURRENT_SOURCE_DIR}/*.h )

ADD_LIBRARY(fontcatalog STATIC
            ${fontcatalog_HEADER_FILES}
            ${fontcatalog_SOURCE_FILES}
          )

TARGET_INCLUDE_DIRECTORIES(fontcatalog PUBLIC $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/..>)

SOURCE_GROUP("src" FILES ${fontcatalog_SOURCE_FILES})
SOURCE_GROUP("include" FILES ${fontcatalog_HEADER_FILES})

SET_TARGET_PROPERTIES(fontcatalog
                    PROPERTIES PUBLIC_HEADER "${headers}")

add_executable(tests fontcatalog_test.cc)

IF(UNIX AND NOT APPLE)
TARGET_LINK_LIBRARIES(tests -Wl,--start-group fontcatalog ${FLYWAVE_LIBRARY_DEPES} -Wl,--end-group)
ELSE()
TARGET_LINK_LIBRARIES(tests fontcatalog ${FLYWAVE_LIBRARY_DEPES})
ENDIF()
