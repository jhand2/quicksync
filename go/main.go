package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

func read_config(config_name string) map[string]string {
	f, err := os.Open("/home/jordan/.quicksync.conf")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	bytes, _ := ioutil.ReadAll(f)

	var conf map[string]map[string]string
	json.Unmarshal(bytes, &conf)
	return conf[config_name]
}

func main() {
	var config = read_config("linux_acc_vm")

	fw := new(FileWatcher)
	fw.Init()
	fw.AddRecursive(config["client_oedir"])
	fw.Start(func(e fsnotify.Event) {
		copy_file(e.Name, e.Name)
	})
}
