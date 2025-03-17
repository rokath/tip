//! @file pack.c
//! @brief This is the tip pack code. Works also without unpack.c.
//! @details Written for ressources constraint embedded devices.
//! This tip code avoids heavy stack usage by using static buffers and is therefore not re-entrant.
//! This implementation is coded for speed in favour RAM usage. 
//! If RAM usage matters, the replace list r could be a bit array at the end of the destination buffer just to mark the unreplacable bytes.
//! In a loop then the packed data can get constructed directly into the destination buffer by searching for the pattern a second time.
//! It is possible to use different tables at the same time, but the code needs to be changed a bit then.
//! @author thomas.hoehenleitner [at] seerose.net

#include <stddef.h>
#include <string.h>
#include <stdio.h>
#include "pack.h"
#include "tip.h"
#include "memmem.h"

#ifndef DEBUG
#define DEBUG 1
#endif

#ifndef TIP_MAX_PATH_COUNT
#define TIP_MAX_PATH_COUNT 100
#endif

#if DEBUG
void printIDPositionTable( void );
void printPath( uint8_t pidx );
void printSrcMap( void );
#endif

static loc_t IDPattern( const uint8_t ** patternAddress, uint8_t id );

size_t tip( uint8_t* dst, const uint8_t * src, size_t len ){
    return tiPack( dst, idTable, src, len );
}

