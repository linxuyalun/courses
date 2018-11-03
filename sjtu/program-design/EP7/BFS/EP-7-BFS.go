package main

import (
	"container/list"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// color prints the blackground color
func printColor(color string) {
	var colorlist = map[string]int{
		"Black":  40,
		"Red":    41,
		"Green":  42,
		"Yellow": 43,
		"Blue":   44,
		"Purple": 45,
		"Azure":  46,
		"White":  47,
	}
	fmt.Printf("%c[%d;%d;%dm  %c[0m", 0x1B, 0, colorlist[color], 30, 0x1B)
}

// generate a random maze
func generateMaze(length int) [][]int {
	var a [][]int
	var b []int
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		a = append(a, b)
		for j := 0; j < length; j++ {
			num := rand.Float32()
			if num > 0.7 {
				a[i] = append(a[i], 1)
			} else {
				a[i] = append(a[i], 0)
			}
		}
	}
	// generater bawn
	for i := 0; i < length; i++ {
		a[0][i] = 1
		a[i][0] = 1
		a[length-1][i] = 1
		a[i][length-1] = 1
	}
	// display the maze
	for _, v := range a {
		for _, value := range v {
			if value == 8 {
				printColor("Red")
			} else if value == 1 {
				printColor("Green")
			} else {
				printColor("White")
			}
		}
		fmt.Printf("\n")
	}

	return a
}

// Point records the coordinate of the current point and its parent point's location in queue
// X,Y is the point coordinate
// Parent is the dequeued point's location in assistant array
type Point struct {
	X      int
	Y      int
	Parent int
}

// mark displays the maze after find the exit
func mark(array [][]int, assistant []Point, point Point) {
	if point.Parent == -1 {
		array[point.X][point.Y] = 8
		for _, v := range array {
			for _, value := range v {
				if value == 8 {
					printColor("Red")
				} else if value == 1 {
					printColor("Green")
				} else {
					printColor("White")
				}
			}
			fmt.Printf("\n")
		}
		os.Exit(0)
	}
	array[point.X][point.Y] = 8
	mark(array, assistant, assistant[point.Parent])

}

// 0 represents the point can go to; 1 represents the wall; 2 represents the point have gone to
func bfs(array [][]int, length, depx, depy, destx, desty int) {
	// use List to simulate the queue operation
	var queue = list.New()
	// root point
	entrance := Point{depx, depy, -1}
	array[depx][depy] = 2
	queue.PushBack(entrance)

	// traversal 8 diretion: →, ↘, ↓, ↙, ←, ↖, ↑
	dx := [8]int{1, 1, 0, -1, -1, -1, 0, 1}
	dy := [8]int{0, 1, 1, 1, 0, -1, -1, -1}
	// assistant records the dequeued points
	assistant := []Point{}
	for queue.Len() != 0 {
		// dequeue
		temp := queue.Remove(queue.Front())
		// type assertion
		curr := temp.(Point)
		assistant = append(assistant, curr)
		pointLocation := len(assistant) - 1

		// find the exit
		if curr.X == destx && curr.Y == desty {
			mark(array, assistant, curr)
		}

		for i := 0; i < 8; i++ {
			next := array[curr.X+dx[i]][curr.Y+dy[i]]
			if next != 1 && next != 2 {
				array[curr.X+dx[i]][curr.Y+dy[i]] = 2
				// enqueue
				nextPoint := Point{curr.X + dx[i], curr.Y + dy[i], pointLocation}
				queue.PushBack(nextPoint)
			}
		}
	}
}

func main() {

	fmt.Println("Input the length:")
	var length int
	fmt.Scan(&length)
	array := generateMaze(length)
	fmt.Printf("\n")
	fmt.Println("Input the ENTRANCE and EXIT:")
	var a, b, c, d int
	fmt.Scan(&a, &b, &c, &d)
	bfs(array, length, a, b, c, d)

}
