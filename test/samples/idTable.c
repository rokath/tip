//! @file idTable.c
//! @brief Generated code - do not edit!

#include <stdint.h>
#include <stddef.h>

//! idTable is sorted by pattern length and pattern count.
//! The pattern position + 1 is the replace id.
//! The generator pattern max size was 4 and the list pattern max size is: 4
const uint8_t idTable[] = { // from ./try.txt
                                 // `ASCII`|  count  id
	  4, 0x31, 0x32, 0x33, 0x31, // `1231`|      1  01
	  4, 0x32, 0x33, 0x31, 0x32, // `2312`|      1  02
	  4, 0x33, 0x31, 0x32, 0x33, // `3123`|      1  03
	  3, 0x31, 0x32, 0x33,       // `123` |      2  04
	  3, 0x32, 0x33, 0x31,       // `231` |      1  05
	  3, 0x33, 0x31, 0x32,       // `312` |      1  06
	  2, 0x31, 0x32,             // `12`  |      2  07
	  2, 0x32, 0x33,             // `23`  |      2  08
	  2, 0x33, 0x31,             // `31`  |      1  09
	  0 // table end marker
};

// tipTableSize is 37.

//   3: (   2)   3, 0x31, 0x32, 0x33,       // `123` 
//   6: (   2)   2, 0x31, 0x32,             // `12`  
//   7: (   2)   2, 0x32, 0x33,             // `23`  

