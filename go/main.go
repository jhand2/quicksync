package main

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
)

func main() {
	fw := new(FileWatcher)
	fw.Init()
	fw.AddRecursive(".")
	fw.Start(func(e fsnotify.Event) {
		fmt.Println(e.Name)
	})
}
