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
	// fmt.Println(r)
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
		trs.status = InsertEntry(entrystr)
	case "Remove":
		trs.status = RemoveEntry(entrystr)
	case "Recoever":
		trs.status = RecoverEntry(entrystr)
	case "GetAll":
		trs.msg = GetAll()
		trs.status = 1
	default:
		fmt.Println("Unknown action!")
	}
}
