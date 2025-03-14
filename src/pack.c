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

size_t tip( uint8_t* dst, const uint8_t * src, size_t len ){
    return tiPack( dst, idTable, src, len );
}

// idPatTable points to param table passed to some functions.
//! This allows using different idTable's than idTable.c 
//! especially for testing and not to have to pass it to all functions. 
static const uint8_t *idPatTable;

//! @brief tiPack encodes src buffer with size len into dst buffer and returns encoded len.
//! @details For the tip encoding it uses the linked idTable.c object.
// - Some bytes groups in the src buffer are replacable with IDs 0x01...0x7f and some not.
// - The rlist r holds the replace information. Additionally the dst buffer is prefilled with IDs from both sides.
// - Example: dst = 5, 6, 0, ..., 0, 2, 3, 4; and rlist = 00111111000111000001111110111
// - ID 5 has 2 and ID 6 4 bytes, so ID 2 and 4 have 3 bytes and ID 3 has 5 bytes.
// - The unreplacable bytes are collected into a buffer.
size_t tiPack( uint8_t * dst, const uint8_t * table, const uint8_t * src, size_t slen ){
    if( slen == 0 || TIP_SRC_BUFFER_SIZE_MAX < slen ){
        return 0;
    }
    size_t dstSize = ((18725ul*slen)>>14)+1;  // The max possible dst size is len*8/7+1 or ((len*65536*8/7)>>16)+1;
    uint8_t * dstLimit = dst + dstSize;
    memset(dst, 0, dstSize);
    idPatTable = table;
    size_t tipSize = buildTiPacket(dst, dstLimit, table, src, slen);
    return tipSize;
}

//! nextIDPatTablePos points to the ID pattern table next pattern position.
static const uint8_t * nextIDPatTablePos = NULL;

//! initGetNextPattern causes getNextPattern to start from 0.
static void initGetNextPattern( const uint8_t * idTbl ){
    idPatTable = idTbl;
    nextIDPatTablePos = idTbl;
}

//! getNextPattern returns next pattern location in pt and size in sz or *sz == 0.
//! @param pt is filled with the replace pattern address if exists.
//! @param sz is filled with the replace size or 0, if not exists.
static void getNextPattern(const uint8_t ** pt, size_t * sz ){
    if( (*sz = *nextIDPatTablePos++) != 0 ){ // a pattern exists here
        *pt = nextIDPatTablePos;
        nextIDPatTablePos += *sz;
        return;
    }
}

//! IDPosTable holds all IDs with their positions occuring in the current src buffer.
IDPosTable_t IDPosTable = {0};

//! insertIDPosSorted inserts id with pos and len into IDPosTable with smallest pos first.
/*static*/ void insertIDPosSorted(uint8_t id, offset_t offset){
    int i;
    int insertFlag = 0;
    for( i = 0; i < IDPosTable.count; i++ ){
        if( offset < IDPosTable.item[i].start ){
            insertFlag = 1;
            break; // insert here
        }
    }
    if( insertFlag ){ 
        IDPosition_t *dst = &IDPosTable.item[i+1];
        IDPosition_t *src = &IDPosTable.item[i];
        size_t size = (IDPosTable.count-i)*sizeof(IDPosition_t);
        memmove(dst, src, size);
    }
    IDPosTable.item[i].id = id;
    IDPosTable.item[i].start = offset;
    IDPosTable.count++;
}

//! newIDPosTable uses idPatTable and parses src buffer for matching pattern
//! and creates a idPosTable specific to the actual src buffer.
//! It adds IDs with offset in a way, that smaller offsets occur first.
STATIC void newIDPosTable(const uint8_t * IDPatTable, const uint8_t * src, size_t slen){
    memset(&IDPosTable, 0, sizeof(IDPosTable));
    initGetNextPattern(IDPatTable);
    for( int id = 1; id < 0x80; id++ ){ // Traverse the ID table. 
        const uint8_t * needle = NULL;
        size_t nlen;
        repeat:
        getNextPattern( &needle, &nlen );
        if( nlen == 0 ){ // End of table reached, if less 127 IDs.
            break; 
        }
        int offset = 0;
        while(offset<slen-1){
            uint8_t * pos = memmem(src+offset, slen-offset, needle, nlen);
            if(pos == NULL){ // Needle not found. 
                id++; // increment "manually"
                goto repeat; // Pattern not found, try next pattern.
            }
            offset_t loc = pos - src;
            insertIDPosSorted(id, loc);
            offset = loc + 1; // We search the identical pattern in the while loop.
            // "xxxxxPPPxxx" - after finding first PP, we need to find the 2nd PP inside PPP.
        }
    }
}

