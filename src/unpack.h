//! @file unpack.h
//! @brief This is the tip user interface.
//! @author thomas.hoehenleitner [at] seerose.net

#ifndef UNPACK_H_
#define UNPACK_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>
#include <stddef.h>

size_t tiu( uint8_t* dst, const uint8_t * src, size_t len );
size_t tiUnpack( uint8_t* dst, const uint8_t * table, const uint8_t * src, size_t slen );

#ifdef __cplusplus
}
#endif

#endif // UNPACK_H_
