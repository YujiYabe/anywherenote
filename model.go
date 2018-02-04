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
	gorm.Model
	NoteStar    int    `json:"note_star"`
	NoteName    string `json:"note_name"`
	NoteAddress string `json:"note_address"`
} //--------------------------------------------

// Note ページの内容
type Note struct {
	gorm.Model
	PageTitle string `json:"page_title"`
	PageBody  string `json:"page_body"`
	PageStar  int    `json:"page_star"`
} //--------------------------------------------

// DataSet DBファイルの情報とノート情報のセット
type DataSet struct {
	NoteDBID         uint   `json:"NoteDBID"`
	NoteDBStar       int    `json:"NoteDBStar"`
	NoteDBName       string `json:"NoteDBName"`
	NoteDBAddress    string `json:"NoteDBAddress"`
	NoteDBUpdateTime time.Time
	List             []Note `json:"list"`
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

// UserConfig
type UserConfig struct {
	IsEnableAppMode     bool          `json:"IsEnableAppMode"`
	WaitSecondLiveCheck time.Duration `json:"WaitSecondLiveCheck"`
	WaitSecondInterval  time.Duration `json:"WaitSecondInterval"`
	UsePortNumber       string        `json:"UsePortNumber"`
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

func setupDB(dbAddress string) *gorm.DB {
	db, err := gorm.Open(useDBMSName, dbAddress)
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

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
			selectPosition.NoteID = conf[0].ID
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

		notedb.Order("updated_at desc").Find(&NoteList)

		DataSetList.NoteDBID = value.ID
		DataSetList.NoteDBStar = value.NoteStar
		DataSetList.NoteDBName = value.NoteName
		DataSetList.NoteDBAddress = value.NoteAddress
		DataSetList.NoteDBUpdateTime = value.UpdatedAt
		DataSetList.List = NoteList

		data2.DataSet = append(data2.DataSet, DataSetList)

		// 選択中のページIDがなければ最新のページIDを返却
		if selectPosition.PageID == 0 {

			if value.ID == selectPosition.NoteID && len(NoteList) != 0 {

				selectPosition.PageID = NoteList[0].ID

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
	notedb.Where("id = ?", pageID).Find(&note)

	// printEventLog("debug2",note.PageBody)
	// printEventLog("debug2",addFile)
	note.PageBody = note.PageBody + addFile

	notedb.Save(&note)

	return nil

} //--------------------------------------------

// addPage →ページ情報を追加
func addPage(rcvArg map[string]string) error {

	//指定したディレクトリにDBファイルを作成
	noteAddress := rcvArg["noteAddress"]

	noteDBAddress := noteAddress + directorySeparator + noteDBName

	// DBのオープン
	notedb := setupDB(noteDBAddress)
	defer notedb.Close()
	notedb.LogMode(isEnableLogMode)

	notedb.AutoMigrate(&Note{})

	var note Note
	t := time.Now().Format(dateTimeFormat)
	s := t + " created"

	note.PageTitle = s
	note.PageBody = ""
	note.PageStar = 3

	// INSERTを実行
	notedb.Create(&note)

	updateNoteFromPage(noteAddress)

	return nil
} //--------------------------------------------

// updateNote →ノート情報を更新
func updateNote(rcvArg map[string]string) error {

	noteName := rcvArg["noteName"]
	noteStar, _ := strconv.Atoi(rcvArg["noteStar"])
	noteID, _ := strconv.Atoi(rcvArg["postNoteID"])

	// 設定DBのオープン
	confdb := setupDB(confDBAddress)
	defer confdb.Close()
	confdb.LogMode(isEnableLogMode)

	var conf Conf
	confdb.Where("id = ?", noteID).First(&conf)

	conf.NoteName = noteName
	conf.NoteStar = noteStar

	confdb.Save(&conf)

	return nil
} //--------------------------------------------

// updatePage →ページ情報を更新
func updatePage(rcvArg map[string]string) error {

	noteAddress := rcvArg["noteAddress"]
	pageID := rcvArg["pageID"]
	pageTitle := rcvArg["pageTitle"]
	pageBody := rcvArg["pageBody"]

	createdAt, _ := time.Parse(dateTimeFormat, rcvArg["createdAt"])
	updatedAt, _ := time.Parse(dateTimeFormat, rcvArg["updatedAt"])

	// printEventLog("debug", updatedAt)

	if pageTitle == "" {
		pageTitle = " "
	}

	if pageBody == "" {
		pageBody = " "
	}

	noteDBAddress := noteAddress + directorySeparator + noteDBName

	// DBのオープン
	notedb := setupDB(noteDBAddress)
	defer notedb.Close()
	notedb.LogMode(isEnableLogMode)

	// __________________________________
	// DB内容取得
	// NoteList := []Note{}

	// notedb.Model(&NoteList).Where("id = ?", pageID).Update(&Note{PageTitle: pageTitle, PageBody: pageBody, CreatedAt: createdAt, UpdatedAt: updatedAt})
	// notedb.Model(&NoteList).Where("id = ?", pageID).Update(&Note{: , : , : })

	NoteList := Note{
		PageTitle: pageTitle,
		PageBody:  pageBody,
		Model:     gorm.Model{CreatedAt: createdAt, UpdatedAt: updatedAt},
	}
	NoteList.UpdatedAt = updatedAt

	printEventLog("debug", updatedAt)

	notedb.Model(Note{}).Where("id = ?", pageID).Update(NoteList) // こっちで+1日になる

	// model = ModelA{
	// 	ID:   100,
	// 	name: "mary",
	// }
	// model.UpdatedAt = time.Now().AddDate(0, 0, 1) // Now + 1日にする.

	// db.Save(&model) // これはUpdatedAtはtime.Now()になってしまう。
	// db.Model(ModelA{}).Where("id = ?", model.ID).Update(model) // こっちで+1日になる

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

	confdb.Where("id = ?", noteID).Delete(&conf)

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
	notedb.Where("id = ?", pageID).Delete(&Note{})

	return nil
} //--------------------------------------------

// updateNoteFromPage のコメントアウト
func updateNoteFromPage(noteAddress string) {
	printEventLog("start", "updateNoteFromPage 終了")

	// 設定DBのオープン
	confdb := setupDB(confDBAddress)
	defer confdb.Close()
	confdb.LogMode(isEnableLogMode)

	var conf Conf

	confdb.Where("note_address = ?", noteAddress).First(&conf)

	conf.NoteAddress = noteAddress
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
