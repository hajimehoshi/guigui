// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 Hajime Hoshi

package basicwidget

import (
	"image"
	"image/color"
	"log/slog"
	"runtime"
	"slices"
	"strings"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/exp/textinput"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/rivo/uniseg"
	"golang.org/x/text/language"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/internal/clipboard"
)

func isKeyRepeating(key ebiten.Key) bool {
	d := inpututil.KeyPressDuration(key)
	// In the current implementation of text, d == 1 might be skipped especially for backspace key.
	// TODO: Fix this.
	if d == 2 {
		return true
	}
	if d < 24 {
		return false
	}
	return (d-24)%4 == 0
}

func findWordBoundaries(text string, idx int) (start, end int) {
	start = idx
	end = idx

	// invalid text
	if text == "" || idx < 0 || idx >= len(text) {
		return start, end
	}

	gr := uniseg.NewGraphemes(text)
	pos := 0

	for gr.Next() {
		nextPos := pos + len(gr.Bytes())
		if idx < nextPos {
			// found it
			break
		}
		pos = nextPos
	}

	// look for start of word
	state := -1
	tempText := text[:pos]
	words := []string{}

	// process all words before our position to find the current word start
	for len(tempText) > 0 {
		var word string
		word, tempText, state = uniseg.FirstWordInString(tempText, state)
		words = append(words, word)
	}

	// look for end of word
	state = -1
	tempText = text[pos:]
	currentWord := ""

	// get the first word segment from our position
	if len(tempText) > 0 {
		currentWord, _, state = uniseg.FirstWordInString(tempText, state)
	}

	// calculate the start and end positions
	if len(words) > 0 {
		// check word boundary or within word
		lastWord := words[len(words)-1]
		if strings.TrimSpace(lastWord) == "" {
			// at boundary, use previous word if it exists
			if len(words) > 1 {
				start = len(text) - len(tempText) - len(lastWord) - len(words[len(words)-2])
				end = start + len(words[len(words)-2])
			}
		} else {
			// within a word
			start = len(text) - len(tempText) - len(lastWord)
			end = start + len(lastWord) + len(currentWord)
		}
	}

	return start, end
}

type TextFilter func(text string, start, end int) (string, int, int)

func DefaultTextColor(context *guigui.Context) color.Color {
	return Color(context.ColorMode(), ColorTypeBase, 0.1)
}

type Text struct {
	guigui.DefaultWidget

	field textinput.Field

	hAlign      HorizontalAlign
	vAlign      VerticalAlign
	color       color.Color
	transparent float64
	locales     []language.Tag
	scaleMinus1 float64
	bold        bool

	widthSet  bool
	width     int
	heightSet bool
	height    int

	selectable           bool
	editable             bool
	multiline            bool
	selectionDragStart   int
	selectionShiftIndex  int
	dragging             bool
	toAdjustScrollOffset bool
	prevFocused          bool

	clickCount    int
	lastClickTick int64

	filter TextFilter

	cursor        textCursor
	scrollOverlay ScrollOverlay

	temporaryClipboard string

	cachedWidth  int
	cachedHeight int
	initOnce     sync.Once

	onEnterPressed func(text string)
}

func (t *Text) SetOnEnterPressed(f func(text string)) {
	t.onEnterPressed = f
}

func (t *Text) resetCachedSize() {
	t.cachedWidth = -1
	t.cachedHeight = -1
}

func (t *Text) Layout(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	t.initOnce.Do(func() {
		t.resetCachedSize()
	})

	if t.selectable || t.editable {
		p := guigui.Position(t)
		p.X -= cursorWidth(context)
		guigui.SetPosition(&t.cursor, p)
		appender.AppendChildWidget(&t.cursor)
	}

	guigui.SetPosition(&t.scrollOverlay, guigui.Position(t))
	appender.AppendChildWidget(&t.scrollOverlay)
}

func (t *Text) SetSelectable(selectable bool) {
	if t.selectable == selectable {
		return
	}
	t.selectable = selectable
	t.selectionDragStart = -1
	t.selectionShiftIndex = -1
	guigui.RequestRedraw(t)
}

func (t *Text) Text() string {
	return t.field.Text()
}

func (t *Text) SetText(text string) {
	start, end := t.field.Selection()
	start = min(start, len(text))
	end = min(end, len(text))
	t.setTextAndSelection(text, start, end, -1)
	t.resetCachedSize()
}

