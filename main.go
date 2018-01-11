package main 

import (
    "os"
    "fmt"
    "os/signal"
    //"runtime"
    "github.com/KenmyZhang/medicine-manual-spider/app"
)

var URL = "http://ypk.39.net/AllCategory"

func main() {
    cleanupDone := make(chan bool)
    catoNameChan := make(chan string, 100)
    diagNameChan := make(chan string, 100)
    drugNumChan  := make(chan string, 100)
    diagNameAndPageChan  := make(chan *app.DiagNameAndPage, 100)    
    go app.GetCato(URL, catoNameChan)
    go app.GetDiag(catoNameChan, diagNameChan)
    go app.GetDrugListMaxPage(diagNameChan, diagNameAndPageChan)
    go app.GetDrugNums(diagNameAndPageChan, drugNumChan)
    go SpyAllMedicineManual(drugNumChan, cleanupDone)
    Stop(cleanupDone)

}


func SpyAllMedicineManual(drugNums chan string, cleanupDone chan bool) {
    for {
        select {
            case numStr := <-drugNums:
                go app.SpyMedicineManual(numStr)
            // 从ch中读取到数据
            case <-cleanupDone:
            // 一直没有从ch中读取到数据,但从cleanupDone中读取到了数据            
        }
    }
}

func Stop(cleanupDone chan bool) {
    signalChan := make(chan os.Signal, 1)
    signal.Notify(signalChan, os.Interrupt)
    go func() {
        for _ = range signalChan {
            cleanUp()
            cleanupDone <- true
        }
    }()
    <-cleanupDone
}

func cleanUp() {
    app.MgoSession.Close()
    fmt.Println("清理...\n")
}
