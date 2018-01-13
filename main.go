package main

import (
    "html/template"
    "io"
    "net/http"
    "runtime"
    "time"
    "github.com/labstack/echo"
    "github.com/labstack/echo/middleware"


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
const isEnableAppMode = false // for debug
const waitSecondLiveCheck = 8
const waitSecondInterval = 1 
const dateTimeFormat = "2006-01-02 15:04:05"


const useDBMSName = "sqlite3"
const confDBName  = "conf.db"
const noteDBName  = "note.db"

const dataDirName  = "data"

const usePortNumber = "3000"



// DataSet DBファイルの情報とノート情報のセット
type DataSet struct {
    NoteDBID         uint   `json:"NoteDBID"`
    NoteDBName       string `json:"NoteDBName"`
    NoteDBAddress    string `json:"NoteDBAddress"`
    NoteDBUpdateTime time.Time
    List             []Note `json:"list"`
}


// SelectPosition
type SelectPosition struct {
    NoteID  uint   `json:"NoteID"`
    PageID  uint   `json:"PageID"`

}

// ReturnValue 戻り値とDataSetのセット Key0がリターンコード 
type ReturnValue struct {
    Key0   string           `json:"key0"`
    Key1   []DataSet        `json:"key1"`
    Key2   SelectPosition   `json:"key2"`
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
    e.Static("/public", "data/public")
    e.File("/favicon.ico", "data/public/favicon.ico")

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
    printEventLog("start" , "データ取得 開始")

    var selectPosition = SelectPosition{}

    if c.FormValue("select_note_id") == "" {
        selectPosition.NoteID = 0
    }else{
        selectPosition.NoteID = strconv.Atoi(c.FormValue("select_note_id") )
    }
    
    if c.FormValue("select_page_id") == "" {
        selectPosition.PageID = 0
    }else{
        selectPosition.PageID = c.FormValue("select_page_id") 
    }


    returnjson := getData(selectPosition)

    printEventLog("end" , "データ取得 終了")
    return c.Render(http.StatusOK, "loadpage", returnjson)

} //--------------------------------------------

//HandleAddNotePost は /hello のPost時のHTMLデータ生成処理を行います。
func HandleAddNotePost(c echo.Context) error {
    printEventLog("start" , "ノート追加 開始")

    argMap := make(map[string]string)
    argMap["newNoteName"]    = c.FormValue("new_note_name")
    argMap["newNoteAddress"] = c.FormValue("new_note_address")

    // ★ todo 対象ディレクトリが存在しない場合エラーを返却
    returnvalue := addNote(argMap)
    printEventLog( "returnFuncStatus" , returnvalue )


    // タスク パスが存在するかの確認
    printEventLog( "end" , "ノート追加 終了")
    return nil
} //--------------------------------------------

//HandleAddPagePost は /hello のPost時のHTMLデータ生成処理を行います。
func HandleAddPagePost(c echo.Context) error {
    printEventLog("start" , "ページ追加 開始")

    argMap := make(map[string]string)
    argMap["noteAddress"] = c.FormValue("note_address")

    addPage(argMap)

    printEventLog("end" , "ページ追加 終了")
    return nil
} //--------------------------------------------


//HandleUpdateNotePost は /hello のPost時のHTMLデータ生成処理を行います。
func HandleUpdateNotePost(c echo.Context) error {
    printEventLog("start" , "ノート更新 開始")

    argMap := make(map[string]string)
    argMap["newNoteName"] = c.FormValue("new_note_name")
    argMap["postNoteID"]  = c.FormValue("note_id")
    updateNote(argMap)

    printEventLog("end" , "ノート更新 終了")
    return nil
} //--------------------------------------------

//HandleUpdatePagePost は /hello のPost時のHTMLデータ生成処理を行います。
func HandleUpdatePagePost(c echo.Context) error {
    printEventLog("start" , "ページ更新 開始")

    argMap := make(map[string]string)
    argMap["noteAddress"] = c.FormValue("note_address")
    argMap["pageID"]      = c.FormValue("page_id")
    argMap["pageTitle"]   = c.FormValue("page_title")
    argMap["pageBody"]    = c.FormValue("page_body")

    updatePage(argMap)

    printEventLog("end" , "ページ更新 終了")
    return nil
} //--------------------------------------------

// HandleDeletePagePost のコメントアウト
func HandleDeletePagePost(c echo.Context) error {
    printEventLog("start" , "ページ削除 開始")

    argMap := make(map[string]string)
    argMap["noteAddress"] = c.FormValue("note_address")
    argMap["pageID"]      = c.FormValue("page_id")

    deletePage(argMap)
    
    printEventLog("end" , "ページ削除 終了")
    return nil
} //--------------------------------------------

//HandleDeleteNotePost は /hello のPost時のHTMLデータ生成処理を行います。
func HandleDeleteNotePost(c echo.Context) error {
    printEventLog("start" , "ノート削除 開始")

    argMap := make(map[string]string)
    argMap["postNoteID"] = c.FormValue("note_id")
    deleteNote(argMap)

    printEventLog("end" , "ノート削除 終了")
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
    var baseTemplate = "data/templates/layout.html"
    templates = make(map[string]*template.Template)
    templates["loadpage"]  = template.Must(template.ParseFiles(baseTemplate, "data/templates/loadpage.html"))
} //--------------------------------------------




