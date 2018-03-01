package app

import(
	"github.com/KenmyZhang/medicine-manual-spider/model"
	l4g "github.com/alecthomas/log4go"
	"regexp"	
	"strconv"
	"time"
)

var Yao123Cato = []string{
	`https://www.yao123.com/category_8`,
	`https://www.yao123.com/category_30`,
	`https://www.yao123.com/category_32`,
	`https://www.yao123.com/category_33`,
	`https://www.yao123.com/category_59`,
	`https://www.yao123.com/category_60`,
	`https://www.yao123.com/category_61`,
	`https://www.yao123.com/category_87`,
	`https://www.yao123.com/category_104`,
	`https://www.yao123.com/category_112`,
	`https://www.yao123.com/category_119`,
	`https://www.yao123.com/category_140`,
	`https://www.yao123.com/category_209`,
	`https://www.yao123.com/category_757`}

var (
	yao123TotalPage       = regexp.MustCompile(`<span>共<b class="totalpages">[0-9]+?</b>页`)
	yao123TotalPagePrefix = regexp.MustCompile(`<span>共<b class="totalpages">`)
	yao123TotalPageSuffix = regexp.MustCompile(`</b>页`)

	yao123Num       = regexp.MustCompile(`<a title="" target="_blank" href="/product/[\s\S]+?" class="name">`)
	yao123NumPrefix = regexp.MustCompile(`<a title="" target="_blank" href="/product/`)
	yao123NumSuffix = regexp.MustCompile(`" class="name">`)


	yao123ProductName       = regexp.MustCompile(`<li><span>商品名称：</span><span class="lv">[\s\S]+?</span></li>`)
	yao123ProductNamePrefix = regexp.MustCompile(`<li><span>商品名称：</span><span class="lv">`)
	yao123ProductNameSuffix = regexp.MustCompile(`</span></li>`)

	yao123ApproveNum       = regexp.MustCompile(`<li><span>批准文号：</span><span class="lv">[\s\S]+?</span></li>`)
	yao123ApproveNumPrefix = regexp.MustCompile(`<li><span>批准文号：</span><span class="lv">`)
	yao123ApproveNumSuffix = regexp.MustCompile(`</span></li>`)

	yao123Manualfacturer       = regexp.MustCompile(`<li><span>生产厂家（生产企业）：</span><span class="lv">[\s\S]+?</span></li>`)
	yao123ManualfacturerPrefix = regexp.MustCompile(`<li><span>生产厂家（生产企业）：</span><span class="lv">`)
	yao123ManualfacturerSuffix = regexp.MustCompile(`</span></li>`)

	yao123CurrentSize       = regexp.MustCompile(`<li><span>规格：</span><span class="lv">[\s\S]+?</span>`)
	yao123CurrentSizePrefix = regexp.MustCompile(`<li><span>规格：</span><span class="lv">`)
	yao123CurrentSizeSuffix = regexp.MustCompile(`</span>`)	
)

type CatoNameAndMaxPage struct {
	Catoname string
	MaxPage int
}

func GetMaxPageFromPerCato(catoUrl string, catoAndMaxPageChan chan *CatoNameAndMaxPage) {
	l4g.Debug(catoUrl + " begin")
	respBody, err := httpGet(catoUrl, false)
	if err != nil {
		l4g.Error(catoUrl + ", " + err.Error())
		return
	}
	totalPage := yao123TotalPage.FindString(respBody)
	totalPage = yao123TotalPagePrefix.ReplaceAllString(totalPage, "")
	totalPage = yao123TotalPageSuffix.ReplaceAllString(totalPage, "")	
	maxPage, _ := strconv.Atoi(totalPage)
	if maxPage == 0 {
		maxPage = 1
	}
	catoAndMaxPage := &CatoNameAndMaxPage{}
	catoAndMaxPage.Catoname = catoUrl
	catoAndMaxPage.MaxPage = maxPage 
	catoAndMaxPageChan <- catoAndMaxPage
}

func GetAllCatoMaxPage(catoAndMaxPageChan chan *CatoNameAndMaxPage) {
	for _, cato := range Yao123Cato {
		GetMaxPageFromPerCato(cato, catoAndMaxPageChan)
	}
	close(catoAndMaxPageChan)
}	

func GetAllProductNum(catoAndMaxPageChan chan *CatoNameAndMaxPage, allProductNumChan chan string) {
	for {
		select {
			case catoAndMaxPage, ok :=<- catoAndMaxPageChan:
				if !ok {
					l4g.Info("catoAndMaxPageChan already close")
					return
				}

				l4g.Error("catoAndMaxPage.MaxPage;:", catoAndMaxPage.MaxPage)
				for i := 1; i <= catoAndMaxPage.MaxPage; i++ {
					url := catoAndMaxPage.Catoname + "?pageNum=" + strconv.Itoa(i)
					GetProductNumFromPerPage(url, allProductNumChan)				
				}
			case <-time.After(time.Minute * 1):
				l4g.Error("timeout")
				return
		}
	}
}

