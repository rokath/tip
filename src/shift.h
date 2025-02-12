//! @file shift.h
//! @brief This is the shift interface.

#include <stdint.h>
#include <stddef.h>

size_t shift87bit( uint8_t* lst, const uint8_t * src, size_t len );
size_t shift78bit( uint8_t * dst, const uint8_t * src, size_t slen );