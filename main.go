package main

import (
	"bytes"
	"fmt"
	"github.com/YYYMagic/funny-snake/cmd"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/colornames"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math/rand"
	"os"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
)

const (
	Left = iota
	Right
	Up
	Down
	None
)

var headImg *ebiten.Image
var xiangImg *ebiten.Image
var size int
var xNum, yNum int

func init() {
	var err error
	img, _, err := image.Decode(bytes.NewReader(funny))
	if err != nil {
		fmt.Println("image.Decode error: ", err)
		os.Exit(1)
	}
	headImg = ebiten.NewImageFromImage(img)
	if err != nil {
		panic(err)
	}
	size, _ = headImg.Size()
	xNum = ScreenWidth / size
	yNum = ScreenHeight / size

	img, _, err = image.Decode(bytes.NewReader(xiang))
	if err != nil {
		fmt.Println("image.Decode error: ", err)
		os.Exit(1)
	}
	xiangImg = ebiten.NewImageFromImage(img)
	if err != nil {
		panic(err)
	}

}

type Snake struct {
	body      []image.Point
	direction int
}

func NewSnake(x, y int) *Snake {
	s := &Snake{body: make([]image.Point, 1), direction: None}
	s.body[0].X = x
	s.body[0].Y = y
	return s
}

func (s *Snake) Draw(screen *ebiten.Image) {
	for i := 0; i < len(s.body); i++ {
		if i == 0 {
			opt := &ebiten.DrawImageOptions{}
			opt.GeoM.Translate(float64(s.body[i].X*size), float64(s.body[i].Y*size))
			screen.DrawImage(headImg, opt)
			continue
		}
		ebitenutil.DrawRect(screen, float64(s.body[i].X*size), float64(s.body[i].Y*size), float64(size), float64(size), colornames.Ghostwhite)
	}
}

func (s *Snake) SetDir(ids []ebiten.GamepadID) {
	// keyboard
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		if s.direction != Down {
			s.direction = Up
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		if s.direction != Up {
			s.direction = Down
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		if s.direction != Right {
			s.direction = Left
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		if s.direction != Left {
			s.direction = Right
		}
	}

	for _, p := range ids {
		// standard button
		if ebiten.IsStandardGamepadLayoutAvailable(p) {
			if inpututil.IsStandardGamepadButtonJustPressed(p, ebiten.StandardGamepadButtonRightTop) {
				if s.direction != Down {
					s.direction = Up
				}
			} else if inpututil.IsStandardGamepadButtonJustPressed(p, ebiten.StandardGamepadButtonRightBottom) {
				if s.direction != Up {
					s.direction = Down
				}
			} else if inpututil.IsStandardGamepadButtonJustPressed(p, ebiten.StandardGamepadButtonRightLeft) {
				if s.direction != Right {
					s.direction = Left
				}
			} else if inpututil.IsStandardGamepadButtonJustPressed(p, ebiten.StandardGamepadButtonRightRight) {
				if s.direction != Left {
					s.direction = Right
				}
			}
		}

		// nintendo switch axis
		if inpututil.IsGamepadButtonJustPressed(p, ebiten.GamepadButton16) {
			if s.direction != Down {
				s.direction = Up
			}
		} else if inpututil.IsGamepadButtonJustPressed(p, ebiten.GamepadButton18) {
			if s.direction != Up {
				s.direction = Down
			}
		} else if inpututil.IsGamepadButtonJustPressed(p, ebiten.GamepadButton19) {
			if s.direction != Right {
				s.direction = Left
			}
		} else if inpututil.IsGamepadButtonJustPressed(p, ebiten.GamepadButton17) {
			if s.direction != Left {
				s.direction = Right
			}
		}
	}
}

func (s *Snake) CollideWithPoint(p *image.Point) bool {
	return s.body[0].X == p.X &&
		s.body[0].Y == p.Y
}

func (s *Snake) CollideWithSelf() bool {
	for _, v := range s.body[1:] {
		if s.body[0].X == v.X &&
			s.body[0].Y == v.Y {
			return true
		}
	}
	return false
}

func (s *Snake) CollideWithBounds(p1, p2 *image.Point) bool {
	return s.body[0].X < p1.X ||
		s.body[0].Y < p1.Y ||
		s.body[0].X >= p2.X ||
		s.body[0].Y >= p2.Y
}

func (s *Snake) Growth() {
	s.body = append(s.body, image.Point{
		X: s.body[len(s.body)-1].X,
		Y: s.body[len(s.body)-1].Y,
	})
}

func (s *Snake) ShouldMove(f func() bool) bool {
	return f()
}

func (s *Snake) Move() {
	for i := int64(len(s.body)) - 1; i > 0; i-- {
		s.body[i].X = s.body[i-1].X
		s.body[i].Y = s.body[i-1].Y
	}
	switch s.direction {
	case Left:
		s.body[0].X--
		if s.body[0].X == -1 {
			s.body[0].X = xNum
		}
	case Right:
		s.body[0].X++
		if s.body[0].X == xNum {
			s.body[0].X = 0
		}
	case Down:
		s.body[0].Y++
		if s.body[0].Y == yNum {
			s.body[0].Y = 0
		}
	case Up:
		s.body[0].Y--
		if s.body[0].Y == -1 {
			s.body[0].Y = yNum
		}
	}
}

func (g *Game) Reset() {
	g.apple.X = 2
	g.apple.Y = 2
	g.tMove = 4
	g.snake = NewSnake(xNum/2, yNum/2)
	g.score = 0
	g.level = 1
}

func (g *Game) CheckReset() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.Reset()
	}
	for _, p := range g.gamepadIDs {
		// standard button
		if ebiten.IsStandardGamepadLayoutAvailable(p) {
			if inpututil.IsStandardGamepadButtonJustPressed(p, ebiten.StandardGamepadButtonCenterLeft) {
				g.Reset()
			}
		}
	}
}

type Game struct {
	snake      *Snake
	apple      *image.Point
	stoneCount int
	tCount     int
	tMove      int
	score      int
	level      int

	gamepadIDs []ebiten.GamepadID
}

func (g *Game) Update() error {
	g.tCount++

	g.gamepadIDs = ebiten.AppendGamepadIDs(g.gamepadIDs[:0])
	g.snake.SetDir(g.gamepadIDs)
	g.CheckReset()

	shouldMove := g.snake.ShouldMove(func() bool {
		return g.tCount%g.tMove == 0
	})
	if !shouldMove {
		return nil
	}
	g.tCount = 0

	if g.snake.CollideWithSelf() {
		g.Reset()
		return nil
	}

	if g.snake.CollideWithPoint(g.apple) {
		g.apple.X = rand.Intn(xNum - 1)
		g.apple.Y = rand.Intn(yNum - 1)
		g.snake.Growth()
		if len(g.snake.body) > 10 && len(g.snake.body) <= 20 {
			g.level = 2
			g.tMove = 3
		} else if len(g.snake.body) > 20 {
			g.level = 3
			g.tMove = 2
		} else {
			g.level = 1
		}
		g.score++
	}

	g.snake.Move()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(float64(g.apple.X*size), float64(g.apple.Y*size))
	screen.DrawImage(xiangImg, opt)
	g.snake.Draw(screen)
	if g.snake.direction == None {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("Press up/down/left/right or axis to start"))
	} else {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("Level: %d Score: %d", g.level, g.score))
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	cmd.Execute()

	g := &Game{
		apple: &image.Point{0, 0},
	}
	g.Reset()

	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Funny Snake")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
