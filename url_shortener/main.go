package main

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

// is this too much voodo for our purposes?
func getlinksmap() map[string]string {
	isxod := make(map[string]string)

	isx_js, err := json.Marshal(isxod)
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat("data.json"); os.IsNotExist(err) {
		err := os.WriteFile("data.json", isx_js, 0644)

		if err != nil {
			panic(err)
		}
	}

	isxod, isx_js = nil, nil

	data, err := os.ReadFile("data.json")
	if err != nil {
		panic(err)
	}

	isxod_ag := make(map[string]string)

	if err := json.Unmarshal(data, &isxod_ag); err != nil {
		panic(err)
	}

	return isxod_ag
}

func handlecreate(w http.ResponseWriter, req *http.Request) {
	linksmap := getlinksmap()

	for i, v := range linksmap {
		if i == req.PostFormValue("link") {
			io.WriteString(w, v)
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	link := req.PostFormValue("link")

	encd := base64.StdEncoding.EncodeToString([]byte(link))[:8]

	if len(link) < 8 || (link[0:8] != "https://" && link[0:7] != "http://") {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Bad request: not a link (http or https needed)")
		return
	}

	linksmap[req.PostFormValue("link")] = "/" + encd

	jsonlinks, err := json.Marshal(linksmap)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return
	}

	if err := os.WriteFile("data.json", jsonlinks, 0644); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "/"+encd)
}

func handlestuff(w http.ResponseWriter, req *http.Request) {
	if req.RequestURI == "/favicon.ico" {
		return
	}

	linksmap := getlinksmap()

	for i, v := range linksmap {
		if req.RequestURI == v {
			http.Redirect(w, req, i, http.StatusPermanentRedirect)
		}
	}
}

func main() {
	http.HandleFunc("/", handlestuff)
	http.HandleFunc("/createlink", handlecreate)
	if err := http.ListenAndServe("127.0.0.1:5555", nil); err != nil {
		panic(err)
	}
}
