// Module url provides functions to create, use, and retires wee Urls.
package wee

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Shortener struct {
	repo   *Repository
	logger *log.Logger
}

func NewShortener(r *Repository, l *log.Logger) *Shortener {
	return &Shortener{
		repo:   r,
		logger: l}
}

// FollowUrl expands the weeUrl and redirects the client browser to the full URL.
func (s *Shortener) FollowUrl(c *gin.Context) {
	// this is where we redirect...
	wee := c.Param("weeUrl")

	// THIS SHOULD BE UNNECESSARY if the main routing r.GET pattern is correct
	if wee == "" {
		s.logger.Printf("found redirect route, ROUTER LOGIC IS WRONG")
		c.File("./public/index.html")
		return
	}
	rec, err := s.repo.find(wee)
	if err != nil {
		// this will be a common situation, anytime the weeUrl is mistyped (or just junk)
		s.logger.Printf("could not find record for %s, %v\n", wee, err)
		c.JSON(http.StatusNotFound, gin.H{
			"reason": err.Error(),
		})
		return
	}
	s.logger.Printf("redirecting to %s\n", rec.Url)
	c.Redirect(http.StatusTemporaryRedirect, rec.Url)
}

// ShortenUrl accepts a full URL string, creates record for it, and saves it.
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
		"weeUrl": rec.Tag,
		"token":  rec.Token,
	})
}

// LengthenUrl expands a weeUrl and returns it for display but does not redirect to it.
func (s *Shortener) LengthenUrl(c *gin.Context) {
	var rec *Record
	wee := c.Param("weeUrl")

	rec, err := s.repo.find(wee)
	if err != nil {
		s.logger.Printf("Error, could not find a record for wee URL %s: %v\n", wee, err)
		c.JSON(http.StatusNotFound, gin.H{
			"reason": err.Error(),
		})
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"url": rec.Url,
	})
}

// RetireUrl accepts token, then locates the corresponding record and removes it.
func (s *Shortener) RetireUrl(c *gin.Context) {
	token := c.Param("token")

	err := s.repo.remove(token)
	if err != nil {
		s.logger.Printf("Error, could not remove a record for token %s: %v\n", token, err)
		c.JSON(http.StatusNotFound, gin.H{
			"reason": err.Error(),
		})
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
	})
}
