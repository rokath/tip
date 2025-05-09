/* USER CODE BEGIN Header */
/**
 ******************************************************************************
 * @file           : main.c
 * @brief          : Main program body
 ******************************************************************************
 * @attention
 *
 * Copyright (c) 2023 STMicroelectronics.
 * All rights reserved.
 *
 * This software is licensed under terms that can be found in the LICENSE file
 * in the root directory of this software component.
 * If no LICENSE file comes with this software, it is provided AS-IS.
 *
 ******************************************************************************
 */
/* USER CODE END Header */
/* Includes ------------------------------------------------------------------*/
#include "main.h"
#include "cmsis_os.h"

/* Private includes ----------------------------------------------------------*/
/* USER CODE BEGIN Includes */
#include "trice.h"
#include "ti_pack.h"
#include "ti_unpack.h"
//#include <limits.h> // INT_MAX
/* USER CODE END Includes */

/* Private typedef -----------------------------------------------------------*/
/* USER CODE BEGIN PTD */

/* USER CODE END PTD */

/* Private define ------------------------------------------------------------*/
/* USER CODE BEGIN PD */

/* USER CODE END PD */

/* Private macro -------------------------------------------------------------*/
/* USER CODE BEGIN PM */

/* USER CODE END PM */

/* Private variables ---------------------------------------------------------*/

osThreadId defaultTaskHandle;
osThreadId myTask02_DiagnoHandle;
/* USER CODE BEGIN PV */
/* USER CODE END PV */

/* Private function prototypes -----------------------------------------------*/
void SystemClock_Config(void);
static void MX_GPIO_Init(void);
static void MX_USART2_UART_Init(void);
void StartDefaultTask(void const * argument);
void StartTask02(void const * argument);

/* USER CODE BEGIN PFP */

/* USER CODE END PFP */

/* Private user code ---------------------------------------------------------*/
/* USER CODE BEGIN 0 */
__weak int _close(void) { return -1; }
__weak int _lseek(void) { return -1; }
__weak int _read (void) { return -1; }
__weak int _write(void) { return -1; }


/* USER CODE END 0 */

/**
  * @brief  The application entry point.
  * @retval int
  */
int main(void)
{

  /* USER CODE BEGIN 1 */

  /* USER CODE END 1 */

  /* MCU Configuration--------------------------------------------------------*/

  /* Reset of all peripherals, Initializes the Flash interface and the Systick. */
  HAL_Init();

  /* USER CODE BEGIN Init */
#if !TRICE_OFF
  TriceInit(); // This so early, to allow trice logs inside interrupts from the beginning. Only needed for RTT.
  trice("msg:Hi\n");
//TriceHeadLine("  𝕹𝖀𝕮𝕷𝕰𝕺-G01BRE with tip  ");
#endif
  /* USER CODE END Init */

  /* Configure the system clock */
  SystemClock_Config();

  /* USER CODE BEGIN SysInit */

  /* USER CODE END SysInit */

  /* Initialize all configured peripherals */
  MX_GPIO_Init();
  MX_USART2_UART_Init();
  /* USER CODE BEGIN 2 */

  /* USER CODE END 2 */

  /* USER CODE BEGIN RTOS_MUTEX */
  /* add mutexes, ... */
  /* USER CODE END RTOS_MUTEX */

  /* USER CODE BEGIN RTOS_SEMAPHORES */
  /* add semaphores, ... */
  /* USER CODE END RTOS_SEMAPHORES */

  /* USER CODE BEGIN RTOS_TIMERS */
  /* start timers, add new ones, ... */
  /* USER CODE END RTOS_TIMERS */

  /* USER CODE BEGIN RTOS_QUEUES */
  /* add queues, ... */
  /* USER CODE END RTOS_QUEUES */

  /* Create the thread(s) */
  /* definition and creation of defaultTask */
  osThreadDef(defaultTask, StartDefaultTask, osPriorityNormal, 0, 256);
  defaultTaskHandle = osThreadCreate(osThread(defaultTask), NULL);

  /* definition and creation of myTask02_Diagno */
  osThreadDef(myTask02_Diagno, StartTask02, osPriorityIdle, 0, 256);
  myTask02_DiagnoHandle = osThreadCreate(osThread(myTask02_Diagno), NULL);

  /* USER CODE BEGIN RTOS_THREADS */
  /* add threads, ... */
  /* USER CODE END RTOS_THREADS */

  /* Start scheduler */
  osKernelStart();
  /* We should never get here as control is now taken by the scheduler */
  /* Infinite loop */
  /* USER CODE BEGIN WHILE */
  while (1)
  {
    /* USER CODE END WHILE */

    /* USER CODE BEGIN 3 */
  }
  /* USER CODE END 3 */
}

