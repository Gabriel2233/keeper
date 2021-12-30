package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

var (
	nfCmd = flag.NewFlagSet("nf", flag.ExitOnError)
	rfCmd = flag.NewFlagSet("rf", flag.ExitOnError)
	lfCmd = flag.NewFlagSet("lf", flag.ExitOnError)

	nsCmd = flag.NewFlagSet("ns", flag.ExitOnError)
	rsCmd = flag.NewFlagSet("rs", flag.ExitOnError)
	lsCmd = flag.NewFlagSet("ls", flag.ExitOnError)
)

func parseArgs(flagSet *flag.FlagSet, argsLen int) ([]string, error) {
	flagSet.Parse(os.Args[2:])
	args := flagSet.Args()

	if len(os.Args[2:]) != argsLen {
		return nil, errors.New(fmt.Sprintf("error: expected %d arguments, got %d", argsLen, len(args)))
	}

	return args, nil
}

func NewFolder(store *Store) (int64, error) {
	args, err := parseArgs(nfCmd, 1)
	if err != nil {
		return -1, err
	}

	id, err := store.AddFolder(args[0])
	if err != nil {
		return -1, err
	}

	return id, nil
}

func ListFolders(store *Store) ([]Folder, error) {
	var ret []Folder
	_, err := parseArgs(nfCmd, 0)
	if err != nil {
		return ret, err
	}

	ret, err = store.ListFolders()
	if err != nil {
		return ret, err
	}

	return ret, nil
}

func RemoveFolder(store *Store) error {
	args, err := parseArgs(nfCmd, 1)
	if err != nil {
		return err
	}

	err = store.RemoveFolderByName(args[0])
	if err != nil {
		return err
	}

	return nil
}

func NewSheet(store *Store) (int64, error) {
	args, err := parseArgs(nfCmd, 3)
	if err != nil {
		return -1, err
	}

	folderName, sheetName, alias := args[0], args[1], args[2]

	file, err := ioutil.TempFile("/tmp", "sample.*.txt")
	if err != nil {
		return -1, err
	}
	defer os.Remove(file.Name())

	vimCmd := exec.Command("vim", file.Name())
	vimCmd.Stderr = os.Stderr
	vimCmd.Stdout = os.Stdout
	vimCmd.Stdin = os.Stdin
	vimCmd.Run()

	data := make([]byte, 1024)
	_, err = file.Read(data)
	if err != nil {
		return -1, err
	}

	id, err := store.AddSheet(folderName, sheetName, alias, string(data))
	if err != nil {
		return -1, err
	}

	return id, nil
}

func ListSheetsUnderFolder(store *Store) ([]Sheet, error) {
	var ret []Sheet
	args, err := parseArgs(nfCmd, 1)
	if err != nil {
		return ret, err
	}

	folderName := args[0]

	ret, err = store.ListSheetsInFolder(folderName)
	if err != nil {
		return ret, err
	}

	return ret, nil
}

func RemoveSheet(store *Store) error {
	args, err := parseArgs(nfCmd, 1)
	if err != nil {
		return err
	}

	err = store.RemoveSheetByAlias(args[0])
	if err != nil {
		return err
	}

	return nil
}
