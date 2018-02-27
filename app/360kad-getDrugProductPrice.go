package app

import (
	"regexp"
  "time"
  "strconv"
  "runtime"
  l4g "github.com/alecthomas/log4go"
  "github.com/KenmyZhang/medicine-manual-spider/model"
)

var (
    b_cato    = regexp.MustCompile(`<a href="http://search.360kad.com/\?pageText=[\s\S]+?" target="_blank">[\s\S]+?</a>`)
    b_catoPrefix = regexp.MustCompile(`<a href="http://search.360kad.com/\?pageText=[\s\S]+?" target="_blank">`)
    b_catoSuffix = regexp.MustCompile(`</a>`)

    a_cato    = regexp.MustCompile(`<a href="http://www.360kad.com/Category_[0-9]+/Index.aspx" target="_blank">`)
    a_catoPrefix = regexp.MustCompile(`<a href="`)
    a_catoSuffix = regexp.MustCompile(`" target="_blank">`)

    a_maxPage    = regexp.MustCompile(`>[0-9]+?</a>[\s]+?<a class="Ynext"`) 
    a_maxPagePrefix = regexp.MustCompile(`>`)
    a_maxPageSuffix = regexp.MustCompile(`</a>[\s]+?<a class="Ynext"`)

    b_maxPage       = regexp.MustCompile(`<a class="Ylast" href="/Category_[0-9]+/Index.aspx\?page=[0-9]+">尾页`) 
    b_maxPagePrefix = regexp.MustCompile(`<a class="Ylast" href="/Category_[0-9]+/Index.aspx\?page=`)
    b_maxPageSuffix = regexp.MustCompile(`">尾页`)    

    a_productNum       = regexp.MustCompile(`<span id="ico_[0-9]+?"></span>`) 
    a_productNumPrefix = regexp.MustCompile(`<span id="ico_`)
    a_productNumSuffix = regexp.MustCompile(`"></span>`) 

    b_productNum       = regexp.MustCompile(`<a class="name" target="_blank" href="/product/[0-9]+?.shtml"`) 
    b_productNumPrefix = regexp.MustCompile(`<a class="name" target="_blank" href="/product/`)
    b_productNumSuffix = regexp.MustCompile(`.shtml"`) 

    kadProductName       = regexp.MustCompile(`<li title="[\s\S]+?">商品名称：[\s\S]+?</li>`) 
    kadProductNamePrefix = regexp.MustCompile(`<li title="[\s\S]+?">商品名称：`)
    kadProductNameSuffix = regexp.MustCompile(`</li>`)  

    kadApproveNum       = regexp.MustCompile(`<li title="[\S]+?">批准文号：[\S]+?</li>`) 
    kadApproveNumPrefix = regexp.MustCompile(`<li title="[\S]+?">批准文号：`)
    kadApproveNumSuffix = regexp.MustCompile(`</li>`)  

    kadManualfacturer       = regexp.MustCompile(`<li title="[\S\s]+?">生产企业：[\S\s]+?</li>`) 
    kadManualfacturerPrefix = regexp.MustCompile(`<li title="[\S\s]+?">生产企业：`)
    kadManualfacturerSuffix = regexp.MustCompile(`</li>`)     

    kadCurrentSize       = regexp.MustCompile(`<li title="[\s\S]+?">规格：[\s\S]+?</li>`) 
    kadCurrentSizePrefix = regexp.MustCompile(`<li title="[\s\S]+?">规格：`)
    kadCurrentSizeSuffix = regexp.MustCompile(`</li>`) 

    kadAllSize       = regexp.MustCompile(`<ul class="clearfix">[\s\S]+?<li class="dtl-inf-rur"><a href="javascript:;">[\s\S]+?</a></li>[\s]+?</ul>`)     

    kadOtherSize       = regexp.MustCompile(`<li><a href="/product/[0-9]+?.shtml">[\s\S]+?</a></li>`)
    kadOtherSizePrefix = regexp.MustCompile(`<li><a href="/product/`)
    kadOtherSizeSuffix = regexp.MustCompile(`.shtml">[\s\S]+?</a></li>`)     

/*
                                <ul class="clearfix">
                                                <li><a href="/product/40152.shtml">50ml</a></li>
                                                <li class="dtl-inf-rur"><a href="javascript:;">100ml</a></li>
                                                <li><a href="/product/129109.shtml">70ml/瓶</a></li>
                                                <li><a href="/product/129110.shtml">5ml</a></li>
                                </ul>    

                                <ul class="clearfix">
                                                <li class="dtl-inf-rur"><a href="javascript:;">50ml</a></li>
                                                <li><a href="/product/71613.shtml">100ml</a></li>
                                                <li><a href="/product/129109.shtml">70ml/瓶</a></li>
                                                <li><a href="/product/129110.shtml">5ml</a></li>
                                </ul>                                
*/
    kadPrice       = regexp.MustCompile(`salePrice2 : [0-9]+?,`) 
    kadPricePrefix = regexp.MustCompile(`salePrice2 : `)
    kadPriceSuffix = regexp.MustCompile(`,`) 
)

