package main

import (
	"bytes"
	"fmt"
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
	response += "LINUX http://" + ip + "/flatcar_production_pxe.vmlinuz\r\n"
	response += "APPEND initrd=http://" + ip + "/flatcar_production_pxe_image.cpio.gz rootfstype=tmpfs ignition.config.url=http://" + ip + "/pxe-config.ign flatcar.first_boot=1 console=tty0 flatcar.autologin=tty1\r\n"
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

func dockercompose() {

	nodes, err := ioutil.ReadDir("./docker/")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range nodes {
		if f.Name() != "env" {
			fmt.Println(f.Name())

			out, err := exec.Command("/usr/bin/docker-compose", "-p", f.Name(), "--env-file", "./docker/env", "-H", "tcp://"+f.Name()+":2375", "-f", "docker/"+f.Name(), "up", "-d", "--remove-orphans").CombinedOutput()

			if err != nil {
				fmt.Printf("Error updating:%s Message:%s", f.Name(), err)
				output := string(out[:])
				fmt.Println(output)
			}
		}
	}

}

func main() {

	go func() {
		for {
			dockercompose()
		}
	}()

	http.Handle("/", http.FileServer(http.Dir("/files/")))

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
	http.ListenAndServe(":80", nil)

}
