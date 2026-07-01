import { type SearchResponseResultItem as SearchResponseResultItemDTO } from "../types/dto";

export class SearchResultItem {
  constructor(_data?: SearchResponseResultItemDTO) {}

  toDTO(): SearchResponseResultItemDTO {
    return {};
  }
}
