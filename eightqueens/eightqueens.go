package main

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"
)

type ChessBoard [8][8]bool
type SolType [8]string
type ChessSol []SolType

var count int
var lock sync.Mutex

var solutions ChessSol

func (chessSol ChessSol) Len() int {
	return len(chessSol)
}
func (chessSol ChessSol) Less(i, j int) bool {

	var left, right SolType = chessSol[i], chessSol[j]

	for c := 0; c < 8; c++ {
		if left[c] == right[c] {
			continue
		} else {
			return left[c] < right[c]
		}
	}
	return true
}
func (chessSol ChessSol) Swap(i, j int) {
	chessSol[i], chessSol[j] = chessSol[j], chessSol[i]
}

func (board ChessBoard) display() (ret [8]string) {
	chessLetter := [8]string{"a", "b", "c", "d", "e", "f", "g", "h"}

	for i := 0; i < 8; i++ {
		fmt.Println()
		posstr := ""
		for j := 0; j < 8; j++ {
			if board[i][j] {
				posstr = chessLetter[i] + strconv.Itoa(j+1)
				fmt.Print("q")
			} else {
				fmt.Print(".")
			}
			fmt.Print(" ")
		}
		ret[i] = posstr
	}

	return
}

func (board *ChessBoard) set(row, col int, val bool) bool {
	if row < 0 || row > 7 || col < 0 || col > 7 {
		return false
	}
	board[row][col] = val
	return true
}

func (board ChessBoard) checkAnyCross(row, col int) bool {

	check := func(name string, rowFunc, colFunc func(*int)) bool {
		r, c := row, col
		rowFunc(&r)
		colFunc(&c)

		for r >= 0 && r < 8 && c >= 0 && c < 8 {
			if board[r][c] {
				return true
			}
			rowFunc(&r)
			colFunc(&c)
		}
		return false
	}

	decrement := func(v *int) { *v-- }
	increment := func(v *int) { *v++ }
	nochange := func(v *int) {}

	return check("NE", decrement, decrement) ||
		check("N", decrement, nochange) ||
		check("NW", decrement, increment)
}

// Return current column position for given row
func (board ChessBoard) curColPos(row int) int {
	for c := 0; c < 8; c++ {
		if board[row][c] {
			return c
		}
	}
	return -1
}

// Process name and start & end column
func (board *ChessBoard) process(name string, scol, ecol int) {
	for i, c := 0, scol; i < 8 && i >= 0; i++ {
		emptyrow := true
		for j := c; j < 8; j++ {
			if !board.checkAnyCross(i, j) {
				board.set(i, j, true)
				emptyrow = false
				break
			}
		}

		if !emptyrow && i == 7 {
			firstRowCol := board.curColPos(0)
			if firstRowCol == -1 || firstRowCol == ecol {
				break
			}
			i++
			emptyrow = true
			lock.Lock()
			count++
			displayCount := count
			lock.Unlock()
			fmt.Printf("\n\n%s_%d:\n", name, displayCount)
			sol := board.display()
			fmt.Print(sol)
			solutions = append(solutions, sol)
		}

		if c = 0; emptyrow {
			for k := i - 1; k >= 0; k-- {
				c = board.curColPos(k)
				if board.set(k, c, false) {
					i = k - 1
					c++
					break
				} else {
					c = 0
				}
			}
		}
	}
}

func main() {
	startTime := time.Now()

	result := make(chan bool, 8)

	for i := 0; i < 8; i++ {
		go func(proc int) {
			board := ChessBoard{}
			board.process("G"+strconv.Itoa(proc), proc, proc+1)
			result <- true
		}(i)
	}

	for doit := 0; doit < 8; {
		<-result
		doit++
	}

	var solInterface sort.Interface = solutions
	sort.Sort(solInterface)

	fmt.Println("\n\nSolutions: ")
	for _, v := range solutions {
		fmt.Println(v)
	}

	fmt.Printf("\n\nTotal time taken for solution: %s\n", time.Since(startTime))
}
