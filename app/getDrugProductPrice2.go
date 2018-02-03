package app

import (
    "fmt"
    "regexp"
    "github.com/KenmyZhang/medicine-manual-spider/model"
)

var (
    price2       = regexp.MustCompile(`<input type="hidden" id="price" value="[\s\S]*?">`)
    price2Prefix = regexp.MustCompile(`<input type="hidden" id="price" value="`)
    price2Suffix = regexp.MustCompile(`">`)

    product2       = regexp.MustCompile(`<input type="hidden" id="names" value="[\s\S]*?">`)
    product2Prefix = regexp.MustCompile(`<input type="hidden" id="names" value="`)
    product2Suffix = regexp.MustCompile(`">`)

    approveNum2       = regexp.MustCompile(`批准文号：<a href="http://app1.sfda.gov.cn/datasearch/face3/dir.html">[\s\S]*?</a><br>`)
    approveNum2Prefix = regexp.MustCompile(`批准文号：<a href="http://app1.sfda.gov.cn/datasearch/face3/dir.html">`)
    approveNum2Suffix = regexp.MustCompile(`</a><br>`)

    size2       = regexp.MustCompile(`规 格：[\s\S]+?<br>`) 
    size2Prefix = regexp.MustCompile(`规 格：`)
    size2Suffix = regexp.MustCompile(`<br>`)

    //manufacturer2 = regexp.MustCompile(`生产厂家：<a style="color: #0c69ae;" href="http://www.yaofang.cn/c/search?s_words=[\s\S]+?">`)
    manufacturer2 = regexp.MustCompile(`生产厂家：<a style="color: #0c69ae;" href="http://www.yaofang.cn/c/search\?s_words=[\s\S]*?">`)
    manufacturer2Prefix = regexp.MustCompile(`生产厂家：<a style="color: #0c69ae;" href="http://www.yaofang.cn/c/search\?s_words=`)
    manufacturer2Suffix = regexp.MustCompile(`">`)
    otherChar     = regexp.MustCompile(`&nbsp;`)
)


func GetProductSizeAndPriceRoutine2(numChan chan string, cleanupDone chan bool) {
  for{
    select {
    case num := <-numChan:
      go getProductSizeAndPrice2(num)
    case <-cleanupDone:
    }
  }

}

func getProductSizeAndPrice2(num string) {
    product2SizeAndPrize := &model.ProductSizeAndPrize{}
    url := "http://www.yaofang.cn/goods-" + num + ".html"
    fmt.Println("httpGET  GetDiag" + url + ", begin")
    body, err := httpGet(url, false)
    if err != nil {
       fmt.Println("getDiag", "app.get_diag.http_get.app_error", nil, "num:" + num + ", " + err.Error())
       return
    }
    fmt.Println("httpGET  GetDiag" + url + ", success")

    product2SizeAndPrize.Num = num

    price2Str := price2.FindString(body)
    price2Str = price2Prefix.ReplaceAllString(price2Str, "")
    price2Str = price2Suffix.ReplaceAllString(price2Str, "")
    fmt.Println("门店价格：", price2Str)
    product2SizeAndPrize.Price = price2Str

 
    product2NameStr := product2.FindString(body)
    product2NameStr = product2Prefix.ReplaceAllString(product2NameStr, "")
    product2NameStr = product2Suffix.ReplaceAllString(product2NameStr, "")
    fmt.Println("药品名称",product2NameStr)
    product2SizeAndPrize.Name = product2NameStr
 
    approveNum2Str := approveNum2.FindString(body)
    approveNum2Str = approveNum2Prefix.ReplaceAllString(approveNum2Str, "")
    approveNum2Str = approveNum2Suffix.ReplaceAllString(approveNum2Str, "")
    fmt.Println("批准文号：",approveNum2Str)
    product2SizeAndPrize.ApprovalNumber = approveNum2Str

    currentSizeStr := size2.FindString(body)
    currentSizeStr = size2Prefix.ReplaceAllString(currentSizeStr, "")
    currentSizeStr = size2Suffix.ReplaceAllString(currentSizeStr, "")
    currentSizeStr = otherChar.ReplaceAllString(currentSizeStr, "")
    fmt.Println("规格：", currentSizeStr)
    product2SizeAndPrize.CurrentSize = currentSizeStr 

    product2ManufacturerStr := manufacturer2.FindString(body)
    product2ManufacturerStr = manufacturer2Prefix.ReplaceAllString(product2ManufacturerStr, "")
    product2ManufacturerStr = manufacturer2Suffix.ReplaceAllString(product2ManufacturerStr, "")
    fmt.Println("生产厂商：", product2ManufacturerStr)
    product2SizeAndPrize.Manufacturer = product2ManufacturerStr  

    if product2SizeAndPrize.Name != "" {
        SaveProductSizeAndPrize(product2SizeAndPrize)
    }   
}