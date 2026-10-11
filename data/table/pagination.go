package table

import (
	"fmt"
	"strconv"

	"lazer/laze"
	exception "lazer/error"
)

func getPaginationParams(params map[string][]string) (int, int, int, laze.Exception) {
	limit := 10
	offset := 0
	page := 1
	if rawPageSize, ok := params["pageSize"]; ok {
		pageSize, err := strconv.Atoi(rawPageSize[0])
		if err != nil {
			ex := exception.FromError(err, exception.INTERNALERROR)
			return page, limit, offset, ex
		}
		limit = pageSize
	}
	if rawPage, ok := params["page"]; ok {
		pageInt, err := strconv.Atoi(rawPage[0])
		if err != nil {
			ex := exception.FromError(err, exception.INTERNALERROR)
			return page, limit, offset, ex
		}
		if pageInt < 1 {
			pageInt = 1
		}
		page = pageInt
		offset = (page - 1) * limit
	}
	return page, limit, offset, nil
}

func createPaginationString(limit int, offset int) string {
	paginationStr := fmt.Sprintf(" LIMIT %v OFFSET %v", limit, offset)
	return paginationStr
}

func (table *Table) getPagination(params map[string][]string) (string, int, int, int) {
	page, limit, offset, err := getPaginationParams(params)
	if err != nil {
		fmt.Println(err)
	}
	paginationStr := createPaginationString(limit, offset)
	return paginationStr, page, limit, offset
}