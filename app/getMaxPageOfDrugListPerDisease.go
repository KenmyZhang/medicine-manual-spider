package app

import (
	"regexp"
	"fmt"
	"strings"
)

var (
/*
URL: 
    http://ypk.39.net/search/不孕症-p0/

查找内容：
    <a  target='_self' href='/search/不孕症-p15/' class='last'>尾页</a>

匹配详解：
    

*/
    allPage       = regexp.MustCompile(`<a  target='_self' href='/search/[.\x{4e00}-\x{9fa5}0-9]+-p[0-9]/'[ ]*class='last'>尾页</a>`)
    allPagePrefix = regexp.MustCompile(`<a  target='_self' href='/search/[.\x{4e00}-\x{9fa5}0-9]+-p`)
    allPageSuffix = regexp.MustCompile(`/'[ ]*class='last'>\x{5c3e}\x{9875}</a>`)
)

func GetDrugListMaxPage(diagNameChan chan string, diagNameAndPageChan chan *DiagNameAndPage){
    for{
    	select {
    		case diagName := <-diagNameChan:
	        	url := "http://ypk.39.net/search/" + diagName + "/"
	        	fmt.Println(url + "; GetDrugListMaxPage begin")
	        	body, err := httpGet(url, true)
	        	if  err != nil {
	            	fmt.Println("getDrugListMaxPage" + "app.get_drug_list_max_page.http_get.app.error "+ ",url:" + url + ", " + err.Error())
	        		return
	        	}
	        	fmt.Println(url + "; GetDrugListMaxPage end")

	        	diagNameAndPage := &DiagNameAndPage{}
	        	numPage := ""
	        	numPage = allPage.FindString(body)
	        	numPage = allPagePrefix.ReplaceAllString(numPage, "")
	        	numPage = allPageSuffix.ReplaceAllString(numPage, "")
	        	if strings.TrimSpace(numPage) == "" {
	            	numPage = "0"   
	        	}
	        	diagNameAndPage.Name = diagName
	        	diagNameAndPage.Page = numPage
	        	diagNameAndPageChan <- diagNameAndPage
	        	fmt.Println(url + "; diagNames:" + diagName)
	        	fmt.Println(url + "; MaxPage:" + numPage)
        }
    }
    return
}