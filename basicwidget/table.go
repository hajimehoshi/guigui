// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 Björn De Meyer, Hajime JHshi

package basicwidget

import (
	"image"
	"image/color"
	"slices"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/hajimehoshi/guigui"
)

// Cell is a cell in a column of a Table.
type Cell struct {
	Content    guigui.Widget // Content is the content of the cell. A cell can contain any widget.
	Selectable bool          // Selectable is whether or not the cell can be selected.
	Editable   bool          // Editable is whether or not the cell can be edited.
	Tag        any           // Tag is freely usable data for the cell.
}

func DefaultActiveCellTextColor(context *guigui.Context) color.Color {
	return Color(context.ColorMode(), ColorTypeBase, 1)
}

func DefaultDisabledCellTextColor(context *guigui.Context) color.Color {
	return Color(context.ColorMode(), ColorTypeBase, 0.5)
}

// Column is a column in a Table. Tables have a fixed amount of rows.
// Tables are column-first.
type Column struct {
	guigui.DefaultWidget

	columnFrame columnFrame
	caption     Text

	cells                  []Cell
	selectedCellIndexPlus1 int
	hoveredCellIndexPlus1  int
	showCellBorders        bool
	lastSelectingcellTime  time.Time

	indexToJumpPlus1        int
	dropSrcIndexPlus1       int
	dropDstIndexPlus1       int
	pressStartX             int
	pressStartY             int
	startPressingIndexPlus1 int
	startPressingLeft       bool

	widthSet            bool
	heightSet           bool
	width               int
	height              int
	cachedDefaultWidth  int
	cachedDefaultHeight int

	onCellSelected func(index int)
	Selectable     bool // Selectable is whether or not the column can be selected.
}

func CellPadding(context *guigui.Context) int {
	return UnitSize(context) / 4
}

func ColumnCornerRadius(context *guigui.Context) int {
	return UnitSize(context) / 16
}

func (c *Column) SetOnCellSelected(f func(index int)) {
	c.onCellSelected = f
}

func (c *Column) Layout(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	guigui.SetPosition(&c.columnFrame, guigui.Position(c))
	appender.AppendChildWidget(&c.columnFrame)

	p := guigui.Position(c)
	p.X += ColumnCornerRadius(context) + CellPadding(context)
	p.Y += ColumnCornerRadius(context)

	guigui.SetPosition(&c.caption, p)
	appender.AppendChildWidget(&c.caption)
	_, h := c.caption.Size(context)
	p.Y += h
	p.X += ColumnCornerRadius(context) + CellPadding(context)
	p.Y += ColumnCornerRadius(context)

	for _, cell := range c.cells {
		guigui.SetPosition(cell.Content, p)
		appender.AppendChildWidget(cell.Content)
		_, h := cell.Content.Size(context)
		p.Y += h
	}

	p = guigui.Position(c)
}

func (c *Column) selectedCell() (Cell, bool) {
	idx := c.selectedCellIndex()
	if idx < 0 || idx >= len(c.cells) {
		return Cell{}, false
	}
	return c.cells[idx], true
}

func (c *Column) CellAt(index int) (Cell, bool) {
	if index < 0 || index >= len(c.cells) {
		return Cell{}, false
	}
	return c.cells[index], true
}

func (c *Column) selectedCellIndex() int {
	return c.selectedCellIndexPlus1 - 1
}

func (c *Column) hoveredCellIndex() int {
	return c.hoveredCellIndexPlus1 - 1
}

func (c *Column) SetCells(cells []Cell) {
	c.cells = make([]Cell, len(cells))
	copy(c.cells, cells)
	c.cachedDefaultWidth = 0
	c.cachedDefaultHeight = 0
}

func (c *Column) SetCaption(caption string) {
	c.caption.SetText(caption)
	c.caption.SetBold(caption != "")
}

func (c *Column) Caption() string {
	return c.caption.Text()
}

