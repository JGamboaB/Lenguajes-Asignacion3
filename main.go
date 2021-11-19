package main

//go mod init main.go
//go mod tidy

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/lxn/win"

	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

// Variables
const (
	BAR_WIDTH = 1
	FONT_WIDTH = 8 //8
	FONT_HEIGHT = 16
	MAX_NUMBER_SIZE = 32
)

var (
	width int = int(win.GetSystemMetrics(win.SM_CXSCREEN) / FONT_WIDTH)
	height int = int(win.GetSystemMetrics(win.SM_CYSCREEN) / (FONT_HEIGHT*3))
	m sync.Mutex

	// Charts
	bsChart widgets.BarChart
	ssChart widgets.BarChart
	isChart widgets.BarChart
	qsChart widgets.BarChart
	hsChart widgets.BarChart

	//Values per algorithm
	bsTime time.Duration
	bsSwaps = 0
	bsComparisons = 0
	bsIterations = 0

	ssTime time.Duration
	ssSwaps = 0
	ssComparisons = 0
	ssIterations = 0

	isTime time.Duration
	isSwaps = 0
	isComparisons = 0
	isIterations = 0

	qsTime time.Duration
	qsSwaps = 0
	qsComparisons = 0
	qsIterations = 0

	hsTime time.Duration
	hsSwaps = 0
	hsComparisons = 0
	hsIterations = 0
)

func main(){
	//Pedir n (intervalo 10 a 100)
	//Semilla: entrada de datos por teclado o segundos y decimas del reloj del sistema
	//Arreglo base

	//arregloBase := temporalRANDOM(100)
	fmt.Print("Indique la cantidad de numeros(Se recomienda 100 maximo para una visualizacion correcta): ")

	var size int
	fmt.Scanln(&size)
	if err := ui.Init(); err != nil{
		log.Fatalf("failed to initialize termui: %v", err)
	}

	//Slice
	arregloBase := temporalRANDOM(size) // CAMBIAR TEMPORAL

	//Init
	initBSChart(arregloBase)
	initSSChart(arregloBase)
	initISChart(arregloBase)
	initQSChart(arregloBase)
	initHSChart(arregloBase)

	ui.Render(&bsChart)
	ui.Render(&ssChart)
	ui.Render(&isChart)
	ui.Render(&qsChart)
	ui.Render(&hsChart)

	//Goroutines & Drawing Start
	go bsChartDrawer(arregloBase)
	go ssChartDrawer(arregloBase)
	go isChartDrawer(arregloBase)
	go qsChartDrawer(arregloBase)
	hsChartDrawer(arregloBase)

	//Ending
	fmt.Scanln()
	ui.Close()
}

func temporalRANDOM(n int) []float64{
	arr := make([]float64, n)
	for i := 0; i < n; i++{
		arr[i] = float64(rand.Intn(30))
	}; return arr
}

// / / / / / / Sorting \ \ \ \ \ \ 

// Bubblesort >> New <<

func bubbleSort(arr *[]float64, pair chan []int) {
	initTime := time.Now() // Time Start

	arr2 := *arr
	len := len(arr2)

	//Move through all elements
	for i:= 0; i < len; i++ {
		for j := 0; j < len-i-1; j++ {
			// Move from 0 to len-i-1 and swap if element is greater than the next one
			if arr2[j] > arr2[j+1] {
				arr2[j], arr2[j+1] = arr2[j+1], arr2[j]

				pair <- []int{j, j+1} // Channel
				bsSwaps++
			}; bsComparisons++; bsIterations++
		}; bsIterations++
	}

	//*arr = arr2 // Assign changes to original array
	close(pair)

	endTime := time.Now() // Time End
	bsTime = endTime.Sub(initTime) // Total Time
}

// Selection (changes indexes)

