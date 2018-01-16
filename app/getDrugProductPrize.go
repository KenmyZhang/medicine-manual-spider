package app

import (
    "fmt"
    "regexp"
    "github.com/KenmyZhang/medicine-manual-spider/model"
)

var (
/*
    从下面字符串中查找出价格2698.00元:
    <span class="new-price-tip" data-v-61fefe86>健客价</span><span class="new-price-icon" data-v-61fefe86>￥</span><span data-v-61fefe86>2698.00</span>
*/  
    prize       = regexp.MustCompile(`<span class="new-price-tip" data-v-61fefe86>健客价</span><span class="new-price-icon" data-v-61fefe86>￥</span><span data-v-61fefe86>[0-9.]+</span>`)
    prizePrefix = regexp.MustCompile(`<span class="new-price-tip" data-v-61fefe86>健客价</span><span class="new-price-icon" data-v-61fefe86>￥</span><span data-v-61fefe86>`)
    prizeSuffix = regexp.MustCompile(`</span>`)

/*
    从下面字符串中查找出药品名称：
    <div class="product-name" data-v-61fefe86><span data-v-61fefe86>阿胶(东阿阿胶)</span>
*/
    productName       = regexp.MustCompile(`<div class="product-name" data-v-61fefe86><span data-v-61fefe86>[\x{4e00}-\x{9fa5}\(\) ]+</span>`)
    productNamePrefix = regexp.MustCompile(`<div class="product-name" data-v-61fefe86><span data-v-61fefe86>`)
    productNameSuffix = regexp.MustCompile(`</span>`)

/*
    从下面字符串中查找出药品的批准文号：
    批准文号：
    <span class="introduct-bd" data-v-80bdb96c>国药准字Z37021368</span></p><p class="introduct-hd" data-v-80bdb96c>
*/
    approveNum       = regexp.MustCompile(`批准文号：[\s]+<span class="introduct-bd" data-v-80bdb96c>国药准字[A-Za-z0-9]+</span></p><p class="introduct-hd" data-v-80bdb96c>`)
    approveNumPrefix = regexp.MustCompile(`批准文号：[\s]+<span class="introduct-bd" data-v-80bdb96c>`)
    approveNumSuffix = regexp.MustCompile(`</span></p><p class="introduct-hd" data-v-80bdb96c>`)

/*
从下面字符串中查找出药品的产品规格：
  <dl class="product-sizes" data-v-61fefe86><dt data-v-61fefe86>规格：</dt><dd data-v-61fefe86><ul data-v-61fefe86><li class="product-size active" data-v-61fefe86><a href="/product/6902.html" class="router-link-exact-active router-link-active" data-v-61fefe86>500g</a></li><li class="product-size" data-v-61fefe86><a href="/product/154505.html" data-v-61fefe86>250g</a></li><li class="product-size" data-v-61fefe86><a href="/product/178176.html" data-v-61fefe86>125g</a></li></ul></dd></dl>
*/
    size       = regexp.MustCompile(`(?U)<dl class="product-sizes" data-v-61fefe86><dt data-v-61fefe86>规格：.*</dl>`) 
    sizePrefix = regexp.MustCompile(`<dl class="product-sizes" data-v-61fefe86><dt data-v-61fefe86>规格：`)
    sizeSuffix = regexp.MustCompile(`</dl>`)

/*
从下面字符串中查找出药品的当前的产品规格：
  <dl class="product-sizes" data-v-61fefe86><dt data-v-61fefe86>规格：</dt><dd data-v-61fefe86><ul data-v-61fefe86><li class="product-size active" data-v-61fefe86><a href="/product/6902.html" class="router-link-exact-active router-link-active" data-v-61fefe86>500g</a></li><li class="product-size" data-v-61fefe86><a href="/product/154505.html" data-v-61fefe86>250g</a></li><li class="product-size" data-v-61fefe86><a href="/product/178176.html" data-v-61fefe86>125g</a></li></ul></dd></dl>
*/
    currentSize       = regexp.MustCompile(`(?U)<li class="product-size active"[\s\S]+</li>`) 
    currentSizePrefix = regexp.MustCompile(`(?U)<li class="product-size active"[\s\S]+class="router-link-exact-active router-link-active" data-v-61fefe86>`)
    currentSizeSuffix = regexp.MustCompile(`</a></li>`)

/*
从下面字符串中查找出药品的所有产品规格：
  <dl class="product-sizes" data-v-61fefe86><dt data-v-61fefe86>规格：</dt><dd data-v-61fefe86><ul data-v-61fefe86><li class="product-size active" data-v-61fefe86><a href="/product/6902.html" class="router-link-exact-active router-link-active" data-v-61fefe86>500g</a></li><li class="product-size" data-v-61fefe86><a href="/product/154505.html" data-v-61fefe86>250g</a></li><li class="product-size" data-v-61fefe86><a href="/product/178176.html" data-v-61fefe86>125g</a></li></ul></dd></dl>
*/
    perSize       = regexp.MustCompile(`(?U)<a href=[\s\S]+</a>`) 
    perSizePrefix = regexp.MustCompile(`(?U)<a href=[\s\S]+>`)
    perSizeSuffix = regexp.MustCompile(`</a>`)

)

