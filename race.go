package main

import (
	"math/rand"
	"os"
	"io/ioutil"
	"time"
	"github.com/nsf/termbox-go"
)

// SETUP
const WAIT = 50
const ALLOW_BACK_OVER_LINE = false
const START_CHAR = 'S'
const FINISH_CHAR = 'F'

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
//var outFlg = true

func main() {
	// Termbox init
	termboxErr := termbox.Init()
	if termboxErr != nil {
        panic(termboxErr)
	}
	defer termbox.Close()
	termbox.HideCursor()
	// Track init
	var trackFileName string
	if  len(os.Args) > 1 {
		trackFileName = os.Args[1]
	} else {
		trackFileName = "t1.tr"
	}	
	var boardByte, readfileErr = ioutil.ReadFile(trackFileName)
	if (readfileErr != nil) {
		panic ("Could not open file with track")
	}
	var boardString = string(boardByte)
	board = []rune(boardString)
	tempBoard = make([]rune, len(board))
	getFinishLines()
	// Draw finish line instead of L and R
	drawFinishLine()
	// Put player on start (R)
	var f = fR[0]
	p1.X = f.X; p1.Y = f.Y
	// Draw starting position
	draw(getBoard())
	termbox.Flush()
	// Key reader goroutine
	go checkKey()
	// Main Loop	
	for !checkWin() {
		var move = getMove()
		if (isMoveOk(move)) {
			executeMove(move)
		}
		checkCircle()
		wait(WAIT)
		draw(getBoard())
		termbox.Flush()
	}
}

func drawFinishLine() {
	// gre cez cev board, zamena
	for key, val := range board {
		if val == START_CHAR {
			board[key] = ' '
		}
		if val == FINISH_CHAR {
			board[key] = '|'
		}
	}
}

func checkKey() {
	ev := termbox.PollEvent()
	go checkKey()
	switch ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyCtrlC {
				termbox.Close()
				panic("Don't know how else to exit program:(")
			}
		case termbox.EventError:
			panic(ev.Err)
	}
}

func draw(text string) {
	var x=0; var y=0
	for _, value := range text {
		if value == '\n' {
			x=0
			y++
			continue
		}
		termbox.SetCell(x, y, value, 	termbox.ColorWhite, 										termbox.ColorBlack)
		x++
	}
}
 
func getFinishLines() {
	for key, value := range board {
		if value == START_CHAR {
			var x, y = getPosXY(key)
			fR = append(fR,FinishLine{x,y})
		}
		if value == FINISH_CHAR {
			var x, y = getPosXY(key)
			fL = append(fL,FinishLine{x,y})
		}
	}
}

//TODO per player
var flFlag = false
var frFlag = false
func checkCircle() {
	// na liniji
	if doesArrayContainPosition(fL, getPosIntPl(p1)) {
		if (frFlag == true) {
			circles--
			flFlag = true
			frFlag = false
		} else {
			flFlag = true
		}
	// desno od linije
    } else if doesArrayContainPosition(fR, getPosIntPl(p1)) {
    	if (flFlag == true) {
			circles++
			frFlag = true
			flFlag = false
		} else {
			frFlag = true
		}
	// drugje
    } else {
		flFlag = false
		frFlag = false
	}
}

func doesArrayContainPosition(fll []FinishLine, position int) bool {
	for _, fl := range fll {
		if position == getPosIntFl(fl) {
			return true
		}
	}
	return false
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
	var newPos = getPos(p1.X, p1.Y, move)
	var sym = board[newPos]
	if contains([]rune{' ', '|'}, sym) {
		if !ALLOW_BACK_OVER_LINE {
			return !headingBackOverLine(getPosIntPl(p1), newPos)
		} else {
			return true
		}
	}
	return false
}

//TODO da ne gre niti na crto
func headingBackOverLine(oldPos int, newPos int) bool {
	// v naslednji potezi bi na ciljni crti
	// in zdaj je desno od ciljne crte
	if 	doesArrayContainPosition(fL, newPos) &&
		doesArrayContainPosition(fR, oldPos) {
		return true
	}
    return false
}

func contains(arr []rune, sim rune) bool {
	for _, val := range arr {
		if val == sim {
			return true
		}
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