/**
  * @brief System Clock Configuration
  * @retval None
  */
void SystemClock_Config(void)
{
  /* HSI configuration and activation */
  LL_RCC_HSI_Enable();
  while(LL_RCC_HSI_IsReady() != 1)
  {
  }

  /* Set AHB prescaler*/
  LL_RCC_SetAHBPrescaler(LL_RCC_SYSCLK_DIV_1);

  /* Sysclk activation on the HSI */
  LL_RCC_SetSysClkSource(LL_RCC_SYS_CLKSOURCE_HSI);
  while(LL_RCC_GetSysClkSource() != LL_RCC_SYS_CLKSOURCE_STATUS_HSI)
  {
  }

  /* Set APB1 prescaler*/
  LL_RCC_SetAPB1Prescaler(LL_RCC_APB1_DIV_1);
  /* Update CMSIS variable (which can be updated also through SystemCoreClockUpdate function) */
  LL_SetSystemCoreClock(16000000);

   /* Update the time base */
  if (HAL_InitTick (TICK_INT_PRIORITY) != HAL_OK)
  {
    Error_Handler();
  }
}

/**
  * @brief USART2 Initialization Function
  * @param None
  * @retval None
  */
static void MX_USART2_UART_Init(void)
{

  /* USER CODE BEGIN USART2_Init 0 */

  /* USER CODE END USART2_Init 0 */

  LL_USART_InitTypeDef USART_InitStruct = {0};

  LL_GPIO_InitTypeDef GPIO_InitStruct = {0};

  LL_RCC_SetUSARTClockSource(LL_RCC_USART2_CLKSOURCE_PCLK1);

  /* Peripheral clock enable */
  LL_APB1_GRP1_EnableClock(LL_APB1_GRP1_PERIPH_USART2);

  LL_IOP_GRP1_EnableClock(LL_IOP_GRP1_PERIPH_GPIOA);
  /**USART2 GPIO Configuration
  PA2   ------> USART2_TX
  PA3   ------> USART2_RX
  */
  GPIO_InitStruct.Pin = USART2_TX_Pin;
  GPIO_InitStruct.Mode = LL_GPIO_MODE_ALTERNATE;
  GPIO_InitStruct.Speed = LL_GPIO_SPEED_FREQ_LOW;
  GPIO_InitStruct.OutputType = LL_GPIO_OUTPUT_PUSHPULL;
  GPIO_InitStruct.Pull = LL_GPIO_PULL_NO;
  GPIO_InitStruct.Alternate = LL_GPIO_AF_1;
  LL_GPIO_Init(USART2_TX_GPIO_Port, &GPIO_InitStruct);

  GPIO_InitStruct.Pin = USART2_RX_Pin;
  GPIO_InitStruct.Mode = LL_GPIO_MODE_ALTERNATE;
  GPIO_InitStruct.Speed = LL_GPIO_SPEED_FREQ_LOW;
  GPIO_InitStruct.OutputType = LL_GPIO_OUTPUT_PUSHPULL;
  GPIO_InitStruct.Pull = LL_GPIO_PULL_NO;
  GPIO_InitStruct.Alternate = LL_GPIO_AF_1;
  LL_GPIO_Init(USART2_RX_GPIO_Port, &GPIO_InitStruct);

  /* USART2 interrupt Init */
  NVIC_SetPriority(USART2_LPUART2_IRQn, 3);
  NVIC_EnableIRQ(USART2_LPUART2_IRQn);

  /* USER CODE BEGIN USART2_Init 1 */

  /* USER CODE END USART2_Init 1 */
  USART_InitStruct.PrescalerValue = LL_USART_PRESCALER_DIV1;
  USART_InitStruct.BaudRate = 115200;
  USART_InitStruct.DataWidth = LL_USART_DATAWIDTH_8B;
  USART_InitStruct.StopBits = LL_USART_STOPBITS_1;
  USART_InitStruct.Parity = LL_USART_PARITY_NONE;
  USART_InitStruct.TransferDirection = LL_USART_DIRECTION_TX_RX;
  USART_InitStruct.HardwareFlowControl = LL_USART_HWCONTROL_NONE;
  USART_InitStruct.OverSampling = LL_USART_OVERSAMPLING_16;
  LL_USART_Init(USART2, &USART_InitStruct);
  LL_USART_SetTXFIFOThreshold(USART2, LL_USART_FIFOTHRESHOLD_1_8);
  LL_USART_SetRXFIFOThreshold(USART2, LL_USART_FIFOTHRESHOLD_1_8);
  LL_USART_DisableFIFO(USART2);
  LL_USART_ConfigAsyncMode(USART2);

  /* USER CODE BEGIN WKUPType USART2 */

  /* USER CODE END WKUPType USART2 */

  LL_USART_Enable(USART2);

  /* Polling USART2 initialisation */
  while((!(LL_USART_IsActiveFlag_TEACK(USART2))) || (!(LL_USART_IsActiveFlag_REACK(USART2))))
  {
  }
  /* USER CODE BEGIN USART2_Init 2 */

  /* USER CODE END USART2_Init 2 */

}

