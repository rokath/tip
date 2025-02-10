//! @file unpack.c
//! @brief This is the tip unpack code. Works also without pack.c.
//! @details todo
//! @author thomas.hoehenleitner [at] seerose.net

#include "tipInternal.h"



//! @brief tiu decodes src buffer with size len into dst buffer and returns decoded len.
size_t tiu( uint8_t * dst, const uint8_t * src, size_t slen ){
    if( slen < 6 ){
        memcpy(dst, src, slen);
        return slen;
    }

    // todo
    return 0;// shift78bit( dst, len, src-dst ); // dummy, todo
}
