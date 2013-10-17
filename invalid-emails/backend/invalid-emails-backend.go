package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"github.com/dmotylev/goproperties"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	//"time"
)

const (
	DefaultListenAdr  = ":60080"
	DefaultListenPath = "/"
)

const (
	StatusInvalidSignature  int = 453
	StatusInvalidDataFormat int = 452
)

var (
	producer   *Producer
	config     properties.Properties
	configPath = flag.String("config", "config.cfg", "path to configuration file")
)

func calcMd5(val string) string {
	h := md5.New()
	io.WriteString(h, val)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func isRequestValid(data, hash, salt, secret string) bool {
	return hash == calcMd5(calcMd5(data+secret)+salt)
}

func HandleRequest(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-type", "text/html")
	data := req.FormValue("data")

	if !isRequestValid(
		data,
		req.FormValue("hash"),
		req.FormValue("time"),
		config.String("security.key", ""),
	) {
		http.Error(w, "Invalid signature", StatusInvalidSignature)
		return
	}

	if err := producer.PostTask(data); err != nil {
		http.Error(w, "Failed to publish task", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	io.WriteString(w, "Task successfuly registered")
	io.WriteString(w, "\n")
}

func main() {
	flag.Parse()

	var err error
	config, err = properties.Load(*configPath)

	if err != nil {
		log.Fatal(err)
	}

	producer = NewProducer(config)

	if err = producer.Connect(); err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	router.HandleFunc(config.String("listen.path", DefaultListenPath), HandleRequest).
		Methods("POST")

	http.Handle("/", router)

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
		s := <-sig
		producer.Stop()
		log.Printf("Got %s - exiting", s)
		os.Exit(1)
	}()

	if err = http.ListenAndServe(config.String("listen.addr", ""), nil); err != nil {
		log.Fatal(err)
	}
}
