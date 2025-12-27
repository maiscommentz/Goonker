package ui

import (
	"Goonker/common"
	"fmt"
	"image/color"
	"log"
	"strings"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	// Button configuration
	ButtonWidth        = 200.0
	ButtonHeight       = 60.0
	ButtonCornerRadius = 10.0
	ButtonTextYAnchor  = 0.35

	// Button positions
	MainMenuPlayBtnY        = 200.0
	MainMenuQuitBtnY        = 280.0
	RoomsMenuBackBtnY       = 150.0
	RoomsMenuPlayBotBtnY    = 230.0
	RoomsMenuJoinGameBtnY   = 310.0
	RoomsMenuCreateRoomBtnY = 390.0
	RoomsMenuBtnX           = 50.0

	// Rooms list
	RoomsLineWidth  = 2
	RoomsRowPadding = 10
	RoomsRowGap     = 70.0

	// Spinning wheel parameters
	WheelTintRed   = 0.8
	WheelTintGreen = 0.8
	WheelTintBlue  = 1.0

	// Room ID text height
	WaitingMenuRoomTextY = (float64(WindowHeight) / 2.0) - 100

	// Waiting room text field
	WaitingMenuTextFieldX    = (float64(WindowWidth)-WaitingMenuTextFieldW)/2 + 300
	WaitingMenuTextFieldY    = (float64(WindowHeight)-WaitingMenuTextFieldH)/2 - 200
	WaitingMenuTextFieldW    = 300
	WaitingMenuTextFieldH    = 50
	WaitingMenuTextFieldFont = 16

	// Assets
	FontPath = "font.ttf"
)

type Button struct {
	X, Y, Width, Height float64
	Image               *ebiten.Image
	Text                string
}

type TextField struct {
	X, Y          float64
	Width, Height float64
	Text          string
	Focused       bool
	MaxLength     int
	Image         *ebiten.Image

	cursorVisible bool
	cursorTimer   int
	fontSize      float64
}

type MainMenu struct {
	BtnPlay *Button
	BtnQuit *Button
}

type WaitingMenu struct {
	RotationAngle float64
	RoomId        string
}

type RoomsMenu struct {
	Rooms         []*Room
	RoomIndex     int
	BtnPlayBot    *Button
	BtnCreateRoom *Button
	BtnJoinGame   *Button
	BtnBack       *Button
	RoomField     *TextField
}

type Room struct {
	JoinBtn *Button
	Id      string
	Image   *ebiten.Image
}

type Grid struct {
	Col       int
	BoardData [GridCol][GridCol]common.PlayerID
}

type Cell struct {
	Btn    Button
	Symbol common.PlayerID
}

type Drawable interface {
	Draw(screen *ebiten.Image)
}

// Constructor for the button.
func NewButton(x, y, w, h float64, text string, fontSize float64) *Button {
	b := &Button{
		X: x, Y: y, Width: w, Height: h,
		Text: text,
	}

	dc := gg.NewContext(int(w), int(h))

	dc.DrawRoundedRectangle(0, 0, w, h, ButtonCornerRadius)
	dc.SetHexColor(gridBorderColor)
	dc.Fill()

	dc.SetFontFace(FontFace)
	dc.SetHexColor(gridBackgroundColor)
	dc.DrawStringAnchored(text, w/2, h/2, 0.5, ButtonTextYAnchor)

	b.Image = ebiten.NewImageFromImage(dc.Image())

	return b
}

func NewTextField(x, y, w, h float64, fontSize float64) *TextField {
	tf := &TextField{
		X:         x,
		Y:         y,
		Width:     w,
		Height:    h,
		Text:      "",
		MaxLength: 50,
		fontSize:  fontSize,
	}
	tf.redraw()
	return tf
}

func (tf *TextField) redraw() {
	dc := gg.NewContext(int(tf.Width), int(tf.Height))

	// Background
	dc.DrawRoundedRectangle(0, 0, tf.Width, tf.Height, 5)
	if tf.Focused {
		dc.SetColor(color.RGBA{240, 240, 255, 255})
	} else {
		dc.SetColor(color.RGBA{255, 255, 255, 255})
	}
	dc.Fill()

	// Border
	dc.DrawRoundedRectangle(0, 0, tf.Width, tf.Height, 5)
	dc.SetLineWidth(2)
	if tf.Focused {
		dc.SetColor(color.RGBA{0, 100, 255, 255})
	} else {
		dc.SetColor(color.RGBA{150, 150, 150, 255})
	}
	dc.Stroke()

	// Text
	dc.SetFontFace(FontFace)
	dc.SetColor(color.Black)
	dc.DrawString(tf.Text, 10, tf.Height/2+tf.fontSize/3)

	// Blinking cursor
	if tf.Focused && tf.cursorVisible {
		textWidth, _ := dc.MeasureString(tf.Text)
		cursorX := 10 + textWidth
		cursorY1 := tf.Height/2 - tf.fontSize/2
		cursorY2 := tf.Height/2 + tf.fontSize/2

		dc.SetLineWidth(2)
		dc.SetColor(color.Black)
		dc.DrawLine(cursorX, cursorY1, cursorX, cursorY2)
		dc.Stroke()
	}

	tf.Image = ebiten.NewImageFromImage(dc.Image())
}

