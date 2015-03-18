package controllers

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/pivotalservices/pezinventory/models"
)

//ItemController - controller for searching item information
type ItemController struct {
}

//ListItems - returns a collection of Items
func (c *ItemController) ListItems(render render.Render) {
	i := models.ListItems()
	render.JSON(200, successMessage(i))
}

//GetItem - returns a single Item record
func (c *ItemController) GetItem(params martini.Params, render render.Render) {
	i := models.GetItem(params["id"])
	render.JSON(200, successMessage(i))
}

//GetItemHistory - returns the history for a given Item
func (c *ItemController) GetItemHistory(params martini.Params, render render.Render) {
	render.JSON(501, errorMessage("Not Implemented"))
}
