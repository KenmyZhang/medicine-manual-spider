package app

import (
	"fmt"
    "regexp"
    "strconv"
)

var (
/*
URL:
    http://ypk.39.net/search/不孕症-p0/

查找内容：   
    <strong><a href="/500376/" >

匹配详情：

*/
                                        //<strong><a href="/1000004235/" >                                                  
   drugNum         =  regexp.MustCompile(`<strong><a href="/[0-9a-zA-Z]+/"`)
   drugNumPrefix   =  regexp.MustCompile(`<strong><a href="/`)
   drugNumSuffix   =  regexp.MustCompile(`/"`)
)

type DiagNameAndPage struct {
    Name string
    Page string 
}

func GetDrugNums(diagNameAndPageChan chan *DiagNameAndPage,drugNumChan chan string) {
    for{
    	select {
    		case diagNameAndPage := <-diagNameAndPageChan:
	        page, _ := strconv.Atoi(diagNameAndPage.Page)
    	    for i := 0; i <= page; i++ {
            	url := "http://ypk.39.net/search/" + diagNameAndPage.Name +"-p" +strconv.Itoa(i) + "/"
            	body, err := httpGet(url)
            	if  err != nil {
                    fmt.Println("getDrugNums" + "app.get_drug_nums.http_get.app.error:" + "diagNameAndPage.Name:" +  diagNameAndPage.Name + ",  page:" + strconv.Itoa(i) + "," + err.Error())
            	    return
                }

            	drugNumMatches :=  drugNum.FindAllString(body, -1)
            	for _, drugNum := range drugNumMatches {
                	drugNum = drugNumPrefix.ReplaceAllString(drugNum, "")
                	drugNum = drugNumSuffix.ReplaceAllString(drugNum, "")
                	drugNumChan <- drugNum
                	fmt.Println(url + "; drugNum:" + drugNum)
            	}
        	}
    	}
	}
    return
}