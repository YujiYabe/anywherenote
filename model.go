package main
import (

    "encoding/json"
    "github.com/jinzhu/gorm"
    _ "github.com/mattn/go-sqlite3"
    "os"
    "log"
    "strconv"

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

// getData →ノート情報、ページ情報を取得
func getData( selectPosition SelectPosition ) string {
    // ノート情報の取得

    // 設定DBのオープン
    confdb, err := gorm.Open( useDBMSName , dataDirName + directorySeparator + confDBName )

    confdb.LogMode(true)
    defer confdb.Close()

    if err != nil {
        panic("failed to connect database")
    }

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

    
    confdb.Find(&conf).Order("id desc")

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

        NoteDBAddress := value.Address + directorySeparator + noteDBName

        notedb, err := gorm.Open( useDBMSName , NoteDBAddress )
        notedb.LogMode(true)
        defer notedb.Close()
    
        if err != nil {
            panic("failed to connect database")
        }
        NoteList := []Note{}

        notedb.Find(&NoteList).Order("updated_at desc")


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

// addNote →ノート情報を追加
func addNote( rcvMap map[string]string ) error {

    //指定したディレクトリにDBファイルを作成
    newNoteName    := rcvMap["newNoteName"]
    newNoteAddress := rcvMap["newNoteAddress"]
    
    //フォルダが存在確認
    _, errAddress := os.Stat(newNoteAddress)
    if errAddress != nil {
        return errAddress
    }

    notefullAddress := newNoteAddress + directorySeparator + noteDBName
    _, errFullAddress := os.Stat(notefullAddress)

    //新規ノートDB作成
    if errFullAddress != nil {

    
        notedb, err := gorm.Open( useDBMSName , notefullAddress )
        defer notedb.Close()
    
        if err != nil {
            panic("failed to connect database")
        }
        // Migrate the schema
        notedb.AutoMigrate(&Note{})


        // 新規DBの場合、最初のページ作成
        // argMap := make(map[string]string)
        // argMap["noteAddress"] = newNoteAddress
    
        // addPage( argMap )
    
    }

    //------------------------------
    confdb, err := gorm.Open( useDBMSName , dataDirName + directorySeparator + confDBName )
    defer confdb.Close()
    confdb.LogMode(true)

    if err != nil {
        panic("failed to connect database")
    }

    var conf Conf

    conf.Name    = newNoteName
    conf.Address = newNoteAddress

    // INSERTを実行
    confdb.Create(&conf)

    return nil
} //--------------------------------------------

// addPage →ページ情報を追加
func addPage(rcvMap map[string]string ) error {

    //指定したディレクトリにDBファイルを作成
    noteAddress    := rcvMap["noteAddress"]


    notedb, err := gorm.Open( useDBMSName , noteAddress + directorySeparator +  noteDBName )
    if err != nil {
        panic("failed to connect database")
    }
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
func updateNote( rcvMap map[string]string  ) error {

    newNoteName :=rcvMap[ "newNoteName"]
    postNoteID  :=rcvMap[ "postNoteID"]

    //------------------------------
    confdb, err := gorm.Open( useDBMSName , dataDirName + directorySeparator + confDBName )
    defer confdb.Close()
    confdb.LogMode(true)

    if err != nil {
        panic("failed to connect database2")
    }

    var conf Conf

    var noteID int
    noteID, _ = strconv.Atoi(postNoteID)

    confdb.Where("id = ?", noteID ).First(&conf)

    conf.Name = newNoteName
    confdb.Save(&conf)

    return nil
} //--------------------------------------------

// updatePage →ページ情報を更新
func updatePage( rcvMap map[string]string ) error {

    noteAddress := rcvMap["noteAddress"]
    pageID      := rcvMap["pageID"]     
    pageTitle   := rcvMap["pageTitle"]  
    pageBody    := rcvMap["pageBody"]   

    if pageTitle == "" {
        pageTitle =  " "
    }

    if pageBody == "" {
        pageBody =  " "
    }
    
    notedb, err := gorm.Open("sqlite3", noteAddress + directorySeparator + noteDBName )
    if err != nil {
        panic("failed to connect database")
    }

    defer notedb.Close()
    notedb.LogMode(true)


    // __________________________________
    // DB内容取得
    NoteList := []Note{}

    notedb.Model(&NoteList).Where("id = ?", pageID).Update(&Note{PageTitle: pageTitle, PageBody: pageBody})


    if err != nil {
        panic("failed to connect database")
    }

    updateNoteFromPage(noteAddress)

    return nil
} //--------------------------------------------

// deleteNote →ノート情報を追加
func deleteNote( rcvMap map[string]string ) error {

    // __________________________________
    // ポスト内容取得
    postNoteID  := rcvMap["postNoteID"]

    //------------------------------
    confdb, err := gorm.Open( useDBMSName , dataDirName + directorySeparator + confDBName )
    defer confdb.Close()
    confdb.LogMode(true)

    if err != nil {
        panic("failed to connect database")
    }

    var conf Conf
    var noteID int
    noteID, _ = strconv.Atoi(postNoteID)

    confdb.Where("id = ?", noteID ).Delete(&conf)

    return nil
} //--------------------------------------------

// deletePage →ページ情報を追加
func deletePage( rcvMap map[string]string) error {
//    rcvMap map[string]string 

    // __________________________________
    // ポスト内容取得
    noteAddress := rcvMap["noteAddress"]
    pageID      := rcvMap["pageID"]     

    notedb, err := gorm.Open("sqlite3", noteAddress + directorySeparator + noteDBName )
    if err != nil {
        panic("failed to connect database")
    }

    defer notedb.Close()
    notedb.LogMode(true)

    // __________________________________
    // DB内容取得
    notedb.Where("id = ?", pageID).Delete(&Note{})
    
    if err != nil {
        panic("failed to connect database")
    }

    return nil
} //--------------------------------------------

// updateNoteFromPage のコメントアウト
func updateNoteFromPage(noteAddress string)  {

    confdb, err := gorm.Open( useDBMSName , dataDirName + directorySeparator + confDBName )
    defer confdb.Close()
    confdb.LogMode(true)

    if err != nil {
        panic("failed to connect database")
    }

    var conf Conf

    confdb.Where("address = ?", noteAddress ).First(&conf)

    conf.Address = noteAddress
    confdb.Save(&conf)



} //--------------------------------------------


// ===============================================

// makeConfDB は設定データベースの初期化
func makeConfDB() {

    file, err := os.OpenFile( dataDirName + directorySeparator + confDBName , os.O_WRONLY|os.O_CREATE, 0666 )
    if err != nil {
        //エラー処理
        log.Fatal(err)
    }
    defer file.Close()

    db, err := gorm.Open( useDBMSName , dataDirName + directorySeparator + confDBName  )
    if err != nil {
        panic("failed to connect database")
    }
    // Migrate the schema
    db.AutoMigrate(&Conf{})

}
