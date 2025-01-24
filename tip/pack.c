
#include <strings.h>
#include <stddef.h>
#include "tip.h"

//! @brief replacement_t is a replacement type descriptor.
typedef struct {
    uint16_t bo; // bo is the buffer offset, where replacement size starts.
    uint8_t  sz; // sz is the replacement size (2-8).
    uint8_t  by; // by is the replacement byte 0x01 to 0xff.
} replacement_t;

//! @brief rpp is the replacement list. It cannot get more elements.
static replacement_t rp[TIP_SRC_BUFFER_SIZE_MAX/2];

//! @brief rpInit is called when a new unpacked buffer arrived.
void rpInit( void ){
    rc = 2; // The first 2 element are initialized with a stopper.
    rp[0].bo = 0; // From idx to rp[0].bo is first hey stack.
    rp[0].sz = 0; // Dummy entry to get started.
    // The space between 2 rps is a hay stack.
    rp[1].bo = TIP_SRC_BUFFER_SIZE_MAX; // From idx to rp[0].bo is first hey stack.
    rp[1].sz = 0; // Dummy entry to get started.
}

//! @brief tip encodes src buffer with size len into dst buffer and returns encoded len.
//! @details For the tip encoding it uses the linked tipTable.c object.
size_t tip( uint8_t* dst, uint8_t const * src, size_t len ){
    rpInit();
    for( int i = 0; i < TipTableLength; i++ ){
        // get next needle (the next pattern)
        uint8_t * needle = TipTable[i].pt;
        size_t nlen = TipTable[i].sz;
        // Traverse the already found replacements to find hey stacks to look for the needle.
        int k = 0;
        do{      
            // get next hay stack
            uint8_t hay = src + rp[k].bo + rp[k].sz;
            size_t hlen = rp[k+1].bo - rp[k].bo - rp[k].sz;
            // search the needle
            uint8_t * loc = memmem( hay, hlen, needle, nlen );
            if( loc ){
                uint8_t by = TipTable[i].by;
                uint16_t offset = loc - src; // offset is the needle (=pattern) position.
                rpInsert( by, offset, nlen );
                k++; // The rp insert takes part inside the already processed rps.
            }
            k++;
        }while(hay+hlen>=src+len)
    }

    // Some bytes groups in the src buffer are replacable with 0x01...0xFF and some not.
    // The replacement list rp contains now the replacement information.
    // Lets collect the unreplacable bytes into a buffer now.

    collectUnreplacableBytes( src );
    convertUnreplacableBytes();
    return generateTipPacket( dst );
}

//! @brief rpInsert extends rp in an ordered way.
//! @param by The replacement byte for the location.
//! @param offset The location to be extended with.
//! @param sz The replacement pattern size.
void rpInsert( uint8_t by, uint16_t offset, uint8_t sz ){
    // int i = ri;
    // while( rp[i++].bo < bo );
    // rp[i].bo = bo;
    // rp[i].sz = sz;
    // rc++;
}

//! @brief ur contains all unreplacable bytes from src. It cannot get longer.
//! @details All unreplacable bytes are stretched inside the to 7-bit units.
//! This makes the data a bit longer.
static uint8_t ur[TIP_SRC_BUFFER_SIZE_MAX*8/7+1];

//! @brief urCount is the number of valid bytes inside ur.
static size_t urCount = 0;

static void collectUnreplacableBytes( uint8_t const * src ){
    // todo: Fill ur from rp data 
}

static void convertUnreplacableBytes( void ){
    // todo: Transform ur to 7-bit unit.
}

// generateTipPacket uses rp and ur to build the tip.
static size_t generateTipPacket( uint8_t dst ){
    // todo:
}


