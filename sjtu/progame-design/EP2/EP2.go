package main

import (
	"fmt"
)

func isConvexPolygons(input [][]int) bool {

	node := 0
	flag := []bool{}
	for i := 0; i < len(input); i++ {
		flag = append(flag, false)
	}
	//choose the first node
	flag[node] = true

	for !everyNodeIsLinked(flag) {
		for j := 0; j < len(input); j++ {
			if j == len(input)-1 {
				return false
			}
			if flag[j] {
				continue
			}
			if isOnTheSameSide(input[node], input[j], input) {
				flag[j] = true
				node = j
				break
			}
		}
	}

	return true

}

func everyNodeIsLinked(flag []bool) bool {
	for i := range flag {
		if flag[i] == false {
			return false
		}
	}
	return true
}

func isOnTheSameSide(a, b []int, input [][]int) bool {
	buffer := []int{}
	for i := 0; i < len(input); i++ {
		buffer = append(buffer, input[i][0]*(b[1]-a[1])+input[i][1]*(a[0]-b[0])-a[0]*b[1]+a[1]*b[0])
	}
	//if every node is on the same side, they all have the same symbol
	for i := range buffer {
		for j := i + 1; j < len(buffer); j++ {
			if buffer[i]*buffer[j] < 0 {
				return false
			}
		}
	}
	return true
}

//calculate segment's y=kx+b, return k and b
func calLineFuc(pointA []int, pointB []int) (float64, float64) {
	k := float64((pointB[1] - pointA[1])) / float64((pointB[0] - pointA[0]))
	b := float64(pointB[1]) - k*float64(pointB[0])
	return k, b
}

func main() {
	input := [][]int{{0, 0}, {3, -1}, {4, 1}, {4, 2}, {1, 1}, {3, 2}, {3, 4}, {1, 2}, {0, 2}}
	k := isConvexPolygons(input)
	if k {
		fmt.Printf("Your input is Convex Polygon")
	} else {
		fmt.Printf("Your input is Convax Polygon")
	}

}