/**
  * @brief GPIO Initialization Function
  * @param None
  * @retval None
  */
static void MX_GPIO_Init(void)
{
  LL_GPIO_InitTypeDef GPIO_InitStruct = {0};
/* USER CODE BEGIN MX_GPIO_Init_1 */
/* USER CODE END MX_GPIO_Init_1 */

  /* GPIO Ports Clock Enable */
  LL_IOP_GRP1_EnableClock(LL_IOP_GRP1_PERIPH_GPIOC);
  LL_IOP_GRP1_EnableClock(LL_IOP_GRP1_PERIPH_GPIOF);
  LL_IOP_GRP1_EnableClock(LL_IOP_GRP1_PERIPH_GPIOA);

  /**/
  LL_GPIO_ResetOutputPin(LED_GREEN_GPIO_Port, LED_GREEN_Pin);

  /**/
  GPIO_InitStruct.Pin = LED_GREEN_Pin;
  GPIO_InitStruct.Mode = LL_GPIO_MODE_OUTPUT;
  GPIO_InitStruct.Speed = LL_GPIO_SPEED_FREQ_HIGH;
  GPIO_InitStruct.OutputType = LL_GPIO_OUTPUT_PUSHPULL;
  GPIO_InitStruct.Pull = LL_GPIO_PULL_NO;
  LL_GPIO_Init(LED_GREEN_GPIO_Port, &GPIO_InitStruct);

/* USER CODE BEGIN MX_GPIO_Init_2 */
/* USER CODE END MX_GPIO_Init_2 */
}

/* USER CODE BEGIN 4 */
/* USER CODE END 4 */

/* USER CODE BEGIN Header_StartDefaultTask */
/**
  * @brief  Function implementing the defaultTask thread.
  * @param  argument: Not used
  * @retval None
  */
