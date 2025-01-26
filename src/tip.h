//! @file tip.h
//! @brief This is the tip user interface.
//! @author thomas.hoehenleitner [at] seerose.net

#ifndef TIP_H_
#define TIP_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>

//! TiP is NOT re-entrant or parallel usable, because there are some static objects!
size_t TiP( uint8_t* dst, uint8_t const * src, size_t len );

//! TiU is NOT re-entrant or parallel usable, because there are some static objects!
size_t TiU( uint8_t* dst, uint8_t const * src, size_t len );

#ifdef __cplusplus
}
#endif

#endif // TIP_H_
