package app

import (
    "fmt"
	"log"
    "gopkg.in/mgo.v2"
    //"gopkg.in/mgo.v2/bson"
)

var MedicineManualCollection *mgo.Collection
var MgoSession *mgo.Session

func init() {
    MgoSession, err := mgo.Dial("localhost")
    if err != nil {
        panic(err)
    }

    // Optional Switch the MgoSession to a monotonic behavior.
    MgoSession.SetMode(mgo.Monotonic, true)

    MedicineManualCollection = MgoSession.DB("spider").C("medicine_manuals")
}

func SaveMedicineManual(medicineManual *MedicineManual) {
	err := MedicineManualCollection.Insert(medicineManual)
    if err != nil {
        log.Fatal(err)
        return
    }
    fmt.Println("medicineManual:", medicineManual)
}