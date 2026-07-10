package pagehandler

import (
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers/userhandler"
	"github.com/aportela/doneo/internal/utils"
)

func addRequestToDomain(request addRequest) domain.Page {
	return domain.Page{
		Title: request.Title,
	}
}

func updateRequestToDomain(request updateRequest) domain.Page {
	return domain.Page{
		ID:    request.ID,
		Title: request.Title,
		Body:  request.Body,
	}
}

func domainToResponse(page domain.Page) pageResponse {
	return pageResponse{
		ID:        page.ID,
		Title:     page.Title,
		Body:      page.Body,
		CreatedBy: userhandler.BaseDomainToBaseResponse(page.CreatedBy),
		CreatedAt: page.CreatedAt.UnixMilli(),
		UpdatedAt: utils.TimePtrToInt64Ptr(page.UpdatedAt),
	}
}

func domainArrayToResponseArray(pages []domain.Page) []pageResponse {
	pagesResponse := []pageResponse{}
	for _, page := range pages {
		pagesResponse = append(pagesResponse, domainToResponse(page))
	}
	return pagesResponse
}

func toSearchResponse(pages []domain.Page) searchResponse {
	return searchResponse{
		Pages: domainArrayToResponseArray(pages),
	}
}