/* USER CODE END Header_StartDefaultTask */
void StartDefaultTask(void const * argument)
{
  /* USER CODE BEGIN 5 */
  osDelay(100);
  TRICE_UNUSED(argument)
  TRice("msg:StartDefaultTask\n");


  //  uint64_t sum = 0;
  //  for( uint64_t i = 0; i < 1000000u; i++ ){
  //    sum += i;
  //  }
  //  trice64("Hi, sum is %u\n", sum);
  //BUG!!!TriceHeadLine("  NUCLEO-G0B1RE  123");
//LogTriceConfiguration();


static uint8_t IdT[] = {3, 0xaa, 0xaa, 0xaa, 0}; // idTable

//ok static uint8_t src[] = { 0xd1, 0xaa, 0xaa, 0xaa, 0xd2};
//ok static uint8_t exp[] = { 0xe0, 0x01, 0xd1, 0xd2};

static uint8_t src[] = {0xd1, 0xaa, 0xaa, 0xaa, 0xd2, 0xaa, 0xaa, 0xaa, 0xd3};
static uint8_t exp[] = {0xf0, 0x01, 0xd1, 0x01, 0xd2, 0xd3}; // expected

static uint8_t dst[100] = {0};

//static uint8_t buf[100] = {0};
//static uint8_t pkg[100] = {0};
//static uint8_t src[] = { 0x0d, 0x0a, 0x0d, 0x0a, 0x74, 0x68, 0x65, 0x20, 0x0d, 0x0a }; // ok
//static uint8_t src[] = { 0x55, 0x74, 0x68, 0x65, 0x20, 0x55, 0x74, 0x68, 0x65, 0x20 }; // ok
//static uint8_t src[] = { 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0xaa, 0xbb, 0x38, 0x39, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66 };
//static char* in = "ABCABCABC123";
//static uint8_t src[100] = {0};
//static const uint8_t IdT[] = {5, 0, 0, 0, 0, 0, 3, 0xee, 0xff, 0xaa, 2, 0xaa, 0xbb, 0};
//static       uint8_t src[] = {0xff, 0x00, 0xaa, 0xbb, 0xee, 0xff, 0xaa, 0xbb, 0xcc};

size_t slen =  sizeof(src); //   strlen(in);

// ok newIDPosTable(IdT, src, slen );
// ok trice8B("rd:%02x \n", &IDPosTable.item[0], IDPosTable.count * sizeof(IDPosition_t));
// ok for( int i = 0; i < IDPosTable.count; i++ ){
// ok   trice("%2d: id%2d, start=%d, limit=%d \n", i, IDPosTable.item[i].id, IDPosTable.item[i].start, IDPosLimit(i));
// ok }
// ok createSrcMap(IdT, src, slen );
// ok for( int i = 0; i < srcMap.count; i++ ){
// ok     int plen = srcMap.path[i][0];
// ok     trice("att:path %2d: len %2d: ", i, plen);
// ok     trice8B("msg:%3d\n", &srcMap.path[i][1], plen);
// ok }

trice8B("buf: %02x\n", src, slen);
trice8B("exp: %02x\n", exp, sizeof(exp) ); 
size_t dlen = tiPack( dst, src, slen, 0, 0, 0, IdT );
trice8B("tip: %02x\n", dst, dlen ); 


//size_t plen = tip(pkg, src, slen);
//trice8B("wr:%02x \n", pkg, plen);

//size_t blen = tiu(buf, pkg, plen);
//trice8B("att:%02x \n", buf, blen);


/*
  static uint8_t buf[100] = {0};
  static uint8_t pkg[100] = {0};
//static uint8_t src[] = { 0x0d, 0x0a, 0x0d, 0x0a, 0x74, 0x68, 0x65, 0x20, 0x0d, 0x0a }; // ok
//static uint8_t src[] = { 0x55, 0x74, 0x68, 0x65, 0x20, 0x55, 0x74, 0x68, 0x65, 0x20 }; // ok
  static uint8_t src[] = { 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0xaa, 0xbb, 0x38, 0x39, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66 };
  //static char* in = "ABCABCABC123";
  //static uint8_t src[100] = {0};
  size_t slen =  sizeof(src); //   strlen(in);
  //  memcpy(src,in,slen);
  trice8B("rd:%02x \n", src, slen);

  size_t plen = tip(pkg, src, slen);

  trice8B("wr:%02x \n", pkg, plen);
  
  size_t blen = tiu(buf, pkg, plen);
  trice8B("att:%02x \n", buf, blen);
*/



  /* Infinite loop */
  for(;;)
  {
    /*
    static uint8_t dst[100] = {0};
    static uint8_t src[] = {0xff, 0xff, 0xff, 0xff, 0xff, 0xff};
    size_t slen = sizeof(src);
    size_t dlen = tip(dst, src, slen);
    trice8B("rd:%02x ", src, slen);
    trice8B("wr:%02x ", dst, dlen);
	*/
    osDelay(1000);
  }
  /* USER CODE END 5 */
}

