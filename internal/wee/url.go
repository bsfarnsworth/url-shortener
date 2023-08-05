package wee

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	testDest = "www.nytimes.com/games"
)


// TBD not sure if I can do this,
// I want a this type so it can provide log & db as context
// will it conflict with the gin Context?

type Shortener struct {
	repo  *Repository
	logger *log.Logger
}

func NewShortener(r *Repository, l *log.Logger) *Shortener {
	return &Shortener{
		repo: r,
		logger: l}
}

func (s *Shortener) FollowUrl(c *gin.Context) {
	// this is where we redirect...
	wee := c.Param("weeUrl")
	rec, err := s.repo.find(wee)
	if err != nil {
		if testDest != "" {
			c.Redirect(http.StatusTemporaryRedirect, testDest)
			return
		}
		c.JSON(http.StatusNotFound, gin.H{
			"reason": err.Error(),
		})
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, rec.Url)
}

func (s *Shortener) ShortenUrl(c *gin.Context) {
	url := c.PostForm("fullUrl")
	rec := createRecord(url)
	err := s.repo.add(rec)
	if err != nil {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusInternalServerError, gin.H{
			"reason": err.Error(),
		})
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"weeUrl":   rec.Tag,
		"token": rec.Token,
	})
}

func (s *Shortener) LengthenUrl(c *gin.Context) {
	wee := c.Param("weeUrl")
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"url": "https://https://www.nytimes.com/crosswords/game/mini/" + wee,
	})
}

func (s *Shortener) RetireUrl(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"url": "gone",
	})
}
