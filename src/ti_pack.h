//! @file ti_pack.h
//! @brief This is the tip user interface for packing.
//! @author thomas.hoehenleitner [at] seerose.net

#ifndef TI_PACK_H_
#define TI_PACK_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>
#include <stddef.h>

//! @brief tip encodes src buffer with size len into dst buffer and returns encoded len.
//! @details For the tip encoding it uses the linked idTable.c object.
//! ATTENTION: The pack functions are usable only sequentially!
//! tip is the default user interface.
size_t tip( uint8_t* dst, const uint8_t * src, size_t len );

//! @brief tiPack encodes src buffer with size len into dst buffer and returns encoded len.
//! @details For the tip encoding it uses the passed ID table object.
//! - Some bytes groups in the src buffer are replacable with IDs 0x01...0x7f and some not.
//! - The unreplacable bytes are temporarily collected into a buffer and transformed to free and 
//! set all most significant bits to distinguish them later from the IDs. 
//! - Afterwards the replacement IDs and the transformed unreplacable bytes are mixed according to their locations.
//! ATTENTION: The pack functions are usable only sequentially!
//! tiPack is the extended user interface to use different ID tables in the same project.
size_t tiPack( uint8_t * dst, const uint8_t * table, const uint8_t * src, size_t slen );

#ifdef __cplusplus
}
#endif

#endif // TI_PACK_H_
