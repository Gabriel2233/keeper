package ui

import (
	"fmt"
	"log"
	"os"
	"strings"

	db "github.com/Gabriel2233/keeper/database"
	"github.com/gdamore/tcell/v2"
)

const (
	folderSectionWidth int = 20
	sheetsSectionWidth int = 20
	sheetInfoItemWidth int = 18

	foldersInitialCol int = 0
	sheetsInitialCol  int = 20
	sheetInitialCol   int = 40
)

var (
	defaultFolderStyle  = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorWhite)
	selectedFolderStyle = tcell.StyleDefault.Background(tcell.ColorOrange).Foreground(tcell.ColorBlack)

	defaultSheetStyle  = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorWhite)
	selectedSheetStyle = tcell.StyleDefault.Background(tcell.ColorRed).Foreground(tcell.ColorReset).Bold(true)

	sheetTextAreaStyle        = tcell.StyleDefault.Background(tcell.ColorDodgerBlue).Foreground(tcell.ColorBlack)
	defaultSheetTextAreaStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorWhite)

	helpAreaStyle = tcell.StyleDefault.Background(tcell.ColorLawnGreen).Foreground(tcell.ColorBlack)
)

type Ui struct {
	folders    []db.Folder
	sheets     []db.Sheet
	sheetIndex int

	vCursor int
	hCursor int

	screen tcell.Screen
	store  db.Store
}

func NewUi(store db.Store, folders []db.Folder, sheets []db.Sheet) *Ui {
	screen, err := createScreen()
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}

	return &Ui{
		store:      store,
		screen:     screen,
		folders:    folders,
		sheets:     sheets,
		sheetIndex: -1,
		vCursor:    0,
		hCursor:    0,
	}
}

func (ui *Ui) Draw() {
	w, h := ui.screen.Size()
	folderRow := 0

	if len(ui.folders) == 0 {
		drawTextSection(ui.screen, foldersInitialCol, folderRow, folderSectionWidth, selectedFolderStyle, "No folders", false)
		return
	} else {
		for i, f := range ui.folders {
			if ui.hCursor == 0 && ui.vCursor == i {
				drawTextSection(ui.screen, foldersInitialCol, folderRow, folderSectionWidth, selectedFolderStyle, f.Name, true)
			} else {
				drawTextSection(ui.screen, foldersInitialCol, folderRow, folderSectionWidth, defaultFolderStyle, f.Name, true)
			}

			folderRow += 1
		}
	}

	sheetsRow := 0

	if len(ui.sheets) == 0 {
		drawTextSection(ui.screen, sheetsInitialCol, sheetsRow, sheetsSectionWidth, selectedSheetStyle, "No sheets", false)
	} else {
		for i, s := range ui.sheets {
			if ui.hCursor == 1 && ui.vCursor == i {
				drawTextSection(ui.screen, sheetsInitialCol, sheetsRow, sheetsSectionWidth, selectedSheetStyle, s.Name, true)
			} else {
				drawTextSection(ui.screen, sheetsInitialCol, sheetsRow, sheetsSectionWidth, defaultSheetStyle, s.Name, true)
			}

			sheetsRow += 1
		}
	}

	sheetRow := 0

	if ui.sheetIndex == -1 {
		drawTextSection(ui.screen, sheetInitialCol, sheetRow, w, sheetTextAreaStyle, "No Sheet", false)
	} else {
		id := fmt.Sprintf("Sheet Id: %d", ui.sheets[ui.sheetIndex].Id)
		drawTextSection(ui.screen, sheetInitialCol, sheetRow, sheetInfoItemWidth, sheetTextAreaStyle, id, false)

		name := fmt.Sprintf("Name: %s", ui.sheets[ui.sheetIndex].Name)
		drawTextSection(ui.screen, sheetInitialCol+sheetInfoItemWidth, sheetRow, sheetInfoItemWidth, sheetTextAreaStyle, name, true)

		alias := fmt.Sprintf("Alias: %s", ui.sheets[ui.sheetIndex].Alias)
		drawTextSection(ui.screen, sheetInitialCol+(2*sheetInfoItemWidth), sheetRow, sheetInfoItemWidth, sheetTextAreaStyle, alias, false)

		createdAt := fmt.Sprintf("Created: %s", ui.sheets[ui.sheetIndex].CreatedAt.Format("Mon _2 2021"))
		drawTextSection(ui.screen, sheetInitialCol+(3*sheetInfoItemWidth), sheetRow, w, sheetTextAreaStyle, createdAt, false)
		sheetRow += 2

		drawText(ui.screen, sheetInitialCol, sheetRow, w-1, h-2, defaultSheetStyle, ui.sheets[ui.sheetIndex].Data)
	}

	help := "[help] [move: ↑ → ↓ ←]  [delete item under cursor: del] [exit: ctrl+c]"
	drawTextSection(ui.screen, 0, h-1, w+len(help), helpAreaStyle, help, false)
}

