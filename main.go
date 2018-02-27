package main 

import (
    "os"
    "fmt"
    "os/signal"
    "strconv"
    "time"
    "github.com/KenmyZhang/medicine-manual-spider/app"
    l4g "github.com/alecthomas/log4go"
)

var URL = "http://ypk.39.net/AllCategory"

func main() {
    l4g.AddFilter("stdout", l4g.ERROR, l4g.NewConsoleLogWriter())
    cleanupDone := make(chan bool)
   
//    catoNameChan := make(chan string, 100)
//    diagNameChan := make(chan string, 100)
    //drugNumChan  := make(chan string, 100)
    /*
    diagNameAndPageChan  := make(chan *app.DiagNameAndPage, 100)    
    go app.GetCato(URL, catoNameChan)
    go app.GetDiag(catoNameChan, diagNameChan)
    go app.GetDrugListMaxPage(diagNameChan, diagNameAndPageChan)
    go app.GetDrugNums(diagNameAndPageChan, drugNumChan)
    */

    //go rangeDrugNum(drugNumChan)

    //11111111111111111 app.GetProductSizeAndPriceRoutine(drugNumChan, cleanupDone)
    //22222222222222222 go SpyAllMedicineManual(drugNumChan, cleanupDone)
    //go app.GetProductSizeAndPriceFromYaoFang(drugNumChan, cleanupDone)
    //go app.SpyMedicineProductPriceFromJiaGe()
    go app.SpyProductPriceFrom360kad()
    Stop(cleanupDone)

}

func rangeDrugNum(drugNumChan chan string) {
    /*
    for medicine_manuals
    //for i := 500000; i <= 900000; i ++ {
    //for i := 500000; i >= 0; i-- {
    //  for i := 1000000000; i <= 1000100000; i++ {  
    */

    /*for      product
    //for i := 0; i <= 600000; i++ { 
    //for i := 600000; i <= 1229408 ; i++ { 
    */    
    //for i := 0; i <= 341815; i++ { 
    for i := 157772; i <= 1000000; i++ { 
        time.Sleep(200 * time.Millisecond)
        drugNumChan <- strconv.Itoa(i)
    }
}

func SpyAllMedicineManual(drugNums chan string, cleanupDone chan bool) {
    for {
        select {
            case numStr := <-drugNums:
                time.Sleep(50 * time.Millisecond)
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
