package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/lxn/win"
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

	barNumber := width / BAR_WIDTH - 1
	fmt.Print("Indique la cantidad de numeros(Se recomienda " + strconv.Itoa(barNumber) +" maximo para una visualizacion correcta): ")

	var size int
	fmt.Scanln(&size)
	if err := ui.Init(); err != nil{
		log.Fatalf("failed to initialize termui: %v", err)
	}

	arregloBase := temporalRANDOM(size) // CAMBIAR TEMPORAL
	initBSChart(arregloBase)
	initSSChart(arregloBase)
	ui.Render(&bsChart)
	ui.Render(&ssChart)
	go bsChartDrawer(arregloBase)
	ssChartDrawer(arregloBase)
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
	arr2 := *arr
	len := len(arr2)

	//Move through all elements
	for i:= 0; i < len; i++ {
		for j := 0; j < len-i-1; j++ {
			// Move from 0 to len-i-1 and swap if element is greater than the next one
			if arr2[j] > arr2[j+1] {
				arr2[j], arr2[j+1] = arr2[j+1], arr2[j]
				pair <- []int{j, j+1}
			}
		}
	}

	//*arr = arr2 // Assign changes to original array
	close(pair)
}

// Selection (changes indexes)

func selectionSort(arr *[]int) {
	arr2 := *arr
	len := len(arr2)

	for currentIndex := 0; currentIndex < len-1; currentIndex++ { // Done to all the indexes in the array
		indexMin := currentIndex

		for i := currentIndex + 1; i < len; i++ { // Get the index of the smallest value from the numbers to the right
			if arr2[i] < arr2[indexMin] {
				indexMin = i
			}
		}

		//Swap numbers
		arr2[currentIndex], arr2[indexMin] = arr2[indexMin], arr2[currentIndex]
	}

	*arr = arr2 // Assign changes to original array
}

// Insertion >> New <<

func insertionSort(arr *[] int) {
	arr2 := *arr
	len := len(arr2)

	for i := 1; i < len; i++ {
		key := arr2[i]
		j := i-1

		//Move greater elements of arr[0 .. i-1] to position ahead of current
		for ; j >= 0 && key < arr2[j]; j--{
			arr2[j+1] = arr2[j]
		}

		arr2[j+1] = key
	}

	*arr = arr2 // Assign changes to original array
}

// Quicksort [iterative for drawing]: https://www.geeksforgeeks.org/iterative-quick-sort/

func partition(arr *[]int, low int, high int) int { //
	arr2 := *arr
	pivot := arr2[high]

	i := low - 1

	for j := low; j < high; j++ {
		if arr2[j] <= pivot {
			i++
			arr2[i], arr2[j] = arr2[j], arr2[i] //Gets the lesser values to the left of the pivot
		}
	}

	//Swap pivot with the next element to i
	arr2[i+1], arr2[high] = arr2[high], arr2[i+1]

	*arr = arr2 // Assign changes to original array

	return i + 1 //new pivot
}

func quickSort(arr *[]int){
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

		pivot := partition(arr, low, high) //pivot at correct position

		if pivot-1 > low { //If elements on left push left side to stack
			top++
			stack[top] = low
			top++
			stack[top] = pivot-1
		}

		if pivot+1 < high{
			top++
			stack[top] = pivot+1
			top++
			stack[top] = high
		}
	}
}

// Heapsort >> New << [iteraive just in case]: https://www.geeksforgeeks.org/iterative-heap-sort/

func buildMaxHeap(arr *[]int, n int){
	arr2 := *arr
	for i := 1; i < n; i++{
		if arr2[i] > arr2[(i-1)/2]{ // Child bigger than parent
			j := i

			for arr2[j] > arr2[(j-1)/2]{ //Swap until parent is smaller than child
				arr2[j], arr2[(j-1)/2] = arr2[(j-1)/2], arr2[j]
				j = (j-1)/2
			}
		}
	}

	*arr = arr2
}

func heapSort(arr *[]int){
	n := len(*arr)

	buildMaxHeap(arr, n)
	arr2 := *arr

	for i := n-1; i > 0; i--{
		arr2[0], arr2[i] = arr2[i], arr2[0] //swap first with last
		j, index := 0, 0

		for { 
			index = 2 * j + 1

			if index < (i - 1) && arr2[index] < arr2[index + 1]{
				index++
			}; if index < i && arr2[j] < arr2[index]{
				arr2[j], arr2[index] = arr2[index], arr2[j]
			}; j = index

			if index >= i{
				break
			}
			
		}
	}

	*arr = arr2
}


// / / / / / / Graphic \ \ \ \ \ \

// / / / / Initialize 

func initBSChart(arr []float64){
	bsChart = *widgets.NewBarChart()
	bsChart.Data = arr
	bsChart.Title = "BubbleSort"
	bsChart.SetRect(0, 0, width/2 - 2, height-2)
	bsChart.BarWidth = BAR_WIDTH
	bsChart.BarGap = 0
	bsChart.BarColors = []ui.Color{ui.ColorRed}
	bsChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorRed)} // Can't be seen

	// Indexes
	//bsChart.Labels = generateLabels(arr)
	//bsChart.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorWhite)}
}

func initSSChart(arr []float64){
	ssChart = *widgets.NewBarChart()
	ssChart.Data = arr
	ssChart.Title = "SelectionSort"
	ssChart.SetRect(width/2, 0, width - 4, height-2)
	ssChart.BarWidth = BAR_WIDTH
	ssChart.BarGap = 0
	ssChart.BarColors = []ui.Color{ui.ColorBlue}
	ssChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlue)} // Can't be seen
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

	//End
	bsChart.Title = "BubbleSort-Finalizado-" +
		"Tiempo:"+strconv.FormatInt(bsTime.Milliseconds(),10)+"ms-" +
		"Swaps:"+strconv.Itoa(bsSwaps)+"-" +
		"Comparaciones:"+strconv.Itoa(bsComparisons)+"-"+
		"Iteraciones:"+strconv.Itoa(bsIterations)
	m.Lock()
	ui.Render(&bsChart)
	m.Unlock()
	
	//fmt.Println(slice)
	//fmt.Println(bsChart.Data)
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
	go bubbleSort(&copyArr, pairsChannel)

	//Update Changes in pairs
	for pair := range pairsChannel{
		swap(&ssChart.Data[pair[0]], &ssChart.Data[pair[1]])
		m.Lock()
		ui.Render(&ssChart)
		m.Unlock()
	}

	//End
	ssChart.Title = "BubbleSort-Finalizado-" +
		"Tiempo:"+strconv.FormatInt(ssTime.Milliseconds(),10)+"ms-" +
		"Swaps:"+strconv.Itoa(ssSwaps)+"-" +
		"Comparaciones:"+strconv.Itoa(ssComparisons)+"-"+
		"Iteraciones:"+strconv.Itoa(ssIterations)
	m.Lock()
	ui.Render(&ssChart)
	m.Unlock()
	
	//fmt.Println(slice)
	//fmt.Println(ssChart.Data)
}

// / / / / / / Extra \ \ \ \ \ \

func swap (a *float64, b *float64){
	temp := *a
	*a = *b
	*b = temp
}

func generateLabels(arr []float64) []string {
	var labels = make([]string, len(arr))
	for i := range arr {
		labels[i] = strconv.Itoa(i)
	}; return labels
}