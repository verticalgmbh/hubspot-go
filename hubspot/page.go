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

// NewPage - creates a new page parameter
//
// **Parameters**
//   offset: offset where to start listing (use offset returned by a former page request)
//   count : number of items to return. If <= 0 is specified as count the request default is used
func NewPage(offset int64, count int) *Page {
	return &Page{
		Offset: offset,
		Count:  count}
}

// FindInPages - finds an item in results of a paging mechanism
func FindInPages(pager func(*Page) (*PageResponse, error), predicate func(data interface{}) (bool, error)) (interface{}, error) {
	var page *Page
	for {
		response, err := pager(page)
		if err != nil {
			return nil, err
		}

		for _, item := range response.Data {
			ok, err := predicate(item)
			if err != nil {
				return nil, err
			}

			if ok {
				return item, nil
			}
		}

		if response.HasMore {
			page = NewPage(response.Offset, 0)
		} else {
			break
		}
	}

	return nil, nil
}
