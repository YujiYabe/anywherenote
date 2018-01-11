package main
import (
    "os"
    "log"

)

func checkConfig() {
    //=============================================
    // 初期設定
    // 設定DB読み込み
    _, err := os.Stat(confDBName)
    if err != nil {
        MakeConfDB()
    }

} //--------------------------------------------


func printEvent(flagName string , message string) {
	if flagName == "start" {
		log.Println("______________________________________________________________________________" )
		log.Println(message)
	}else if flagName == "end" {
		log.Println(message)
		log.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^" )
	}




} //--------------------------------------------
