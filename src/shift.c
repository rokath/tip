//! @file shift.c
//! @author thomas.hoehenleitner [at] seerose.net

#include "shift.h"

//! shift87bit transforms slen 8-bit bytes in src to 7-bit units.
//! @param src is the bytes source buffer.
//! @param slen is the 8-bit byte count.
//! @param lst is the last address inside the dst buffer.
//! @retval is count of 7-bit bytes after operation. 
//! @details The dst buffer is filled from the end.Thas allows to do an in-buffer conversion.
//! The destination address is computable afterwards: dst = lim - retval.
//! lim is allowed to be "close" behind buf + slen, thus making in-place conversion possible.
//! Example: slen=17, limit=24
//!       (src) <---            slen=17                   --->(u8)
//! slen=17: b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 __ __ __ __ __ __ __
//! ret =20: __ __ __ __ m7 b7 b7 b7 m7 b7 b7 b7 b7 b7 b7 b7 m7 b7 b7 b7 b7 b7 b7 b7
//!                   (dst) <---                ret=20                       --->(lst)
//! In dst all MSBits are set to 1, to avoid any zeroes.
//! The data are processed from the end.
size_t shift87bit( uint8_t* lst, uint8_t * const src, size_t slen ){
    uint8_t * u8 = src + slen; // first address behind src buffer
    uint8_t * dst = lst; // destination address
    while( src < u8 ){
        uint8_t msb = 0x80;
        for( int i = 1; i < 8; i++ ){
            u8--; // next value address
            uint8_t ms = 0x80 & *u8; // most significant bit                i     12345678
            msb |= ms >> i; // Store most significant bit at bit position:  8 -> _76543210 
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


//! shift78bit transforms slen 7-bit bytes in src to 8-bit units in dst.
//! @param src is a byte buffer.
//! @param slen is the 7-bit byte count.
//! @param dst is the destination buffer. It is allowed to be equal src.
//! @retval is count 8-bit bytes
//! @details buf is filled from the end (=buf+limit)
//! Example: slen=20, limit=24
//!       (src)<---               slen=20                       --->(lst)     
//! slen=20: m7 b7 b7 b7 m7 b7 b7 b7 b7 b7 b7 b7 m7 b7 b7 b7 b7 b7 b7 b7
//! ret =17: b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8
//!       (dst)<---               dlen=17               --->(ptr)
//! dlen = slen*7/8
size_t shift78bit( uint8_t * dst, const uint8_t * src, size_t slen ){
    size_t dlen = (7*slen)>>3;
    uint8_t * ptr = dst + dlen - 1; // ptr is last address in dst buffer
    uint8_t * lst = (uint8_t *)src + slen - 1; // lst is last address in source buffer.

    while( src <= lst - 7 ){
        uint8_t msbyte = 0x7f & *(lst-7); // remove 1 in msb _100 0000 == 0x40
        for( int i = 0; i < 7; i++ ){ 
            uint8_t bits6_0 = 0x7f & *lst--; // _111 1111 == 0x7f
            uint8_t mask = 0x40 >> i;        // _100 0000
            uint8_t b7bit = msbyte & mask; // _100 0000 & _100 0000 == 0x40
            b7bit = b7bit ? 0x80 : 0;
            *ptr-- = b7bit | bits6_0;
        }
        lst--; // Skip over already processed msbyte.
    }
    if( lst <= src){
        return dlen;
    }
    // Now we have one msbyte and 1-6 b7 bytes left.
    uint8_t msbyte = 0x7f & *src;
    size_t cnt = lst - src; // cnt of remaining 1-6 b7 bytes
    for( int i = 0; i < cnt; i++ ){ 
        uint8_t bits6_0 = 0x7f & *lst--; // _111 1111 == 0x7f
        uint8_t mask = 0x40 >> i;        // _100 0000
        uint8_t b7bit = msbyte & mask; // _100 0000 & _100 0000 == 0x40
        b7bit = b7bit ? 0x80 : 0;
        *ptr-- = b7bit | bits6_0;
     }
    return dlen;
}
