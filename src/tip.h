//! @file tip.h
//! @brief This is the tip user interface.

#include <stdint.h>

//! TiP is NOT re-entrant or parallel usable, because there are some static objects!
size_t TiP( uint8_t* dst, uint8_t const * src, size_t len );

//! TiU is NOT re-entrant or parallel usable, because there are some static objects!
size_t TiU( uint8_t* dst, uint8_t const * src, size_t len );