func (c *Column) SetCell(cell Cell, index int) {
	c.cells[index] = cell
}

func (c *Column) AddCell(cell Cell, index int) {
	c.cells = slices.Insert(c.cells, index, cell)
}

func (c *Column) RemoveCell(index int) {
	c.cells = slices.Delete(c.cells, index, index+1)
}

func (c *Column) Movecell(from int, to int) {
	moveItemInSlice(c.cells, from, 1, to)
}

func (c *Column) SetSelectedCellIndex(index int) {
	if index < 0 || index >= len(c.cells) {
		index = -1
	}
	if c.selectedCellIndex() != index {
		c.selectedCellIndexPlus1 = index + 1
		guigui.RequestRedraw(c)
	}
	if c.onCellSelected != nil {
		c.onCellSelected(index)
	}
}

func (c *Column) JumpTocellIndex(index int) {
	if index < 0 || index >= len(c.cells) {
		return
	}
	c.indexToJumpPlus1 = index + 1
}

func (c *Column) setHoveredCellIndex(index int) {
	if index < 0 || index >= len(c.cells) {
		index = -1
	}
	if c.hoveredCellIndex() == index {
		return
	}
	c.hoveredCellIndexPlus1 = index + 1
	if c.isHoveringVisible() {
		guigui.RequestRedraw(c)
	}
}

func (c *Column) ShowCellBorders(show bool) {
	if c.showCellBorders == show {
		return
	}
	c.showCellBorders = true
	guigui.RequestRedraw(c)
}

func (c *Column) isHoveringVisible() bool {
	return true
}

func (c *Column) calcDropDstIndex(context *guigui.Context) int {
	_, y := ebiten.CursorPosition()
	for i := range c.cells {
		if r := c.cellRect(context, i); y < (r.Min.Y+r.Max.Y)/2 {
			return i
		}
	}
	return len(c.cells)
}

func (c *Column) HandleInput(context *guigui.Context) guigui.HandleInputResult {
	if x, y := ebiten.CursorPosition(); image.Pt(x, y).In(guigui.VisibleBounds(c)) {
		y -= RoundedCornerRadius(context)
		y -= guigui.Position(c).Y

		index := -1
		var cy int
		for i, cell := range c.cells {
			_, h := cell.Content.Size(context)
			if cy <= y && y < cy+h {
				index = i
				break
			}
			cy += h
		}
		c.setHoveredCellIndex(index)
		if index >= 0 && index < len(c.cells) {
			left := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
			right := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight)

			switch {
			case left || right:
				if !c.cells[index].Selectable {
					return guigui.HandleInputByWidget(c)
				}

				wasFocused := guigui.IsFocused(c)
				guigui.Focus(c)
				if c.selectedCellIndex() != index || !wasFocused {
					c.SetSelectedCellIndex(index)
					c.lastSelectingcellTime = time.Now()
				}
				c.pressStartX = x
				c.pressStartY = y
				if right { // TODO: send event
				}
				c.startPressingIndexPlus1 = index + 1
				c.startPressingLeft = left

			case ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft):
				// TODO: send event

			case inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft):
				if c.selectedCellIndex() == index && c.startPressingLeft && time.Since(c.lastSelectingcellTime) > 400*time.Millisecond {
					// TODO: send event
				}
				c.pressStartX = 0
				c.pressStartY = 0
				c.startPressingIndexPlus1 = 0
				c.startPressingLeft = false
			}

			return guigui.HandleInputByWidget(c)
		}
		c.dropSrcIndexPlus1 = 0
		c.pressStartX = 0
		c.pressStartY = 0
	} else {
		c.setHoveredCellIndex(-1)
	}

	return guigui.HandleInputResult{}
}

func (c *Column) Update(context *guigui.Context) error {
	idx := c.indexToJumpPlus1 - 1
	if idx >= 0 {
		c.indexToJumpPlus1 = 0
	}

	return nil
}

