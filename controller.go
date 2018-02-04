package main

import (
	"html/template"
	"io"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"encoding/json"
	// "fmt"
	"io/ioutil"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/skratchdot/open-golang/open"
)

// グローバル変数
var (
	recieveString      string // ブラウザからのハートビート受け取り
	directorySeparator = "/"  // linux separator
	confDBAddress      string // 設定DBパス
	userConfig         = new(UserConfig)
)

// 定数
const useDBMSName = "sqlite3"          // 使用DBMS
const dataDirName = "data"             // htmlやjs、DBファイル等格納先ディレクトリ
const confDBName = "config.db"         // ローカル設定DBファイル名
const userConfFile = "userconfig.json" // ローカル設定DBファイル名

const fileDirName = "file"   // ページ内にアップロードするディレクトリ名
const noteDBName = "note.db" // ノート保存先DBファイル名

const dateTimeFormat = "2006-01-02 15:04:05" // 日時フォーマット

// レイアウト適用済のテンプレートを保存するmap
var templates map[string]*template.Template

// Template はHTMLテンプレートを利用するためのRenderer Interface
type Template struct{}

// Render はHTMLテンプレートにデータを埋め込んだ結果をWriterに書き込み
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return templates[name].ExecuteTemplate(w, "layout.html", data)
} //--------------------------------------------

// 初期処理
func init() {

	// ディレクトリセパレータ
	if runtime.GOOS == "windows" {
		directorySeparator = "\\"
	}

	// ユーザ設定ファイルの読み込み
	userConfigFile, err := ioutil.ReadFile(dataDirName + directorySeparator + "public" + directorySeparator + userConfFile)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(userConfigFile, &userConfig)

	confDBAddress = dataDirName + directorySeparator + confDBName

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
	// e.Use(middleware.Logger())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	e.Use(middleware.Recover())

	// 静的ファイルのパスを設定
	e.Static("/public", "data/public")
	e.File("/favicon.ico", "data/public/favicon.ico")

	// ノートアドレスとノートIDを紐づけて静的ファイルパスを設定
	rValue := getAllNoteAddress()

	for _, value := range rValue {
		v := uint64(value.ID)
		noteID := strconv.FormatUint(v, 10)

		e.Static("/"+noteID, value.NoteAddress+directorySeparator+fileDirName)
	}

	// 各ルーティングに対するハンドラを設定
	e.GET("/", LoadPageGet)
	e.GET("/livecheck", LiveCheckGet)
	e.POST("/addnote", AddNotePost)
	e.POST("/addpage", AddPagePost)
	e.POST("/updatenote", UpdateNotePost)
	e.POST("/updatepage", UpdatePagePost)
	e.POST("/deletenote", DeleteNotePost)
	e.POST("/deletepage", DeletePagePost)
	e.POST("/uploadfile", UploadFilePost)

	// アプリモードが有効であればローカルのブラウザを開く
	// ブラウザからのハートビートを受け取る
	// 無効であればサーバモードとして起動
	if userConfig.IsEnableAppMode {
		open.Run("http://localhost:" + userConfig.UsePortNumber)
		go calcTime()
	}

	// サーバーを開始
	e.Logger.Fatal(e.Start(":" + userConfig.UsePortNumber))
} //--------------------------------------------

