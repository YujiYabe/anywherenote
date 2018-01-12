package main
import (
    "os"
    "log"
    "time"

)

func checkConfig() {
    //=============================================
    // 初期設定
    // 設定DB読み込み
    _, err := os.Stat( dataDirName + directorySeparator + confDBName )
    if err != nil {
        MakeConfDB()
    }

} //--------------------------------------------


func printEventLog(flagName string , message string) {

    separateLine := "------------------------------------------------------------------------------"
    if flagName == "start" {
		log.Println( separateLine )
		log.Println(message)
	}else if flagName == "end" {
		log.Println(message)
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


