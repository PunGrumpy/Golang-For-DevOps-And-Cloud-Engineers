package main

import "fmt"

func main() {
	var arr1 [7]int = [7]int{1, 2, 3, 4, 5, 6, 7}
	fmt.Printf("arr1: %v\n", arr1)
	fmt.Printf("%d %d\n", len(arr1), cap(arr1))
	var arr2 []int = arr1[1:3]
	fmt.Printf("arr2: %v\n", arr2)
	fmt.Printf("%d %d\n", len(arr2), cap(arr2))
	arr2 = arr2[0 : len(arr2)+2]
	fmt.Printf("arr2: %v\n", arr2)
	fmt.Printf("%d %d\n", len(arr2), cap(arr2))
	for k := range arr2 {
		arr2[k] += 1
	}
	fmt.Printf("arr2: %v\n", arr2)
	fmt.Printf("%d %d\n", len(arr2), cap(arr2))
	fmt.Printf("arr1: %v\n", arr1)

	var arr3 []int = []int{1, 2, 3}
	fmt.Printf("arr3: %v\n", arr3)
	fmt.Printf("%d %d\n", len(arr3), cap(arr3))
	arr3 = append(arr3, 4)
	fmt.Printf("arr3: %v\n", arr3)
	fmt.Printf("%d %d\n", len(arr3), cap(arr3))
	arr3 = append(arr3, 5)
	fmt.Printf("arr3: %v\n", arr3)
	fmt.Printf("%d %d\n", len(arr3), cap(arr3))

	arr4 := make([]int, 3)
	fmt.Printf("arr4: %v\n", arr4)
	fmt.Printf("%d %d\n", len(arr4), cap(arr4))
}
