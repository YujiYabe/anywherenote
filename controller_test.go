package main

import (
	"io/ioutil"
	"runtime"
	"testing"

	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

const (
	o = true
	x = false
)

var (
	isEnableTestLoadPage   = o
	isEnableTestLiveCheck  = o
	isEnableTestCreateNote = o
	isEnableTestUpdateNote = o
	isEnableTestCreatePage = o
	isEnableTestUpdatePage = o
	isEnableTestDeletePage = o
	isEnableTestDeleteNote = o

	noteID1      = "1"
	noteAddress1 = "C:\\workspace\\gotest"

	noteName1 = "create_note"
	noteName2 = "change_note"

	cmmnUpdate1 = "2018/01/01 01:01:01"
	cmmnUpdate2 = "2018/02/02 01:01:01"
	cmmnStar1   = "3"
	cmmnStar2   = "2"

	pageID1 = "1"

	pageTitle1 = ""
	pageBody1  = ""

	pageTitle2 = "second_page_title"
	pageBody2  = "second_page_body"

	oNoteForm1 = map[string]string{"note_id": noteID1, "note_address": noteAddress1, "note_star": cmmnStar1, "note_name": noteName1, "note_update": cmmnUpdate1}
	oNoteForm2 = map[string]string{"note_id": noteID1, "note_address": noteAddress1, "note_star": cmmnStar2, "note_name": noteName2, "note_update": cmmnUpdate2}

	xNoteForm1 = map[string]string{"note_id": "1", "note_address": "C:\\workspace\\nothingpath", "note_name": "create_note", "note_update": "2018/01/01 01:01:01"}

	oPageForm1 = map[string]string{"note_address": noteAddress1, "page_id": pageID1, "page_star": cmmnStar1, "page_title": pageTitle1, "page_body": pageBody1, "page_update": cmmnUpdate1}
	oPageForm2 = map[string]string{"note_address": noteAddress1, "page_id": pageID1, "page_star": cmmnStar2, "page_title": pageTitle2, "page_body": pageBody2, "page_update": cmmnUpdate1}

	oExpectedString1 = "\"{\\\"RtnCode\\\":\\\"0\\\",\\\"DataSet\\\":[{\\\"NoteID\\\":1,\\\"NoteStar\\\":2,\\\"NoteName\\\":\\\"change_note\\\",\\\"NoteAddress\\\":\\\"C:\\\\\\\\workspace\\\\\\\\gotest\\\",\\\"NoteUpdate\\\":\\\"\\\",\\\"list\\\":[{\\\"PageID\\\":1,\\\"PageStar\\\":3,\\\"PageTitle\\\":\\\"\\\",\\\"PageBody\\\":\\\"\\\",\\\"PageUpdate\\\":\\\"\\\"}]}],\\\"SlctPst\\\":{\\\"NoteID\\\":1,\\\"PageID\\\":1}}\""
	oExpectedString2 = "\"{\\\"RtnCode\\\":\\\"0\\\",\\\"DataSet\\\":[{\\\"NoteID\\\":1,\\\"NoteStar\\\":2,\\\"NoteName\\\":\\\"change_note\\\",\\\"NoteAddress\\\":\\\"C:\\\\\\\\workspace\\\\\\\\gotest\\\",\\\"NoteUpdate\\\":\\\"\\\",\\\"list\\\":[{\\\"PageID\\\":1,\\\"PageStar\\\":3,\\\"PageTitle\\\":\\\"second_page_title\\\",\\\"PageBody\\\":\\\"second_page_body\\\",\\\"PageUpdate\\\":\\\"\\\"}]}],\\\"SlctPst\\\":{\\\"NoteID\\\":1,\\\"PageID\\\":1}}\""
)

func init() {

	// ディレクトリセパレータ
	if runtime.GOOS == "windows" {
		directorySeparator = "\\"
	}

	// 設定ファイルの読み込み
	userConfigFile, err := ioutil.ReadFile(dataDirName + directorySeparator + "public" + directorySeparator + userConfFile)
	if err != nil {
		panic(err)
	}

	json.Unmarshal(userConfigFile, &userConfig)

	//テスト実行前にファイルを削除
	testDBAddress := oNoteForm1["note_address"] + directorySeparator + noteDBName

	// if err1 == nil{
	if err1 := osCheckFile(testDBAddress); err1 == nil {
		if err2 := os.Remove(testDBAddress); err2 != nil {
			log.Println(err2)
		}
	}

	confDBAddress := dataDirName + directorySeparator + confDBName

	// if err1 == nil{
	if err1 := osCheckFile(confDBAddress); err1 == nil {
		if err2 := os.Remove(confDBAddress); err2 != nil {
			log.Println(err2)
		}
	}

	for {
		time.Sleep(1 * time.Millisecond)
		// log.Println( "waiting" )

		if err1 := osCheckFile(confDBAddress); err1 != nil {
			break
		}

	}

	loadTemplates()
	checkConfig()

} //--------------------------------------------

func TestLoadPage(t *testing.T) {
	if isEnableTestLoadPage {
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
}

func TestLiveCheck(t *testing.T) {
	if isEnableTestLiveCheck {
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
}

func TestCreateNote(t *testing.T) {
	if isEnableTestCreateNote {

		// Setup
		e := echo.New()

		f := make(url.Values)
		f.Set("note_name", oNoteForm1["note_name"])
		f.Set("note_address", oNoteForm1["note_address"])
		f.Set("note_update", oNoteForm1["note_update"])

		req := httptest.NewRequest(echo.POST, "/createnote", strings.NewReader(f.Encode()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Assertions
		if assert.NoError(t, CreateNotePost(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			// assert.Equal(t, oNoteForm1, rec.Body.String())
		}
	}
}
func TestUpdateNote(t *testing.T) {
	if isEnableTestUpdateNote {
		// Setup
		e := echo.New()

		f := make(url.Values)

		f.Set("note_id", oNoteForm2["note_id"])
		f.Set("note_name", oNoteForm2["note_name"])
		f.Set("note_star", oNoteForm2["note_star"])

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
}

func TestCreatePage(t *testing.T) {
	if isEnableTestCreatePage {

		// Setup
		e := echo.New()

		f := make(url.Values)

		f.Set("note_id", oNoteForm1["note_id"])
		f.Set("note_address", oNoteForm1["note_address"])

		req := httptest.NewRequest(echo.POST, "/CreatePagePost", strings.NewReader(f.Encode()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Assertions
		if assert.NoError(t, CreatePagePost(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, oExpectedString1, rec.Body.String())
		}
	}
}

func TestUpdatePage(t *testing.T) {
	if isEnableTestUpdatePage {

		// Setup
		e := echo.New()

		f := make(url.Values)
		f.Set("note_address", oPageForm2["note_address"])
		f.Set("page_id", oPageForm2["page_id"])
		f.Set("page_title", oPageForm2["page_title"])
		f.Set("page_body", oPageForm2["page_body"])
		f.Set("created_at", oPageForm2["created_at"])
		f.Set("updated_at", oPageForm2["updated_at"])

		req := httptest.NewRequest(echo.POST, "/UpdatePagePost", strings.NewReader(f.Encode()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Assertions
		if assert.NoError(t, UpdatePagePost(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, oExpectedString2, rec.Body.String())
		}
	}
}

func TestDeletePage(t *testing.T) {
	if isEnableTestDeletePage {

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
}

func TestDeleteNote(t *testing.T) {
	if isEnableTestDeleteNote {

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
}
