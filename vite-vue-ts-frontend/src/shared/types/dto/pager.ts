export type PagerRequest = {
  currentPage: number;
  resultsPage: number;
};

export type PagerResponse = {
  enabled: boolean;
  currentPage: number;
  resultsPage: number;
  totalPages: number;
  totalResults: number;
};
