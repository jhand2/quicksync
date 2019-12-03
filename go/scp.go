package main

import (
	"fmt"
	"log"
	"os"
	"syscall"

	scp "github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

type ScpCopier struct {
	client scp.Client
}

func (s *ScpCopier) Init(username string, privkey string, server string) {
	// Get password for privkey
	fmt.Print("Enter Password: ")
	pw, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatal("Could not read password. Error: ", err)
	}

	conf, err := auth.PrivateKeyWithPassphrase(username, pw, privkey, ssh.InsecureIgnoreHostKey())
	if err != nil {
		log.Fatal("Could not use private key. Error: ", err)
	}

	s.client = scp.NewClient(server, &conf)
}

func (s *ScpCopier) copy_file(src string, dst string) {
	fmt.Printf("Copying %s to %s\n", src, dst)
	err := s.client.Connect()
	if err != nil {
		log.Fatal("Could not open ssh connection. Error: ", err)
	}

	f, err := os.Open(src)
	if err != nil {
		log.Fatal("Could not open file. Error: ", err)
	}

	defer s.client.Close()
	defer f.Close()

	// TODO: Change file permissions?
	err = s.client.CopyFile(f, dst, "0655")
	if err != nil {
		log.Fatal("Could not copy file. Error: ", err)
	}
}