func (t *Text) SetFilter(filter TextFilter) {
	t.filter = filter
	t.applyFilter()
}

func (t *Text) selectAll() {
	t.setTextAndSelection(t.field.Text(), 0, len(t.field.Text()), -1)
}

func (t *Text) setTextAndSelection(text string, start, end int, shiftIndex int) {
	t.selectionShiftIndex = shiftIndex
	if start > end {
		start, end = end, start
	}

	if s, e := t.field.Selection(); t.field.Text() == text && s == start && e == end {
		return
	}
	t.field.SetTextAndSelection(text, start, end)
	t.toAdjustScrollOffset = true
	guigui.RequestRedraw(t)
	t.resetCachedSize()
}

func (t *Text) SetLocales(locales []language.Tag) {
	if slices.Equal(t.locales, locales) {
		return
	}

	t.locales = append([]language.Tag(nil), locales...)
	guigui.RequestRedraw(t)
	t.resetCachedSize()
}

func (t *Text) SetBold(bold bool) {
	if t.bold == bold {
		return
	}

	t.bold = bold
	guigui.RequestRedraw(t)
	t.resetCachedSize()
}

func (t *Text) SetScale(scale float64) {
	if t.scaleMinus1 == scale-1 {
		return
	}

	t.scaleMinus1 = scale - 1
	guigui.RequestRedraw(t)
	t.resetCachedSize()
}

func (t *Text) SetHorizontalAlign(align HorizontalAlign) {
	if t.hAlign == align {
		return
	}

	t.hAlign = align
	guigui.RequestRedraw(t)
	t.resetCachedSize()
}

func (t *Text) SetVerticalAlign(align VerticalAlign) {
	if t.vAlign == align {
		return
	}

	t.vAlign = align
	guigui.RequestRedraw(t)
	t.resetCachedSize()
}

func (t *Text) SetColor(color color.Color) {
	if equalColor(t.color, color) {
		return
	}

	t.color = color
	guigui.RequestRedraw(t)
}

func (t *Text) SetOpacity(opacity float64) {
	if 1-t.transparent == opacity {
		return
	}

	t.transparent = 1 - opacity
	guigui.RequestRedraw(t)
}

func (t *Text) SetEditable(editable bool) {
	if t.editable == editable {
		return
	}

	if editable {
		t.selectionDragStart = -1
		t.selectionShiftIndex = -1
	}
	t.editable = editable
	guigui.RequestRedraw(t)
}

func (t *Text) SetScrollable(context *guigui.Context, scrollable bool) {
	if scrollable {
		guigui.Show(&t.scrollOverlay)
	} else {
		guigui.Hide(&t.scrollOverlay)
	}
}

func (t *Text) IsMultiline() bool {
	return t.multiline
}

func (t *Text) SetMultiline(multiline bool) {
	if t.multiline == multiline {
		return
	}

	t.multiline = multiline
	guigui.RequestRedraw(t)
	t.resetCachedSize()
}

func (t *Text) textBounds(context *guigui.Context) image.Rectangle {
	offsetX, offsetY := t.scrollOverlay.Offset()

	b := guigui.Bounds(t)

	tw, _ := text.Measure(t.textToDraw(), t.face(context), t.lineHeight(context))
	if b.Dx() < int(tw) {
		b.Max.X = b.Min.X + int(tw)
	}

	th := t.textHeight(context, t.textToDraw())
	switch t.vAlign {
	case VerticalAlignTop:
		b.Max.Y = b.Min.Y + th
	case VerticalAlignMiddle:
		h := b.Dy()
		b.Min.Y += (h - th) / 2
		b.Max.Y = b.Min.Y + th
	case VerticalAlignBottom:
		b.Min.Y = b.Max.Y - th
	}

	b = b.Add(image.Pt(int(offsetX), int(offsetY)))
	return b
}

func (t *Text) face(context *guigui.Context) text.Face {
	size := FontSize(context) * (t.scaleMinus1 + 1)
	weight := text.WeightMedium
	if t.bold {
		weight = text.WeightBold
	}
	locales := append([]language.Tag(nil), t.locales...)
	locales = context.AppendLocales(locales)
	var liga bool
	if !t.selectable && !t.editable {
		liga = true
	}
	return fontFace(size, weight, liga, locales)
}

