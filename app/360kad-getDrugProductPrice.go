package app

import (
	"regexp"
	"fmt"
  "time"
  "strconv"
  //"github.com/KenmyZhang/medicine-manual-spider/model"
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
)

type UrlNameAndMaxPage struct {
  UrlName string
  MaxPage int
}

func GetAllCatoNames(a_urlNameChan, b_urlNameChan chan string){
	url := `http://www.360kad.com/dymhh/allclass.shtml`
  
  fmt.Println("GetAllCatoNames:" + url + " begin")
  respBody, err := httpGet(url, false)	
  if err != nil {
   	fmt.Println("ERROR GetAllCatoNames:" + url + ", " + err.Error())
   	return
  }   
  fmt.Println("GetAllCatoNames:" +url + " end") 

  catomatches :=  b_cato.FindAllString(respBody, -1)
  for _, catomatch := range catomatches {
      catomatch = b_catoPrefix.ReplaceAllString(catomatch, "")
      catomatch = b_catoSuffix.ReplaceAllString(catomatch, "")
      b_urlName := "http://search.360kad.com/?Pagetext=" + catomatch
      fmt.Println("b",b_urlName)
      a_urlNameChan <- b_urlName
  }

  a_catomatches :=  a_cato.FindAllString(respBody, -1)
  for _, a_catomatch := range a_catomatches {
      a_catomatch = a_catoPrefix.ReplaceAllString(a_catomatch, "")
      a_catomatch = a_catoSuffix.ReplaceAllString(a_catomatch, "")
      fmt.Println("a",a_catomatch)
      b_urlNameChan <- a_catomatch
  }

  return
}

func GetMaxPageOfPerCato_A(a_urlNameChan chan string, a_urlNameAndMaxPageChan chan *UrlNameAndMaxPage) {
  for {
    select {
      case urlName := <-a_urlNameChan:
        fmt.Println("GetMaxPageOfPerCato_A:" + urlName + " begin")
        respBody, err := httpGet(urlName, false)  
        if err != nil {
          fmt.Println("ERROR GetMaxPageOfPerCato_A:" + urlName + ", " + err.Error())
          return
        }   
        fmt.Println("GetMaxPageOfPerCato_A:" +urlName + " end")
        a_max := a_maxPage.FindString(respBody)
        a_max = a_maxPageSuffix.ReplaceAllString(a_max, "")
        a_max = a_maxPagePrefix.ReplaceAllString(a_max, "")
        a_urlNameAndMaxPage := &UrlNameAndMaxPage{}        
        a_urlNameAndMaxPage.UrlName = urlName
        fmt.Println("a_max",a_max)
        max, _ :=  strconv.Atoi(a_max)
        if max == 0 {
           max = 1
        }
        a_urlNameAndMaxPage.MaxPage = max
        a_urlNameAndMaxPageChan <- a_urlNameAndMaxPage
        fmt.Println("a_urlNameAndMaxPage",a_urlNameAndMaxPage)

      case <-time.After(time.Minute * 2):
        fmt.Println("ERROR GetMaxPageOfPerCato_A")
        return
    }
  }
}

func GetMaxPageOfPerCato_B(b_urlNameChan chan string, b_urlNameAndMaxPageChan chan *UrlNameAndMaxPage) {
  for {
    select {
      case urlName := <-b_urlNameChan:
        fmt.Println("GetMaxPageOfPerCato_B:" + urlName + " begin")
        respBody, err := httpGet(urlName, false)  
        if err != nil {
          fmt.Println("ERROR GetMaxPageOfPerCato_B:" + urlName + ", " + err.Error())
          return
        }   
        fmt.Println("GetMaxPageOfPerCato_B:" +urlName + " end")
        b_max := b_maxPage.FindString(respBody)
        b_max = b_maxPageSuffix.ReplaceAllString(b_max, "")
        b_max = b_maxPagePrefix.ReplaceAllString(b_max, "")
        b_urlNameAndMaxPage := &UrlNameAndMaxPage{}
        b_urlNameAndMaxPage.UrlName = urlName
        fmt.Println("b_max",b_max)
        max, _ :=  strconv.Atoi(b_max)
        if max == 0 {
           max = 1
        }
        b_urlNameAndMaxPage.MaxPage = max
        b_urlNameAndMaxPageChan <- b_urlNameAndMaxPage
        fmt.Println("b_urlNameAndMaxPage",b_urlNameAndMaxPage)

      case <-time.After(time.Minute * 2):
        fmt.Println("ERROR GetMaxPageOfPerCato_B")
        return
    }
  }
}

func GetAllPageOfPerCato_A(a_urlNameAndMaxPageChan chan *UrlNameAndMaxPage) {
  for {
    select {
      case a_urlNameAndMaxPage := <-a_urlNameAndMaxPageChan:
        for i := 1; i < a_urlNameAndMaxPage.MaxPage; i++ {
          url := a_urlNameAndMaxPage.UrlName + "&pageIndex=" + strconv.Itoa(i)
          respBody, err := httpGet(url, false)  
          if err != nil {
            fmt.Println("ERROR GetMaxPageOfPerCato_B:" + url + ", " + err.Error())
            return
          } 
          fmt.Println("respBody:",respBody) 
        }
      case <-time.After(time.Minute * 2):
        fmt.Println("ERROR GetAllPageOfPerCato_A")
        return
    }
  }
}

func GetAllPageOfPerCato_B(b_urlNameAndMaxPageChan chan *UrlNameAndMaxPage) {
  for {
    select {
      case b_urlNameAndMaxPage := <-b_urlNameAndMaxPageChan:
        for i := 1; i < b_urlNameAndMaxPage.MaxPage; i++ {
          url := b_urlNameAndMaxPage.UrlName + "?page=" + strconv.Itoa(i)
          respBody, err := httpGet(url, false)  
          if err != nil {
            fmt.Println("ERROR GetAllPageOfPerCato_B:" + url + ", " + err.Error())
            return
          }
          fmt.Println("respBody:",respBody)       
        }
      case <-time.After(time.Minute * 2):
        fmt.Println("ERROR GetAllPageOfPerCato_B")
        return
    }
  }
}

func SpyProductPriceFrom360kad() {
  fmt.Println("begin")
  a_urlNameChan := make(chan string, 1000)
  b_urlNameChan := make(chan string, 1000)
  a_urlNameAndMaxPageChan := make(chan *UrlNameAndMaxPage, 2000)
  b_urlNameAndMaxPageChan := make(chan *UrlNameAndMaxPage, 2000)
  go GetAllCatoNames(a_urlNameChan, b_urlNameChan)
  go GetMaxPageOfPerCato_A(a_urlNameChan, a_urlNameAndMaxPageChan)
  go GetMaxPageOfPerCato_B(b_urlNameChan, b_urlNameAndMaxPageChan) 
  go GetAllPageOfPerCato_A(a_urlNameAndMaxPageChan)
  go GetAllPageOfPerCato_B(b_urlNameAndMaxPageChan)
}