package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"path"
)

func uploaderHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.FormValue("user_id")
	file, header, err := r.FormFile("avatar_file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	filename := path.Join("avatars", userId+path.Ext(header.Filename))
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, "Successfull")
}
