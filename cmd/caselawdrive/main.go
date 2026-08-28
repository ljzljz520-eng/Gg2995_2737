package main

import (
	"fmt"
	"lawdrive/internal/storage"
	"os"
)

func main() {
	path := "lawdrive.db"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	store, err := storage.Open(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer store.Close()
	cases, err := store.Cases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("lawdrive ready: %d cases\n", len(cases))
}
