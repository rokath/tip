//! @file ti_unpack.h
//! @brief This is the tip user interface.
//! @author thomas.hoehenleitner [at] seerose.net

#ifndef TI_UNPACK_H_
#define TI_UNPACK_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>
#include <stddef.h>

//! @brief tiu decodes src buffer with size len into dst buffer and returns decoded len.
size_t tiu( uint8_t* dst, const uint8_t * src, size_t len );

//! @brief tiUnpack decodes src buffer with size slen into dst buffer and returns decoded dlen.
//! @details For the tip decoding it uses the passed idTable object.
size_t tiUnpack( uint8_t* dst, const uint8_t * table, const uint8_t * src, size_t slen );

#ifdef __cplusplus
}
#endif

#endif // TI_UNPACK_H_