func selectionSort(arr *[]float64, pair chan []int) {
	initTime := time.Now() // Time Start

	arr2 := *arr
	len := len(arr2)

	for currentIndex := 0; currentIndex < len-1; currentIndex++ { // Done to all the indexes in the array
		indexMin := currentIndex

		for i := currentIndex + 1; i < len; i++ { // Get the index of the smallest value from the numbers to the right
			if arr2[i] < arr2[indexMin] {
				indexMin = i
			}; ssComparisons++; ssIterations++
		}

		//Swap numbers
		arr2[currentIndex], arr2[indexMin] = arr2[indexMin], arr2[currentIndex]
		pair <- []int{currentIndex, indexMin} // Channel

		ssSwaps++
		ssIterations++
	}

	*arr = arr2 // Assign changes to original array
	close(pair)

	endTime := time.Now() // Time End
	ssTime = endTime.Sub(initTime) // Total Time
}

// Insertion >> New <<

func insertionSort(arr *[]float64, oneWay chan []int) {
	initTime := time.Now() // Time Start

	arr2 := *arr
	len := len(arr2)

	for i := 1; i < len; i++ {
		key := arr2[i]
		j := i-1

		//Move greater elements of arr[0 .. i-1] to position ahead of current
		isComparisons++ //isComparison????
		for ; j >= 0 && key < arr2[j]; j--{
			oneWay <- []int{j+1, int(arr2[j])} // Channel
			arr2[j+1] = arr2[j]

			isSwaps++
			isIterations++
		}

		oneWay <- []int{j+1, int(key)} // Channel
		arr2[j+1] = key

		isSwaps++
		isIterations++
	}
	

	*arr = arr2 // Assign changes to original array
	close(oneWay)
	
	endTime := time.Now() // Time End
	isTime = endTime.Sub(initTime) // Total Time
}

// Quicksort [iterative for drawing]: https://www.geeksforgeeks.org/iterative-quick-sort/

func partition(arr *[]float64, low int, high int, pair chan []int) int { //
	arr2 := *arr
	pivot := arr2[high]

	i := low - 1

	for j := low; j < high; j++ {
		if arr2[j] <= pivot {
			i++
			arr2[i], arr2[j] = arr2[j], arr2[i] //Gets the lesser values to the left of the pivot
			pair <- []int{i, j} // Channel

			qsSwaps++
		}; qsComparisons++; qsIterations++
	}

	//Swap pivot with the next element to i
	arr2[i+1], arr2[high] = arr2[high], arr2[i+1]
	pair <- []int{i+1, high} // Channel

	qsSwaps++

	*arr = arr2 // Assign changes to original array

	return i + 1 //new pivot
}

func quickSort(arr *[]float64, pair chan []int){
	initTime := time.Now() // Time Start

	low := 0; high := len(*arr) - 1
	stack := make([]int, high+1) //Auxiliary stack

	top := -1 //Top of stack

	//Push high & low
	top++
	stack[top] = low
	top++
	stack[top] = high

	for top >= 0 {
		//Pop high & low
		high = stack[top]
		top--
		low = stack[top]
		top--

		pivot := partition(arr, low, high, pair) //pivot at correct position

		if pivot-1 > low { //If elements on left push left side to stack
			top++
			stack[top] = low
			top++
			stack[top] = pivot-1
		}; qsComparisons++

		if pivot+1 < high{
			top++
			stack[top] = pivot+1
			top++
			stack[top] = high
		}; qsComparisons++

		qsIterations++
	}

	close(pair)

	endTime := time.Now() // Time End
	qsTime = endTime.Sub(initTime) // Total Time
}

// Heapsort >> New << [iteraive just in case]: https://www.geeksforgeeks.org/iterative-heap-sort/

