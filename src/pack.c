8//! @file pack.c
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

/*static*/ replace_t * buildReplaceList(int * rcount, const uint8_t * table, const uint8_t * src, size_t slen);
static size_t collectUnreplacableBytes( uint8_t * dst, replace_t * rlist, int rcount, const uint8_t * src );
/*static*/ size_t shift87bit( uint8_t* lst, const uint8_t * src, size_t slen );
static void initGetNextPattern( const uint8_t * table );
static void getNextPattern(const uint8_t ** pt, size_t * sz );
static replace_t * newReplaceList(offset_t slen);
static void replaceableListInsert( replace_t * r, int * rcount, int k, uint8_t by, offset_t offset, uint8_t sz );
static size_t generateTipPacket( uint8_t * dst, uint8_t * u7, uint32_t u7Size, replace_t * rlist, int rcount );

size_t tip( uint8_t* dst, const uint8_t * src, size_t len ){
    size_t dstLimit = ((18725*len)>>14)+1;  // The max possible dst size is len*8/7+1 or ((len*65536*8/7)>>16)+1;
    return tiPack( dst, dstLimit, idTable, src, len );
}

//! @brief tiPack encodes src buffer with size len into dst buffer and returns encoded len.
//! @details For the tip encoding it uses the linked idTable.c object.
// - Some bytes groups in the src buffer are replacable with IDs 0x01...0x7f and some not.
// - The rlist r holds the replace information. Additionally the dst buffer is prefilled with IDs from both sides.
// - Example: dst = 5, 6, 0, ..., 0, 2, 3, 4; and rlist = 00111111000111000001111110111
// - ID 5 has 2 and ID 6 4 bytes, so ID 2 and 4 have 3 bytes and ID 3 has 5 bytes.
// - The unreplacable bytes are collected into a buffer.
size_t tiPack( uint8_t* dst, size_t dstLimit, const uint8_t * table, const uint8_t * src, size_t slen ){
    if( slen > TIP_SRC_BUFFER_SIZE_MAX ){
        return 0;
    }
    memset(dst, 0, dstLimit);
    uint8_t * rlist = newReplaceBits();
    setReplaceBits(dst, dstLimit, rlist, table, src, slen);
    // All unreplacable bytes are stretched inside to 7-bit units. This makes the data a bit longer.
    static uint8_t ur[TIP_SRC_BUFFER_SIZE_MAX*8u/7u+1]; 
    size_t ubSize = getUnreplacableBytes(ur, rlist, src );
    uint8_t * urLimit = &ur[sizeof(ur)]; // first address after the ur buffer.
    size_t urSize = shift87bit( urLimit - 1, ur, ubSize );
    uint8_t * u7 = urLimit - urSize;
    size_t tipSize = buildTipPacket( dst, dstLimit, urSize, rlist );
    return tipSize;
}

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

//! @brief newReplaceBits is called when a new unpacked buffer arrived.
//! @details It returns always the same static object to avoid memory allocation.
//! @retval is a pointer to the replace list.
//! If the N-th bit in the list is set, it signals that the N-th byte in the unpacked buffer is replacable.
static uint8_t * newReplaceBits(void){
    static uint8_t list[TIP_SRC_BUFFER_SIZE_MAX/8];
    memset(list, 0, sizeof(list));
    return list;
};

// setBits sets in bitList the bits from pos to pos+len (exclusive).
// Example: uint8_t x = 0; setbits(&x, 3, 2); -> x == 0b00011000.
void setBits( uint8_t * bitList, size_t pos, size_t len ){
    // todo
}

