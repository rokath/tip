//! @file idTable.c
//! @brief Generated code - do not edit!

#include <stdint.h>
#include <stddef.h>

//! idTable is sorted by pattern length and pattern count.
//! The pattern position + 1 is the replace id.
//! The generator pattern max size was 4 and the list pattern max size is: 4
const uint8_t idTable[] = { // from ./try.txt
                                 // `ASCII`|  count  id
	  4, 0x31, 0x32, 0x33, 0x31, // `1231`|      2  01
	  4, 0x32, 0x33, 0x31, 0x32, // `2312`|      2  02
	  4, 0x33, 0x31, 0x32, 0x33, // `3123`|      2  03
	  3, 0x31, 0x32, 0x33,       // `123` |      3  04
	  3, 0x32, 0x33, 0x31,       // `231` |      2  05
	  3, 0x33, 0x31, 0x32,       // `312` |      2  06
	  2, 0x31, 0x32,             // `12`  |      3  07
	  2, 0x32, 0x33,             // `23`  |      3  08
	  2, 0x33, 0x31,             // `31`  |      2  09
	  0 // table end marker
};

// tipTableSize is 37.

//   0: (   2)   4, 0x31, 0x32, 0x33, 0x31, // `1231`
//   1: (   2)   4, 0x32, 0x33, 0x31, 0x32, // `2312`
//   2: (   2)   4, 0x33, 0x31, 0x32, 0x33, // `3123`
//   3: (   3)   3, 0x31, 0x32, 0x33,       // `123` 
//   4: (   2)   3, 0x32, 0x33, 0x31,       // `231` 
//   5: (   2)   3, 0x33, 0x31, 0x32,       // `312` 
//   6: (   3)   2, 0x31, 0x32,             // `12`  
//   7: (   3)   2, 0x32, 0x33,             // `23`  
//   8: (   2)   2, 0x33, 0x31,             // `31`  

