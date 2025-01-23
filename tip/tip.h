
#include <stdint.h>
#include "tipConfig.h"

typedef struct{
    uint8_t rp; // replacement byte
    uint16_t pt[2]; // 2-byte pattern
} tip2_t;

typedef struct{
    uint8_t rp; // replacement byte
    uint16_t pt[3]; // 3-byte pattern
} tip3_t;

typedef struct{
    uint8_t rp; // replacement byte
    uint8_t pt[4]; // 4-byte pattern
} tip4_t;

typedef struct{
    uint8_t rp; // replacement byte
    uint8_t pt[8]; // 8-byte pattern
} tip8_t;

extern tip2_t T2[];
extern tip3_t T3[];
extern tip4_t T4[];
extern tip8_t T8[];

extern const int T2cnt;
extern const int T3cnt;
extern const int T4cnt;
extern const int T8cnt;
