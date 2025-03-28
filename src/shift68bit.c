
#include <stdint.h>

//! shift68bit transforms slen 6-bit bytes in src to 8-bit units in dst.
//! @param src is a byte buffer.
//! @param slen is the 6-bit byte count.
//! @param dst is the destination buffer. It is allowed to be equal src for in-place conversion.
//! @retval is count 8-bit bytes
//! Example: slen=7
//!      (src)<--  slen=7   -->  
//! slen=7: m6 b6 b6 m6 b6 b6 b6
//! ret =5: b8 b8 b8 b8 b8 
size_t shift68bit( uint8_t * dst, const uint8_t * src, size_t slen ){
    uint8_t * pb8 = dst;
    const uint8_t * slim = src + slen;

    // m6len is alway 3, but the very first one can be 1 or 2 or 3:
    // When slen is            2 3 4  6 7 8  a b c  e f 0 ...
    // then dlen is            1 2 3  4 5 6  7 8 9  a b c ... 
    // The first m6 byte count 1 2 3  1 2 3  1 2 3  1 2 3 ... (m6len)
    // The slen 2 lsbits are   2 3 0  2 3 0  2 3 0  2 3 0 ... (slen2lsb)
    uint8_t slen2lsb = slen & 3;
    uint8_t m6len = slen2lsb ? slen2lsb - 1 : 3;

    while( src < slim){
        uint8_t m6 = *src++;            // c1 == 11.00_00_01
        for( int i=2*m6len; i>0; i-=2){
            uint8_t b8 = 0x3f & *src++;
            uint8_t mask = 3 << (i-2); // 00110000 00001100 00000011 (i:6,4,2)
            uint8_t cmp =  mask & m6;  // 00cc0000 0000cc00 000000cc (i:6,4,2)
            cmp <<= (8-i); // cc000000:   <<= 2    <<= 4    <<=6
            b8 |= cmp;
            *pb8++ = b8;
        }
        m6len = 3; 
    }
    return pb8 - dst;
}
