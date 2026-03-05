cmake_minimum_required(VERSION 3.20)
project(ahd_illustrator_plugin_bridge LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

option(AHD_ILLUSTRATOR_USE_SDK "Enable Illustrator SDK integration hooks" OFF)
set(ILLUSTRATOR_SDK_ROOT "" CACHE PATH "Optional path to the Adobe Illustrator SDK")

add_library(ahd_illustrator_plugin_bridge SHARED
  src/BridgeHandlers.cpp
  src/PluginMain.cpp
)

target_include_directories(ahd_illustrator_plugin_bridge
  PRIVATE
    ${CMAKE_CURRENT_SOURCE_DIR}/src
)

target_compile_definitions(ahd_illustrator_plugin_bridge
  PRIVATE
    AHD_PLUGIN_NAME="AHDIllustrator"
)

if(AHD_ILLUSTRATOR_USE_SDK)
  if(NOT ILLUSTRATOR_SDK_ROOT)
    message(FATAL_ERROR "Set ILLUSTRATOR_SDK_ROOT when AHD_ILLUSTRATOR_USE_SDK=ON")
  endif()
  target_compile_definitions(ahd_illustrator_plugin_bridge PRIVATE AHD_ILLUSTRATOR_USE_SDK=1)
  target_include_directories(ahd_illustrator_plugin_bridge PRIVATE ${ILLUSTRATOR_SDK_ROOT})
endif()