type UrlNameAndMaxPage struct {
  UrlName string
  MaxPage int
}

func GetAllCatoNames(a_urlNameChan, b_urlNameChan chan string){
	url := `http://www.360kad.com/dymhh/allclass.shtml`
  l4g.Debug(url + " begin")
  respBody, err := httpGet(url, false)	
  if err != nil {
   	l4g.Error(url + ", " + err.Error())
   	return
  }   
  l4g.Debug(url + " end") 

  catomatches :=  b_cato.FindAllString(respBody, -1)
  for _, catomatch := range catomatches {
      catomatch = b_catoPrefix.ReplaceAllString(catomatch, "")
      catomatch = b_catoSuffix.ReplaceAllString(catomatch, "")
      b_urlName := "http://search.360kad.com/?Pagetext=" + catomatch
      l4g.Debug("b",b_urlName)
      a_urlNameChan <- b_urlName
  }

  a_catomatches :=  a_cato.FindAllString(respBody, -1)
  for _, a_catomatch := range a_catomatches {
      a_catomatch = a_catoPrefix.ReplaceAllString(a_catomatch, "")
      a_catomatch = a_catoSuffix.ReplaceAllString(a_catomatch, "")
      l4g.Debug("a",a_catomatch)
      b_urlNameChan <- a_catomatch
  }
  return
}

func GetMaxPageOfPerCato_A(a_urlNameChan chan string, a_urlNameAndMaxPageChan chan *UrlNameAndMaxPage) {
  for {
    select {
      case urlName := <-a_urlNameChan:
        l4g.Debug(urlName + " begin")
        respBody, err := httpGet(urlName, false)  
        if err != nil {
          l4g.Error(urlName + ", " + err.Error())
          return
        }   
        l4g.Debug(urlName + " end")
        a_max := a_maxPage.FindString(respBody)
        a_max = a_maxPageSuffix.ReplaceAllString(a_max, "")
        a_max = a_maxPagePrefix.ReplaceAllString(a_max, "")
        a_urlNameAndMaxPage := &UrlNameAndMaxPage{}        
        a_urlNameAndMaxPage.UrlName = urlName
        max, _ :=  strconv.Atoi(a_max)
        if max == 0 {
           max = 1
        }
        a_urlNameAndMaxPage.MaxPage = max
        a_urlNameAndMaxPageChan <- a_urlNameAndMaxPage
        l4g.Debug("a_urlNameAndMaxPage",a_urlNameAndMaxPage)
        runtime.Gosched()

      case <-time.After(time.Minute * 2):
        l4g.Error("timeout")
        return
    }
  }
}

