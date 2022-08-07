package main

import (
	"flag"
	"log"
	"os"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	sshx "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
	"golang.org/x/crypto/ssh"
)

var caseNF = flag.String("case", "", "Case name.")

func main() {
	flag.Parse()

	var auth transport.AuthMethod
	var repo *git.Repository
	var repoURL = "git@github.com:gebv/tmp.git"
	var repoURLHttps = "https://github.com/gebv/tmp.git"

	switch *caseNF {
	case "1":
		log.Println("auth via ssh.PublicKeysCallback")
		auth = authViaSSHAgent()
	case "2", "":
		log.Println("nil auth method")
		auth = authEmptyMethod()
	case "3":
		log.Println("auth via ssh.Password")
		auth = authViaSSHPassword()
	case "4":
		log.Println("auth via ssh.PublicKeys")
		auth = authViaPublicKeysFromFile()
	case "5":
		repoURL = repoURLHttps
		log.Println("auth via http (nil auth method)")
		auth = nil
	default:
		log.Fatalf("not supported case name %q", *caseNF)
	}

	repo, err := git.Clone(
		memory.NewStorage(),
		nil,
		&git.CloneOptions{
			URL:        repoURL,
			NoCheckout: true,
			Progress:   os.Stdout,
			Auth:       auth,
		})
	if err != nil {
		log.Fatalf("failed clone: %v", err)
	}

	remote, err := repo.Remote("origin")
	if err != nil {
		log.Fatalf("failed remote origin: %v", err)
	}

	refList, err := remote.List(&git.ListOptions{
		Auth: auth,
	})
	if err != nil {
		log.Fatalf("failed get list: %v", err)
	}

	result := []string{}
	for _, ref := range refList {
		if ref.Name().IsBranch() {
			result = append(result, "branch/"+ref.Name().Short()+"/"+ref.Hash().String())
			continue
		}
		if ref.Name().IsTag() {
			result = append(result, "tag/"+ref.Name().Short()+"/"+ref.Hash().String())
			continue
		}
	}

	log.Println("list refs:", result)
}

func authViaSSHAgent() *sshx.PublicKeysCallback {
	auth, err := sshx.NewSSHAgentAuth("git")
	if err != nil {
		log.Fatalf("failed get ssh agent: %v", err)
	}
	auth.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	return auth
}

func authViaPublicKeysFromFile() *sshx.PublicKeys {
	auth, err := sshx.NewPublicKeysFromFile("git", "/root/.ssh/id_rsa", "")
	if err != nil {
		log.Fatalf("failed configure from public key: %v", err)
	}
	auth.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	return auth
}

func authEmptyMethod() transport.AuthMethod {
	return nil
}

func authViaSSHPassword() *sshx.Password {
	return &sshx.Password{
		User: "git",
		HostKeyCallbackHelper: sshx.HostKeyCallbackHelper{
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		},
	}
}
