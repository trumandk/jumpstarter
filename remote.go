package main

import (
        "fmt"
      "net/http"
//      "os"
        "io/ioutil"
        "golang.org/x/crypto/ssh"
)

func commandSSH(ip string, command string) {

go func() {
        client, session, err := connectToHost("core", ip + ":22")
        if err != nil {
                //panic(err)
		fmt.Println(err)
        }
        out, err := session.CombinedOutput(command)
        if err != nil {
                //panic(err)
		fmt.Println(err)
        }
        fmt.Println(string(out))
        client.Close()
        }()

}

func reboot(w http.ResponseWriter, req *http.Request) {
        client, session, err := connectToHost("core", "192.168.1.132:22")
        if err != nil {
		fmt.Println(err)
                //panic(err)
        }
        out, err := session.CombinedOutput("sudo reboot")
        if err != nil {
		fmt.Println(err)
        }
        fmt.Println(string(out))
        client.Close()
}
func poweroff(w http.ResponseWriter, req *http.Request) {
        client, session, err := connectToHost("core", "192.168.1.132:22")
        if err != nil {
		fmt.Println(err)
        }
        out, err := session.CombinedOutput("sudo poweroff")
        if err != nil {
		fmt.Println(err)
        }
        fmt.Println(string(out))
        client.Close()
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
                Auth: []ssh.AuthMethod{PublicKeyFile("/.ssh/id_rsa")},
                //Auth: []ssh.AuthMethod{ssh.Password(pass)},
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

