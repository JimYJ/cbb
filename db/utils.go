package db

import (
	"fmt"
	"math"
	"strconv"
)

// Pagina 分页
func Pagina(pageSize, pageNo, totalCount string) (string, int) {
	pageNoInt, err := strconv.Atoi(pageNo)
	pageSizeInt, err2 := strconv.Atoi(pageSize)
	totalCountInt, err3 := strconv.Atoi(totalCount)
	if err != nil || err2 != nil || err3 != nil {
		return "", 0
	}
	if pageNoInt < 1 {
		pageNoInt = 1
	}
	if pageSizeInt < 1 {
		pageSizeInt = 10
	}
	pageTotal := int(math.Ceil(float64(totalCountInt) / float64(pageSizeInt)))
	if pageNoInt > pageTotal {
		pageNoInt = pageTotal
	}
	if pageSizeInt > totalCountInt {
		pageSizeInt = totalCountInt
	}
	pageStart := (pageNoInt - 1) * pageSizeInt
	pageEnd := pageNoInt * pageSizeInt
	if pageEnd > totalCountInt {
		pageEnd = totalCountInt
	}
	// log.Println(pageNoInt, pageSizeInt, pageStart, pageTotal, pageEnd)
	sql := fmt.Sprintf("limit %d,%d", pageStart, pageEnd-pageStart)
	return sql, pageTotal
}