func buildMaxHeap(arr *[]float64, n int, pair chan []int){
	arr2 := *arr
	for i := 1; i < n; i++{
		if arr2[i] > arr2[(i-1)/2]{ // Child bigger than parent
			j := i

			for arr2[j] > arr2[(j-1)/2]{ //Swap until parent is smaller than child
				arr2[j], arr2[(j-1)/2] = arr2[(j-1)/2], arr2[j]
				pair <- []int{j, (j-1)/2} // Channel
				j = (j-1)/2

				hsSwaps++
			}; hsIterations++
		}; hsComparisons++

		hsIterations++
	}

	*arr = arr2
}

func heapSort(arr *[]float64, pair chan []int){
	initTime := time.Now() // Time Start

	n := len(*arr)

	buildMaxHeap(arr, n, pair)
	arr2 := *arr

	for i := n-1; i > 0; i--{
		arr2[0], arr2[i] = arr2[i], arr2[0] //swap first with last
		pair <- []int{0, i} // Channel
		j, index := 0, 0

		hsSwaps++

		for {
			hsIterations++

			index = 2 * j + 1

			if index < (i - 1) && arr2[index] < arr2[index + 1]{
				index++
			}; hsComparisons+=2
			if index < i && arr2[j] < arr2[index]{
				arr2[j], arr2[index] = arr2[index], arr2[j]
				pair <- []int{j, index} // Channel

				hsSwaps++
			}; j = index
			hsComparisons++
			if index >= i{
				break
			}
			
		}; hsIterations++
	}

	*arr = arr2
	close(pair)

	endTime := time.Now() // Time End
	hsTime = endTime.Sub(initTime) // Total Time
}


// / / / / / / Graphic \ \ \ \ \ \

// / / / / Initialize 

func initBSChart(arr []float64){
	bsChart = *widgets.NewBarChart()
	bsChart.Data = arr
	bsChart.BarWidth = BAR_WIDTH
	bsChart.BarGap = 0
	
	//Changes per Chart
	bsChart.Title = "BubbleSort"
	bsChart.SetRect(0, 0, width/2 - 2, height-2)
	bsChart.BarColors = []ui.Color{ui.ColorRed}
	bsChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorRed)} // Can't be seen

	// Indexes
	//bsChart.Labels = generateLabels(arr)
	//bsChart.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorWhite)}
}

func initSSChart(arr []float64){
	ssChart = *widgets.NewBarChart()
	ssChart.Data = arr
	ssChart.BarWidth = BAR_WIDTH
	ssChart.BarGap = 0

	//Changes per Chart
	ssChart.Title = "SelectionSort"
	ssChart.SetRect(width/2, 0, width - 4, height-2)
	ssChart.BarColors = []ui.Color{ui.ColorCyan}
	ssChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorCyan)} // Can't be seen
}

func initISChart(arr []float64){
	isChart = *widgets.NewBarChart()
	isChart.Data = arr
	isChart.BarWidth = BAR_WIDTH
	isChart.BarGap = 0

	//Changes per Chart
	isChart.Title = "InsertionSort"
	isChart.SetRect(0, height-2, width/2 - 2, height*2-4)
	isChart.BarColors = []ui.Color{ui.ColorGreen}
	isChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorGreen)} // Can't be seen
}

func initQSChart(arr []float64){
	qsChart = *widgets.NewBarChart()
	qsChart.Data = arr
	qsChart.BarWidth = BAR_WIDTH
	qsChart.BarGap = 0

	//Changes per Chart
	qsChart.Title = "QuickSort"
	qsChart.SetRect(width/2, height-2, width - 4, height*2-4)
	qsChart.BarColors = []ui.Color{ui.ColorMagenta}
	qsChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorMagenta)} // Can't be seen
}

func initHSChart(arr []float64){
	hsChart = *widgets.NewBarChart()
	hsChart.Data = arr
	hsChart.BarWidth = BAR_WIDTH
	hsChart.BarGap = 0

	//Changes per Chart
	hsChart.Title = "HeapSort"
	hsChart.SetRect(0, height*2-4, width-4, height*3-1)
	hsChart.BarColors = []ui.Color{ui.ColorYellow}
	hsChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorYellow)} // Can't be seen
}