func GetMaxPageOfPerCato_B(b_urlNameChan chan string, b_urlNameAndMaxPageChan chan *UrlNameAndMaxPage) {
  for {
    select {
      case urlName := <-b_urlNameChan:
        l4g.Debug(urlName + " begin")
        respBody, err := httpGet(urlName, false)  
        if err != nil {
          l4g.Error(urlName + ", " + err.Error())
          return
        }   
        l4g.Debug(urlName + " end")
        b_max := b_maxPage.FindString(respBody)
        b_max = b_maxPageSuffix.ReplaceAllString(b_max, "")
        b_max = b_maxPagePrefix.ReplaceAllString(b_max, "")
        b_urlNameAndMaxPage := &UrlNameAndMaxPage{}
        b_urlNameAndMaxPage.UrlName = urlName
        l4g.Debug("b_max",b_max)
        max, _ :=  strconv.Atoi(b_max)
        if max == 0 {
           max = 1
        }
        b_urlNameAndMaxPage.MaxPage = max
        b_urlNameAndMaxPageChan <- b_urlNameAndMaxPage
        l4g.Debug("b_urlNameAndMaxPage",b_urlNameAndMaxPage)
        runtime.Gosched()

      case <-time.After(time.Minute * 2):
        l4g.Error("timeout")
        return
    }
  }
}

func GetAllPageOfPerCato_A(a_urlNameAndMaxPageChan chan *UrlNameAndMaxPage, numChan_A chan string) {
  for {
    select {
      case a_urlNameAndMaxPage := <-a_urlNameAndMaxPageChan:
        for i := 1; i < a_urlNameAndMaxPage.MaxPage; i++ {
          url := a_urlNameAndMaxPage.UrlName + "&pageIndex=" + strconv.Itoa(i)         
          GetPerpageOfNum_A(url, numChan_A)
        }
        runtime.Gosched()
      case <-time.After(time.Minute * 2):
        l4g.Error("timeout")
        return
    }
  }
}

func GetAllPageOfPerCato_B(b_urlNameAndMaxPageChan chan *UrlNameAndMaxPage, numChan_B chan string) {
  for {
    select {
      case b_urlNameAndMaxPage := <-b_urlNameAndMaxPageChan:
        for i := 1; i < b_urlNameAndMaxPage.MaxPage; i++ {
          url := b_urlNameAndMaxPage.UrlName + "?page=" + strconv.Itoa(i)
          GetPerpageOfNum_B(url, numChan_B)
        }
        runtime.Gosched()
      case <-time.After(time.Minute * 2):
        l4g.Error("timeout")
        return
    }
  }
}

func GetPerpageOfNum_A(url_A string, numChan_A chan string) {
    respBody, err := httpGet(url_A, false)  
    if err != nil {
      l4g.Error(url_A + ", " + err.Error())
      return
    }
    a_productNums := a_productNum.FindAllString(respBody, -1)
    for _, num := range a_productNums {
      num = a_productNumPrefix.ReplaceAllString(num, "")
      num = a_productNumSuffix.ReplaceAllString(num, "")
      numChan_A <- num
    }   
}

func GetPerpageOfNum_B(url_B string, numChan_B chan string) {
    respBody, err := httpGet(url_B, false)  
    if err != nil {
      l4g.Debug(url_B + ", " + err.Error())
      return
    }  
    b_productNums := b_productNum.FindAllString(respBody, -1)
    for _, num := range b_productNums {
      num = b_productNumPrefix.ReplaceAllString(num, "")
      num = b_productNumSuffix.ReplaceAllString(num, "")
      numChan_B <- num
    } 
}

