package app

import (
    "regexp"
    "fmt"
)

var (
/*
URL: 
      http://ypk.39.net/AllCategory

查找内容：
                        <dt>
                                <a rel="canonical" href="/nanke/" target="_blank">
                                    头痛用药</a>
                        </dt>

匹配详解： 
    \s                   匹配任何空白字符，包括空格、制表符、换页符等等。等价于 [ \f\n\r\t\v]。
    \s{1,}               匹配任何空白字符，一个以上
        [a-z]                匹配 'a' 到 'z' 范围内的任意小写字母字符。
    [a-z]+               匹配 'a' 到 'z' 范围内的任意小写字母字符,一个以上
        [\x{4e00}-\x{9fa5}]  匹配中文字符
*/
    cato       = regexp.MustCompile(`<dt>\s{1,}<a rel="canonical" href="/[a-z]+/" target="_blank">(\s{1,}[.\x{4e00}-\x{9fa5}]+)</a>`)
    catoPrefix = regexp.MustCompile(`<dt>\s{1,}<a rel="canonical" href="/[a-z]+/" target="_blank">\s{1,}`)
    catoSuffix = regexp.MustCompile(`\s{0,}</a>`)
)


func GetCato(url string, catoNameChan chan string) {
    fmt.Println("httpGET GetCato" + url + "begin")
    body, err := httpGet(url, true)
    if err != nil {
        fmt.Println("app.getCato.http_get.app_error, allCategoryUrl:" + url + ", " + err.Error())
        return
    }
    fmt.Println("httpGET GetCato" + url + "end")
    matches :=  cato.FindAllString(body, -1)
    fmt.Println(url + ";GetCato matches" + "begin")
    for index, cato := range matches {
      cato = catoPrefix.ReplaceAllString(cato, "")
      cato = catoSuffix.ReplaceAllString(cato, "")
      catoNameChan <- cato
      fmt.Println("index:", index) 
      fmt.Println("cato:", cato)
    }
    fmt.Println(url + ";GetCato matches" + "end")
    return
}
