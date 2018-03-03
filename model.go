package main

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	// "os"
	// "log"
)

// Conf 初期読み込み用設定
type Conf struct {
	NoteID      uint   `json:"NoteID" gorm:"primary_key"`
	NoteStar    int    `json:"NoteStar"`
	NoteName    string `json:"NoteName"`
	NoteAddress string `json:"NoteAddress"`
	NoteUpdate  string `json:"NoteUpdate"`
} //--------------------------------------------

// Note ページの内容
type Note struct {
	PageID     uint   `json:"PageID" gorm:"primary_key"`
	PageStar   int    `json:"PageStar"`
	PageTitle  string `json:"PageTitle"`
	PageBody   string `json:"PageBody"`
	PageUpdate string `json:"PageUpdate"`
} //--------------------------------------------

// DataSet DBファイルの情報とノート情報のセット
type DataSet struct {
	Conf
	List []Note `json:"list"`
}

// SelectPosition 選択中情報
type SelectPosition struct {
	NoteID uint `json:"NoteID"`
	PageID uint `json:"PageID"`
}

// ReturnValue 戻り値とDataSetのセット RtnCodeがリターンコード
type ReturnValue struct {
	RtnCode string         `json:"RtnCode"`
	DataSet []DataSet      `json:"DataSet"`
	SlctPst SelectPosition `json:"SlctPst"`
} //--------------------------------------------

// UserConfig json形式のユーザ設定ファイル
type UserConfig struct {
	IsEnableAppMode     bool          `json:"IsEnableAppMode"`
	WaitSecondLiveCheck time.Duration `json:"WaitSecondLiveCheck"`
	WaitSecondInterval  time.Duration `json:"WaitSecondInterval"`
	UsePortNumber       string        `json:"UsePortNumber"`
}

