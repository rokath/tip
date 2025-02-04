//! @file shift.c
//! @author thomas.hoehenleitner [at] seerose.net

#include "shift.h"

//! shift87bit transforms len 8-bit bytes in src to 7-bit units.
//! @param src is the bytes source buffer.
//! @param len is the 8-bit byte count.
//! @param lst is the last address inside the dst buffer.
//! @retval is count of 7-bit bytes after operation. 
//! @details The dst buffer is filled from the end because we do not know its exact size in advance.
//! The destination address is computable afterwards: dst = lim - retval.
//! lim is allowed to be "close" behind buf + len, thus making in-place conversion possible.
//! Example: len=17, limit=24
//!       (src) <---              17                    --->  (u8)              (lst)
//! len=17: b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 __ __ __ __ __ __ __
//! ret=20: __ __ __ __ m7 b7 b7 b7 m7 b7 b7 b7 b7 b7 b7 b7 m7 b7 b7 b7 b7 b7 b7 b7
//!                   (dst) <---                     20                       ---> 
//! In dst all MSBits are set to 1, to avoid any zeroes.
size_t shift87bit( uint8_t* lst, uint8_t * src, size_t len ){
    uint8_t * u8 = src + len;
    uint8_t * dst = lst;
    while( src < u8 ){
        uint8_t msb = 0x80;
        for( int i = 1; i < 8; i++ ){
            u8--;
            msb |= (0x80 & *u8)>>i; // Store the MSB of the current last byte at bit position
            *dst-- = (0x7F & *u8) | 0x80; // the last byte 7 LSBs and set MSB=1 to the end
            if(src == u8){
                break;
            }
        }
        *dst-- = msb;
        msb = 0x80;
    }
    return lst - dst;
}


//! shift78bit transforms len 7-bit bytes in src to 8-bit units in dst.
//! @param src is a byte buffer.
//! @param len is the 7-bit byte count.
//! @param dst is the destination buffer. It is allowed to be equal src.
//! @retval is count 8-bit bytes
//! @details buf is filled from the end (=buf+limit)
//! Example: len=20, limit=24
//!       (src) <---              20                              --->      
//! len=20: m7 b7 b7 b7 m7 b7 b7 b7 b7 b7 b7 b7 m7 b7 b7 b7 b7 b7 b7 b7
//! ret=17: b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8
//!       (dst) <---              17                     --->
size_t shift78bit( uint8_t * dst, uint8_t * src, size_t len ){
    uint8_t * ptr = dst; 
    uint8_t * lst = src + len;
    int msbitCount = (7*len>>3)%7; // See TipUserManual.md for details.
    msbitCount = msbitCount ? msbitCount : 7; // 0 --> 7 
    while( src < lst ){
        uint8_t msbyte = 0x7f & *src++;
        for( int i = 0; i <= 6; i++ ){ 
            uint8_t b87 = 0x7f & *src++; // _010 0000
            uint8_t mask = 0x01 << (8-i); // _100 0000
            uint8_t b8bit = msbyte & mask; // _100 0000 
            b8bit = b8bit ? 0x80 : 0;
            *ptr++ = b8bit | b87;
            msbitCount--; // after getting 0, msbitCount gets negative and that´s fine
            if( !msbitCount || src == lst ){
                break;
            }
        }
    }
    return ptr - dst;
}
