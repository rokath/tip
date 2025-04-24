
#include "tipInternal.h"

uint8_t const * IDTable = idTable;

// tipInit is usable to change settings. Invalid values are ignored silently.
void tipInit( 
	unsigned urcb, // unreplacableContainerBits, only 6 and 7 are accepted
	unsigned id1Max, // ID1Max, can be up to 191 for unreplacableContainerBits == 6 or 127 for unreplacableContainerBits == 7.
	unsigned id1Count, // ID1Count, used ID1 values. The remaining space is for indirect indexing.
	uint8_t * const idt ){// idTable
 
	if( urcb == 6 || urcb == 7 ){
		unreplacableContainerBits = urcb;
	}
	if( unreplacableContainerBits == 6 && id1Max <= 191 
	  || unreplacableContainerBits == 7 && id1Max <= 127){
		ID1Max = id1Max;
	}
	ID1Count =  (ID1Max - ID1Count) * 255;
	MaxID = ID1Count + (ID1Max - ID1Count) * 255;
	if (idt) {
		IDTable = idt;
	}
	uint8_t * p = IDTable;
	maxPatternlength = 0;
	LastID = 0;
	while( *p ){
		LastID++;
		maxPatternlength = *p < maxPatternlength ? maxPatternlength : *p;
		p += *p + 1;
	}
}