// Constructor for the room.
func NewRoom(id string) *Room {
	// Room row dimensions
	width := WindowWidth / 2
	height := TitleFontSize

	// Initialize room
	room := &Room{
		Id: id,
	}

	dc := gg.NewContext(width, height)

	// Draw bottom separator line
	dc.SetHexColor(gridBorderColor)
	dc.SetLineWidth(RoomsLineWidth)
	dc.DrawLine(0, float64(height), float64(width), float64(height))
	dc.Stroke()

	// Load font for room name
	if err := dc.LoadFontFace(FontPath, TextFontSize); err != nil {
		log.Printf("Error loading font: %v", err)
	}

	// Draw Room Name (Left aligned)
	dc.SetHexColor(gridBorderColor)
	dc.DrawStringAnchored(id, RoomsRowPadding, float64(height)/3, 0.0, 0.5)

	room.Image = ebiten.NewImageFromImage(dc.Image())

	// Initialize Join Button
	room.JoinBtn = NewButton(0, 0, ButtonWidth/3, ButtonHeight/2, "Join", TextFontSize)

	return room
}

// Constructor for the main menu.
func NewMainMenu() *MainMenu {
	menu := &MainMenu{}

	// Center buttons
	centerX := (float64(WindowWidth) - ButtonWidth) / 2

	// Create buttons
	menu.BtnPlay = NewButton(centerX, MainMenuPlayBtnY, ButtonWidth, ButtonHeight, "Play", SubtitleFontSize)
	menu.BtnQuit = NewButton(centerX, MainMenuQuitBtnY, ButtonWidth, ButtonHeight, "Quit", SubtitleFontSize)

	return menu
}

// Constructor for the play menu.
func NewRoomsMenu() *RoomsMenu {
	menu := &RoomsMenu{}

	// Create buttons
	menu.BtnBack = NewButton(RoomsMenuBtnX, RoomsMenuBackBtnY, ButtonWidth, ButtonHeight, "Back", SubtitleFontSize)
	menu.BtnPlayBot = NewButton(RoomsMenuBtnX, RoomsMenuPlayBotBtnY, ButtonWidth, ButtonHeight, "Against Bot", SubtitleFontSize)
	menu.BtnCreateRoom = NewButton(RoomsMenuBtnX, RoomsMenuCreateRoomBtnY, ButtonWidth, ButtonHeight, "Create Room", SubtitleFontSize)
	menu.BtnJoinGame = NewButton(RoomsMenuBtnX, RoomsMenuJoinGameBtnY, ButtonWidth, ButtonHeight, "Join Game", SubtitleFontSize)

	// Create textfield
	menu.RoomField = NewTextField(WaitingMenuTextFieldX, WaitingMenuTextFieldY, WaitingMenuTextFieldW, WaitingMenuTextFieldH, WaitingMenuTextFieldFont)

	return menu
}

func (tf *TextField) Update() {
	// Handle click
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		wasFocused := tf.Focused
		tf.Focused = float64(mx) >= tf.X && float64(mx) <= tf.X+tf.Width &&
			float64(my) >= tf.Y && float64(my) <= tf.Y+tf.Height

		if wasFocused != tf.Focused {
			tf.redraw()
		}
	}

	if !tf.Focused {
		return
	}

	needsRedraw := false

	// Blinking cursor
	tf.cursorTimer++
	if tf.cursorTimer >= 30 {
		tf.cursorVisible = !tf.cursorVisible
		tf.cursorTimer = 0
		needsRedraw = true
	}

	// Handle input
	runes := ebiten.AppendInputChars(nil)
	if len(runes) > 0 {
		for _, r := range runes {
			if len(tf.Text) < tf.MaxLength {
				tf.Text += string(r)
			}
		}
		needsRedraw = true
		tf.cursorVisible = true
		tf.cursorTimer = 0
	}

	// Backspace
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && len(tf.Text) > 0 {
		tf.Text = tf.Text[:len(tf.Text)-1]
		needsRedraw = true
		tf.cursorVisible = true
		tf.cursorTimer = 0
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		tf.Focused = false
		needsRedraw = true
	}

	if needsRedraw {
		tf.redraw()
	}
}

