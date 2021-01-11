package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	InitRedis()
	e1 := &Entry{"Me", "NNN", "Nauthor", 9.0, "https://a.com", "This is synopsis", "The comments", "2020-1", "Bv"}
	estr1, _ := json.Marshal(e1)

	e2 := &Entry{"Me2", "NNN", "Nauthor", 9.0, "https://a.com", "This is synopsis", "The comments", "2020-1", "Bv"}
	estr2, _ := json.Marshal(e2)

	e3 := &Entry{"Me3", "NNN", "Nauthor", 9.0, "https://a.com", "This is synopsis", "The comments", "2020-1", "Bv"}
	estr3, _ := json.Marshal(e3)

	InsertEntry(string(estr1))
	InsertEntry(string(estr2))
	InsertEntry(string(estr3))

	RemoveEntry(string(estr3))
	fmt.Println(GetAll())

	err := http.ListenAndServe("0.0.0.0:80", &rdsHandler{})
	if err != nil {
		panic(err)
	}

	return
}
