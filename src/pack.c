//! @file pack.c
//! @brief This is the tip pack code. Works also without unpack.c.
//! @details Written for ressources constraint embedded devices.
//! This tip code avoids heavy stack usage by using static buffers and is therefore not re-entrant.
//! This implementation is coded for speed in favour RAM usage. 
//! If RAM usage matters, the replace list r could be a bit array at the end of the destination buffer just to mark the unreplacable bytes.
//! In a loop then the packed data can get constructed directly into the destination buffer by searching for the pattern a second time.
//! It is possible to use different tables at the same time, but the code needs to be changed a bit then.
//! @author thomas.hoehenleitner [at] seerose.net

#include <string.h>
#include "pack.h"
#include "tip.h"
#include "memmem.h"

static const uint8_t * nextTablePos = 0;
static unsigned int nextID = 1;

//! initGetNextPattern causes getNextPattern to start from 0.
static void initGetNextPattern( const uint8_t * table ){
    nextTablePos = table;
    nextID = 1;
}

//! getNextPattern returns next pattern location in pt and size in sz or *sz == 0.
//! @param pt is filled with the replace pattern address if exists.
//! @param sz is filled with the replace size or 0, if not exists.
static void getNextPattern(const uint8_t ** pt, size_t * sz ){
    if( (*sz = *nextTablePos++) != 0 ){ // a pattern exists here
        *pt = nextTablePos;
        nextTablePos += *sz;
        return;
    }
}

size_t tip( uint8_t* dst, const uint8_t * src, size_t len ){
    return tiPack( dst, idTable, src, len );
}

//! @brief tiPack encodes src buffer with size len into dst buffer and returns encoded len.
//! @details For the tip encoding it uses the linked idTable.c object.
// - Some bytes groups in the src buffer are replacable with IDs 0x01...0x7f and some not.
// - The rlist r holds the replace information. Additionally the dst buffer is prefilled with IDs from both sides.
// - Example: dst = 5, 6, 0, ..., 0, 2, 3, 4; and rlist = 00111111000111000001111110111
// - ID 5 has 2 and ID 6 4 bytes, so ID 2 and 4 have 3 bytes and ID 3 has 5 bytes.
// - The unreplacable bytes are collected into a buffer.
size_t tiPack( uint8_t * dst, const uint8_t * table, const uint8_t * src, size_t slen ){
    size_t dstSize = ((18725*slen)>>14)+1;  // The max possible dst size is len*8/7+1 or ((len*65536*8/7)>>16)+1;
    uint8_t * dstLimit = dst + dstSize;
    if( slen > TIP_SRC_BUFFER_SIZE_MAX ){
        return 0;
    }
    memset(dst, 0, dstSize);
    size_t tipSize = buildTiPacket(dst, dstLimit, table, src, slen);
    return tipSize;
}


int matchingIDpositions = 0;

typedef struct{
    uint8_t ID;
    int pos;
} IDposition_t;

static IDposition_t IDpos[100];

void addMatch( uint8_t id, int offset){
    IDpos[matchingIDpositions].ID = id;
    IDpos[matchingIDpositions].pos = offset;
    matchingIDpositions++;
}

size_t buildTiPacket(uint8_t * dst, uint8_t * dstLimit, const uint8_t * table, const uint8_t * src, size_t slen){
    size_t pkgSize = 0;   // final ti packgae size
    initGetNextPattern(table);
    for( int id = 1; id < 0x80; id++ ){ // traverse the ID table. It is sorted by decreasing pattern length.    
        const uint8_t * needle = NULL;
        size_t nlen;
        getNextPattern( &needle, &nlen );
        if( nlen == 0 ){ // end of table if less 127 IDs
            break; 
        }
        uint8_t * pos = memmem(src, slen, needle, nlen);
        if(pos == NULL){
            continue;
        }
        int offset = pos - src;
        addMatch(id, offset);
    }
    for( int i = 0; i < matchingIDpositions; i++ ){
        // todo: find arrangements
    }
    return pkgSize;
}


