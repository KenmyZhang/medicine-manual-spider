package app

import (
	"github.com/KenmyZhang/medicine-manual-spider/model"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"time"
  	l4g "github.com/alecthomas/log4go"
)

var MedicineCato = []string{"x-GanMao/", "x-KeChuan/", "x-GanDan/", "x-ChangWei/",
	"x-FuKe/", "x-NanKe/", "x-GuKe/", "x-ZhiTong/",
	"x-WaiShang/", "x-ZhongLiu/", "x-PiFuBing/", "x-WuGuanKe/",
	"x-XiaoYanYao/", "x-MiNiaoXi/", "x-JingShenXi/", "x-BuYiLei/",
	"x-GaoXueYa/", "x-TangNiaoBing/", "x-XueYeXi/",
	"x-MianYiLi/", "x-KangGuoMin/", "x-WeiFenLei/",
	"z-GanMao/", "z-KeChuan/", "z-QingRe/", "z-AnShen/",
	"z-FuKe/", "z-NanKe/", "z-GanDan/", "z-ChangWei/",
	"z-XueYeXi/", "z-ZhongLiu/", "z-KouQiang/",
	"z-YanKe/", "z-GuKe/", "z-WaiShang/", "z-TouTong/",
	"z-MiNiao/", "z-TangNiaoBing/", "z-PiFuBing/",
	"z-ErBiHou/", "z-BuYiLei/", "z-WeiFenLei/"}

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

	a_manufacturer       = regexp.MustCompile(`<p>生产厂家：<u>[\s\S]+?</u></p>`)
	a_manufacturerPrefix = regexp.MustCompile(`<p>生产厂家：<u>`)
	a_manufacturerSuffix = regexp.MustCompile(`</u></p>`)

	a_price       = regexp.MustCompile(`零售价格：[0-9.]+元`)
	a_pricePrefix = regexp.MustCompile(`零售价格：`)
	a_priceSuffix = regexp.MustCompile(`元`)
)

type MedicineUrlNameAndPage struct {
	MedicineUrlName string
	MaxPage         int
}

type MedicineUrlNameAndNum struct {
	MedicineUrlName string
	Num             string
}

func RangeMedicineCato(arr []string, catoChan chan string) {
	for _, val := range arr {
		l4g.Debug("cato-", val)
		catoChan <- val
	}
	close(catoChan)
	l4g.Info("close catoChan")
}

func GetMaxPageOfMedicine(catoChan chan string, medicineUrlNameAndPageChan chan MedicineUrlNameAndPage) {
	for {
		select {
		case catoName, ok := <-catoChan:
			if !ok {
				close(medicineUrlNameAndPageChan)
				l4g.Info("catoChan already close")
				return
			} else {
				url := `https://www.315jiage.cn/` + catoName
				l4g.Debug(url + " begin")
				respBody, err := httpGet(url, false)
				if err != nil {
					l4g.Error(url + ", " + err.Error())
					continue
				}
				l4g.Debug(url + " end")
				a_totalPageStr := a_totalPage.FindString(respBody)
				a_totalPageStr = a_totalPagePrefix.ReplaceAllString(a_totalPageStr, "")
				a_totalPageStr = a_totalPageSuffix.ReplaceAllString(a_totalPageStr, "")
				maxPage, _ := strconv.Atoi(a_totalPageStr)
				if maxPage == 0 {
					maxPage = 1
				}
				medicineUrlNameAndPage := MedicineUrlNameAndPage{MedicineUrlName: url, MaxPage: maxPage}
				l4g.Debug(medicineUrlNameAndPage)
				runtime.Gosched()
				medicineUrlNameAndPageChan <- medicineUrlNameAndPage
			}
		case <-time.After(time.Minute * 5):
			l4g.Error("timeout")
			return
		}
	}
}

