//! @file ti_pack.c
//! @brief This is the tip pack code. Works also without unpack.c.
//! @details Written for ressources constraint embedded devices.
//! This code avoids heavy stack usage by using static buffers and is therefore not re-entrant.
//! This implementation is coded for speed in favour RAM usage. 
//! It is possible to use different tables in the same project, but only sequentially.

#include <stddef.h>
#include <string.h>
#include <stdio.h>
#include "tipInternal.h"
#include "memmem.h"

//! IDPosition_t describes a src buffer position which is replacable with an id.
typedef struct {
    uint8_t id;  // id of pattern found in src
    loc_t start; // id pattern start in src
} IDPosition_t;

//! IDPosTable_t holds all src buffer specific ID positions.
typedef struct {
    loc_t count; //! count is the number of items inside IDPosTable. In cannot exceed TIP_SRC_BUFFER_SIZE_MAX-1.
    IDPosition_t item[TIP_SRC_BUFFER_SIZE_MAX-1];
} IDPosTable_t;

//! path_t is a src buffer specific possible ID positions sequence consisting of index values into the IDPosTable.
typedef struct {
    loc_t count; // count is the amount of position indicies in this path. In cannot exceed TIP_SRC_BUFFER_SIZE_MAX/2.
    loc_t posIndex[TIP_SRC_BUFFER_SIZE_MAX/2]; // posIndex is an item index in the IDPosTable_t and cannot exeed TIP_SRC_BUFFER_SIZE_MAX-1.
} path_t;

//! srcMap_t holds all so far possible paths for the src buffer.
//! To limit its neede size, after each added IDPosTable idx, obsolete paths are removed.
typedef struct {
    unsigned count; // count is the actual path count in srcMap.
    path_t path[TIP_MAX_PATH_COUNT]; // Each path contains a count and count indexes into the IDPosTable.
} srcMap_t;

static size_t buildTiPacket(uint8_t * dst, uint8_t * dstLimit, const uint8_t * table, const uint8_t * src, size_t slen); // forward declaration

#if DEBUG
static void printIDPositionTable( void );              // forward declaration
//static void printSrcMap( void );
#endif
static void printPath( char * prefix, unsigned pidx ); // forward declaration

//! @brief idPatTable points to a parameter "table" passed to some functions.
//! @details This allows using different idTable's than idTable.c 
//! especially for testing and not to have to pass it to all functions.
//! The ID table has max 127 IDs with pattern, each max 255 bytes long.
//! ATTENTION: The pack functions are usable only sequentially!
static const uint8_t * idPatTable = idTable;

size_t tip( uint8_t * dst, const uint8_t * src, size_t len ){ // default user interface
    return tiPack( dst, idTable, src, len );
}

size_t tiPack( uint8_t * dst, const uint8_t * table, const uint8_t * src, size_t slen ){ // extended user interface
    if( TIP_SRC_BUFFER_SIZE_MAX < slen ){
#if DEBUG
        printf( "slen %llu > %u TIP_SRC_BUFFER_SIZE_MAX is invalid\n", slen, TIP_SRC_BUFFER_SIZE_MAX );
#endif
        return 0;
    }
    size_t dstSize = ((18725ul*slen)>>14)+1;  // The max possible dst size is len*8/7+1 or ((len*65536*8/7)>>16)+1.
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
static IDPosTable_t IDPosTable = {0};

//! insertIDPosSorted inserts id with pos and len into IDPosTable with smallest pos first.
static void insertIDPosSorted(uint8_t id, loc_t offset){
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
static void createIDPosTable(const uint8_t * IDPatTable, const uint8_t * src, size_t slen){
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
            insertIDPosSorted(id, loc); // We search the identical pattern in the while loop.
            offset = loc + 1; // "xxxxxPPPxxx" - after finding first PP, we need to find the 2nd PP inside PPP.
        }
    }    
#if DEBUG
    printIDPositionTable();
#endif
}

//! srcMap holds all possible paths for current src buffer.
//! - cnt, idx, idx, ...
//! -   3,  17,   5,  4, // a path with 3 IDpos
static srcMap_t srcMap = {0};