func (c *Column) cellYFromIndex(context *guigui.Context, index int) int {
	y := RoundedCornerRadius(context)
	_, ch := c.caption.Size(context)
	y += ch
	for i, cell := range c.cells {
		if i == index {
			break
		}
		_, wh := cell.Content.Size(context)
		y += wh
	}
	return y
}

func (c *Column) cellRect(context *guigui.Context, index int) image.Rectangle {
	p := guigui.Position(c)
	w, h := c.Size(context)
	b := image.Rectangle{
		Min: p,
		Max: p.Add(image.Pt(w, h)),
	}
	padding := CellPadding(context)
	b.Min.X += RoundedCornerRadius(context) + padding
	b.Max.X -= RoundedCornerRadius(context) + padding
	b.Min.Y += c.cellYFromIndex(context, index)

	_, ih := c.cells[index].Content.Size(context)
	b.Max.Y = b.Min.Y + ih
	return b
}

func (c *Column) selectedCellColor(context *guigui.Context) color.Color {
	if c.selectedCellIndex() < 0 || c.selectedCellIndex() >= len(c.cells) {
		return nil
	}
	if guigui.IsFocused(c) {
		return Color(context.ColorMode(), ColorTypeAccent, 0.5)
	}
	return Color(context.ColorMode(), ColorTypeBase, 0.8)
}

func (c *Column) Draw(context *guigui.Context, dst *ebiten.Image) {
	clr := Color(context.ColorMode(), ColorTypeBase, 1)

	p := guigui.Position(c)
	w, h := c.Size(context)
	bounds := image.Rectangle{
		Min: p,
		Max: p.Add(image.Pt(w, h)),
	}
	DrawRoundedRect(context, dst, bounds, clr, RoundedCornerRadius(context))

	// Draw cell borders.
	if c.showCellBorders && len(c.cells) > 0 {
		p := guigui.Position(c)
		w, _ := c.Size(context)
		y := float32(p.Y) + float32(ColumnCornerRadius(context))
		for i, cell := range c.cells {
			_, ih := cell.Content.Size(context)
			y += float32(ih)
			if i == c.selectedCellIndex() || i+1 == c.selectedCellIndex() {
				continue
			}
			if i == len(c.cells)-1 {
				continue
			}
			x0 := p.X + ColumnCornerRadius(context)
			x1 := p.X + w - ColumnCornerRadius(context)
			width := 1 * float32(context.Scale())
			clr := Color(context.ColorMode(), ColorTypeBase, 0.5)
			vector.StrokeLine(dst, float32(x0), y, float32(x1), y, width, clr, false)
		}
	}

	if clr := c.selectedCellColor(context); clr != nil && c.selectedCellIndex() >= 0 && c.selectedCellIndex() < len(c.cells) {
		r := c.cellRect(context, c.selectedCellIndex())
		r.Min.X -= RoundedCornerRadius(context)
		r.Max.X += RoundedCornerRadius(context)
		if r.Overlaps(guigui.VisibleBounds(c)) {
			DrawRoundedRect(context, dst, r, clr, ColumnCornerRadius(context))
		}
	}

	if c.isHoveringVisible() && c.hoveredCellIndex() >= 0 && c.hoveredCellIndex() < len(c.cells) && c.cells[c.hoveredCellIndex()].Selectable {
		r := c.cellRect(context, c.hoveredCellIndex())
		r.Min.X -= RoundedCornerRadius(context)
		r.Max.X += RoundedCornerRadius(context)
		if r.Overlaps(guigui.VisibleBounds(c)) {
			clr := Color(context.ColorMode(), ColorTypeBase, 0.9)
			DrawRoundedRect(context, dst, r, clr, ColumnCornerRadius(context))
		}
	}
}

