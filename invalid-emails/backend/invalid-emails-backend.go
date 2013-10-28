package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"github.com/dmotylev/goproperties"
	"github.com/gorilla/mux"
	"io"
	"log"
	"log/syslog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
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
	producer        *Producer
	config          properties.Properties
	configPath      = flag.String("config", "", "path to configuration file")
	showVersion     = flag.Bool("version", false, "show application version and exit")
	showHelp        = flag.Bool("help", false, "show help")
	version         string
	applicationName string
	wg              sync.WaitGroup
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
	wg.Add(1)
	defer wg.Done()
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

func usage() {
	fmt.Printf("%s parameters:\n", applicationName)
	fmt.Println("\t--config=XXX\t - specify path to config file (required)")
	fmt.Println("\t--help\t\t - show this message")
	fmt.Println("\t--version\t - show application version and exit")
}

func checkArgs() {
	if *showVersion {
		fmt.Printf("%s version %s\n", applicationName, version)
		os.Exit(0)
	}
	if *showHelp {
		usage()
		os.Exit(0)
	}
	if *configPath == "" {
		usage()
		log.Fatal("Please specify path to configuration file")
	}
}

func main() {
	flag.Parse()

	w, err := syslog.New(syslog.LOG_INFO, applicationName)
	if err != nil {
		log.Fatalf("connecting to syslog: %s", err)
	}
	log.SetOutput(w)

	checkArgs()

	config, err = properties.Load(*configPath)

	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp", config.String("listen.addr", ""))
	if err != nil {
		log.Fatal(err)
	}

	producer = NewProducer(config)

	if err = producer.Connect(); err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	log.Printf("listening on %s%s", config.String("listen.addr", ""), config.String("listen.path", DefaultListenPath))
	router.HandleFunc(config.String("listen.path", DefaultListenPath), HandleRequest).
		Methods("POST")

	http.Handle("/", router)

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
		s := <-sig
		listener.Close()
		log.Printf("Got %s - exiting", s)
	}()

	if err = http.Serve(listener, nil); err != nil {
		log.Print(err)
	}
	log.Print("Waiting active requests for being finished")
	wg.Wait()
	producer.Stop()
	log.Print("Exited")
}
