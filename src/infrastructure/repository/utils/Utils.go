package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gbrayhan/microservices-go/src/domain"
	"gorm.io/gorm"
)

func ComplementSearch(r *gorm.DB, sortBy string, sortDirection string, limit int64, offset int64, filters map[string][]string, dateRangeFilters []domain.DateRangeFilter, searchText string, searchColumns []string, columnMapping map[string]string) (query *gorm.DB, err error) {
	if r == nil {
		return nil, nil
	}

	query = r
	if sortBy != "" {
		orderClause := fmt.Sprintf("%s %s", columnMapping[sortBy], sortDirection)
		query = query.Order(orderClause).Limit(int(limit)).Offset(int(offset))
	} else {
		query = query.Limit(int(limit)).Offset(int(offset))
	}

	if len(filters) > 0 {
		filters = UpdateFilterKeys(filters, columnMapping)
		for key, values := range filters {
			query = query.Where(fmt.Sprintf("%s IN (?)", key), values)
		}
	}

	if len(dateRangeFilters) > 0 {
		for i := range dateRangeFilters {
			if newFieldName, ok := columnMapping[dateRangeFilters[i].Field]; ok {
				dateRangeFilters[i].Field = newFieldName
			}
		}
		for _, filter := range dateRangeFilters {
			query = query.Where(fmt.Sprintf("%s BETWEEN ? AND ?", filter.Field), filter.Start, filter.End)
		}
	}

	if searchText != "" {
		var orConditions []string
		for _, column := range searchColumns {
			orConditions = append(orConditions, fmt.Sprintf("%s LIKE '%%%s%%'", column, searchText))
		}
		searchQuery := fmt.Sprintf("AND (%s)", strings.Join(orConditions, " OR "))
		query = query.Where(fmt.Sprintf("1=1 %s", searchQuery))
	}
	return
}

func UpdateFilterKeys(filters map[string][]string, columnMapping map[string]string) map[string][]string {
	updatedFilters := make(map[string][]string)
	for key, value := range filters {
		if updatedKey, ok := columnMapping[key]; ok {
			updatedFilters[updatedKey] = value
		} else {
			updatedFilters[key] = value
		}
	}
	return updatedFilters
}

func ApplyFilters(columnMapping map[string]string, filters map[string][]string, dateRangeFilters []domain.DateRangeFilter, searchText string, searchColumns []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		query := db
		if len(filters) > 0 {
			filters = UpdateFilterKeys(filters, columnMapping)
			for key, values := range filters {
				query = query.Where(fmt.Sprintf("%s IN (?)", key), values)
			}
		}
		if len(dateRangeFilters) > 0 {
			for _, filter := range dateRangeFilters {
				if newFieldName, ok := columnMapping[filter.Field]; ok {
					filter.Field = newFieldName
				}
				query = query.Where(fmt.Sprintf("%s BETWEEN ? AND ?", filter.Field), filter.Start, filter.End)
			}
		}
		if searchText != "" && len(searchColumns) > 0 {
			var orConditions []string
			var args []interface{}
			for _, column := range searchColumns {
				orConditions = append(orConditions, fmt.Sprintf("%s LIKE ?", column))
				args = append(args, "%"+searchText+"%")
			}
			searchQuery := fmt.Sprintf("(%s)", strings.Join(orConditions, " OR "))
			query = query.Where(searchQuery, args...)
		}
		return query
	}
}

func IsZeroValue(value any) bool {
	return reflect.DeepEqual(value, reflect.Zero(reflect.TypeOf(value)).Interface())
}
