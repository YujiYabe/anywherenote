package main

import (
    "html/template"
    "io"
    "net/http"
    "os"
    "runtime"
    "log"
    "time"
    "github.com/labstack/echo"
    "github.com/labstack/echo/middleware"
    // "github.com/jinzhu/gorm"
    // _ "github.com/mattn/go-sqlite3"

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


const useDBMSName = "sqlite3"
const confDBName  = "conf.db"
const noteDBName  = "note.db"


const usePortNumber = "3000"



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

    e.GET("/", HandleLoadPageGet)
    e.POST("/addnote", HandleAddNotePost)
    e.POST("/addpage", HandleAddPagePost)

    e.POST("/updatenote", HandleUpdateNotePost)
    e.POST("/updatepage", HandleUpdatePagePost)

    e.POST("/deletenote", HandleDeleteNotePost)
    e.POST("/deletepage", HandleDeletePagePost)

    e.GET("/livecheck", HandleLiveCheckGet)
 

    open.Run("http://localhost:" + usePortNumber)

	go calcTime()


    // サーバーを開始
    e.Logger.Fatal(e.Start(":" + usePortNumber))
} //--------------------------------------------


// HandleLoadPageGet のコメントアウト
func HandleLoadPageGet(c echo.Context) error {
    printEvent("start" , "データ取得 開始")

    returnjson := getData()

    printEvent("end" , "データ取得 終了")
    return c.Render(http.StatusOK, "loadpage", returnjson)

} //--------------------------------------------

//HandleAddNotePost は /hello のPost時のHTMLデータ生成処理を行います。
func HandleAddNotePost(c echo.Context) error {

    printEvent("start" , "ノート追加 開始")

    argMap := make(map[string]string)
    argMap["newNoteName"]    = c.FormValue("new_note_name")
    argMap["newNoteAddress"] = c.FormValue("new_note_address")

    addNote(argMap)

    // タスク パスが存在するかの確認
    log.Println("ノート追加 終了")
    log.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^" )
    return nil
    
} //--------------------------------------------

//HandleAddPagePost は /hello のPost時のHTMLデータ生成処理を行います。
func HandleAddPagePost(c echo.Context) error {
    log.Println("______________________________________________________________________________" )
    log.Println("ページ追加 開始")

    argMap := make(map[string]string)
    argMap["noteAddress"] = c.FormValue("note_address")

    addPage(argMap)

    log.Println("ページ追加 終了")
    log.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^" )

    return nil

} //--------------------------------------------


//HandleUpdateNotePost は /hello のPost時のHTMLデータ生成処理を行います。
func HandleUpdateNotePost(c echo.Context) error {
    log.Println("______________________________________________________________________________" )
    log.Println("ノート更新 開始")

    argMap := make(map[string]string)
    argMap["newNoteName"] = c.FormValue("new_note_name")
    argMap["postNoteID"] = c.FormValue("note_id")
    updateNote(argMap)

    log.Println("ノート更新 終了")
    log.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^" )

    return nil
} //--------------------------------------------

//HandleUpdatePagePost は /hello のPost時のHTMLデータ生成処理を行います。
func HandleUpdatePagePost(c echo.Context) error {
    log.Println("______________________________________________________________________________" )
    log.Println("ページ更新 開始")

    argMap := make(map[string]string)
    argMap["noteAddress"] = c.FormValue("note_address")
    argMap["pageID"]      = c.FormValue("page_id")
    argMap["pageTitle"]   = c.FormValue("page_title")
    argMap["pageBody"]    = c.FormValue("page_body")

    updatePage(argMap)


    log.Println("ページ更新 終了")
    log.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^" )

    return nil
} //--------------------------------------------

// HandleDeletePagePost のコメントアウト
func HandleDeletePagePost(c echo.Context) error {
    log.Println("______________________________________________________________________________" )
    log.Println("ページ削除 開始")


    argMap := make(map[string]string)
    argMap["noteAddress"] = c.FormValue("note_address")
    argMap["pageID"]      = c.FormValue("page_id")

    deletePage(argMap)
    
    log.Println("ページ削除 終了")
    log.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^" )

    return nil

} //--------------------------------------------

//HandleDeleteNotePost は /hello のPost時のHTMLデータ生成処理を行います。
func HandleDeleteNotePost(c echo.Context) error {
    log.Println("______________________________________________________________________________" )
    log.Println("ノート削除 開始")

    argMap := make(map[string]string)
    argMap["postNoteID"] = c.FormValue("note_id")

    log.Println("ノート削除 終了")
    log.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^" )

    return nil
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






// 各HTMLテンプレートに共通レイアウトを適用した結果を保存します（初期化時に実行）。
func loadTemplates() {
    var baseTemplate = "templates/layout.html"
    templates = make(map[string]*template.Template)
    templates["loadpage"]       = template.Must(template.ParseFiles(baseTemplate, "templates/loadpage.html"))
} //--------------------------------------------



//--------------------------------------------
func endProcess() {
    os.Exit(0)
} //--------------------------------------------


//--------------------------------------------
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
    } //--------------------------------------------