func (t *Text) lineHeight(context *guigui.Context) float64 {
	return LineHeight(context) * (t.scaleMinus1 + 1)
}

func (t *Text) HandleInput(context *guigui.Context) guigui.HandleInputResult {
	if !t.selectable && !t.editable {
		return guigui.HandleInputResult{}
	}

	textBounds := t.textBounds(context)

	face := t.face(context)
	cursorPosition := image.Pt(ebiten.CursorPosition())
	if t.dragging {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			idx := textIndexFromPosition(textBounds, cursorPosition, t.field.Text(), face, t.lineHeight(context), t.hAlign, t.vAlign)
			if idx < t.selectionDragStart {
				t.setTextAndSelection(t.field.Text(), idx, t.selectionDragStart, -1)
			} else {
				t.setTextAndSelection(t.field.Text(), t.selectionDragStart, idx, -1)
			}
		}
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			t.dragging = false
			t.selectionDragStart = -1
		}
		return guigui.HandleInputByWidget(t)
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if cursorPosition.In(guigui.VisibleBounds(t)) {
			if ebiten.Tick()-t.lastClickTick < int64(ebiten.TPS()/2) {
				t.clickCount++
			} else {
				t.clickCount = 1
			}

			switch t.clickCount {
			case 1:
				t.dragging = true
				idx := textIndexFromPosition(textBounds, cursorPosition, t.field.Text(), face, t.lineHeight(context), t.hAlign, t.vAlign)
				t.selectionDragStart = idx
				if start, end := t.field.Selection(); start != idx || end != idx {
					t.setTextAndSelection(t.field.Text(), idx, idx, -1)
				}
			case 2:
				text := t.field.Text()
				idx := textIndexFromPosition(textBounds, cursorPosition, text, face, t.lineHeight(context), t.hAlign, t.vAlign)
				start, end := findWordBoundaries(text, idx)
				// TODO: `selectionDragEnd` needed to emulate Chrome's behavior.
				t.selectionDragStart = start
				t.setTextAndSelection(text, start, end, -1)
			case 3:
				t.selectAll()
			}

			guigui.Focus(t)
			t.lastClickTick = ebiten.Tick()
			return guigui.HandleInputByWidget(t)
		}
		guigui.Blur(t)
	}

	if !guigui.IsFocused(t) {
		if t.field.IsFocused() {
			t.field.Blur()
			guigui.RequestRedraw(t)
		}
		return guigui.HandleInputResult{}
	}
	t.field.Focus()

	if !t.editable && !t.selectable {
		return guigui.HandleInputResult{}
	}

	start, _ := t.field.Selection()
	var processed bool
	if x, _, bottom, ok := textPosition(textBounds, t.field.Text(), start, face, t.lineHeight(context), t.hAlign, t.vAlign); ok {
		var err error
		processed, err = t.field.HandleInput(int(x), int(bottom))
		if err != nil {
			slog.Error(err.Error())
			processed = false
		}
	}
	if processed {
		guigui.RequestRedraw(t)
		t.adjustScrollOffset(context)
		return guigui.HandleInputByWidget(t)
	}

	// Do not accept key inputs when compositing.
	if _, _, ok := t.field.CompositionSelection(); ok {
		return guigui.HandleInputResult{}
	}

	// TODO: Use WebAPI to detect OS is runtime.GOOS == "js"
	isDarwin := runtime.GOOS == "darwin"

	if t.editable {
		switch {
		case inpututil.IsKeyJustPressed(ebiten.KeyEnter):
			if t.multiline {
				start, end := t.field.Selection()
				text := t.field.Text()[:start] + "\n" + t.field.Text()[end:]
				t.setTextAndSelection(text, start+len("\n"), start+len("\n"), -1)
			}
			t.applyFilter()
			// TODO: This is not reached on browsers. Fix this.
			if t.onEnterPressed != nil {
				t.onEnterPressed(t.field.Text())
			}
		case isKeyRepeating(ebiten.KeyBackspace) ||
			isDarwin && ebiten.IsKeyPressed(ebiten.KeyControl) && isKeyRepeating(ebiten.KeyH):
			start, end := t.field.Selection()
			if start != end {
				text := t.field.Text()[:start] + t.field.Text()[end:]
				t.setTextAndSelection(text, start, start, -1)
			} else if start > 0 {
				text, pos := backspaceOnClusters(t.field.Text(), face, start)
				t.setTextAndSelection(text, pos, pos, -1)
			}
		case !isDarwin && ebiten.IsKeyPressed(ebiten.KeyControl) && isKeyRepeating(ebiten.KeyD) ||
			isDarwin && ebiten.IsKeyPressed(ebiten.KeyControl) && isKeyRepeating(ebiten.KeyD):
			// Delete
			start, end := t.field.Selection()
			if start != end {
				text := t.field.Text()[:start] + t.field.Text()[end:]
				t.setTextAndSelection(text, start, start, -1)
			} else if isDarwin && end < len(t.field.Text()) {
				text, pos := deleteOnClusters(t.field.Text(), face, end)
				t.setTextAndSelection(text, pos, pos, -1)
			}

		case !isDarwin && ebiten.IsKeyPressed(ebiten.KeyControl) && isKeyRepeating(ebiten.KeyX) ||
			isDarwin && ebiten.IsKeyPressed(ebiten.KeyMeta) && isKeyRepeating(ebiten.KeyX):
			// Cut
			start, end := t.field.Selection()
			if start != end {
				if err := clipboard.WriteAll(t.field.Text()[start:end]); err != nil {
					slog.Error(err.Error())
					return guigui.HandleInputResult{}
				}
				text := t.field.Text()[:start] + t.field.Text()[end:]
				t.setTextAndSelection(text, start, start, -1)
			}
		case !isDarwin && ebiten.IsKeyPressed(ebiten.KeyControl) && isKeyRepeating(ebiten.KeyV) ||
			isDarwin && ebiten.IsKeyPressed(ebiten.KeyMeta) && isKeyRepeating(ebiten.KeyV):
			// Paste
			start, end := t.field.Selection()
			ct, err := clipboard.ReadAll()
			if err != nil {
				slog.Error(err.Error())
				return guigui.HandleInputResult{}
			}
			text := t.field.Text()[:start] + ct + t.field.Text()[end:]
			t.setTextAndSelection(text, start+len(ct), start+len(ct), -1)
		}
	}

	switch {
	case isKeyRepeating(ebiten.KeyLeft) ||
		isDarwin && ebiten.IsKeyPressed(ebiten.KeyControl) && isKeyRepeating(ebiten.KeyB):
		start, end := t.field.Selection()
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			if t.selectionShiftIndex == end {
				pos := prevPositionOnClusters(t.field.Text(), face, end)
				t.setTextAndSelection(t.field.Text(), start, pos, pos)
			} else {
				pos := prevPositionOnClusters(t.field.Text(), face, start)
				t.setTextAndSelection(t.field.Text(), pos, end, pos)
			}
		} else {
			if start != end {
				t.setTextAndSelection(t.field.Text(), start, start, -1)
			} else if start > 0 {
				pos := prevPositionOnClusters(t.field.Text(), face, start)
				t.setTextAndSelection(t.field.Text(), pos, pos, -1)
			}
		}
	case isKeyRepeating(ebiten.KeyRight) ||
		isDarwin && ebiten.IsKeyPressed(ebiten.KeyControl) && isKeyRepeating(ebiten.KeyF):
		start, end := t.field.Selection()
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			if t.selectionShiftIndex == start {
				pos := nextPositionOnClusters(t.field.Text(), face, start)
				t.setTextAndSelection(t.field.Text(), pos, end, pos)
			} else {
				pos := nextPositionOnClusters(t.field.Text(), face, end)
				t.setTextAndSelection(t.field.Text(), start, pos, pos)
			}
		} else {
			if start != end {
				t.setTextAndSelection(t.field.Text(), end, end, -1)
			} else if start < len(t.field.Text()) {
				pos := nextPositionOnClusters(t.field.Text(), face, start)
				t.setTextAndSelection(t.field.Text(), pos, pos, -1)
			}
		}
	case isKeyRepeating(ebiten.KeyUp) ||
		isDarwin && ebiten.IsKeyPressed(ebiten.KeyControl) && isKeyRepeating(ebiten.KeyP):
		lh := t.lineHeight(context)
		shift := ebiten.IsKeyPressed(ebiten.KeyShift)
		var moveEnd bool
		start, end := t.field.Selection()
		idx := start
		if shift && t.selectionShiftIndex == end {
			idx = end
			moveEnd = true
		}
		if x, y0, y1, ok := textPosition(textBounds, t.field.Text(), idx, face, lh, t.hAlign, t.vAlign); ok {
			y := (y0+y1)/2 - lh
			idx := textIndexFromPosition(textBounds, image.Pt(int(x), int(y)), t.field.Text(), face, lh, t.hAlign, t.vAlign)
			if shift {
				if moveEnd {
					t.setTextAndSelection(t.field.Text(), start, idx, idx)
				} else {
					t.setTextAndSelection(t.field.Text(), idx, end, idx)
				}
			} else {
				t.setTextAndSelection(t.field.Text(), idx, idx, -1)
			}
		}
	case isKeyRepeating(ebiten.KeyDown) ||
		isDarwin && ebiten.IsKeyPressed(ebiten.KeyControl) && isKeyRepeating(ebiten.KeyN):
		lh := t.lineHeight(context)
		shift := ebiten.IsKeyPressed(ebiten.KeyShift)
		var moveStart bool
		start, end := t.field.Selection()
		idx := end
		if shift && t.selectionShiftIndex == start {
			idx = start
			moveStart = true
		}
		if x, y0, y1, ok := textPosition(textBounds, t.field.Text(), idx, face, lh, t.hAlign, t.vAlign); ok {
			y := (y0+y1)/2 + lh
			idx := textIndexFromPosition(textBounds, image.Pt(int(x), int(y)), t.field.Text(), face, lh, t.hAlign, t.vAlign)
			if shift {
				if moveStart {
					t.setTextAndSelection(t.field.Text(), idx, end, idx)
				} else {
					t.setTextAndSelection(t.field.Text(), start, idx, idx)
				}
			} else {
				t.setTextAndSelection(t.field.Text(), idx, idx, -1)
			}
		}
	case isDarwin && ebiten.IsKeyPressed(ebiten.KeyControl) && isKeyRepeating(ebiten.KeyA):
		idx := 0
		start, end := t.field.Selection()
		if i := strings.LastIndex(t.field.Text()[:start], "\n"); i >= 0 {
			idx = i + 1
		}
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			t.setTextAndSelection(t.field.Text(), idx, end, idx)
		} else {
			t.setTextAndSelection(t.field.Text(), idx, idx, -1)
		}
	case isDarwin && ebiten.IsKeyPressed(ebiten.KeyControl) && isKeyRepeating(ebiten.KeyE):
		idx := len(t.field.Text())
		start, end := t.field.Selection()
		if i := strings.Index(t.field.Text()[end:], "\n"); i >= 0 {
			idx = end + i
		}
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			t.setTextAndSelection(t.field.Text(), start, idx, idx)
		} else {
			t.setTextAndSelection(t.field.Text(), idx, idx, -1)
		}
	case !isDarwin && ebiten.IsKeyPressed(ebiten.KeyControl) && isKeyRepeating(ebiten.KeyA) ||
		isDarwin && ebiten.IsKeyPressed(ebiten.KeyMeta) && isKeyRepeating(ebiten.KeyA):
		t.selectAll()
	case !isDarwin && ebiten.IsKeyPressed(ebiten.KeyControl) && isKeyRepeating(ebiten.KeyC) ||
		isDarwin && ebiten.IsKeyPressed(ebiten.KeyMeta) && isKeyRepeating(ebiten.KeyC):
		// Copy
		start, end := t.field.Selection()
		if start != end {
			if err := clipboard.WriteAll(t.field.Text()[start:end]); err != nil {
				slog.Error(err.Error())
				return guigui.HandleInputResult{}
			}
		}
	case isDarwin && ebiten.IsKeyPressed(ebiten.KeyControl) && isKeyRepeating(ebiten.KeyK):
		// 'Kill' the text after the cursor or the selection.
		start, end := t.field.Selection()
		if start == end {
			end = strings.Index(t.field.Text()[start:], "\n")
			if end < 0 {
				end = len(t.field.Text())
			} else {
				end += start
			}
		}
		t.temporaryClipboard = t.field.Text()[start:end]
		text := t.field.Text()[:start] + t.field.Text()[end:]
		t.setTextAndSelection(text, start, start, -1)
	case isDarwin && ebiten.IsKeyPressed(ebiten.KeyControl) && isKeyRepeating(ebiten.KeyY):
		// 'Yank' the killed text.
		if t.temporaryClipboard != "" {
			start, _ := t.field.Selection()
			text := t.field.Text()[:start] + t.temporaryClipboard + t.field.Text()[start:]
			t.setTextAndSelection(text, start+len(t.temporaryClipboard), start+len(t.temporaryClipboard), -1)
		}
	}

	return guigui.HandleInputResult{}
}

