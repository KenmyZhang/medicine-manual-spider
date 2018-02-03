package model

type ProductSizeAndPrize struct {
    Id               string `bson:"_id" json:"id"`
    Name             string `bson:"name" json:"name"`
    Price            string `bson:"price" json:"price"`
    Num              string `bson:"num" json:"num"`
    Manufacturer     string `bson:"manufacturer" json:"manufacturer"`
    ApprovalNumber   string `bson:"approvalNumber" json:"approvalNumber"`
    CurrentSize      string `bson:"currentSize" json:"currentSize"`     
    AllSize          string `bson:"allSize" json:"allSize"`
    CreateAt         int64  `bson:"createAt" json:"createAt"`
    UpdateAt         int64  `bson:"updateAt" json:"updateAt`
}

func (u *ProductSizeAndPrize) PreSave() {
    if u.Id == "" {
        u.Id = NewId()
    }

    u.CreateAt = GetMillis()
    u.UpdateAt = u.CreateAt
}
