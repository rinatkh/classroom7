package main

import (
	"fmt"

	"github.com/rinatkh/classroom7/internal/simplesort"
)

func main() {
	fmt.Println(simplesort.Example())
	input := []int{5, 1, 4, 2}
	fmt.Printf("Bubble, состояния: %v\n", simplesort.BubbleTrace(append([]int(nil), input...)))
	fmt.Printf("Selection, состояния: %v\n", simplesort.SelectionTrace(append([]int(nil), input...)))
	fmt.Printf("Insertion, состояния: %v\n", simplesort.InsertionTrace(append([]int(nil), input...)))
}
