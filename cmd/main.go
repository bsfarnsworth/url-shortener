package main

import (
	"log"
	"net/http"
	"os"
	
	//  "github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"url-shortener/internal/wee"
)

const (
	// bring these in from config instead of defining here
	dbtype = "sqlite3"
	dbfile = "./WeeRepo.db"	// this needs to reside on a volume in deploy machnine
	port = ":3000"
)

func main() {
	// Provide the same logger to all
	logger := log.New(os.Stderr, "[wee] ", log.Lshortfile)

	// Connect to our repository of URLs
	repo := wee.NewRepository(dbtype, dbfile, logger)
	err := repo.Connect()
	defer repo.Disconnect()
	if err != nil {
		logger.Printf("Fatal unable to connect to repository %s\n", dbfile)
		os.Exit(1)
	}
	
	// Shortener is the toolkit
	shorty := wee.NewShortener(repo, logger)
	
	// Set the router as the default one shipped with Gin
	router := gin.Default()
	
	// Serve frontend static files
	//	router.Use(static.Serve("/", static.LocalFile("./views", true)))
	
	// Setup route group for the API
	api := router.Group("")
	{
		api.POST("/api/v1/shorten", shorty.ShortenUrl)
		api.GET("/api/v1/lengthen/:weeUrl", shorty.LengthenUrl)
		api.GET("/api/v1/retire/:token", shorty.RetireUrl)

		api.GET("/:weeUrl", shorty.FollowUrl)

		// wth? unless useful as echo nuke this
		// -> we want to direct this to home page
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H {
				"message": "pong",
				})
		})
	}
	
	// Start and run the server
	router.Run(port)
}
