package main

import (
	"bytes"
	"fmt"
	"github.com/go-git/go-git"
	"github.com/go-git/go-git/plumbing/transport/ssh"
	"github.com/pin/tftp"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

func defaultFile(ip string) *bytes.Buffer {
	var response string
	response = "default flatcar-http\r\n"
	response += "prompt 1\r\n"
	response += "timeout 15\r\n\r\n"

	response += "LABEL flatcar-http\r\n"
	response += "LINUX http://" + ip + "/files/flatcar_production_pxe.vmlinuz\r\n"
	response += "APPEND initrd=http://" + ip + "/files/flatcar_production_pxe_image.cpio.gz rootfstype=tmpfs ignition.config.url=http://" + ip + "/files/pxe-config.ign flatcar.first_boot=1 console=tty0 flatcar.autologin=tty1\r\n"
	buf := bytes.NewBufferString(response)
	return buf
}

func readHandler(filename string, r io.ReaderFrom) error {

	fmt.Printf("open: %s\n", filename)
	if strings.Contains(filename, "default") {
		ip := r.(tftp.RequestPacketInfo).LocalIP().String()
		fmt.Printf("Generate default with ip:%s \n", ip)
		n, err := r.ReadFrom(defaultFile(ip))
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
	publicKeys, err := ssh.NewPublicKeysFromFile("git", "/root/.ssh/id_rsa", "")

	if err != nil {
		fmt.Printf("generate publickeys failed :%s", err)
		return
	}
	//	fmt.Printf("pubkey :%s", publicKeys)

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

func dockerGitUpdate() {
	r, err := git.PlainOpen("/git/")
	if err != nil {
		fmt.Printf("plain open :%s", err)
	}
	w, err := r.Worktree()
	if err != nil {
		fmt.Printf("worktree error :%s", err)
	}
	publicKeys, err := ssh.NewPublicKeysFromFile("git", "/root/.ssh/id_rsa", "")

	if err != nil {
		fmt.Printf("generate publickeys failed :%s", err)
		return
	}
	err = w.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth:       publicKeys,
		Progress:   os.Stdout,
	})
	if err != nil {
		fmt.Printf("pull :%s", err)
	}
	ref, err := r.Head()
	fmt.Printf("ref :%s", ref)
	if err != nil {
		fmt.Printf("head :%s", err)
	}
}

func dockercompose() {

	nodes, err := ioutil.ReadDir("/git/docker/")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range nodes {
		if f.Name() != "env" && f.Name() != "all" && dockerOnline(f.Name()) {
			fmt.Println(f.Name())

			out, err := exec.Command("/usr/bin/docker-compose", "-p", f.Name(), "--env-file", "/git/docker/env", "-H", "ssh://core@"+f.Name(), "-f", "/git/docker/"+f.Name(), "-f", "/git/docker/all", "up", "-d", "--remove-orphans").CombinedOutput()

			if err != nil {
				fmt.Printf("Error updating:%s Message:%s", f.Name(), err)
			}
			output := string(out[:])
			fmt.Println(output)
		}
	}

}

/*
func status(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "<html>")
	fmt.Fprintf(w, "<head><meta http-equiv=\"refresh\" content=\"5\"></head>")
	fmt.Fprintf(w, "<body>")
	fmt.Fprintf(w, "<center>")
	nodes, err := ioutil.ReadDir("/git/docker/")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "<table border=1>")
	fmt.Fprintf(w, "<tr>")
	fmt.Fprintf(w, "<th>IP</th>")
	fmt.Fprintf(w, "<th>Ping</th>")
	fmt.Fprintf(w, "</tr>")

	for _, f := range nodes {
		if f.Name() != "env" {
			result, time := pingTest(f.Name())
			//	fmt.Fprintf(w, "ip:%s result:%s ping:%s<br>\n", f.Name(), result, time)
			if result {
				//fmt.Fprintf(w, "<tr>")
				fmt.Fprintf(w, "<tr style=\"background-color:#00FF00\">")
			} else {
				fmt.Fprintf(w, "<tr style=\"background-color:#FF0000\">")
			}
			fmt.Fprintf(w, "<td>%s</td>", f.Name())
			fmt.Fprintf(w, "<td>%s</td>", time)
			fmt.Fprintf(w, "</tr>")

		}
	}
}
*/
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
	mux.HandleFunc("/", status)
	mux.Handle("/files/", http.StripPrefix("/files", fileServer))
	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":80", mux)
	log.Fatal(err)

}
