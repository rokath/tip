//! @file tipTable.c
//! @brief Generated code - do not edit!

#include "tipTable.h"

char* s = "hi";
uint8_t f[2] = { 1, 2 };


tip_t t = { 
    .sz=8, 
    .pt={0xff, 0xff, 0xee, 0xcc, 0xfa, 0xaf, 0x00, 0xaa} 
};

tip_t U8 = { 0x52, 8, {0xff, 0xff, 0xee, 0xcc, 0xfa, 0xaf, 0x00, 0xaa} };
tip_t U3 = { 0x58, 3, {0xaf, 0x00, 0xaa} };

// tipTable is sorted by pattern length and pattern count.
// The pattern position + 1 is the replacement id
static uint8_t tipTable[] = {
     8, 0xff, 0xff, 0xee, 0xcc, 0xfa, 0xaf, 0x00, 0xaa, // cnt =12345
     8, 0xff, 0xff, 0xee, 0xcc, 0xfa, 0xaf, 0x00, 0xbb, // cnt =  123
     3, 0xff, 0xff, 0xee,
     0 // table end marker
};


#include <stddef.h>



//! getPatternFromId returns a pointer or NULL. 
static int getPatternFromId( uint8_t id, uint8_t ** pt, size_t * sz ){
    uint8_t sz = tipTable[0];
    uint8_t p = &tipTable[1];
    unsigned int idx = 0;
    for( int i = 0; i < sizeof(tipTable); ){
        sz = tipTable[i++];
        if( sz ){
            idx++;
            if( idx == id ){ // found
                
            i += sz;

        }
        p = &tipTable[i];
    }
    r
    for(int i = 0; i < TipTableLength){
        if( i == id ){
            *pt = TipTable[i].pt;
            *sz = TipTable[i].sz;
            return 1;
        }
    }
    return 0;
}