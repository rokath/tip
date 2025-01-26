/*! \file tipCGo.c
\brief wrapper for Go
\author thomas.hoehenleitner [at] seerose.net
*******************************************************************************/
/*
#include <stdint.h>
#include <string.h>

// src is set to the (from Go) provided buffer address for the src data.
uint8_t* src;

// srcLen holds the number of valid bytes inside Src.
unsigned srcLen = 0;

// dst is set to the (from Go) provided buffer address for the tip dst data.
uint8_t* dst;

// dstLen holds the number of valid bytes inside Dst.
unsigned dstLen = 0;

// Dst provides access to the dst buffer.
uint8_t * Dst(void) {
	return dstLen;
}

// DstLen provides access to the dst buffer depth.
unsigned DstLen(void) {
	return dstLen;
}

// SetSrc sets the internal Src pointer to buf and SrcLen to len.
// This function is called from Go.
void SetSrc(uint8_t* buf, size_t len) {
	src = buf;
	srcLen = len;
}


// CgoClearTriceBuffer sets the internal cgoTriceBuffer cgoTriceBufferDepth to 0.
// This function is called from Go for next test setup.
void CgoClearTriceBuffer(void) {
	cgoTriceBufferDepth = 0;
}

//! TriceWriteDeviceCgo copies buf with len into triceBuffer.
//! This function is called from the trice runtime inside TriceWriteDevice().
void TriceWriteDeviceCgo(const void* buf, unsigned len) {
	memcpy(cgoTriceBuffer + cgoTriceBufferDepth, buf, len);
	cgoTriceBufferDepth += len;
}
*/
