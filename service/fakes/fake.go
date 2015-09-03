package fakes

import (
	"encoding/json"

	"github.com/pivotal-pez/pezinventory/service"
	"github.com/pivotal-pez/pezinventory/service/integrations"
	"gopkg.in/mgo.v2"
)

//FakeNewCollectionDialer -
func FakeNewCollectionDialer(c []pezinventory.InventoryItem) func(url, dbname, collectionname string) (col integrations.Collection, err error) {
	b, err := json.Marshal(c)
	if err != nil {
		panic("shit is broken")
	}

	return func(url, dbname, collectionname string) (col integrations.Collection, err error) {
		col = &FakeCollection{
			Data: b,
		}
		return
	}
}

//FakeCollection -
type FakeCollection struct {
	mgo.Collection
	Data  []byte
	Error error
}

//Close -
func (s *FakeCollection) Close() {

}

//Wake -
func (s *FakeCollection) Wake() {

}

//Find -- finds all records matching given selector
func (s *FakeCollection) Find(selector interface{}, result interface{}) (err error) {
	err = json.Unmarshal(s.Data, result)
	return
}

//FindAndModify -
func (s *FakeCollection) FindAndModify(selector interface{}, update interface{}, result interface{}) (info *mgo.ChangeInfo, err error) {
	return
}

//UpsertID -
func (s *FakeCollection) UpsertID(id interface{}, result interface{}) (changInfo *mgo.ChangeInfo, err error) {
	return
}

//FindOne -
func (s *FakeCollection) FindOne(id string, result interface{}) (err error) {
	err = json.Unmarshal(s.Data, result)
	return
}
