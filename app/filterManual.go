package app

import (
    "fmt"
    "strings"
    "regexp"
    "os"
    "bytes"
    "time"
    "encoding/base32"
    "github.com/pborman/uuid"
)

type MedicineManual struct {
    Id               string `bson:"_id" json:"id"`
    Company          string `bson:"company" json:"company"`        //公司
    Address          string `bson:"address" json:"address"`        //地址
    Telephone        string `bson:"telephone" json:"telephone"`    //联系电话
    DrugName         string `bson:"drugName" json:"drugName"`      //药品名称
    DrugNum          string `bson:"drugNum" json:"drugNum"`      //药品number
    Ingredients      string `bson:"ingredients" json:"ingredients"`//成份
    Indications      string `bson:"indications" json:"indications"`//适应症
    Functions        string `bson:"functions" json:"functions"`    //功能主治
    Usage            string `bson:"usage" json:"usage"`            //用法用量
    AdverseReactions string `bson:"adverseReactions" json:"adverseReactions"`  //不良反应
    Precautions      string `bson:"precautions" json:"precautions"`  //注意事项
    SpecialPopulationMedication string `bson:"specialPopulationMedication" json:"specialPopulationMedication"`  //特殊人群用药
    Pharmacological  string `bson:"pharmacological" json:"pharmacological"` //药理作用
    Forbidens        string `bson:"forbidens" json:"forbidens"` //禁忌
    Interactions     string `bson:"interactions" json:"interactions"` //药物相互作用
    Storage          string `bson:"storage" json:"storage"` //用法用量
    ApprovalNumber   string `bson:"approvalNumber" json:"approvalNumber"`//批准文号
    Manufacturer     string `bson:"manufacturer" json:"manufacturer"` //生产企业
    CreateAt         int64  `bson:"createAt" json:"createAt"`
    UpdateAt         int64  `bson:"updateAt" json:"updateAt`
}

