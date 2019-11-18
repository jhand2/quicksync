package main

import "fmt"
import "path/filepath"
import "log"
import "os"
import "regexp"
import "github.com/fsnotify/fsnotify"

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

		return nil;
	}

	// starting at the root of the project, walk each file/directory searching for
	// directories
	if err := filepath.Walk(path, addDir); err != nil {
		fmt.Println("ERROR", err)
	}
}

func (f *FileWatcher) isIgnored(path string) bool {
	// TODO: Make ignore regex extensible
	match, err := regexp.MatchString(`(\.swp)|(\.swx)|(~$)`, path)
	if err != nil {
		log.Fatal(err)
	}

	return match
}

func (f *FileWatcher) handleEvent(e fsnotify.Event) {
	if !f.isIgnored(e.Name) {
		if e.Op == fsnotify.Write {
			log.Println("event:", e)
			log.Println(e.Name)
		}
	}
}

func (f *FileWatcher) Start() {
	defer f.w.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-f.w.Events:
				if !ok {
					return
				}

				f.handleEvent(event)
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
