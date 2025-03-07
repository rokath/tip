//! @file tipConfig.h
//! @brief This is the tip configuration.

#ifndef TIP_CONFIG_H_
#define TIP_CONFIG_H_

#ifdef __cplusplus
extern "C" {
#endif

//! TIP_SRC_BUFFER_SIZE_MAX is the maximun allowed input data buffer length.
//! Its size has influence on the statically allocated RAM.
//! Must be a multiple of 8.
#define TIP_SRC_BUFFER_SIZE_MAX 31*8u // 248 bytes (max65528)

#ifdef __cplusplus
}
#endif

#endif // TIP_CONFIG_H_