var (
    company          = regexp.MustCompile(`<li class="company">(.+)</li>`)
    address          = regexp.MustCompile(`<li class="address">(.+)</li>`)
    telephone        = regexp.MustCompile(`<li class="telephone">(.+)</li>`)

    findDrugName     = regexp.MustCompile(`\x{3010}([\s]*)\x{836f}\x{54c1}\x{540d}\x{79f0}([\s]*)\x{3011}([^\x{3010}]{0,})`)
    drugName         = regexp.MustCompile(`\x{3010}([\s]*)\x{836f}\x{54c1}\x{540d}\x{79f0}([\s]*)\x{3011}`)
    findIngredients  = regexp.MustCompile(`\x{3010}([\s]*)\x{6210}\x{4efd}([\s]*)\x{3011}([^\x{3010}]{0,})`)
    ingredients      = regexp.MustCompile(`\x{3010}([\s]*)\x{6210}\x{4efd}\x{3011}`)

    //Indications 适应症
    //function 功能主治
    findIndication   = regexp.MustCompile(`\x{3010}([\s]*)\x{9002}\x{5e94}\x{75c7}([\s]*)\x{3011}([^\x{3010}]{0,})`)
    indication       = regexp.MustCompile(`\x{3010}([\s]*)\x{9002}\x{5e94}\x{75c7}([\s]*)\x{3011}`)
    findFunction     = regexp.MustCompile(`\x{3010}([\s]*)\x{529f}\x{80fd}\x{4e3b}\x{6cbb}([\s]*)\x{3011}([^\x{3010}]{0,})`)
    function         = regexp.MustCompile(`\x{3010}([\s]*)\x{529f}\x{80fd}\x{4e3b}\x{6cbb}([\s]*)\x{3011}`)
    findUsage        = regexp.MustCompile(`\x{3010}([\s]*)\x{7528}\x{6cd5}\x{7528}\x{91cf}([\s]*)\x{3011}([^\x{3010}]{0,})`)
    usage            = regexp.MustCompile(`\x{3010}([\s]*)\x{7528}\x{6cd5}\x{7528}\x{91cf}([\s]*)\x{3011}`)
    adverseReactions = regexp.MustCompile(`\x{3010}([\s]*)\x{4e0d}\x{826f}\x{53cd}\x{5e94}([\s]*)\x{3011}`)
    findAdverse      = regexp.MustCompile(`\x{3010}([\s]*)\x{4e0d}\x{826f}\x{53cd}\x{5e94}([\s]*)\x{3011}([^\x{3010}]{0,})`)
    attention        = regexp.MustCompile(`\x{3010}([\s]*)\x{6ce8}\x{610f}\x{4e8b}\x{9879}([\s]*)\x{3011}`)
    findAttention    = regexp.MustCompile(`\x{3010}([\s]*)\x{6ce8}\x{610f}\x{4e8b}\x{9879}([\s]*)\x{3011}([^\x{3010}]{0,})`)
    forbiden         = regexp.MustCompile(`\x{3010}([\s]*)\x{7981}\x{5fcc}([\s]*)\x{3011}`)
    findForbiden     = regexp.MustCompile(`\x{3010}([\s]*)\x{7981}\x{5fcc}([\s]*)\x{3011}([^\x{3010}]{0,})`)
    findSpecialPopulationMedication = regexp.MustCompile(`\x{3010}([\s]*)\x{7279}\x{6b8a}\x{4eba}\x{7fa4}\x{7528}\x{836f}([\s]*)\x{3011}([^\x{3010}]{0,})`)
    //特殊人群用药
    specialPopulationMedication     = regexp.MustCompile(`\x{3010}([\s]*)\x{7279}\x{6b8a}\x{4eba}\x{7fa4}\x{7528}\x{836f}\x{3011}`)
    //药理作用
    findPharmacologicalEffects      = regexp.MustCompile(`\x{3010}([\s]*)\x{836f}\x{7406}\x{4f5c}\x{7528}([\s]*)\x{3011}([^\x{3010}]{0,})`)
    pharmacologicalEffects          = regexp.MustCompile(`\x{3010}([\s]*)\x{836f}\x{7406}\x{4f5c}\x{7528}([\s]*)\x{3011}`)

    interactions     = regexp.MustCompile(`\x{3010}([\s]*)\x{836f}\x{7269}\x{76f8}\x{4e92}\x{4f5c}\x{7528}([\s]*)\x{3011}`)
    findInteractions = regexp.MustCompile(`\x{3010}([\s]*)\x{836f}\x{7269}\x{76f8}\x{4e92}\x{4f5c}\x{7528}([\s]*)\x{3011}([^\x{3010}]{0,})`)
    storage          = regexp.MustCompile(`\x{3010}([\s]*)\x{8d2e}\x{85cf}([\s]*)\x{3011}`)
    findStorage      = regexp.MustCompile(`\x{3010}([\s]*)\x{8d2e}\x{85cf}([\s]*)\x{3011}([^\x{3010}]{0,})`)
    approvement             = regexp.MustCompile(`\x{3010}([\s]*)\x{6279}\x{51c6}\x{6587}\x{53f7}([\s]*)\x{3011}`)
    judgeApprovement        = regexp.MustCompile(`\x{3010}([\s]*)\x{6279}\x{51c6}\x{6587}\x{53f7}([\s]*)\x{3011}([\s\S]*)\x{3010}`)
    findApprovementWithEnd  = regexp.MustCompile(`\x{3010}([\s]*)\x{6279}\x{51c6}\x{6587}\x{53f7}([\s]*)\x{3011}([^\x{3010}]{0,})`)
    findApprovementNoEnd    = regexp.MustCompile(`\x{3010}([\s]*)\x{6279}\x{51c6}\x{6587}\x{53f7}([\s]*)\x{3011}([\s\S\x{4e00}-\x{9fa5}]*)\x{5173}\x{6ce8}\x{6216}\x{8054}\x{7cfb}\x{6211}\x{4eec}`)
    manufacturer     = regexp.MustCompile(`\x{3010}([\s]*)\x{751f}\x{4ea7}\x{4f01}\x{4e1a}([\s]*)\x{3011}`)
    findManufacturer = regexp.MustCompile(`\x{3010}([\s]*)\x{751f}\x{4ea7}\x{4f01}\x{4e1a}([\s]*)\x{3011}([\s\S]{0,})\x{5173}\x{6ce8}\x{6216}\x{8054}\x{7cfb}\x{6211}\x{4eec}`)
    contact          = regexp.MustCompile(`\x{5173}\x{6ce8}\x{6216}\x{8054}\x{7cfb}\x{6211}\x{4eec}`)
    endSign          = regexp.MustCompile(`\x{3010}`)
    wrap             = regexp.MustCompile(`\s`)
)

