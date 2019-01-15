package main

import (
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
			if num > 0.6 {
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

// 0 represents the point can go to; 1 represents the wall; 2 represents the point have gone to; 8 represents the path
func dfs(array [][]int, length, currx, curry, destx, desty int) {

	array[currx][curry] = 8

	// find the way
	if currx == destx && curry == desty {
		for _, v := range array {
			for _, value := range v {
				if value == 8 {
					printColor("Red")
				} else if value == 1 {
					printColor("Green")
				} else if value == 2 {
					printColor("Azure")
				} else {
					printColor("White")
				}
			}
			fmt.Printf("\n")
		}
		os.Exit(0)
	}

	// traversal 8 diretion: →, ↘, ↓, ↙, ←, ↖, ↑
	dx := [8]int{1, 1, 0, -1, -1, -1, 0, 1}
	dy := [8]int{0, 1, 1, 1, 0, -1, -1, -1}
	for i := 0; i < 8; i++ {
		if array[currx+dx[i]][curry+dy[i]] != 8 && array[currx+dx[i]][curry+dy[i]] != 1 && array[currx+dx[i]][curry+dy[i]] != 2 {
			dfs(array, length, currx+dx[i], curry+dy[i], destx, desty)
		}
	}

	// backtracking
	array[currx][curry] = 2
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
	dfs(array, length, a, b, c, d)
}
