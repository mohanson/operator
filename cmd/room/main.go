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
	aesData = map[string]map[string]string{
		"MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCC2zNRwzvLNwHCmoConDvnk/LAzh0O7XyGKrpc0cLNBp+HRX1BL05LPHN+F5mTfcK2aU2XhRSe+4eLh6Cdvrd4brKKgQomHlqK08mSgfqpV12Y/sa83fPcuUxjdn96n6Yac1X0NfE4A1ysTaSJykaj2U4QoPzV1AxG6YCDAmBxZwIDAQAB": {
			"AESKey": "EGZaXkrQuNwjt16ZpPfhdjjuNLju3Lq425Xu7E+Ml4gfpdjAFFJYw39hcLx0uisZ4oM+7knvQrcq/k+12CQH3hvnoXBNfuhd+MFxa+TEqmfPpVGD7mpMajFNJar+R8c5bXyz0ErfFFqhk/uZKYujf8WNmvQ9+yyOWOGITpclhbY=",
			"AESIV":  "ZR6ICXpy4ZIkEQwwVhEFP20uUJOlNeMt9yhd+FlGxfQcbbgd/4aQCtdQSKl0L+BNfHgMCMVqMnlPIZaHCd5GKZrF3A+ocYMwERMwsCVq7qleErdwtl1TNhc3tkvl0VMONGP9jxZj/z4DCtFQHr8WUQ3pQAxhkrqwUO9vBiIVMkA=",
		},
		"MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCIkbkbzfnX7jDhlwBhYBbQ5AGbeUapzKnZLEqG3mys8uRnoJVqfAf+dTeyIh7k/koupzU8XeJoSRvjpiwzYLr6qqpALRaPoThSfx6Gaos1rtOLQ5/MUlCbR5Pn2OyyXMlstb3d/mFZ/QXod5+GruLlbFnUjLR8orstQ4tXpYiV5QIDAQAB": {
			"AESKey": "ZhSa4YtsvGTZmpfHt5y2R9HtjPqmNeH7SRIH7X900FGDinKXtohFr48LrlwRQ6XDt3amdOR2VzBD/pnC3fT7JjqLxKudOcJ6a+76x+h3Tg31DHLRdiiidZvdvJkXnpKNfTlSh3oXVOV79RDwycEuRaQiCEe40xHVfeBmB/9BHEk=",
			"AESIV":  "UlZN/tIoEPlwWrvE4yjmxrIdNNeQnOUMQdr/0eLMiekpT3WI3zN78nOX9k1cY9Ke/1XnhmbfEs6v9lIByL6eFwh02B/f0TjJ43B5bgioGvsDvtZuxG/z/+ibw7rg3qrnFBoNZwRqnQTyuveWHYdLXlWuPpZhbTf1hQwSZ2Yaa/M=",
		},
		"MIGeMA0GCSqGSIb3DQEBAQUAA4GMADCBiAKBgHbnHwsRtNSSUyec8MtuENn73SS/RlTAzv4ZDk3PN7Rkt4+yyDw2hkw36CYGaY9bm/mPWWOEY8aiow3AprXPBi5ui/sCuBeaLK6Zcnm9CLWHo/BX53Wlsn9T9bCrIsoy6zNMZrOk/gDF90pS0EbCMkyCMkjbS2EiaNTRl8b6l0uZAgMBAAE=": {
			"AESKey": "cmYIfaff9rd1yMigNLXmA42vVEy+DRXVL9mQdeQ0vskzx2Ertxw2LxcoeU1uAPjrfXuX9DfoeeOVYToUvICmwXdzYg0IZa9NPei/Kl6EorhH9o6fKAes4RZUcQvrDc0gF70g8EqQ0d2IV8SQBErsYf/PCuiesBdrgUp2lIk11lw=",
			"AESIV":  "If1SVkB7dXKkX4qW5zjlqABBSm5zwSsCN0dJGRZ35n8/uIEqoC2N4E9wDgcBid38GcBQlGxEl4mFgZ3pf1duz/ntV1yqWiIkSOE03cE8CuTIrbLgcQ5lNMj0KeuRyeay1wqDjD04reTegxsePNKgudki9F0ig3o3S8GmO/gQl7s=",
		},
		"MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCW+HUNviIs5muuqPhA3mRtDxDB4xN6jGS32AQf59waoZagyUqNQigVJY3ICgKAOEjX8xOxmhzsw2Y3sCVR0Tg4jb0U2lPqz+kxbzkKgbJDa501rxmN/kMGOnx4DP/0lW+F54hSEjMXKTE3cXHHoXhcWs77l/C6af33hw6zJvCntQIDAQAB": {
			"AESKey": "eVWtk4TlYhnmAsPN5Ihi/MBIX80ej6ixI1aqRUrvtTa7K9gPu61iE3JSaRF+krD4tsik22iCduQbkM1ozplry+e5DCZ8mLR4ZOZQnP976y7Q5S6TR0xvrJtEU1dJKm/usQ85rXIAWrV3SG9oXP/d8f6DtzujN1qLu8J8csFt670=",
			"AESIV":  "EbXdkJ05hlvKW+mcFxp/Mzpcn0sZAPGp7QkURjyJwTmqDSidDCMuezg5u9vLLEom2yOXC2o7m4uCuYILnV8zPotgcyVJMNwUw3EOLpdrl4Xq9AdXvEdb/81waX75pkYNEY5Q19AqkdljOIc/62Q28T7ioBcY7ZxPejQ2e5jZ384=",
		},
		"MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDy76ynBNE84XniSlVQ7FL7h/qb0RbwsOz4XNqYB+GGnHgNOlTj2B/ZMmZNJ8kmcm39U4crilj/U0JWeboMsbAL0jwS9OVsFGTWvHNhRnJWmSjttP2nJLpPiS8KJw4mtjLKTXs9mtLGQ2FXxj+p63jU+O9XOypAPwg+YR+N1VlAUwIDAQAB": {
			"AESKey": "enTJmaR5jjfn2oZazlrTXtYUyeqzqryizEgyaKx98R5XVfy7pAemE3ce87Agi5rycw4TzAVb52Y3IdKh0u99NPXzfpJWWs8YLetetCuNn3lAsqc2iWawnGOfHr4iYBnNfaqqpOODupyJTC4N0TCzZi1Sgd23SywafsM4DQ9PgOs=",
			"AESIV":  "tNJdybc+cTAD1q6qqJ5Qoqzng9dl8/e4DUrhcGn1/Ab5qbpNhaGGyjeTCZ1klXDD6SWVIAbzPy8OkGxtBw9KoFbav43ivRWFvr2mrtDgb6VNZA6Hg0QAHnfoix5A+hTPyk8syoKmatDNoHM7ix9cRTii+jINyESfOPLZCML7ZSU=",
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
