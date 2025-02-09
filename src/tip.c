//! @file tip.c
//! @brief Contains common tip code!
//! @author thomas.hoehenleitner [at] seerose.net

#include "tipInternal.h"

//! getPatternFromId seaches in testTable for id and returns its pattern location in pt.
//! @param pt is filled with the replacement pattern address if id was found.
//! @param sz is filled with the replacement size or 0, if id was not found.
//! @param id is the replacement byte. Valid values for id are 1...0x7f.
//! @param table is the pattern table.
void getPatternFromId( const uint8_t ** pt, size_t * sz, uint8_t id, const uint8_t * table ){
    unsigned int idx = 0;
    while( (*sz = *table++) && (*sz)){  // a pattern exists here
        if( ++idx == id ){ // id found
            *pt = table;
            return;
        }
        table += *sz;
    }
}

static const uint8_t * nextTablePos = 0;
static unsigned int nextID = 1;

//! initGetNextPattern causes getNextPattern to start from 0.
void initGetNextPattern( const uint8_t * table ){
    nextTablePos = table;
    nextID = 1;
}

//! getNextPattern returns next pattern location in pt and size in sz or *sz == 0.
//! @param pt is filled with the replacement pattern address if exists.
//! @param sz is filled with the replacement size or 0, if not exists.
void getNextPattern(const uint8_t ** pt, size_t * sz ){
    if( (*sz = *nextTablePos++) != 0 ){ // a pattern exists here
        *pt = nextTablePos;
        nextTablePos += *sz;
        return;
    }
}