func GetPerProductByNum(numChan chan string) {
  for {
    select {
      case num, ok :=<- numChan:
        if !ok {
          l4g.Error("Finsh A")
          return
        }
        url := "http://www.360kad.com/product/" + "40152" + ".shtml"
        l4g.Debug(url + "  begin") 
        respBody, err := httpGet(url, false)  
        if err != nil {
           l4g.Error(url + ", " + err.Error())
        } 
        l4g.Debug(url + "  end")  
       
        productSizeAndPrize := &model.ProductSizeAndPrize{}

        productSizeAndPrize.Num = num

        a_productNameStr := kadProductName.FindString(respBody)
        a_productNameStr = kadProductNamePrefix.ReplaceAllString(a_productNameStr, "")
        a_productNameStr = kadProductNameSuffix.ReplaceAllString(a_productNameStr, "")
        l4g.Debug("药品名称",a_productNameStr)   
        productSizeAndPrize.Name = a_productNameStr

        a_approveNumStr := kadApproveNum.FindString(respBody)
        a_approveNumStr = kadApproveNumPrefix.ReplaceAllString(a_approveNumStr, "")
        a_approveNumStr = kadApproveNumSuffix.ReplaceAllString(a_approveNumStr, "")
        l4g.Debug("批准文号：",a_approveNumStr)
        productSizeAndPrize.ApprovalNumber = a_approveNumStr

        a_sizeStr := kadCurrentSize.FindString(respBody)
        a_sizeStr = kadCurrentSizePrefix.ReplaceAllString(a_sizeStr, "")
        a_sizeStr = kadCurrentSizeSuffix.ReplaceAllString(a_sizeStr, "")
        l4g.Debug("规格：", a_sizeStr)
        productSizeAndPrize.CurrentSize = a_sizeStr        

        a_manufacturerStr := kadManualfacturer.FindString(respBody)
        a_manufacturerStr  = kadManualfacturerPrefix.ReplaceAllString(a_manufacturerStr, "")
        a_manufacturerStr  = kadManualfacturerSuffix.ReplaceAllString(a_manufacturerStr, "")
        l4g.Debug("生产厂商：", a_manufacturerStr)
        productSizeAndPrize.Manufacturer = a_manufacturerStr 

        a_priceStr := kadPrice.FindString(respBody)
        a_priceStr = kadPricePrefix.ReplaceAllString(a_priceStr, "")
        a_priceStr = kadPriceSuffix.ReplaceAllString(a_priceStr, "")
        l4g.Debug("价格：", a_priceStr)
        productSizeAndPrize.Price = a_priceStr
        l4g.Debug(url)
        
        a_kadAllSize := kadAllSize.FindString(respBody)
        otherSizeNums := kadOtherSize.FindAllString(a_kadAllSize, -1)
        for _, num := range otherSizeNums {
          num = kadOtherSizePrefix.ReplaceAllString(num, "")
          num = kadOtherSizeSuffix.ReplaceAllString(num, "")
          l4g.Debug("num",num)
          numChan <- num
        } 

        if productSizeAndPrize.Name != "" {
          SaveProductSizeAndPrize(productSizeAndPrize)
        }
      case <-time.After(time.Minute * 1):
        l4g.Error("timeout")
        return
    }
  }
}

func SpyProductPriceFrom360kad() {
  l4g.Debug("begin")
  a_urlNameChan := make(chan string, 1000)
  b_urlNameChan := make(chan string, 1000)
  a_urlNameAndMaxPageChan := make(chan *UrlNameAndMaxPage, 2000)
  b_urlNameAndMaxPageChan := make(chan *UrlNameAndMaxPage, 2000)
  a_numChan := make(chan string, 10000)
  b_numChan := make(chan string, 10000)
  go GetAllCatoNames(a_urlNameChan, b_urlNameChan)
  go GetMaxPageOfPerCato_A(a_urlNameChan, a_urlNameAndMaxPageChan)
  go GetMaxPageOfPerCato_B(b_urlNameChan, b_urlNameAndMaxPageChan) 
  go GetAllPageOfPerCato_A(a_urlNameAndMaxPageChan, a_numChan)
  go GetAllPageOfPerCato_B(b_urlNameAndMaxPageChan, b_numChan)
  go GetPerProductByNum(a_numChan)
  go GetPerProductByNum(b_numChan)
}