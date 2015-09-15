package pezinventory

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pivotal-pez/pezinventory/service/integrations"

	"gopkg.in/mgo.v2/bson"
)

//Lease wraps the leases collection.
type Lease struct {
	ID                bson.ObjectId          `bson:"_id,omitempty" json:"id"`
	InventoryItemID   bson.ObjectId          `bson:"inventory_item_id,omitempty" json:"inventory_item_id"`
	User              string                 `json:"user"`
	DurationDays      int                    `json:"duration_days"`
	StartDate         string                 `json:"start_date"`
	EndDate           string                 `json:"end_date"`
	Status            string                 `json:"status"`
	Attributes        map[string]interface{} `json:"attributes"`
	PrivateAttributes map[string]interface{} `json:"private_attributes,omitempty"`
}

//RedactedLease wraps the leases collection omitting private attributes.
type RedactedLease struct {
	ID              bson.ObjectId          `bson:"_id,omitempty" json:"id"`
	InventoryItemID bson.ObjectId          `bson:"inventory_item_id,omitempty" json:"inventory_item_id"`
	User            string                 `json:"user"`
	DurationDays    int                    `json:"duration_days"`
	StartDate       string                 `json:"start_date"`
	EndDate         string                 `json:"end_date"`
	Status          string                 `json:"status"`
	Attributes      map[string]interface{} `json:"attributes"`
}

//FindLeaseByIDHandler will return a redacted lease record for the given ID.
func FindLeaseByIDHandler(collection integrations.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		collection.Wake()

		id := mux.Vars(req)["id"]
		if id == "" {
			Formatter().JSON(w, http.StatusBadRequest, errorMessage("lease id must be specified"))
			return
		}

		lease := new(RedactedLease)
		if err := collection.FindOne(id, lease); err == nil {
			Formatter().JSON(w, http.StatusOK, successMessage(lease))
		} else {
			log.Println("lease lookup failed")
			Formatter().JSON(w, http.StatusNotFound, errorMessage(err.Error()))
		}
	}
}

//LeaseInventoryItemHandler creates a new lease record against an available InventoryItem
//and calls dispenser to provision that InventoryItem to the requestor.
//
//Unless supplied, the StartDate and EndDate values will be calculated according to the
//time of the invocation and the DurationDays value.  If DurationDays is not supplied, it
//will default to 14.
//
//NOTE: The call to dispenser is not yet implemented.
func LeaseInventoryItemHandler(ic integrations.Collection, lc integrations.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var obj Lease
		decoder := json.NewDecoder(req.Body)

		err := decoder.Decode(&obj)
		if err != nil {
			Formatter().JSON(w, http.StatusBadRequest, errorMessage(err.Error()))
			return
		}

		if obj.InventoryItemID == "" {
			Formatter().JSON(w, http.StatusBadRequest, errorMessage("inventory_item_id must be specified"))
			return
		}

		if obj.StartDate == "" || obj.EndDate == "" {
			if obj.DurationDays <= 0 {
				obj.DurationDays = 14
			}
			epoch := time.Now()
			obj.StartDate = epoch.String()
			obj.EndDate = epoch.AddDate(0, 0, obj.DurationDays).String()
		}

		err = InventoryItemReservingStatus(obj.InventoryItemID, ic)
		if err != nil {
			Formatter().JSON(w, http.StatusNotFound, errorMessage(err.Error()))
			return
		}

		lc.Wake()
		obj.ID = bson.NewObjectId()
		_, err = lc.UpsertID(obj.ID, obj)
		if err != nil {
			e := InventoryItemAvailableStatus(obj.InventoryItemID, ic)
			if e != nil {
				log.Printf("Could not release Inventory Item %s", obj.InventoryItemID.Hex())
			}
			Formatter().JSON(w, http.StatusInternalServerError, errorMessage(err.Error()))
			return
		}

		err = InventoryItemLeasedStatus(obj.InventoryItemID, obj.ID, ic)
		if err != nil {
			Formatter().JSON(w, http.StatusInternalServerError, errorMessage(err.Error()))
			return
		}

		Formatter().JSON(w, http.StatusOK, obj)
	}
}
