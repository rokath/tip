//! @file ti_unpack.c
//! @brief This is the tip unpack code. Works also without pack.c.
//! @details todo
//! @author thomas.hoehenleitner [at] seerose.net

#include <string.h>
#include "ti_unpack.h"
#include "tipInternal.h"

static int collectUTBytes( uint8_t * dst, const uint8_t * src, size_t slen );
/*static*/ size_t shift78bit( uint8_t * dst, const uint8_t * src, size_t slen );
static size_t reconvertBits( uint8_t * lst, const uint8_t * src, size_t slen );
static size_t restorePacket( uint8_t * dst, const uint8_t * table, const uint8_t * u8, size_t u8len, const uint8_t * src, size_t slen );
static size_t getPatternFromId( uint8_t * pt, const uint8_t * table, uint8_t id );

size_t tiu( uint8_t * dst, const uint8_t * src, size_t slen ){
    return tiUnpack(dst, idTable, src, slen );
}

uint8_t uT8[TIP_SRC_BUFFER_SIZE_MAX*8u/7u+1]; // todo
uint8_t u8[TIP_SRC_BUFFER_SIZE_MAX]; // todo

size_t tiUnpack( uint8_t* dst, const uint8_t * table, const uint8_t * src, size_t slen ){
    int uTlen = collectUTBytes( uT8, src, slen );
    size_t u8len;
#if OPTIMIZE_UNREPLACABLES
    if (uTlen <= 0 ) { // Unrplacable byte optimisation was possible.
        u8len = -uTlen;
        memcpy( u8, uT8, u8len );
    } else { // Otherwise the last byte is an unreplacable and not the only one and there is at least one ID.
        u8len = reconvertBits( u8, uT8, uTlen ); // Optimization was not possible.
    }
#else // #if OPTIMIZE_UNREPLACABLES
    u8len = reconvertBits( u8, uT8, uTlen );
#endif // #else // #if OPTIMIZE_UNREPLACABLES
    size_t dlen = restorePacket( dst, table, u8, u8len, src, slen );
    return dlen;
}

// collectUTBytes copies all bytes with msbit=1 into dst and returns their count.
static int collectUTBytes( uint8_t * dst, const uint8_t * src, size_t slen ){
    uint8_t * p = dst;
    for( int i = 0; i < slen; i++ ){
        if(UNREPLACABLE_MASK & src[i]){
            *p++ = src[i];
        }
    }
    int count = p - dst;
#if OPTIMIZE_UNREPLACABLES // cases like III or IIU or UUIIIUII 
    if ( (count <= 1) // TiP packet has no or max one unrplacable byte.
      || (*(src+slen-1) <= DIRECT_ID_MAX ) ) {// TiP packet ends with an ID.
        count = -count; // Unreplacable bytes optimisation was possible.
    }
#endif // #if OPTIMIZE_UNREPLACABLES
    return count;
}

#if UNREPLACABLE_BIT_COUNT == 7

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
/*static*/ size_t shift78bit( uint8_t * dst, const uint8_t * src, size_t slen ){
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

#endif

#if UNREPLACABLE_BIT_COUNT == 6
//! shift68bit transforms slen 6-bit bytes in src to 8-bit units in dst.
//! @param src is a byte buffer.
//! @param slen is the 6-bit byte count.
//! @param dst is the destination buffer. It is allowed to be equal src for in-place conversion.
//! @retval is count 8-bit bytes
//! Example: slen=7
//!      (src)<--  slen=7   -->  
//! slen=7: m6 b6 b6 m6 b6 b6 b6
//! ret =5: b8 b8 b8 b8 b8 
/*static*/ size_t shift68bit( uint8_t * dst, const uint8_t * src, size_t slen ){
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
#endif


//! reconvertBits transmutes slen n-bit bytes in src to 8-bit units in dst.
//! @param src is a byte buffer.
//! @param slen is the n-bit byte count.
//! @param dst is the destination buffer. It is allowed to be equal src for in-place conversion.
//! @retval is count 8-bit bytes
//! @details buf is filled from the end (=buf+limit)
static size_t reconvertBits( uint8_t * dst, const uint8_t * src, size_t slen ){
    #if UNREPLACABLE_BIT_COUNT == 7
        return shift78bit( dst, src, slen );
    #endif
    #if UNREPLACABLE_BIT_COUNT == 6
        return shift68bit( dst, src, slen );
    #endif
}

//! restorePacket reconstructs original data using src, slen, u8, u8len and table into dst and returns the count.
static size_t restorePacket( uint8_t * dst, const uint8_t * table, const uint8_t * u8, size_t u8len, const uint8_t * src, size_t slen ){
    uint8_t * p = dst;
    for( int i = 0; i < slen; i++ ){
        if( UNREPLACABLE_MASK & src[i] ){ // an uT8 byte
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
