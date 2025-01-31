//! @file unpack.c
//! @brief This is the tip unpack code. Works also without pack.c.
//! @details todo
//! @author thomas.hoehenleitner [at] seerose.net

#include "tipInternal.h"

static size_t shift78bit( uint8_t * buf, size_t len, size_t limit );

//! @brief tiu decodes src buffer with size len into dst buffer and returns decoded len.
size_t tiu( uint8_t * dst, const uint8_t * src, size_t len ){
    if( len < 16 ){
        memcpy(dst, src, len);
        return len;
    }

    // todo
    return shift78bit( dst, len, src-dst ); // dummy, todo
}

//! shift78bit transforms len 7-bit bytes in buf to 8-bit units.
//! @param buf is a byte buffer. It is destroyed during operation.
//! @param len is the 7-bit byte count.
//! @param limit is the max byte count fitting into buf (limit >= len)
//! @retval is count 8-bit bytes
//! @details buf is filled from the end (=buf+limit)
//! The destination is computable afterwards: dst = buf + limit - retval.
//! Example: len=20, limit=24
//!       (buf) <---              20                             --->  [n7]        [n8]
//! len=20: m7 b7 b7 b7 m7 b7 b7 b7 b7 b7 b7 b7 m7 b7 b7 b7 b7 b7 b7 b7 __ __ __ __ 
//! ret=17: __ __ __ __ __ __ __ b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8
//!                            (dst) <---              17                     --->
static size_t shift78bit( uint8_t * buf, size_t len, size_t limit ){
    //int n7 = len; // n7 data index limit.
    //for( int n8 = limit; n8 > 0; ){ // n8 is buf index limit
    //    uint8_t msb = 0x7f & buf[n7-8];
    //    for( int i = 7; i > 0; i-- && n7 > 0 ){
    //        uint8_t m = (msb>>i)<<8;
    //        buf[--n8] = m | buf[--n7]; // the last byte 7 LSBs and MSB=1 to the end
    //    }
    //}
    return 123;
}