//! IDPatternLength returns pattern length of id. The max pattern length is 255.
static uint8_t IDPatternLength( uint8_t id ){
    const uint8_t * next = idPatTable;
    for( int i = 1; i < id; i++ ){
        next += 1 + *next;
    }
    uint8_t len = *next;
    return len;
}

// //! IDPattern returns pattern address of id. The max pattern length is 255.
// static uint8_t * IDPatternAddress( uint8_t id ){
//     const uint8_t * next = idPatTable;
//     for( int i = 1; i < id; i++ ){
//         next += 1 + *next;
//     }
//     return next;
// }

static void printPatternAsASCII( uint8_t id ){
    const uint8_t * next = idPatTable;
    for( int i = 1; i < id; i++ ){
        next += 1 + *next;
    }
    uint8_t len = *next++;
    const uint8_t * pat = next;
    char msg[256] = {0};
    for( int i = 0; i < len; i++ ){
        char c = ' ';
        if( 32 <= pat[i] && pat[i] < 128 ){
            c = (char)(pat[i]);
        }
        sprintf( msg+i, "%c", c);
    }
    printf( msg );
}

//! forkPath extends srcMap with a copy of path pidx and returns index of copy.
static unsigned forkPath( unsigned pidx ){
    memcpy(&srcMap.path[srcMap.count], &srcMap.path[pidx], sizeof(path_t));
    return srcMap.count++;
}

//! appendPosTableIndexToPath appends position table index pti to pidx.
static void appendPosTableIndexToPath( unsigned pidx, loc_t pti ){
    path_t * path = &srcMap.path[pidx];
    loc_t last = path->count++;
    path->posIndex[last] = pti;
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
static int IDPosAppendableToPath( path_t pathIndex, uint8_t idPos ){
    path_t pathIdPosCount = srcMap.path[pathIndex][0];
    uint8_t lastIdPos = srcMap.path[pathIndex][pathIdPosCount];
    if( IDPosLimit(lastIdPos) <= IDPosTable.item[idPos].start ){
        return 1;
    }
    return 0;
}
*/

        // shrinkSrcMap removes unneded paths. Exampöe:
        // src: ABCDABCD
        //   0: ABCD     - delete
        //   1: ABC      - delete
        //   2: AB       - delete
        //   3: AB  ABCD -    delete
        //   4: ABC ABCD -    delete
        //   5: ABCDABCD
        //   6:     ABC - already not possible!
        // - After an idx went thru srcMap: ABCD 
        //   - If several paths contain same idx, remove those where idx is not last:
        //     - 0&5 -> 5, 1&4 -> 4, 2&3 -> 3 
        //   - If several paths end with same idx, keep only biggest pathPatternSize.
        //     - 3&4&5 -> 5
        void shrinkSrcMap( void ){
        }

void createSrcMap(const uint8_t * table, const uint8_t * src, size_t slen){
    createIDPosTable(table, src, slen); // Get all ID positions in src ordered by increasing offset.
    memset(&srcMap, 0, sizeof(srcMap)); // Start with no path (PathCount=0).
    for( unsigned pti = 0; pti < IDPosTable.count; pti++ ){ // Loop over IDPosition table.
        IDPosition_t idPos = IDPosTable.item[pti]; // For each next idPos nnn:
        loc_t nnn_start = idPos.start;
        int IDPosAppended = 0;
        unsigned srcMapCount = srcMap.count;
        for( int k = 0; k <  srcMapCount; k++ ){ // Loop over all so far existing paths.
            if( srcMap.count > TIP_MAX_PATH_COUNT ){ // Create no new paths for this src buffer.
#if DEBUG
                printf( "srcMap is full (%d paths)\n", TIP_MAX_PATH_COUNT);
#endif
                IDPosAppended = 1; // Do not add any further paths.
                break;
            }
            path_t path = srcMap.path[k];                  // path is next path in srcMap.
            uint8_t pcnt = path.count;                     // pcnt is the number od IDPosTable indices in this path.
            uint8_t idx = path.posIndex[pcnt-1];           // idx is last IDPosTable index in this path.
            IDPosition_t lastIdPos = IDPosTable.item[idx]; // lastIdPos is the last (referenced) idPos in this path.
            uint8_t lll_Id = lastIdPos.id;                 // lll_Id is the ID we got from the IDPosTable.
            loc_t lll_start = lastIdPos.start;             // lll_start ist the start position of the pattern.
            uint8_t lll_len = IDPatternLength( lll_Id );   // lll_len is the length of the pattern to check.
            loc_t lll_limit = lll_start + lll_len;         // lll-limit is the first free position in this path.
            //case
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
                if( srcMap.count < TIP_MAX_PATH_COUNT ){
                    unsigned n = forkPath(k);
                    appendPosTableIndexToPath(n, pti);
#if DEBUG
                    printPath("fork: ", n);
#endif

                }else{
                    appendPosTableIndexToPath(k, pti);
#if DEBUG
                    printPath("appd: ", k);
#endif
                }
                IDPosAppended = 1; // idx is appended to at least one path now. 

            }
        }
        if( !IDPosAppended ){ // pti did not fit to any path, so lets create a new path for it.
            if( srcMap.count < TIP_MAX_PATH_COUNT ){ // Create no new paths for this src buffer.  
                unsigned nextIdx = srcMap.count;
                path_t * path = &srcMap.path[nextIdx];
                path->count = 1;         // One position table index is now in this new path.
                path->posIndex[0] = pti; // Write the pti (the first is naturally 0)
                srcMap.count++;          // We have one more path now.
#if DEBUG
                printPath(" new: ", nextIdx);
#endif
            }else{
#if DEBUG
                printf( "no new path possible: srcMap is full (%d paths)\n", srcMap.count);
#endif
            }
        }
        shrinkSrcMap();
    }
}

