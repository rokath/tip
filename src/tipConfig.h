//! @file tipConfig.h
//! @brief This is the tip configuration.
//! @author thomas.hoehenleitner [at] seerose.net

#ifndef TIP_CONFIG_H_
#define TIP_CONFIG_H_

#ifdef __cplusplus
extern "C" {
#endif

//! TIP_SRC_BUFFER_SIZE_MAX is the maximun allowed input data buffer length.
#define TIP_SRC_BUFFER_SIZE_MAX 1024

typedef uint16_t offset_t; //!< uint16_t allows to process up to 65 KB buffers.

#ifdef __cplusplus
}
#endif

#endif // TIP_CONFIG_H_