// idPatTable points to param table passed to some functions.
//! This allows using different idTable's than idTable.c 
//! especially for testing and not to have to pass it to all functions. 
static const uint8_t *idPatTable = idTable;

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
/*static*/ void insertIDPosSorted(uint8_t id, loc_t offset){
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

//! createIDPosTable uses idPatTable and parses src buffer for matching pattern
//! and creates a idPosTable specific to the actual src buffer.
//! It adds IDs with offset in a way, that smaller offsets occur first.
STATIC void createIDPosTable(const uint8_t * IDPatTable, const uint8_t * src, size_t slen){
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
            loc_t loc = pos - src;
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

//! IDPatternLength returns pattern length of id. 
static loc_t IDPatternLength( uint8_t id ){
    const uint8_t * next = idPatTable;
    for( int i = 1; i < id; i++ ){
        next += 1 + *next;
    }
    uint8_t len = *next;
    return len;
}

//! IDPattern writes pattern address of id and returns pattern length.. 
static loc_t IDPattern( const uint8_t ** patternAddress, uint8_t id ){
    const uint8_t * next = idPatTable;
    for( int i = 1; i < id; i++ ){
        next += 1 + *next;
    }
    uint8_t len = *next++;
    *patternAddress = next;
    return len;
}

//! forkPath extends srcMap with a copy of path pidx and returns index of copy.
static uint8_t forkPath( uint8_t pidx ){
    uint8_t psize = srcMap.path[pidx][0] + 1; // path: cnt, idx, idx, ...
    memcpy(srcMap.path[srcMap.count], srcMap.path[pidx], psize);
    return srcMap.count++;
}

//! appendPosTableIndexToPath appends position table index pti to pidx.
static void appendPosTableIndexToPath( uint8_t pidx, uint8_t pti ){
    uint8_t cnt = srcMap.path[pidx][0]; // cnt is the idx count in the (forked) path.
    uint8_t idx = cnt + 1;              // next free place
    srcMap.path[pidx][idx] = pti;       // write pti
    srcMap.path[pidx][0] = cnt + 1;     // one more idx
    #if DEBUG
    printf( "forked:" );
    printPath(pidx);
    #endif
}
/*
//! IDPosLimit returns first offset after ID position idx.
STATIC loc_t IDPosLimit(uint8_t idx){
    uint8_t id = IDPosTable.item[idx].id;
    loc_t len = IDPatternLength( id );
    loc_t limit = IDPosTable.item[idx].start + len;
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
*/

void createSrcMap(const uint8_t * table, const uint8_t * src, size_t slen){
    createIDPosTable(table, src, slen); // Get all ID positions in src ordered by increasing offset.

#if DEBUG
    printIDPositionTable();
#endif

    memset(&srcMap, 0, sizeof(srcMap)); // Start with no path (PathCount=0).
    for( int pti = 0; pti < IDPosTable.count; pti++ ){ // Loop over IDPosition table.
        IDPosition_t idPos = IDPosTable.item[pti]; // For each next idPos nnn:
        uint8_t nnn_id = idPos.id;
        loc_t nnn_start = idPos.start;
        loc_t nnn_len = IDPatternLength( nnn_id );
        loc_t nnn_limit = nnn_start + nnn_len;
        int IDPosAppended = 0;

        #if DEBUG
        char msg[200];
        int mlen = 0;
        int mlenMax = 0;
        #endif

        int srcMapCount =  srcMap.count;
        for( int k = 0; k <  srcMapCount; k++ ){ // Loop over all so far existing paths.
            if( srcMapCount > TIP_MAX_PATH_COUNT ){ // Create no new paths for this src buffer.

                #if DEBUG
                printf( "srcMap is full (%d paths)", TIP_MAX_PATH_COUNT);
                #endif

                IDPosAppended = 1; // Do not add any further paths.
                break;
            }

            uint8_t * path = srcMap.path[k]; // path is next path in srcMap.
            uint8_t pcnt = path[0]; // pcnt is the number od IDPosTable indices in this path.

            //  #if DEBUG
            //  uint8_t * pidx = path+1; // pidx is start of pcnt IDPosTable indices.
            //  for( int p = 0; p < pcnt; p++ ){
            //      IDPosition_t PathIdPos = IDPosTable.item[p];
            //      uint8_t ppp_Id = PathIdPos.id;
            //      loc_t ppp_start = PathIdPos.start;
            //      loc_t ppp_len = IDPatternLength( ppp_Id );
            //      loc_t ppp_limit = ppp_start + ppp_len;
            //  }
            //  #endif

            uint8_t idx = path[pcnt]; // idx is last IDPosTable index in this path.
            IDPosition_t lastIdPos = IDPosTable.item[idx]; // lastIdPos is the last (referenced) idPos in this path.
            uint8_t lll_Id = lastIdPos.id;
            loc_t lll_start = lastIdPos.start;
            loc_t lll_len = IDPatternLength( lll_Id );
            loc_t lll_limit = lll_start + lll_len;

            // case
            //   path: ppp...lll        - path k                                    | comment  | action
            // 0 patt:   nnn            - new pattern lays complete before          | error    | ignore pattern, take next path
            // 1 patt:     nnnN         - new pattern overlaps only start           | error    | ignore pattern, take next path
            // 2 patt:     nnNNN        - new pattern overlaps start and ends equal | error    | ignore pattern, take next path
            // 3 patt:     nnNNNnn      - new pattern overlaps full                 | error    | ignore pattern, take next path
            // 4 patt:       NNN        - new pattern matches exactly               | error    | ignore pattern, take next path
            // 5 patt:       NN         - new pattern matches start and is shorter  | possible | cannot append to this path
            // 6 patt:        N         - new pattern lays coplete inside           | possible | cannot append to this path
            // 7 patt:        NN        - new pattern matches end and is shorter    | possible | cannot append to this path
            // 8 patt:       NNNnn      - new pattern overlaps end and starts equal | possible | cannot append to this path
            // 9 patt:         Nnn      - new pattern overlaps only end             | possible | cannot append to this path
            //10 patt:          nnnn    - new pattern lays complete after           | possible | fork path k and append pattern to forked

            if( lll_limit <= nnn_start ){ // case 10
                #if DEBUG
                mlen = 0; // sprintf( msg, "case 10 !!! - addPatternToForked: lll_limit %d <= %d nnn_start - new pattern lays complete after", lll_limit, nnn_start );
                #endif
                uint8_t n = forkPath(k); 
                appendPosTableIndexToPath(n, pti);
                IDPosAppended = 1; // idx is appended to at least one path now. 
            }
            else if( nnn_limit <= lll_start ){ // case 0
                #if DEBUG
                mlen = sprintf( msg, "case 0 - unexpected, error: nnn_limit %d <= %d lll_start - new pattern lays complete before", nnn_limit, lll_start); // IDPosTable is sorted by rising positions
                #endif
            }
            else if( nnn_start < lll_start && nnn_limit < lll_limit ){ // case 1
                #if DEBUG
                mlen = sprintf( msg, "case 1 - unexpected, error: nnn_start %d < %d lll_start && nnn_limit %d < %d lll_limit - new pattern overlaps only start", nnn_start, lll_start, nnn_limit, lll_limit ); // IDPosTable is sorted by rising positions.
                #endif
            }
            else if( nnn_start < lll_start && nnn_limit == lll_limit ){ // case 2
                #if DEBUG
                mlen = sprintf( msg, "case 2 - unexpected, error: nnn_start %d < %d lll_start && nnn_limit %d == %d lll_limit - new pattern overlaps only start", nnn_start, lll_start, nnn_limit, lll_limit); // IDPosTable is sorted by rising positions.
                #endif
            }
            else if( nnn_start < lll_start && nnn_limit > lll_limit ){ // case 3
                #if DEBUG
                mlen = sprintf( msg, "case 3 - unexpected, error: nnn_start %d < %d lll_start && nnn_limit %d > %d lll_limit - new pattern overlaps full", nnn_start, lll_start, nnn_limit, lll_limit ); // IDPosTable is sorted by rising positions.
                #endif
            }
            else if( nnn_start == lll_start && nnn_limit == lll_limit ){ // case 4
                #if DEBUG
                mlen = sprintf( msg, "case 4 - unexpected, error: nnn_start %d == %d lll_start && nnn_limit %d == %d lll_limit - new pattern matches exactly", nnn_start, lll_start, nnn_limit, lll_limit ); // IDPosTable has uniq idPos values.
                #endif
            }
            else if( nnn_start == lll_start && nnn_limit < lll_limit ){ // case 5 
                #if DEBUG
                mlen = 0; // sprintf( msg, "case 5 -  cannot append to this path: nnn_start %d == %d lll_start && nnn_limit %d < %d lll_limit - new pattern matches start and is shorter", nnn_start, lll_start, nnn_limit, lll_limit ); // cannot append to this path
                #endif
            }
            else if( nnn_start > lll_start && nnn_limit < lll_limit ){ // case 6
                #if DEBUG
                mlen = 0; // sprintf( msg, "case 6 -  cannot append to this path: nnn_start %d > %d lll_start && nnn_limit %d < %d lll_limit - new pattern lays completely inside",  nnn_start, lll_start, nnn_limit, lll_limit ); // cannot append to this path
                #endif
            }
            else if( nnn_start > lll_start && nnn_limit == lll_limit ){ // case 7
                #if DEBUG
                mlen = 0; // sprintf( msg, "case 7 -  cannot append to this path: nnn_start %d > %d lll_start && nnn_limit %d == %d lll_limit - new pattern matches end and is shorter", nnn_start, lll_start, nnn_limit, lll_limit); // cannot append to this path
                #endif
            }
            else if( nnn_start == lll_start && nnn_limit > lll_limit ){ // case 8
                #if DEBUG
                mlen = sprintf( msg, "case 8 -  cannot append to this path: nnn_start %d == %d lll_start && nnn_limit %d > %d lll_limit - new pattern matches start and is longer", nnn_start,  lll_start, nnn_limit, lll_limit ); // cannot append to this path
                #endif
            }
            else if( nnn_start > lll_start && nnn_limit > lll_limit ){ // case 9
                #if DEBUG
                mlen = 0; // sprintf( msg, "case 9 -  cannot append to this path: nnn_start %d > %d lll_start && nnn_limit %d > %d lll_limit - new pattern overlaps end", nnn_start, lll_start, nnn_limit, lll_limit ); // cannot append to this path
                #endif
            }
            #if DEBUG
            mlenMax = mlen > mlenMax ? mlen : mlenMax;
            if( mlen > 0 ){
                printf( "%s (mlenMax %d)\n", msg, mlenMax );
            }
            #endif
        }
        if( !IDPosAppended ){ 
            int nextIdx = srcMap.count;
            #if DEBUG
            printf( "Create a new path%3d with pit%3d (id%3d, loc%3d)\n", srcMap.count, pti, idPos.id, idPos.start );
            #endif
            srcMap.path[nextIdx][0] = 1;   // one IDPos in this new path
            srcMap.path[nextIdx][1] = pti; // the pti (the first is naturally 0)
            srcMap.count++;                // one more path
        }
    }
}

//! IDPosLength returns first offset after ID position idx.
STATIC loc_t IDPosLength(uint8_t idx){
    uint8_t id = IDPosTable.item[idx].id;
    loc_t len = IDPatternLength( id );
    return len;
}

//! pathPatternSize returns sum of all pattern lengths in path pidx.
loc_t pathPatternSize( int pidx ){
    if( srcMap.count <= pidx ){
        return 0;
    }
    uint8_t * path = srcMap.path[pidx];
    uint8_t pathIdxCount = path[0];
    uint8_t * idx = &(path[1]);
    loc_t sum = 0;
    for( int i = 0; i < pathIdxCount; i++ ){
        sum += IDPosLength(*idx++);
    }
    return sum;
}

//! MinDstLengthPath returns srcMap path index which results in a shortest package.
uint8_t MinDstLengthPath(void){
    loc_t maxSum = 0;
    uint8_t pathIndex = 0;
    for( int i = 0; i < srcMap.count; i++ ){
        loc_t psiz = pathPatternSize(i);
        if( psiz > maxSum ){
            maxSum = psiz;
            pathIndex = i;
        }
    } // Maximum psiz results in smallest tip size.
    return pathIndex;
}

//! selectUnreplacableBytes coppies all unreplacable bytes from src to dst.
//! It uses idPatTable and the path index pidx in the actual srcMap, which is linked to IDPosTable.
size_t selectUnreplacableBytes( uint8_t * dst, uint8_t pidx, const uint8_t * src, size_t slen ){
    uint8_t * path = srcMap.path[pidx]; // This is the path we use. 
    uint8_t count = path[0]; // The path contains: count, IDPosTable index, IDPosTable index, ...
    const uint8_t * srcNext = src; // next position for src buffer read
    uint8_t * dstNext = dst;       // next position for dst buffer write
    uint8_t * tidx = path+1; // Here are starting the IDPosTable indices.
    loc_t u8sum = 0;
    for( int i = 0; i < count; i++ ){
        IDPosition_t idPos = IDPosTable.item[tidx[i]];
        uint8_t id = idPos.id;
        const uint8_t * patternFrom = src + idPos.start; // pattern start in src buffer
        loc_t u8len = patternFrom - srcNext; // count of unreplacable bytes
        memcpy( dstNext, srcNext, u8len );
        loc_t patlen = IDPatternLength( id );
        srcNext += patlen + u8len;
        dstNext += u8len;
        u8sum += u8len;
    }
    size_t rest = slen - (srcNext - src); // total - pattern sum
    memcpy( dstNext, srcNext, rest );
    dstNext += rest;
    size_t len = dstNext - dst;
    return len;
}

//! createOutput uses the u7 buffer and pidx to intermix transformed unreplacable bytes and pattern IDs.
//! It uses idPatTable and the path index pidx in the actual srcMap, which is linked to IDPosTable.
size_t createOutput( uint8_t * dst, uint8_t pidx, const uint8_t * u7src, size_t u7len, const uint8_t * src ){
    uint8_t * path = srcMap.path[pidx]; // This is the path we use. 
    uint8_t count = path[0]; // The path contains: count, IDPosTable index, IDPosTable index, ...
    const uint8_t * srcNext = src; // next position for src buffer read
    uint8_t * dstNext = dst;       // next position for dst buffer write
    uint8_t * tidx = path+1; // Here are starting the IDPosTable indices.
    const uint8_t * u7Next = u7src;
    for( int i = 0; i < count; i++ ){
        IDPosition_t idPos = IDPosTable.item[tidx[i]];
        uint8_t id = idPos.id;
        const uint8_t * patternFrom = src + idPos.start; // pattern start in src buffer
        loc_t u8len = patternFrom - srcNext; // count of unreplacable bytes
        loc_t patlen = IDPatternLength( id );
        srcNext += patlen + u8len;
        memcpy( dstNext, u7Next, u8len ); // Copy u8len bytes from u7src buffer.
        u7Next += u8len;
        dstNext += u8len;
        *dstNext++ = id; // Write the pattern replace id.
    }
    size_t rest = u7len - (u7Next - u7src);
    memcpy( dstNext, u7Next, rest );
    dstNext += rest;
    size_t len = dstNext - dst;
    return len;
}


//! shift87bit transforms slen 8-bit bytes in src to 7-bit units.
//! @param src is the bytes source buffer.
//! @param slen is the 8-bit byte count.
//! @param lst is the last address inside the dst buffer.
//! @retval is count of 7-bit bytes after operation. 
//! @details The dst buffer is filled from the end.That allows to do an in-buffer conversion.
//! The destination address is computable afterwards: dst = lim - retval.
//! lst is allowed to be "close" behind buf + slen, thus making in-place conversion possible.
//! Example: slen=17, lst=src+24-1
//!       (src) <---            slen=17                   --->(u8)
//! slen=17: b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 __ __ __ __ __ __ __
//! ret =20: __ __ __ __ m7 b7 b7 b7 m7 b7 b7 b7 b7 b7 b7 b7 m7 b7 b7 b7 b7 b7 b7 b7
//!                   (dst) <---                ret=20                       --->(lst)
//! In dst all MSBits are set to 1, to avoid any zeroes.
//! The data are processed from the end.
/*static*/ size_t shift87bit( uint8_t* lst, const uint8_t * src, size_t slen ){
    const uint8_t * u8 = src + slen; // first address behind src buffer
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

size_t buildTiPacket(uint8_t * dst, uint8_t * dstLimit, const uint8_t * table, const uint8_t * src, size_t slen){
    createSrcMap(table, src, slen);

#if DEBUG
    printSrcMap();
#endif

    uint8_t pidx = MinDstLengthPath(); // find minimum line

#if DEBUG
    printf( "MinDstLengthPath: %u\n", pidx );
#endif

    memset(dst, 0, dstLimit-dst);
    loc_t u8Count = selectUnreplacableBytes(dst, pidx, src, slen );

#if DEBUG
   printf( "u8Count: %ul\n", u8Count );
#endif

    loc_t u7Count = shift87bit( dstLimit-1, dst, u8Count );
    uint8_t * u7src = dstLimit - u7Count;

#if DEBUG
    printf( "u7Count: %ul\n", u7Count );
#endif

    size_t pkgSize = createOutput( dst, pidx, u7src, u7Count, src );
    return pkgSize; // final ti package size
}


#if DEBUG

void printPath( uint8_t pidx ){
    uint8_t * path = srcMap.path[pidx]; 
    int plen = path[0];
    loc_t psiz = pathPatternSize(pidx);
    printf( "path%3d: pattern sum size%3d, plen%3d: ", pidx, psiz, plen);
    for( int k = 0; k < plen; k++ ){
        printf( "idx%3d, ", path[k+1]);
    }
    printf( "\n" );
}

void printSrcMap( void ){
    printf( "srcMap: -----------\n");
    for( int i = 0; i < srcMap.count; i++ ){
        printPath(i);
    }
    printf( "-----------\n");
}

void printIDPositionTable( void ){
    printf( "IDPositionTable:\n");
    printf(" idx | id  | pos | ASCII\n");
    printf("-----|-----|-----|------\n");
    for( int idx = 0; idx < IDPosTable.count; idx++ ){
        uint8_t id = IDPosTable.item[idx].id;
        loc_t loc = IDPosTable.item[idx].start;
        const uint8_t * pattern;
        loc_t length = IDPattern( &pattern, id);
        uint8_t s[100] = {0};
        memcpy(s, pattern, length);
        printf(" %3d | %3d | %3d | '%s' \n", idx, id, loc, s);
    }
    printf("------------------------\n");
}

#endif
