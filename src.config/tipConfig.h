//! @file tipConfig.h
//! @brief This is the tip configuration.

#ifndef TIP_CONFIG_H_
#define TIP_CONFIG_H_

#ifdef __cplusplus
extern "C" {
#endif

//! TIP_SRC_BUFFER_SIZE_MAX is the maximun allowed input data buffer length.
//! Its size has influence on the statically allocated RAM.
#define TIP_SRC_BUFFER_SIZE_MAX 65528u // bytes (max 65528u)

//! TIP_MAX_PATH_COUNT is the max allowed path count.
//! Its size has influence on the statically allocated RAM.
#define TIP_MAX_PATH_COUNT 20000

#ifdef __cplusplus
}
#endif

#endif // TIP_CONFIG_H_
