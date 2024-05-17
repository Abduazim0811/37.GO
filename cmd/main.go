package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`
}

type Album struct {
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums []Album

func loadConfig(filename string) (Config, error) {
	var config Config
	file, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	return config, err
}

func loadAlbums(filename string) ([]Album, error) {
	var albums []Album
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&albums)
	return albums, err
}

func filterAlbums(query string) []Album {
	var result []Album
	for _, album := range albums {
		if strings.Contains(strings.ToLower(album.Title), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(album.Artist), strings.ToLower(query)) ||
			strings.Contains(fmt.Sprintf("%.2f", album.Price), query) {
			result = append(result, album)
		}
	}
	return result
}

func albumsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("filter")
	var result []Album
	if query != "" {
		result = filterAlbums(query)
	} else {
		result = albums
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func main() {
	config, err := loadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	albums, err = loadAlbums("albums.json")
	if err != nil {
		log.Fatalf("Error loading albums: %v", err)
	}

	http.HandleFunc("/albums", albumsHandler)
	address := fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port)
	log.Printf("Server starting on %s\n", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
