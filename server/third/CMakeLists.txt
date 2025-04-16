cmake_minimum_required(VERSION 3.10)
project(HikvisionProject)

# 设置 C++ 标准
set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)
add_compile_options(-DLIBUS_NO_SSL -O3)

# 引入 FetchContent
include(FetchContent)

# 下载 uWebSockets
FetchContent_Declare(
    uWebSockets
    GIT_REPOSITORY https://github.com/uNetworking/uWebSockets.git
    GIT_TAG v20.74.0
)
FetchContent_MakeAvailable(uWebSockets)

# 设置包含目录
include_directories(
    ${PROJECT_SOURCE_DIR}/src
    ${PROJECT_SOURCE_DIR}/hik/include
    ${PROJECT_SOURCE_DIR}/build/_deps/uwebsockets-src/src
    ${PROJECT_SOURCE_DIR}/build/_deps/uwebsockets-src/src/uSockets/src
    ${PROJECT_SOURCE_DIR}/build/_deps/uwebsockets-src/src/examples
    ${PROJECT_SOURCE_DIR}/build/_deps/uwebsockets-src/examples
    ${PROJECT_SOURCE_DIR}/build/_deps/uwebsockets-src/examples/helpers
    ${PROJECT_SOURCE_DIR}/build/_deps/uwebsockets-src/uSockets/src
)

# 调试：打印路径
message(STATUS "uWebSockets_SOURCE_DIR: ${uWebSockets_SOURCE_DIR}")
message(STATUS "USOCKETS_LIBRARY: ${USOCKETS_LIBRARY}")

# 源代码文件
set(SOURCE_FILES
    main.cpp
    ${PROJECT_SOURCE_DIR}/src/server.cpp
)

# 创建可执行文件
add_executable(HikvisionApp ${SOURCE_FILES})

# 查找 libuv
# 查找 libuv
find_library(UV_LIBRARY NAMES uv libuv)
if(NOT UV_LIBRARY)
    message(FATAL_ERROR "libuv not found. Please install libuv1-dev.")
endif()

# 查找 OpenSSL
find_library(SSL_LIBRARY NAMES ssl libssl)
if(NOT SSL_LIBRARY)
    message(FATAL_ERROR "OpenSSL not found. Please install libssl-dev.")
endif()

find_library(CRYPTO_LIBRARY NAMES crypto libcrypto)
if(NOT CRYPTO_LIBRARY)
    message(FATAL_ERROR "OpenSSL crypto not found. Please install libssl-dev.")
endif()

# 查找 zlib
find_library(ZLIB_LIBRARY NAMES z zlib)
if(NOT ZLIB_LIBRARY)
    message(FATAL_ERROR "zlib not found. Please install zlib1g-dev.")
endif()

# 设置 uSockets 库路径
set(USOCKETS_LIBRARY /home/zks/Projects/private/car-pai/server/third/build/_deps/uwebsockets-src/uSockets/uSockets.a)

# 链接 uSockets 和 libuv
target_link_libraries(HikvisionApp
${USOCKETS_LIBRARY}
${UV_LIBRARY}
${SSL_LIBRARY}
${CRYPTO_LIBRARY}
${ZLIB_LIBRARY}
)