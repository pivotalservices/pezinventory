package integrations

import "gopkg.in/mgo.v2"

type (
	//Collection - an interface representing a trimmed down collection object
	Collection interface {
		Wake()
		Close()
		Find(selector interface{}, result interface{}) (err error)
		FindOne(id string, result interface{}) (err error)
		UpsertID(selector interface{}, update interface{}) (info *mgo.ChangeInfo, err error)
		FindAndModify(selector interface{}, update interface{}, target interface{}) (info *mgo.ChangeInfo, err error)
		Count() (int, error)
	}

	//CollectionRepo - mgo collection adaptor
	CollectionRepo struct {
		Col     *mgo.Collection
		session *mgo.Session
	}

	//CollectionDialer - a funciton type to dial for collections
	CollectionDialer func(url string, dbname string, collectionname string) (collection Collection, err error)
)
