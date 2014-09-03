Race
====

Simple character racing game prototype written in Go. Definitely unfinished. But if nothing else displays a use of channels and termbox-go library, which is basically stripped down ncurses for Go.

![Alt text](/doc/race.png?raw=true "Screenshot")
### How to run:
```
sudo apt-get install golang
mkdir "$HOME"/go-packages
export GOPATH="$HOME/go-packages"
go get github.com/nsf/termbox-go
cd "$HOME"
git clone https://github.com/gto76/race.git
cd race/src
go run race.go
```
