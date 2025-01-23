
#include <strings.h>
#include <stddef.h>
#include "tip.h"


typedef struct {
    uint8_t rp;
    uint8_t * addr;
    size_t len;
} replacement_t;

//! rp8 is the 8-byte replacement list. It cannot get longer.
static replacement_t rp[TIP_SRC_BUFFER_SIZE_MAX/2];
static int rCnt = 0; //!< c8 is the 8-byte replacement count

//! rp8 is the 2-byte replacement list. It cannot get longer.
static replacement_t rp2[TIP_SRC_BUFFER_SIZE_MAX/2];
static int c2 = 0; //!< c8 is the 8-byte replacement count


/*
 * The memmem() function finds the start of the first occurrence of the
 * substring 'needle' of length 'nlen' in the memory area 'haystack' of
 * length 'hlen'.
 *
 * The return value is a pointer to the beginning of the sub-string, or
 * NULL if the substring is not found.
 */ // https://stackoverflow.com/questions/2188914/c-searching-for-a-string-in-a-file
void *memmem(const void *haystack, size_t hlen, const void *needle, size_t nlen)
{
    int needle_first;
    const void *p = haystack;
    size_t plen = hlen;

    if (!nlen) {
        return NULL;
    }
    needle_first = *(unsigned char *)needle;
    // https://en.cppreference.com/w/c/string/byte/memchr
    while (plen >= nlen && (p = memchr(p, needle_first, plen - nlen + 1)))
    {
        if (!memcmp(p, needle, nlen))
            return (void *)p;
        p++;
        plen = hlen - (p - haystack);
    }
    return NULL;
}




//! T8Scan searces in buf with size blen for T8 pattern.
//! If a match was found, rp8 is extended with c8 increment.
void T8Scan( uint8_t * buf, size_t blen ){
    for( int i = 0; i < T8cnt; i++ ){ // iterate over T8
        uint8_t* s = buf; // search location
    next:
        uint8_t * loc = memmem( s, blen - (s - buf), T8[i].pt, 8 );
        if( loc ){
            rp8[c8].addr = loc;
            rp8[c8++].rp = T8[i].rp;
            s = loc + 8; // next search start
            goto next;
        }
    }
}


//! T2Scan searces in buf with size blen for T2 pattern.
//! If a match was found, rp2 is extended with c2 increment.
void T2Scan( uint8_t * buf, size_t blen ){
    for( int i = 0; i < T2cnt; i++ ){ // iterate over T8
        uint8_t* s = buf; // search location
    next:
        uint8_t * loc = memmem( s, blen - (s - buf), T2[i].pt, 2 );
        if( loc ){
            rp2[c2].addr = loc;
            rp2[c2++].rp = T2[i].rp;
            s = loc + 2; // next search start
            goto next;
        }
    }
}



// TipPack tip encodes src buffer with size len into dst buffer and returns encoded len.
// For the tip encoding it uses the linked tipTable.c object.
size_t TipPack( void* dst, void* src, size_t len ){
    //memset(flags, 0, sizeof(flags));
    c8 = c2 = 0;
    uint8_t* p = dst;
    T8Scan( src, len );
    for( int i = 0; i < c8; i++ ){
        T2Scan( rp8[i].addr, )
    }
    return 0;
}