//! srcMap holds all possible paths for current src buffer.
//! - cnt, idx, idx, ...
//! -   3,  17,   5,  4, // a path with 3 IDpos
STATIC srcMap_t srcMap = {0};

//! initSrcMap resets path table.
static void initSrcMap( void ){
    memset(&srcMap, 0, sizeof(srcMap));
}

//! IDPatternLength returns pattern length of id. 
static offset_t IDPatternLength( uint8_t id ){
    const uint8_t * next = idPatTable;
    for( int i = 1; i < id; i++ ){
        next += 1 + *next;
    }
    uint8_t len = *next;
    return len;
}

//! IDPosLimit returns first offset after ID position idx.
STATIC offset_t IDPosLimit(uint8_t idx){
    uint8_t id = IDPosTable.item[idx].id;
    offset_t len = IDPatternLength( id );
    offset_t limit = IDPosTable.item[idx].start + len;
    return limit;
}

//! IDPosAppendableToPath checks if pathIndex limit is small enough to append IDPos.
//! \param pathIndex is the path to check.
//! \param IDPosIdx is the ID position inside IDPosTable.
static int IDPosAppendableToPath( uint8_t pathIndex, uint8_t idPos ){
    uint8_t pathIdPosCount = srcMap.path[pathIndex][0];
    uint8_t lastIdPos = srcMap.path[pathIndex][pathIdPosCount];
    if( IDPosLimit(lastIdPos) <= IDPosTable.item[idPos].start ){
        return 1;
    }
    return 0;
}

//! forkPath extends srcMap with a copy of path pidx and returns index of copy.
static uint8_t forkPath( uint8_t pidx ){
    uint8_t psize = srcMap.path[pidx][0] + 1;
    memcpy(srcMap.path[srcMap.count], srcMap.path[pidx], psize);
    return srcMap.count++;
}

//! appendIDPos appends idpos to pidx.
static void appendIDPos( uint8_t pidx, uint8_t idpos ){
    uint8_t cnt = srcMap.path[pidx][0]; // pidx idpos count
    uint8_t idx = cnt + 1;           // next free place
    srcMap.path[pidx][idx] = idpos;     // write idpos
    srcMap.path[pidx][0] = cnt + 1;     // one more idpos
}

void createSrcMap(const uint8_t * table, const uint8_t * src, size_t slen){
    newIDPosTable(table, src, slen); // Get all ID positions in src ordered by increasing offset.
    initSrcMap();                 // Start with no path (PathCount=0).
    for( int idPos = 0; idPos < IDPosTable.count; idPos++ ){ // Loop over IDPosition table for each IDPos.
        int IDPosAppended = 0;
        for( int k = srcMap.count - 1; k >= 0; k-- ){ // Loop over all so far existing paths from the end.
            if( srcMap.count > TIP_MAX_PATH_COUNT ){ // Create no new paths for this buffer.
                IDPosAppended = 1;
                break;
            }
            if( IDPosAppendableToPath(k, idPos) ){ // ID position idPos fits to path k.
                uint8_t n = forkPath(k); 
                appendIDPos(n,idPos);
                IDPosAppended = 1;
            }
        }
        if( !IDPosAppended ){ // Create a new path with idPos.
            int nextIdx = srcMap.count;
            srcMap.path[nextIdx][0] = 1;     // one IDPos in this new path
            srcMap.path[nextIdx][1] = idPos; // the IDPos (the first is naurally 0)
            srcMap.count++;                  // one more path
        }
    }
}

size_t buildTiPacket(uint8_t * dst, uint8_t * dstLimit, const uint8_t * table, const uint8_t * src, size_t slen){
    createSrcMap(table, src, slen);

    size_t pkgSize = 0;   // final ti packgae size

    // find minimum line
    // create output
    return pkgSize;
}

/*
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
*/
