package main

import (
	"fmt"
	"net/http"
)

type TransRes struct {
	status int
	msg    string
}

type rdsHandler struct {
}

// ServeHTTP : handle http request
func (hdl *rdsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var trs = TransRes{0, ""}

	defer func() {
		fmt.Fprintf(w, "Status:%d\n", trs.status)
		fmt.Fprintf(w, trs.msg)
	}()

	r.ParseForm()
	action := r.FormValue("Action")
	entrystr := r.FormValue("Entry")

	switch action {
	case "Insert":
		InsertEntry(entrystr)
	case "Remove":
		RemoveEntry(entrystr)
	case "Recoever":
		RecoverEntry(entrystr)
	case "GetAll":
		trs.msg = GetAll()
		trs.status = 1
	default:
		fmt.Println("Unknown action!")
	}
}
