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

//Config represents a configuration of handler
type Config struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Json   string `json:"json"`
}

func init() {
	flag.Parse()
}

func main() {
	f, err := os.Open(*conf)
	if err != nil {
		log.Fatalf("Error opening file %s", *conf)
	}

	var confs []Config
	err = json.NewDecoder(f).Decode(&confs)
	if err != nil {
		log.Fatalf("Error reading file %v", *conf)
	}

	log.Printf("Stub Server running at :%d with conf file %s", *port, *conf)
	mux := http.NewServeMux()

	for _, conf := range confs {
		mux.HandleFunc(conf.Path, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != conf.Method {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			} else {
				w.Header().Add("Content-Type", "application/json")
				j, err := os.Open(conf.Json)
				if err != nil {
					log.Fatalf("Error on open file %v", err)
				}
				b, err := ioutil.ReadAll(j)
				if err != nil {
					log.Fatalf("Error on read file %v", err)
				}
				w.Write(b)
			}
		})
	}

	err = http.ListenAndServe(fmt.Sprintf(":%d", *port), mux)
	if err == nil {
		log.Fatal("Failed to run server: ", err)
	}
}