func (c *Column) defaultWidth(context *guigui.Context) int {
	if c.cachedDefaultWidth > 0 {
		return c.cachedDefaultWidth
	}
	var w int
	w, _ = c.caption.Size(context)
	for _, cell := range c.cells {
		iw, _ := cell.Content.Size(context)
		w = max(w, iw)
	}
	w += 2*ColumnCornerRadius(context) + 2*CellPadding(context)
	c.cachedDefaultWidth = w
	return w
}

func (c *Column) defaultHeight(context *guigui.Context) int {
	if c.cachedDefaultHeight > 0 {
		return c.cachedDefaultHeight
	}

	var h int
	_, h = c.caption.Size(context)
	h += ColumnCornerRadius(context)
	for _, w := range c.cells {
		_, wh := w.Content.Size(context)
		h += wh
	}
	h += ColumnCornerRadius(context)
	c.cachedDefaultHeight = h
	return h
}

func (c *Column) Size(context *guigui.Context) (int, int) {
	var w, h int
	if c.widthSet {
		w = c.width
	} else {
		w = c.defaultWidth(context)
	}
	if c.heightSet {
		h = c.height
	} else {
		h = c.defaultHeight(context)
	}
	return w, h
}

func (c *Column) SetSize(width, height int) {
	c.width = width
	c.widthSet = true
	c.height = height
	c.heightSet = true
}

func (c *Column) SetWidth(width int) {
	c.width = width
	c.widthSet = true
}

func (c *Column) SetHeight(height int) {
	c.height = height
	c.heightSet = true
}

func (c *Column) ResetWidth() {
	c.widthSet = false
	c.width = 0
}

func (c *Column) ResetHeight() {
	c.heightSet = false
	c.height = 0
}

type columnFrame struct {
	guigui.DefaultWidget
}

func (c *columnFrame) Draw(context *guigui.Context, dst *ebiten.Image) {
	border := RoundedRectBorderTypeInset
	p := guigui.Position(c)
	w, h := c.Size(context)
	bounds := image.Rectangle{
		Min: p,
		Max: p.Add(image.Pt(w, h)),
	}
	clr := Color2(context.ColorMode(), ColorTypeBase, 0.7, 0)
	borderWidth := float32(1 * context.Scale())
	DrawRoundedRectBorder(context, dst, bounds, clr, ColumnCornerRadius(context), borderWidth, border)
}

func (c *columnFrame) Size(context *guigui.Context) (int, int) {
	return guigui.Parent(c).Size(context)
}

// Table is a simple static table of columns with cells.
type Table struct {
	guigui.DefaultWidget

	tableFrame    columnFrame
	scrollOverlay ScrollOverlay

	columns                  []Column
	selectedColumnIndexPlus1 int
	hoveredColumnIndexPlus1  int
	showColumnBorders        bool
	lastSelectingColumnTime  time.Time

	indexToJumpPlus1        int
	dropSrcIndexPlus1       int
	dropDstIndexPlus1       int
	pressStartX             int
	pressStartY             int
	startPressingIndexPlus1 int
	startPressingLeft       bool

	widthSet            bool
	heightSet           bool
	width               int
	height              int
	cachedDefaultWidth  int
	cachedDefaultHeight int

	onColumnSelected func(index int)
}

type tableFrame struct {
	guigui.DefaultWidget
}

func (t *tableFrame) Draw(context *guigui.Context, dst *ebiten.Image) {
	border := RoundedRectBorderTypeInset
	p := guigui.Position(t)
	w, h := t.Size(context)
	bounds := image.Rectangle{
		Min: p,
		Max: p.Add(image.Pt(w, h)),
	}
	clr := Color2(context.ColorMode(), ColorTypeBase, 0.7, 0)
	borderWidth := float32(1 * context.Scale())
	DrawRoundedRectBorder(context, dst, bounds, clr, RoundedCornerRadius(context), borderWidth, border)
}

func (t *tableFrame) Size(context *guigui.Context) (int, int) {
	return guigui.Parent(t).Size(context)
}

