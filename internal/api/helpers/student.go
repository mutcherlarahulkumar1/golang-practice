package helpers

import (
	"fmt"
	"net/http"
	"strings"
)

func AddSortParams(request *http.Request, query string) string {
	sortParams := request.URL.Query()["sortby"]
	if len(sortParams) > 0 {
		query += " ORDER BY "
		for i, param := range sortParams {
			parts := strings.Split(param, ":")

			field, order := parts[0], parts[1]
			if i > 0 {
				query += ","
			}
			query += " " + field + " " + order
		}

	}
	return query
}

func AddFilters(request *http.Request, query string, args []any) (string, []any) {
	params := map[string]string{
		"firstname": "firstName",
		"lastname":  "lastName",
		"email":     "email",
		"class":     "class",
	}

	i := 1
	for param, dbField := range params {
		value := request.URL.Query().Get(param)
		if value != "" {
			query += fmt.Sprintf(" AND %s = $%d", dbField, i)
			args = append(args, value)
			i++
		}
	}
	return query, args
}
