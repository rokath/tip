//! @file mmemmem.c
//! @brief From // https://stackoverflow.com/questions/2188914/c-searching-for-a-string-in-a-file

#include <string.h>
#include "memmem.h"

/** 
 * @brief The memmem() function finds the start of the first occurrence of the
 * substring 'needle' of length 'nlen' in the memory area 'haystack' of
 * length 'hlen'.
 *
 * @details The return value is a pointer to the beginning of the sub-string, or
 * NULL if the substring is not found.
 */
void *Memmem(const void *haystack, size_t hlen, const void *needle, size_t nlen)
{
    int needle_first;
    const void *p = haystack;
    size_t plen = hlen;

    if (!nlen) {
        return NULL;
    }
    needle_first = *(unsigned char *)needle;

    while (plen >= nlen && (p = memchr(p, needle_first, plen - nlen + 1)))
    {
        if (!memcmp(p, needle, nlen))
            return (void *)p;
        p++;
        plen = hlen - (p - haystack);
    }
    return NULL;
}