func ColumnPadding(context *guigui.Context) int {
	return UnitSize(context) / 4
}

func (t *Table) SetOnColumnSelected(f func(index int)) {
	t.onColumnSelected = f
}

func (t *Table) Layout(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	guigui.SetPosition(&t.tableFrame, guigui.Position(t))
	appender.AppendChildWidget(&t.tableFrame)

	_, offsetY := t.scrollOverlay.Offset()
	p := guigui.Position(t)
	p.X += RoundedCornerRadius(context) + ColumnPadding(context)
	p.Y += RoundedCornerRadius(context) + int(offsetY)

	for _, column := range t.columns {
		guigui.SetPosition(&column, p)
		appender.AppendChildWidget(&column)
		w, _ := column.Size(context)
		p.X += w
	}

	p = guigui.Position(t)
	guigui.SetPosition(&t.scrollOverlay, p)
	appender.AppendChildWidget(&t.scrollOverlay)
}

func (t *Table) selectedColumn() (Column, bool) {
	idx := t.selectedColumnIndex()
	if idx < 0 || idx >= len(t.columns) {
		return Column{}, false
	}
	return t.columns[idx], true
}

func (t *Table) ColumnAt(index int) (Column, bool) {
	if index < 0 || index >= len(t.columns) {
		return Column{}, false
	}
	return t.columns[index], true
}

func (t *Table) selectedColumnIndex() int {
	return t.selectedColumnIndexPlus1 - 1
}

func (t *Table) hoveredColumnIndex() int {
	return t.hoveredColumnIndexPlus1 - 1
}

func (t *Table) SetColumns(columns []Column) {
	t.columns = make([]Column, len(columns))
	copy(t.columns, columns)
	t.cachedDefaultWidth = 0
	t.cachedDefaultHeight = 0
}

func (t *Table) SetColumn(column Column, index int) {
	t.columns[index] = column
}

func (t *Table) AddColumn(column Column, index int) {
	t.columns = slices.Insert(t.columns, index, column)
}

func (t *Table) RemoveColumn(index int) {
	t.columns = slices.Delete(t.columns, index, index+1)
}

func (t *Table) MoveColumn(from int, to int) {
	moveItemInSlice(t.columns, from, 1, to)
}

func (t *Table) SetSelectedColumnIndex(index int) {
	if index < 0 || index >= len(t.columns) {
		index = -1
	}
	if t.selectedColumnIndex() != index {
		t.selectedColumnIndexPlus1 = index + 1
		guigui.RequestRedraw(t)
	}
	if t.onColumnSelected != nil {
		t.onColumnSelected(index)
	}
}

func (t *Table) JumpToRowIndex(index int) {
	if index < 0 || index >= len(t.columns) {
		return
	}
	t.indexToJumpPlus1 = index + 1
}

func (t *Table) setHoveredColumnIndex(index int) {
	if index < 0 || index >= len(t.columns) {
		index = -1
	}
	if t.hoveredColumnIndex() == index {
		return
	}
	t.hoveredColumnIndexPlus1 = index + 1
	if t.isHoveringVisible() {
		guigui.RequestRedraw(t)
	}
}

func (t *Table) ShowColumnBorders(show bool) {
	if t.showColumnBorders == show {
		return
	}
	t.showColumnBorders = true
	guigui.RequestRedraw(t)
}

func (t *Table) isHoveringVisible() bool {
	return true
}

func (t *Table) calcDropDstIndex(context *guigui.Context) int {
	_, y := ebiten.CursorPosition()
	for i := range t.columns {
		if r := t.ColumnRect(context, i); y < (r.Min.Y+r.Max.Y)/2 {
			return i
		}
	}
	return len(t.columns)
}

