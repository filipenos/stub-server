package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

//Variables assigned by the user
var (
	port = flag.Int("port", 8090, "Port that server run")
	conf = flag.String("conf", "conf.json", "JSON file with configs")
)

var confs []Config

//Config represents a configuration of handler
type Config struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Code   int    `json:"code"`
	Body   string `json:"body"`
	Type   string `json:"type"`
}

func init() {
	flag.Parse()
}

func main() {
	f, err := os.Open(*conf)
	if err != nil {
		log.Fatalf("Error opening file %s", *conf)
	}

	err = json.NewDecoder(f).Decode(&confs)
	if err != nil {
		log.Fatalf("Error reading file %v", *conf)
	}

	log.Printf("Stub Server running at :%d with conf file %s", *port, *conf)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		served := false
		for _, conf := range confs {
			if r.URL.Path == conf.Path {
				conf.serve(w, r)
				served = true
			}
		}
		if !served {
			http.Error(w, "Not found", http.StatusNotFound)
		}
	})
	err = http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err == nil {
		log.Fatal("Failed to run server: ", err)
	}
}

//serve func that serve path with a conf
func (conf *Config) serve(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving %v %v with %#v", r.Method, r.URL, conf)
	if r.Method != conf.Method {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	} else {
		if conf.Code != 0 && conf.Code != http.StatusOK {
			w.WriteHeader(conf.Code)
		}
		switch conf.Type {
		case "":
			w.Write([]byte(conf.Body))
		case "json":
			w.Header().Add("Content-Type", "application/json")
			j, err := os.Open(conf.Body)
			if err != nil {
				log.Fatalf("Error on open file %v", err)
			}
			b, err := ioutil.ReadAll(j)
			if err != nil {
				log.Fatalf("Error on read file %v", err)
			}
			w.Write(b)
		}
	}
}
