//! @file pack.h
//! @brief This is the tip user interface.
//! @author thomas.hoehenleitner [at] seerose.net

#ifndef PACK_H_
#define PACK_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>
#include <stddef.h>

size_t tip( uint8_t* dst, const uint8_t * src, size_t len );
size_t tiPack( uint8_t * dst, const uint8_t * table, const uint8_t * src, size_t slen );
size_t buildTiPacket(uint8_t * dst, uint8_t * dstLimit, const uint8_t * table, const uint8_t * src, size_t slen);

#ifdef __cplusplus
}
#endif

#endif // PACK_H_
