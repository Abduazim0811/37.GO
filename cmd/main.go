package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []Album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
}

func main() {
	// Load configuration
	var config Config
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	router.GET("/albums", getAlbums)

	router.Run(":" + strconv.Itoa(config.Server.Port))
}
func getAlbums(c *gin.Context) {
	title := c.Query("title")
	artist := c.Query("artist")
	price := c.Query("price")

	var filteredAlbums []Album
	for _, album := range albums {
		if (title == "" || strings.Contains(strings.ToLower(album.Title), strings.ToLower(title))) &&
			(artist == "" || strings.Contains(strings.ToLower(album.Artist), strings.ToLower(artist))) &&
			(price == "" || fmt.Sprintf("%.2f", album.Price) == price) {
			filteredAlbums = append(filteredAlbums, album)
		}
	}

	c.IndentedJSON(http.StatusOK, filteredAlbums)
}
