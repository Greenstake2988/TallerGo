package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	//
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albums slice to seed record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func getAlbums(c *gin.Context) {

	db, err := gorm.Open(sqlite.Open("albums.sqlite"), &gorm.Config{})

	if err != nil {
		return
	}

	var albums []album
	result := db.Find(&albums)

	if result.Error != nil {
		return
	}
	c.IndentedJSON(200, albums)
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {

	db, err := gorm.Open(sqlite.Open("albums.sqlite"), &gorm.Config{})

	if err != nil {
		return
	}

	var newAlbum album

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Add the new album to the slice.
	//albums = append(albums, newAlbum)
	db.Create(&newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	db, err := gorm.Open(sqlite.Open("albums.sqlite"), &gorm.Config{})

	if err != nil {
		return
	}
	var album album

	err = db.First(&album, id).Error

	if err != nil {
		return
	}
	c.IndentedJSON(http.StatusOK, &album)
}

func main() {
	db, err := gorm.Open(sqlite.Open("albums.sqlite"), &gorm.Config{})

	if err != nil {
		return
	}

	db.AutoMigrate(&album{})

	db.Create(&albums)

	r := gin.Default()
	r.GET("/albums", getAlbums)
	r.POST("/albums", postAlbums)
	r.GET("/albums/:id", getAlbumByID)

	r.Run() // listen and serve on 0.0.0.0:8080
}
