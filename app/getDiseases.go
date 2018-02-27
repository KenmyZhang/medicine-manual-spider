package app

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	/*
	   URL:
	       http://ypk.39.net/EffectSearch.aspx?ef=妇科用药
	   查找内容：
	       <li><a target="_blank" href="/search/产后血瘀/">产后血瘀</a></li>

	   <li><a href="/biyun/" target="_blank">避孕药</a></li>

	   匹配详解：

	*/
	diag       = regexp.MustCompile(`<li><a target="_blank" href="/search/[.\x{4e00}-\x{9fa5}]+/">`)
	diagPrefix = regexp.MustCompile(`<li><a target="_blank" href="/search/`)
	diagSuffix = regexp.MustCompile(`/">`)
)

func GetDiag(catoNameChan chan string, diagNameChan chan string) {
	for {
		select {
		case urlName := <-catoNameChan:
			url := "http://ypk.39.net/EffectSearch.aspx?ef=" + strings.TrimSpace(urlName)
			fmt.Println("httpGET  GetDiag" + url + "begin")
			body, err := httpGet(url, true)
			if err != nil {
				fmt.Println("getDiag", "app.get_diag.http_get.app_error", nil, "urlName:"+urlName+", "+err.Error())
				return
			}
			fmt.Println("httpGET  GetDiag" + url + "begin")
			fmt.Println("match  GetDiag" + url + "begin")
			diagMatches := diag.FindAllString(body, -1)
			for _, diagName := range diagMatches {
				diagName = diagPrefix.ReplaceAllString(diagName, "")
				diagName = diagSuffix.ReplaceAllString(diagName, "")
				diagNameChan <- diagName
			}
			fmt.Println("match  GetDiag" + url + "end")
		}
	}
	return
}
