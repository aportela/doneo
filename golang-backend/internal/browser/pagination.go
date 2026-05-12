package browser

import (
	"math"
)

const (
	DefaultResultsPage = 20
	MaxResultsPage     = 500
)

type Params struct {
	CurrentPage int
	ResultsPage int
}

func (p Params) getCurrentPage() int {
	if p.CurrentPage <= 0 {
		return 1
	}
	return p.CurrentPage
}

func (p Params) getResultsPage() int {
	if p.ResultsPage < 0 {
		return DefaultResultsPage
	}

	if p.ResultsPage > MaxResultsPage {
		return MaxResultsPage
	}

	return p.ResultsPage
}

func (p Params) Limit() int {
	return p.getResultsPage()
}

func (p Params) Offset() int {
	return (p.getCurrentPage() - 1) * p.getResultsPage()
}

func (p Params) Enabled() bool {
	return p.ResultsPage > 0
}

type Result struct {
	CurrentPage  int
	ResultsPage  int
	TotalResults int
	TotalPages   int
}

func NewResult(params Params, totalResults int) Result {
	resultsPage := params.getResultsPage()
	currentPage := params.getCurrentPage()

	totalPages := 0
	if resultsPage > 0 && totalResults > 0 {
		totalPages = int(math.Ceil(float64(totalResults) / float64(resultsPage)))
	}

	return Result{
		CurrentPage:  currentPage,
		ResultsPage:  resultsPage,
		TotalResults: totalResults,
		TotalPages:   totalPages,
	}
}
