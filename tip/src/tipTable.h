
#include <stdint.h>

typedef struct {
    uint8_t by; //!< by is the replacement byte. 
    uint8_t sz; //!< size is the pt size.
    uint8_t * pt; //!< pt is the pointer to the replacement pattern.
} tip_t; 

extern tip_t TipTable[];

extern int const TipTableLength;
