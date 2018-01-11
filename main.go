package main

import (
    "html/template"
    "io"
    "net/http"
    "os"
    "runtime"
    "log"
	"encoding/json"
    "time"
    "strconv"
    "github.com/labstack/echo"
    "github.com/labstack/echo/middleware"
    "github.com/jinzhu/gorm"
    _ "github.com/mattn/go-sqlite3"
    "github.com/skratchdot/open-golang/open"

)
// グローバル変数
var (
    // チャネル
    // liveChannel = make(chan string)

    recieveString string
    directorySeparator  = "/" // linux separator
)

// 定数
const isEnableAppMode = true // for debug
const waitSecondLiveCheck = 8
const waitSecondInterval = 1 
const dateTimeFormat = "2006-01-02 15:04:05"

const settingDBName = "setting.db"
const noteDBName = "note.db"

const usePortNumber = "3000"
const useDBMSName = "sqlite3"

// Setting 初期読み込み用設定
type Setting struct {
    gorm.Model
    Name    string  `json:"setting_name"`
    Address string  `json:"setting_address"`
} //--------------------------------------------



// Note オンラインストレージのノート
type Note struct {
    gorm.Model
    PageTitle   string `json:"page_title"`
    PageBody    string `json:"page_body"`
} //--------------------------------------------


// DataSet DBファイルの情報とノート情報のセット
type DataSet struct {
	NoteDbID         uint   `json:"NoteDbID"`
	NoteDbName       string `json:"NoteDbName"`
	NoteDbAddress    string `json:"NoteDbAddress"`
	NoteDbUpdateTime time.Time
    List             []Note `json:"list"`
}

// ReturnValue 戻り値とDataSetのセット Key0がリターンコード 
type ReturnValue struct {
    Key0   string    `json:"key0"`
	Key1   []DataSet `json:"key1"`
} //--------------------------------------------



// レイアウト適用済のテンプレートを保存するmap
var templates map[string]*template.Template

// Template はHTMLテンプレートを利用するためのRenderer Interface
type Template struct {}

// Render はHTMLテンプレートにデータを埋め込んだ結果をWriterに書き込み
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return templates[name].ExecuteTemplate(w, "layout.html", data)
} //--------------------------------------------

// 初期化
func init() {
    if runtime.GOOS == "windows" {
        directorySeparator = "\\"
    }
    
    loadTemplates()
    checkConfig()

} //--------------------------------------------

func main() {
    // Echoのインスタンスを生成
    e := echo.New()

    // テンプレートを利用するためのRendererの設定
    t := &Template{}
    e.Renderer = t

    // ミドルウェアを設定
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    // 静的ファイルのパスを設定
    e.Static("/public", "public")
    e.File("/favicon.ico", "public/favicon.ico")

    // 各ルーティングに対するハンドラを設定

    e.POST("/addnote", HandleAddNotePost)
    e.POST("/addpage", HandleAddPagePost)

    e.POST("/updatenote", HandleUpdateNotePost)
    e.POST("/updatepage", HandleUpdatePagePost)

    e.POST("/deletepage", HandleDeletePagePost)
    e.POST("/deletenote", HandleDeleteNotePost)

    e.GET("/livecheck", HandleLiveCheckGet)
    e.GET("/", HandleLoadPageGet)
 

    open.Run("http://localhost:" + usePortNumber)

	go calcTime()


    // サーバーを開始
    e.Logger.Fatal(e.Start(":" + usePortNumber))
} //--------------------------------------------

// HandleLiveCheckGet のコメントアウト
func HandleLiveCheckGet(c echo.Context) error {
    // //現在時刻取得
    t := time.Now()

    //現在からn秒後の時刻を取得
    afterTime := t.Add(time.Duration( waitSecondLiveCheck ) * time.Second).Format(dateTimeFormat)

    recieveString = afterTime


    //チャネルへ時刻情報を送信
    // liveChannel <- c.FormValue("expireLiveTime")



    return nil
} //--------------------------------------------

// HandleDeletePagePost のコメントアウト
func HandleDeletePagePost(c echo.Context) error {
    log.Println("==============================================================================" )
    log.Println("==============================================================================" )
    log.Println("ページ削除 開始")

	// __________________________________
	// ポスト内容取得
    noteAddress := c.FormValue("note_address")
    pageID      := c.FormValue("page_id")

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


    log.Println("ページ削除 終了")
    log.Println("==============================================================================" )
    log.Println("==============================================================================" )

    return nil

} //--------------------------------------------


