package main

import (
	"fmt"
	"github.com/go-git/go-git"
	"github.com/go-git/go-git/plumbing/object"
	"net/http"
)

func gitWeb(w http.ResponseWriter, req *http.Request) {
	menu(w, req)
	r, err := git.PlainOpen("/git/")
	if err != nil {
		fmt.Printf("plain open :%s", err)
	}
	/*
		work, err := r.Worktree()
		if err != nil {
			fmt.Printf("worktree error :%s", err)
		}*/
	ref, err := r.Head()
	if err != nil {
		fmt.Printf("head error :%s", err)
	}
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		fmt.Printf("log error :%s", err)
	}

	fmt.Fprintf(w, "<pre><code>\n")
	err = cIter.ForEach(func(c *object.Commit) error {
		//fmt.Println(string(c.String()))
		fmt.Fprintf(w, "%s\n", c.String())

		return nil
	})
	fmt.Fprintf(w, "</code></pre>\n")

	if err != nil {
		fmt.Printf("cIter.ForEach error :%s", err)
	}

	/*
	   	value, _ := sjson.Set(ignitionFile, "storage.files.0.contents.source", "http://"+req.Host+"/files/stat")
	           w.Header().Set("Content-Type", "application/json")
	           w.Write([]byte(value))
	           log.Println("Server PXE booted:", req.Host)
	*/
}
