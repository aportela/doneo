export type EmptyResponse = {};

export type UserTimerResponse = {
  id: string;
  summary: string;
  startedAt: number;
  finishedAt: number | null;
};

export type SearchResponse = {
  userTimers: UserTimerResponse[];
};
