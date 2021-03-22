package main

import (
	"math/rand"
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
	Dead  bool
	Snake []Tuple
	Dir   int
	Food  []Tuple
	Turn  int64
}

func (b *Box) Move() {
	b.Turn++
	w, h := b.W-2, b.H-2
	head := b.Snake[len(b.Snake)-1]
	prev := Tuple{-1, -1}
	if len(b.Snake) > 1 {
		prev = b.Snake[len(b.Snake)-2]
	}
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
	if prev.X != -1 && prev.Y != -1 {
		if t.X == prev.X && t.Y == prev.Y {
			b.Dead = true
			return
		}
	}
	if t.X < 0 || t.X >= w || t.Y < 0 || t.Y >= h {
		b.Dead = true
		return
	}
	var ate bool
	for i, f := range b.Food {
		if f.X == t.X && f.Y == t.Y {
			ate = true
			copy(b.Food[i:], b.Food[i+1:])
			b.Food = b.Food[:len(b.Food)-1]
			break
		}
	}
	if ate {
		b.Snake = append(b.Snake, t)
	} else {
		copy(b.Snake, b.Snake[1:])
		b.Snake[len(b.Snake)-1] = t
	}
	if len(b.Food) == 0 || rand.Intn(100) == 0 {
	food:
		for {
			x, y := rand.Intn(w), rand.Intn(h)
			for _, t := range b.Snake {
				if x == t.X && y == t.Y {
					continue food
				}
			}
			b.Food = append(b.Food, Tuple{x, y})
			break
		}
	}
	if len(b.Food) > 1 {
		if b.Turn%100 == 0 {
			b.Food = b.Food[1:]
		}
	}
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
	for i := range b.Food {
		f := b.Food[i]
		termbox.SetCell(f.X+1, f.Y+1, 'X', termbox.ColorWhite, termbox.ColorBlack)
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
			if b.Dead {
				break loop
			}
			b.Draw()
			termbox.Flush()
			time.Sleep(50 * time.Millisecond)
		}
	}
}
