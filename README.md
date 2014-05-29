race
====
![Alt text](/doc/race.png?raw=true "Screenshot")

Just a simple character racing game against erratic computer opponent. Defenately unfinished. But if nothing else displays a use of termbox-go library, which is basicaly stripped down ncurses for go.

### How to run:
```
sudo apt-get install golang
mkdir "$HOME"/go-packages
export GOPATH="$HOME/go-packages"
go get github.com/nsf/termbox-go
mkdir "$HOME"/race
cd "$HOME"
git clone https://github.com/gto76/race.git
cd race
go run race.go
```
