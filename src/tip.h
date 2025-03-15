//! @file tip.h
//! @brief This is the tip configuration and tip internal common interface.
//! @author thomas.hoehenleitner [at] seerose.net

#ifndef TIP_H_
#define TIP_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>
#include "tipConfig.h"

#ifndef TIP_SRC_BUFFER_SIZE_MAX
#error "needs to #define TIP_SRC_BUFFER_SIZE_MAX 248u // bytes (max65535)"
#endif // #ifndef TIP_SRC_BUFFER_SIZE_MAX

#if TIP_SRC_BUFFER_SIZE_MAX & 7 
#error "needs to be a multiple of 8"
#endif // #if TIP_SRC_BUFFER_SIZE_MAX & 7 

#if TIP_SRC_BUFFER_SIZE_MAX > 256u*1024u*1024u
#error invalid TIP_SRC_BUFFER_SIZE_MAX value
#elif TIP_SRC_BUFFER_SIZE_MAX > 0xfff8u
typedef uint32_t offset_t;
#elif TIP_SRC_BUFFER_SIZE_MAX > 0xf8u
typedef uint16_t offset_t;
#else
typedef uint8_t offset_t;
#endif

#include "pack.h"
#include "unpack.h"

extern const uint8_t idTable[];

#ifdef __cplusplus
}
#endif

#endif // TIP_H_
