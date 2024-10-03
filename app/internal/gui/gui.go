package gui

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"strings"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/LoreviQ/ChessAnalysis/app/internal/database"
)

type GUI struct {
	window *app.Window
	ops    *op.Ops
	theme  *chessAnalysisTheme

	// Components
	header  *header
	sidebar *sidebar
	board   *board

	// Database
	db *database.Database
}

type chessAnalysisTheme struct {
	giouiTheme      *material.Theme
	chessBoardTheme *chessBoardTheme
}

type chessBoardTheme struct {
	themeName     string
	square1Colour color.NRGBA
	square2Colour color.NRGBA
	player1Colour color.NRGBA
	player2Colour color.NRGBA
	pieces        *map[string]paint.ImageOp
}

func NewTheme() *chessAnalysisTheme {
	return &chessAnalysisTheme{
		giouiTheme:      material.NewTheme(),
		chessBoardTheme: NewChessBoardTheme(""),
	}
}

func NewChessBoardTheme(theme string) *chessBoardTheme {
	switch theme {
	default: // chess.com theme
		themeName := "chess.com"
		imageMap, err := loadImages(themeName)
		if err != nil {
			return nil
		}
		return &chessBoardTheme{
			themeName:     themeName,
			square1Colour: color.NRGBA{234, 236, 206, 255},
			square2Colour: color.NRGBA{114, 148, 82, 255},
			player1Colour: color.NRGBA{255, 255, 255, 255},
			player2Colour: color.NRGBA{64, 61, 57, 255},
			pieces:        imageMap,
		}
	}
}

// Returns a GUI struct
func NewGUI(width, height int, db *database.Database) *GUI {
	// Default window size
	if width == 0 {
		width = 2000
	}
	if height == 0 {
		height = 1200
	}

	// Create window
	w := new(app.Window)
	w.Option(app.Size(unit.Dp(width), unit.Dp(height)))
	w.Option(app.Title("Chess Analysis"))
	ops := new(op.Ops)
	th := NewTheme()
	th.giouiTheme.Palette.Bg = color.NRGBA{48, 46, 42, 255}
	th.giouiTheme.Palette.Fg = color.NRGBA{255, 255, 255, 255}

	// define components
	g := &GUI{
		window: w,
		ops:    ops,
		theme:  th,
		db:     db,
	}
	g.header = newHeader(g)
	g.sidebar = newSidebar(g)
	g.board = newBoard(g)

	return g
}

// CreateGUI creates the GUI
func (g *GUI) CreateGUI() {
	go g.draw()
	app.Main()
}

// Main event loop
func (g *GUI) draw() error {
	for {
		switch e := g.window.Event().(type) {

		// Re-render app
		case app.FrameEvent:
			gtx := app.NewContext(g.ops, e)
			g.Layout(gtx)
			e.Frame(gtx.Ops)
		// Exit app
		case app.DestroyEvent:
			os.Exit(0)
		}
	}
}

func (g *GUI) Layout(gtx layout.Context) layout.Dimensions {
	// Set Background
	rect := image.Rectangle{
		Max: image.Point{
			X: gtx.Constraints.Max.X,
			Y: gtx.Constraints.Max.Y,
		},
	}
	paint.FillShape(gtx.Ops, g.theme.giouiTheme.Palette.Bg, clip.Rect(rect).Op())
	return layout.Flex{Axis: layout.Vertical, Spacing: 0}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return g.header.Layout(gtx)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Spacing: 0}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return g.sidebar.Layout(gtx)
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return g.board.Layout(gtx)
				}),
			)
		}),
	)
}

func loadImages(themeName string) (*map[string]paint.ImageOp, error) {
	pieces := make(map[string]paint.ImageOp)
	dir := filepath.Join("assets", "images", themeName)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".png") {
			imageName := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))
			img, err := loadImage(path)
			if err != nil {
				return err
			}
			pieces[imageName] = img
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &pieces, nil
}

func loadImage(filename string) (paint.ImageOp, error) {
	file, err := os.Open(filename)
	if err != nil {
		return paint.ImageOp{}, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return paint.ImageOp{}, err
	}

	return paint.NewImageOp(img), nil
}
