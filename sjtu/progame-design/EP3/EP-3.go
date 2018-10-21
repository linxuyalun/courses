package main

import (
	"fmt"
)

func main() {
	grid := 3
	square := make([][]int, grid)
	for row := 0; row < grid; row++ {
		square[row] = make([]int, grid)
	}
	// i is row number
	i := 0
	// j is column number
	j := (grid - 1) / 2
	fmt.Println(square[i][j])
	// first number in square
	square[i][j] = 1
	// number is all the numbers in the squere
	number := 2
	for number != grid*grid+1 {
		tempI := (i + grid - 1) % grid
		tempJ := (j + 1) % grid
		if square[tempI][tempJ] == 0 {
			i = tempI
			j = tempJ
			square[i][j] = number
			number++
			continue
		}
		for square[i][j] != 0 {
			i = (i + 1) % grid
		}
		square[i][j] = number
		number++
	}
	fmt.Println(square)

	// solve the magic square with some blank cells
	// Actually, there's some laws in 3 order magic square: each number can be calculated by other two specific numbers
	var a [9]int
	a[1] = 6
	a[2] = 10
	a[3] = 4
	fillAllBlanks := false
	for !fillAllBlanks {

		if a[0] == 0 {
			if a[5] != 0 && a[7] != 0 {
				a[0] = (a[5] + a[7]) / 2
			}
			if a[4] != 0 && a[8] != 0 {
				a[0] = 2*a[4] - a[8]
			}
		}

		if a[1] == 0 {
			if a[4] != 0 && a[7] != 0 {
				a[1] = 2*a[4] - a[7]
			}
			if a[3] != 0 && a[8] != 0 {
				a[1] = 2*a[8] - a[3]
			}
			if a[5] != 0 && a[6] != 0 {
				a[1] = 2*a[6] - a[5]
			}
		}

		if a[2] == 0 {
			if a[3] != 0 && a[7] != 0 {
				a[2] = (a[3] + a[7]) / 2
			}
			if a[4] != 0 && a[6] != 0 {
				a[2] = 2*a[4] - a[6]
			}
		}

		if a[3] == 0 {
			if a[4] != 0 && a[5] != 0 {
				a[3] = 2*a[4] - a[5]
			}
			if a[1] != 0 && a[8] != 0 {
				a[3] = 2*a[8] - a[1]
			}
			if a[7] != 0 && a[2] != 0 {
				a[3] = 2*a[2] - a[7]
			}
		}

		if a[4] == 0 {
			if a[0] != 0 && a[8] != 0 {
				a[4] = (a[0] + a[8]) / 2
			}
			if a[1] != 0 && a[7] != 0 {
				a[4] = (a[1] + a[7]) / 2
			}
			if a[1] != 0 && a[7] != 0 {
				a[4] = (a[1] + a[7]) / 2
			}
			if a[2] != 0 && a[6] != 0 {
				a[4] = (a[2] + a[6]) / 2
			}
			if a[3] != 0 && a[6] != 0 {
				a[4] = (a[3] + a[6]) / 2
			}
			if a[3] != 0 && a[5] != 0 {
				a[4] = (a[3] + a[5]) / 2
			}
		}

		if a[5] == 0 {
			if a[3] != 0 && a[4] != 0 {
				a[5] = 2*a[4] - a[3]
			}
			if a[0] != 0 && a[7] != 0 {
				a[5] = 2*a[0] - a[7]
			}
			if a[1] != 0 && a[6] != 0 {
				a[5] = 2*a[6] - a[1]
			}
		}

		if a[6] == 0 {
			if a[1] != 0 && a[5] != 0 {
				a[6] = (a[1] + a[5]) / 2
			}
			if a[4] != 0 && a[2] != 0 {
				a[6] = 2*a[4] - a[2]
			}
		}

		if a[7] == 0 {
			if a[1] != 0 && a[4] != 0 {
				a[7] = 2*a[4] - a[1]
			}
			if a[3] != 0 && a[2] != 0 {
				a[7] = 2*a[2] - a[3]
			}
			if a[0] != 0 && a[5] != 0 {
				a[7] = 2*a[0] - a[5]
			}
		}

		if a[8] == 0 {
			if a[1] != 0 && a[3] != 0 {
				a[8] = (a[1] + a[3]) / 2
			}
			if a[4] != 0 && a[0] != 0 {
				a[0] = 2*a[4] - a[0]
			}
		}

		for i := range a {
			if a[i] == 0 {
				fillAllBlanks = false
				break
			}
			fillAllBlanks = true
		}
	}

	fmt.Println(a)
}
