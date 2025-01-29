//! @file tipConfig.h
//! @brief This is the tip configuration.
//! @author thomas.hoehenleitner [at] seerose.net

#ifndef TIP_CONFIG_H_
#define TIP_CONFIG_H_

#ifdef __cplusplus
extern "C" {
#endif

//! TIP_SRC_BUFFER_SIZE_MAX is the maximun allowed input data buffer length.
#define TIP_SRC_BUFFER_SIZE_MAX 0xffffffff

//! TIP_PATTERN_SIZE_MAX is the max replacement pattern size.
//! The tipTable generator uses this value.
#define TIP_PATTERN_SIZE_MAX 8

#ifdef __cplusplus
}
#endif

#endif // TIP_CONFIG_H_