func GetAllMedicineNumFromOneCato(medicineUrlNameAndPageChan chan MedicineUrlNameAndPage, 
  medicineUrlNameAndNumChan chan *MedicineUrlNameAndNum, f *os.File) {
	for {
		select {
		case medicineUrlNameAndPage, ok := <-medicineUrlNameAndPageChan:
			if !ok {
				l4g.Info("medicineUrlNameAndPageChan already close")
				close(medicineUrlNameAndNumChan)
				return
			}

			for i := 1; i <= medicineUrlNameAndPage.MaxPage; i++ {
				nums := GetAllMedicineNumFromOnePage(medicineUrlNameAndPage.MedicineUrlName, i, medicineUrlNameAndNumChan)
				for _, num := range nums {
					_, _ = f.WriteString(medicineUrlNameAndPage.MedicineUrlName + num + ".htm" + "\n")
				}
			}
			f.Sync()
      		runtime.Gosched()
		case <-time.After(time.Minute * 1):
			l4g.Error("timeout")
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

	l4g.Debug(url + " begin")
	respBody, err := httpGet(url, false)
	if err != nil {
		l4g.Error(url + ", " + err.Error())
		return nil
	}
	l4g.Debug(url + " end")

	drugNumMatches := a_num.FindAllString(respBody, -1)
	var nums []string
	for _, drugNum := range drugNumMatches {
		drugNum = a_numPrefix.ReplaceAllString(drugNum, "")
		drugNum = a_numSuffix.ReplaceAllString(drugNum, "")
		l4g.Debug(drugNum)
		medicineUrlNameAndNum := &MedicineUrlNameAndNum{}
		medicineUrlNameAndNum.MedicineUrlName = urlName
		medicineUrlNameAndNum.Num = drugNum
		nums = append(nums, drugNum)
		medicineUrlNameAndNumChan <- medicineUrlNameAndNum
	}
	return nums
}

func GetOneMedcine(medicineUrlNameAndNumChan chan *MedicineUrlNameAndNum) {
	for {
		select {
		case medicineUrlNameAndNum, ok := <-medicineUrlNameAndNumChan:
			if !ok {
				l4g.Info("medicineUrlNameAndNumChan already close")
				return
			}
			url := medicineUrlNameAndNum.MedicineUrlName + medicineUrlNameAndNum.Num + `.htm`
			l4g.Debug(url + "  end")
			respBody, err := httpGet(url, false)
			if err != nil {
				l4g.Error(url + ", " + err.Error())
        continue
			}
			l4g.Debug(url + " begin")

			productSizeAndPrize := &model.ProductSizeAndPrize{}

			productSizeAndPrize.Num = medicineUrlNameAndNum.Num

			a_productNameStr := a_productName.FindString(respBody)
			a_productNameStr = a_productNamePrefix.ReplaceAllString(a_productNameStr, "")
			a_productNameStr = a_productNameSuffix.ReplaceAllString(a_productNameStr, "")
			l4g.Debug("药品名称", a_productNameStr)
			productSizeAndPrize.Name = a_productNameStr

			a_approveNumStr := a_approveNum.FindString(respBody)
			a_approveNumStr = a_approveNumPrefix.ReplaceAllString(a_approveNumStr, "")
			a_approveNumStr = a_approveNumSuffix.ReplaceAllString(a_approveNumStr, "")
			l4g.Debug("批准文号：", a_approveNumStr)
			productSizeAndPrize.ApprovalNumber = a_approveNumStr

			a_sizeStr := a_size.FindString(respBody)
			a_sizeStr = a_sizePrefix.ReplaceAllString(a_sizeStr, "")
			a_sizeStr = a_sizeSuffix.ReplaceAllString(a_sizeStr, "")
			l4g.Debug("规格：", a_sizeStr)
			productSizeAndPrize.CurrentSize = a_sizeStr

			a_manufacturerStr := a_manufacturer.FindString(respBody)
			a_manufacturerStr = a_manufacturerPrefix.ReplaceAllString(a_manufacturerStr, "")
			a_manufacturerStr = a_manufacturerSuffix.ReplaceAllString(a_manufacturerStr, "")
			l4g.Debug("生产厂商：", a_manufacturerStr)
			productSizeAndPrize.Manufacturer = a_manufacturerStr

			a_priceStr := a_price.FindString(respBody)
			a_priceStr = a_pricePrefix.ReplaceAllString(a_priceStr, "")
			a_priceStr = a_priceSuffix.ReplaceAllString(a_priceStr, "")
			l4g.Debug("价格：", a_priceStr)
			productSizeAndPrize.Price = a_priceStr

			if productSizeAndPrize.Name != "" {
				SaveProductSizeAndPrize(productSizeAndPrize)
			}
      		runtime.Gosched()
		case <-time.After(time.Minute * 1):
			l4g.Error("timeout")
			return
		}
	}
}

func SpyMedicineProductPriceFromJiaGe(f *os.File) {
	catoChan := make(chan string, 100)
	medicineUrlNameAndPageChan := make(chan MedicineUrlNameAndPage, 10000)
	medicineUrlNameAndNumChan := make(chan *MedicineUrlNameAndNum, 60000)
	go RangeMedicineCato(MedicineCato, catoChan)
	go GetMaxPageOfMedicine(catoChan, medicineUrlNameAndPageChan)
	go GetAllMedicineNumFromOneCato(medicineUrlNameAndPageChan, medicineUrlNameAndNumChan, f)
	go GetOneMedcine(medicineUrlNameAndNumChan)
}