func (t *Text) adjustScrollOffset(context *guigui.Context) {
	t.updateContentSize(context)

	start, end, ok := t.selectionToDraw()
	if !ok {
		return
	}
	text := t.textToDraw()

	tb := t.textBounds(context)
	face := t.face(context)
	bounds := guigui.Bounds(t)
	if x, _, y, ok := textPosition(tb, text, end, face, t.lineHeight(context), t.hAlign, t.vAlign); ok {
		var dx, dy float64
		if max := float64(bounds.Max.X); x > max {
			dx = max - x
		}
		if max := float64(bounds.Max.Y); y > max {
			dy = max - y
		}
		t.scrollOverlay.SetOffsetByDelta(dx, dy)
	}
	if x, y, _, ok := textPosition(tb, text, start, face, t.lineHeight(context), t.hAlign, t.vAlign); ok {
		var dx, dy float64
		if min := float64(bounds.Min.X); x < min {
			dx = min - x
		}
		if min := float64(bounds.Min.Y); y < min {
			dy = min - y
		}
		t.scrollOverlay.SetOffsetByDelta(dx, dy)
	}
}

func (t *Text) textToDraw() string {
	return t.field.TextForRendering()
}

func (t *Text) selectionToDraw() (start, end int, ok bool) {
	s, e := t.field.Selection()
	if !t.editable {
		return s, e, true
	}
	if !guigui.IsFocused(t) {
		return s, e, true
	}
	cs, ce, ok := t.field.CompositionSelection()
	if !ok {
		return s, e, true
	}
	// When cs == ce, the composition already started but any conversion is not done yet.
	// In this case, put the cursor at the end of the composition.
	// TODO: This behavior might be macOS specific. Investigate this.
	if cs == ce {
		return s + ce, e + ce, true
	}
	return 0, 0, false
}

