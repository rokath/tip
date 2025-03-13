//! @file pack.h
//! @brief This is the tip user interface.
//! @author thomas.hoehenleitner [at] seerose.net

#ifndef PACK_H_
#define PACK_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>
#include <stddef.h>
#include "tip.h"

typedef uint8_t offset_t;

offset_t IDPosLimit(uint8_t i);

size_t tip( uint8_t* dst, const uint8_t * src, size_t len );
size_t tiPack( uint8_t * dst, const uint8_t * table, const uint8_t * src, size_t slen );
size_t buildTiPacket(uint8_t * dst, uint8_t * dstLimit, const uint8_t * table, const uint8_t * src, size_t slen);

//! IDPosition_t could get smaller, by storing only the offset, but this way the code is faster.
typedef struct{
    uint8_t id;     // id of pattern found in src
    offset_t start; // id pattern start
} IDPosition_t;

typedef struct {
    int count; //! count is the number of items inside IDPosTable.
    IDPosition_t item[TIP_SRC_BUFFER_SIZE_MAX-1];
} IDPosTable_t;

//! IDPosTable holds all IDs with their positions occuring in the current src buffer.
extern IDPosTable_t IDPosTable;

void newIDPosTable(const uint8_t * IDPatTable, const uint8_t * src, size_t slen);

#ifdef __cplusplus
}
#endif

#endif // PACK_H_
