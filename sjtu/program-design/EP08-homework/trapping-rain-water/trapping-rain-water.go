package main

import (
	"container/list"
	"fmt"
)

func main() {

	var b = [8]int{3, 4, 5, 3, 4, 2, 6, 0}
	// ans is the sum of water, i is the index of array
	var ans, i int

	// use List to simulate stack
	var stack = list.New()
	// init
	stack.PushFront(0)
	i++

	for i < len(b) {
		if stack.Len() != 0 && b[i] > b[stack.Front().Value.(int)] {
			// pop
			bottomIndex := stack.Remove(stack.Front()).(int)
			var height, length int
			if stack.Len() != 0 {
				if b[stack.Front().Value.(int)] > b[i] {
					height = b[i] - b[bottomIndex]
				} else {
					height = b[stack.Front().Value.(int)] - b[bottomIndex]
				}
			} else {
				continue
			}
			length = i - stack.Front().Value.(int) - 1
			ans += length * height
		} else {
			// push
			stack.PushFront(i)
			i++

		}
	}

	fmt.Println(ans)

}
