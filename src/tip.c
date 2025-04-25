//! @file tip.c
//! @brief This is optional wrapper code mainly for testing.

#include "tipInternal.h"

// tipInit is usable to change settings. Invalid values are ignored silently.
static void tipInit( 
	unsigned urcb,         //!< @param unreplacableContainerBits, only 6 and 7 are accepted
	unsigned id1Count,     //!< @param ID1Count, used ID1 values of ID1Max. The remaining space is for indirect indexing.
	uint8_t const * idt ){ //!< @param idTable
 
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
	uint8_t const * p = IDTable;
	maxPatternlength = 0;
	LastID = 0;
	while( *p ){
		LastID++;
		maxPatternlength = *p < maxPatternlength ? maxPatternlength : *p;
		p += *p + 1;
	}
}

static unsigned urcb;        //!< stored unreplacableContainerBits, only 6 and 7 are accepted
static unsigned id1Count;    //!< stored ID1Count, used ID1 values. The remaining space is for indirect indexing.
static uint8_t const * idt;  //!< stored idTable

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

//! @brief tiPack encodes src buffer with size len into dst buffer and returns encoded len.
//! @details For the tip encoding it uses the passed ubc, id1Count and ID table.
//! @param dst is the destination buffer.
//! @param src is the source buffer.
//! @param slen is the valid data length inside source buffer.
//! @param ubc is the unreplacable container bits count, only 6 and 7 are accepted.
//! @param id1Count is the used ID1 values from ID1Max. The remaining space is for indirect indexing.
//! @param idt is the used id table.
//! @retval length of tip packet
size_t tiPack( uint8_t * dst, uint8_t const * src,size_t slen, unsigned ubc, unsigned id1Count, uint8_t const * idt ){
	tipStoreGlobals();
    tipInit(ubc, id1Count, idt);
    size_t plen = tip( dst, src, slen );
	tipRestoreGlobals();
	return plen;
}

//! @brief tiUnpack decodes src buffer with size len into dst buffer and returns encoded len.
//! @details For the tip decoding it uses the passed ubc, id1Count and ID table.
//! @param dst is the destination buffer.
//! @param src is the source buffer.
//! @param slen is the valid data length inside source buffer.
//! @param ubc is the unreplacable container bits count, only 6 and 7 are accepted.
//! @param id1Count is the used ID1 values from ID1Max. The remaining space is for indirect indexing.
//! @param idt is the used id table.
//! @retval length of untip packet
size_t tiUnpack( uint8_t* dst, const uint8_t * src, size_t slen, unsigned ubc, unsigned id1Count, const uint8_t * idt ){
	tipStoreGlobals();
    tipInit(ubc, id1Count, idt);
    size_t plen = tiu( dst, src, slen );
	tipRestoreGlobals();
	return plen;
}
