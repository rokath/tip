//! @file tipInternal.h
//! @brief This is the tip internal interface.

#ifndef TIP_INTERNAL_H_
#define TIP_INTERNAL_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>
#include "tipConfig.h" // project specific file

size_t shift87bit( uint8_t * lst, const uint8_t * src, size_t slen );
size_t shift78bit( uint8_t * dst, const uint8_t * src, size_t slen );

size_t shift86bit( uint8_t * lst, const uint8_t * src, size_t slen );
size_t shift68bit( uint8_t * dst, const uint8_t * src, size_t slen );

#ifndef TIP_SRC_BUFFER_SIZE_MAX
#error "needs to #define TIP_SRC_BUFFER_SIZE_MAX 248u // bytes (max65535)"
#endif // #ifndef TIP_SRC_BUFFER_SIZE_MAX

#if TIP_SRC_BUFFER_SIZE_MAX & 7 
#error "needs to be a multiple of 8"
#endif // #if TIP_SRC_BUFFER_SIZE_MAX & 7 

#if TIP_SRC_BUFFER_SIZE_MAX > 256u*1024u*1024u
#error invalid TIP_SRC_BUFFER_SIZE_MAX value
#elif TIP_SRC_BUFFER_SIZE_MAX > 0xfff8u
typedef uint32_t loc_t;
#elif TIP_SRC_BUFFER_SIZE_MAX > 0xf8u
typedef uint16_t loc_t;
#else
typedef uint8_t loc_t;
#endif

#if UNREPLACABLE_BIT_COUNT == 7
#define UNREPLACABLE_MASK 0x80
#define DIRECT_ID_MAX 127
#else
#define UNREPLACABLE_MASK 0xC0
#define DIRECT_ID_MAX 191
#endif

///////////////////////////////////////////////////////////////////////////////////////////////////
// Exported for target tests
//

//
///////////////////////////////////////////////////////////////////////////////////////////////////

#include "ti_pack.h"
#include "ti_unpack.h"

extern const uint8_t idTable[];

#ifdef __cplusplus
}
#endif

#endif // TIP_INTERNAL_H_
