# Makefile

BUILDFILES = build.sh Makefile mcu.mak flags.mak gcc0.mak gcc1.mak clang0.mak clang1.mak
TIP_VERBOSE = # -v

######################################
# target
######################################
TARGET = G0B1


######################################
# building variables
######################################
# debug build?
TIP_DEBUG = 1


######################################
# source
######################################
# C sources
C_SOURCES =  \
Core/Src/main.c \
Core/Src/app_freertos.c \
Core/Src/stm32g0xx_it.c \
Core/Src/stm32g0xx_hal_msp.c \
Core/Src/stm32g0xx_hal_timebase_tim.c \
Drivers/STM32G0xx_HAL_Driver/Src/stm32g0xx_ll_gpio.c \
Drivers/STM32G0xx_HAL_Driver/Src/stm32g0xx_ll_exti.c \
Drivers/STM32G0xx_HAL_Driver/Src/stm32g0xx_ll_rcc.c \
Drivers/STM32G0xx_HAL_Driver/Src/stm32g0xx_ll_utils.c \
Drivers/STM32G0xx_HAL_Driver/Src/stm32g0xx_ll_crs.c \
Drivers/STM32G0xx_HAL_Driver/Src/stm32g0xx_hal_rcc.c \
Drivers/STM32G0xx_HAL_Driver/Src/stm32g0xx_hal_rcc_ex.c \
Drivers/STM32G0xx_HAL_Driver/Src/stm32g0xx_hal_flash.c \
Drivers/STM32G0xx_HAL_Driver/Src/stm32g0xx_hal_flash_ex.c \
Drivers/STM32G0xx_HAL_Driver/Src/stm32g0xx_hal_gpio.c \
Drivers/STM32G0xx_HAL_Driver/Src/stm32g0xx_hal_dma.c \
Drivers/STM32G0xx_HAL_Driver/Src/stm32g0xx_hal_dma_ex.c \
Drivers/STM32G0xx_HAL_Driver/Src/stm32g0xx_ll_dma.c \
Drivers/STM32G0xx_HAL_Driver/Src/stm32g0xx_hal_pwr.c \
Drivers/STM32G0xx_HAL_Driver/Src/stm32g0xx_hal_pwr_ex.c \
Drivers/STM32G0xx_HAL_Driver/Src/stm32g0xx_hal_cortex.c \
Drivers/STM32G0xx_HAL_Driver/Src/stm32g0xx_hal.c \
Drivers/STM32G0xx_HAL_Driver/Src/stm32g0xx_hal_exti.c \
Drivers/STM32G0xx_HAL_Driver/Src/stm32g0xx_hal_tim.c \
Drivers/STM32G0xx_HAL_Driver/Src/stm32g0xx_hal_tim_ex.c \
Drivers/STM32G0xx_HAL_Driver/Src/stm32g0xx_ll_usart.c \
Core/Src/system_stm32g0xx.c \
Middlewares/Third_Party/FreeRTOS/Source/croutine.c \
Middlewares/Third_Party/FreeRTOS/Source/event_groups.c \
Middlewares/Third_Party/FreeRTOS/Source/list.c \
Middlewares/Third_Party/FreeRTOS/Source/queue.c \
Middlewares/Third_Party/FreeRTOS/Source/stream_buffer.c \
Middlewares/Third_Party/FreeRTOS/Source/tasks.c \
Middlewares/Third_Party/FreeRTOS/Source/timers.c \
Middlewares/Third_Party/FreeRTOS/Source/CMSIS_RTOS/cmsis_os.c \
Middlewares/Third_Party/FreeRTOS/Source/portable/MemMang/heap_4.c \
Middlewares/Third_Party/FreeRTOS/Source/portable/GCC/ARM_CM0/port.c

#Core/Src/sysmem.c 
#Core/Src/syscalls.c

