package main


var CONFIG = map[string]string{
    "client_oedir": "/home/jordan/code/openenclave",
    "client_privkey": "/home/jordan/.ssh/id_rsa",
    "target_username": "jorhand",
    "target_ip": "13.68.192.102",
    "target_oedir": "/home/jorhand/openenclave",
}

func main() {
	fw := new(FileWatcher)
	fw.Init()
	fw.AddRecursive(".")
	fw.Start()
}
