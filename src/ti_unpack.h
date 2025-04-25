//! @file ti_unpack.h
//! @brief This is the tip user interface.

#ifndef TI_UNPACK_H_
#define TI_UNPACK_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>
#include <stddef.h>

//! @brief tiu decodes src buffer with size len into dst buffer and returns decoded len.
//! @details For the decoding it uses the linked idTable.c object.
size_t tiu( uint8_t* dst, const uint8_t * src, size_t len );

#ifdef __cplusplus
}
#endif

#endif // TI_UNPACK_H_
