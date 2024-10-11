package gui

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/LoreviQ/ChessAnalysis/app/internal/database"
	"github.com/LoreviQ/ChessAnalysis/app/internal/eval"
)

type GUI struct {
	window *app.Window
	ops    *op.Ops
	theme  *chessAnalysisTheme

	// Components
	header       *header
	sidebar      *sidebar
	board        *Board
	settingsMenu *settingsMenu

	// Database
	db *database.Database

	// Engine
	eng *eval.Engine
}

type chessAnalysisTheme struct {
	themeName       string
	giouiTheme      *material.Theme
	chessBoardTheme *chessBoardTheme
	bg              color.NRGBA
	fg              color.NRGBA
	contrastBg      color.NRGBA
	contrastFg      color.NRGBA
	text            color.NRGBA
	textMuted       color.NRGBA
}

type chessBoardTheme struct {
	square1 color.NRGBA
	square2 color.NRGBA
	player1 color.NRGBA
	player2 color.NRGBA
	pieces  map[string]*image.Image
}

type config struct {
	EnginePath string `json:"EnginePath"`
	SyzygyPath string `json:"SyzygyPath"`
	Movetime   int    `json:"Movetime"`
	Depth      int    `json:"Depth"`
	Threads    int    `json:"Threads"`
	Hash       int    `json:"Hash"`
	MultiPV    int    `json:"MultiPV"`
}

func NewTheme(theme string) *chessAnalysisTheme {
	switch theme {
	case "lichess.org":
		return &chessAnalysisTheme{
			themeName:       "lichess.org",
			giouiTheme:      material.NewTheme(),
			chessBoardTheme: NewChessBoardTheme("lichess.org"),
			bg:              color.NRGBA{223, 221, 216, 255},
			fg:              color.NRGBA{90, 90, 90, 255},
			contrastBg:      color.NRGBA{255, 255, 255, 255},
			contrastFg:      color.NRGBA{246, 246, 246, 255},
			text:            color.NRGBA{77, 77, 77, 255},
			textMuted:       color.NRGBA{90, 90, 90, 122},
		}
	case "HotDogStand":
		return &chessAnalysisTheme{
			themeName:       "HotDogStand",
			giouiTheme:      material.NewTheme(),
			chessBoardTheme: NewChessBoardTheme("HotDogStand"),
			bg:              color.NRGBA{0, 0, 0, 255},
			fg:              color.NRGBA{255, 255, 255, 255},
			contrastBg:      color.NRGBA{255, 0, 0, 255},
			contrastFg:      color.NRGBA{0, 0, 0, 255},
			text:            color.NRGBA{255, 255, 255, 255},
			textMuted:       color.NRGBA{255, 255, 255, 122},
		}
	default: // chess.com theme
		return &chessAnalysisTheme{
			themeName:       "chess.com",
			giouiTheme:      material.NewTheme(),
			chessBoardTheme: NewChessBoardTheme("chess.com"),
			bg:              color.NRGBA{48, 46, 42, 255},
			fg:              color.NRGBA{255, 255, 255, 255},
			contrastBg:      color.NRGBA{39, 37, 35, 255},
			contrastFg:      color.NRGBA{30, 30, 27, 255},
			text:            color.NRGBA{255, 255, 255, 255},
			textMuted:       color.NRGBA{255, 255, 255, 122},
		}
	}
}

func NewChessBoardTheme(theme string) *chessBoardTheme {
	switch theme {
	case "lichess.org":
		imageMap, err := loadImages(theme)
		if err != nil {
			imageMap, err = loadImages("chess.com")
			if err != nil {
				panic("Failed to load piece images")
			}
		}
		return &chessBoardTheme{
			square1: color.NRGBA{235, 218, 183, 255},
			square2: color.NRGBA{172, 138, 102, 255},
			player1: color.NRGBA{255, 255, 255, 255},
			player2: color.NRGBA{0, 0, 0, 255},
			pieces:  imageMap,
		}
	case "HotDogStand":
		imageMap, err := loadImages(theme)
		if err != nil {
			imageMap, err = loadImages("chess.com")
			if err != nil {
				panic("Failed to load piece images")
			}
		}
		return &chessBoardTheme{
			square1: color.NRGBA{255, 0, 0, 255},
			square2: color.NRGBA{255, 255, 0, 255},
			player1: color.NRGBA{255, 0, 0, 255},
			player2: color.NRGBA{255, 255, 0, 255},
			pieces:  imageMap,
		}
	default: // chess.com theme
		theme = "chess.com"
		imageMap, err := loadImages(theme)
		if err != nil {
			return nil
		}
		return &chessBoardTheme{
			square1: color.NRGBA{234, 236, 206, 255},
			square2: color.NRGBA{114, 148, 82, 255},
			player1: color.NRGBA{255, 255, 255, 255},
			player2: color.NRGBA{64, 61, 57, 255},
			pieces:  imageMap,
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
	th := NewTheme("")

	// load settings
	defaultSettings := &config{
		EnginePath: "",
		SyzygyPath: "",
		Movetime:   1000,
		Depth:      20,
		Threads:    4,
		Hash:       128,
		MultiPV:    1,
	}
	settings, err := loadConfig()
	if err != nil {
		settings = defaultSettings
	}
	settings = setDefaults(settings, defaultSettings)

	// define components
	g := &GUI{
		window: w,
		ops:    ops,
		theme:  th,
		db:     db,
	}
	g.eng, _ = eval.InitializeStockfish(settings.EnginePath, settings.SyzygyPath, settings.Movetime, settings.Depth, settings.Threads, settings.Hash, settings.MultiPV)
	g.header = newHeader(g)
	g.board = newBoard(g, nil)
	g.sidebar = newSidebar(g)
	g.settingsMenu = newSettingsMenu(g)

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
	paint.FillShape(gtx.Ops, g.theme.bg, clip.Rect(rect).Op())
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
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
		}),
		// Header Dropdown Menus
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			themeButton := g.header.buttons[0]
			return themeButton.layoutDropDown(gtx, func(gtx layout.Context) layout.Dimensions {
				offset := layout.Inset{
					Top:  unit.Dp(float32(g.header.size.Y) / +1),
					Left: unit.Dp(1),
				}
				return offset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min = image.Point{}
					menu := component.Menu(g.theme.giouiTheme, themeButton.menu)
					menu.SurfaceStyle.Fill = g.theme.contrastFg
					return menu.Layout(gtx)
				})
			})
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return g.settingsMenu.Layout(gtx)
		}),
	)
}

func loadImages(themeName string) (map[string]*image.Image, error) {
	pieces := make(map[string]*image.Image)
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
	return pieces, nil
}

func loadImage(filename string) (*image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return &img, nil
}

func loadConfig() (*config, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &config{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// saveConfig function to save the current settings to the JSON file
func saveConfig(cfg *config) error {
	filePath := "config.json"
	// Marshal the config map into JSON
	JSON, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}
	// Write the updated JSON back to the file
	if err := os.WriteFile(filePath, JSON, 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}
	return nil
}

func setDefaults(settings, defaultSettings *config) *config {
	settingsValue := reflect.ValueOf(settings).Elem()
	defaultSettingsValue := reflect.ValueOf(defaultSettings).Elem()

	for i := 0; i < settingsValue.NumField(); i++ {
		field := settingsValue.Field(i)
		if isZeroValue(field) {
			field.Set(defaultSettingsValue.Field(i))
		}
	}

	return settings
}

func isZeroValue(v reflect.Value) bool {
	return v.Interface() == reflect.Zero(v.Type()).Interface()
}
