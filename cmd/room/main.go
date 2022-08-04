package main

import (
	"encoding/base64"
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
	aesData = map[string]map[string]string{
		"MIGeMA0GCSqGSIb3DQEBAQUAA4GMADCBiAKBgGFuKDQYmXVffJARKc5WI/JHDSj+AYMjiUAvBDFNNEtzXiM00lk6Mk5kE9zMwRLQ/UYBu/cXRgE+fANLkw6h2b8lGmyndPAgZNM5Md3Lg+EXEK9iUuRhfdn41CJ5vwmvXGeaRbkOVEEPxJTma+ajpHUQbqCkUgJrmIVpqvQbww07AgMBAAE=": {
			"AESKey": "CIVVs/B9NdfCnQs3jHunZ4UB0O7iphk2tBjMEKF+GflLOq4O//99x2IJb1LI6CkSh4UNBp/Mwsq+oisB+UMdnmFA1F+NKFy+mLg+BUgA8M4njRxzslaTSsR6pOVHMrBiYzBw0KMf2yYSCmLb6N5+LIxC7EIZDR8iImXUjS5ZuHs=",
			"AESIV":  "QioZf5pxexahIVwYiWOFeflbo+06KunySiUOAq0+RrtyNdBAsDEaqmZ3hfp7uVrJpuvbsYCRxkPIYacOWmoTCufiH5W8iaAyxPmDLynhIRpMWhpjmktR2iSznw7q+70QtF1pqD2nAPTik4IPMKa4ZHeEJayWZcJrDcw9JXnceFQ=",
		},
		"MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCsgt845/684lcmp1kdGFyoxmqzg2Gq2aC5Snmrq5Y4uR7RrGh+cPVVUYAfvQLsxD3lZmwDRZn6BsMv2tUZkl3ibC6oWFWHHKPVL9qkoeTxPSvMGqPNYMsw+nt9Io6WjhDJT5C0uLWnFYVN2QE0UbtIYnDowIndN30WhZKydcdYQQIDAQAB": {
			"AESKey": "nQnd4IhKQDZLyqAjAk6yFgh9kItqVrAfHN1El08DIjXx1aqGXqynTFadpEe8J/AskVWSUZa6fGG1Dvrf30mz50X/PtYAVJK3nrodvHBSq9cNKu9iJfYs3x0DfBjTqy4QfxXP/4VK1tOLX0YUyo715tykoXbwZU4S6HWuJDGd+dQ=",
			"AESIV":  "byqqEbHzj2eUgOqW2nOjFw3NDhhaF8LN0Hrik3wxDvnr/lwMnI5uB5gVTZkjTpyeVR48WhW4oIVZOzDDoVdeIhz9ki1rXqGrPvL+6a1GWi23hO461rwlsOogDHf3DF8V/08jZpf+1dkq+SpJWHDf2Tc5pvf975RDSIvYhqVGhlc=",
		},
		"MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDJs5TLCjWcOJlyiSbO1+wPkpGXlbOuNn7L0czOyUgkBSU/DSVcT1bv7l1Ja/hqHY0ntP5wpbtECpmCiBPFv80uB4nj7UFeWP7Jn4eSX7CB1HJIoJ3O9OcNH/Zg2xYQPeScLo/0X17DgBxjU8JuN4MZWmPTuKRzuWxkXO5kYCbfRwIDAQAB": {
			"AESKey": "eiMgr1FpzLivvkYHdFaDZMn/F5elPOxzqufI8VmAeKhRPA4Qp4Y5LerjKf0MdBSbF7axxbD+V6Jv78Z+Cd2Q+10BrXrKzTuWRBb2s/60EmADBEqKCAv3fWHiePzP5ZdU19Rmhfs889LewxSVOUlDDqx/vtP+v7tKHc3cTZDQeXk=",
			"AESIV":  "vbdxdaqxKMEGXBxj4/mWgcUMc13UH62o/Z+q4jOmiTE7Bzifwq955ncJun+R9BCvvkUcPhTt/f86TDcBnzOxkbxrSM8bbYM+i4ruuKTAHA80FPBmi2zQ+wbFDnu24e+TfJrsyDUfhZsI5Sy8mKdVFnwOgrRkDom7T3HFlWHjB7Y=",
		},
		"MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCoEFb1nI858+z7HFF8jejiWWQFvu1WhsVX7uLhbkOVPNPwIVG9Nz6A5hI65FR3suEGabAv6Qb+qLJoaAxtG77Y56UUGTfVexR41HCXNm72l7KvkdELjyEzkhkmdY2tHqXOI3zOz0ahU4+YOXHWnqLDsRYRawhxS0YtQY9QhhA4XwIDAQAB": {
			"AESKey": "cvzDqbfvhXLZAdXPrxpCkeOw8MomUqQqYvmvfSD6v283oEBwmscE5LJN1yEYJLIit698T30Q9zIATH2RqikvWseKBDdfLFQbiZs7BVRwqXPKQDMjV92pNTPOEYA15/DwMoS8yvlEYHjt1dBDfKejAVgdHQBSiBG+M+UbBd98tTg=",
			"AESIV":  "JfGWFeL/VgXXxNc6YRRjnwj2rOiJaXxmAdBAJFM5KYsWHga5jUMHZ8LxAkGVaMj3PIQ3yfp17iXEE3JDQZ95qLMmZfk7WSashT+muCLlFqSBCpT9yFHtmn/goCHGTeL9L+OR5vdWMb1vx8CKHtxGBgbgexMNcovsyQbK/0JRfqs=",
		},
		"MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCciFhUbwCGA4Tr2cXxrl8Ue91mBPZ1CBnkMKL1nwnx2WeUXQcCOB6XQuIR5jfqm8A3mY+pTiE1ngjNzE/1riFLmcI8cW+pSxGkJaSCBTNo175se8AKzH7mvpYPsb4+lcm8sEfZtkPBh9BuXuVg80378TRRiPPGYcL4LGL4aJf8zwIDAQAB": {
			"AESKey": "OZaEeeIIHOQWMvegTVLLks2PRfjuPUQ39C1GuB1tIRq2lVFvG+pWU+YfbeOq72vqvDeYSjnPlCp+tiSsTK3MkgmRDKiB7+EHC2FlHjU8v70daYQkE/+lJOwZQrYqfxGrSCPxRVWnPeQlzz/j4rIZRbeklvuEJVKVIPUVicYnb4k=",
			"AESIV":  "Jyo8O6NKJihOLq1TpvCHG0kwhJFypvdZig5fr6K40HSIFk1PQqDuywR7MlwLsgUQ/GmqWkD9q2vigsejCyO44CUzN5CwG6MBQb9FrB8nRM00d3KML5ClTEM74G3dYtAJQM2GXjoiT8WFcGeosIB7hW30qEgFK+IvvLT9I6rabCs=",
		},
	}
)

func getAesEnvelop(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

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
		"AESKey": base64.StdEncoding.EncodeToString(aesKey),
		"AESIV":  base64.StdEncoding.EncodeToString(aesIV),
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
