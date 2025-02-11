//! @file pack.c
//! @brief This is the tip pack code. Works also without unpack.c.
//! @details Written for ressources constraint embedded devices.
//! This tip code avoids heavy stack usage by using static buffers and is therefore not re-entrant.
//! This implementation is coded for speed in favour RAM usage. 
//! If RAM usage matters, the replace list r could be a bit array at the end of the destination buffer just to mark the unreplacable bytes.
//! In a loop then the packed data can get constructed directly into the destination buffer by searching for the pattern a second time.
//! It is possible to use different tables at the same time, but the code needs to be changed a bit then.
//! @author thomas.hoehenleitner [at] seerose.net

#include "tip.h"
#include "tipInternal.h"
#include "trice.h"

size_t tip( uint8_t* dst, const uint8_t * src, size_t len ){
    return TiPack( dst, idTable, src, len );
}

//! @brief tip encodes src buffer with size len into dst buffer and returns encoded len.
//! @details For the tip encoding it uses the linked idTable.c object.
// - Some bytes groups in the src buffer are replacable with IDs 0x01...0x7f and some not.
// - The replace list r holds the replace information.
// - The unreplacable bytes are collected into a buffer.
size_t TiPack( uint8_t* dst, const uint8_t * table, const uint8_t * src, size_t slen ){
    //if( slen < 6 ){
    //    memcpy(dst, src, slen);
    //    return slen;
    //}
    int rcount;
    replace_t * rlist = buildReplaceList(&rcount, table, src, slen);
    for( int i = 0; i < rcount; i++ ){
        trice( "d:%2d: bo=%d, sz=%d, id=%02x\n", i, rlist[i].bo, rlist[i].sz, rlist[i].id);
    }
    
    // All unreplacable bytes are stretched inside to 7-bit units. This makes the data a bit longer.
    static uint8_t ur[TIP_SRC_BUFFER_SIZE_MAX*8/7+1]; 
    size_t ubSize = collectUnreplacableBytes( ur, rlist, rcount, src );
    uint8_t * urLimit = &ur[sizeof(ur)]; // first address after the ur buffer.
    size_t urSize = shift87bit( urLimit - 1, ur, ubSize );
    uint8_t * u7 = urLimit - urSize;
    size_t tipSize = generateTipPacket( dst, u7, urSize, rlist, rcount );
    return tipSize;
}

//! @brief newReplacableList is called when a new unpacked buffer arrived.
//! @details It returns always the same static object to avoid memory allocation.
//! @param slen is the source buffer size.
//! @retval is a pointer to the replace list.
replace_t * newReplaceList(size_t slen){
    //static replaceList_t r; // replace list
    static replace_t list[TIP_SRC_BUFFER_SIZE_MAX/2 + 2]; //!< The whole src buffer could be replacable with 2-byte pattern.
    // The first 2 elements are initialized as boders.
    list[0].bo = 0; // byte offset start
    list[0].sz = 0; // size
    list[0].id = 0; // no replacement
    // From (r[0].bo + r[0].sz) to r[1].bo is the first hey stack.
    list[1].bo = slen; // byte offset limit
    list[1].sz = 0; // needed as end marker
    list[1].id = 0; // no replacement
    return list;
};

replace_t * buildReplaceList(int * rcount, const uint8_t * table, const uint8_t * src, size_t slen){
    replace_t * rlist = newReplaceList(slen);
    *rcount = 2;
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
            repeat:
            hay = src + rlist[k].bo + rlist[k].sz;
            hlen = rlist[k+1].bo - rlist[k].bo - rlist[k].sz;  
            uint8_t * loc = memmem( hay, hlen, needle, nlen ); // search hay for the needle
            if( loc ){ // found, id is the replace byte.
                offset_t offset = loc - src; // offset is the needle (=pattern) position.
                replaceableListInsert( rlist, rcount, k, id, offset, nlen );
                goto repeat; // Same k and same  (needle) needs processing again but in the next hay stack.
            } // The r insert takes part inside the already processed rs.
            k++;
        }while( hay+hlen < src+slen );
    }
    return rlist;
}

// generateTipPacket uses r and u to build the tip.
//! @param dst start of result data
//! @param u7 start of buffer with 7 lsbits btes
//! @param u7Size count of 7 lsbits bytes
//! @param rl replace list
//! @retval length of tip packet
size_t generateTipPacket( uint8_t * dst, uint8_t * u7, size_t u7Size, replace_t* rlist, int rcount ){ 
    size_t tipSize = 0;
    int k = 0;  // Traverse rlist to find relacement pattern.
    do { // r->list[k] is done here, we need to fill the space and insert r[k+1] pattern.
        int uBytes = rlist[k+1].bo - (rlist[k].bo + rlist[k].sz);
        while(u7Size > 0 && uBytes > 0){
            // Each inserted u7 byte is also a place holder for a u8 byte.
            // u7 count is >= u8 count, sowe can cover all u8 positions.
            // The u7 we have more, we append ant the end.
            *dst++ = *u7++;
            uBytes--;
            u7Size--;
            tipSize++;
        }
        k++;
        size_t sz = rlist[k].sz; // Size of next replace.
        if( sz == 0 ){
            continue; // no more replaces, but some unreplacable may still exist.
        }
        *dst++ = rlist[k].id;
        tipSize++;
    }while(k < rcount-1);
    while(u7Size > 0){ // append remaining u7 bytes
        *dst++ = *u7++;
        u7Size--;
        tipSize++;
    }
    return tipSize;
}

//! @brief replaceableListInsert extends r in an ordered way.
//! @param r ist the replace list.
//! @param k is the position after where to insert.
//! @param id is the replace byte for the location.
//! @param offset is the location to be extended with.
//! @param sz is the replace pattern size.
void replaceableListInsert(replace_t * rlist, int * rcount, int k, uint8_t id, offset_t offset, uint8_t sz){
    k++;
    memmove( &(rlist[k+1]), &(rlist[k]), (*rcount-k)*sizeof(replace_t));
    rlist[k].id = id;
    rlist[k].bo = offset;
    rlist[k].sz = sz;
    (*rcount)++;
}

//! collectUnreplacableBytes uses information in rl to construct dst (->u) from src.
//! @param dst is destination address.
//! @param r is the replace list. Its holes are the unreplacable bytes information.
//! @param src is the data buffer containing repacable and unreplacable bytes.
//! @retval is the dst size.
size_t collectUnreplacableBytes( uint8_t * dst, replace_t * rlist, int rcount, const uint8_t * src ){
    size_t dstCount = 0;
    for( int k = 0; k < rcount - 1; k++ ){
        offset_t offset = rlist[k].bo + rlist[k].sz;
        const uint8_t * addr = src + offset;
        size_t len = rlist[k+1].bo - offset; // gap
        memcpy( dst + dstCount, addr, len );
        dstCount += len;
    }
    return dstCount;
}
