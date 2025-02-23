//! @file unpack.c
//! @brief This is the tip unpack code. Works also without pack.c.
//! @details todo
//! @author thomas.hoehenleitner [at] seerose.net

#include <string.h>
#include "unpack.h"
#include "tip.h"

static size_t collectU7Bytes( uint8_t * dst, const uint8_t * src, size_t slen );
/*static*/ size_t shift78bit( uint8_t * dst, const uint8_t * src, size_t slen );
static size_t restorePacket( uint8_t * dst, const uint8_t * table, const uint8_t * u8, size_t u8len, const uint8_t * src, size_t slen );
static size_t getPatternFromId( uint8_t * pt, const uint8_t * table, uint8_t id );

//! @brief tiu decodes src buffer with size len into dst buffer and returns decoded len.
size_t tiu( uint8_t * dst, const uint8_t * src, size_t slen ){
    return tiUnpack(dst, idTable, src, slen );
}

//! @brief tiUnpack decodes src buffer with size slen into dst buffer and returns decoded dlen.
//! @details For the tip decoding it uses the passed idTable object.
size_t tiUnpack( uint8_t* dst, const uint8_t * table, const uint8_t * src, size_t slen ){
    static uint8_t u78[TIP_SRC_BUFFER_SIZE_MAX*8u/7u+1]; // todo
    size_t u7len = collectU7Bytes( u78, src, slen );

    static uint8_t u8[TIP_SRC_BUFFER_SIZE_MAX]; // todo
    size_t u8len = shift78bit( u8, u78, u7len );

    size_t dlen = restorePacket( dst, table, u8, u8len, src, slen );
    return dlen;
}

// collectU7Bytes copies all bytes with msbit=1 into dst and returns their count.
static size_t collectU7Bytes( uint8_t * dst, const uint8_t * src, size_t slen ){
    uint8_t * p = dst;
    for( int i = 0; i < slen; i++ ){
        if(0x80 & src[i]){
            *p++ = src[i];
        }
    }
    return p - dst;
}

//! shift78bit transforms slen 7-bit bytes in src to 8-bit units in dst.
//! @param src is a byte buffer.
//! @param slen is the 7-bit byte count.
//! @param dst is the destination buffer. It is NOT allowed to be equal src for in-place conversion.
//! @retval is count 8-bit bytes
//! @details buf is filled from the end (=buf+limit)
//! Example: slen=20, limit=24
//!       (src)<---               slen=20                       --->(lst)     
//! slen=20: m7 b7 b7 b7 m7 b7 b7 b7 b7 b7 b7 b7 m7 b7 b7 b7 b7 b7 b7 b7
//! ret =17: b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8
//!       (dst)<---               dlen=17               --->(ptr)
//! dlen = slen*7/8
/*static*/ size_t shift78bit( uint8_t * dst, const uint8_t * src, size_t slen ){
    size_t dlen = (7*slen)>>3;
    uint8_t * ptr = dst + dlen - 1; // ptr is last address in dst buffer
    uint8_t * lst = (uint8_t *)src + slen - 1; // lst is last address in source buffer.

    while( src <= lst - 7 ){
        uint8_t msbyte = 0x7f & *(lst-7); // remove 1 in msb _100 0000 == 0x40
        for( int i = 0; i < 7; i++ ){ 
            uint8_t bits6_0 = 0x7f & *lst--; // _111 1111 == 0x7f
            uint8_t mask = 0x40 >> i;        // _100 0000
            uint8_t b7bit = msbyte & mask;   // _100 0000 & _100 0000 == 0x40
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

//! restorePacket reconstructs original data using src, slen, u8, u8len and table into dst and returns the count.
static size_t restorePacket( uint8_t * dst, const uint8_t * table, const uint8_t * u8, size_t u8len, const uint8_t * src, size_t slen ){
    uint8_t * p = dst;
    for( int i = 0; i < slen; i++ ){
        if( 0x80 & src[i] ){ // an u78 byte
            if( u8len > 0){
                *p++ = *u8++;
                u8len--;
            }
        }else{ // an id
            size_t sz = getPatternFromId( p, table, src[i] );
            p += sz;
        }
    }
    return p - dst;
}

//! getPatternFromId seaches in testTable for id.
//! @param pt is filled with the replace pattern if id was found.
//! @param table is the pattern table.
//! @param id is the replace byte. Valid values for id are 1...0x7f.
//! @retval is the pattern size or 0, if id was not found.
static size_t getPatternFromId( uint8_t * pt, const uint8_t * table, uint8_t id ){
    size_t sz;
    unsigned int idx = 0x01;
    while( (sz = *table++) && sz){  // a pattern exists here
        if( idx == id ){ // id found
            memcpy(pt, table, sz);
            return sz;
        }
        idx++;
        table += sz;
    }
    return 0;
}
