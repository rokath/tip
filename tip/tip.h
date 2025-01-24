
#include <stdint.h>
#include "tipConfig.h"
#include "tipTable.h"

//! tip is NOT re-entrant or parallel usable, because there are some static objects!
size_t tip( uint8_t* dst, uint8_t const * src, size_t len );

//! tiu is NOT re-entrant or parallel usable, because there are some static objects!
size_t tiu( uint8_t* dst, uint8_t const * src, size_t len );
