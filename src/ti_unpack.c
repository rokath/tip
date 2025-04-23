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
#if OPTIMIZE_UNREPLACABLES == 1
    if (uTlen <= 0 ) { // Unrplacable byte optimisation was possible.
        u8len = -uTlen;
        memcpy( u8, uT8, u8len );
    } else { // Otherwise the last byte is an unreplacable and not the only one and there is at least one ID.
        u8len = reconvertBits( u8, uT8, uTlen ); // Optimization was not possible.
    }
#else // #if OPTIMIZE_UNREPLACABLES == 1
    u8len = reconvertBits( u8, uT8, uTlen );
#endif // #else // #if OPTIMIZE_UNREPLACABLES == 1
    size_t dlen = restorePacket( dst, table, u8, u8len, src, slen );
    return dlen;
}

//! isID checks, if at p an ID or ID begin occures.
//! @retval 0: no ID
//! @retval 1: primary ID
//! @retval 2: secondary ID (id1 at p[0] and id2 at p[1])
static int isID( const uint8_t * p){
    if (*p <= ID1Count){
        return 1;
    }else if (*p <= ID1Max){
        return 2;
    }
    return 0;
}

#if OPTIMIZE_UNREPLACABLES == 1
/*
//! endsWithID checks, if src ends with an ID.
//! @retval 0: no ID
//! @retval 1: primary ID
//! @retval 2: secondary ID (id1 at p[0] and id2 at p[1])
static int endsWithID(const uint8_t * src, size_t slen){
    if (slen == 0){
        return 0;
    }
    if (slen == 1){
        uint8_t by = src[slen-1];
        if (by <= ID1Count){
            return 1;
        } else if (by <= ID1Max){
            for(;;); // unexpected, option for further optimization?
        } // The last src byte is > ID1Count, what would expect one more byte. 
        return 0;    
    }
	// tip: 80bc 79d1 7fc7 7fc1 7fb8            f0 7979   a8 7fea 7ffe 7e89 7c6e 7fd0 790d 50 7fec 7f72                           e8bcaf 7c09 79ad 7d7b     f0
	//      uuuu id   id   id   id              uu id     uu id   id   id   id   id   id   ID id   id                             uuuuuu id   id   id       uu // uu len=8
    // At this point, isID sees 0x7bf0 and returns a false positive answer.
    // We need to process from the buffer start!
    int result = isID( src + slen - 2); // slen >= 2
    if (result == 2) {
        return 2;
    }
    return isID( src + slen - 1); 
}
*/
#endif // #if OPTIMIZE_UNREPLACABLES == 1

// collectUTBytes copies all bytes with msbit=1 into dst and returns their count.
static int collectUTBytes( uint8_t * dst, const uint8_t * src, size_t slen ){
    uint8_t * p = dst;
    int lastByteIsId;
    for( int i = 0; i < slen; i++ ){
        int x = isID(src+i);
        if (x==0){
            *p++ = src[i]; // collect
            lastByteIsId = 0;
        }else{
            lastByteIsId = 1;
            if (x==2){
                i++; // indirect ID, ignore next byte (secondary ID)
            } // else primary ID, do nothing
        }
    }
    int count = p - dst;
#if OPTIMIZE_UNREPLACABLES == 1 // cases like III or IIU or UUIIIUII 
    if (count <= 1) { // TiP packet has no or max one unrplacable byte, cases like III or IIU
        count = -count; // Unreplacable bytes optimisation was possible.
    }else if (lastByteIsId) {// TiP packet ends with an ID.
        count = -count; // Unreplacable bytes optimisation was possible.
    }
#endif // #if OPTIMIZE_UNREPLACABLES == 1
    return count;
}

//! reconvertBits transmutes slen n-bit bytes in src to 8-bit units in dst.
//! @param src is a byte buffer.
//! @param slen is the n-bit byte count.
//! @param dst is the destination buffer. It is allowed to be equal src for in-place conversion.
//! @retval is count 8-bit bytes
//! @details buf is filled from the end (=buf+limit)
static size_t reconvertBits( uint8_t * dst, const uint8_t * src, size_t slen ){
    if (unreplacableContainerBits == 7){
        return shift78bit( dst, src, slen );
    }else if (unreplacableContainerBits == 6){
        return shift68bit( dst, src, slen );
    }else{
        for(;;);
    }
}

//! restorePacket reconstructs original data using src, slen, u8, u8len and table into dst and returns the count.
static size_t restorePacket( uint8_t * dst, const uint8_t * table, const uint8_t * u8, size_t u8len, const uint8_t * src, size_t slen ){
    uint8_t * p = dst;
    for( int i = 0; i < slen; i++ ){ // TODO: slen-1
        int x = isID(src+i);
        if(x == 1 ){ // primary ID
            size_t sz = getPatternFromId( p, table, src[i] );
            p += sz;
        } else if(x == 2) { // indirect ID + secondary ID
            uint8_t id1 = src[i++];
            uint8_t id2 = src[i];
            // See in tipTable.go func tipPackageIDs() and TiP Usermanual Appendix.
            int offs = ID1Count + 1;
            int id =(id1-offs)*255 + id2 - 1 + offs; // == 255*id1 - 254*offs + id2 - 1
            size_t sz = getPatternFromId( p, table, (id_t)id);
            p += sz;
        } else if (u8len) {
            *p++ = *u8++; // collect restored unreplacable byte
            u8len--;
        }
    }
    return p - dst;
}

//! getPatternFromId seaches in testTable for id.
//! @param pt is filled with the replace pattern if id was found.
//! @param table is the pattern table.
//! @param id is the replacement. Valid values for id are 1...MaxID.
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
