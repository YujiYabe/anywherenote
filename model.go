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
func getData() string {
    // ノート情報の取得

    confdb, err := gorm.Open( useDBMSName , dataDirName + directorySeparator + confDBName )
	// DBログモードon
	confdb.LogMode(true)

    defer confdb.Close()

    if err != nil {
        panic("failed to connect database")
    }

    var count = 0
    confdb.Table("confs").Count(&count)
    
    //設定DBの中にレコードがなければ処理をスキップ

    if count == 0 { 
        // handle the error    
        returnmap := map[string]string{
            "key0": "1",
            "key1": "",
            "key2": "",
        }
        jsonreturnmap, err := json.Marshal(returnmap)
        if err != nil {
            panic("not convert json array")
        }
    
        stringjsonreturnmap := string(jsonreturnmap)
        return stringjsonreturnmap

        // return nil
    }
    
    var conf []Conf
    confdb.Table("confs").Order("updated_at desc").Find(&conf)

    var data2 = ReturnValue{}
    var DataSetList = DataSet{}
    data2.Key0 = "0"

    for _ ,value := range conf {

        NoteDbAddress := value.Address + directorySeparator + noteDBName

        notedb, err := gorm.Open( useDBMSName , NoteDbAddress )
        notedb.LogMode(true)
        defer notedb.Close()
    
        if err != nil {
            panic("failed to connect database")
        }
        NoteList := []Note{}
    
        notedb.Table("notes").Order("updated_at desc").Find(&NoteList)
    
        DataSetList.NoteDbID = value.ID
        DataSetList.NoteDbName = value.Name
        DataSetList.NoteDbAddress = value.Address
        DataSetList.NoteDbUpdateTime = value.UpdatedAt 
        DataSetList.List = NoteList 

        data2.Key1 = append(data2.Key1, DataSetList)
    
    }

    jsonreturnmap, err := json.Marshal(data2)
    if err != nil {
        panic("not convert json array")
    }
    stringjsonreturnmap := string(jsonreturnmap)

    return stringjsonreturnmap
} //--------------------------------------------

// addNote →ノート情報を追加
func addNote( argMap map[string]string ) error {

	//指定したディレクトリにDBファイルを作成
	newNoteName    := argMap["newNoteName"]
	newNoteAddress := argMap["newNoteAddress"]
	
    notefullAddress := newNoteAddress + directorySeparator + noteDBName

    file, err := os.OpenFile( notefullAddress , os.O_WRONLY|os.O_CREATE, 0666 )
    if err != nil {
        //エラー処理
        log.Fatal(err)
    }
    defer file.Close()

    notedb, err := gorm.Open( useDBMSName , notefullAddress )
    defer notedb.Close()

    if err != nil {
        panic("failed to connect database")
    }
    // Migrate the schema
    notedb.AutoMigrate(&Note{})


    //------------------------------
    confdb, err := gorm.Open( useDBMSName , dataDirName + directorySeparator + confDBName )
    defer confdb.Close()
	confdb.LogMode(true)

    if err != nil {
        panic("failed to connect database2")
    }

    var conf Conf

    conf.Name = newNoteName
    conf.Address = newNoteAddress


    // INSERTを実行
    confdb.Create(&conf)


	return nil
} //--------------------------------------------

// addPage →ページ情報を追加
func addPage(argMap map[string]string ) error {

	//指定したディレクトリにDBファイルを作成
	noteAddress    := argMap["noteAddress"]


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
func updateNote( argMap map[string]string  ) error {
//	
	// __________________________________
	// ポスト内容取得
    newNoteName :=argMap[ "newNoteName"]
    postNoteID  :=argMap[ "postNoteID"]

    //------------------------------
    confdb, err := gorm.Open( useDBMSName , dataDirName + directorySeparator + confDBName )
    defer confdb.Close()
	confdb.LogMode(true)

    if err != nil {
        panic("failed to connect database2")
    }

    var conf Conf
    // noteID := int(noteID)
    var noteID int
    noteID, _ = strconv.Atoi(postNoteID)

    confdb.Where("id = ?", noteID ).First(&conf)

    conf.Name = newNoteName
    confdb.Save(&conf)

	return nil
} //--------------------------------------------

// updatePage →ページ情報を更新
func updatePage( argMap map[string]string ) error {
//	 

	// __________________________________
	// ポスト内容取得
    noteAddress := argMap["noteAddress"]
    pageID      := argMap["pageID"]     
    pageTitle   := argMap["pageTitle"]  
    pageBody    := argMap["pageBody"]   

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
func deleteNote( argMap map[string]string ) error {

	// __________________________________
	// ポスト内容取得
    postNoteID  := argMap["postNoteID"]

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
func deletePage( argMap map[string]string) error {
//	argMap map[string]string 

	// __________________________________
	// ポスト内容取得
    noteAddress := argMap["noteAddress"]
    pageID      := argMap["pageID"]     

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

// MakeConfDb は設定データベースの初期化
func MakeConfDB() {


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
