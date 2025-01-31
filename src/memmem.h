//! @file memmem.h
//! @brief This is the memmem interface.

#include <stddef.h>

void *memmem(const void *haystack, size_t hlen, const void *needle, size_t nlen);
int MemmemOffset(const void *haystack, size_t hlen, const void *needle, size_t nlen);