func (t *Text) compositionSelectionToDraw() (uStart, cStart, cEnd, uEnd int, ok bool) {
	if !t.editable {
		return 0, 0, 0, 0, false
	}
	if !guigui.IsFocused(t) {
		return 0, 0, 0, 0, false
	}
	s, _ := t.field.Selection()
	cs, ce, ok := t.field.CompositionSelection()
	if !ok {
		return 0, 0, 0, 0, false
	}
	// When cs == ce, the composition already started but any conversion is not done yet.
	// In this case, assume the entire region is the composition.
	// TODO: This behavior might be macOS specific. Investigate this.
	l := t.field.UncommittedTextLengthInBytes()
	if cs == ce {
		return s, s, s + l, s + l, true
	}
	return s, s + cs, s + ce, s + l, true
}

func (t *Text) Update(context *guigui.Context) error {
	guigui.Hide(&t.scrollOverlay)

	if !t.prevFocused && guigui.IsFocused(t) {
		t.field.Focus()
		t.cursor.resetCounter()
		start, end := t.field.Selection()
		if start < 0 || end < 0 {
			t.selectAll()
		}
	} else if t.prevFocused && !guigui.IsFocused(t) {
		t.applyFilter()
	}

	if t.toAdjustScrollOffset && !guigui.VisibleBounds(t).Empty() {
		t.adjustScrollOffset(context)
		t.toAdjustScrollOffset = false
	}

	t.prevFocused = guigui.IsFocused(t)

	return nil
}

