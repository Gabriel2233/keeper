package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/Gabriel2233/keeper/cmd"
	db "github.com/Gabriel2233/keeper/database"
	"github.com/Gabriel2233/keeper/ui"
)

func main() {
    store, err := db.NewStore("./store.db")
    if err != nil {
        panic(err)
    }

    folders, err := store.ListFolders()
    must(err)

    sheets, err := store.ListSheetsInFolder(folders[0].Name)
    must(err)

    if len(sheets) == 0 {
        sheets = []db.Sheet{}
    }

    if len(os.Args) == 1 {
        ui := ui.NewUi(*store, folders, sheets)
        ui.Loop()
    }

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

    switch os.Args[1] {
    case "nf":
        id, err := cmd.NewFolder(store)
        must(err)

        fmt.Printf("created new folder with id %d", id)
    case "rf":
        err := cmd.RemoveFolder(store) 
        must(err)

        fmt.Println("folder removed successfully")
    case "ns":
        id, err := cmd.NewSheet(store)
        must(err)

        fmt.Printf("created new sheet with id %d", id)
    case "rs":
        err := cmd.RemoveSheet(store) 
        must(err)

        fmt.Println("sheet removed successfully")
    }

    w.Flush()
}

func usage() {
    fmt.Println("Usage: ")
}

func must(err error) {
    if err != nil {
        fmt.Printf("error: %s\n", err)
        os.Exit(1)
    }
}
