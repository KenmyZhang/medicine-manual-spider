package app

import (
    "fmt"
	"log"
    "gopkg.in/mgo.v2"
    "github.com/KenmyZhang/medicine-manual-spider/model"
    //"gopkg.in/mgo.v2/bson"
)

var MedicineManualCollection *mgo.Collection
var MedicineProductCollection *mgo.Collection
var MgoSession *mgo.Session

func init() {
    MgoSession, err := mgo.Dial("localhost")
    if err != nil {
        panic(err)
    }

    // Optional Switch the MgoSession to a monotonic behavior.
    MgoSession.SetMode(mgo.Monotonic, true)

    MedicineManualCollection = MgoSession.DB("spider").C("medicine_manuals")
    MedicineProductCollection = MgoSession.DB("spider").C("medicine_products")

    MedicineProductCollection.EnsureIndex(mgo.Index{
        Key:    []string{"CurrentSize", "Name"},
        Unique: true,
    })


}

func SaveMedicineManual(medicineManual *model.MedicineManual) {
    medicineManual.PreSave()    
	err := MedicineManualCollection.Insert(medicineManual)
    if err != nil {
        log.Fatal(err)
        return
    }
}

func SaveProductSizeAndPrize(productSizeAndPrize *model.ProductSizeAndPrize) {
    productSizeAndPrize.PreSave()   
    fmt.Println(productSizeAndPrize) 
    err := MedicineProductCollection.Insert(productSizeAndPrize)
    if err != nil {
        log.Fatal(err)
        return
    }
}