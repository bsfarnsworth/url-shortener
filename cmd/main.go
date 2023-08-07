package main

import (
	"log"
	"os"
	
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
	
	// Set the router as the default one shipped with Gin.
	// React pages will be served with SPA server pattern, discussed in
        //   https://github.com/gin-gonic/gin/issues/3109
	router := gin.Default()
	router.SetTrustedProxies(nil)

	// To keep the minimal path for the shortener's Follow function
	// accept exact "/" first but "/<url>" next;
	// NoRoute will route everything else into React;

	router.GET("/:weeUrl", shorty.FollowUrl)

	router.GET("/api/v1/lengthen/:weeUrl", shorty.LengthenUrl)
	router.GET("/api/v1/retire/:token", shorty.RetireUrl)
	router.POST("/api/v1/shorten", shorty.ShortenUrl)

	router.NoRoute(func(c *gin.Context) {
		c.File("./public/index.html")
	})

	// Start and run the server
	router.Run(port)
}
