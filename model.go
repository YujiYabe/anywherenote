package main
import (

    "encoding/json"
    "github.com/jinzhu/gorm"
    _ "github.com/mattn/go-sqlite3"
    "strconv"

    // "os"
    // "log"

)
// Conf 初期読み込み用設定
type Conf struct {
    gorm.Model
    Name    string  `json:"conf_name"`
    Address string  `json:"conf_address"`
} //--------------------------------------------


// Note ページの内容
type Note struct {
    gorm.Model
    PageTitle   string `json:"page_title"`
    PageBody    string `json:"page_body"`
} //--------------------------------------------






func setupDB( dbAddress string ) *gorm.DB {
	db, err := gorm.Open( useDBMSName , dbAddress )
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

// getData →ノート情報、ページ情報を取得
func getData( selectPosition SelectPosition ) string {

    // 設定DBのオープン

    confdb := setupDB( confDBAddress )
	defer confdb.Close()
    confdb.LogMode(true)


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

    
    confdb.Find( &conf ).Order( "id desc" )

    // 選択中のノートIDがなければ最新のノートIDを返却
    if selectPosition.NoteID == 0 {
        if len(conf) != 0 {
            selectPosition.NoteID = conf[0].ID
        }
    }
    
    var data2 = ReturnValue{}
    var DataSetList = DataSet{}

    data2.RtnCode = "0"

    for _ ,value := range conf {

        noteDBAddress := value.Address + directorySeparator + noteDBName

        // DBのオープン
        notedb := setupDB( noteDBAddress )
        defer notedb.Close()
        notedb.LogMode(true)


        NoteList := []Note{}

        notedb.Order("updated_at desc").Find( &NoteList )

        DataSetList.NoteDBID         = value.ID
        DataSetList.NoteDBName       = value.Name
        DataSetList.NoteDBAddress    = value.Address
        DataSetList.NoteDBUpdateTime = value.UpdatedAt 
        DataSetList.List             = NoteList 

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




// addPage →ページ情報を追加
func addPage(rcvArg map[string]string ) error {

    //指定したディレクトリにDBファイルを作成
    noteAddress    := rcvArg["noteAddress"]

    noteDBAddress := noteAddress + directorySeparator +  noteDBName

    // DBのオープン
    notedb := setupDB( noteDBAddress )
    defer notedb.Close()
    notedb.LogMode(true)

    notedb.AutoMigrate(&Note{})

    var note Note

    note.PageTitle = ""
    note.PageBody  = ""

    // INSERTを実行
    notedb.Create(&note)

    updateNoteFromPage(noteAddress) 


    return nil
} //--------------------------------------------

// updateNote →ノート情報を更新
func updateNote( rcvArg map[string]string  ) error {

    noteName    :=rcvArg[ "noteName" ]
    noteID, _   :=strconv.Atoi( rcvArg[ "postNoteID" ] )   

    // 設定DBのオープン
    confdb := setupDB( confDBAddress )
	defer confdb.Close()
    confdb.LogMode(true)

    var conf Conf
    confdb.Where("id = ?", noteID ).First( &conf )

    conf.Name = noteName
    
    confdb.Save(&conf)

    return nil
} //--------------------------------------------

// updatePage →ページ情報を更新
func updatePage( rcvArg map[string]string ) error {

    noteAddress := rcvArg["noteAddress"]
    pageID      := rcvArg["pageID"]     
    pageTitle   := rcvArg["pageTitle"]  
    pageBody    := rcvArg["pageBody"]   

    if pageTitle == "" {
        pageTitle =  " "
    }

    if pageBody == "" {
        pageBody =  " "
    }
    
    noteDBAddress := noteAddress + directorySeparator + noteDBName

    // DBのオープン
    notedb := setupDB( noteDBAddress )
    defer notedb.Close()
    notedb.LogMode(true)

    // __________________________________
    // DB内容取得
    NoteList := []Note{}

    notedb.Model( &NoteList ).Where( "id = ?", pageID ).Update( &Note{ PageTitle: pageTitle, PageBody: pageBody })

    return nil
} //--------------------------------------------

// deleteNote →ノート情報を追加
func deleteNote( rcvArg map[string]string ) error {

    // __________________________________
    // ポスト内容取得
    noteID, _   := strconv.Atoi( rcvArg["postNoteID"] ) 

    // 設定DBのオープン
    confdb := setupDB( confDBAddress )
	defer confdb.Close()
    confdb.LogMode(true)

    var conf Conf

    confdb.Where( "id = ?" , noteID ).Delete( &conf )

    return nil
} //--------------------------------------------

// deletePage →ページ情報を追加
func deletePage( rcvArg map[string]string ) error {

    // __________________________________
    // ポスト内容取得
    noteAddress := rcvArg["noteAddress"]
    pageID      := rcvArg["pageID"]     

    noteDBAddress := noteAddress + directorySeparator + noteDBName
    // DBのオープン
    notedb := setupDB( noteDBAddress )
    defer notedb.Close()
    notedb.LogMode(true)

    // __________________________________
    // DB内容取得
    notedb.Where("id = ?", pageID).Delete(&Note{})

    return nil
} //--------------------------------------------

// updateNoteFromPage のコメントアウト
func updateNoteFromPage( noteAddress string )  {

    // 設定DBのオープン
    confdb := setupDB( confDBAddress )
	defer confdb.Close()
    confdb.LogMode(true)

    var conf Conf

    confdb.Where("address = ?", noteAddress ).First(&conf)

    conf.Address = noteAddress
    confdb.Save(&conf)

} //--------------------------------------------


//--------------------------------------------
func dbApplyType( targetDBAddress string, targetStruct interface{} ) {
    db := setupDB( targetDBAddress )
	defer db.Close()
    db.LogMode(true)

    // Migrate the schema
    db.AutoMigrate( targetStruct )
    
} //--------------------------------------------
