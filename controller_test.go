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
	if assert.NoError(t, LiveCheckGet(c)) {
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

	req := httptest.NewRequest(echo.POST, "/AddNotePost", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)


	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, AddNotePost(c)) {
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

	req := httptest.NewRequest(echo.POST, "/AddPagePost", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)


	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, AddPagePost(c)) {
		assert.Equal(t,  http.StatusCreated, rec.Code)
		// assert.Equal(t, userJSON, rec.Body.String())
	}
	// if assert.NoError(t, h.getUser(c)) {
	// 	assert.Equal(t, http.StatusOK, rec.Code)
	// 	assert.Equal(t, userJSON, rec.Body.String())
	// }
}



func TestUpdateNote(t *testing.T) {

	// Setup
	e := echo.New()

	userForm := make(map[string]string)
	userForm["note_name"]    = "forxxTest"
    userForm["note_id"]      = "6"

	f := make(url.Values)
	f.Set("note_name", userForm["note_name"])
	f.Set("note_id", userForm["note_id"])

	req := httptest.NewRequest(echo.POST, "/UpdateNotePost", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)


	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, UpdateNotePost(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// assert.Equal(t, userJSON, rec.Body.String())
	}
	// if assert.NoError(t, h.getUser(c)) {
	// 	assert.Equal(t, http.StatusOK, rec.Code)
	// 	assert.Equal(t, userJSON, rec.Body.String())
	// }
}


func TestUpdatePage(t *testing.T) {

	// Setup
	e := echo.New()

	userForm := make(map[string]string)
	userForm["note_address"]    = "C:\\gotest"
    userForm["page_id"]         = "1"
	userForm["page_title"]      = "testTitzzzle"
    userForm["page_body"]       = "testBozzzdy"

	f := make(url.Values)
	f.Set("note_address" , userForm["note_address"])
	f.Set("page_id"		 , userForm["page_id"])
	f.Set("page_title"	 , userForm["page_title"])
	f.Set("page_body"	 , userForm["page_body"])

	req := httptest.NewRequest(echo.POST, "/UpdatePagePost", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)


	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, UpdatePagePost(c)) {
		assert.Equal(t, http.StatusCreated , rec.Code)
		// assert.Equal(t, userJSON, rec.Body.String())
	}
	// if assert.NoError(t, h.getUser(c)) {
	// 	assert.Equal(t, http.StatusOK, rec.Code)
	// 	assert.Equal(t, userJSON, rec.Body.String())
	// }
}