/* USER CODE BEGIN Header_StartTask02 */
/**
* @brief Function implementing the myTask02_Diagno thread.
* @param argument: Not used
* @retval None
*/
/* USER CODE END Header_StartTask02 */
void StartTask02(void const * argument)
{
  /* USER CODE BEGIN StartTask02 */
  TRICE_UNUSED(argument)
  TRice("msg:StartTask02:Diagnostics and TriceTransfer\n" );
  /* Infinite loop */
  for(;;)
  {
#if !TRICE_OFF

//#if TRICE_DIAGNOSTICS == 1
//    static int i = INT_MAX;
//    if( i++ > 100 ){
//      i = 0;
//      TriceLogDiagnosticData();
//    }
//#endif // #if TRICE_DIAGNOSTICS == 1

    //TriceTransfer();
    osDelay(100);

#if TRICE_BUFFER == TRICE_RING_BUFFER && TRICE_RING_BUFFER_OVERFLOW_WATCH == 1
    WatchRingBufferMargins();
#endif // #if TRICE_RING_BUFFER_OVERFLOW_WATCH == 1

#endif // #if !TRICE_OFF
  }
  /* USER CODE END StartTask02 */
}

/**
  * @brief  Period elapsed callback in non blocking mode
  * @note   This function is called  when TIM17 interrupt took place, inside
  * HAL_TIM_IRQHandler(). It makes a direct call to HAL_IncTick() to increment
  * a global variable "uwTick" used as application time base.
  * @param  htim : TIM handle
  * @retval None
  */
void HAL_TIM_PeriodElapsedCallback(TIM_HandleTypeDef *htim)
{
  /* USER CODE BEGIN Callback 0 */

  /* USER CODE END Callback 0 */
  if (htim->Instance == TIM17) {
    HAL_IncTick();
  }
  /* USER CODE BEGIN Callback 1 */

  /* USER CODE END Callback 1 */
}

/**
  * @brief  This function is executed in case of error occurrence.
  * @retval None
  */
void Error_Handler(void)
{
  /* USER CODE BEGIN Error_Handler_Debug */
  /* User can add his own implementation to report the HAL error return state */
  __disable_irq();
  while (1)
  {
  }
  /* USER CODE END Error_Handler_Debug */
}

#ifdef  USE_FULL_ASSERT
/**
  * @brief  Reports the name of the source file and the source line number
  *         where the assert_param error has occurred.
  * @param  file: pointer to the source file name
  * @param  line: assert_param error line source number
  * @retval None
  */
void assert_failed(uint8_t *file, uint32_t line)
{
  /* USER CODE BEGIN 6 */
  /* User can add his own implementation to report the file name and line number,
     ex: printf("Wrong parameters value: file %s on line %d\r\n", file, line) */
  /* USER CODE END 6 */
}
#endif /* USE_FULL_ASSERT */