// UploadFilePost のコメントアウト
func UploadFilePost(c echo.Context) error {
	// Read form fields

	noteAddress := c.FormValue("note_address")
	pageID := c.FormValue("page_id")

	//-----------
	// Read file
	//-----------

	// Source
	file, err := c.FormFile(fileDirName)
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination // ノート配下のfileディレクトリに作成yyyymmdd-hhmmss形式
	dst, err := os.Create(noteAddress + directorySeparator + fileDirName + directorySeparator + file.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	// return c.HTML(http.StatusOK, fmt.Sprintf("<p>File %s uploaded successfully with fields </p>", file.Filename))

	//対象のページ内容を呼び出し、対象のファイルを追記した内容で再書き込み

	addFile := "<br>"

	imgExtList := []string{".bmp", ".gif", ".png", ".jpg", ".jpeg", ".BMP", ".GIF", ".PNG", ".JPG", ".JPEG"}
	pos := strings.LastIndex(file.Filename, ".")

	if contains(imgExtList, file.Filename[pos:]) {
		addFile = addFile + "<img class='image_style' src='//note_id///" + file.Filename + "'>"
	} else {
		// <div class="sam13"><a href="">テスト</a></div>

		// addFile = addFile + "<div class='file_style btn btn-warning'>" + file.Filename + "</div>"

		// addFile = addFile + "<a href='//note_id//' class='file_style'>" + file.Filename + "</a>"
		// addFile = addFile + "<a href='//note_id//' onclick='file_download(this);'>" + file.Filename + "</a>"
		// addFile = addFile + "<a  href='//note_id///" + file.Filename + "' onclick='fileDownload(this);'>" + file.Filename + "</a>"
		// addFile = addFile + "<a  href='.///note_id///" + file.Filename + "' download='" + file.Filename + "'>" + file.Filename + "</a>"

		addFile = addFile + "<a  href='//note_id///" + file.Filename + "' class='fileDownload'><div contenteditable='false' class='btn btn-warning'>" + file.Filename + "</div></a>"

	}
	addFile = addFile + "<br><br><br>"

	sndArg := make(map[string]string)
	sndArg["noteAddress"] = noteAddress
	sndArg["pageID"] = pageID
	sndArg["addFile"] = addFile

	addFileToPage(sndArg)

	var selectPosition = SelectPosition{}

	selectPosition.NoteID = convertStringToUint(c.FormValue("note_id"))
	selectPosition.PageID = convertStringToUint(c.FormValue("page_id"))

	returnValue := getData(selectPosition)

	printEventLog("end", "ページ更新 終了")

	return c.JSON(http.StatusCreated, returnValue)

}

// LoadPageGet のコメントアウト
func LoadPageGet(c echo.Context) error {
	printEventLog("start", "データ取得 開始")

	var selectPosition = SelectPosition{}

	selectPosition.NoteID = 0
	selectPosition.PageID = 0

	returnjson := getData(selectPosition)

	printEventLog("end", "データ取得 終了")
	return c.Render(http.StatusOK, "loadpage", returnjson)

} //--------------------------------------------

//AddNotePost のコメントアウト
// ディレクトリの存在確認
// DBファイルの存在確認
// ┣ファイルが存在しなければ新たに作成
// 設定DBにノートDBの情報を追加
func AddNotePost(c echo.Context) error {
	printEventLog("start", "ノート追加 開始")

	noteAddress := c.FormValue("note_address")

	// ディレクトリ 存在確認
	err := osCheckFile(noteAddress)
	if err != nil {
		return err
	}

	// DBファイルの存在確認
	notefullAddress := noteAddress + directorySeparator + noteDBName

	err2 := osCheckFile(notefullAddress)

	// ファイルが存在しなければ新たに作成
	if err2 != nil {
		dbApplyType(notefullAddress, &Note{})
	}

	filefullAddress := noteAddress + directorySeparator + fileDirName
	// ファイル用ディレクトリが存在するかを確認
	err3 := osCheckFile(filefullAddress)
	if err3 != nil {
		if err4 := os.Mkdir(filefullAddress, 0777); err != nil {
			printEventLog("error", err4)
		}
	}

	//____________________________________
	// 設定DBに追加
	confdb := setupDB(confDBAddress)
	defer confdb.Close()
	confdb.LogMode(true)

	var conf Conf

	conf.NoteName = c.FormValue("note_name")
	conf.NoteAddress = c.FormValue("note_address")
	conf.NoteStar = 3

	// INSERTを実行
	confdb.Create(&conf)

	// ★追加したアドレスのスタティックパスを追加
	// v := uint64(value.ID)
	// noteID := strconv.FormatUint(v, 10)
	// e.Static("/"+noteID, value.NoteAddress + directorySeparator + fileDirName)

	printEventLog("end", "ノート追加 終了")
	return nil

} //--------------------------------------------

// AddPagePost のコメントアウト
// 対象のノートにページを追加
// 返却jsonに対象のノートIDを追加
// 最新のデータセットを取得
// jsonを返却
func AddPagePost(c echo.Context) error {
	printEventLog("start", "ページ追加 開始")

	//対象のノートにページを追加
	sndArg := make(map[string]string)
	sndArg["noteAddress"] = c.FormValue("note_address")

	addPage(sndArg)

	// 返却jsonに対象のノートIDを追加
	var selectPosition = SelectPosition{}
	selectPosition.NoteID = convertStringToUint(c.FormValue("note_id"))

	// 最新のデータセットを取得
	returnValue := getData(selectPosition)

	printEventLog("end", "ページ追加 終了")

	return c.JSON(http.StatusCreated, returnValue)

} //--------------------------------------------

// UpdateNotePost のコメントアウト
// ノート情報の更新（設定DB内の対象ノート名の変更）
func UpdateNotePost(c echo.Context) error {
	printEventLog("start", "ノート更新 開始")

	sndArg := make(map[string]string)
	sndArg["postNoteID"] = c.FormValue("note_id")
	sndArg["noteName"] = c.FormValue("note_name")
	sndArg["noteStar"] = c.FormValue("note_star")
	updateNote(sndArg)

	printEventLog("end", "ノート更新 終了")
	return nil
} //--------------------------------------------

// UpdatePagePost  のコメントアウト
// ページの更新
// 設定DBにある更新したノートの更新日時を変更
// 返却jsonに対象のノートIDを追加
func UpdatePagePost(c echo.Context) error {
	printEventLog("start", "ページ更新 開始")

	// ページの更新
	sndArg := make(map[string]string)
	sndArg["noteAddress"] = c.FormValue("note_address")
	sndArg["pageID"] = c.FormValue("page_id")
	sndArg["pageTitle"] = c.FormValue("page_title")
	sndArg["pageBody"] = c.FormValue("page_body")

	updatePage(sndArg)

	// 設定DBにある更新したノートの更新日時を変更
	updateNoteFromPage(c.FormValue("note_address"))

	// 返却jsonに対象のノートIDを追加
	var selectPosition = SelectPosition{}
	selectPosition.NoteID = convertStringToUint(c.FormValue("note_id"))

	returnValue := getData(selectPosition)

	printEventLog("end", "ページ更新 終了")

	return c.JSON(http.StatusCreated, returnValue)

} //--------------------------------------------

// DeleteNotePost  のコメントアウト
// ノート情報の削除
func DeleteNotePost(c echo.Context) error {
	printEventLog("start", "ノート削除 開始")

	sndArg := make(map[string]string)
	sndArg["postNoteID"] = c.FormValue("note_id")
	deleteNote(sndArg)

	printEventLog("end", "ノート削除 終了")
	return nil
} //--------------------------------------------

// DeletePagePost のコメントアウト
// ページ情報の削除
func DeletePagePost(c echo.Context) error {
	printEventLog("start", "ページ削除 開始")

	sndArg := make(map[string]string)
	sndArg["noteAddress"] = c.FormValue("note_address")
	sndArg["pageID"] = c.FormValue("page_id")

	deletePage(sndArg)

	// 返却jsonに対象のノートIDを追加
	var selectPosition = SelectPosition{}
	selectPosition.NoteID = convertStringToUint(c.FormValue("note_id"))

	returnValue := getData(selectPosition)

	printEventLog("end", "ページ削除 終了")

	return c.JSON(http.StatusCreated, returnValue)

} //--------------------------------------------

// LiveCheckGet のコメントアウト
func LiveCheckGet(c echo.Context) error {
	// //現在時刻取得
	t := time.Now()

	//現在からn秒後の時刻を取得
	afterTime := t.Add(time.Duration(userConfig.WaitSecondLiveCheck) * time.Second).Format(dateTimeFormat)

	recieveString = afterTime

	//チャネルへ時刻情報を送信
	// liveChannel <- c.FormValue("expireLiveTime")

	return nil
} //--------------------------------------------

// 各HTMLテンプレートに共通レイアウトを適用した結果を保存します（初期化時に実行）。
func loadTemplates() {
	var baseTemplate = "data/template/layout.html"
	templates = make(map[string]*template.Template)
	templates["loadpage"] = template.Must(template.ParseFiles(baseTemplate, "data/template/loadpage.html"))
} //--------------------------------------------