// HandleUpdateNoteFromPage のコメントアウト
func HandleUpdateNoteFromPage(noteAddress string)  {
    log.Println("==============================================================================" )
    log.Println("==============================================================================" )
    log.Println("ノート更新日時更新 開始")


    settingdb, err := gorm.Open("sqlite3", settingDBName )
    defer settingdb.Close()
	settingdb.LogMode(true)

    if err != nil {
        panic("failed to connect database")
    }

    var setting Setting

    settingdb.Where("address = ?", noteAddress ).First(&setting)

    setting.Address = noteAddress
    settingdb.Save(&setting)


    log.Println("ノート更新日時更新 終了")
    log.Println("==============================================================================" )
    log.Println("==============================================================================" )

} //--------------------------------------------

//HandleUpdatePagePost は /hello のPost時のHTMLデータ生成処理を行います。
func HandleUpdatePagePost(c echo.Context) error {
    log.Println("==============================================================================" )
    log.Println("==============================================================================" )
    log.Println("ページ更新 開始")

	// __________________________________
	// ポスト内容取得
    noteAddress := c.FormValue("note_address")
    pageID      := c.FormValue("page_id")
    pageTitle   := c.FormValue("page_title")
    pageBody    := c.FormValue("page_body")

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

    HandleUpdateNoteFromPage(noteAddress)

    log.Println("ページ更新 終了")
    log.Println("==============================================================================" )
    log.Println("==============================================================================" )

    return nil
} //--------------------------------------------


//HandleUpdateNotePost は /hello のPost時のHTMLデータ生成処理を行います。
func HandleUpdateNotePost(c echo.Context) error {
    log.Println("==============================================================================" )
    log.Println("==============================================================================" )
    log.Println("ノート更新 開始")

	// __________________________________
	// ポスト内容取得
    newNoteName := c.FormValue("new_note_name")
    postNoteID  := c.FormValue("note_id")

    
    //------------------------------
    settingdb, err := gorm.Open( useDBMSName , settingDBName )
    defer settingdb.Close()
	settingdb.LogMode(true)

    if err != nil {
        panic("failed to connect database2")
    }

    var setting Setting
    // noteID := int(noteID)
    var noteID int
    noteID, _ = strconv.Atoi(postNoteID)

    settingdb.Where("id = ?", noteID ).First(&setting)

    setting.Name = newNoteName
    settingdb.Save(&setting)



    log.Println("ノート更新 終了")
    log.Println("==============================================================================" )
    log.Println("==============================================================================" )

    return nil
} //--------------------------------------------



//HandleDeleteNotePost は /hello のPost時のHTMLデータ生成処理を行います。
func HandleDeleteNotePost(c echo.Context) error {
    log.Println("==============================================================================" )
    log.Println("==============================================================================" )
    log.Println("ノート削除 開始")

	// __________________________________
	// ポスト内容取得
    postNoteID  := c.FormValue("note_id")

    //------------------------------
    settingdb, err := gorm.Open( useDBMSName , settingDBName )
    defer settingdb.Close()
	settingdb.LogMode(true)

    if err != nil {
        panic("failed to connect database")
    }

    var setting Setting
    var noteID int
    noteID, _ = strconv.Atoi(postNoteID)

    settingdb.Where("id = ?", noteID ).Delete(&setting)

    log.Println("ノート削除 終了")
    log.Println("==============================================================================" )
    log.Println("==============================================================================" )

    return nil
} //--------------------------------------------




//HandleAddNotePost は /hello のPost時のHTMLデータ生成処理を行います。
func HandleAddNotePost(c echo.Context) error {
    log.Println("==============================================================================" )
    log.Println("==============================================================================" )
    log.Println("ノート追加 開始")

    //指定したディレクトリにDBファイルを作成
    newNoteName    := c.FormValue("new_note_name")
    newNoteAddress := c.FormValue("new_note_address")

    // タスク パスが存在するかの確認


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
    settingdb, err := gorm.Open( useDBMSName , settingDBName )
    defer settingdb.Close()
	settingdb.LogMode(true)

    if err != nil {
        panic("failed to connect database2")
    }

    var setting Setting

    setting.Name = newNoteName
    setting.Address = newNoteAddress


    // INSERTを実行
    settingdb.Create(&setting)


    log.Println("ノート追加 終了")
    log.Println("==============================================================================" )
    log.Println("==============================================================================" )
    return nil
    
} //--------------------------------------------





