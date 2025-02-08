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
    replacement_t  list[TIP_SRC_BUFFER_SIZE_MAX/2 + 2]; //!< list is the replacement list. It cannot get more elements. The space between 2 replacemts is a hay stack or finally unreplacable.
    int count; //!< count is the actual replace count inside replaceList.
} replace_t;

//! @details All unreplacable bytes are stretched inside to 7-bit units. This makes the data a bit longer.
typedef struct {
    uint8_t buffer[TIP_SRC_BUFFER_SIZE_MAX*8/7+1];  //!< buffer holds all unreplacable bytes from src. It cannot get longer.
    uint8_t * last; //! last is the last address inside the unreplacable bytes buffer. ( = &(buffer[sizeof(buffer)-1]); )
} unreplacable_t;

typedef struct {
    uint8_t const length
    uint8_t const * const sequence;
} pattern_t;

typedef struct {
    pattern_t * pattern;
    size_t const size;
} idTable_t;



void getPatternFromId( uint8_t id, uint8_t ** pt, size_t * sz );
void restartPattern(void);
void getNextPattern(uint8_t ** pt, size_t * sz );

#ifdef __cplusplus
}
#endif

#endif // TIP_INTERNAL_H_