func GetProductNumFromPerPage(url string, allProductNumChan chan string) {
	respBody, err := httpGet(url, false)
	if err != nil {
		l4g.Error(url + ", " + err.Error())
		return
	}	

    nums := yao123Num.FindAllString(respBody, -1)
    for _, num := range nums {
      num = yao123NumPrefix.ReplaceAllString(num, "")
      num = yao123NumSuffix.ReplaceAllString(num, "")
      allProductNumChan <- num
    }
}

func GetPerProductPrice(allProductNumChan chan string) {
	for {
		select{
			case num, ok := <-allProductNumChan:
				if !ok {
					l4g.Info("allProductNumChan already close")
					return				
				}
				FilterProductPrice(num)
				return
			case <-time.After(time.Minute * 1):
				l4g.Error("timeout")
				return
		}
	}
}

func FilterProductPrice(numStr string) {
	//url := `https://www.yao123.com/product/` + numStr
	url := `http://www.ehaoyao.com/product-33030.html`
	l4g.Debug(url + "  begin")
	respBody, err := httpGet(url, false)
	if err != nil {
		l4g.Error(url + ", " + err.Error())
		return
	}	
	l4g.Debug(url + "  end")
	l4g.Error("respBody", respBody)

//	_, _ = f.WriteString(numStr + "\n")
//	f.Sync()		

	productSizeAndPrize := &model.ProductSizeAndPrize{}

	productSizeAndPrize.Num = numStr
	l4g.Debug("numStr:", numStr)

	a_productNameStr := yao123ProductName.FindString(respBody)
	a_productNameStr = yao123ProductNamePrefix.ReplaceAllString(a_productNameStr, "")
	a_productNameStr = yao123ProductNameSuffix.ReplaceAllString(a_productNameStr, "")
	l4g.Error("药品名称", a_productNameStr)
	productSizeAndPrize.Name = a_productNameStr

	a_approveNumStr := yao123ApproveNum.FindString(respBody)
	a_approveNumStr = yao123ApproveNumPrefix.ReplaceAllString(a_approveNumStr, "")
	a_approveNumStr = yao123ApproveNumSuffix.ReplaceAllString(a_approveNumStr, "")
	l4g.Error("批准文号：", a_approveNumStr)
	productSizeAndPrize.ApprovalNumber = a_approveNumStr

	a_sizeStr := yao123CurrentSize.FindString(respBody)	
	a_sizeStr = yao123CurrentSizePrefix.ReplaceAllString(a_sizeStr, "")
	a_sizeStr = yao123CurrentSizeSuffix.ReplaceAllString(a_sizeStr, "")
	l4g.Error("规格2：", a_sizeStr)
	productSizeAndPrize.CurrentSize = a_sizeStr

	a_manufacturerStr := yao123Manualfacturer.FindString(respBody)
	a_manufacturerStr = yao123ManualfacturerPrefix.ReplaceAllString(a_manufacturerStr, "")
	a_manufacturerStr = yao123ManualfacturerSuffix.ReplaceAllString(a_manufacturerStr, "")
	l4g.Error("生产厂商：", a_manufacturerStr)
	productSizeAndPrize.Manufacturer = a_manufacturerStr

//	a_priceStr := yao123Price.FindString(respBody)
//	a_priceStr = yao123PricePrefix.ReplaceAllString(a_priceStr, "")
//	a_priceStr = yao123PriceSuffix.ReplaceAllString(a_priceStr, "")
//	l4g.Error("价格：", a_priceStr)
//	productSizeAndPrize.Price = a_priceStr

	//        a_yao123AllSize := yao123AllSize.FindString(respBody)
	//        otherSizeNums := yao123OtherSize.FindAllString(a_yao123AllSize, -1)
	//        for _, num := range otherSizeNums {
	//          num = yao123OtherSizePrefix.ReplaceAllString(num, "")
	//          num = yao123OtherSizeSuffix.ReplaceAllString(num, "")
	//          l4g.Debug("num",num)
	//          numChan <- num
	//        }
	//if productSizeAndPrize.Name != "" {
	//	SaveProductSizeAndPrize(productSizeAndPrize)
	//}
}

func SpyProductPriceFromYao123() {
	catoAndMaxPageChan := make(chan *CatoNameAndMaxPage, 1000)
	allProductNumChan := make(chan string, 10000)
	go GetAllCatoMaxPage(catoAndMaxPageChan)
	go GetAllProductNum(catoAndMaxPageChan, allProductNumChan)
	go GetPerProductPrice(allProductNumChan)
}