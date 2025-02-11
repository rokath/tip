//! @file unpack.c
//! @brief This is the tip unpack code. Works also without pack.c.
//! @details todo
//! @author thomas.hoehenleitner [at] seerose.net

#include "tipInternal.h"

//! @brief tiu decodes src buffer with size len into dst buffer and returns decoded len.
size_t tiu( uint8_t * dst, const uint8_t * src, size_t slen ){
    return tiUnpack(dst, idTable, src, slen );
}

//! @brief tiUnpack decodes src buffer with size slen into dst buffer and returns decoded dlen.
//! @details For the tip decoding it uses the passed idTable object.
size_t tiUnpack( uint8_t* dst, const uint8_t * table, const uint8_t * src, size_t slen ){
    if( slen < 6 ){
        memcpy(dst, src, slen);
        return slen;
    }
    static uint8_t u7[256]; // todo
    size_t u7len = collectU7Bytes( u7, src, slen );

    static uint8_t u8[256]; // todo
    size_t u8len = shift78bit( u8, u7, u7len );

    size_t dlen = restorePacket( dst, table, u8, u8len, src, slen );
    return dlen;
}

// collectU7Bytes copies all bytes with msbit=1 into dst and returns their count.
size_t collectU7Bytes( uint8_t * dst, const uint8_t * src, size_t slen ){
    uint8_t * p = dst;
    for( int i = 0; i < slen; i++ ){
        if(0x80 & src[i]){
            *p++ = src[i];
        }
    }
    return p - dst;
}

//! restorePacket reconstructs original data using src, slen, u8, u8len and table into dst and returns the count.
size_t restorePacket( uint8_t * dst, const uint8_t * table, const uint8_t * u8, size_t u8len, const uint8_t * src, size_t slen ){
    uint8_t * p = dst;
    for( int i = 0; i < slen; i++ ){
        if( 0x80 & src[i] ){ // an u7 byte
            if( u8len-- > 0){
                *p++ = *u8++;
            }
        }else{ // an id
            size_t sz = getPatternFromId( p, table, src[i] );
            p += sz;
        }
    }
    return p - dst;
}