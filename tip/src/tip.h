
#include <stdint.h>
#include "tipConfig.h"
#include "tipTable.h"

//! tip is NOT re-entrant or parallel usable, because there are some static objects!
size_t TiP( uint8_t* dst, uint8_t const * src, size_t len );

//! tiu is NOT re-entrant or parallel usable, because there are some static objects!
size_t TiU( uint8_t* dst, uint8_t const * src, size_t len );

//! @brief replacement_t is a replacement type descriptor.
typedef struct {
    offset_t bo; // bo is the buffer offset, where replacement size starts.
    uint8_t  sz; // sz is the replacement size (2-255).
    uint8_t  by; // by is the replacement byte 0x01 to 0xff.
} replacement_t;
