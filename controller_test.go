package main

import (
    "testing"

	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert" 
	"net/url"
)


func TestLiveCheck(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/livecheck", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, HandleLiveCheckGet(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// assert.Equal(t, userJSON, rec.Body.String())
	}
	// if assert.NoError(t, h.getUser(c)) {
	// 	assert.Equal(t, http.StatusOK, rec.Code)
	// 	assert.Equal(t, userJSON, rec.Body.String())
	// }
}


func TestAddNote(t *testing.T) {

	// Setup
	e := echo.New()

	userForm := make(map[string]string)
	userForm["new_note_name"]    = "forTest"
    userForm["new_note_address"] = "C:\\gotest"

	f := make(url.Values)
	f.Set("new_note_name", userForm["new_note_name"])
	f.Set("new_note_address", userForm["new_note_address"])

	req := httptest.NewRequest(echo.POST, "/HandleAddNotePost", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)


	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, HandleAddNotePost(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// assert.Equal(t, userJSON, rec.Body.String())
	}
	// if assert.NoError(t, h.getUser(c)) {
	// 	assert.Equal(t, http.StatusOK, rec.Code)
	// 	assert.Equal(t, userJSON, rec.Body.String())
	// }
}

func TestAddPage(t *testing.T) {

	// Setup
	e := echo.New()

	userForm := make(map[string]string)
	userForm["note_id"]    = "6"
    userForm["note_address"] = "C:\\gotest"

	f := make(url.Values)
	f.Set("note_id", userForm["note_id"])
	f.Set("note_address", userForm["note_address"])

	req := httptest.NewRequest(echo.POST, "/HandleAddPagePost", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)


	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, HandleAddPagePost(c)) {
		assert.Equal(t,  http.StatusCreated, rec.Code)
		// assert.Equal(t, userJSON, rec.Body.String())
	}
	// if assert.NoError(t, h.getUser(c)) {
	// 	assert.Equal(t, http.StatusOK, rec.Code)
	// 	assert.Equal(t, userJSON, rec.Body.String())
	// }
}