//HandleAddPagePost は /hello のPost時のHTMLデータ生成処理を行います。
func HandleAddPagePost(c echo.Context) error {
    log.Println("==============================================================================" )
    log.Println("==============================================================================" )
    log.Println("ページ追加 開始")

    //指定したディレクトリにDBファイルを作成
    noteAddress   := c.FormValue("note_address")

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

    HandleUpdateNoteFromPage(noteAddress) 

    log.Println("ページ追加 終了")
    log.Println("==============================================================================" )
    log.Println("==============================================================================" )

    return nil

} //--------------------------------------------





func checkConfig() {
    //=============================================
    // 初期設定
    // 設定DB読み込み
    _, err := os.Stat(settingDBName)
    if err != nil {
        MakeSettingDb()
    }

    // 共有ファイルの場所情報を取得
	db, err := gorm.Open( useDBMSName , settingDBName)
	if err != nil {
	  panic("failed to connect database")
	}
	defer db.Close()
	var setting Setting
	db.First(&setting) // find product with id 1


} //--------------------------------------------




// ===============================================

// MakeSettingDb は設定データベースの初期化
func MakeSettingDb() {
    log.Println("==============================================================================" )
    log.Println("==============================================================================" )
    log.Println("設定DB作成 開始")

    file, err := os.OpenFile( settingDBName , os.O_WRONLY|os.O_CREATE, 0666 )
    if err != nil {
        //エラー処理
        log.Fatal(err)
    }
    defer file.Close()
    db, err := gorm.Open( useDBMSName , settingDBName )
    if err != nil {
        panic("failed to connect database")
    }
    // Migrate the schema
    db.AutoMigrate(&Setting{})

    log.Println("設定DB作成 終了")
    log.Println("==============================================================================" )
    log.Println("==============================================================================" )
}

// 各HTMLテンプレートに共通レイアウトを適用した結果を保存します（初期化時に実行）。
func loadTemplates() {
    var baseTemplate = "templates/layout.html"
    templates = make(map[string]*template.Template)
    templates["loadpage"]       = template.Must(template.ParseFiles(baseTemplate, "templates/loadpage.html"))
} //--------------------------------------------


// HandleLoadPageGet は /hello_form のGet時のHTMLデータ生成処理を行います。
func HandleLoadPageGet(c echo.Context) error {
    // return nil

    // ノート情報の取得
    settingdb, err := gorm.Open( useDBMSName , settingDBName )
	// DBログモードon
	settingdb.LogMode(true)

    defer settingdb.Close()

    if err != nil {
        panic("failed to connect database")
    }

    var count = 0
    settingdb.Table("settings").Count(&count)
    
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
        return c.Render(http.StatusOK, "loadpage", stringjsonreturnmap)

        // return nil
    }
    
    var setting []Setting
    settingdb.Table("settings").Order("updated_at desc").Find(&setting)

    var data2 = ReturnValue{}
    var DataSetList = DataSet{}
    data2.Key0 = "0"

    for _ ,value := range setting {

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

    return c.Render(http.StatusOK, "loadpage", stringjsonreturnmap)

} //--------------------------------------------

// ==============================================
func endProcess() {
    os.Exit(0)
}// ==============================================


// ==============================================
func calcTime() {


    //ブラウザが起動して、pingを発行するまでn秒まつ
    time.Sleep( waitSecondLiveCheck * time.Second)

    for {
        time.Sleep( waitSecondInterval * time.Second)

        t := time.Now()
        beforeTime := t.Add(time.Duration(1) * time.Second).Format(dateTimeFormat)
        now, _ := time.Parse(dateTimeFormat ,beforeTime)


        old, _ := time.Parse(dateTimeFormat , recieveString)

        if !old.After(now) {    // old <= now --- ! old > now
            if isEnableAppMode {
                log.Println("アプリ終了")
                endProcess()
            }
        }
    }
}// ==============================================



