
#include <stddef.h>
#include <stdint.h>

//! @brief shift86bit transforms slen 8-bit bytes in src to 7-bit units.
//! @param src is the bytes source buffer.
//! @param slen is the 8-bit byte count.
//! @param lst is the last address inside the dst buffer.
//! @retval is count of 6-bit bytes after operation. 
//! @details The dst buffer is filled from the end. That allows to do an in-buffer conversion.
//! The destination address is computable afterwards: dst = lim - retval.
//! lst is allowed to be "close" behind buf + slen, thus making in-place conversion possible.
//! Example: slen=17, lst=src+24-1
//!      (src)<---              slen=17                    --->(u8)
//! slen=17: b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 __ __ __ __ __ __ __ __ __ __ __ 
//! ret =23: __ __ __ __ __ m6 b6 b6 m6 b6 b6 b6 m6 b6 b6 b6 m6 b6 b6 b6 m6 b6 b6 b6 m6 b6 b6 b6
//!                        (dst) <---                      ret=23                       --->(lst)
//! In dst all MSBits 7&6 are set to 1, to avoid any zeroes.
//! The data are processed from the end.
size_t shift86bit( uint8_t * lst, const uint8_t * src, size_t slen ){
    const uint8_t * u8 = src + slen; // first address behind src buffer
    uint8_t * dst = lst;             // destination address
    while( src < u8 ){
        uint8_t msb = 0xc0;
        for( int i = 1; i < 4; i++ ){
            u8--;                    // next value address
            uint8_t ms = 0xc0 & *u8; // most significant bits 7&6 
            msb |= ms >> (2*(4-i)); // Store most significant bits 7&6 at bit position 5&4, 3&2, 1&0
            *dst-- = 0xc0 | *u8; // Copy to the end: the last byte 6 LSBs and set MSB 7&6 
            if(src == u8){
                break;
            }
        }
        *dst-- = msb;
        msb = 0xc0;
    }
    return lst - dst;
}

