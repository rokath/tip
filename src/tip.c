//! @file tip.c
//! @brief This is optional wrapper code mainly for testing.

#include "tipInternal.h"

// tipInit is usable to change settings. Invalid values are ignored silently.
void tipInit( 
	unsigned urcb,         //!< @param unreplacableContainerBits, only 6 and 7 are accepted
	unsigned id1Count,     //!< @param ID1Count, used ID1 values of ID1Max. The remaining space is for indirect indexing.
	uint8_t * const idt ){ //!< @param idTable
 
	if( urcb == 6 || urcb == 7 ){
		unreplacableContainerBits = urcb;
	}
	if( unreplacableContainerBits == 6 ){
		ID1Max = 191; 
	} else if (unreplacableContainerBits == 7){
		ID1Max = 127;
	} 
	if (id1Count <= ID1Max){
		ID1Count = id1Count;
	}else if (ID1Count > ID1Max) {
		ID1Count = ID1Max;
	}
	unsigned ID2Count =  (ID1Max - ID1Count) * 255;
	MaxID = ID1Count + ID2Count;
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

static unsigned urcb;        //!< @param unreplacableContainerBits, only 6 and 7 are accepted
static unsigned id1Count;    //!< @param ID1Count, used ID1 values. The remaining space is for indirect indexing.
static uint8_t const * idt;  //!< @param idTable

//! @brief tipStoreGlobals reads global varibles (to restore them later).
static void tipStoreGlobals( void ) {
	urcb = unreplacableContainerBits;
	id1Count = ID1Count;    
	idt = IDTable;
}

//! @brief tipRestoreGlobals restores global varibles (after storing them).
static void tipRestoreGlobals( void ) {
	unreplacableContainerBits = urcb;
	ID1Count = id1Count;    
	IDTable = idt;
}

//! tiPack2 stores and initializes global variables, performs the tip function and restores global variables.
//! @retval length of tip packet
size_t tiPack2( 
    uint8_t * dst,               //!< @param destination buffer
    uint8_t const * src,         //!< @param source buffer
    size_t slen,                 //!< @param valid data length inside source buffer
    unsigned urcb,               //!< @param unreplacableContainerBits, only 6 and 7 are accepted
	unsigned id1Count,           //!< @param ID1Count, used ID1 values. The remaining space is for indirect indexing.
	uint8_t const * const idt ){ //!< @param idTable

	tipStoreGlobals();
    tipInit(urcb, id1Count, idt);
    size_t plen = tip( dst, src, slen );
	tipRestoreGlobals();
	return plen;
}
