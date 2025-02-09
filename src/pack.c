//! @file pack.c
//! @brief This is the tip pack code. Works also without unpack.c.
//! @details Written for ressources constraint embedded devices.
//! This tip code avoids heavy stack usage by using static buffers and is therefore not re-entrant.
//! This implementation is coded for speed in favour RAM usage. 
//! If RAM usage matters, the replacement list r could be a bit array at the end of the destination buffer just to mark the unreplacable bytes.
//! In a loop then the packed data can get constructed directly into the destination buffer by searching for the pattern a second time.
//! It is possible to use different tables at the same time, but the code needs to be changed a bit then.
//! @author thomas.hoehenleitner [at] seerose.net

#include "tipInternal.h"

size_t tip( uint8_t* dst, const uint8_t * src, size_t len ){
    return TiPack( dst, idTable, src, len );
}

//! @brief tip encodes src buffer with size len into dst buffer and returns encoded len.
//! @details For the tip encoding it uses the linked idTable.c object.
// - Some bytes groups in the src buffer are replacable with IDs 0x01...0x7f and some not.
// - The replacement list r holds the replacement information.
// - The unreplacable bytes are collected into a buffer.
size_t TiPack( uint8_t* dst, const uint8_t * table, const uint8_t * src, size_t slen ){
    if( slen < 16 ){
        memcpy(dst, src, slen);
        return slen;
    }   
    replaceList_t * r = buildReplacementList(table, src, slen);
    // All unreplacable bytes are stretched inside to 7-bit units. This makes the data a bit longer.
    static uint8_t ur[TIP_SRC_BUFFER_SIZE_MAX*8/7+1]; 
    size_t ubSize = collectUnreplacableBytes( ur, r, src );
    uint8_t * urLast = &(ur[sizeof(ur)-1]); // last address inside the ur buffer.
    size_t urSize = shift87bit( urLast, ur, ubSize );
    uint8_t * u7 = urLast - urSize;
    size_t tipSize = generateTipPacket( dst, u7, urSize, r );
    return tipSize;
}

//! @brief newReplacableList is called when a new unpacked buffer arrived.
//! @details It returns always the same staitc object to avoid memory allocation.
//! @param slen is the source buffer size.
//! @retval r is a pointer to the replacement list.
replaceList_t * newReplacableList(size_t slen){
    static replaceList_t r; // replace list
    r.count = 2; // The first 2 elements are initialized as boders.
    r.list[0].bo = 0; // byte offset
    r.list[0].sz = 0; // r[0].by is never used. 
    // From (r[0].bo + r[0].sz) to r[1].bo is the first hey stack.
    r.list[1].bo = slen; // byte offset
    r.list[1].sz = 0; // needed as end marker. r[1].by is never used.
    return &r;
};

replaceList_t * buildReplacementList( const uint8_t * table, const uint8_t * src, size_t slen){
    replaceList_t * r = newReplacableList(slen);
    initGetNextPattern(table);
    for( int id = 1; id < 0x80; id++ ){ // traverse te table.
        // get biggest needle (the next pattern)
        const uint8_t * needle = NULL;
        size_t nlen;
        getNextPattern( &needle, &nlen );
        if( nlen == 0 ){ // end of table if less 127 IDs
            break; 
        }
        // Traverse r to find hey stacks.
        int k = 0;
        const uint8_t * hay;
        size_t hlen;
        do{ // get next hay stack
            hay = src + r->list[k].bo + r->list[k].sz;
            hlen = r->list[k+1].bo - r->list[k].bo - r->list[k].sz;  
            uint8_t * loc = memmem( hay, hlen, needle, nlen ); // search hay for the needle
            if( loc ){ // found, id is the replacement byte.
                offset_t offset = loc - src; // offset is the needle (=pattern) position.
                replaceableListInsert( r, k, id, offset, nlen );
                k--; // Same k needs processing again.
            } // The r insert takes part inside the already processed rs.
            k++;
        }while( hay+hlen < src+slen );
    }
    return r;
}

// generateTipPacket uses r and u to build the tip.
//! @param dst start of result data
//! @param u7 start of buffer with 7 lsbits btes
//! @param u7Size count of 7 lsbits bytes
//! @param rl replacement list
//! @retval length of tip packet
size_t generateTipPacket( uint8_t * dst, uint8_t * u7, size_t u7Size, replaceList_t* r ){ 
    size_t tipSize = 0;
    int k = 0;  // Traverse r to find relacement pattern.
    do { // r->list[k] is done here, we need to fill the space and insert r[k+1] pattern.
        int uBytes = r->list[k+1].bo - (r->list[k].bo + r->list[k].sz);
        while(u7Size-- && uBytes--){
            // Each inserted u7 byte is also a place holder for a u8 byte.
            // u7 count is >= u8 count, sowe can cover all u8 positions.
            // The u7 we have more, we append ant the end.
            *dst++ = *u7++;
            tipSize++;
        }
        size_t sz = r->list[k+1].sz; // Size of next replacement.
        if( sz == 0 ){
            k++; // no more replacements, but some unreplacable may still exist.
            continue;
        }
        *dst++ = r->list[k+1].id;
        tipSize++;
    }while(k < r->count -1);
    return tipSize;
}

//! @brief replaceableListInsert extends r in an ordered way.
//! @param r ist the replacement list.
//! @param k is the position after where to insert.
//! @param id is the replacement byte for the location.
//! @param offset is the location to be extended with.
//! @param sz is the replacement pattern size.
void replaceableListInsert(replaceList_t * r, int k, uint8_t id, offset_t offset, uint8_t sz){
    k++;
    r->count++;
    memmove( &(r->list[k+1]), &(r->list[k]), (r->count-k)*sizeof(replaceList_t));
    r->list[k].id = id;
    r->list[k].bo = offset;
    r->list[k].sz = sz;
}

//! collectUnreplacableBytes uses information in rl to construct dst (->u) from src.
//! @param dst is destination address.
//! @param r is the replacement list. Its holes are the unreplacable bytes information.
//! @param src is the data buffer containing repacable and unreplacable bytes.
//! @retval is the dst size.
size_t collectUnreplacableBytes( uint8_t * dst, replaceList_t * r, const uint8_t * src ){
    size_t dstCount = 0;
    for( int k = 0; k < r->count - 1; k++ ){
        offset_t offset = r->list[k].bo + r->list[k].sz;
        const uint8_t * addr = src + offset;
        size_t len = r->list[k+1].bo - offset; // gap
        memcpy( dst + dstCount, addr, len );
        dstCount += len;
    }
    return dstCount;
}
