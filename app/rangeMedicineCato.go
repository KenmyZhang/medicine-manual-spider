package app

import (
	"strconv"
	"regexp"
	"fmt"
  "time"
  "os"
  "runtime"
  "github.com/KenmyZhang/medicine-manual-spider/model"
)

var MedicineCato = []string{"x-GanMao/","x-KeChuan/","x-GanDan/","x-ChangWei/",
	"x-FuKe/","x-NanKe/","x-GuKe/","x-ZhiTong/",
	"x-WaiShang/","x-ZhongLiu/","x-PiFuBing/","x-WuGuanKe/",
	"x-XiaoYanYao/","x-MiNiaoXi/","x-JingShenXi/","x-BuYiLei/",
	"x-GaoXueYa/","x-TangNiaoBing/","x-XueYeXi/",
	"x-MianYiLi/","x-KangGuoMin/","x-WeiFenLei/",
	"z-GanMao/","z-KeChuan/","z-QingRe/","z-AnShen/",
	"z-FuKe/","z-NanKe/", "z-GanDan/","z-ChangWei/",
	"z-XueYeXi/","z-ZhongLiu/","z-KouQiang/",
	"z-YanKe/","z-GuKe/","z-WaiShang/","z-TouTong/",
	"z-MiNiao/","z-TangNiaoBing/","z-PiFuBing/",
	"z-ErBiHou/","z-BuYiLei/","z-WeiFenLei/"}

var (
    a_totalPage       = regexp.MustCompile(`var pager=new iwmsPager[\s\S]+?true`)
    a_totalPagePrefix = regexp.MustCompile(`var pager=new iwmsPager[\s\S]+?,`)
    a_totalPageSuffix = regexp.MustCompile(`,true`)

    a_num       = regexp.MustCompile(`<div class="title text-oneline"><a href="../[a-zA-Z-]+?/[0-9]+.htm" target="_blank">`)
    a_numPrefix = regexp.MustCompile(`<div class="title text-oneline"><a href="../[a-zA-Z-]+?/`)
    a_numSuffix = regexp.MustCompile(`.htm" target="_blank">`)    

    a_productName       = regexp.MustCompile(`<p>药品名称：<span itemprop="name"><u>[\s\S]+?</u>`)
    a_productNamePrefix = regexp.MustCompile(`<p>药品名称：<span itemprop="name"><u>`)
    a_productNameSuffix = regexp.MustCompile(`</u>`)    

    a_approveNum       = regexp.MustCompile(`<p>批准文号：<u><a href="[\S]+?">[\s\S]+?</a>`)
    a_approveNumPrefix = regexp.MustCompile(`<p>批准文号：<u><a href="[\S]+?">`)
    a_approveNumSuffix = regexp.MustCompile(`</a>`)

    a_size       = regexp.MustCompile(`<p>规格：<u>[\s\S]+?</u>`) 
    a_sizePrefix = regexp.MustCompile(`<p>规格：<u>`)
    a_sizeSuffix = regexp.MustCompile(`</u>`)

    a_manufacturer = regexp.MustCompile(`<p>生产厂家：<u>[\s\S]+?</u></p>`)
    a_manufacturerPrefix = regexp.MustCompile(`<p>生产厂家：<u>`)
    a_manufacturerSuffix = regexp.MustCompile(`</u></p>`)

    a_price = regexp.MustCompile(`零售价格：[0-9.]+元`)
    a_pricePrefix = regexp.MustCompile(`零售价格：`)
    a_priceSuffix = regexp.MustCompile(`元`)

)

type MedicineUrlNameAndPage struct {
	MedicineUrlName string
	MaxPage int
}

type MedicineUrlNameAndNum struct {
  MedicineUrlName string
  Num string
}

func RangeMedicineCato(arr []string, catoChan chan string) {
	for _, val := range arr {
    fmt.Println("cato-",val)
		catoChan <- val
	}
	close(catoChan)
  fmt.Println("close catoChan")
}	

