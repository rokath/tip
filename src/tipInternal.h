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

//! id_t could be uint8_t to safe memory, when no indirect IDs are used.
typedef uint16_t id_t; 

//! UnreplacableContainerBits is an in idTable.c generated value.
extern const unsigned unreplacableContainerBits; 

//! ID1Max is an in idTable.c generated value.
extern const unsigned ID1Max;

//! ID1Count is an in idTable.c generated value.
extern const unsigned ID1Count;

//! MaxID is an in idTable.c generated value.
extern const unsigned MaxID;

//! LastID is an in idTable.c generated value.
extern const unsigned LastID;

//! idTable is an in idTable.c generated value.
extern const uint8_t idTable[];

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
