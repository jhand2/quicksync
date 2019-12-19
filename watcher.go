package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/fsnotify/fsnotify"
)

type WatchCallback func(e fsnotify.Event)

type FileWatcher struct {
	w *fsnotify.Watcher
}

func (f *FileWatcher) Init() {
	var err error
	f.w, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
}

func (f *FileWatcher) Add(path string) {
	f.w.Add(path)
}

func (f *FileWatcher) AddRecursive(path string) {
	addDir := func(path string, fi os.FileInfo, err error) error {
		if fi.Mode().IsDir() {
			f.w.Add(path)
		}

		return nil
	}

	// starting at the root of the project, walk each file/directory searching for
	// directories
	if err := filepath.Walk(path, addDir); err != nil {
		fmt.Println("ERROR", err)
	}
}

func (f *FileWatcher) isIgnored(e fsnotify.Event) bool {
	// TODO: Make ignore regex extensible
	match, err := regexp.MatchString(`(\.swp)|(\.swx)|(~$)|(/.git/)`, e.Name)
	if err != nil {
		log.Fatal(err)
	}

	if match {
		return true
	}

	// TODO: Make Event type filter configurable
	if e.Op != fsnotify.Write {
		return true
	}

	return false
}

func (f *FileWatcher) handleEvent(e fsnotify.Event, callback WatchCallback) {
	if !f.isIgnored(e) {
		callback(e)
	}
}

func (f *FileWatcher) Start(callback WatchCallback) {
	defer f.w.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-f.w.Events:
				if !ok {
					return
				}

				f.handleEvent(event, callback)
			case err, ok := <-f.w.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Wait until done listening
	<-done
}