func GetMaxPageOfMedicine(catoChan chan string, medicineUrlNameAndPageChan chan MedicineUrlNameAndPage) {
  for {
    select {
      case catoName, ok := <-catoChan: 
     		if !ok {
          close(medicineUrlNameAndPageChan)
     			fmt.Println("end of medicineUrlNameAndPageChan")
     			return
     		} else {
     			url := `https://www.315jiage.cn/` +catoName
          fmt.Println("GetMaxPageOfMedicine:" + url + " begin")                
  				respBody, err := httpGet(url, false)
  				if err != nil {
      			fmt.Println("ERROR GetMaxPageOfMedicine:" + url + ", " + err.Error())
      			return
  				}
          fmt.Println("GetMaxPageOfMedicine: " + url + " end")
			    a_totalPageStr := a_totalPage.FindString(respBody)
 					a_totalPageStr = a_totalPagePrefix.ReplaceAllString(a_totalPageStr,"")
 					a_totalPageStr = a_totalPageSuffix.ReplaceAllString(a_totalPageStr,"")
 					maxPage, _ := strconv.Atoi(a_totalPageStr)
 					medicineUrlNameAndPage := MedicineUrlNameAndPage{MedicineUrlName:url, MaxPage:maxPage}
          fmt.Println("GetMaxPageOfMedicine:medicineUrlNameAndPage-", medicineUrlNameAndPage)
          runtime.Gosched()
 					medicineUrlNameAndPageChan <- medicineUrlNameAndPage
 			}
      case <-time.After(time.Minute * 5):
        fmt.Println("ERROR GetMaxPageOfMedicine timeout")
        return
  	}  
  }
}

