package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

var random string

func PongHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-HelloHttp-Instance", random)
	w.Write([]byte("PONG"))
}

func LogRequestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("")
	fmt.Println("Proto", r.Proto)
	fmt.Println("TransferEncoding", r.TransferEncoding)
	fmt.Println("Close", r.Close)
	fmt.Println("Host", r.Host)
	fmt.Println("RemoteAddr", r.RemoteAddr)
	for k, v := range r.Header {
		fmt.Println("Header", k, v)
	}

	w.Header().Set("X-HelloHttp-Instance", random)
	w.Write([]byte("PONG"))
}

func ClientHandler(w http.ResponseWriter, r *http.Request) {
	urlStr := r.Header.Get("X-Req-URL")
	if urlStr == "" {
		w.WriteHeader(400)
		w.Write([]byte("missing X-Req-URL"))
		return
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("X-HelloHttp-Instance", random)
	httputil.NewSingleHostReverseProxy(u).ServeHTTP(w, r)
}

func init() {
	bs := make([]byte, 4)
	rand.Read(bs)
	random = hex.EncodeToString(bs)
}

func main() {
	for _, env := range os.Environ() {
		fmt.Println(env)
	}

	http.HandleFunc("/", PongHandler)
	http.HandleFunc("/log", LogRequestHandler)
	http.HandleFunc("/req", ClientHandler)

	fmt.Println("listening on 3000")
	http.ListenAndServe(":3000", nil)
}
