//! @file pack.c
//! @brief This is the tip pack code. Works also without unpack.c.
//! @details Written for ressources constraint embedded devices.
//! This tip() code avoids heavy stack usage by using static buffers and is therefore not re-entrant.
//! This implementation is optimized for speed in favour RAM usage. 
//! If RAM usage matters, the replacement list r could be a bit array at the end of the destination buffer just to mark the unreplacable bytes.
//! In a loop then the packed data can get constructed directly into the destination buffer by searching for the pattern a second time.
//! @author thomas.hoehenleitner [at] seerose.net

#include "tipInternal.h"

static inline void rInit(size_t len);
static inline void rInsert( int k, uint8_t by, offset_t offset, uint8_t sz );
//static size_t shift87bit( uint8_t * buf, size_t len, size_t limit );
static void collectUnreplacableBytes( uint8_t const * src );
static size_t generateTipPacket( uint8_t * dst );

//! @brief r is the replacement list. It cannot get more elements.
//! The space between 2 rs is a hay stack.
static replacement_t r[TIP_SRC_BUFFER_SIZE_MAX/2 + 2];

 //! @brief rCount is replacement count inside r.
static int rCount;

//! @brief u holds all unreplacable bytes from src. It cannot get longer.
//! @details All unreplacable bytes are stretched inside to 7-bit units.
//! This makes the data a bit longer.
static uint8_t u[TIP_SRC_BUFFER_SIZE_MAX*8/7+1];

//! @brief uCount is the number of valid bytes inside u.
static size_t uCount = 0;

//! @brief tip encodes src buffer with size len into dst buffer and returns encoded len.
//! @details For the tip encoding it uses the linked tipTable.c object.
size_t tip( uint8_t* dst, uint8_t const * src, size_t len ){
    if( len < 16 ){
        memcpy(dst, src, len);
        return len;
    }
    resetPattern();
    rInit(len);
    for( int id = 1; id < 0x7f; id++ ){
        // get biggest needle (the next pattern)
        uint8_t * needle = NULL;
        size_t nlen;
        getNextPattern( &needle, &nlen );
        if( nlen == 0 ){
            break; 
        }
        // Traverse r to find hey stacks.
        int k = 0;
        uint8_t const * hay;
        size_t hlen;
        do{ // get next hay stack
            hay = src + r[k].bo + r[k].sz;
            hlen = r[k+1].bo - r[k].bo - r[k].sz;
            // search the needle
            uint8_t * loc = memmem( hay, hlen, needle, nlen );
            if( loc ){ // found, id is the replacement byte.
                offset_t offset = loc - src; // offset is the needle (=pattern) position.
                rInsert(k, id, offset, nlen );
                k--; // Same k needs processing again.
            } // The r insert takes part inside the already processed rs.
            k++;
        }while( hay+hlen < src+len );
    }
    // Some bytes groups in the src buffer are replacable with 0x01...0xFF and some not.
    // The replacement list r contains now the replacement information.
    // Lets collect the unreplacable bytes into a buffer now.
    collectUnreplacableBytes( src );
   // uCount = shift87bit( u, uCount, sizeof(u) );
    return generateTipPacket( dst );
}

// generateTipPacket uses r and u to build the tip.
static size_t generateTipPacket( uint8_t * dst ){ 
    uint8_t * u7 = u + sizeof(u) - uCount; // unreplacable to 7-bit converted bytes
    // Traverse r to find relacement pattern.
    int k = 0;
    do { // r[k] is done here, we need to fill the space and insert r[k+1] pattern.
        int uBytes = r[k+1].bo - (r[k].bo + r[k].sz);
        while(uCount && uBytes--){
            // Each inserted u7 byte is also a place holder for a u8 byte.
            // u7 count is >= u8 count, so we can cover all u8 positions.
            // The u7 we have more, we append ant the end.
            *dst++ = 0x80 | *u7++; // Set msb as unreplacable marker.
        }
        size_t sz = r[k+1].sz; // Size of next replacement.
        if( sz == 0 ){
            k++; // no more replacements, but some unreplacable may still exist.
            continue;
        }
        uint8_t * pt = NULL;
        getPatternFromId( r[k+1].by, &pt, &sz );
        while( sz-- ){
            *dst++ = *pt++;
        }
    }while(k < rCount -1);
    return 123;
}

//! @brief rInit is called when a new unpacked buffer arrived.
//! @param len is the source buffer size.
static inline void rInit(size_t len){
    rCount = 2; // The first 2 elements are initialized as boders.
    r[0].bo = 0;
    r[0].sz = 0; // r[0].by is never used. 
    // From (r[0].bo + r[0].sz) to r[1].bo is the first hey stack.
    r[1].bo = len;
    r[1].sz = 0; // needed as end marker. r[1].by is never used. 
};

//! @brief rInsert extends r in an ordered way.
//! @param k is the position after where to insert.
//! @param by is the replacement byte for the location.
//! @param offset is the location to be extended with.
//! @param sz is the replacement pattern size.
static inline void rInsert( int k, uint8_t by, offset_t offset, uint8_t sz ){
    k++;
    rCount++;
    memmove( &r[k+1], &r[k], (rCount-k)*sizeof(replacement_t));
    r[k].by = by;
    r[k].bo = offset;
    r[k].sz = sz;
}


// collectUnreplacableBytes uses information in r to construct u from src.
static void collectUnreplacableBytes( uint8_t const * src ){
    for( int k = 0; k < rCount - 1; k++ ){
        offset_t offset = r[k].bo + r[k].sz;
        uint8_t const * addr = src + offset;
        size_t len = r[k+1].bo - offset; // gap
        memcpy( u + uCount, addr, len );
        uCount += len;
    }
}
