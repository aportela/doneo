export const PAGER_DEFAULT_RESULTS_PAGE = 16;

export interface PagerQuery {
  currentPage: number;
  resultsPage: number;
}

export interface PagerResult extends PagerQuery {
  enabled: boolean;
  totalPages: number;
  totalResults: number;
}

export interface Pagination {
  currentPage: number;
  resultsPage: number;
  totalPages: number;
  totalResults: number;
}
