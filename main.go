package main

import (
	"fmt"
)

func main(){
	//Pedir n (intervalo 10 a 100)
	//Semilla: entrada de datos por teclado o segundos y decimas del reloj del sistema -> convertido a num entre 0 y 599
	//Arreglo base

	arr := []int{-12,32,99,67,-83,123, 1, -8923, 23}
	heapSort(&arr)
	fmt.Println(arr)

}

/*
Función generadora de arreglos de tamaño n que contengan números pseudo-aleatorios 
obtenidos mediante el método de congruencia lineal multiplicativa, a partir de una semilla dada. La semilla 
deberá ser un número primo entre 11 y 101. Los valores generados deben ser convertidos al intervalo 0 .. 29. 
El período debe ser ≥ 2048. n puede ser cualquier número en el intervalo 10 .. 100
*/

func RandArray(n int, seed int, k int, m int) []int {
	
	// Validating "n"
	
	/*
	if n < 10 || n > 100 { \ \ \ Cambio
		fmt.Println("El valor n es incorrecto")
		return nil
	}
	*/
	
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

	// Validating "m"
	if m < 2048 {
		fmt.Println("El valor m es incorrecto")
		return nil
	}
	
	arr := make([]int, n)
	a := 8*k + 3	// 8k + 5 can also be used
	//first := seed

	for i := 0; i < n; i++ {	// Generating the Array
		num := (a * seed) % m	// Main Algorithm, X = (a * [seed or previous number]) % m
		num = num % 30	// Changed to 0..30 \ \ \ Cambio

		/*
		if num == first && i != 0{	// We can stop the loop to avoid repeating the pattern
			arr = arr[:i]
			break
		}
		*/

		arr[i] = num
		seed = num	// Seed is now the previous number
	}

	fmt.Println("Resultado: ", arr)

	return arr
}

// / / / / / / Sorting \ \ \ \ \ \

// Bubblesort >> New <<

func bubbleSort(arr *[]int) {
	arr2 := *arr
	len := len(arr2)

	//Move through all elements
	for i:= 0; i < len; i++ {
		for j := 0; j < len-i-1; j++ {
			// Move from 0 to len-i-1 and swap if element is greater than the next one
			if arr2[j] > arr2[j+1] {
				arr2[j], arr2[j+1] = arr2[j+1], arr2[j]
			}
		}
	}

	*arr = arr2 // Assign changes to original array
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

// Quicksort [iterative for drawing]

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


func quickSort(arr *[]int){ //based on: https://www.geeksforgeeks.org/iterative-quick-sort/
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

