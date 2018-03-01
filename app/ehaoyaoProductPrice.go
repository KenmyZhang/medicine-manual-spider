package app

import (
	l4g "github.com/alecthomas/log4go"
	"github.com/KenmyZhang/medicine-manual-spider/model"
	"regexp"
	"strconv"
	"runtime"
	"os"
)

var eHaoYaoCatoArr = []string{
	`c367-s66.html`,
	`c367-s67.html`,
	`c367-s68.html`,
	`c367-s69.html`,
	`c367-s70.html`,
	`c367-s71.html`,
	`c367-s72.html`,
	`c367-s73.html`,
	`c367-s74.html`,
	`c367-s75.html`,
	`c367-s76.html`,
	`c367-s77.html`,
	`c367-s78.html`,
	`c367-s79.html`,
	`c367-s80.html`,
	`c367-s81.html`,
	`c367-s82.html`,
}

var (
	eHaoYaoTotalPage       = regexp.MustCompile(`<i class="countpage">[0-9]+?</i>`)
	eHaoYaoTotalPagePrefix = regexp.MustCompile(`<i class="countpage">`)
	eHaoYaoTotalPageSuffix = regexp.MustCompile(`</i>`)

	eHaoYaoNum       = regexp.MustCompile(`<a target="_blank" href="/product-[0-9]+?.html">[\s]+?<div class="name" title="`)
	eHaoYaoNumPrefix = regexp.MustCompile(`<a target="_blank" href="/`)
	eHaoYaoNumSuffix = regexp.MustCompile(`.html">[\s]+?<div class="name" title="`)	

	eHaoYaoProductName       = regexp.MustCompile(`<tr><td class="td_key">产品名称</td><td class="td_val">[\s\S]+?</td></tr>`)
	eHaoYaoProductNamePrefix = regexp.MustCompile(`<tr><td class="td_key">产品名称</td><td class="td_val">`)
	eHaoYaoProductNameSuffix = regexp.MustCompile(`</td></tr>`)

	eHaoYaoApproveNum       = regexp.MustCompile(`<tr class="bgColor"><td class="td_key">批准文号</td><td class="td_val">[\s\S]+?</td></tr>`)
	eHaoYaoApproveNumPrefix = regexp.MustCompile(`<tr class="bgColor"><td class="td_key">批准文号</td><td class="td_val">`)
	eHaoYaoApproveNumSuffix = regexp.MustCompile(`</td></tr>`)

	eHaoYaoManualfacturer       = regexp.MustCompile(`<tr class="bgColor"><td class="td_key">生产厂家</td><td class="td_val">[\s\S]+?</td></tr>`)
	eHaoYaoManualfacturerPrefix = regexp.MustCompile(`<tr class="bgColor"><td class="td_key">生产厂家</td><td class="td_val">`)
	eHaoYaoManualfacturerSuffix = regexp.MustCompile(`</td></tr>`)

	eHaoYaoCurrentSize       = regexp.MustCompile(`<tr class="bgColor"><td class="td_key">规格</td><td class="td_val">[\s\S]+?</td></tr>`)
	eHaoYaoCurrentSizePrefix = regexp.MustCompile(`<tr class="bgColor"><td class="td_key">规格</td><td class="td_val">`)
	eHaoYaoCurrentSizeSuffix = regexp.MustCompile(`</td></tr>`)		

	eHaoYaoPrice       = regexp.MustCompile(`<span class="price lFloat"><i>¥</i><em>[\s\S]+?</em></span>`)
	eHaoYaoPricePrefix = regexp.MustCompile(`<span class="price lFloat"><i>¥</i><em>`)
	eHaoYaoPriceSuffix = regexp.MustCompile(`</em></span>`)	
)

type EHaoYaoUrlNameAndMaxPage struct {
	UrlName string
	MaxPage int
}

func GetAllMaxPage(eHaoYaoUrlNameAndMaxPageChan chan EHaoYaoUrlNameAndMaxPage) {
	for _, cato := range eHaoYaoCatoArr {
		GetMaxPageFromCatoPage(cato, eHaoYaoUrlNameAndMaxPageChan)
		runtime.Gosched()
	}
}

func GetMaxPageFromCatoPage(cato string, eHaoYaoUrlNameAndMaxPageChan chan EHaoYaoUrlNameAndMaxPage) {
	url := `http://www.ehaoyao.com/products/` + cato
	l4g.Debug(url + " begin")
	respBody, err := httpGet(url, false)
	if err != nil {
		l4g.Error(url + ", " + err.Error())
		return
	}
	l4g.Debug(url + " end")
	totalPage := eHaoYaoTotalPage.FindString(respBody)
	totalPage = eHaoYaoTotalPagePrefix.ReplaceAllString(totalPage, "")
	totalPage = eHaoYaoTotalPageSuffix.ReplaceAllString(totalPage, "")
	total, _ := strconv.Atoi(totalPage)
	if total == 0 {
		total = 1
	}
	eHaoYaoUrlNameAndMaxPage := EHaoYaoUrlNameAndMaxPage{}
	eHaoYaoUrlNameAndMaxPage.MaxPage = total
	eHaoYaoUrlNameAndMaxPage.UrlName = url
	eHaoYaoUrlNameAndMaxPageChan <- eHaoYaoUrlNameAndMaxPage
}

