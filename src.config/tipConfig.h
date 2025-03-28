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

#define VERBOSE 0
#define DEBUG 0

//! OPTIMIZE_UNREPLACABLES allows to reduce the TiP packet size in some special cases.
//! It is a selectable option just for tests and should be enabled always.
#define OPTIMIZE_UNREPLACABLES 1

//! UNREPLACABLE_BIT_COUNT can be set to 6 or to 7.
//! * With 7 bits the TiP packets can get max 14% longer and max 127 primary pattern IDs possible.
//! * With 6 bits the TiP packets can get max 33% longer and max 191 primary pattern IDs possible.
//! * See TiP user manual for more information.
#define UNREPLACABLE_BIT_COUNT 7


#ifdef __cplusplus
}
#endif

#endif // TIP_CONFIG_H_
