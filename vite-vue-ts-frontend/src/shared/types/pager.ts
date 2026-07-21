export const PAGER_DEFAULT_RESULTS_PAGE = 10;
export const PAGER_DEFAULT_RESULTS_PAGE_NO_PAGINATION = 0;

export interface PagerQuery {
  enabled: boolean;
  currentPage: number;
  resultsPage: number;
}

export interface PagerResult extends PagerQuery {
  totalPages: number;
  totalResults: number;
}

export interface Pagination {
  enabled: boolean;
  currentPage: number;
  resultsPage: number;
  totalPages: number;
  totalResults: number;
}
