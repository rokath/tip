//! @file ti_unpack.c
//! @brief This is the tip unpack code. Works also without pack.c.
//! @details todo
//! @author thomas.hoehenleitner [at] seerose.net

#include <string.h>
#include "ti_unpack.h"
#include "tipInternal.h"

static int collectUTBytes( uint8_t * dst, const uint8_t * src, size_t slen );
static size_t reconvertBits( uint8_t * lst, const uint8_t * src, size_t slen );
static size_t restorePacket( uint8_t * dst, const uint8_t * table, const uint8_t * u8, size_t u8len, const uint8_t * src, size_t slen );
static size_t getPatternFromId( uint8_t * pt, const uint8_t * table, id_t id );

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
        if(src[i] <= ID1Count ){
            // primary ID, do nothing
        } else if(src[i] <= ID1Max) {
            i++; // indirect ID, ignore next byte (secondary ID)
        } else {
            *p++ = src[i]; // collect
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
        if(src[i] <= ID1Count ){ // primary ID
            size_t sz = getPatternFromId( p, table, src[i] );
            p += sz;
        } else if(src[i] <= ID1Max) { // indirect ID + secondary ID
            uint8_t indirectID = src[i++];
            uint8_t secondaryID = src[i];
            // Example for understanding the id computation with ID1Count=124 and ID1Max=127:
            // ID1                 1:                       =   1
            // ID1               ...:                       = ...
            // ID1 = ID1Count    124:                       = 124
            // indirectID        125: (0*255) + 0...254 - 0 = 125...379 <- level 0
            // indirectID        126: (1*255) + 0...254 - 1 = 380...508 <- level 1
            // indirectID=ID1Max 127: (2*255) + 0...254 - 2 = 509...763 <- level 2
            int level = indirectID - ID1Count - 1;
            id_t id = level * 255 + (secondaryID-1) - level;
            size_t sz = getPatternFromId( p, table, id);
            p += sz;
        } else {
            *p++ = src[i]; // collect
        }
    }
    return p - dst;
}

//! getPatternFromId seaches in testTable for id.
//! @param pt is filled with the replace pattern if id was found.
//! @param table is the pattern table.
//! @param id is the replace byte. Valid values for id are 1...0x7f.
//! @retval is the pattern size or 0, if id was not found.
static size_t getPatternFromId( uint8_t * pt, const uint8_t * table, id_t id ){
    size_t sz;
    id_t idx = 0x01;
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