func SpyMedicineManual(drugNum string) {
    medicineManual := &MedicineManual{}
    url := "http://ypk.39.net/"+ drugNum + "/manual"  
    respBody, err := httpGet(url)
    if err != nil {
        fmt.Println("SpyMedicineManual url:" + url + ", " + err.Error())
        return
    }

    medicineManual.DrugNum = drugNum

    companyMatches := company.FindString(respBody)
    companyRst := strings.TrimSpace(strings.Trim(strings.Trim(companyMatches,`<li class="company">`),`</`))  

    addressMatches := address.FindString(respBody)
    addressRst := strings.TrimSpace(strings.Trim(strings.Trim(addressMatches,`<li class="address">`),`</`))

    telephoneMatches := telephone.FindString(respBody)
    telephoneRst := strings.TrimSpace(strings.Trim(strings.Trim(telephoneMatches,`<li class="telephone">`),`</`))


    //将HTML标签全转换成小写  
    re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
    respBody = re.ReplaceAllStringFunc(respBody, strings.ToLower)
    //去除STYLE  
    re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
    respBody = re.ReplaceAllString(respBody, "")
    //去除SCRIPT  
    re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
    respBody = re.ReplaceAllString(respBody, "")
    //去除所有尖括号内的HTML代码，并换成换行符  
    re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
    respBody = re.ReplaceAllString(respBody, "\n")

    //去除连续的换行符  
    re, _ = regexp.Compile("\\s{1,}")
    respBody = re.ReplaceAllString(respBody, "\n")

    //查找药名
    drugNameMatches := findDrugName.FindString(respBody)
    //去除前缀  
    drugNameMatches = drugName.ReplaceAllString(drugNameMatches, "")
    //去除后缀
    drugNameMatches = endSign.ReplaceAllString(drugNameMatches, "")
    drugNameMatches = wrap.ReplaceAllString(drugNameMatches, " ")           

    //查找成分
    ingredientsMatches := findIngredients.FindString(respBody) 
    ingredientsMatches = ingredients.ReplaceAllString(ingredientsMatches, "")
    ingredientsMatches = endSign.ReplaceAllString(ingredientsMatches, "")
    ingredientsMatches = wrap.ReplaceAllString(ingredientsMatches, "")    

    //查找功能主治
    functionMatches := findFunction.FindString(respBody)
    functionMatches = function.ReplaceAllString(functionMatches, "")
    functionMatches = endSign.ReplaceAllString(functionMatches, "")
    functionMatches = wrap.ReplaceAllString(functionMatches, "")    

    //查找适应症
    indicationMatches := findIndication.FindString(respBody)
    indicationMatches = indication.ReplaceAllString(indicationMatches, "")
    indicationMatches = endSign.ReplaceAllString(indicationMatches, "")
    indicationMatches = wrap.ReplaceAllString(indicationMatches, "")   

    //查找用法
    usageMatches := findUsage.FindString(respBody)
    usageMatches = usage.ReplaceAllString(usageMatches, "")
    usageMatches = endSign.ReplaceAllString(usageMatches, "")
    usageMatches = wrap.ReplaceAllString(usageMatches, "")           

    //查找不良反应
    adverseReactionsMatches := findAdverse.FindString(respBody)
    adverseReactionsMatches = adverseReactions.ReplaceAllString(adverseReactionsMatches, "")
    adverseReactionsMatches = endSign.ReplaceAllString(adverseReactionsMatches, "")
    adverseReactionsMatches = wrap.ReplaceAllString(adverseReactionsMatches, "")  


    //查找禁忌
    forbidenMatches := findForbiden.FindString(respBody)
    forbidenMatches = forbiden.ReplaceAllString(forbidenMatches, "")
    forbidenMatches = endSign.ReplaceAllString(forbidenMatches, "")
    forbidenMatches = wrap.ReplaceAllString(forbidenMatches, "")    

    //查找特殊人群用药
    specialPopulationMedicationMatches := findSpecialPopulationMedication.FindString(respBody)
    specialPopulationMedicationMatches = specialPopulationMedication.ReplaceAllString(specialPopulationMedicationMatches, "")
    specialPopulationMedicationMatches = endSign.ReplaceAllString(specialPopulationMedicationMatches, "")
    specialPopulationMedicationMatches = wrap.ReplaceAllString(specialPopulationMedicationMatches, " ")    


    //查找注意事项
    attentionMatches := findAttention.FindString(respBody)
    attentionMatches = attention.ReplaceAllString(attentionMatches, "")
    attentionMatches = endSign.ReplaceAllString(attentionMatches, "")
    attentionMatches = wrap.ReplaceAllString(attentionMatches, "")     


    //查找相互作用
    interactionMatches := findInteractions.FindString(respBody)
    interactionMatches = interactions.ReplaceAllString(interactionMatches, "")
    interactionMatches = endSign.ReplaceAllString(interactionMatches, "")
    interactionMatches = wrap.ReplaceAllString(interactionMatches, "")      

    //药理作用
    pharmacologicalMatches := findPharmacologicalEffects.FindString(respBody)
    pharmacologicalMatches = pharmacologicalEffects.ReplaceAllString(pharmacologicalMatches, "")
    pharmacologicalMatches = endSign.ReplaceAllString(pharmacologicalMatches, "")
    pharmacologicalMatches = wrap.ReplaceAllString(pharmacologicalMatches, "")

    //查找贮藏方式
    storageMatches := findStorage.FindString(respBody)
    storageMatches = storage.ReplaceAllString(storageMatches, "")
    storageMatches = endSign.ReplaceAllString(storageMatches, "")
    storageMatches = wrap.ReplaceAllString(storageMatches, "")

    //查找许可证号
    approvementMatches := "" 
    if judgeApprovement.MatchString(respBody) {
        approvementMatches = findApprovementWithEnd.FindString(respBody)
        approvementMatches = approvement.ReplaceAllString(approvementMatches, "")
        approvementMatches = endSign.ReplaceAllString(approvementMatches, "")
    } else {
        approvementMatches = findApprovementNoEnd.FindString(respBody)
        approvementMatches = approvement.ReplaceAllString(approvementMatches, "")
        approvementMatches = contact.ReplaceAllString(approvementMatches, "")
    }
    approvementMatches = wrap.ReplaceAllString(approvementMatches, "")

    //查找生产企业
    manufacturerMatches := findManufacturer.FindString(respBody)
    manufacturerMatches = manufacturer.ReplaceAllString(manufacturerMatches, "")
    manufacturerMatches = contact.ReplaceAllString(manufacturerMatches, "")
    manufacturerMatches = wrap.ReplaceAllString(manufacturerMatches, " ")


    if strings.TrimSpace(drugNameMatches) == "" {
        fmt.Println("drug name is null, num:", drugNum)
        return
    }
    var f *os.File
    pg := ""
    if len(drugNum) >= 6 {
        pg = drugNum[:len(drugNum)-5]
    } else {
        pg = "0"
    }
    filename := "./manual/manual" + pg
    if !checkFileIsExist(filename) {
        err = os.MkdirAll(filename, os.ModePerm)
        if err != nil {
            fmt.Println(filename + " mkdir error:" + err.Error())
            return
        }
    }

    f, err = os.Create(filename + "/" + drugNum + ".txt")
    if err != nil {
        fmt.Println("SpyMedicineManual:" + err.Error())
        return
    }
    defer f.Close() 

    fmt.Println("drug name is not null, num:", drugNum)

    _, _ = f.WriteString("药品编号："+ drugNum + "\n\n")
    _, _ = f.WriteString("公司："+ companyRst + "\n\n")
    _, _ = f.WriteString("地址：" + addressRst + "\n\n")
    _, _ = f.WriteString("联系电话：" + telephoneRst + "\n\n")
    _, _ = f.WriteString("药品名称：" + drugNameMatches + "\n\n")
    _, _ = f.WriteString("成份：" + ingredientsMatches + "\n\n")
    _, _ = f.WriteString("功能主治：" + functionMatches + "\n\n")
    _, _ = f.WriteString("适应症：" + indicationMatches + "\n\n")
    _, _ = f.WriteString("用法用量：" + usageMatches + "\n\n")
    _, _ = f.WriteString("不良反应：" + adverseReactionsMatches + "\n\n")
    _, _ = f.WriteString("禁忌：" + forbidenMatches + "\n\n")
    _, _ = f.WriteString("特殊人群用药：" + specialPopulationMedicationMatches + "\n\n")
    _, _ = f.WriteString("注意事项：" + attentionMatches + "\n\n")
    _, _ = f.WriteString("药物相互作用：" + interactionMatches + "\n\n")
    _, _ = f.WriteString("药理作用：" + pharmacologicalMatches + "\n\n")
    _, _ = f.WriteString("生产企业：" + manufacturerMatches + "\n\n")    
    _, _ = f.WriteString("批准文号：" + approvementMatches + "\n\n")    
    _, _ = f.WriteString("贮藏：" + storageMatches + "\n\n")        

    f.Sync()

    medicineManual.PreSave()                                 
    medicineManual.Company = companyRst                               
    medicineManual.Address = addressRst                              
    medicineManual.Telephone = telephoneRst                           
    medicineManual.DrugName = drugNameMatches                           
    medicineManual.DrugNum = drugNum                            
    medicineManual.Ingredients = ingredientsMatches                          
    medicineManual.Indications = indicationMatches                          
    medicineManual.Functions = functionMatches                            
    medicineManual.Usage = usageMatches                              
    medicineManual.AdverseReactions = adverseReactionsMatches                    
    medicineManual.Precautions =  attentionMatches                       
    medicineManual.SpecialPopulationMedication = specialPopulationMedicationMatches         
    medicineManual.Pharmacological = pharmacologicalMatches
    medicineManual.Forbidens = forbidenMatches                           
    medicineManual.Storage = storageMatches   
    medicineManual.ApprovalNumber = approvementMatches
    medicineManual.Manufacturer =  manufacturerMatches   
    medicineManual.Interactions = interactionMatches 

    SaveMedicineManual(medicineManual)

//    if _, err := SaveMedicineManual(medicineManual); err != nil {
//        l4g.Error("medicine name: " + drugNameMatches +  " , " + err.Error())
//    }

    return
}


/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func checkFileIsExist(filename string) bool {
    var exist = true
    if _, err := os.Stat(filename); os.IsNotExist(err) {
        exist = false
    }
    return exist
}


func (u *MedicineManual) PreSave() {
    if u.Id == "" {
        u.Id = NewId()
    }

    u.CreateAt = GetMillis()
    u.UpdateAt = u.CreateAt
}


var encoding = base32.NewEncoding("ybndrfg8e234fdfsxot1uwisza345h769")

func NewId() string {
    var b bytes.Buffer
    encoder := base32.NewEncoder(encoding, &b)
    encoder.Write(uuid.NewRandom())
    encoder.Close()
    b.Truncate(26) // removes the '==' padding
    return b.String()
}

// GetMillis is a convience method to get milliseconds since epoch.
func GetMillis() int64 {
    return time.Now().UnixNano() / int64(time.Millisecond)
}
