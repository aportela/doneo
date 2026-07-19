export interface PagerQuery {
  currentPage: number;
  resultsPage: number;
}

export interface PagerResult extends PagerQuery {
  enabled: boolean;
  totalPages: number;
  totalResults: number;
}
