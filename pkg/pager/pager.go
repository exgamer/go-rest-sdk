package pager

import "math"

type Pager struct {
	TotalItems       int
	Page             int
	ItemsPerPage     int
	From             int
	To               int
	LastPage         int
	CurrentPageItems int
}

func (pager *Pager) init() {
	if pager.TotalItems == 0 {
		pager.TotalItems = 30
	}

	var pge = 0
	if pager.Page > 0 {
		pge = pager.Page - 1
	}

	pager.From = pager.ItemsPerPage * pge
	if pager.CurrentPageItems == 0 {
		pager.From = 0
	}

	pager.To = pager.From + pager.CurrentPageItems
	pager.From += 1
	pager.LastPage = int(math.Ceil(float64(pager.TotalItems) / float64(pager.ItemsPerPage)))
}

func (pager *Pager) AsMap() map[string]int {
	pager.init()
	params := map[string]int{}
	params["total_items"] = pager.TotalItems
	params["page"] = pager.Page
	params["items_per_page"] = pager.ItemsPerPage
	params["from"] = pager.From
	params["to"] = pager.To
	params["last_page"] = pager.LastPage
	params["current_page_items"] = pager.CurrentPageItems

	return params
}
