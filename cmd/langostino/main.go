package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/mikelangelon/langostino/pkg/keys"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Mode int

const (
	ModeTitle Mode = iota
	ModeGame
	ModeGameOver

	screenWidth   = 640
	screenHeight  = 480
	tileSize      = 32
	titleFontSize = fontSize * 1.5
	fontSize      = 24
	smallFontSize = fontSize / 2

	boxSize = 10
)

const (
	title = "Langostino"
)

var (
	titleArcadeFont font.Face
	arcadeFont      font.Face
)

func init() {
	tt, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	titleArcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    titleFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	arcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	mode Mode

	position position
	speed    position

	danger position
}

type position struct {
	x int
	y int
}

func (g *Game) Update() error {
	switch g.mode {
	case ModeTitle:
		if keys.IsAnyKeyJustPressed() {
			g.mode = ModeGame
		}
	default:
		if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
			g.speed = position{
				x: -5,
				y: 0,
			}
		} else if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
			g.speed = position{
				x: 5,
				y: 0,
			}
		} else if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
			g.speed = position{
				x: 0,
				y: -5,
			}
		} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
			g.speed = position{
				x: 0,
				y: 5,
			}
		}
		temp := g.position
		temp.x += g.speed.x
		temp.y += g.speed.y
		if temp.x <= 0 {
			temp.x = 0
		} else if temp.x > (screenWidth -boxSize) {
			temp.x =  (screenWidth -boxSize)
		}
		if temp.y <= 0 {
			temp.y = 0
		} else if temp.y > (screenHeight - boxSize) {
			temp.y = (screenHeight - boxSize)
		}
		g.position = temp

		if g.hit() {
			g.mode = ModeTitle
		}
	}
	return nil
}

func (g *Game) hit() bool{
	if !(g.position.x <= g.danger.x && g.position.x + 10 >= g.danger.x) {
		return false;
	}

	if !(g.position.y <= g.danger.y && g.position.y + 10 >= g.danger.y) {
		return false
	}

	return true
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x80, 0xa0, 0xc0, 0xff})
	switch g.mode {
	case ModeTitle:
		x := (screenWidth - len(title)*titleFontSize) / 2
		text.Draw(screen, title, titleArcadeFont, x, 4*titleFontSize, color.White)
		for i, l := range []string{"PRESS ANY KEY OR BUTTON", "OR TOUCH SCREEN"} {
			x := (screenWidth - len(l)*fontSize) / 2
			text.Draw(screen, l, arcadeFont, x, (i+10)*fontSize, color.White)
		}
	default:
		ebitenutil.DrawRect(screen, float64(g.position.x), float64(g.position.y), boxSize, boxSize, color.Black)
		ebitenutil.DrawRect(screen, float64(g.danger.x), float64(g.danger.y), boxSize, boxSize, color.White)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func NewGame() *Game {
	g := &Game{
		position: position{
			x: 100,
			y: 100,
		},
		danger: position{
			x:0,
			y:50,
		},
	}
	return g
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Langostino (Demo)")
	if err := ebiten.RunGame(NewGame()); err != nil {
		panic(err)
	}
}
