package hubspot

// Page - parameters for a page to be listed
type Page struct {
	Offset int64
	Count  int
}

// PageResponse - response of a list page request
type PageResponse struct {
	Data    []interface{}
	Offset  int64
	HasMore bool
}
