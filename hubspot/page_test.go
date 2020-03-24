package hubspot

import (
	"testing"

	"github.com/go-errors/errors"
	"github.com/stretchr/testify/require"
)

func TestFindInPages(t *testing.T) {
	processedpages := 0
	processeditems := 0
	item, err := FindInPages(func(page *Page) (*PageResponse, error) {
		processedpages++
		var offset int64
		var data []interface{}

		if page == nil {
			offset = 1
			data = append(data, int64(8), int64(22), int64(102), int64(7739))
		} else {
			offset = page.Offset
			data = append(data, 8+offset, 22+offset, 102+offset, 7739+offset)
			offset++
		}
		return &PageResponse{
			Offset:  offset,
			Data:    data,
			HasMore: offset < 8,
		}, nil
	},
		func(item interface{}) (bool, error) {
			processeditems++
			number, ok := item.(int64)
			if !ok {
				return false, errors.Errorf("Invalid item type")
			}
			return number == 105, nil
		})

	require.NoError(t, err)
	require.Equal(t, int64(105), item)
	require.Equal(t, 4, processedpages)
	require.Equal(t, 15, processeditems)
}
