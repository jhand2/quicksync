package main

import "fmt"
import "time"
import "github.com/radovskyb/watcher"


func watch_files() {
	w := watcher.New()
	w.SetMaxEvents(1)
	w.FilterOps(watcher.Rename, watcher.Move, watcher.Create, watcher.Remove)

	// Set up listener
	go func() {
		for {
			select {
			case event := <-w.Event:	
				fmt.Println(event) // Print the event's info.
			case err := <-w.Error:
				fmt.Println(err)
				return
			case <-w.Closed:
				return
			}
		}
	}()

	// Watch current folder for changes.
	if err := w.AddRecursive("."); err != nil {
		fmt.Println(err)
		return;
	}

	// Start listening
	go func() {
		w.Wait()
	}()

	// Start the watching process - it'll check for changes every 100ms.
	if err := w.Start(time.Millisecond * 100); err != nil {
		fmt.Println(err)
		return;
	}
}

func main() {
	watch_files()
}
