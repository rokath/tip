//! @file tipTable.c
//! @brief Generated code - do not edit!

#include <stddef.h>
#include <stdint.h>

#include "tipInternal.h"

// tipTable is sorted by pattern length and pattern count.
// The pattern position + 1 is the replacement id.
uint8_t tipTable[] = {
  // sz, pattern                                           id  count
     8, 0xff, 0xff, 0xee, 0xcc, 0xfa, 0xaf, 0x00, 0xaa, // 01  12345
     8, 0xff, 0xff, 0xee, 0xcc, 0xfa, 0xaf, 0x00, 0xbb, // 02    123
     5, 0xff, 0xff, 0xee, 0xcc, 0xfa,                   // 03   1234
     3, 0xff, 0xff, 0xee,                               // 04   9012
     0 // table end marker
};

const size_t tipTableSize = 25;
