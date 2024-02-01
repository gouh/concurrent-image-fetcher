package responses

import (
	"concurrent-image-fetcher/internal/models"
	"concurrent-image-fetcher/internal/requests"
)

func GetImagesResponse(request requests.PaginationRequest, totalPages *int, items *[]models.AppFile) *CommonResponse {
	itemsInPage := 0
	if items != nil {
		itemsCount := *items
		itemsInPage = len(itemsCount)
	}

	return &CommonResponse{
		Meta: Meta{
			Page:        request.Page,
			PageSize:    request.PageSize,
			ItemsInPage: &itemsInPage,
			TotalPages:  totalPages,
			Error:       nil,
		},
		Data: *items,
	}
}