func GetAllMedicineNumFromOneCato(medicineUrlNameAndPageChan chan MedicineUrlNameAndPage, medicineUrlNameAndNumChan chan *MedicineUrlNameAndNum) {
	for {
		select {
			case medicineUrlNameAndPage, ok := <-medicineUrlNameAndPageChan:
        if !ok {
          fmt.Println("end of medicineUrlNameAndNumChan")
          close(medicineUrlNameAndNumChan)
          return
        }

        f, err := os.OpenFile("./all_medicine_num.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm|os.ModeTemporary)
        if err != nil {
          fmt.Println("ERROR GetAllMedicineNumFromOneCato:" + err.Error())
          return
        }
        defer f.Close()

				for i := 1; i <= medicineUrlNameAndPage.MaxPage; i++{
					nums := GetAllMedicineNumFromOnePage(medicineUrlNameAndPage.MedicineUrlName, i, medicineUrlNameAndNumChan)
          for _, num := range nums {
            _, _ = f.WriteString(num +"\n")   
          }         
				}
        f.Sync() 
			case  <-time.After(time.Minute * 5):
        fmt.Println("ERROR GetAllMedicineNumFromOneCato timeout")
				return
		}
	}
}

func GetAllMedicineNumFromOnePage(urlName string, index int, medicineUrlNameAndNumChan chan *MedicineUrlNameAndNum) []string {
  var url string
	if index == 1 {
		url = urlName
	} else {
		url = urlName + `defaultp` + strconv.Itoa(index) + `.htm`
	}
  
  fmt.Println("GetAllMedicineNumFromOnePage:" + url + " begin")
  respBody, err := httpGet(url, false)	
  if err != nil {
   	fmt.Println("ERROR GetAllMedicineNumFromOnePage:" + url + ", " + err.Error())
   	return nil
  }   
  fmt.Println("GetAllMedicineNumFromOnePage:" +url + " end") 

  drugNumMatches :=  a_num.FindAllString(respBody, -1)
  var nums []string
  for _, drugNum := range drugNumMatches {
    drugNum = a_numPrefix.ReplaceAllString(drugNum, "")
    drugNum = a_numSuffix.ReplaceAllString(drugNum, "")
    fmt.Println("GetAllMedicineNumFromOnePage:drugNum-" + drugNum) 
    medicineUrlNameAndNum := &MedicineUrlNameAndNum{}
    medicineUrlNameAndNum.MedicineUrlName = urlName
    medicineUrlNameAndNum.Num = drugNum
    nums = append(nums, drugNum)
    medicineUrlNameAndNumChan <- medicineUrlNameAndNum
    runtime.Gosched()
  }
  return nums
}

func GetOneMedcine(medicineUrlNameAndNumChan chan *MedicineUrlNameAndNum) {
  for {
    select {
      case medicineUrlNameAndNum, ok :=<- medicineUrlNameAndNumChan:
        if !ok {
          fmt.Println("end of medicineUrlNameAndNumChan, Finsh ALL")
          return
        }
  	    url := medicineUrlNameAndNum.MedicineUrlName + medicineUrlNameAndNum.Num + `.htm`
        fmt.Println("GetOneMedcine:" + url + "  end") 
        respBody, err := httpGet(url, false)	
        if err != nil {
      	   fmt.Println("ERROR SpyMedicineManual:" + url + ", " + err.Error())
      	   return
        }	
        fmt.Println("GetOneMedcine:" + url + " begin") 
       
        productSizeAndPrize := &model.ProductSizeAndPrize{}

        productSizeAndPrize.Num = medicineUrlNameAndNum.Num

       	a_productNameStr := a_productName.FindString(respBody)
       	a_productNameStr = a_productNamePrefix.ReplaceAllString(a_productNameStr, "")
       	a_productNameStr = a_productNameSuffix.ReplaceAllString(a_productNameStr, "")
       	fmt.Println("药品名称",a_productNameStr)   
        productSizeAndPrize.Name = a_productNameStr

    	  a_approveNumStr := a_approveNum.FindString(respBody)
    	  a_approveNumStr = a_approveNumPrefix.ReplaceAllString(a_approveNumStr, "")
    	  a_approveNumStr = a_approveNumSuffix.ReplaceAllString(a_approveNumStr, "")
    	  fmt.Println("批准文号：",a_approveNumStr)
        productSizeAndPrize.ApprovalNumber = a_approveNumStr

    	  a_sizeStr := a_size.FindString(respBody)
    	  a_sizeStr = a_sizePrefix.ReplaceAllString(a_sizeStr, "")
    	  a_sizeStr = a_sizeSuffix.ReplaceAllString(a_sizeStr, "")
    	  fmt.Println("规格：", a_sizeStr)
        productSizeAndPrize.CurrentSize = a_sizeStr        

    	  a_manufacturerStr := a_manufacturer.FindString(respBody)
    	  a_manufacturerStr = a_manufacturerPrefix.ReplaceAllString(a_manufacturerStr, "")
    	  a_manufacturerStr = a_manufacturerSuffix.ReplaceAllString(a_manufacturerStr, "")
    	  fmt.Println("生产厂商：", a_manufacturerStr)
        productSizeAndPrize.Manufacturer = a_manufacturerStr 

   	    a_priceStr := a_price.FindString(respBody)
    	  a_priceStr = a_pricePrefix.ReplaceAllString(a_priceStr, "")
    	  a_priceStr = a_priceSuffix.ReplaceAllString(a_priceStr, "")
    	  fmt.Println("价格：", a_priceStr)
        productSizeAndPrize.Price = a_priceStr

        if productSizeAndPrize.Name != "" {
          SaveProductSizeAndPrize(productSizeAndPrize)
        }
      case <-time.After(time.Minute * 5):
        fmt.Println("ERROR GetOneMedcine timeout")
        return
    }
  }
}

func SpyMedicineProductPriceFromJiaGe() {
	catoChan := make(chan string, 100)
	medicineUrlNameAndPageChan := make(chan MedicineUrlNameAndPage, 10000)
	medicineUrlNameAndNumChan  := make(chan *MedicineUrlNameAndNum, 60000)
	go RangeMedicineCato(MedicineCato, catoChan)
	go GetMaxPageOfMedicine(catoChan, medicineUrlNameAndPageChan)
	go GetAllMedicineNumFromOneCato(medicineUrlNameAndPageChan, medicineUrlNameAndNumChan)
	go GetOneMedcine(medicineUrlNameAndNumChan)
}


