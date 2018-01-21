package main

import (
    "testing"

	"net/http"
	"net/http/httptest"
	"strings"
	"os"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert" 
	"net/url"
	"log"
	// "time"

)

var (
	oNoteForm1 = map[string]string{ "note_id":"1", "note_address":"C:\\gotest", "note_name":"create_note" }
	oNoteForm2 = map[string]string{ "note_id":"1", "note_address":"C:\\gotest", "note_name":"change_note_name" }
	oPageForm1 = map[string]string{ 
		"note_address":"C:\\gotest", 
		"page_id":"1", 
		"page_title":"first_page_title", 
		"page_body":"first_page_body",
	 }

	 xNoteForm1 = map[string]string{ "note_id":"1", "note_address":"C:\\nothingpath", "note_name":"create_note" }

)


func init() {

	testFilePath := oNoteForm1["note_address"] + "\\note.db"

	// if err1 == nil{
	if err1 := osCheckFile( testFilePath ); err1 == nil {
		if err2 := os.Remove( testFilePath ); err2 != nil {
			log.Println( err2 )
		}
	}

	confFilePath := "data\\conf.db"
	// if err1 == nil{
	if err1 := osCheckFile( confFilePath ); err1 == nil {
		if err2 := os.Remove( confFilePath ); err2 != nil {
			log.Println( err2 )
		}
	}
	// time.Sleep( 10 * time.Second )

} //--------------------------------------------



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

	f := make(url.Values)
	f.Set("note_name", oNoteForm1["note_name"])
	f.Set("note_address", oNoteForm1["note_address"])

	req := httptest.NewRequest(echo.POST, "/AddNotePost", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)


	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, AddNotePost(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// assert.Equal(t, userJSON, rec.Body.String())
	}

}

func TestAddPage(t *testing.T) {

	// Setup
	e := echo.New()


	f := make(url.Values)

	f.Set("note_id", oNoteForm1["note_id"])
	f.Set("note_address", oNoteForm1["note_address"])

	req := httptest.NewRequest(echo.POST, "/AddPagePost", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)


	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, AddPagePost(c)) {
		assert.Equal(t,  http.StatusCreated, rec.Code)
		// assert.Equal(t, userJSON, rec.Body.String())
	}
}



func TestUpdateNote(t *testing.T) {

	// Setup
	e := echo.New()

	f := make(url.Values)

	f.Set("note_id", oNoteForm2["note_id"])
	f.Set("note_name", oNoteForm2["note_name"])


	req := httptest.NewRequest(echo.POST, "/UpdateNotePost", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)


	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, UpdateNotePost(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// assert.Equal(t, userJSON, rec.Body.String())
	}
}


func TestUpdatePage(t *testing.T) {

	// Setup
	e := echo.New()

	f := make(url.Values)
	f.Set("note_address" , oPageForm1["note_address"]	)
	f.Set("page_id"		 , oPageForm1["page_id"]		)
	f.Set("page_title"	 , oPageForm1["page_title"]		)
	f.Set("page_body"	 , oPageForm1["page_body"]		)


	req := httptest.NewRequest(echo.POST, "/UpdatePagePost", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)


	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, UpdatePagePost(c)) {
		assert.Equal(t, http.StatusCreated , rec.Code)
		// assert.Equal(t, userJSON, rec.Body.String())
	}
}

func TestDeletePage(t *testing.T) {

	// Setup
	e := echo.New()

	f := make(url.Values)

	f.Set("note_address", oPageForm1["note_address"])
	f.Set("page_id", oPageForm1["page_id"])


	req := httptest.NewRequest(echo.POST, "/DeletePagePost", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)


	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, DeletePagePost(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		// assert.Equal(t, userJSON, rec.Body.String())
	}
}


func TestDeleteNote(t *testing.T) {

	// Setup
	e := echo.New()

	f := make(url.Values)

	f.Set("note_id", oNoteForm1["note_id"])


	req := httptest.NewRequest(echo.POST, "/DeleteNotePost", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)


	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, DeleteNotePost(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// assert.Equal(t, userJSON, rec.Body.String())
	}
}
