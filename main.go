package main

import (
	"consistent-hashing/config"
	"consistent-hashing/hash"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

var hashRing = hash.NewHashRing()

func main() {

	for _, server := range config.BootstrapServers {
		hashRing.AddServer(server)
	}

	http.HandleFunc("/add-server", AddServer)
	http.HandleFunc("/remove-server", RemoveServer)
	http.HandleFunc("/get-data", GetData)
	http.HandleFunc("/add-data", AddData)

	log.Fatal(http.ListenAndServe(":3030", nil))
}

func AddServer(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	server := struct {
		ServerName string
	}{}

	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "error while parsing data", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &server)

	if err != nil {
		http.Error(w, "error while parsing data", http.StatusInternalServerError)
		return
	}

	err = hashRing.AddServer(server.ServerName)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("server added to the hash ring"))
}

func RemoveServer(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	server := struct {
		ServerName string
	}{}

	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "error while parsing data", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &server)

	if err != nil {
		http.Error(w, "error while parsing data", http.StatusInternalServerError)
		return
	}

	err = hashRing.RemoveServer(server.ServerName)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("server removed from the hash ring"))
}

func GetData(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	query := r.URL.Query()

	val, err := hashRing.GetValue(query.Get("key"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := json.Marshal(val)

	if err != nil {
		http.Error(w, "error while fetching data", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(res))
}

func AddData(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data := struct {
		Key   string
		Value interface{}
	}{}

	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "error while parsing data", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &data)

	if err != nil {
		http.Error(w, "error while parsing data", http.StatusInternalServerError)
		return
	}

	added := hashRing.AddData(data.Key, data)

	if !added {
		http.Error(w, "error while adding new data", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("data added successfully"))
}
