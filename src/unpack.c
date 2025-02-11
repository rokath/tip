//! @file unpack.c
//! @brief This is the tip unpack code. Works also without pack.c.
//! @details todo
//! @author thomas.hoehenleitner [at] seerose.net

#include "tipInternal.h"

//! @brief tiu decodes src buffer with size len into dst buffer and returns decoded len.
size_t tiu( uint8_t * dst, const uint8_t * src, size_t slen ){
    return tiUnpack(dst, idTable, src, slen );
}

//! @brief tiUnpack decodes src buffer with size slen into dst buffer and returns decoded dlen.
//! @details For the tip decoding it uses the passed idTable object.
size_t tiUnpack( uint8_t* dst, const uint8_t * table, const uint8_t * src, size_t slen ){
    if( slen < 6 ){
        memcpy(dst, src, slen);
        return slen;
    }   
    return 0; // todo
}
