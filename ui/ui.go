package ui

import (
	"fmt"
	"log"
	"os"

	db "github.com/Gabriel2233/keeper/database"
	"github.com/gdamore/tcell/v2"
)

/*
   draw for nf
   state ->
       currentFolderIndex
       currentSheetIndex

   events ->
       up, down: move focused view item
*/

var (
	defaultFolderStyle  = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorWhite)
	selectedFolderStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorBlue)

	defaultSheetStyle  = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorWhite)
	selectedSheetStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorGreen)
)

type Ui struct {
	folders []db.Folder
	sheets  []db.Sheet

	curViewIdx   int
	curFolderIdx int
	curSheetIdx  int

	screen tcell.Screen
	store  db.Store
}

func NewUi(store db.Store, folders []db.Folder, sheets []db.Sheet) *Ui {
	screen, err := createScreen()
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}

	return &Ui{
		store:        store,
		screen:       screen,
		folders:      folders,
		sheets:       sheets,
		curViewIdx:   0,
		curFolderIdx: 0,
		curSheetIdx:  0,
	}
}

func (ui *Ui) Draw() {
	if len(ui.folders) == 0 {
		drawBox(ui.screen, 0, 0, 15, 3, defaultFolderStyle, "No folders")
	} else {
        x1, y1 := 1, 1
        x2, y2 := x1 + 14, y1 + 3

		for i, f := range ui.folders {

			if ui.curFolderIdx == i {
				drawBox(ui.screen, x1, y1, x2, y2, selectedFolderStyle, f.Name)
			} else {
				drawBox(ui.screen, x1, y1, x2, y2, defaultFolderStyle, f.Name)
			}

			y1 += 4
			y2 += 4
		}
	}

	if len(ui.sheets) == 0 {
		drawBox(ui.screen, 40, 0, 65, 3, defaultSheetStyle, "No sheets")
	} else {
        x1, y1 := 40, 1
        x2, y2 := x1 + 14, y1 + 3

		for i, s := range ui.sheets {
			if ui.curSheetIdx == i {
				drawBox(ui.screen, x1, y1, x2, y2, selectedSheetStyle, s.Name)
			} else {
				drawBox(ui.screen, x1, y1, x2, y2, defaultSheetStyle, s.Name)
			}

			y1 += 4
			y2 += 4
		}
	}
}

func (ui *Ui) MoveRight() {
	if ui.curViewIdx+1 > 2 {
		ui.curViewIdx = 0
		return
	}

	ui.curViewIdx += 1
}

func (ui *Ui) MoveLeft() {
	if ui.curViewIdx-1 < 0 {
		ui.curViewIdx = 2
		return
	}

	ui.curViewIdx -= 1
}

func (ui *Ui) MoveUp() {
	if ui.curViewIdx == 0 {
		// Move down in the folders view and get the next item id
		if ui.curFolderIdx-1 < 0 {
			ui.curFolderIdx = len(ui.folders) - 1
		} else {
			ui.curFolderIdx -= 1
		}
		// fetch new sheets for that folder
		nextFolderSheets, err := ui.store.ListSheetsInFolder(ui.folders[ui.curFolderIdx].Name)
		if err != nil {
			ui.screen.Fini()
			fmt.Printf("error: %s\n", err)
			os.Exit(1)
		}

		ui.sheets = nextFolderSheets
		// redraw the screnn with the new items
	} else if ui.curViewIdx == 1 {
		// Move down in the sheets view and get the next item id
		if ui.curSheetIdx-1 < 0 {
			ui.curSheetIdx = len(ui.sheets) - 1
		} else {
			ui.curSheetIdx -= 1
		}
	} else {
		return
	}
}

func (ui *Ui) MoveDown() {
	if ui.curViewIdx == 0 {
		// Move down in the folders view and get the next item id
		if ui.curFolderIdx+1 >= len(ui.folders) {
			ui.curFolderIdx = 0
		} else {
			ui.curFolderIdx += 1
		}
		// fetch new sheets for that folder
		nextFolderSheets, err := ui.store.ListSheetsInFolder(ui.folders[ui.curFolderIdx].Name)
		if err != nil {
			ui.screen.Fini()
			fmt.Printf("error: %s\n", err)
			os.Exit(1)
		}

		ui.sheets = nextFolderSheets
		// redraw the screnn with the new items
	} else if ui.curViewIdx == 1 {
		// Move down in the sheets view and get the next item id
		if ui.curSheetIdx+1 >= len(ui.sheets) {
			ui.curSheetIdx = 0
		} else {
			ui.curSheetIdx += 1
		}
		// redraw the screen with the new items
	} else {
		return
	}
}

func (ui *Ui) Loop() {
	ui.screen.SetStyle(tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorRed))

	for {
	    ui.Draw()
	    ui.screen.Show()
		ev := ui.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
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
		if row > y2 {
			break
		}
	}
}

func drawBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

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
