package main

import (
	"context"
        "fmt"
        "io/ioutil"
        "log"
        "net/http"
		//"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)


func dockerStat(w http.ResponseWriter, server string){
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithHost("tcp://" + server + ":2375"), client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	/*
		containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

        //	fmt.Fprintf(w, "<table>")
	var running, total uint64
	for _, container := range containers {
	total++	
		if container.State == "running" {
			running++
		}*/
       // fmt.Fprintf(w, "<tr>")
       /*
                fmt.Fprintf(w, "<td>%s</td>", container.ID)
                fmt.Fprintf(w, "<td>%s</td>", container.Status)
                fmt.Fprintf(w, "<td>%s</td>", container.State)
                fmt.Fprintf(w, "<td>%s</td>", container.Names)
		fmt.Println(container.Names)
       // fmt.Fprintf(w, "</tr>")
	}
		*/
          //      fmt.Fprintf(w, "<td>%d</td>", running)
          //      fmt.Fprintf(w, "<td>%d</td>", total)
        //	fmt.Fprintf(w, "</table>")
		info, err2 := cli.Info(ctx)
	if err2 != nil {
		panic(err2)
	}
                fmt.Fprintf(w, "<td>%d</td>", info.Containers)
                fmt.Fprintf(w, "<td>%d</td>", info.ContainersRunning)
                fmt.Fprintf(w, "<td>%d</td>", info.ContainersPaused)
                fmt.Fprintf(w, "<td>%d</td>", info.ContainersStopped)
                fmt.Fprintf(w, "<td>%d</td>", info.MemTotal)
                fmt.Fprintf(w, "<td>%s</td>", info.SystemTime)
		/*stats, err3 := cli.Stat(ctx)
	if err3 != nil {
		panic(err3)
	}
	*/

}

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
        fmt.Fprintf(w, "<th>Containers</th>")
        fmt.Fprintf(w, "<th>Running</th>")
        fmt.Fprintf(w, "<th>Paused</th>")
        fmt.Fprintf(w, "<th>Stopped</th>")
        fmt.Fprintf(w, "<th>MemTotal</th>")
        fmt.Fprintf(w, "<th>SystemTime</th>")
        fmt.Fprintf(w, "</tr>")

        for _, f := range nodes {
                if f.Name() != "env" {
                        result, time := pingTest(f.Name())
                        //      fmt.Fprintf(w, "ip:%s result:%s ping:%s<br>\n", f.Name(), result, time)
                        if result {
                                //fmt.Fprintf(w, "<tr>")
                                fmt.Fprintf(w, "<tr style=\"background-color:#00FF00\">")
                        } else {
                                fmt.Fprintf(w, "<tr style=\"background-color:#FF0000\">")
                        }
                        fmt.Fprintf(w, "<td>%s</td>", f.Name())
                        fmt.Fprintf(w, "<td>%s</td>", time)
        //                fmt.Fprintf(w, "<td>")
                        if result {
			dockerStat(w, f.Name())
			}
          //              fmt.Fprintf(w, "</td>")
                        fmt.Fprintf(w, "</tr>")

                }
        }
}


