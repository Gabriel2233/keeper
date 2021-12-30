package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
)

const (
    FOLDER_CELL_WIDTH_STEP int = 50
    FOLDER_CELL_HEIGHT_STEP int = 3
)

var (
    defaultStyleListFolders = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorWhite)
    selectedStyleListFolders = tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorReset)

    
    defaultStyleListSheets = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorWhite)
    selectedStyleListSheets = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorGreen)
)

type ListFoldersState struct {
    screen tcell.Screen
    folders []Folder
    currentFolderIndex int

    defaultStyle tcell.Style
    selectedFolderStyle tcell.Style
}

func NewListFoldersState(s tcell.Screen, folders []Folder, currentFolderIndex int) *ListFoldersState {
    return &ListFoldersState{
        screen: s,
        folders: folders,
        currentFolderIndex: currentFolderIndex,

        defaultStyle: defaultStyleListFolders,
        selectedFolderStyle: selectedStyleListFolders,
    }
}

func (lfs *ListFoldersState) Draw() {
    x1, y1 := 1, 1
    x2, y2 := x1 + FOLDER_CELL_WIDTH_STEP, y1 + FOLDER_CELL_HEIGHT_STEP

    _, h := lfs.screen.Size()

    numFolders := 0

    for i, f := range lfs.folders {
        data := fmt.Sprintf("name: %s | created: %s | id: %d", f.Name, f.CreatedAt.Format("Mon _2 21"), f.Id)
        if lfs.currentFolderIndex == i {
            drawBox(lfs.screen, false, x1, y1, x2, y2, lfs.selectedFolderStyle, data)
        } else {
            drawBox(lfs.screen, false, x1, y1, x2, y2, lfs.defaultStyle, data)
        }

        y1 += FOLDER_CELL_HEIGHT_STEP
        y2 += FOLDER_CELL_HEIGHT_STEP
        numFolders += 1

        if numFolders > h / 18 {
            x1 += FOLDER_CELL_WIDTH_STEP
            x2 += FOLDER_CELL_WIDTH_STEP
            numFolders = 0

            y1, y2 = 1, 4
        }
    }
}

func (lfs *ListFoldersState) MoveUp() {
    if lfs.currentFolderIndex - 1 < 0 {
        lfs.currentFolderIndex = len(lfs.folders) - 1
        return
    }

    lfs.currentFolderIndex -= 1
}

func (lfs *ListFoldersState) MoveDown() {
    if lfs.currentFolderIndex + 1 >= len(lfs.folders) {
        lfs.currentFolderIndex = 0
        return
    }

    lfs.currentFolderIndex += 1
}

func (lfs *ListFoldersState) Loop() {
    lfs.screen.SetStyle(defaultStyleListFolders)
    lfs.Draw()
    for {
        lfs.screen.Show()
        ev := lfs.screen.PollEvent()
        switch ev := ev.(type) {
		case *tcell.EventResize:
			lfs.screen.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				lfs.screen.Fini()
				os.Exit(0)
			} else if ev.Rune() == 'C' || ev.Rune() == 'c' {
				lfs.screen.Clear()
            } else if ev.Rune() == 'j' || ev.Key() == tcell.KeyDown {
				lfs.screen.Clear()
                lfs.MoveDown()
                lfs.Draw()
            } else if ev.Rune() == 'k' || ev.Key() == tcell.KeyUp {
				lfs.screen.Clear()
                lfs.MoveUp()
                lfs.Draw()
            } else if ev.Key() == tcell.KeyEnter {
                id := lfs.Select()
                lfs.screen.Fini()

                fmt.Println("chose folder with id: ", id)
            }
        }
    }
}

func (lfs *ListFoldersState) Select() int64 {
    return lfs.folders[lfs.currentFolderIndex].Id
}

type ListSheetsState struct {
    screen tcell.Screen
    sheets []Sheet
    currentSheetIndex int

    defaultStyle tcell.Style
    selectedSheetStyle tcell.Style
}

func NewListSheetsState(s tcell.Screen, sheets []Sheet, currentSheetIndex int) *ListSheetsState {
    return &ListSheetsState{
        screen: s,
        sheets: sheets,
        currentSheetIndex: currentSheetIndex,

        defaultStyle: defaultStyleListSheets,
        selectedSheetStyle: selectedStyleListSheets,
    }
}

func setupScreen() (tcell.Screen, error) {
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

func drawBox(s tcell.Screen, hasBorder bool, x1, y1, x2, y2 int, style tcell.Style, text string) {
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

	if hasBorder {
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
	}

	drawText(s, x1+1, y1+1, x2-1, y2-1, style, text)
}
