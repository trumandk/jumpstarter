package main

import (
	//"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types"
	"log"
	"net/http"
	"os/exec"
)
	//var ctx = context.Background()
	var cli, errDocker = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

func containers(w http.ResponseWriter, req *http.Request) {



	containerList, err5 := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err5 != nil {
		panic(err5)
	}
	b, err := json.MarshalIndent(containerList, "", " ")
	if err != nil {
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func main() {
	if errDocker != nil {
		panic(errDocker)
	}
        exec.Command("/bin/mkdir","/var/lib/docker").CombinedOutput()
        out, err2 := exec.Command("/usr/bin/mount","-t","ext4","/dev/sda1","/var/lib/docker").CombinedOutput()

        if err2 != nil {
                fmt.Printf("Error Message:%s", err2)
        }
       output := string(out[:])
        fmt.Println(output)


	mux := http.NewServeMux()
	mux.HandleFunc("/", status)
	mux.HandleFunc("/containers", containers)
	err := http.ListenAndServe(":4", mux)
	log.Fatal(err)
}
