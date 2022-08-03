package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"net/http"

	"github.com/godump/doa"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

var (
	aesData = map[string]map[string]string{}
)

func getAesEnvelop(w http.ResponseWriter, r *http.Request) {
	rsa := r.URL.Query().Get("rsaPublicKey")
	if v, ok := aesData[rsa]; ok {
		w.Write(doa.Try(json.Marshal(v)))
		return
	}
	aesKey := make([]byte, 32)
	rand.Read(aesKey)
	aesIV := make([]byte, 16)
	rand.Read(aesIV)
	ret := map[string]string{
		"AESKey": hex.EncodeToString(aesKey),
		"AESIV":  hex.EncodeToString(aesIV),
	}
	aesData[rsa] = ret
	w.Write(doa.Try(json.Marshal(ret)))
}

func main() {
	flag.Parse()
	hub := newHub()
	go hub.run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/getAesEnvelop", getAesEnvelop)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
