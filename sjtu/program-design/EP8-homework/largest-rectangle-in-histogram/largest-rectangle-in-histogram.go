package main

import (
	"container/list"
	"fmt"
)

// For explanation, please see https://www.youtube.com/watch?v=RVIh0snn4Qc
func main() {
	histogram := [8]int{1, 2, 3, 4, 5, 3, 3, 6}

	// use List to simulate stack
	var stack = list.New()
	// init
	stack.PushBack(0)
	var i = 1
	var max int

	for stack.Len() != 0 {
		if i == len(histogram) || histogram[i] < histogram[stack.Front().Value.(int)] {
			// pop
			index := stack.Remove(stack.Front()).(int)
			length := i - index
			if length*histogram[index] > max {
				max = length * histogram[index]
			}
		} else {
			// push
			stack.PushFront(i)
			i++
		}
	}

	fmt.Println(max)

}