//! IDPosLength returns first offset after ID position idx.
static uint8_t IDPosLength(uint8_t idx){
    uint8_t id = IDPosTable.item[idx].id;
    uint8_t len = IDPatternLength( id );
    return len;
}

//! pathPatternSize returns sum of all pattern lengths in path pidx.
static loc_t pathPatternSize( unsigned pidx ){
    if( pidx >= srcMap.count ){
#if DEBUG
        printf( "ERROR: pidx %u >= %u srcMap.count\n", pidx, srcMap.count);
#endif
        return 0;
    }
    path_t path = srcMap.path[pidx];
    loc_t pathIdxCount = path.count;
    loc_t sum = 0;
    for( int i = 0; i < pathIdxCount; i++ ){
        loc_t idx = path.posIndex[i];
        sum += IDPosLength(idx);
    }
    return sum;
}

//! MinDstLengthPath returns srcMap path index which results in a shortest package.
//! TODO: Faster implemtation
static unsigned MinDstLengthPath(void){
    loc_t maxSum = 0;
    unsigned pidx = 0;
    for( unsigned i = 0; i < srcMap.count; i++ ){
        loc_t psiz = pathPatternSize(i);
        if( psiz > maxSum ){
            maxSum = psiz;
            pidx = i;
        }
    }
    // printPath( "first       psiz max ", pidx );
    // Maximum psiz results in smallest tip size only for min path length.
    loc_t minLen = TIP_SRC_BUFFER_SIZE_MAX;
    for( unsigned i = 0; i < srcMap.count; i++ ){
        if( maxSum == pathPatternSize(i) ){
            // printPath( "alternative psiz max ", i );
            loc_t plen = srcMap.path[i].count;
            if( plen < minLen ){
                minLen = plen;
                pidx = i;
                // printf( "pidx %u, minlen %u\n", pidx, minLen );
            }
        }
    }
#if DEBUG
    printf( "%u paths\n", srcMap.count );
    printPath( "selected    plen min ", pidx );
#endif
    return pidx;
}

