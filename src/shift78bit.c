
#include <stdint.h>

//! shift78bit transforms slen 7-bit bytes in src to 8-bit units in dst.
//! @param src is a byte buffer.
//! @param slen is the 7-bit byte count.
//! @param dst is the destination buffer. It is allowed to be equal src for in-place conversion.
//! @retval is count 8-bit bytes
//! Example: slen=7
//!       (src)<--       slen=11        -->  
//! slen=11: M7 A7 B7 m7 b7 b7 b7 b7 b7 b7 b7
//! ret = 9: A8 B8 b8 b8 b8 b8 b8 b8 b8
//! M7 == 0b100000AB, A is the msb of A8 and B is the msb of B8
size_t shift78bit( uint8_t * dst, const uint8_t * src, size_t slen ){
    uint8_t * pb8 = dst;
    const uint8_t * slim = src + slen;

    // m7len is alway 7, but the very first one can be 1 - 7:
    // When slen is            2 3 4 5 6 7 8  a b c d e f 0 ...
    // then dlen is            1 2 3 4 5 6 7  8 9 a b c d e ... 
    // The first m7 byte count 1 2 3 4 5 6 7  1 2 3 4 5 6 7 ... (m7len)
    // The slen 3 lsbits are   2 3 4 5 6 7 0  2 3 4 5 6 7 0 ... (m7len+1)
    uint8_t slen3lsb = slen & 7;
    uint8_t m7len = slen3lsb ? slen3lsb - 1 : 7;

    while( src < slim){
        uint8_t m7 = *src++; 
        for( int i=m7len; i>0; i--){
            uint8_t b8 = 0x7f & *src++; // get 7 lsb
            uint8_t mask = 1 << (i-1);
            uint8_t msb = mask & m7 ? 0x80 : 0;
            b8 |= msb;
            *pb8++ = b8;
        }
        m7len = 7; 
    }
    return pb8 - dst;
}