// setReplaceBits starts with buf=src and tries to find biggest matching pattern from table at buf AND bufLimit-nlen.
// If a pattern was found at buf, buf is incremented by found pattern size.
// If a pattern was found at bufLimit-nlen, bufLimit is decremented by found pattern size.
// If a pattern was found at front or back we start over.
// If none found, we increment buf by 1 and decrement bufLimit and start over.
// It is possible, the same pattern is found again at the same place, but we do not care, we just set the bits.
void setReplaceBits(uint8_t * dst, uint8_t * dstLimit, uint8_t * rlist, const uint8_t * table, const uint8_t * src, size_t slen){
    const uint8_t * buf = src;
    const uint8_t * bufLimit = src + slen;
    int frontSearch = 1;
    int backSearch = 1;
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
            offset_t offset = buf - src; // relative pattern position
            setBits( rlist, offset, nlen );
            *dst++ = id;
            buf += nlen;
            if( !(buf < src + slen)){
                frontSearch = 0;
            }
        }
        if( backSearch && 0 == strncmp((void*)(bufLimit-nlen), (void*)needle, nlen) ){ // match at buf back
            backMatch = 1;
            offset_t offset = bufLimit-nlen - src; // relative pattern position
            setBits( rlist, offset, nlen );
            *--dstLimit = id;
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
    // Arriving here means that no table pattern fits to buf or bufLimit
    u = *buf++;
    *dst++ = 0x80 & u;
    if( !(buf < src + slen)){
        frontSearch = 0;
    }
    u = *bufLimit--;
    *--dstLimit = 0x80 & u;
    if( !(src < bufLimit) ){
        backSearch = 0;
    }
    if( frontSearch || backSearch ){
        goto repeat;
    }
}

// generateTipPacket uses r and u to build the tip.
//! @param dst start of result data
//! @param u7 start of buffer with 7 lsbits btes
//! @param u7Size count of 7 lsbits bytes
//! @param rl replace list
//! @retval length of tip packet
static size_t buildTipPacket( uint8_t * dst, uint8_t * u7, uint32_t u7Size, replace_t* rlist ){ 
    size_t packetSize = 0;

    uint8_t * rlist = newReplaceBits();
    setReplaceBits(rlist, table, src, slen);

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
            packetSize++;
        }
        k++;
        uint8_t sz = rlist[k].sz; // Size of next replace.
        if( sz == 0 ){
            continue; // no more replaces, but some unreplacable may still exist.
        }
        *dst++ = rlist[k].id;
        packetSize++;
    }while(k < rcount-1);
    while(u7Size > 0){ // append remaining u7 bytes
        *dst++ = *u7++;
        u7Size--;
        packetSize++;
    }
    return packetSize;
}

//! @brief replaceableListInsert extends r in an ordered way.
//! @param rlist ist the replace list.
//! @param k is the rlist position after where to insert.
//! @param id is the replace byte for the location.
//! @param offset is the location to be extended with.
//! @param sz is the replace pattern size.
static void replaceableListInsert(replace_t * rlist, int * rcount, int k, uint8_t id, offset_t offset, uint8_t sz){
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
static size_t collectUnreplacableBytes( uint8_t * dst, replace_t * rlist, int rcount, const uint8_t * src ){
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

//! shift87bit transforms slen 8-bit bytes in src to 7-bit units.
//! @param src is the bytes source buffer.
//! @param slen is the 8-bit byte count.
//! @param lst is the last address inside the dst buffer.
//! @retval is count of 7-bit bytes after operation. 
//! @details The dst buffer is filled from the end.Thas allows to do an in-buffer conversion.
//! The destination address is computable afterwards: dst = lim - retval.
//! lim is allowed to be "close" behind buf + slen, thus making in-place conversion possible.
//! Example: slen=17, limit=24
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

#if 1 // 0 side

// buildReplaceList starts with biggest table pattern and searches for matches.
replace_t * buildReplaceList(int * rcount, const uint8_t * table, const uint8_t * src, size_t slen){
    replace_t * rlist = newReplaceList(slen);
    *rcount = 2;
    initGetNextPattern(table);
    for( int id = 1; id < 0x80; id++ ){ // traverse the ID table.
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
                offset_t offset = loc - src; // offset is the relative needle (=pattern) position.
                replaceableListInsert( rlist, rcount, k, id, offset, nlen );
                goto repeat; // Same k and same  (needle) needs processing again but in the next hay stack.
            } // The r insert takes part inside the already processed rs.
            k++;
        }while( hay+hlen < src+slen );
    }
    return rlist;
}

#endif

#if 0 // 1 side

