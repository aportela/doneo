import type { PagerQuery, PagerResult } from "../types/pager";

export class ClientPager implements PagerResult {
  enabled: boolean;
  currentPage: number;
  resultsPage: number;
  totalPages: number;
  totalResults: number;

  constructor(enabled: boolean, currentPage: number, resultsPage: number) {
    this.enabled = enabled;
    this.currentPage = currentPage;
    this.resultsPage = resultsPage;
    this.totalPages = 0;
    this.totalResults = 0;
  }

  toPagerQuery = (): PagerQuery => {
    return {
      currentPage: this.currentPage,
      resultsPage: this.resultsPage,
    };
  };

  fromPagerResult = (pagerResult: PagerResult) => {
    this.enabled = pagerResult.enabled;
    this.currentPage = pagerResult.currentPage;
    this.resultsPage = pagerResult.resultsPage;
    this.totalPages = pagerResult.totalPages;
    this.totalResults = pagerResult.totalResults;
  };
}