func GetNumFromPerpage(eHaoYaoUrlNameAndMaxPageChan chan EHaoYaoUrlNameAndMaxPage, allProductNumChan chan string, f *os.File) {
	for {
		select {
		case eHaoYaoUrlNameAndMaxPage, ok :=<-eHaoYaoUrlNameAndMaxPageChan:
			if !ok {
				l4g.Info("eHaoYaoUrlNameAndMaxPageChan already close")
				close(allProductNumChan)
				return
			}
			for i := 1; i<= eHaoYaoUrlNameAndMaxPage.MaxPage; i++ {
				url := eHaoYaoUrlNameAndMaxPage.UrlName + `?type=1&page=` + strconv.Itoa(i)
				l4g.Debug(url + " begin")
				respBody, err := httpGet(url, false)
				if err != nil {
					l4g.Error(url + ", " + err.Error())
					return
				}
				l4g.Error(url + " end")
				nums := eHaoYaoNum.FindAllString(respBody, -1)
				for _, num := range nums {
					num = eHaoYaoNumPrefix.ReplaceAllString(num, "")
					num = eHaoYaoNumSuffix.ReplaceAllString(num, "")
					allProductNumChan<- num
					_, _ = f.WriteString(num + "\n")
					f.Sync()					
				}
			}
			runtime.Gosched()
		}

	}
}

func GetProductPriceFromPerNum(allProductNumChan chan string) {
	for {
		select {
		case num, ok :=<-allProductNumChan:
			if !ok {
				l4g.Info("allProductNumChan already close")
				return
			}
			url := `http://www.ehaoyao.com/` + num + `.html`
			l4g.Debug(url + " begin")
			respBody, err := httpGet(url, false)
			if err != nil {
				l4g.Error(url + ", " + err.Error())
				return
			}
			productSizeAndPrize := &model.ProductSizeAndPrize{}

			productSizeAndPrize.Num = num
			l4g.Debug("num:", num)

			a_productNameStr := eHaoYaoProductName.FindString(respBody)
			a_productNameStr = eHaoYaoProductNamePrefix.ReplaceAllString(a_productNameStr, "")
			a_productNameStr = eHaoYaoProductNameSuffix.ReplaceAllString(a_productNameStr, "")
			l4g.Debug("药品名称", a_productNameStr)
			productSizeAndPrize.Name = a_productNameStr

			a_approveNumStr := eHaoYaoApproveNum.FindString(respBody)
			a_approveNumStr = eHaoYaoApproveNumPrefix.ReplaceAllString(a_approveNumStr, "")
			a_approveNumStr = eHaoYaoApproveNumSuffix.ReplaceAllString(a_approveNumStr, "")
			l4g.Debug("批准文号：", a_approveNumStr)
			productSizeAndPrize.ApprovalNumber = a_approveNumStr

			a_sizeStr := eHaoYaoCurrentSize.FindString(respBody)	
			a_sizeStr = eHaoYaoCurrentSizePrefix.ReplaceAllString(a_sizeStr, "")
			a_sizeStr = eHaoYaoCurrentSizeSuffix.ReplaceAllString(a_sizeStr, "")
			l4g.Debug("规格2：", a_sizeStr)
			productSizeAndPrize.CurrentSize = a_sizeStr

			a_manufacturerStr := eHaoYaoManualfacturer.FindString(respBody)
			a_manufacturerStr = eHaoYaoManualfacturerPrefix.ReplaceAllString(a_manufacturerStr, "")
			a_manufacturerStr = eHaoYaoManualfacturerSuffix.ReplaceAllString(a_manufacturerStr, "")
			l4g.Debug("生产厂商：", a_manufacturerStr)
			productSizeAndPrize.Manufacturer = a_manufacturerStr

			a_priceStr := eHaoYaoPrice.FindString(respBody)
			a_priceStr = eHaoYaoPricePrefix.ReplaceAllString(a_priceStr, "")
			a_priceStr = eHaoYaoPriceSuffix.ReplaceAllString(a_priceStr, "")
			l4g.Debug("价格：", a_priceStr)
			productSizeAndPrize.Price = a_priceStr

			if productSizeAndPrize.Name != "" {
				SaveProductSizeAndPrize(productSizeAndPrize)
			}
			runtime.Gosched()
		}

	}
}

func SpyProductPriceFromEHaoYao(f *os.File) {
	eHaoYaoUrlNameAndMaxPageChan := make(chan EHaoYaoUrlNameAndMaxPage, 1000)
	allProductNumChan := make(chan string, 20000)
	go GetAllMaxPage(eHaoYaoUrlNameAndMaxPageChan)
	go GetNumFromPerpage(eHaoYaoUrlNameAndMaxPageChan, allProductNumChan, f)
	go GetProductPriceFromPerNum(allProductNumChan)
}

