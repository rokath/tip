//! @file tipInternal.h
//! @brief Contains not exported common declarations.
//! @author thomas.hoehenleitner [at] seerose.net

#ifndef TIP_INTERNAL_H_
#define TIP_INTERNAL_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <stddef.h>
#include <stdint.h>
#include <string.h>
#include "memmem.h"
#include "shift.h"
#include "tip.h"
#include "tipConfig.h"

#if TIP_SRC_BUFFER_SIZE_MAX > 0xffffffff

#error invalid TIP_SRC_BUFFER_SIZE_MAX value

#elif TIP_SRC_BUFFER_SIZE_MAX > 0xffffff

typedef uint32_t offset_t;

#elif TIP_SRC_BUFFER_SIZE_MAX > 0xffff

typedef uint32_t offset_t;
//typedef struct offset_tag {
//    unsigned offset_tag : 24;
//} __attribute__((packed))offset_t;

#elif TIP_SRC_BUFFER_SIZE_MAX > 0xff

typedef uint16_t offset_t;

#else

typedef uint8_t offset_t;

#endif

//! @brief replace_t is a replacement type descriptor.
typedef struct {
    offset_t bo; // bo is the buffer offset, where replacement bytes starts.
    uint8_t  sz; // sz is the replacement size (2-255).
    uint8_t  id; // id is the replacement byte 0x01 to 0x7f.
} replacement_t;

typedef struct {
    replacement_t  list[TIP_SRC_BUFFER_SIZE_MAX/2 + 2]; //!< list is the replacement list. It cannot get more elements. 
                                                        //! The space between 2 replacemts is a hay stack or finally unreplacable.
    int count; //!< count is the actual replace count inside replaceList.
} replaceList_t;

extern const uint8_t idTable[];

void getPatternFromId( const uint8_t ** pt, size_t * sz, uint8_t id, const uint8_t * table );
void initGetNextPattern( const uint8_t * table );
void getNextPattern(const uint8_t ** pt, size_t * sz );
replaceList_t * newReplacableList(size_t slen);
void replaceableListInsert( replaceList_t * r, int k, uint8_t by, offset_t offset, uint8_t sz );
size_t collectUnreplacableBytes( uint8_t * dst, replaceList_t * r, const uint8_t * src );
size_t generateTipPacket( uint8_t * dst, uint8_t * u7, size_t uSize, replaceList_t * r );
replaceList_t * buildReplacementList( const uint8_t * table, const uint8_t * src, size_t slen);
size_t TiPack( uint8_t* dst, const uint8_t * table, const uint8_t * src, size_t slen );

#ifdef __cplusplus
}
#endif

#endif // TIP_INTERNAL_H_