func setupDB(dbAddress string) *gorm.DB {
	db, err := gorm.Open(useDBMSName, dbAddress)
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func getAllNoteAddress() []Conf {
	confdb := setupDB(confDBAddress)
	defer confdb.Close()
	confdb.LogMode(isEnableLogMode)

	var conf []Conf
	// confdb.Select("id, NoteAddress").Find(&conf)
	confdb.Find(&conf)
	return conf

} //--------------------------------------------

// getData →ノート情報、ページ情報を取得
func getData(selectPosition SelectPosition) string {

	// 設定DBのオープン

	confdb := setupDB(confDBAddress)
	defer confdb.Close()
	confdb.LogMode(isEnableLogMode)

	var count int
	var conf []Conf
	confdb.Find(&conf).Count(&count)

	//設定DBの中にレコードがなければ処理をスキップ
	if count == 0 {
		// handle the error
		returnmap := map[string]string{
			"RtnCode": "1",
			"DataSet": "",
			"SlctPst": "",
		}
		jsonreturnmap, err := json.Marshal(returnmap)
		if err != nil {
			panic("not convert json array")
		}

		stringjsonreturnmap := string(jsonreturnmap)
		return stringjsonreturnmap
	}

	confdb.Order("note_star desc, note_name asc").Find(&conf)
	// confdb.Find(&conf).Order("note_star desc, id desc")
	// "age desc, name"
	// 選択中のノートIDがなければ最新のノートIDを返却
	if selectPosition.NoteID == 0 {
		if len(conf) != 0 {
			selectPosition.NoteID = conf[0].NoteID
		}
	}

	var data2 = ReturnValue{}
	var DataSetList = DataSet{}

	data2.RtnCode = "0"

	for _, value := range conf {

		noteDBAddress := value.NoteAddress + directorySeparator + noteDBName

		// DBのオープン
		notedb := setupDB(noteDBAddress)
		defer notedb.Close()
		notedb.LogMode(isEnableLogMode)

		NoteList := []Note{}

		// notedb.Order("updated_at desc").Find(&NoteList)
		// notedb.Select("page_id, page_star, page_title, page_body, page_update").Order("page_star desc, page_title asc").Find(&NoteList)
		notedb.Order("page_star desc, page_title asc").Find(&NoteList)

		DataSetList.NoteID = value.NoteID
		DataSetList.NoteStar = value.NoteStar
		DataSetList.NoteName = value.NoteName
		DataSetList.NoteAddress = value.NoteAddress
		DataSetList.NoteUpdate = value.NoteUpdate
		DataSetList.List = NoteList

		data2.DataSet = append(data2.DataSet, DataSetList)

		// 選択中のページIDがなければ最新のページIDを返却
		if selectPosition.PageID == 0 {

			if value.NoteID == selectPosition.NoteID && len(NoteList) != 0 {

				selectPosition.PageID = NoteList[0].PageID

			}
		}
	}

	data2.SlctPst = selectPosition

	jsonreturnmap, err := json.Marshal(data2)

	if err != nil {
		panic("not convert json array")
	}

	stringjsonreturnmap := string(jsonreturnmap)

	return stringjsonreturnmap

} //--------------------------------------------

func addFileToPage(rcvArg map[string]string) error {

	noteAddress := rcvArg["noteAddress"]
	pageID := rcvArg["pageID"]
	addFile := rcvArg["addFile"]
	// printEventLog("debug",noteAddress)
	// printEventLog("debug",pageID)

	noteDBAddress := noteAddress + directorySeparator + noteDBName

	// DBのオープン
	notedb := setupDB(noteDBAddress)
	defer notedb.Close()
	notedb.LogMode(isEnableLogMode)
	// __________________________________
	// DB内容取得
	note := Note{}
	notedb.Where("page_id = ?", pageID).Find(&note)

	// printEventLog("debug2",note.PageBody)
	// printEventLog("debug2",addFile)
	note.PageBody = note.PageBody + addFile

	notedb.Save(&note)

	return nil

} //--------------------------------------------
// createPage →ページ情報を追加
func createNote(rcvArg map[string]string) error {

	noteName := rcvArg["noteName"]
	noteAddress := rcvArg["noteAddress"]
	noteUpdate := rcvArg["noteUpdate"]

	//____________________________________
	// 設定DBに追加
	confdb := setupDB(confDBAddress)
	defer confdb.Close()
	confdb.LogMode(isEnableLogMode)

	var conf Conf

	conf.NoteName = noteName
	conf.NoteAddress = noteAddress
	conf.NoteStar = 3
	conf.NoteUpdate = noteUpdate

	// INSERTを実行
	confdb.Create(&conf)

	return nil
} //--------------------------------------------

// createPage →ページ情報を追加
func createPage(rcvArg map[string]string) error {

	noteAddress := rcvArg["noteAddress"]
	pageUpdate := rcvArg["pageUpdate"]

	noteDBAddress := noteAddress + directorySeparator + noteDBName

	// DBのオープン
	notedb := setupDB(noteDBAddress)
	defer notedb.Close()
	notedb.LogMode(isEnableLogMode)

	notedb.AutoMigrate(&Note{})

	var note Note
	// t := time.Now().Format(dateTimeFormat)
	// s := t + " created"
	// note.PageTitle = s

	note.PageTitle = ""
	note.PageBody = ""
	note.PageStar = 3
	note.PageUpdate = pageUpdate

	notedb.Create(&note)

	sndArg := make(map[string]string)

	sndArg["noteAddress"] = noteAddress
	sndArg["noteUpdate"] = pageUpdate

	updateNoteFromPage(sndArg)

	return nil
} //--------------------------------------------

// updateNote →ノート情報を更新
func updateNote(rcvArg map[string]string) error {

	noteName := rcvArg["noteName"]
	noteStar, _ := strconv.Atoi(rcvArg["noteStar"])
	noteID, _ := strconv.Atoi(rcvArg["postNoteID"])
	noteUpdate := rcvArg["noteUpdate"]

	// 設定DBのオープン
	confdb := setupDB(confDBAddress)
	defer confdb.Close()
	confdb.LogMode(isEnableLogMode)
	printEventLog("debug", noteID)

	var conf Conf
	confdb.Where("note_id = ?", noteID).First(&conf)

	conf.NoteName = noteName
	conf.NoteStar = noteStar
	conf.NoteUpdate = noteUpdate

	confdb.Save(&conf)

	return nil
} //--------------------------------------------

// updatePage →ページ情報を更新
func updatePage(rcvArg map[string]string) error {

	noteAddress := rcvArg["noteAddress"]
	pageID := rcvArg["pageID"]
	pageTitle := rcvArg["pageTitle"]
	pageBody := rcvArg["pageBody"]
	pageUpdate := rcvArg["pageUpdate"]
	pageStar, _ := strconv.Atoi(rcvArg["pageStar"])

	// printEventLog("debug", updatedAt)

	if pageTitle == "" {
		pageTitle = " "
	}

	if pageBody == "" {
		pageBody = " "
	}

	NoteList := Note{
		PageTitle:  pageTitle,
		PageBody:   pageBody,
		PageStar:   pageStar,
		PageUpdate: pageUpdate,
	}

	noteDBAddress := noteAddress + directorySeparator + noteDBName

	// DBのオープン
	notedb := setupDB(noteDBAddress)
	defer notedb.Close()
	notedb.LogMode(isEnableLogMode)

	notedb.Model(Note{}).Where("page_id = ?", pageID).Update(NoteList)

	sndArg := make(map[string]string)

	sndArg["noteAddress"] = noteAddress
	sndArg["noteUpdate"] = pageUpdate

	updateNoteFromPage(sndArg)

	return nil
} //--------------------------------------------

// deleteNote →ノート情報を追加
func deleteNote(rcvArg map[string]string) error {

	// __________________________________
	// ポスト内容取得
	noteID, _ := strconv.Atoi(rcvArg["postNoteID"])

	// 設定DBのオープン
	confdb := setupDB(confDBAddress)
	defer confdb.Close()
	confdb.LogMode(isEnableLogMode)

	var conf Conf

	confdb.Where("note_id = ?", noteID).Delete(&conf)

	return nil
} //--------------------------------------------

// deletePage →ページ情報を追加
func deletePage(rcvArg map[string]string) error {

	// __________________________________
	// ポスト内容取得
	noteAddress := rcvArg["noteAddress"]
	pageID := rcvArg["pageID"]

	noteDBAddress := noteAddress + directorySeparator + noteDBName
	// DBのオープン
	notedb := setupDB(noteDBAddress)
	defer notedb.Close()
	notedb.LogMode(isEnableLogMode)

	// __________________________________
	// DB内容取得
	notedb.Where("page_id = ?", pageID).Delete(&Note{})

	return nil
} //--------------------------------------------

// updateNoteFromPage のコメントアウト
func updateNoteFromPage(rcvArg map[string]string) {
	printEventLog("start", "updateNoteFromPage 終了")

	noteAddress := rcvArg["noteAddress"]
	noteUpdate := rcvArg["noteUpdate"]

	// 設定DBのオープン
	confdb := setupDB(confDBAddress)
	defer confdb.Close()
	confdb.LogMode(isEnableLogMode)

	var conf Conf

	confdb.Where("note_address = ?", noteAddress).First(&conf)

	conf.NoteAddress = noteAddress
	conf.NoteUpdate = noteUpdate
	confdb.Save(&conf)
	printEventLog("end", "updateNoteFromPage 終了")
} //--------------------------------------------

//--------------------------------------------
func dbApplyType(targetDBAddress string, targetStruct interface{}) {
	db := setupDB(targetDBAddress)
	defer db.Close()
	db.LogMode(isEnableLogMode)

	// Migrate the schema
	db.AutoMigrate(targetStruct)

} //--------------------------------------------