// / / / / Drawing 

func bsChartDrawer(slice []float64){
	// / / bsChart.Data = slice
	bsChart.Data = make([]float64, len(slice))
	copy(bsChart.Data, slice)

	//Copy used in BubbleSort
	copyArr := make([]float64, len(slice))
	copy(copyArr, bsChart.Data)

	//Channel
	pairsChannel := make(chan []int, 1000)
	go bubbleSort(&copyArr, pairsChannel)

	//Update Changes in pairs
	for pair := range pairsChannel{
		swap(&bsChart.Data[pair[0]], &bsChart.Data[pair[1]])
		m.Lock()
		ui.Render(&bsChart)
		m.Unlock()
	}

	playSound()

	//End
	bsChart.Title = "BubbleSort-Finalizado-" +
		"Tiempo:"+strconv.FormatInt(bsTime.Milliseconds(),10)+"ms-" +
		"Swaps:"+strconv.Itoa(bsSwaps)+"-" +
		"Comparaciones:"+strconv.Itoa(bsComparisons)+"-"+
		"Iteraciones:"+strconv.Itoa(bsIterations)
	m.Lock()
	ui.Render(&bsChart)
	m.Unlock()
}

func ssChartDrawer(slice []float64){
	// / / ssChart.Data = slice
	ssChart.Data = make([]float64, len(slice))
	copy(ssChart.Data, slice)

	//Copy used in SelectionSort
	copyArr := make([]float64, len(slice))
	copy(copyArr, ssChart.Data)

	//Channel
	pairsChannel := make(chan []int, 1000)
	go selectionSort(&copyArr, pairsChannel)

	//Update Changes in pairs
	for pair := range pairsChannel{
		swap(&ssChart.Data[pair[0]], &ssChart.Data[pair[1]])
		m.Lock()
		ui.Render(&ssChart)
		m.Unlock()
	}

	playSound()

	//End
	ssChart.Title = "SelectionSort-Finalizado-" +
		"Tiempo:"+strconv.FormatInt(ssTime.Milliseconds(),10)+"ms-" +
		"Swaps:"+strconv.Itoa(ssSwaps)+"-" +
		"Comparaciones:"+strconv.Itoa(ssComparisons)+"-"+
		"Iteraciones:"+strconv.Itoa(ssIterations)
	m.Lock()
	ui.Render(&ssChart)
	m.Unlock()
}

func isChartDrawer(slice []float64){
	// / / isChart.Data = slice
	isChart.Data = make([]float64, len(slice))
	copy(isChart.Data, slice)

	//Copy used in SelectionSort
	copyArr := make([]float64, len(slice))
	copy(copyArr, isChart.Data)

	//Channel
	oneWayChannel := make(chan []int, 1000)
	go insertionSort(&copyArr, oneWayChannel)

	//Update Changes in pairs
	for oneWay := range oneWayChannel{
		isChart.Data[oneWay[0]] = float64(oneWay[1])//isChart.Data[oneWay[1]]
		m.Lock()
		ui.Render(&isChart)
		m.Unlock()
	}

	playSound()

	//End
	isChart.Title = "InsertionSort-Finalizado-" +
		"Tiempo:"+strconv.FormatInt(isTime.Milliseconds(),10)+"ms-" +
		"Swaps:"+strconv.Itoa(isSwaps)+"-" +
		"Comparaciones:"+strconv.Itoa(isComparisons)+"-"+
		"Iteraciones:"+strconv.Itoa(isIterations)
	m.Lock()
	ui.Render(&isChart)
	m.Unlock()
}

