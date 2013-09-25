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
const DEFAULT_TRACK = "t1.tr"
const FINISH_LINE = '|'

// Board
var board []rune
var tempBoard = make([]rune, len(board))

// Player
type Player struct {
	Symbol	rune
	X 		int
	Y		int
	laps 	int//= 0
	flFlag  bool//= false
	frFlag 	bool//= false
}
var pl1 =	Player{2015, 3, 3, 0, false, false}
var pl2 = Player{2016, 3, 3, 0, false, false}
var players = []Player{pl1, pl2}

// Channels
var channels []chan int
var globalChangeChanel chan int

// Controler
func NewControler(pl *Player, moveFromInputer chan int) {
	for {
		move := <-moveFromInputer
		pl.moveIfPossible(move)
	}
}

// Inputers
func NewRandomInputer(moveToControler chan int) {
	for {
		wait(WAIT)
		move := getRandomMove()
		moveToControler <- move
	}
}

func NewKeyInputer(moveToControler chan int, moveFromKeyListenerControler chan int) {
	for {
		var move = <- moveFromKeyListenerControler
		moveToControler <- move
	}
}

// Finish Line
type FinishLine struct {
	X 		int
	Y		int
}
var fL []FinishLine
var fR []FinishLine

//type KeyGroup string
type KeyGroup int
const (
   	ARROW_KEYS = 0
	ASWD_KEYS = 1
)

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
		trackFileName = DEFAULT_TRACK
	}	
	var boardByte, readfileErr = ioutil.ReadFile(trackFileName)
	if (readfileErr != nil) {
		panic ("Could not open file with track")
	}
	var boardString = string(boardByte)
	board = []rune(boardString)
	tempBoard = make([]rune, len(board))
	getFinishLines()
	drawFinishLine()
	putPlayersOnStart()
	// Draw starting position
	draw(getBoard())
	termbox.Flush()

	// Controlers and imputers	
	globalChangeChanel = make(chan int)

	//connectModules()	
	// create channels
	for i := 0; i < len(players); i++ {
		var chanel = make(chan int)
		channels = append(channels, chanel)
	}
	// create modules
	go NewControler(&players[0], channels[0])
	go NewRandomInputer(channels[0])
	//go NewControler(&players[1], channels[1])
	//go NewRandomInputer(channels[1])
	var keyListener2KeyInputer = make(chan int)
	go NewControler(&players[1], channels[1])
	go NewKeyInputer(channels[1], keyListener2KeyInputer)

	// Key reader goroutine
	var keyGroup2ChannelMaper = map[KeyGroup] chan int {
		ARROW_KEYS: keyListener2KeyInputer,	
	} 
	go listenToKeys(keyGroup2ChannelMaper)
	
	// Main Loop	
	for !checkWin() {
		<- globalChangeChanel
		brd := getBoard()
		draw(brd)
		termbox.Flush()
	}
}

/*
func connectModules() {
	// create channels
	for i := 0; i < len(players); i++ {
		var chanel = make(chan int)
		channels = append(channels, chanel)
	}
	// create modules
	go NewControler(&players[0], channels[0])
	go NewRandomInputer(channels[0])

	//go NewControler(&players[1], channels[1])
	//go NewRandomInputer(channels[1])
	
	var keyListener2KeyInputer = make(chan int)
	go NewControler(&players[1], channels[1])
	go NewKeyInputer(channels[1], keyListener2KeyInputer)
	
	/*
	for i, ch := range channels {
		go NewControler(&players[i], ch)
		go NewRandomInputer(ch)
	}
	*/
//}

func putPlayersOnStart() {
	i := 0
	for _, pl := range players {
		if len(fR) == i {
			i = 0
		}
		sL := fR[i]
		i++
		pl.X = sL.X 
		pl.Y = sL.Y
	}
}

func (pl *Player) moveIfPossible(move int) {
	if (pl.isMoveOk(move)) {
		pl.executeMove(move)
	}
}

func drawFinishLine() {
	// gre cez cev board, zamena
	for key, val := range board {
		if val == START_CHAR {
			board[key] = ' '
		}
		if val == FINISH_CHAR {
			board[key] = FINISH_LINE
		}
	}
}



/*
const (
        _           = iota // ignore first value by assigning to blank identifier
        KB ByteSize = 1 << (10 * iota)
        MB
        GB
        TB
        PB
        EB
        ZB
        YB
)
const KeyGroup = (
	ARROW_KEYS = 0
	ASWD_KEYS = 1
)
*/

//TODO change to key listener module
func listenToKeys(chnls map[KeyGroup]chan int) {
	ev := termbox.PollEvent()
	// When it gets event it recusively calls itself,
	// so it listens for new events as soon as possible.
	go listenToKeys(chnls)
	switch ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyCtrlC {
				termbox.Close()
				panic("Don't know how else to exit program:(")
			}
			ch, present := chnls[ARROW_KEYS]
			if present {
				if ev.Key == termbox.KeyArrowUp {
					ch <- 1
				}
				if ev.Key == termbox.KeyArrowRight {
					ch <- 2
				}
				if ev.Key == termbox.KeyArrowDown {
					ch <- 3
				}
				if ev.Key == termbox.KeyArrowLeft {
					ch <- 4
				}
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

func (pl *Player) checkCircle() {
	// na liniji
	if doesArrayContainPosition(fL, pl.getPosInt()) {
		if (pl.frFlag == true) {
			pl.laps--
			pl.flFlag = true
			pl.frFlag = false
		} else {
			pl.flFlag = true
		}
	// desno od linije
    } else if doesArrayContainPosition(fR, pl.getPosInt()) {
    	if (pl.flFlag == true) {
			pl.laps++
			pl.frFlag = true
			pl.flFlag = false
		} else {
			pl.frFlag = true
		}
	// drugje
    } else {
		pl.flFlag = false
		pl.frFlag = false
	}
}

func doesArrayContainPosition(fll []FinishLine, position int) bool {
	for _, fl := range fll {
		if position == fl.getPosInt() {
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
	for _, pl := range players {
		pl.drawPlayerOnBoard(tempBoard)
	}
	return string(tempBoard)
}

func (pl *Player) drawPlayerOnBoard(boardLoc []rune) {
	var posInt = pl.getPosInt()
	boardLoc[posInt] = pl.Symbol
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

func (pl *Player) getPosInt() int {
	return getPosInt(pl.X, pl.Y)
}

func (fl *FinishLine) getPosInt() int {
	return getPosInt(fl.X, fl.Y)
}

func (pl *Player) isMoveOk(move int) bool {
	var newPos = getPos(pl.X, pl.Y, move)
	var sym = board[newPos]
	if contains([]rune{' ', '|'}, sym) {
		if !ALLOW_BACK_OVER_LINE {
			return !headingBackOverLine(pl.getPosInt(), newPos)
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

func (pl *Player) executeMove(move int) {
	if move == 1 {
		pl.Y--
	}	
	if move == 2 {
		pl.X++
	}
	if move == 3 {
		pl.Y++
	}
	if move == 4 {
		pl.X--
	}
	pl.checkCircle()
	globalChangeChanel <- 1
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
func getRandomMove() int {
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