// buildReplaceList starts with src and tries to find biggest matching pattern from table at buf = src.
// If one was found, buf is incremented by found pattern size and we start over.
// If none found, we increment buf by 1 and start over.
replace_t * buildReplaceList(int * rcount, const uint8_t * table, const uint8_t * src, size_t slen){
    replace_t * rlist = newReplaceList(slen);
    *rcount = 2;
    int k = 0;
    const uint8_t * buf = src;
repeat:
    initGetNextPattern(table);
    for( int id = 1; id < 0x80; id++ ){ // traverse the ID table.    
        // get the next pattern
        const uint8_t * needle = NULL;
        size_t nlen;
        getNextPattern( &needle, &nlen );
        if( nlen == 0 ){ // end of table if less 127 IDs
            break; 
        }
        if( strncmp((void*)buf, (void*)needle, nlen) ){ // no match
            continue; // with next pattern
        }
        // found, id is the replace byte.
        offset_t offset = buf - src; // relative pattern position
        replaceableListInsert( rlist, rcount, k, id, offset, nlen );
        k++;
        buf += nlen;
        if( !(buf < src + slen) ){
            return rlist;
        }
    }
    buf++;
    if( buf < src + slen ){
        goto repeat;
    }
    return rlist;
}

#endif

#if 0 // 2 sides

// replaceableListIndex checks. if at offset a pattern with size sz is already known with id.
// If so, -1 is returned.
// Otherwise the k is returned after which the insertion is ok.
static int replaceableListIndex(replace_t * rlist, int * rcount, uint8_t id, offset_t offset, uint8_t sz){
    for( int k = 1; k < *rcount; k++ ){
        if( rlist[k].bo < offset ){
            continue;
        }
        if( rlist[k].bo == offset ){
            //while( rlist[k].sz != sz );
            //while( rlist[k].id != id );
            return -1;
        }
        return k-1; // rlist[k].bo > offset
    }
    //for(;;);
    return -1;
}

// buildReplaceList starts with src and tries to find biggest matching pattern from table at buf = src.
// If one was found, buf is incremented by found pattern size and we start over.
// If none found, we increment buf by 1 and start over.
/*static*/ replace_t * buildReplaceList(int * rcount, const uint8_t * table, const uint8_t * src, size_t slen){
    replace_t * rlist = newReplaceList(slen);
    *rcount = 2;
    int k = 0;
    const uint8_t * buf = src;
    const uint8_t * bufLimit = src + slen;
    int startSearch = 1;
    int endSearch = 1;
repeat:
    initGetNextPattern(table);
    for( int id = 1; id < 0x80; id++ ){ // traverse the ID table.    
        // get the next pattern
        const uint8_t * needle = NULL;
        size_t nlen;
        getNextPattern( &needle, &nlen );
        if( nlen == 0 ){ // end of table if less 127 IDs
            break; 
        }
        if( startSearch && 0 == strncmp((void*)buf, (void*)needle, nlen) ){ // match at buf start  
            offset_t offset = buf - src; // relative pattern position
            k = replaceableListIndex( rlist, rcount, id, offset, nlen );
            if( k != -1 ){
                replaceableListInsert( rlist, rcount, k, id, offset, nlen ); // id is the replace byte.
                //k++;
            }
            buf += nlen;
            if( !(buf < src + slen)){
                startSearch = 0;
            }
            goto repeat;
        }
        if( endSearch && 0 == strncmp((void*)(bufLimit-nlen), (void*)needle, nlen) ){ // match at buf end
            offset_t offset = bufLimit-nlen - src; // relative pattern position
            k = replaceableListIndex( rlist, rcount, id, offset, nlen );
            if( k != -1 ){
                replaceableListInsert( rlist, rcount, k, id, offset, nlen ); // id is the replace byte.
                //k++;
            }
            bufLimit -= nlen;
            if( !(src < bufLimit) ){
                endSearch = 0;
            }
            goto repeat;            
        }
        continue; // with next pattern  
    }
    buf++;
    bufLimit--;
    if( !(buf < src + slen)){
        startSearch = 0;
    }
    if( !(src < bufLimit) ){
        endSearch = 0;
    }
    if( startSearch || endSearch ){
        goto repeat;
    }
    return rlist;
}

#endif
#if 0

//! @brief newReplacableList is called when a new unpacked buffer arrived.
//! @details It returns always the same static object to avoid memory allocation.
//! @param slen is the source buffer size.
//! @retval is a pointer to the replace list.
static replace_t * newReplaceList(offset_t slen){
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

#endif