func (t *Text) applyFilter() {
	if t.filter != nil {
		start, end := t.field.Selection()
		text, start, end := t.filter(t.field.Text(), start, end)
		t.setTextAndSelection(text, start, end, -1)
	}
}

func (t *Text) updateContentSize(context *guigui.Context) {
	w, h := t.TextSize(context)
	t.scrollOverlay.SetContentSize(w, h)
}

func (t *Text) Draw(context *guigui.Context, dst *ebiten.Image) {
	textBounds := t.textBounds(context)
	if !textBounds.Overlaps(guigui.VisibleBounds(t)) {
		return
	}

	text := t.textToDraw()
	face := t.face(context)

	if start, end, ok := t.selectionToDraw(); ok {
		var tailIndices []int
		for i, r := range text[start:end] {
			if r != '\n' {
				continue
			}
			tailIndices = append(tailIndices, start+i)
		}
		tailIndices = append(tailIndices, end)

		headIdx := start
		for _, idx := range tailIndices {
			x0, top0, bottom0, ok0 := textPosition(textBounds, text, headIdx, face, t.lineHeight(context), t.hAlign, t.vAlign)
			x1, _, _, ok1 := textPosition(textBounds, text, idx, face, t.lineHeight(context), t.hAlign, t.vAlign)
			if ok0 && ok1 {
				x := float32(x0)
				y := float32(top0)
				width := float32(x1 - x0)
				height := float32(bottom0 - top0)
				vector.DrawFilledRect(dst, x, y, width, height, Color(context.ColorMode(), ColorTypeAccent, 0.8), false)
			}
			headIdx = idx + 1
		}
	}

	if uStart, cStart, cEnd, uEnd, ok := t.compositionSelectionToDraw(); ok {
		// Assume that the composition is always in the same line.
		if strings.Contains(text[uStart:uEnd], "\n") {
			slog.Error("composition text must not contain '\\n'")
		}
		{
			x0, _, bottom0, ok0 := textPosition(textBounds, text, uStart, face, t.lineHeight(context), t.hAlign, t.vAlign)
			x1, _, _, ok1 := textPosition(textBounds, text, uEnd, face, t.lineHeight(context), t.hAlign, t.vAlign)
			if ok0 && ok1 {
				x := float32(x0)
				y := float32(bottom0) - float32(cursorWidth(context))
				w := float32(x1 - x0)
				h := float32(cursorWidth(context))
				vector.DrawFilledRect(dst, x, y, w, h, Color(context.ColorMode(), ColorTypeAccent, 0.8), false)
			}
		}
		{
			x0, _, bottom0, ok0 := textPosition(textBounds, text, cStart, face, t.lineHeight(context), t.hAlign, t.vAlign)
			x1, _, _, ok1 := textPosition(textBounds, text, cEnd, face, t.lineHeight(context), t.hAlign, t.vAlign)
			if ok0 && ok1 {
				x := float32(x0)
				y := float32(bottom0) - float32(cursorWidth(context))
				w := float32(x1 - x0)
				h := float32(cursorWidth(context))
				vector.DrawFilledRect(dst, x, y, w, h, Color(context.ColorMode(), ColorTypeAccent, 0.4), false)
			}
		}
	}

	var clr color.RGBA
	if t.color != nil {
		clr = color.RGBAModel.Convert(t.color).(color.RGBA)
	} else {
		clr = color.RGBAModel.Convert(DefaultTextColor(context)).(color.RGBA)
	}
	if t.transparent > 0 {
		opacity := 1 - t.transparent
		clr = color.RGBA{
			R: byte(float64(clr.R) * opacity),
			G: byte(float64(clr.G) * opacity),
			B: byte(float64(clr.B) * opacity),
			A: byte(float64(clr.A) * opacity),
		}
	}
	drawText(textBounds, dst, text, face, t.lineHeight(context), t.hAlign, t.vAlign, clr)
}