func (tf *TextField) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(tf.X, tf.Y)
	screen.DrawImage(tf.Image, opts)
}

// Draw the button to the screen.
func (b *Button) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(b.X, b.Y)
	screen.DrawImage(b.Image, opts)
}

// Draw the main menu to the screen.
func (m *MainMenu) Draw(screen *ebiten.Image) {
	screen.DrawImage(MainMenuImage, nil)
	m.BtnPlay.Draw(screen)
	m.BtnQuit.Draw(screen)
}

// Draw the rooms menu to the screen.
func (m *RoomsMenu) Draw(screen *ebiten.Image) {
	screen.DrawImage(RoomsMenuImage, nil)
	m.BtnPlayBot.Draw(screen)
	m.BtnCreateRoom.Draw(screen)
	m.BtnJoinGame.Draw(screen)
	m.BtnBack.Draw(screen)
	m.RoomField.Draw(screen)

	searchedId := m.RoomField.Text

	for i, room := range m.Rooms {
		if searchedId == "" || strings.HasPrefix(room.Id, searchedId) {
			room.Draw(screen, i)
		}
	}

	if len(m.Rooms) == 0 {
		listX := float64(RoomsMenuBtnX + ButtonWidth + RoomsMenuBtnX)
		listY := float64(RoomsMenuBackBtnY)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(listX, listY)
		screen.DrawImage(NoRoomsImage, op)
	}
}

// Draw the room at the specified index
func (r *Room) Draw(screen *ebiten.Image, index int) {
	// Constants for layout
	listX := float64(RoomsMenuBtnX + ButtonWidth + RoomsMenuBtnX)
	listY := float64(RoomsMenuBackBtnY)
	yPos := listY + float64(index)*RoomsRowGap

	// Draw the Row Image (Name + Line)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(listX, yPos)
	screen.DrawImage(r.Image, op)

	// Update Button Position
	btnX := listX + listX + ButtonWidth/2
	r.JoinBtn.X = btnX
	r.JoinBtn.Y = yPos

	// Draw Button
	r.JoinBtn.Draw(screen)
}

// Draw the waiting menu to the screen.
func (waitingMenu *WaitingMenu) Draw(screen *ebiten.Image) {
	screen.DrawImage(WaitingMenuImage, nil)
	screenCenterX := float64(WindowWidth) / 2.0
	screenCenterY := float64(WindowHeight) / 2.0

	// Draw the spinning wheel
	w := WheelImage.Bounds().Dx()
	h := WheelImage.Bounds().Dy()
	halfW := float64(w) / 2.0
	halfH := float64(h) / 2.0

	wheelOpt := &ebiten.DrawImageOptions{}

	wheelOpt.GeoM.Translate(-halfW, -halfH)
	wheelOpt.GeoM.Rotate(waitingMenu.RotationAngle)
	wheelOpt.GeoM.Translate(screenCenterX, screenCenterY)
	wheelOpt.ColorScale.Scale(WheelTintRed, WheelTintGreen, WheelTintBlue, 1)

	screen.DrawImage(WheelImage, wheelOpt)

	// Draw the text
	waitingRoomText := fmt.Sprintf("Room ID : %s", waitingMenu.RoomId)
	textOpt := &text.DrawOptions{}

	textWidth, _ := text.Measure(waitingRoomText, GameFont, textOpt.LineSpacing)
	x := (screenCenterX - (textWidth / 2))

	textOpt.GeoM.Translate(x, WaitingMenuRoomTextY)

	textOpt.ColorScale.ScaleWithColor(color.Black)

	text.Draw(screen, waitingRoomText, GameFont, textOpt)
}

// Check if a button is clicked.
func (b *Button) IsClicked() bool {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		fmx, fmy := float64(mx), float64(my)

		return fmx >= b.X && fmx <= b.X+b.Width &&
			fmy >= b.Y && fmy <= b.Y+b.Height
	}
	return false
}

// Get the clicked cell.
func (g *Grid) OnClick() (int, int, bool) {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return -1, -1, false
	}

	mx, my := ebiten.CursorPosition()

	gridW, gridH := GridImage.Bounds().Dx(), GridImage.Bounds().Dy()

	offsetX := (WindowWidth - gridW) / 2
	offsetY := (WindowHeight - gridH) / 2

	localX := mx - offsetX
	localY := my - offsetY

	if localX < 0 || localY < 0 || localX >= gridW || localY >= gridH {
		return -1, -1, false
	}

	cellSize := gridW / GridCol

	cellX := localX / cellSize
	cellY := localY / cellSize

	return cellX, cellY, true
}
