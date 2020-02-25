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

func print_selected_config(config map[string]string) {
	config_marshalled, err := json.MarshalIndent(config, "", " ")
	if err == nil {
		fmt.Println("Using config", string(config_marshalled))
	}
}

func make_target_path(config map[string]string, source_path string) string {
	client_dir := config["client_dir"]
	suffix := strings.TrimPrefix(source_path, client_dir)

	var sep = ""

	// Only add a path separator if config["target_dir"] does not already end
	// with one
	if !strings.HasSuffix(config["target_dir"], string(os.PathSeparator)) &&
		!strings.HasPrefix(suffix, string(os.PathSeparator)) {
		sep = string(os.PathSeparator)
	}

	return config["target_dir"] + sep + suffix
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: quicksync <config_name>")
		os.Exit(1)
	}

	// Get name of configuration from args and read it from user's config file
	var config = read_config(os.Args[1])
	print_selected_config(config)

	// Create scp client
	scp := new(ScpCopier)
	scp.Init(config["target_username"],
		config["client_privkey"],
		config["target_ip"]+":22")

	// Copy uncommitted files with modifications
	git := new(GitClient)
	git.Init(config["client_dir"])
	matches, err := git.Status()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println()
	for _, m := range matches {
		src_name := config["client_dir"] + string(os.PathSeparator) + m
		target_name := make_target_path(config, m)
		err = scp.copy_file(src_name, target_name)
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println()

	// Create and start file watcher
	fw := new(FileWatcher)
	fw.Init()

	watch_dir := config["client_dir"]
	fmt.Printf("Watching directory %s...\n", watch_dir)
	fw.AddRecursive(watch_dir)
	fw.Start(func(e fsnotify.Event) {
		fname := make_target_path(config, e.Name)
		// Intentionally ignore errors in copying the file
		scp.copy_file(e.Name, fname)
	})
}
