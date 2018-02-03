package model

import (
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

func (u *MedicineManual) PreSave() {
    if u.Id == "" {
        u.Id = NewId()
    }

    u.CreateAt = GetMillis()
    u.UpdateAt = u.CreateAt
}
