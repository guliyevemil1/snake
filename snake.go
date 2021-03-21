package main

import (
	"time"

	"github.com/nsf/termbox-go"
)

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func listen(ch chan<- termbox.Event) {
	for {
		ev := termbox.PollEvent()
		ch <- ev
	}
}

type Tuple struct {
	X, Y int
}

type Box struct {
	W, H  int
	Snake []Tuple
	Dir   int
}

func (b *Box) Move() {
	head := b.Snake[len(b.Snake)-1]
	var t Tuple
	switch b.Dir {
	case 0:
		t.X, t.Y = head.X+1, head.Y
	case 1:
		t.X, t.Y = head.X, head.Y+1
	case 2:
		t.X, t.Y = head.X-1, head.Y
	case 3:
		t.X, t.Y = head.X, head.Y-1
	}
	if t.X < 0 {
		t.X += b.W - 2
	}
	if t.Y < 0 {
		t.Y += b.H - 2
	}
	t.X, t.Y = t.X%(b.W-2), t.Y%(b.H-2)
	b.Snake = append(b.Snake[1:], t)
}

func (b *Box) DrawBorder() {
	for x := 0; x < b.W; x++ {
		termbox.SetCell(x, 0, ' ', termbox.ColorWhite, termbox.ColorWhite)
		termbox.SetCell(x, b.H-1, ' ', termbox.ColorWhite, termbox.ColorWhite)
	}
	for y := 0; y < b.H; y++ {
		termbox.SetCell(0, y, ' ', termbox.ColorWhite, termbox.ColorWhite)
		termbox.SetCell(b.W-1, y, ' ', termbox.ColorWhite, termbox.ColorWhite)
	}
}

func (b *Box) Draw() {
	b.DrawBorder()
	for i := range b.Snake {
		s := b.Snake[i]
		termbox.SetCell(s.X+1, s.Y+1, 'O', termbox.ColorWhite, termbox.ColorBlack)
	}
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	ch := make(chan termbox.Event)
	go listen(ch)

	var b Box

	b.Snake = []Tuple{{0, 0}}

	b.W, b.H = termbox.Size()

loop:
	for {
		select {
		case ev := <-ch:
			switch ev.Type {
			case termbox.EventResize:
				break loop
			case termbox.EventKey:
				if ev.Key == termbox.KeyCtrlC {
					break loop
				}
				switch ev.Ch {
				case 'w':
					b.Dir = 3
				case 'a':
					b.Dir = 2
				case 's':
					b.Dir = 1
				case 'd':
					b.Dir = 0
				}
			}
		default:
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			b.Move()
			b.Draw()
			termbox.Flush()
			time.Sleep(50 * time.Millisecond)
		}
	}
}
