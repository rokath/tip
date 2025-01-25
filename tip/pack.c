//! @file pack.c
//! @brief Written for ressources constraint embedded devices. This tip() code is not re-entrant avoiding heavy stack usage.
//! @details This implementation is optimized for speed in favour RAM usage. 
//! If RAM usage matters r could be a bit array at the end of the destination buffer just to mark the unreplacable bytes.
//! In a loop then the packed data can get constructed directly into the destination buffer by searching for the pattern a second time.

#include <strings.h>
#include <stddef.h>
#include "tip.h"

static size_t shift87bit( uint8_t * buf, size_t len, size_t limit );
static inline void rInit(size_t len);
static inline void rInsert( int k, uint8_t by, offset_t offset, uint8_t sz );
static void collectUnreplacableBytes( uint8_t const * src );

//! @brief r is the replacement list. It cannot get more elements.
//! The space between 2 rs is a hay stack.
static replacement_t r[TIP_SRC_BUFFER_SIZE_MAX/2 + 2];
static int rc; //!< replacement count

//! @brief u holds all unreplacable bytes from src. It cannot get longer.
//! @details All unreplacable bytes are stretched inside to 7-bit units.
//! This makes the data a bit longer.
static uint8_t u[TIP_SRC_BUFFER_SIZE_MAX*8/7+1];

//! @brief uCount is the number of valid bytes inside u.
static size_t uCount = 0;

//! @brief tip encodes src buffer with size len into dst buffer and returns encoded len.
//! @details For the tip encoding it uses the linked tipTable.c object.
size_t tip( uint8_t* dst, uint8_t const * src, size_t len ){
    rInit(len);
    for( int i = 0; i < TipTableLength; i++ ){
        // get needle (the next pattern)
        uint8_t * needle = TipTable[i].pt;
        size_t nlen = TipTable[i].sz;
        // Traverse r to find hey stacks.
        int k = 0;
        uint8_t const * hay;
        size_t hlen;
        do{ // get next hay stack
            hay = src + r[k].bo + r[k].sz;
            hlen = r[k+1].bo - r[k].bo - r[k].sz;
            // search the needle
            uint8_t * loc = memmem( hay, hlen, needle, nlen );
            if( loc ){ // found
                uint8_t by = TipTable[i].by; // by is the replacement byte.
                offset_t offset = loc - src; // offset is the needle (=pattern) position.
                rInsert(k, by, offset, nlen );
                k--; // Same k needs processing again.
            } // The r insert takes part inside the already processed rs.
            k++;
        }while( hay+hlen < src+len );
    }
    // Some bytes groups in the src buffer are replacable with 0x01...0xFF and some not.
    // The replacement list r contains now the replacement information.
    // Lets collect the unreplacable bytes into a buffer now.
    collectUnreplacableBytes( src );
    uCount = shift87bit( u, uCount, sizeof(u) );
    return generateTipPacket( dst );
}

//! getReplacementPattern returns a pointer or NULL. 
static uint8_t * getReplacementPattern( uint8_t by ){
    for(int i = 0; i < TipTableLength){
        if( i == by ){
            return TipTable[i].pt;
        }
    }
    return NULL
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
        uint8_t sz = r[k+1].sz; // Size of next replacement.
        if( sz == 0 ){
            k++; // no more replacements, but some unreplacable may still exist.
            continue;
        }
        uint8_t * pt = getReplacementPattern( r[k+1].by );
        while( sz-- ){
            *dst++ = *pt++;
        }
    }while(k < rc -1);
}


//! shift87bit transforms len 8-bit bytes in buf to 7-bit units.
//! @param buf is a byte buffer. It is destroyed during operation.
//! @param len is the 8-bit byte count.
//! @param limit is the max byte count fitting into buf (limit > len*8/7)
//! @retval is count 7-bit bytes
//! @details buf is filled from the end (=buf+limit)
//! The destination is computable afterwards: dst = buf + limit - retval.
//! Example: len=17, limit=24
//!       (buf) <---              17                    --->  [n8]                 [n7]
//! len=17: b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 b8 __ __ __ __ __ __ __
//! ret=20: __ __ __ __ m7 b7 b7 b7 m7 b7 b7 b7 b7 b7 b7 b7 m7 b7 b7 b7 b7 b7 b7 b7
//!                   (dst) <---                     20                       --->
static size_t shift87bit( uint8_t * buf, size_t len, size_t limit ){
    int n7 = limit; // n7 is buf index limit.
    for( int n8 = len; n8 > 0; ){ // n8 is buf data limit
        uint8_t msb = 0;
        for( int i = 7; i > 0; i-- && n8 > 0){
            msb |= (0x80 & buf[--n8])>>i; // Store the MSB of the current last byte at bit position
            buf[--n7] = 0x80 | buf[n8]; // the last byte 7 LSBs and MSB=1 to the end
        }
    }
}

//! @brief rInit is called when a new unpacked buffer arrived.
//! @param len is the source buffer size.
static inline void rInit(size_t len){
    rc = 2; // The first 2 elements are initialized as boders.
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
    rc++;
    memmove( &r[k+1], &r[k], (rc-k)*sizeof(replacement_t));
    r[k].by = by;
    r[k].bo = offset;
    r[k].sz = sz;
}

static void collectUnreplacableBytes( uint8_t const * src ){
    for( int k = 0; k < rc - 1; k++ ){
        offset_t offset = r[k].bo + r[k].sz;
        uint8_t * addr = src + offset;
        size_t len = r[k+1].bo - offset; // gap
        memcpy( u + uCount, addr, len );
        uCount += len;
    }
}

