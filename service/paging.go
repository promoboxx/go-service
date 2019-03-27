package service

import (
	"fmt"
	"net/http"
	"strconv"
)

// PagingParams represents paging and sorting parameter values
type PagingParams struct {
	PageSize   *int32
	PageNumber *int32
	SortField  *string
	SortAsc    *bool
}

// ParsePagingParams retrieves paging params from the request, allows for whitelisting sort field values
func ParsePagingParams(r *http.Request, sortFieldsWhitelist []string) (PagingParams, error) {
	paging := PagingParams{}

	// page number
	pageNumber, err := Int32PointerFromQueryParam(r, "page_number")
	if err != nil {
		return paging, err
	}
	paging.PageNumber = pageNumber

	// page size
	pageSize, err := Int32PointerFromQueryParam(r, "page_size")
	if err != nil {
		return paging, err
	}
	paging.PageSize = pageSize

	// sort asc
	sortAscStr := r.URL.Query().Get("sort_asc")
	if len(sortAscStr) > 0 {
		sortAsc, err := strconv.ParseBool(sortAscStr)
		if err != nil {
			return paging, err
		}
		paging.SortAsc = &sortAsc
	}

	// sort field
	sortField := r.URL.Query().Get("sort_field")
	if len(sortField) > 0 {
		if !contains(sortFieldsWhitelist, sortField) {
			sortFieldErr := fmt.Errorf("invalid sort field %s", sortField)
			return paging, sortFieldErr
		}
		paging.SortField = &sortField
	}

	return paging, nil
}

// Contains tells whether a contains x.
func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
