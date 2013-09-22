package main

import (
	"math/rand"
	"os"
	"io/ioutil"
	"time"
)

// SETUP
const WAIT = 100

// Board
var board []rune
var tempBoard = make([]rune, len(board))

// Player
type Player struct {
	Symbol	rune
	X 		int
	Y		int
}
var p1 = Player{2015, 3, 3}	

// Finish Line
type FinishLine struct {
	X 		int
	Y		int
}
var fL []FinishLine
var fR []FinishLine

// Counter
var circles int = 0
var outFlg = true

func main() {
	var trackFileName string
	if  len(os.Args) > 1 {
		trackFileName = os.Args[1]
	} else {
		trackFileName = "t1.tr"
	}	
	var boardByte, err = ioutil.ReadFile(trackFileName)
	if (err != nil) {
		panic ("Could not open file with track")
	}
	var boardString = string(boardByte)
	board = []rune(boardString)
	tempBoard = make([]rune, len(board))
	println(getBoard())
	getFinishLines()
	
	for !checkWin() {
		var move = getMove()
		if (isMoveOk(move)) {
			executeMove(move)
		}
		checkCircle()
		wait(WAIT)
		brd:=getBoard()
		clearScr()
		println(brd, circles)
		clearScr()
	}
}
 
func getFinishLines() {
	for key, value := range board {
		if value == 'L' {
			var x, y = getPosXY(key)
			fL = append(fL,FinishLine{x,y})
		}
		if value == 'R' {
			var x, y = getPosXY(key)
			fR = append(fR,FinishLine{x,y})
		}
	}
}

var conHeight int = 11 //env.Getenv("COLUMNS")

var flFlag = false
var frFlag = false

func checkCircle() {
	if doesArrayContainPosition(getPosIntPl(p1), fL) {
		if (frFlag == true) {
			circles--
			flFlag = true
			frFlag = false
		} else {
			flFlag = true
		}
    } else if doesArrayContainPosition(getPosIntPl(p1), fR) {
    	if (flFlag == true) {
			circles++
			frFlag = true
			flFlag = false
		} else {
			frFlag = true
		}
    } else {
		flFlag = false
		frFlag = false
	}
}

func doesArrayContainPosition(position int, fll []FinishLine) bool {
	for _, fl := range fll {
		if position == getPosIntFl(fl) {
			return true
		}
	}
	return false
}

func printCircles() {
	println(circles)
	clearScr()
}

func clearScr() {
	for i:=0; i < conHeight; i++ {
		println()
	}
}

func toString(pl Player) string {
	return string(pl.Symbol) +" "+ string(pl.X) +" "+ string(pl.Y)
}

func getBoard() string {
	copy(tempBoard[:], board[0:len(board)])
	drawPlayerOnBoard(p1, tempBoard)
	return string(tempBoard)
}

func drawPlayerOnBoard(pl Player, board []rune) {
	var posInt = getPosInt(pl.X, pl.Y) 
	board[posInt] = pl.Symbol
}

func getPosXY(pos int) (int, int) {
	var x = 0
	var y = 0
	for i := 0; i<=pos; i++ {
		x++
		if board[i] == '\n' {
			y++
			x = 0
		}
	}
	return x, y
}

func getPosInt(x int, y int) int {
	for i := 0; i<len(board); i++ {
		if board[i] == '\n' {
			y--
		}
		if y == 0 {
			x--
		}
		if y == 0 && x == 0 {
			return i
		}
	}	
	panic("getPosInt: wrong coordinates")
}

func getPosIntPl(pl Player) int {
	return getPosInt(pl.X, pl.Y)
}

func getPosIntFl(fl FinishLine) int {
	return getPosInt(fl.X, fl.Y)
}

func isMoveOk(move int) bool {
	var pos = getPos(p1.X, p1.Y, move)
	var sym = board[pos]
	if sym != rune('.') {
		return true
	}
	return false
}

func executeMove(move int) {
	if move == 1 {
		p1.Y--
	}	
	if move == 2 {
		p1.X++
	}
	if move == 3 {
		p1.Y++
	}
	if move == 4 {
		p1.X--
	}
}

func getPos(x int, y int, move int) int {
	if move == 1 {
		return getPosInt(x, y-1)
	}	
	if move == 2 {
		return getPosInt(x+1, y)
	}
	if move == 3 {
		return getPosInt(x, y+1)
	}
	if move == 4 {
		return getPosInt(x-1, y) 
	}
	panic("Invalid move value!")
}

//1:up 2:right 3:down 4: left
func getMove() int {
    var rt = rand.Perm(4)
	var rrt = rt[1]+1
	return rrt
}

func wait(sec int) {
	var ssec = time.Duration(sec*1000000)
	time.Sleep(ssec)
}

func checkWin() bool{
   return false
}

