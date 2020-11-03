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
/*
func sshCommand(w http.ResponseWriter, req *http.Request) {
	commands, commandok := req.URL.Query()["command"]
	ips, ipok := req.URL.Query()["ip"]

	 if commandok && ipok {
		 command := commands[0]
		 ip := ips[0]
       // 	fmt.Fprintf(w, "skod command=%s ip%s", command, ip)
		 commandSSH(ip, command)
		 http.Redirect(w, req, "/", 304)
    	}
}
*/
func status(w http.ResponseWriter, req *http.Request) {


        fmt.Fprintf(w, "<head><meta http-equiv=\"refresh\" content=\"5\"></head>")
        fmt.Fprintf(w, "<link rel=\"stylesheet\" href=\"files/bootstrap.css\">")
        fmt.Fprintf(w, "<body>")
        fmt.Fprintf(w, "<center>")
        nodes, err := ioutil.ReadDir("/git/docker/")
        if err != nil {
                log.Fatal(err)
        }
        fmt.Fprintf(w, "<table class=\"table table-hover\">")
        fmt.Fprintf(w, "<thead class=\"thead-dark\">")
        fmt.Fprintf(w, "<tr>")
        fmt.Fprintf(w, "<th scope=col></th>")
        fmt.Fprintf(w, "<th scope=col>IP</th>")
        fmt.Fprintf(w, "<th scope=col>Ping</th>")
        fmt.Fprintf(w, "<th scope=col>Containers</th>")
        fmt.Fprintf(w, "<th scope=col>Running</th>")
        fmt.Fprintf(w, "<th scope=col>Paused</th>")
        fmt.Fprintf(w, "<th scope=col>Stopped</th>")
        fmt.Fprintf(w, "<th scope=col>MemTotal</th>")
        fmt.Fprintf(w, "<th scope=col>SystemTime</th>")
        fmt.Fprintf(w, "</tr>")
        fmt.Fprintf(w, "</thead>")
        fmt.Fprintf(w, "<tbody>")

        for _, f := range nodes {
                if f.Name() != "env" {
                        result, time := pingTest(f.Name())
			online := false
			if result {
				online = dockerOnline(f.Name())
			}
                        if result && online {
                                fmt.Fprintf(w, "<tr>\n")
                                //fmt.Fprintf(w, "<tr style=\"background-color:#00FF00\">\n")
                        } else if result {
                                fmt.Fprintf(w, "<tr style=\"background-color:#F4D942\">\n")
			}else{
                                fmt.Fprintf(w, "<tr style=\"background-color:#FF0000\">\n")
                        }
                                fmt.Fprintf(w, "<td>\n")
                        fmt.Fprintf(w, "<a href=\"/ssh?ip=%s&command=sudo reboot\" class=\"btn btn-warning\">Reboot</a>\n", f.Name())
                        fmt.Fprintf(w, "<a href=\"/ssh?ip=%s&command=sudo poweroff\" class=\"btn btn-danger\">Poweroff</a>\n", f.Name())
                                fmt.Fprintf(w, "</td>\n")
                        fmt.Fprintf(w, "<td>%s</td>\n", f.Name())
                        fmt.Fprintf(w, "<td>%s</td>\n", time)
        //                fmt.Fprintf(w, "<td>")
                        if online {
			dockerStat(w, f.Name())
			}
          //              fmt.Fprintf(w, "</td>")
                        fmt.Fprintf(w, "</tr>")

                }
        }
        fmt.Fprintf(w, "</tbody>")
        fmt.Fprintf(w, "</table>")
}


