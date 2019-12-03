package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

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
	if len(os.Args) < 2 {
		fmt.Println("Usage: quicksync <config_name>")
		os.Exit(1)
	}

	var config = read_config(os.Args[1])

	fw := new(FileWatcher)
	fw.Init()
	watch_dir := config["client_oedir"]
	fmt.Printf("Watching directory %s\n", watch_dir)

	scp := new(ScpCopier)
	scp.Init(config["target_username"],
		config["client_privkey"],
		config["target_ip"]+":22")

	fw.AddRecursive(watch_dir)
	fw.Start(func(e fsnotify.Event) {
		suffix := strings.TrimPrefix(e.Name, watch_dir)
		fname := config["target_oedir"] + suffix
		scp.copy_file(e.Name, fname)
	})
}
