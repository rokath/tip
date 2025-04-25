//! @file tipDefaultConfig.h
//! @brief This is the default tip configuration.

#ifndef TIP_DEFAULT_CONFIG_H_
#define TIP_DEFAULT_CONFIG_H_

#ifdef __cplusplus
extern "C" {
#endif

#ifndef TIP_SRC_BUFFER_SIZE_MAX
//! TIP_SRC_BUFFER_SIZE_MAX is the maximun allowed input data buffer length.
//! Its size has influence on the statically allocated RAM.
//! Must be a multiple of 8.
#define TIP_SRC_BUFFER_SIZE_MAX 104u // bytes (max 65528u)
#endif

#ifndef TIP_MAX_PATH_COUNT
//! TIP_MAX_PATH_COUNT is the max allowed path count.
//! Its size has influence on the statically allocated RAM.
//! In a first run on each src buffer position a 2-byte pattern could match.
//! In a first run on the first src buffer position all IDPosition table entries coud match.
//! TODO: Find right logic.
#define TIP_MAX_PATH_COUNT (TIP_SRC_BUFFER_SIZE_MAX/4)
#endif

//  #ifndef TIP_VERBOSE
//  #define TIP_VERBOSE 0
//  #endif
//  
//  #ifndef TIP_DEBUG
//  #define TIP_DEBUG 0
//  #endif
//  
//  #ifndef OPTIMIZE_UNREPLACABLES
//  //! OPTIMIZE_UNREPLACABLES allows to reduce the TiP packet size in some special cases.
//  //! It is an option just for tests and should be enabled always.
//  #define OPTIMIZE_UNREPLACABLES 1
//  #endif

#ifdef __cplusplus
}
#endif

#endif // TIP_DEFAULT_CONFIG_H_
