package main

import (
    "testing"
    "runtime"
    "io/ioutil"

	"net/http"
	"net/http/httptest"
	"strings"
	"os"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert" 
	"net/url"
	"log"
	"time"
	"encoding/json"

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

	// oExpectedString = "\"{\\\"RtnCode\\\":\\\"0\\\",\\\"DataSet\\\":[{\\\"NoteDBID\\\":1,\\\"NoteDBName\\\":\\\"create_note\\\",\\\"NoteDBAddress\\\":\\\"C:\\\\\\\\gotest\\\",\\\"NoteDBUpdateTime\\\":\\\"2018-01-22T08:13:03.2009784+09:00\\\",\\\"list\\\":[{\\\"ID\\\":1,\\\"CreatedAt\\\":\\\"2018-01-22T08:13:03.1869419+09:00\\\",\\\"UpdatedAt\\\":\\\"2018-01-22T08:13:03.1869419+09:00\\\",\\\"DeletedAt\\\":null,\\\"page_title\\\":\\\"\\\",\\\"page_body\\\":\\\"\\\"}]}],\\\"SlctPst\\\":{\\\"NoteID\\\":1,\\\"PageID\\\":1}}\""

)


func init() {

    // ディレクトリセパレータ
    if runtime.GOOS == "windows" {
        directorySeparator = "\\"
    }

    // 設定ファイルの読み込み
    file, err := ioutil.ReadFile( dataDirName + directorySeparator + "config.json" )
    if err != nil {
        panic(err)
    }

    json.Unmarshal(file, &userConfig)



	//テスト実行前にファイルを削除
    testDBAddress := oNoteForm1["note_address"] + "\\note.db"

	// if err1 == nil{
	if err1 := osCheckFile( testDBAddress ); err1 == nil {
		if err2 := os.Remove( testDBAddress ); err2 != nil {
			log.Println( err2 )
		}
	}

    confDBAddress := dataDirName + directorySeparator + confDBName

	// if err1 == nil{
	if err1 := osCheckFile( confDBAddress ); err1 == nil {
		if err2 := os.Remove( confDBAddress ); err2 != nil {
			log.Println( err2 )
		}
	}
	time.Sleep( 5 * time.Second )



    loadTemplates()
    checkConfig()
	
} //--------------------------------------------


func TestLoadPage(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, LiveCheckGet(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// assert.Equal(t, userJSON, rec.Body.String())
	}
}

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
		// assert.Equal(t, oNoteForm1, rec.Body.String())
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
		assert.Equal(t, http.StatusCreated, rec.Code)
		// assert.Equal(t, oExpectedString, rec.Body.String())
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
