//! @file idTable.c
//! @brief Generated code - do not edit!

#include <stdint.h>
#include <stddef.h>

//! idTable is sorted by pattern length and pattern count.
//! The pattern position + 1 is the replace id.
//! The generator pattern max size was 10 and the list pattern max size is: 10
const uint8_t idTable[] = { // from ./try.txt
                                                                     // `ASCII     `|  count  id
	 10, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, // `0123456789`|      3  01
	 10, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, // `1234567890`|      3  02
	 10, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, // `2345678901`|      2  03
	 10, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, // `3456789012`|      2  04
	 10, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, // `4567890123`|      2  05
	 10, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, // `5678901234`|      2  06
	 10, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, // `6789012345`|      2  07
	 10, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, // `7890123456`|      2  08
	 10, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, // `8901234567`|      2  09
	 10, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, // `9012345678`|      2  0a
	  9, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38,       // `012345678` |      3  0b
	  9, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39,       // `123456789` |      3  0c
	  9, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30,       // `234567890` |      3  0d
	  9, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31,       // `345678901` |      2  0e
	  9, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32,       // `456789012` |      2  0f
	  9, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33,       // `567890123` |      2  10
	  9, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34,       // `678901234` |      2  11
	  9, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35,       // `789012345` |      2  12
	  9, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36,       // `890123456` |      2  13
	  9, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37,       // `901234567` |      2  14
	  8, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37,             // `01234567`  |      3  15
	  8, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38,             // `12345678`  |      3  16
	  8, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39,             // `23456789`  |      3  17
	  8, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30,             // `34567890`  |      3  18
	  8, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31,             // `45678901`  |      2  19
	  8, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32,             // `56789012`  |      2  1a
	  8, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33,             // `67890123`  |      2  1b
	  8, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34,             // `78901234`  |      2  1c
	  8, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35,             // `89012345`  |      2  1d
	  8, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36,             // `90123456`  |      2  1e
	  7, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36,                   // `0123456`   |      3  1f
	  7, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37,                   // `1234567`   |      3  20
	  7, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38,                   // `2345678`   |      3  21
	  7, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39,                   // `3456789`   |      3  22
	  7, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30,                   // `4567890`   |      3  23
	  7, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31,                   // `5678901`   |      2  24
	  7, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32,                   // `6789012`   |      2  25
	  7, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33,                   // `7890123`   |      2  26
	  7, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34,                   // `8901234`   |      2  27
	  7, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35,                   // `9012345`   |      2  28
	  6, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35,                         // `012345`    |      3  29
	  6, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36,                         // `123456`    |      3  2a
	  6, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37,                         // `234567`    |      3  2b
	  6, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38,                         // `345678`    |      3  2c
	  6, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39,                         // `456789`    |      3  2d
	  6, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30,                         // `567890`    |      3  2e
	  6, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31,                         // `678901`    |      2  2f
	  6, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32,                         // `789012`    |      2  30
	  6, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33,                         // `890123`    |      2  31
	  6, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34,                         // `901234`    |      2  32
	  5, 0x30, 0x31, 0x32, 0x33, 0x34,                               // `01234`     |      3  33
	  5, 0x31, 0x32, 0x33, 0x34, 0x35,                               // `12345`     |      3  34
	  5, 0x32, 0x33, 0x34, 0x35, 0x36,                               // `23456`     |      3  35
	  5, 0x33, 0x34, 0x35, 0x36, 0x37,                               // `34567`     |      3  36
	  5, 0x34, 0x35, 0x36, 0x37, 0x38,                               // `45678`     |      3  37
	  5, 0x35, 0x36, 0x37, 0x38, 0x39,                               // `56789`     |      3  38
	  5, 0x36, 0x37, 0x38, 0x39, 0x30,                               // `67890`     |      3  39
	  5, 0x37, 0x38, 0x39, 0x30, 0x31,                               // `78901`     |      2  3a
	  5, 0x38, 0x39, 0x30, 0x31, 0x32,                               // `89012`     |      2  3b
	  5, 0x39, 0x30, 0x31, 0x32, 0x33,                               // `90123`     |      2  3c
	  4, 0x30, 0x31, 0x32, 0x33,                                     // `0123`      |      3  3d
	  4, 0x31, 0x32, 0x33, 0x34,                                     // `1234`      |      3  3e
	  4, 0x32, 0x33, 0x34, 0x35,                                     // `2345`      |      3  3f
	  4, 0x33, 0x34, 0x35, 0x36,                                     // `3456`      |      3  40
	  4, 0x34, 0x35, 0x36, 0x37,                                     // `4567`      |      3  41
	  4, 0x35, 0x36, 0x37, 0x38,                                     // `5678`      |      3  42
	  4, 0x36, 0x37, 0x38, 0x39,                                     // `6789`      |      3  43
	  4, 0x37, 0x38, 0x39, 0x30,                                     // `7890`      |      3  44
	  4, 0x38, 0x39, 0x30, 0x31,                                     // `8901`      |      2  45
	  4, 0x39, 0x30, 0x31, 0x32,                                     // `9012`      |      2  46
	  3, 0x30, 0x31, 0x32,                                           // `012`       |      3  47
	  3, 0x31, 0x32, 0x33,                                           // `123`       |      3  48
	  3, 0x32, 0x33, 0x34,                                           // `234`       |      3  49
	  3, 0x33, 0x34, 0x35,                                           // `345`       |      3  4a
	  3, 0x34, 0x35, 0x36,                                           // `456`       |      3  4b
	  3, 0x35, 0x36, 0x37,                                           // `567`       |      3  4c
	  3, 0x36, 0x37, 0x38,                                           // `678`       |      3  4d
	  3, 0x37, 0x38, 0x39,                                           // `789`       |      3  4e
	  3, 0x38, 0x39, 0x30,                                           // `890`       |      3  4f
	  3, 0x39, 0x30, 0x31,                                           // `901`       |      2  50
	  2, 0x30, 0x31,                                                 // `01`        |      3  51
	  2, 0x31, 0x32,                                                 // `12`        |      3  52
	  2, 0x32, 0x33,                                                 // `23`        |      3  53
	  2, 0x33, 0x34,                                                 // `34`        |      3  54
	  2, 0x34, 0x35,                                                 // `45`        |      3  55
	  2, 0x35, 0x36,                                                 // `56`        |      3  56
	  2, 0x36, 0x37,                                                 // `67`        |      3  57
	  2, 0x37, 0x38,                                                 // `78`        |      3  58
	  2, 0x38, 0x39,                                                 // `89`        |      3  59
	  2, 0x39, 0x30,                                                 // `90`        |      3  5a
	  0 // table end marker
};