//! selectUnreplacableBytes coppies all unreplacable bytes from src to dst.
//! It uses idPatTable and the path index pidx in the actual srcMap, which is linked to IDPosTable.
static size_t selectUnreplacableBytes( uint8_t * dst, unsigned pidx, const uint8_t * src, size_t slen ){
    path_t path = srcMap.path[pidx]; // This is the path we use. 
    loc_t count = path.count;        // The path contains count IDPosTable indicies.
    const uint8_t * srcNext = src;   // next position for src buffer read
    uint8_t * dstNext = dst;         // next position for dst buffer write
    loc_t * tidx = path.posIndex;    // Here are starting the IDPosTable indices.
    loc_t u8sum = 0;
    for( int i = 0; i < count; i++ ){
        IDPosition_t idPos = IDPosTable.item[tidx[i]];
        uint8_t id = idPos.id;
        const uint8_t * patternFrom = src + idPos.start; // pattern start in src buffer
        loc_t u8len = patternFrom - srcNext; // count of unreplacable bytes
        memcpy( dstNext, srcNext, u8len );
        uint8_t patlen = IDPatternLength( id );
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
static size_t createOutput( uint8_t * dst, unsigned pidx, const uint8_t * u7src, size_t u7len, const uint8_t * src ){
    path_t path = srcMap.path[pidx]; // This is the path we use. 
    loc_t plen = path.count;        // The path contains: count, IDPosTable index, IDPosTable index, ...
    const uint8_t * srcNext = src;   // next position for src buffer read
    uint8_t * dstNext = dst;         // next position for dst buffer write
    //loc_t * tidx = path.posIndex;    // Here are starting the IDPosTable indices.
    const uint8_t * u7Next = u7src;
    // printf( "PATH %d\n", pidx );
    // printPath( "DEBUG:", pidx );
    for( int i = 0; i < plen; i++ ){
        loc_t idx = path.posIndex[i];
        IDPosition_t idPos = IDPosTable.item[idx];
        uint8_t id = idPos.id;
        const uint8_t * patternFrom = src + idPos.start; // pattern start in src buffer
        loc_t u8len = patternFrom - srcNext; // count of unreplacable bytes
        uint8_t patlen = IDPatternLength( id );
        srcNext += patlen + u8len;
        // printf( "i %u: idx %u, id %u, start %u\n", i, idx, id, idPos.start );
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
static size_t shift87bit( uint8_t * lst, const uint8_t * src, size_t slen ){
    const uint8_t * u8 = src + slen; // first address behind src buffer
    uint8_t * dst = lst;             // destination address
    while( src < u8 ){
        uint8_t msb = 0x80;
        for( int i = 1; i < 8; i++ ){
            u8--;                    // next value address
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

static size_t buildTiPacket(uint8_t * dst, uint8_t * dstLimit, const uint8_t * table, const uint8_t * src, size_t slen){
    createSrcMap(table, src, slen);
    unsigned pidx = MinDstLengthPath(); // find minimum line
    printPath( "SELECT:", pidx );
    memset(dst, 0, dstLimit-dst);
    unsigned u8Count = selectUnreplacableBytes(dst, pidx, src, slen );
    unsigned u7Count = shift87bit( dstLimit-1, dst, u8Count );
    uint8_t * u7src = dstLimit - u7Count;

#if DEBUG
    printf( "MinDstLengthPath: %u, u8Count: %u, u7Count: %u\n", pidx, u8Count, u7Count );
#endif

    size_t pkgSize = createOutput( dst, pidx, u7src, u7Count, src );
    return pkgSize; // final ti package size
}

static void printPath( char * prefix, unsigned pidx ){
    path_t path = srcMap.path[pidx]; 
    int plen = path.count;
    loc_t psiz = pathPatternSize(pidx);
    printf( "%s%6u: psum%3d, pcnt%3d, idx:", prefix, pidx, psiz, plen);
    for( int k = 0; k < plen; k++ ){
        loc_t idx = path.posIndex[k];
        printf( " %d", idx);
    }
    loc_t last = 0;
    printf( ", pat:" );
    for( int k = 0; k < plen; k++ ){
        loc_t idx = path.posIndex[k];
        IDPosition_t idPos = IDPosTable.item[idx];
        uint8_t id = idPos.id;
        if( last < idPos.start ){
            printf( "~" );
        }else{
            printf( " " );
        }
        last = idPos.start + IDPatternLength( id );
        printPatternAsASCII(id);
    }
    
    printf( "\n" );
}

#if DEBUG

//  static void printSrcMap( void ){
//      printf( "srcMap: -----------\n");
//      for( unsigned i = 0; i < srcMap.count; i++ ){
//          printPath("      ", i);
//      }
//      printf( "-----------\n");
//  }

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

static void printIDPositionTable( void ){
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
