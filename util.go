package main
import (
    "os"
    "log"
    "time"

)


func osCheckFile( targetFile string ) error {

    _, err := os.Stat( targetFile )

    return err

} //--------------------------------------------



func checkConfig() {
    //=============================================
    // 初期設定
    // 設定DB読み込み
    err := osCheckFile( confDBAddress )

    if err != nil {
        makeConfDB()
    }

} //--------------------------------------------


func printEventLog(flagName string , message interface{}) {

    separateLine := "------------------------------------------------------------------------------"
    if flagName == "start" {
        log.Println( separateLine )
        log.Println(message)

    }else if flagName == "end" {
        log.Println(message)
        log.Println( separateLine )

    }else if flagName == "returnFuncStatus" {
        log.Println( separateLine )
        log.Println(message)
        log.Println( separateLine )
        
    }else{
        // for debug
        log.Println( separateLine )
        log.Println( separateLine )
        log.Println( separateLine )
        log.Println(flagName)
        log.Println(message)
        log.Println( separateLine )
        log.Println( separateLine )
        log.Println( separateLine )
    }

} //--------------------------------------------


//--------------------------------------------
func endProcess() {
    os.Exit(0)
} //--------------------------------------------

//--------------------------------------------
func calcTime() {




    //ブラウザが起動して、pingを発行するまでn秒まつ
    time.Sleep( userConfig.WaitSecondLiveCheck * time.Second )

    for {
        time.Sleep( userConfig.WaitSecondInterval * time.Second )

        t := time.Now()
        beforeTime := t.Add( time.Duration(1) * time.Second).Format(dateTimeFormat)
        now, _ := time.Parse(dateTimeFormat ,beforeTime)

        old, _ := time.Parse(dateTimeFormat , recieveString)

        if !old.After(now) {    // old <= now --- ! old > now
            if userConfig.IsEnableAppMode {
                log.Println("アプリ終了")
                endProcess()
            }
        }
    }
} //--------------------------------------------


//--------------------------------------------
func osCreateFile( targetFileName string ) {
    file, err := os.OpenFile( targetFileName , os.O_WRONLY|os.O_CREATE , 0666 )

    if err != nil {
        //エラー処理
        log.Fatal(err)
    }
    defer file.Close()

} //--------------------------------------------

