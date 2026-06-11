export type EmptyResponse = {};

export type TimerResponse = {
  id: string;
  summary: string;
  startedAt: number;
  finishedAt: number | null;
};

export type SearchResponse = {
  timers: TimerResponse[];
};
