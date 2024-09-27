package data

import (
	"strings"

	"jordi.tfg.rewrite/internal/validator"
)

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafeList []string
}

func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0 && f.Page <= 10000000, "page", "page must be betwen 1-1000000")
	v.Check(f.PageSize > 0 && f.PageSize <= 100, "pageSize", "pageSize must be betwen 1-1000000")
	v.Check(validator.PermittedValue(f.Sort, f.SortSafeList), "sort", "provide a valid sort value")

}

func (f Filters) SortColumn() string {
	for _, safeValue := range f.SortSafeList {
		if safeValue == f.Sort {
			return strings.TrimPrefix(safeValue, "-")
		}
	}

	panic("Unsafe sort parameter")
}

func (f Filters) SortOrder() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}
