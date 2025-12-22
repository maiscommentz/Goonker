package ui

import (
	"Goonker/common"
	"log"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	RoomsMenuCreateRoomBtnY = 310.0
	RoomsMenuBtnX           = 50.0

	// Spinning wheel parameters
	WheelTintRed   = 0.8
	WheelTintGreen = 0.8
	WheelTintBlue  = 1.0

	// Assets
	FontPath = "client/assets/font.ttf"
)

type Button struct {
	X, Y, Width, Height float64
	Image               *ebiten.Image
	Text                string
}

type MainMenu struct {
	BtnPlay *Button
	BtnQuit *Button
}

type WaitingMenu struct {
	RotationAngle float64
}

type RoomsMenu struct {
	Rooms         []*Room
	RoomIndex     int
	BtnPlayBot    *Button
	BtnCreateRoom *Button
	BtnBack       *Button
}

type Room struct {
	JoinBtn     *Button
	Id          string
	PlayerCount int
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
func NewButton(x, y, w, h float64, text string) *Button {
	b := &Button{
		X: x, Y: y, Width: w, Height: h,
		Text: text,
	}

	dc := gg.NewContext(int(w), int(h))

	dc.DrawRoundedRectangle(0, 0, w, h, ButtonCornerRadius)
	dc.SetHexColor(gridBorderColor)
	dc.Fill()

	if err := dc.LoadFontFace("client/assets/font.ttf", SubtitleFontSize); err != nil {
		log.Printf("Error loading font: %v", err)
	}
	dc.SetHexColor(gridBackgroundColor)
	dc.DrawStringAnchored(text, w/2, h/2, 0.5, ButtonTextYAnchor)

	b.Image = ebiten.NewImageFromImage(dc.Image())

	return b
}

// Constructor for the room.
func NewRoom(id string, playerCount int) *Room {
	room := &Room{
		Id:          id,
		PlayerCount: playerCount,
	}

	return room
}

// Constructor for the main menu.
func NewMainMenu() *MainMenu {
	menu := &MainMenu{}

	// Center buttons
	centerX := (float64(WindowWidth) - ButtonWidth) / 2

	// Create buttons
	menu.BtnPlay = NewButton(centerX, MainMenuPlayBtnY, ButtonWidth, ButtonHeight, "Play")
	menu.BtnQuit = NewButton(centerX, MainMenuQuitBtnY, ButtonWidth, ButtonHeight, "Quit")

	return menu
}

// Constructor for the play menu.
func NewRoomsMenu() *RoomsMenu {
	menu := &RoomsMenu{}

	// Create buttons
	menu.BtnBack = NewButton(RoomsMenuBtnX, RoomsMenuBackBtnY, ButtonWidth, ButtonHeight, "Back")
	menu.BtnPlayBot = NewButton(RoomsMenuBtnX, RoomsMenuPlayBotBtnY, ButtonWidth, ButtonHeight, "Against Bot")
	menu.BtnCreateRoom = NewButton(RoomsMenuBtnX, RoomsMenuCreateRoomBtnY, ButtonWidth, ButtonHeight, "Create Room")

	return menu
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
	m.BtnBack.Draw(screen)

	for i, room := range m.Rooms {
		room.Draw(screen, i)
	}
}

// Draw the room at the specified index
func (r *Room) Draw(screen *ebiten.Image, index int) {
	// List on the right side
	// TODO: it doenst show for the moment
	listX := float64(WindowWidth/2) + 50.0
	startY := 150.0
	gap := 70.0

	r.JoinBtn = NewButton(listX, startY+float64(index)*gap, ButtonWidth, ButtonHeight, r.Id)
	r.JoinBtn.Draw(screen)
}

// Draw the waiting menu to the screen.
func (waitingMenu *WaitingMenu) Draw(screen *ebiten.Image) {
	screen.DrawImage(WaitingMenuImage, nil)

	w := WheelImage.Bounds().Dx()
	h := WheelImage.Bounds().Dy()
	halfW := float64(w) / 2.0
	halfH := float64(h) / 2.0

	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(-halfW, -halfH)

	op.GeoM.Rotate(waitingMenu.RotationAngle)

	screenCenterX := float64(WindowWidth) / 2.0
	screenCenterY := float64(WindowHeight) / 2.0
	op.GeoM.Translate(screenCenterX, screenCenterY)

	op.ColorScale.Scale(WheelTintRed, WheelTintGreen, WheelTintBlue, 1)

	screen.DrawImage(WheelImage, op)
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
