package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net/http"
)

func sshCommand(w http.ResponseWriter, req *http.Request) {
	commands, commandok := req.URL.Query()["command"]
	ips, ipok := req.URL.Query()["ip"]

	if commandok && ipok {
		command := commands[0]
		ip := ips[0]
		go func() {
			client, session, err := connectToHost("core", ip+":22")
			if err != nil {
				fmt.Println(err)
				return
			}
			out, err := session.CombinedOutput(command)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(string(out))
			client.Close()
		}()
		http.Redirect(w, req, "/", 304)
	}
}
func sshCommandOutput(w http.ResponseWriter, req *http.Request) {
	commands, commandok := req.URL.Query()["command"]
	ips, ipok := req.URL.Query()["ip"]

	if commandok && ipok {
		command := commands[0]
		ip := ips[0]
		client, session, err := connectToHost("core", ip+":22")
		if err != nil {
			fmt.Println(err)
			return
		}
		out, err := session.CombinedOutput(command)
		if err != nil {
			fmt.Println(err)
			return
		}
		//fmt.Println(string(out))
		w.Write([]byte(out))
		client.Close()
	}
}

func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

func connectToHost(user, host string) (*ssh.Client, *ssh.Session, error) {

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{PublicKeyFile("/root/.ssh/id_rsa")},
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	client, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		return nil, nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, nil, err
	}

	return client, session, nil
}
