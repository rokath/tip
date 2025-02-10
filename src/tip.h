//! @file tip.h
//! @brief This is the tip user interface.
//! @author thomas.hoehenleitner [at] seerose.net

#ifndef TIP_H_
#define TIP_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>
#include <stddef.h>

//! tip is NOT re-entrant or parallel usable, because there are some static objects!
size_t tip( uint8_t* dst, const uint8_t * src, size_t len );

//! tiu is NOT re-entrant or parallel usable, because there are some static objects!
size_t tiu( uint8_t* dst, const uint8_t * src, size_t len );

#ifdef __cplusplus
}
#endif

#endif // TIP_H_
