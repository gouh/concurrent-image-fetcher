package requests

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type PaginationRequest struct {
	Page     *int
	PageSize *int
}

func (pagination *PaginationRequest) GetDefaultPaginationParams(c *gin.Context) PaginationRequest {
	page, pageSize := 1, 30

	if p, ok := c.GetQuery("page"); ok {
		if pInt, err := strconv.Atoi(p); err == nil && pInt > 0 {
			page = pInt
		}
	}

	if ps, ok := c.GetQuery("pageSize"); ok {
		if psInt, err := strconv.Atoi(ps); err == nil && psInt > 0 {
			pageSize = psInt
		}
	}

	return PaginationRequest{Page: &page, PageSize: &pageSize}
}
