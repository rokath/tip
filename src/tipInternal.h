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

//! @brief replacement_t is a replacement type descriptor.
typedef struct {
    offset_t bo; // bo is the buffer offset, where replacement size starts.
    uint8_t  sz; // sz is the replacement size (2-255).
    uint8_t  by; // by is the replacement byte 0x01 to 0xff.
} replacement_t;

extern uint8_t tipTable[];
extern const size_t tipTableSize;

void getPatternFromId( uint8_t id, uint8_t ** pt, size_t * sz );

#ifdef __cplusplus
}
#endif

#endif // TIP_INTERNAL_H_
