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
#define TIP_SRC_BUFFER_SIZE_MAX 31*8u // 248 bytes (max 65528u)

//! TIP_MAX_PATH_COUNT is the max allowed path count.
//! Its size has influence on the statically allocated RAM.
//! In a first run on each src buffer position a 2-byte pattern could match.
//! In a first run on the first src buffer position all IDPosition table entries coud match.
//! TODO: Find right logic.
#define TIP_MAX_PATH_COUNT (2*(TIP_SRC_BUFFER_SIZE_MAX-1))

#ifdef __cplusplus
}
#endif

#endif // TIP_CONFIG_H_