func (t *Text) Size(context *guigui.Context) (int, int) {
	if t.cachedWidth >= 0 && t.cachedHeight >= 0 {
		return t.cachedWidth, t.cachedHeight
	}

	w, h := t.width, t.height
	if !t.widthSet || !t.heightSet {
		tw, th := t.TextSize(context)
		if !t.widthSet {
			w = tw
		}
		if !t.heightSet {
			h = th
		}
	}
	t.cachedWidth = w
	t.cachedHeight = h
	return w, h
}

func (t *Text) TextSize(context *guigui.Context) (int, int) {
	w, _ := text.Measure(t.textToDraw(), t.face(context), t.lineHeight(context))
	w *= t.scaleMinus1 + 1
	h := t.textHeight(context, t.textToDraw())
	return int(w), h
}

func (t *Text) textHeight(context *guigui.Context, str string) int {
	// The text is already shifted by (lineHeight - (m.HAscent + m.Descent)) / 2.
	return int(t.lineHeight(context) * float64(strings.Count(str, "\n")+1))
}

func (t *Text) SetSize(width, height int) {
	t.widthSet = true
	t.heightSet = true
	t.width = width
	t.height = height
}

func (t *Text) SetWidth(width int) {
	t.widthSet = true
	t.width = width
}

func (t *Text) SetHeight(height int) {
	t.heightSet = true
	t.height = height
}

