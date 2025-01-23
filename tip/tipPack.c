
#include <strings.h>
#include <stddef.h>
#include "tip.h"


typedef struct {
    int first;
    int limit;
    uint8_t rp;
} replacement_t;

//! rpl is the replacement list. It cannot get longer.
static replacement_t rpl[TIP_SRC_BUFFER_SIZE_MAX/2];
static int rC = 0; //!< rC is the replacement count

//! T8Scan searces in p with size len for T8 pattern.
//! If a match was found, rpl is extended with rC increment.
void T8Scan( uint8_t * p, size_t len ){
    uint8_t* s = p; // search index
    for( int i = 0; i < T8cnt; i++ ){ // iterate over T8
        for( int n = 0; n < 8; n++) { // iterate over pt
            uint8_t b =  T8[i].pt[n]; // next pt byte
            while( s + 8 < p + len ){ // iterate over s
                uint8_t x = *s;
                if( x != b ){ // first pt byte not this s byte
                    s++;
                    continue; // compare next s byte with first pt byte
                }
                // Here the first pt byte matches the actual s byte.

            }
        }   
   }

                    rpl[rC].first = k;
                    k += 8;
                    rpl[rC].limit = k;
                    rpl[rC++].rp = T8[i].rp;


}


//! byteIndex returns first b offset in p or len.
//! retval 0: first byte in p is equal b or len was 0.
//! retval len: B not found in [p,p+len)
int byteOffset( uint8_t * p, size_t len, uint8_t b ){
    for( int i = 0; i < len; i++ ){
        if( b == *p++ ){
            return i;
        }
    }
    return len; // -1?
}



// TipPack tip encodes src buffer with size len into dst buffer and returns encoded len.
// For the tip encoding it uses the linked tipTable.c object.
size_t TipPack( void* dst, void* src, size_t len ){
    //memset(flags, 0, sizeof(flags));
    rC = 0;
    uint8_t* p = dst;
    T8Scan( src, len );
    return 0;
}
