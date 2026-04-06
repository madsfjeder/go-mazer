package draw

import (
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
