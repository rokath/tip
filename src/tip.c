//! @file tip.c
//! @brief Contains common tip code!
//! @author thomas.hoehenleitner [at] seerose.net

#include "tipInternal.h"

//! getPatternFromId seaches in testTable for id and returns its pattern location in pt.
//! @param id is the replacement byte. Valid values for id are 1...0x7f.
//! @param pt is filled with the replacement pattern address if id was found.
//! @param sz is filled with the replacement size or 0, if id was not found.
void getPatternFromId( uint8_t id, uint8_t ** pt, size_t * sz ){
    unsigned int idx = 0;
    for( size_t i = 0; i < tipTableSize; ){
        *sz = tipTable[i++];
        if( *sz ){ // a pattern exists here
            idx++;
            if( idx == id ){ // id found
                *pt = &tipTable[i];
                return;
            }
            i += *sz;
        }
    }
}

static uint8_t nextTipTablePos = 0;
static unsigned int nextID = 1;

//! resetPattern causes getNextPattern to start from 0.
void resetPattern(void){
    nextTipTablePos = 0;
    nextID = 1;
}

//! getNextPattern returns next pattern location in pt and size in sz or *sz == 0.
//! @param pt is filled with the replacement pattern address if exists.
//! @param sz is filled with the replacement size or 0, if not exists.
void getNextPattern(uint8_t ** pt, size_t * sz ){
    *sz = tipTable[nextTipTablePos++];
    if( *sz ){ // a pattern exists here
        *pt = &tipTable[nextTipTablePos];
        nextTipTablePos += *sz;
        return;
    }
}
