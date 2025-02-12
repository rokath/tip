//! @file tip.h
//! @brief This is the tip configuration and tip internal common interface.
//! @author thomas.hoehenleitner [at] seerose.net

#ifndef TIP_H_
#define TIP_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>

//! TIP_SRC_BUFFER_SIZE_MAX is the maximun allowed input data buffer length.
//! Its size has influence on the statically allocated RAM.
#define TIP_SRC_BUFFER_SIZE_MAX 0x000000C8 // 200 bytes

#if TIP_SRC_BUFFER_SIZE_MAX > 0xffffffff
#error invalid TIP_SRC_BUFFER_SIZE_MAX value
#elif TIP_SRC_BUFFER_SIZE_MAX > 0xffffff
typedef uint32_t offset_t;
#elif TIP_SRC_BUFFER_SIZE_MAX > 0xffff
typedef uint32_t offset_t;
#elif TIP_SRC_BUFFER_SIZE_MAX > 0xff
typedef uint16_t offset_t;
#else
typedef uint8_t offset_t;
#endif

extern const uint8_t idTable[];

//! @brief replace_t is a replace type descriptor.
typedef struct {
    offset_t bo; // bo is the buffer offset, where replace bytes starts. It holds the list element count on index 0 instea
    uint8_t  sz; // sz is the replace size (2-255).
    uint8_t  id; // id is the replace byte 0x01 to 0x7f.
} replace_t;


#ifdef __cplusplus
}
#endif

#endif // TIP_H_
