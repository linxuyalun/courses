package main

import (
	"container/list"
	"fmt"
)

func main() {

	var b = [5]int{3, 1, 0, 1, 5}
	// ans is the sum of water, i is the index of array
	var ans, i int

	// use List to simulate stack
	var stack = list.New()
	// init
	stack.PushFront(0)
	i++

	for i < len(b) {
		if stack.Len() != 0 && b[i] > stack.Front().Value.(int) {
			// pop
			bottomIndex := stack.Remove(stack.Front()).(int)
			length := i - bottomIndex
			var height int
			if stack.Len() != 0 {
				if b[stack.Front().Value.(int)] > b[i] {
					height = b[i]
				} else {
					height = b[stack.Front().Value.(int)]
				}
			} else {
				continue
			}
			ans += length * height
		} else {
			// push
			stack.PushFront(b[i])
			i++

		}
	}

	fmt.Println(ans)

}