// tipTableSize is 631.

//   0: (   3)  10, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, // `0123456789`
//   1: (   3)  10, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, // `1234567890`
//   2: (   2)  10, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, // `2345678901`
//   3: (   2)  10, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, // `3456789012`
//   4: (   2)  10, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, // `4567890123`
//   5: (   2)  10, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, // `5678901234`
//   6: (   2)  10, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, // `6789012345`
//   7: (   2)  10, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, // `7890123456`
//   8: (   2)  10, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, // `8901234567`
//   9: (   2)  10, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, // `9012345678`
//  10: (   3)   9, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38,       // `012345678` 
//  11: (   3)   9, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39,       // `123456789` 
//  12: (   3)   9, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30,       // `234567890` 
//  13: (   2)   9, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31,       // `345678901` 
//  14: (   2)   9, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32,       // `456789012` 
//  15: (   2)   9, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33,       // `567890123` 
//  16: (   2)   9, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34,       // `678901234` 
//  17: (   2)   9, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35,       // `789012345` 
//  18: (   2)   9, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36,       // `890123456` 
//  19: (   2)   9, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37,       // `901234567` 
//  20: (   3)   8, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37,             // `01234567`  
//  21: (   3)   8, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38,             // `12345678`  
//  22: (   3)   8, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39,             // `23456789`  
//  23: (   3)   8, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30,             // `34567890`  
//  24: (   2)   8, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31,             // `45678901`  
//  25: (   2)   8, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32,             // `56789012`  
//  26: (   2)   8, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33,             // `67890123`  
//  27: (   2)   8, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34,             // `78901234`  
//  28: (   2)   8, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35,             // `89012345`  
//  29: (   2)   8, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36,             // `90123456`  
//  30: (   3)   7, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36,                   // `0123456`   
//  31: (   3)   7, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37,                   // `1234567`   
//  32: (   3)   7, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38,                   // `2345678`   
//  33: (   3)   7, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39,                   // `3456789`   
//  34: (   3)   7, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30,                   // `4567890`   
//  35: (   2)   7, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31,                   // `5678901`   
//  36: (   2)   7, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32,                   // `6789012`   
//  37: (   2)   7, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33,                   // `7890123`   
//  38: (   2)   7, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34,                   // `8901234`   
//  39: (   2)   7, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35,                   // `9012345`   
//  40: (   3)   6, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35,                         // `012345`    
//  41: (   3)   6, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36,                         // `123456`    
//  42: (   3)   6, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37,                         // `234567`    
//  43: (   3)   6, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38,                         // `345678`    
//  44: (   3)   6, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39,                         // `456789`    
//  45: (   3)   6, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30,                         // `567890`    
//  46: (   2)   6, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31,                         // `678901`    
//  47: (   2)   6, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32,                         // `789012`    
//  48: (   2)   6, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33,                         // `890123`    
//  49: (   2)   6, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34,                         // `901234`    
//  50: (   3)   5, 0x30, 0x31, 0x32, 0x33, 0x34,                               // `01234`     
//  51: (   3)   5, 0x31, 0x32, 0x33, 0x34, 0x35,                               // `12345`     
//  52: (   3)   5, 0x32, 0x33, 0x34, 0x35, 0x36,                               // `23456`     
//  53: (   3)   5, 0x33, 0x34, 0x35, 0x36, 0x37,                               // `34567`     
//  54: (   3)   5, 0x34, 0x35, 0x36, 0x37, 0x38,                               // `45678`     
//  55: (   3)   5, 0x35, 0x36, 0x37, 0x38, 0x39,                               // `56789`     
//  56: (   3)   5, 0x36, 0x37, 0x38, 0x39, 0x30,                               // `67890`     
//  57: (   2)   5, 0x37, 0x38, 0x39, 0x30, 0x31,                               // `78901`     
//  58: (   2)   5, 0x38, 0x39, 0x30, 0x31, 0x32,                               // `89012`     
//  59: (   2)   5, 0x39, 0x30, 0x31, 0x32, 0x33,                               // `90123`     
//  60: (   3)   4, 0x30, 0x31, 0x32, 0x33,                                     // `0123`      
//  61: (   3)   4, 0x31, 0x32, 0x33, 0x34,                                     // `1234`      
//  62: (   3)   4, 0x32, 0x33, 0x34, 0x35,                                     // `2345`      
//  63: (   3)   4, 0x33, 0x34, 0x35, 0x36,                                     // `3456`      
//  64: (   3)   4, 0x34, 0x35, 0x36, 0x37,                                     // `4567`      
//  65: (   3)   4, 0x35, 0x36, 0x37, 0x38,                                     // `5678`      
//  66: (   3)   4, 0x36, 0x37, 0x38, 0x39,                                     // `6789`      
//  67: (   3)   4, 0x37, 0x38, 0x39, 0x30,                                     // `7890`      
//  68: (   2)   4, 0x38, 0x39, 0x30, 0x31,                                     // `8901`      
//  69: (   2)   4, 0x39, 0x30, 0x31, 0x32,                                     // `9012`      
//  70: (   3)   3, 0x30, 0x31, 0x32,                                           // `012`       
//  71: (   3)   3, 0x31, 0x32, 0x33,                                           // `123`       
//  72: (   3)   3, 0x32, 0x33, 0x34,                                           // `234`       
//  73: (   3)   3, 0x33, 0x34, 0x35,                                           // `345`       
//  74: (   3)   3, 0x34, 0x35, 0x36,                                           // `456`       
//  75: (   3)   3, 0x35, 0x36, 0x37,                                           // `567`       
//  76: (   3)   3, 0x36, 0x37, 0x38,                                           // `678`       
//  77: (   3)   3, 0x37, 0x38, 0x39,                                           // `789`       
//  78: (   3)   3, 0x38, 0x39, 0x30,                                           // `890`       
//  79: (   2)   3, 0x39, 0x30, 0x31,                                           // `901`       
//  80: (   3)   2, 0x30, 0x31,                                                 // `01`        
//  81: (   3)   2, 0x31, 0x32,                                                 // `12`        
//  82: (   3)   2, 0x32, 0x33,                                                 // `23`        
//  83: (   3)   2, 0x33, 0x34,                                                 // `34`        
//  84: (   3)   2, 0x34, 0x35,                                                 // `45`        
//  85: (   3)   2, 0x35, 0x36,                                                 // `56`        
//  86: (   3)   2, 0x36, 0x37,                                                 // `67`        
//  87: (   3)   2, 0x37, 0x38,                                                 // `78`        
//  88: (   3)   2, 0x38, 0x39,                                                 // `89`        
//  89: (   3)   2, 0x39, 0x30,                                                 // `90`        

