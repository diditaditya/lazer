package table

import (
	"fmt"
	"strconv"

	"lazer/laze"
	exception "lazer/error"
)

func getPaginationParams(params map[string][]string) (int, int, laze.Exception) {
	limit := 10
	offset := 0
	if rawPageSize, ok := params["pageSize"]; ok {
		pageSize, err := strconv.Atoi(rawPageSize[0])
		if err != nil {
			ex := exception.FromError(err, exception.INTERNALERROR)
			return limit, offset, ex
		}
		limit = pageSize
	}
	if rawPage, ok := params["page"]; ok {
		page, err := strconv.Atoi(rawPage[0])
		if err != nil {
			ex := exception.FromError(err, exception.INTERNALERROR)
			return limit, offset, ex
		}
		if page < 1 {
			page = 1
		}
		offset = (page - 1) * limit
	}
	return limit, offset, nil
}

func (table *Table) createPaginationString(params map[string][]string) string {
	limit, offset, err := getPaginationParams(params)
	if err != nil {
		fmt.Println(err)
	}
	paginationStr := fmt.Sprintf(" LIMIT %v OFFSET %v", limit, offset)
	return paginationStr
}