//! buildTiPacket starts with buf=src and tries to find biggest matching pattern from table at buf AND bufLimit-nlen.
//! If a pattern was found at buf, buf is incremented by found pattern size.
//! If a pattern was found at bufLimit-nlen, bufLimit is decremented by found pattern size.
//! If a pattern was found at front or back we start over.
//! If none found, we increment buf by 1 (if possible) and decrement bufLimit (if possible) and start over.
//! It is possible, the same pattern is found again at the same place, but we do not care.
//! Why searching from 2 sides:
//! - ABC12ABC: table: C12,ABC,12A,12 -> ABC, 12A, uuu = 5 bytes, when only front search.
//! - ABC12ABC: table: C12,ABC,12A,12 -> uuu, C12, ABC = 5 bytes, when only back search.
//! - ABC12ABC: table: C12,ABC,12A,12 -> ABC, 12, ABC = 3 bytes, when front and back search, but how to match?
//! 2 possibilities:
//! - ABC 12A uuu        is front search result.
//! -       uuu C12 ABC  is back search result.
//! - If we subtract, we get a remaining 12
size_t buildTiPacket0(uint8_t * dst, uint8_t * dstLimit, const uint8_t * table, const uint8_t * src, size_t slen){
    const uint8_t * buf = src;             // src front pointer
    const uint8_t * bufLimit = src + slen; // src back pointer
    uint8_t * pkg = dst;                   // pkg front ponter
    uint8_t * pkgLimit = dstLimit;         // pkg back pointer
    int frontSearch = 1;                   // front search flag
    int backSearch = 1;                    // back search flag
    uint8_t u;            // unreplacable byte (not covered by a matching pattern)
    uint8_t msb;          // most significant bit of u
    uint8_t u7f = 0x80;   // collected bits 7 of u in front
    uint8_t u7b = 0x80;   // collected bits 7 of u in back
    size_t cu7f = 0;      // count of collected u7f bits
    size_t cu7b = 0;      // count of collected u7b bits
    size_t pkgSize = 0;   // final ti packgae size

    ///////////////
    return pkgSize;
    ///////////////

repeat:
    initGetNextPattern(table);
    for( int id = 1; id < 0x80; id++ ){ // traverse the ID table. It is sorted by decreasing pattern length.    
        int frontMatch = 0;
        int backMatch = 0;
        const uint8_t * needle = NULL;
        size_t nlen;
        getNextPattern( &needle, &nlen );
        if( nlen == 0 ){ // end of table if less 127 IDs
            break; 
        }
        if( frontSearch && 0 == strncmp((void*)buf, (void*)needle, nlen) ){ // match at buf front
            frontMatch = 1;
            *pkg++ = id; // write id
            buf += nlen; // adjust front pointer
            if( !(buf < src + slen)){
                frontSearch = 0;
            }
        }
        if( backSearch && 0 == strncmp((void*)(bufLimit-nlen), (void*)needle, nlen) ){ // match at buf back
            backMatch = 1;
            *--pkgLimit = id;
            bufLimit -= nlen;
            if( !(src < bufLimit) ){
                backSearch = 0;
            }          
        }
        if( frontMatch || backMatch ){
            goto repeat; // start over
        }
        // continue with next pattern  
    }
    // Arriving here means, that no table pattern fits to front or back.
    if( frontSearch && !backSearch){ // back search done, go on with front search
        u = *buf++;
        msb = 0x80 & u;
        u7f |= msb>>++cu7f;
        *pkg++ = 0x80 | u; // store 7 lsb and set msb
        if (cu7f == 7){
            *pkg++ = u7f;
            cu7f = 0;
            u7f = 0x80; // set msb already here
        }
        if( !(buf < src + slen)){
            frontSearch = 0;
        }
        goto repeat;
    }
    if( !frontSearch && backSearch ){ // front search done, go on with back search
        u = *bufLimit--;
        msb = 0x80 & u;
        u7b |= msb>>++cu7b;
        *--pkgLimit = 0x80 | u; // store 7 lsb and set msb
        if (cu7b == 7){
            *--pkgLimit = u7b;
            cu7b = 0;
            u7b = 0x80; // set msb already here
        }
        if( !(src < bufLimit) ){
            backSearch = 0;
        }
        goto repeat;
    }
    if( frontSearch && backSearch ){
        // Here it is not known, if we better reduce front or back now.
        // Reducing both sides may be wrong.
        // We should try one and the other independently and check what is better.
        
        ///////
        // how?
        ///////

        goto repeat;
    }
    // Arriving here means, that buf is >= src+slen and bufLimit is <= src.
    // In dst starts the packet front and its limit is pkg.
    // At dstLimit ends the packet back and its start is pkgLimit.
    // We need to unite u7f and u7b and to move the package end to touch the package start.

    ///////
    // todo
    ///////

    return pkgSize;
}
