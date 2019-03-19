package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"time"
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

func SizeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-HelloHttp-Instance", random)

	byteSizeStr := r.URL.Query().Get("byte_size")
	if byteSizeStr == "" {
		w.WriteHeader(400)
		w.Write([]byte("missing byte_size query var"))
	}

	byteSize, err := strconv.Atoi(byteSizeStr)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("strconv.Atoi failed"))
	}

	for i := 0; i < byteSize; i++ {
		w.Write([]byte{byte(97 + i%26)})
	}
}

func DurationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-HelloHttp-Instance", random)

	durationStr := r.URL.Query().Get("duration")
	if durationStr == "" {
		w.WriteHeader(400)
		w.Write([]byte("missing duration query var"))
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("time.ParseDuration failed"))
	}

	time.Sleep(duration)
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
	http.HandleFunc("/client", ClientHandler)
	http.HandleFunc("/size", SizeHandler)
	http.HandleFunc("/duration", DurationHandler)

	fmt.Println("listening on 3000")
	http.ListenAndServe(":3000", nil)
}