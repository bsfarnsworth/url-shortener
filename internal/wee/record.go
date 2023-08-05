package wee

import (
	"math/rand"

	"github.com/google/uuid"
)

// design constraints
const (
	tagSize = 7
)

// Record operations
// (should be independent of the db/storage, other than the WeeRec definition)

func token() string {
	return "untoken"
}

func weeUrl() string {
	return "unWeeUrl"
}

func fullUrl() string {
	return "unFullUrl"
}

func createRecord(fullUrl string) *Record {
	var rec Record
	
	// sanitize the Url
	// TBD
	rec.Url = fullUrl
	
	// create a unique tag
	var syms = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
 
    s := make([]rune, tagSize)
    for i := range s {
        s[i] = syms[rand.Intn(len(syms))]
    }
	rec.Tag = string(s)
	
	// create a token
	rec.Token = uuid.New().String()
	
	return &rec
}