func (ui *Ui) MoveRight() {
	if ui.hCursor == 0 && len(ui.sheets) == 0 {
		return
	}

	if ui.hCursor == 0 {
		ui.hCursor = 1
	} else {
		ui.hCursor = 0
	}

	if ui.hCursor == 1 && len(ui.sheets) > 0 {
		ui.vCursor = 0
		ui.sheetIndex = ui.vCursor
	}
}

func (ui *Ui) MoveLeft() {
	if ui.hCursor == 0 && len(ui.sheets) == 0 {
		return
	}

	if ui.hCursor == 0 {
		ui.hCursor = 1
	} else {
		ui.hCursor = 0
	}

	if ui.hCursor == 1 && len(ui.sheets) > 0 {
		ui.vCursor = 0
		ui.sheetIndex = ui.vCursor
	}
}

func (ui *Ui) MoveUp() {
	ui.vCursor -= 1
	if ui.vCursor < 0 && ui.hCursor == 0 {
		ui.vCursor = len(ui.folders) - 1
	} else if ui.vCursor < 0 && ui.hCursor == 1 {
		ui.vCursor = len(ui.sheets) - 1
	}

	if ui.hCursor == 0 {
		id := ui.folders[ui.vCursor].Id
		updatedSheets, _ := ui.store.ListSheetsInFolder(id)
		ui.sheets = updatedSheets
		ui.sheetIndex = -1
	} else {
		ui.sheetIndex = ui.vCursor
	}
}

func (ui *Ui) MoveDown() {
	ui.vCursor += 1
	if ui.hCursor == 0 && ui.vCursor >= len(ui.folders) {
		ui.vCursor = 0
	} else if ui.hCursor == 1 && ui.vCursor >= len(ui.sheets) {
		ui.vCursor = 0
	}

	if ui.hCursor == 0 {
		id := ui.folders[ui.vCursor].Id
		updatedSheets, _ := ui.store.ListSheetsInFolder(id)
		ui.sheets = updatedSheets
		ui.sheetIndex = -1
	} else {
		ui.sheetIndex = ui.vCursor
	}
}

func (ui *Ui) Delete() {
	if ui.hCursor == 0 {
		_ = ui.store.RemoveFolderById(ui.folders[ui.vCursor].Id)
		ui.folders = append(ui.folders[:ui.vCursor], ui.folders[ui.vCursor+1:]...)

		if len(ui.folders) > 0 {
			ui.vCursor -= 1
		}
		return
	}

	_ = ui.store.RemoveSheetById(ui.sheets[ui.vCursor].Id)
	ui.sheets = append(ui.sheets[:ui.vCursor], ui.sheets[ui.vCursor+1:]...)

	if len(ui.sheets) > 0 {
		ui.vCursor -= 1
		ui.sheetIndex = ui.vCursor
	} else {
		ui.sheetIndex = -1
		ui.hCursor = 0
		ui.vCursor = 0
	}
}

func (ui *Ui) Loop() {
	ui.screen.SetStyle(tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorWhite))

	for {
		ui.Draw()
		ui.screen.Show()
		ev := ui.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			ui.screen.Clear()
			ui.screen.Sync()
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyCtrlC:
				ui.screen.Fini()
				os.Exit(0)
			case tcell.KeyRight:
				ui.screen.Clear()
				ui.MoveRight()
			case tcell.KeyLeft:
				ui.screen.Clear()
				ui.MoveLeft()
			case tcell.KeyUp:
				ui.screen.Clear()
				ui.MoveUp()
			case tcell.KeyDown:
				ui.screen.Clear()
				ui.MoveDown()
			case tcell.KeyDelete:
				ui.screen.Clear()
				ui.Delete()
			}
		}
	}
}

func createScreen() (tcell.Screen, error) {
	s, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	if err = s.Init(); err != nil {
		return nil, err
	}

	return s, nil
}

func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if r == '\n' {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

func drawTextSection(s tcell.Screen, startCol, startRow, width int, style tcell.Style, text string, constrainStr bool) {
	col := startCol
	row := startRow

	s.SetContent(col, row, ' ', nil, style)
	col++
	width--

	if len(text) > 40 && constrainStr {
		var b strings.Builder
		for i, r := range text {
			if i >= 16 {
				b.WriteString("...")
				break
			}
			b.WriteRune(r)
		}
		text = b.String()
	}

	width -= len(text)

	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
	}

	for i := 0; i < width; i++ {
		s.SetContent(col, row, ' ', nil, style)
		col++
	}
}

func drawBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	for row := y1; row <= y2; row++ {
		for col := x1; col <= x2; col++ {
			s.SetContent(col, row, ' ', nil, style)
		}
	}

	for col := x1; col <= x2; col++ {
		s.SetContent(col, y1, tcell.RuneHLine, nil, style)
		s.SetContent(col, y2, tcell.RuneHLine, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		s.SetContent(x1, row, tcell.RuneVLine, nil, style)
		s.SetContent(x2, row, tcell.RuneVLine, nil, style)
	}

	if y1 != y2 && x1 != x2 {
		s.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
		s.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
		s.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
		s.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
	}

	drawText(s, x1+1, y1+1, x2-1, y2-1, style, text)
}
