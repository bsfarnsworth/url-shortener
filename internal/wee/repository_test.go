package wee

import (
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	tType = "sqlite3"
	tData = ":memory:"
	
	tag1 = "NOMATCH"
	tag2 = "TestTag1"

	token1 = "NOTATOKEN"
)

var (
	tRepo *Repository
	tLogger *log.Logger
)

func TestSetup(t *testing.T) {
	// for sharing amongst the tests
	tLogger = log.New(os.Stderr, "[repository test] ", log.Lshortfile)
	tRepo = nil
}

func TestConnect(t *testing.T) {
	require := require.New(t)

	tRepo = NewRepository(tType, tData, tLogger)
	err := tRepo.Connect()
	require.Nil(err, fmt.Sprintf("Repo failed to connect: %v", err))
	require.NotNil(tRepo, fmt.Sprintf("Repo connection is nil."))
}

func TestCreate(t *testing.T) {
	require := require.New(t)

	require.NotNil(tRepo, fmt.Sprintf("Repo connection is nil."))

	// create should never fail since Connect already created the table
	err := tRepo.create()
	require.Nil(err, fmt.Sprintf("Repo failed to create table: %v", err))
}

func TestAdd(t *testing.T) {
	require := require.New(t)

	require.NotNil(tRepo, fmt.Sprintf("Repo connection is nil."))

	var tRec Record
	tRec.Tag = "TestTag1"
	tRec.Url = "TestURL1"
	tRec.Token = "TestToken1"
	err := tRepo.add(&tRec)
	require.Nil(err, fmt.Sprintf("Repo add operation failed: %v", err))
}

func TestFind(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	require.NotNil(tRepo, fmt.Sprintf("Repo connection is nil."))

	var err error
	var pRec *Record
	
	// relying on successful result of TestAdd

	pRec, err = tRepo.find(tag1)
	assert.NotNil(err, fmt.Sprintf("False match on tag %s failed to cause error", tag1))
	if err == nil {
		assert.NotEqual(tag1, pRec.Tag, fmt.Sprintf("Unexpected false match on tag %s", tag1))
	}

	pRec, err = tRepo.find(tag2)
	assert.Nil(err, fmt.Sprintf("Unexpected error finding tag %s, %v", tag2, err))
	if err == nil {
		assert.Equal(tag2, pRec.Tag, fmt.Sprintf("Failure to find tag %s", tag2))
	}
}

func TestRemove(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	var err error
	var pRec *Record

	// still relying on successful result of TestAdd

	require.NotNil(tRepo, fmt.Sprintf("Repo connection is nil."))
	pRec, err = tRepo.find(tag2)
	require.Nil(err, fmt.Sprintf("Unexpected error on tag %s, %v", tag2, err))

	err = tRepo.remove(token1)
	assert.NotNil(err, fmt.Sprintf("Unexpected success in removing a non-existent record for token %s, %v", token1, err))
	
	err = tRepo.remove(pRec.Token)
	assert.Nil(err, fmt.Sprintf("Unexpected error when removing record for tag %s, %v", tag2, err))

	// just to be sure...
	pRec, err = tRepo.find(tag2)
	assert.NotNil(err, fmt.Sprintf("Unexpectedly found record that was removed for tag %s, %v", tag2, err))
}

func TestDisconnect(t *testing.T) {
	assert := assert.New(t)

	err := tRepo.Disconnect()
	assert.Nil(err, fmt.Sprintf("Repo reported error when disconnecting: %v", err))
}