func (t *Table) HandleInput(context *guigui.Context) guigui.HandleInputResult {
	if x, y := ebiten.CursorPosition(); image.Pt(x, y).In(guigui.VisibleBounds(t)) {
		_, offsetY := t.scrollOverlay.Offset()
		y -= RoundedCornerRadius(context)
		y -= guigui.Position(t).Y
		y -= int(offsetY)
		index := -1
		var cy int
		for i, column := range t.columns {
			_, h := column.Size(context)
			if cy <= y && y < cy+h {
				index = i
				break
			}
			cy += h
		}
		t.setHoveredColumnIndex(index)
		if index >= 0 && index < len(t.columns) {
			left := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
			right := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight)

			switch {
			case left || right:
				if !t.columns[index].Selectable {
					return guigui.HandleInputByWidget(t)
				}

				wasFocused := guigui.IsFocused(t)
				guigui.Focus(t)
				if t.selectedColumnIndex() != index || !wasFocused {
					t.SetSelectedColumnIndex(index)
					t.lastSelectingColumnTime = time.Now()
				}
				t.pressStartX = x
				t.pressStartY = y
				if right { // TODO: send event
				}
				t.startPressingIndexPlus1 = index + 1
				t.startPressingLeft = left

			case ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft):
				// TODO: send event

			case inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft):
				if t.selectedColumnIndex() == index && t.startPressingLeft && time.Since(t.lastSelectingColumnTime) > 400*time.Millisecond {
					// TODO: send event
				}
				t.pressStartX = 0
				t.pressStartY = 0
				t.startPressingIndexPlus1 = 0
				t.startPressingLeft = false
			}

			return guigui.HandleInputByWidget(t)
		}
		t.dropSrcIndexPlus1 = 0
		t.pressStartX = 0
		t.pressStartY = 0
	} else {
		t.setHoveredColumnIndex(-1)
	}

	return guigui.HandleInputResult{}
}

func (t *Table) Update(context *guigui.Context) error {
	_, h := t.Size(context)
	t.scrollOverlay.SetContentSize(t.defaultWidth(context), h)

	idx := t.indexToJumpPlus1 - 1
	if idx >= 0 {
		y := t.rowYFromIndex(context, idx) - RoundedCornerRadius(context)
		t.scrollOverlay.SetOffset(0, float64(-y))
		t.indexToJumpPlus1 = 0
	}

	return nil
}

/*
func (t *Table) Update(context *guigui.Context) error {
	_, h := t.Size(context)
	t.scrollOverlay.SetContentSize(t.defaultWidth(context), h)

	idx := t.indexToJumpPlus1 - 1
	if idx >= 0 {
		x := t.columnXFromIndex(context, idx) - RoundedCornerRadius(context)
		t.scrollOverlay.SetOffset(float64(-x), 0)
		t.indexToJumpPlus1 = 0
	}

	return nil
}*/

func (t *Table) rowYFromIndex(context *guigui.Context, index int) int {
	x := RoundedCornerRadius(context)
	for _, column := range t.columns {
		cm := 0
		for i, cell := range column.cells {
			if i == index {
				break
			}
			_, ch := cell.Content.Size(context)
			cm = max(cm, ch)
		}
		x += cm
	}
	return x
}

func (t *Table) columnXFromIndex(context *guigui.Context, index int) int {
	x := RoundedCornerRadius(context)
	for i, column := range t.columns {
		if i == index {
			break
		}
		cw, _ := column.Size(context)
		x += cw
	}
	return x
}

func (t *Table) ColumnRect(context *guigui.Context, index int) image.Rectangle {
	_, offsetY := t.scrollOverlay.Offset()
	p := guigui.Position(t)
	w, h := t.Size(context)
	b := image.Rectangle{
		Min: p,
		Max: p.Add(image.Pt(w, h)),
	}
	padding := ColumnPadding(context)
	b.Min.X += RoundedCornerRadius(context) + padding
	b.Max.X -= RoundedCornerRadius(context) + padding
	b.Min.X += t.columnXFromIndex(context, index)

	b.Min.Y += RoundedCornerRadius(context) + padding
	b.Min.Y += int(offsetY)
	b.Max.Y -= RoundedCornerRadius(context) + padding

	iw, _ := t.columns[index].Size(context)
	b.Max.X = b.Min.X + iw
	return b
}