func (t *Text) ResetSize() {
	t.widthSet = false
	t.heightSet = false
}

func (t *Text) CursorShape(context *guigui.Context) (ebiten.CursorShapeType, bool) {
	if t.selectable || t.editable {
		return ebiten.CursorShapeText, true
	}
	return 0, false
}

func (t *Text) cursorPosition(context *guigui.Context) (x, top, bottom float64, ok bool) {
	if !guigui.IsFocused(t) {
		return 0, 0, 0, false
	}
	if !t.editable {
		return 0, 0, 0, false
	}
	start, end := t.field.Selection()
	if start < 0 {
		return 0, 0, 0, false
	}
	if end < 0 {
		return 0, 0, 0, false
	}

	textBounds := t.textBounds(context)
	if !textBounds.Overlaps(guigui.VisibleBounds(t)) {
		return 0, 0, 0, false
	}

	_, e, ok := t.selectionToDraw()
	if !ok {
		return 0, 0, 0, false
	}

	text := t.textToDraw()
	face := t.face(context)
	return textPosition(textBounds, text, e, face, t.lineHeight(context), t.hAlign, t.vAlign)
}

func cursorWidth(context *guigui.Context) int {
	return int(2 * context.Scale())
}

func (t *Text) cursorBounds(context *guigui.Context) image.Rectangle {
	x, top, bottom, ok := t.cursorPosition(context)
	if !ok {
		return image.Rectangle{}
	}
	w := cursorWidth(context)
	return image.Rect(int(x)-w/2, int(top), int(x)+w/2, int(bottom))
}

type textCursor struct {
	guigui.DefaultWidget

	counter    int
	prevShown  bool
	prevX      float64
	prevTop    float64
	prevBottom float64
	prevOK     bool
}

func (t *textCursor) resetCounter() {
	t.counter = 0
}

func (t *textCursor) Update(context *guigui.Context) error {
	text := guigui.Parent(t).(*Text)
	x, top, bottom, ok := text.cursorPosition(context)
	if t.prevX != x || t.prevTop != top || t.prevBottom != bottom || t.prevOK != ok {
		t.resetCounter()
	}
	t.prevX = x
	t.prevTop = top
	t.prevBottom = bottom
	t.prevOK = ok

	t.counter++
	if r := t.shouldRenderCursor(context, text); t.prevShown != r {
		t.prevShown = r
		// TODO: This is not efficient. Improve this.
		guigui.RequestRedraw(t)
	}
	return nil
}

func (t *textCursor) shouldRenderCursor(context *guigui.Context, text *Text) bool {
	offset := ebiten.TPS() / 2
	if t.counter > offset && (t.counter-offset)%ebiten.TPS() >= ebiten.TPS()/2 {
		return false
	}
	if _, _, _, ok := text.cursorPosition(context); !ok {
		return false
	}
	s, e, ok := text.selectionToDraw()
	if !ok {
		return false
	}
	if s != e {
		return false
	}
	return true
}

func (t *textCursor) Draw(context *guigui.Context, dst *ebiten.Image) {
	textWidget := guigui.Parent(t)
	if !t.shouldRenderCursor(context, textWidget.(*Text)) {
		return
	}
	b := textWidget.(*Text).cursorBounds(context)
	vector.DrawFilledRect(dst, float32(b.Min.X), float32(b.Min.Y), float32(b.Dx()), float32(b.Dy()), Color(context.ColorMode(), ColorTypeAccent, 0.4), false)
}

func (t *textCursor) IsPopup() bool {
	return true
}

func (t *textCursor) Size(context *guigui.Context) (int, int) {
	w, h := guigui.Parent(t).Size(context)
	return w + 2*cursorWidth(context), h
}
