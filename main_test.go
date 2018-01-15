package main

import (
    "testing"

	"net/http"
	"net/http/httptest"
	// "strings"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert" 
)

// func TestElephantFeed(t *testing.T) {
//     expect := "Grass"
//     actual := "Grass"

//     if expect != actual {
//         t.Errorf("%s != %s", expect, actual)
//     }
// }

func TestGetItem(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// c.SetPath("/users/:email")
	// c.SetParamNames("email")
	// c.SetParamValues("jon@labstack.com")
	// h := &handler{mockDB}

	// Assertions
	if assert.NoError(t, HandleLoadPageGet(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// assert.Equal(t, userJSON, rec.Body.String())
	}
	// if assert.NoError(t, h.getUser(c)) {
	// 	assert.Equal(t, http.StatusOK, rec.Code)
	// 	assert.Equal(t, userJSON, rec.Body.String())
	// }
}

