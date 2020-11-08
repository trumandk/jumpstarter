package main

import (
	"bytes"
	"fmt"
	"github.com/go-git/go-git"
	"github.com/go-git/go-git/plumbing/object"
	"github.com/go-git/go-git/plumbing/transport"
	"github.com/go-git/go-git/plumbing/transport/ssh"
	"github.com/pin/tftp"
	"github.com/tidwall/sjson"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

func MyPublicKeys() transport.AuthMethod {
	publicKeys, err := ssh.NewPublicKeysFromFile("git", "/root/.ssh/id_rsa", "")

	if err != nil {
		panic(err)
	}
	return publicKeys
}

var publicKeys = MyPublicKeys()

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func checkDockerExist(ip string) {
	go func() {
		dockerGitUpdate()
		newDockerFile := "/git/docker/" + ip
		if !fileExists(newDockerFile) {
			fmt.Printf("Generate docker file for ip:%s \n", ip)
			f, err := os.OpenFile(newDockerFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatal(err)
			}
			_, err2 := f.WriteString("version: '3.3'\r\n")
			if err2 != nil {
				log.Fatal(err2)
			}
			dockerGitCommit("docker/" + ip)
		}
	}()

}

func defaultFile(ip string) *bytes.Buffer {
	var response string
	response = "default flatcar-http\r\n"
	response += "prompt 1\r\n"
	response += "timeout 15\r\n\r\n"

	response += "LABEL flatcar-http\r\n"
	response += "LINUX http://" + ip + "/files/flatcar_production_pxe.vmlinuz\r\n"
	response += "APPEND initrd=http://" + ip + "/files/flatcar_production_pxe_image.cpio.gz rootfstype=tmpfs ignition.config.url=http://" + ip + "/ignition flatcar.first_boot=1 console=tty0 flatcar.autologin=tty1\r\n"
	buf := bytes.NewBufferString(response)
	return buf
}

func readHandler(filename string, r io.ReaderFrom) error {

	fmt.Printf("open: %s\n", filename)
	if strings.Contains(filename, "default") {
		ip := r.(tftp.RequestPacketInfo).LocalIP().String()
		ipRemote := r.(tftp.OutgoingTransfer).RemoteAddr()
		fmt.Printf("Generate default with ip:%s \n", ip)
		n, err := r.ReadFrom(defaultFile(ip))
		checkDockerExist(ipRemote.IP.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return err
		}
		fmt.Printf("%d bytes sent\n", n)
		return nil
	}
	file, err := os.Open(filename)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}

	if t, ok := r.(tftp.OutgoingTransfer); ok {
		if fi, err := file.Stat(); err == nil {
			t.SetSize(fi.Size())
		}
	}

	n, err := r.ReadFrom(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}
	fmt.Printf("%d bytes sent\n", n)
	return nil
}

func dockerInitGit() {
	err := os.RemoveAll("/git/")
	if err != nil {
		fmt.Printf("Remove folder :%s", err)
	}

	r, err := git.PlainClone("/git/", false, &git.CloneOptions{
		URL:      os.Getenv("GIT_CLUSTER"),
		Auth:     publicKeys,
		Progress: os.Stdout,
	})
	if err != nil {
		fmt.Printf("Git error :%s", err)
	}
	ref, err := r.Head()
	fmt.Printf("ref :%s", ref)
	if err != nil {
		fmt.Printf("head :%s", err)
	}
	fmt.Printf("ref :%s", ref)

}
func dockerGitCommit(filename string) {
	r, err := git.PlainOpen("/git/")
	if err != nil {
		fmt.Printf("plain open :%s", err)
	}
	w, err := r.Worktree()
	if err != nil {
		fmt.Printf("worktree error :%s", err)
	}
	w.Add(filename)

	commit, err2 := w.Commit("Auto-commit server:"+filename, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "JumpStarter",
			Email: "jumpstarter@jumpstarter.io",
			When:  time.Now(),
		},
	})
	if err2 != nil {
		fmt.Printf("Commit error:%s", err2)
	}
	obj, err3 := r.CommitObject(commit)
	if err3 != nil {
		fmt.Printf("CommitObject :%s", err3)
	}
	fmt.Println(obj)

	err5 := r.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       publicKeys,
		Progress:   os.Stdout,
	})
	if err5 != nil {
		fmt.Printf("push :%s", err5)
	}
	ref, err := r.Head()
	fmt.Printf("ref :%s", ref)
	if err != nil {
		fmt.Printf("head :%s", err)
	}
}

func dockerGitUpdate() {
	r, err := git.PlainOpen("/git/")
	if err != nil {
		fmt.Printf("plain open :%s", err)
	}
	w, err := r.Worktree()
	if err != nil {
		fmt.Printf("worktree error :%s", err)
	}

	err = w.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth:       publicKeys,
		Progress:   os.Stdout,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		fmt.Printf("pull error :%s", err)
	}
}

func dockerRun(ip string, file string) {
	out, err := exec.Command("/usr/bin/docker-compose", "-p", file, "--log-level", "CRITICAL", "--env-file", "/git/docker/env", "-H", "ssh://core@"+ip, "-f", "/git/docker/"+file, "up", "-d", "--remove-orphans").CombinedOutput()

	if err != nil {
		fmt.Printf("Error updating:%s Message:%s", ip, err)
	}
	output := string(out[:])
	if len(output) > 0 {
		fmt.Println(output)
	}
}

func dockercompose() {

	nodes, err := ioutil.ReadDir("/git/docker/")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range nodes {
		if f.Name() != "env" && f.Name() != "all" && dockerOnline(f.Name()) {
			dockerRun(f.Name(), "all")
			dockerRun(f.Name(), f.Name())
		}
	}

}
func IgnitionFile() string {
	pubKeyFile, err := os.Open("/root/.ssh/id_rsa.pub")
	if err != nil {
		panic(err)
	}
	pubKey, _ := ioutil.ReadAll(pubKeyFile)
	jsonFile, err := os.Open("/pxe-config.ign")
	if err != nil {
		panic(err)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)

	value, _ := sjson.Set(string(byteValue), "passwd.users.0.sshAuthorizedKeys.0", string(pubKey))
	defer jsonFile.Close()
	return value
}

var ignitionFile = IgnitionFile()

func ignitionWeb(w http.ResponseWriter, req *http.Request) {
	value, _ := sjson.Set(ignitionFile, "storage.files.0.contents.source", "http://"+req.Host+"/files/stat")
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(value))
	log.Println("Server PXE booted:", req.Host)

}

func main() {
	dockerInitGit()
	go func() {
		for {
			dockerGitUpdate()
			dockercompose()
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		// use nil in place of handler to disable read or write operations
		s := tftp.NewServer(readHandler, nil)
		s.SetTimeout(5 * time.Second)  // optional
		err := s.ListenAndServe(":69") // blocks until s.Shutdown() is called
		if err != nil {
			fmt.Fprintf(os.Stdout, "server: %v\n", err)
			os.Exit(1)
		}
	}()

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("/files"))
	mux.HandleFunc("/ssh", sshCommand)
	mux.HandleFunc("/ignition", ignitionWeb)
	mux.HandleFunc("/containers", containers)
	mux.HandleFunc("/status", status)
	mux.HandleFunc("/", servers)
	mux.Handle("/files/", http.StripPrefix("/files", fileServer))
	log.Println("Starting server on :80")
	err := http.ListenAndServe(":80", mux)
	log.Fatal(err)

}
