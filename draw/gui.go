package draw

import (
	"strconv"

	"maze/config"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type ElementBounds struct {
	width  int32
	height int32
}

type GuiElement interface {
	Render(xPos, yPos float32)
	Bounds() ElementBounds
}

type baseElement struct {
	width  int32
	height int32
}

func (b *baseElement) Bounds() ElementBounds {
	return ElementBounds{
		width:  b.width,
		height: b.height,
	}
}

type Button struct {
	baseElement
	text    *string
	onClick func()
}

func (b *Button) Render(xPos, yPos float32) {
	clicked := gui.Button(
		rl.NewRectangle(
			xPos,
			yPos,
			float32(b.width),
			float32(b.height)),
		*b.text,
	)

	if clicked {
		b.onClick()
	}
}

type Slider struct {
	baseElement
	textLeft  string
	textRight string
	value     float32
}

func (s *Slider) Render(xPos, yPos float32) {
	interval := gui.Slider(
		rl.NewRectangle(
			xPos+50,
			yPos,
			// Account for the labels also taking up space
			float32(s.width-100),
			float32(s.height),
		),
		s.textLeft,
		s.textRight,
		s.value,
		0,
		300,
	)

	s.value = interval
}

type Toggle struct {
	baseElement
	text   string
	active *bool
}

func (t *Toggle) Render(xPos, yPos float32) {
	active := gui.Toggle(
		rl.NewRectangle(
			xPos,
			yPos,
			float32(t.width),
			float32(t.height),
		),
		t.text,
		*t.active,
	)

	*t.active = active
}

type Dropdown struct {
	baseElement
	text           string
	active         *int32
	previousActive int32
	editMode       *bool
	onChange       func()
}

func (s *Dropdown) Render(xPos, yPos float32) {
	editMode := gui.DropdownBox(
		rl.NewRectangle(
			xPos,
			yPos,
			float32(s.width),
			float32(s.height),
		),
		s.text,
		s.active,
		*s.editMode,
	)

	if editMode {
		*s.editMode = !*s.editMode
	}

	if *s.active != s.previousActive {
		s.previousActive = *s.active
		s.onChange()
	}
}

func DrawGui(elements []GuiElement, stats Stats) {
	textColor := rl.NewColor(0, 0, 56, 125)
	buttonsContainerHeight := int32(25)
	padding := config.Padding
	xPos := padding
	yPos := padding

	rl.DrawRectangle(
		xPos,
		yPos,
		config.MenuBarWidth,
		config.MenuBarHeight-padding,
		backgroundColor,
	)

	rl.DrawRectangleLinesEx(
		rl.NewRectangle(
			0,
			0,
			float32(config.Width),
			float32(buttonsContainerHeight+(padding*2)),
		),
		2,
		borderColor,
	)

	rl.DrawRectangleLinesEx(
		rl.NewRectangle(
			0,
			float32(buttonsContainerHeight+padding),
			float32(config.Width),
			float32(buttonsContainerHeight),
		),
		2,
		borderColor,
	)

	runTimeText := "Total run time: " + formatTime(stats.runTimeMicroseconds) + "ms"
	totalStepsText := "Total steps: " + strconv.Itoa(stats.totalSteps)
	solutionStepsText := "Solution steps: " + strconv.Itoa(stats.solutionSteps)

	statsOffset := buttonsContainerHeight + padding
	statTextOffset := statsOffset + 8

	rl.DrawTextEx(
		font,
		totalStepsText,
		rl.Vector2{
			X: float32(padding + 5),
			Y: float32(statTextOffset),
		}, 14, 1, textColor)

	rl.DrawRectangle(
		padding+140,
		statsOffset,
		2,
		buttonsContainerHeight,
		borderColor,
	)

	rl.DrawTextEx(
		font,
		solutionStepsText,
		rl.Vector2{
			X: 150,
			Y: float32(statTextOffset),
		}, 14, 1, textColor)

	rl.DrawRectangle(
		300,
		statsOffset,
		2,
		buttonsContainerHeight,
		borderColor,
	)

	rl.DrawTextEx(
		font,
		runTimeText,
		rl.Vector2{
			X: 308,
			Y: float32(statTextOffset),
		}, 14, 1, textColor)

	rl.DrawRectangle(
		500,
		statsOffset,
		2,
		buttonsContainerHeight,
		borderColor,
	)

	xOffset := xPos + padding
	yOffset := yPos + padding

	for _, element := range elements {
		element.Render(
			float32(xOffset),
			float32(yOffset),
		)
		xOffset += element.Bounds().width + padding
	}
}
