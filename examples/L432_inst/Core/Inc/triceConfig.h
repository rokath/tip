/*! \file triceConfig.h
\author Thomas.Hoehenleitner [at] seerose.net
*******************************************************************************/

#ifndef TRICE_CONFIG_H_
#define TRICE_CONFIG_H_

#ifdef __cplusplus
extern "C" {
#endif

//! TRICE_CLEAN, if found inside triceConfig.h, is modified by the Trice tool to silent editor warnings in the cleaned state.
#define TRICE_CLEAN 1 // Do not define this at an other place! But you can delete this here.

// hardware specific trice lib settings
#include <stdint.h>

extern volatile uint32_t * const DWT_CONTROL;
extern volatile uint32_t * const DWT_CYCCNT ;
extern volatile uint32_t * const DEMCR;
extern volatile uint32_t * const LAR; // lock access register

#define TriceStamp16 (*DWT_CYCCNT) // @64MHz wraps after a bit more than 1ms (MCU clocks) 
#define TriceStamp32 (*DWT_CYCCNT) // @64MHz -> 1 µs, wraps after 2^32 µs ~= 1.2 hours

#define TRICE_BUFFER TRICE_STATIC_BUFFER

// Windows: trice log -p jlink -args "-Device STM32F030R8" -pf none -prefix off -hs off -d16 -showID "deb:%5d" -i ../../demoTIL.json -li ../../demoLI.json
// Unix:   ./RTTLogUnix.sh or manually:
// 		Terminal 1: rm -f ./temp/trice.bin && JLinkRTTLogger -Device STM32G0B1RE -If SWD -Speed 4000 -RTTChannel 0 ./temp/trice.bin
//      Terminal 2: touch ./temp/trice.bin && trice log -p FILE -args ./temp/trice.bin -pf none -prefix off -hs off -d16 -ts ms -i ../../demoTIL.json -li ../../demoLI.json
#define TRICE_DIRECT_OUTPUT 1
#define TRICE_DIRECT_SEGGER_RTT_32BIT_WRITE 1
#define BUFFER_SIZE_UP (256) // "TRICE_DIRECT_BUFFER_SIZE"

//#include "cmsis_gcc.h"
//#define TRICE_ENTER_CRITICAL_SECTION { uint32_t primaskstate = __get_PRIMASK(); __disable_irq(); {
//#define TRICE_LEAVE_CRITICAL_SECTION } __set_PRIMASK(primaskstate); }

void TriceHeadLine(char* name);
void LogTriceConfiguration(void);

#ifdef __cplusplus
}
#endif

#endif /* TRICE_CONFIG_H_ */