# trice lib
C_SOURCES += \
../../../trice/src/cobsDecode.c \
../../../trice/src/cobsEncode.c \
../../../trice/src/SEGGER_RTT.c \
../../../trice/src/tcobsv1Decode.c \
../../../trice/src/tcobsv1Encode.c \
../../../trice/src/trice.c \
../../../trice/src/trice8.c \
../../../trice/src/trice16.c \
../../../trice/src/trice32.c \
../../../trice/src/trice64.c \
../../../trice/src/triceUart.c \
../../../trice/src/triceAuxiliary.c \
../../../trice/src/triceDoubleBuffer.c \
../../../trice/src/triceRingBuffer.c \
../../../trice/src/triceStackBuffer.c \
../../../trice/src/triceStaticBuffer.c \
../../../trice/src/xtea.c 

# examples
C_SOURCES += \
../../../trice/examples/exampleData/triceLogDiagData.c \
../../../trice/examples/exampleData/triceExamples.c \

C_SOURCES += \
../../src/idTable.c \
../../src/memmem.c \
../../src/shift68bit.c \
../../src/shift78bit.c \
../../src/shift87bit.c \
../../src/shift86bit.c \
../../src/ti_pack.c \
../../src/ti_unpack.c 

# ASM sources
ASM_SOURCES =  \
startup_stm32g0b1xx.s

include gcc0.mak
include clang0.mak
include mcu.mak

#######################################
# CFLAGS
#######################################


# macros for gcc
# AS defines
AS_DEFS = 

# C defines
C_DEFS =  \
-DUSE_FULL_LL_DRIVER \
-DSTM32G0B1xx \
-DUSE_HAL_DRIVER \
$(TRICE_FLAGS)


# AS includes
AS_INCLUDES =  \
-ICore/Inc

# C includes
C_INCLUDES =  \
-ICore/Inc \
-IDrivers/STM32G0xx_HAL_Driver/Inc \
-IDrivers/STM32G0xx_HAL_Driver/Inc/Legacy \
-IMiddlewares/Third_Party/FreeRTOS/Source/include \
-IMiddlewares/Third_Party/FreeRTOS/Source/CMSIS_RTOS \
-IMiddlewares/Third_Party/FreeRTOS/Source/portable/GCC/ARM_CM0 \
-IDrivers/CMSIS/Device/ST/STM32G0xx/Include \
-IDrivers/CMSIS/Include

C_INCLUDES += \
-I../../../trice/src \
-I../../src

include flags.mak

.PHONY: all gcc clean log version

all: gcc # gcc clang

clang: $(CLANG_BUILD)/$(TARGET).elf $(CLANG_BUILD)/$(TARGET).hex $(CLANG_BUILD)/$(TARGET).bin
	@:

gcc:     $(GCC_BUILD)/$(TARGET).elf   $(GCC_BUILD)/$(TARGET).hex   $(GCC_BUILD)/$(TARGET).bin
	@:

log:
	trice l -p JLINK -args="-Device STM32G0B1RE -if SWD -Speed 4000 -RTTChannel 0" -pf none -ts ms -d16

#######################################
# build the application
#######################################
include gcc1.mak
include clang1.mak

clean:
	@rm -fR $(GCC_BUILD) $(CLANG_BUILD)
	
version:
	@which arm-none-eabi-gcc
	@arm-none-eabi-gcc --version
	@which $(CLANG_CC)
	@$(CLANG_CC) --version
 
# openocdflashgcc: all
# 	openocd -f interface/jlink.cfg -f target/stm32g0x.cfg -c "program $(GCC_BUILD)/$(TARGET).elf verify reset exit"
# 
# openocdflashclang: all
# 	openocd -f interface/jlink.cfg -f target/stm32g0x.cfg -c "program $(CLANG_BUILD)/$(TARGET).elf verify reset exit"
# 
# seggerflashgcc: all
# 	JFlash -openprj G0B1.jflash -open $(GCC_BUILD)/G0B1.hex,0x08000000 -auto -startapp -exit
# 
# seggerflashclang: all
# 	JFlash -openprj G0B1.jflash -open $(CLANG_BUILD)/G0B1.hex,0x08000000 -auto -startapp -exit
# 
# # See UM0892
# stflashgcc: all
# 	ST-LINK_CLI.exe -c SWD FREQ=4000 UR -P $(GCC_BUILD)/G0B1.hex -V
# 
# # See UM0892
# stflashclang: all
# 	ST-LINK_CLI.exe -c SWD FREQ=4000 UR -P $(CLANG_BUILD)/G0B1.hex -V
