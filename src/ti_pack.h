//! @file ti_pack.h
//! @brief This is the tip user interface for packing.

#ifndef TI_PACK_H_
#define TI_PACK_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>
#include <stddef.h>

//! @brief tip encodes src buffer with size len into dst buffer and returns encoded len.
//! @details For the encoding it uses the linked idTable.c object.
size_t tip( uint8_t* dst, const uint8_t * src, size_t len );

#ifdef __cplusplus
}
#endif

#endif // TI_PACK_H_