func (t *Table) selectedColumnColor(context *guigui.Context) color.Color {
	if t.selectedColumnIndex() < 0 || t.selectedColumnIndex() >= len(t.columns) {
		return nil
	}
	if guigui.IsFocused(t) {
		return Color(context.ColorMode(), ColorTypeAccent, 0.5)
	}
	return Color(context.ColorMode(), ColorTypeBase, 0.8)
}

func (t *Table) Draw(context *guigui.Context, dst *ebiten.Image) {
	clr := Color(context.ColorMode(), ColorTypeBase, 1)

	p := guigui.Position(t)
	w, h := t.Size(context)
	bounds := image.Rectangle{
		Min: p,
		Max: p.Add(image.Pt(w, h)),
	}
	DrawRoundedRect(context, dst, bounds, clr, RoundedCornerRadius(context))

	if clr := t.selectedColumnColor(context); clr != nil && t.selectedColumnIndex() >= 0 && t.selectedColumnIndex() < len(t.columns) {
		r := t.ColumnRect(context, t.selectedColumnIndex())
		r.Min.X -= RoundedCornerRadius(context)
		r.Max.X += RoundedCornerRadius(context)
		if r.Overlaps(guigui.VisibleBounds(t)) {
			DrawRoundedRect(context, dst, r, clr, RoundedCornerRadius(context))
		}
	}

	if t.isHoveringVisible() && t.hoveredColumnIndex() >= 0 && t.hoveredColumnIndex() < len(t.columns) && t.columns[t.hoveredColumnIndex()].Selectable {
		r := t.ColumnRect(context, t.hoveredColumnIndex())
		r.Min.X -= RoundedCornerRadius(context)
		r.Max.X += RoundedCornerRadius(context)
		if r.Overlaps(guigui.VisibleBounds(t)) {
			clr := Color(context.ColorMode(), ColorTypeBase, 0.9)
			DrawRoundedRect(context, dst, r, clr, RoundedCornerRadius(context))
		}
	}
}

func (t *Table) defaultWidth(context *guigui.Context) int {
	if t.cachedDefaultWidth > 0 {
		return t.cachedDefaultWidth
	}
	var w int
	for _, column := range t.columns {
		iw, _ := column.Size(context)
		w = +iw
	}
	w += 2*RoundedCornerRadius(context) + 2*ColumnPadding(context)
	t.cachedDefaultWidth = w
	return w
}

func (t *Table) defaultHeight(context *guigui.Context) int {
	if t.cachedDefaultHeight > 0 {
		return t.cachedDefaultHeight
	}

	var h int
	h += RoundedCornerRadius(context)
	for _, column := range t.columns {
		_, wh := column.Size(context)
		h = max(h, wh)
	}
	h += RoundedCornerRadius(context)
	t.cachedDefaultHeight = h
	return h
}

func (t *Table) Size(context *guigui.Context) (int, int) {
	var w, h int
	if t.widthSet {
		w = t.width
	} else {
		w = t.defaultWidth(context)
	}
	if t.heightSet {
		h = t.height
	} else {
		h = t.defaultHeight(context)
	}
	return w, h
}

func (t *Table) SetSize(width, height int) {
	t.width = width
	t.widthSet = true
	t.height = height
	t.heightSet = true
}

func (t *Table) SetWidth(width int) {
	t.width = width
	t.widthSet = true
}

func (t *Table) SetHeight(height int) {
	t.height = height
	t.heightSet = true
}

func (t *Table) ResetWidth() {
	t.widthSet = false
	t.width = 0
}

func (t *Table) ResetHeight() {
	t.heightSet = false
	t.height = 0
}
