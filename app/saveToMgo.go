package app

import (
	"fmt"
	"github.com/KenmyZhang/medicine-manual-spider/model"
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
	l4g "github.com/alecthomas/log4go"
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
	MedicineManualCollection.EnsureIndex(mgo.Index{
		Key:    []string{"drugNum"},
		Unique: true,
	})

	MedicineProductCollection = MgoSession.DB("spider").C("medicine_products")

	MedicineProductCollection.EnsureIndex(mgo.Index{
		Key:    []string{"currentSize", "name", "manufacturer"},
		Unique: true,
	})
}

func SaveMedicineManual(medicineManual *model.MedicineManual) {
	medicineManual.PreSave()
	err := MedicineManualCollection.Insert(medicineManual)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func SaveProductSizeAndPrize(productSizeAndPrize *model.ProductSizeAndPrize) {
	productSizeAndPrize.PreSave()
	err := MedicineProductCollection.Insert(productSizeAndPrize)
	if err != nil {
		l4g.Error(err.Error())
		return
	}
}