func qsChartDrawer(slice []float64){
	// / / qsChart.Data = slice
	qsChart.Data = make([]float64, len(slice))
	copy(qsChart.Data, slice)

	//Copy used in SelectionSort
	copyArr := make([]float64, len(slice))
	copy(copyArr, qsChart.Data)

	//Channel
	pairsChannel := make(chan []int, 1000)
	go quickSort(&copyArr, pairsChannel)

	//Update Changes in pairs
	for pair := range pairsChannel{
		swap(&qsChart.Data[pair[0]], &qsChart.Data[pair[1]])
		m.Lock()
		ui.Render(&qsChart)
		m.Unlock()
	}

	playSound()

	//End
	qsChart.Title = "QuickSort-Finalizado-" +
		"Tiempo:"+strconv.FormatInt(qsTime.Milliseconds(),10)+"ms-" +
		"Swaps:"+strconv.Itoa(qsSwaps)+"-" +
		"Comparaciones:"+strconv.Itoa(qsComparisons)+"-"+
		"Iteraciones:"+strconv.Itoa(qsIterations)
	m.Lock()
	ui.Render(&qsChart)
	m.Unlock()
}

func hsChartDrawer(slice []float64){
	// / / hsChart.Data = slice
	hsChart.Data = make([]float64, len(slice))
	copy(hsChart.Data, slice)

	//Copy used in SelectionSort
	copyArr := make([]float64, len(slice))
	copy(copyArr, hsChart.Data)

	//Channel
	pairsChannel := make(chan []int, 1000)
	go heapSort(&copyArr, pairsChannel)

	//Update Changes in pairs
	for pair := range pairsChannel{
		swap(&hsChart.Data[pair[0]], &hsChart.Data[pair[1]])
		m.Lock()
		ui.Render(&hsChart)
		m.Unlock()
	}

	playSound()

	//End
	hsChart.Title = "HeapSort-Finalizado-" +
		"Tiempo:"+strconv.FormatInt(hsTime.Milliseconds(),10)+"ms-" +
		"Swaps:"+strconv.Itoa(hsSwaps)+"-" +
		"Comparaciones:"+strconv.Itoa(hsComparisons)+"-"+
		"Iteraciones:"+strconv.Itoa(hsIterations)
	m.Lock()
	ui.Render(&hsChart)
	m.Unlock()
}

// / / / / / / Extra \ \ \ \ \ \

func swap (a *float64, b *float64){
	temp := *a
	*a = *b
	*b = temp
}

func playSound(){
	f, _ := os.Open("w.mpeg")
	streamer, format, _ := mp3.Decode(f)
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(streamer)
}

func randomSlice(size int, seed int, k int, period int) []float64 {
	
	// Validating "size"
	
	if size < 10 || size > 100 {
		fmt.Println("El valor size es incorrecto")
		return nil
	}
	
	// Validating Seed
	if seed < 11 || seed > 101 {
		fmt.Println("El valor de la semilla es incorrecto")
		return nil
	}

	for i := 2; i < seed; i++ {	// Prime Number
		if seed % i == 0 {
			fmt.Println("El valor de la semilla es incorrecto")
			return nil
		}
	}

	// Validating "k"
	if k < 0 {
		fmt.Println("El valor k es incorrecto")
		return nil
	}

	// Validating "period"
	if period < 2048 {
		fmt.Println("El valor m es incorrecto")
		return nil
	}
	
	var slice = make([]float64, size)
	multiplier := 8*k + 3	// 8k + 5 can also be used

	/*
	first := (multiplier * seed) % period
	first = first % 30
	*/

	for i := 0; i < size; i++ {	// Generating the Array
		num := (multiplier * seed) % period	// Main Algorithm, X = (multiplier * [seed or previous number]) % period
		num = num % 30	// Changed to 0..29

		/*
		if num == first && i != 0{	// We can stop the loop to avoid repeating the pattern
			slice = slice[:i]
			break
		}
		*/
		
		slice[i] = float64(num)
		seed = num	// Seed is now the previous number
	}

	fmt.Println("Resultado: ", slice)

	return slice
}