func GetProductSizeAndPriceRoutine(numChan chan string, cleanupDone chan bool) {
  for{
    select {
    case num := <-numChan:
      go getProductSizeAndPrice(num)
    case <-cleanupDone:
    }
  }

}

func getProductSizeAndPrice(num string) {
    productSizeAndPrize := &model.ProductSizeAndPrize{}
    url := "https://m.jianke.com/product/" + num + ".html"
    fmt.Println("httpGET  GetDiag" + url + ", begin")
    body, err := httpGet(url, false)
    if err != nil {
       fmt.Println("getDiag", "app.get_diag.http_get.app_error", nil, "num:" + num + ", " + err.Error())
       return
    }
    fmt.Println("httpGET  GetDiag" + url + ", success")

    productSizeAndPrize.Num = num

    fmt.Println("body：", body)

    prizeStr := prize.FindString(body)
    prizeStr = prizePrefix.ReplaceAllString(prizeStr, "")
    prizeStr = prizeSuffix.ReplaceAllString(prizeStr, "")
    productSizeAndPrize.Price = prizeStr
    fmt.Println("门店价格：", prizeStr)

    productNameStr := productName.FindString(body)
    productNameStr = productNamePrefix.ReplaceAllString(productNameStr, "")
    productNameStr = productNameSuffix.ReplaceAllString(productNameStr, "")
    fmt.Println("药品名称", productNameStr)
    productSizeAndPrize.Name = productNameStr

    approveNumStr := approveNum.FindString(body)
    approveNumStr = approveNumPrefix.ReplaceAllString(approveNumStr, "")
    approveNumStr = approveNumSuffix.ReplaceAllString(approveNumStr, "")
    fmt.Println("批准文号：",approveNumStr)
    productSizeAndPrize.ApprovalNumber = approveNumStr

    currentSizeStr := currentSize.FindString(body)
    currentSizeStr = currentSizePrefix.ReplaceAllString(currentSizeStr, "")
    currentSizeStr = currentSizeSuffix.ReplaceAllString(currentSizeStr, "")
    fmt.Println("规格：", currentSizeStr)
    productSizeAndPrize.CurrentSize = currentSizeStr

    sizeStr := size.FindString(body)
    sizeStr = sizePrefix.ReplaceAllString(sizeStr, "")
    sizeStr = sizeSuffix.ReplaceAllString(sizeStr, "")
    matches :=  perSize.FindAllString(sizeStr, -1)
    for _, match := range matches {
      match = perSizePrefix.ReplaceAllString(match, "")
      match = perSizeSuffix.ReplaceAllString(match, "")
      productSizeAndPrize.AllSize = match
      fmt.Println("cato:", match)
    }
    if productSizeAndPrize.Name != "" {
        SaveProductSizeAndPrize(productSizeAndPrize)
    }
}