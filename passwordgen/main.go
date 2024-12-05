package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
)

func genpass(w http.ResponseWriter, req *http.Request) {
	lenStr := req.URL.Query().Get("len")
	lenInt, err := strconv.Atoi(lenStr)
	if err != nil {
		panic(err)
	}

	alphabetaZ := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers09 := "0123456789"
	specialChars := "!@#$%^&*()_-+=<>?/{}~"
	allChars := []rune(alphabetaZ + numbers09 + specialChars)

	if lenInt <= 3 || lenInt >= 100 {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Your password must be at least 3 characters long and 100 characters long at most\n")
		return
	}

	pass := []rune{}

	for i := 0; i < lenInt; i++ {
		randInd := rand.Intn(len(allChars))
		pass = append(pass, allChars[randInd])
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, fmt.Sprintf("%s\n", string(pass)))
}

func main() {
	http.HandleFunc("/genpass", genpass)
	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		panic(err) // who needs log.Fatal()?
	}
}
