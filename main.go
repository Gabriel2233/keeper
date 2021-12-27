package main

import (
	"fmt"
	"os"
	"text/tabwriter"
)

// keeper nf <FOLDER>
// keeper lf
// keeper rf <FOLDER> [i]

// keper ns <FOLDER> <NAME> <ALIAS>
// Data:

// keeper ls <FOLDER>
// keeper rs <SHEET> | <ALIAS>

func main() {
    store, err := NewStore("./store.db")
    if err != nil {
        panic(err)
    }

    if len(os.Args) == 1 {
        usage()
        os.Exit(0)
    }

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

    switch os.Args[1] {
    case "nf":
        id, err := NewFolder(store)
        if err != nil {
            fmt.Printf("error: failed to create folder, reason: \n%s\n", err)
            os.Exit(1)
        }

        fmt.Printf("created new folder with id %d", id)
    case "lf":
        folders, err := ListFolders(store)
        if err != nil {
            fmt.Printf("error: failed to list folders, reason: \n%s\n", err)
            os.Exit(1)
        }

        fmt.Println("Folder list:")
        for _, f := range folders {
            fmt.Fprintf(w, "Id: %d\tName: %s\tCreated At: %s\n", f.Id, f.Name, f.CreatedAt.Format("Mon Jan _2"))
        }
    case "rf":
        err := RemoveFolder(store) 
        if err != nil {
            fmt.Printf("error: failed to remove folder, reason: \n%s\n", err)
            os.Exit(1)
        }

        fmt.Println("folder removed successfully")
    case "ns":
        id, err := NewSheet(store)
        if err != nil {
            fmt.Printf("error: failed to create sheet, reason: \n%s\n", err)
            os.Exit(1)
        }

        fmt.Printf("created new sheet with id %d", id)
    case "ls":
        sheets, err := ListSheetsUnderFolder(store)
        if err != nil {
            fmt.Printf("error: failed to list sheets, reason: \n%s\n", err)
            os.Exit(1)
        }

        for _, s := range sheets {
            fmt.Fprintf(w, "Id: %d\tName: %s\tAlias: %s\tCreated At: %s\n", s.Id, s.Name, s.Alias, s.CreatedAt.Format("Mon Jan _2"))
        }
    case "rs":
        err := RemoveSheet(store) 
        if err != nil {
            fmt.Printf("error: failed to remove sheet, reason: \n%s\n", err)
            os.Exit(1)
        }

        fmt.Println("sheet removed successfully")
    }

    w.Flush()
}

func usage() {
    fmt.Println("Usage: ")
}
