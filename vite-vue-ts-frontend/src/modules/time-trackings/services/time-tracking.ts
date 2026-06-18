import { axiosInstance } from "../../../api/client";

import type {
  AddRequest,
  UpdateRequest,
  TimeTrackingResponse,
  SearchResponse,
} from "../types/dto";

export const noteService = {
  async addTaskTimeTracking(
    projectId: string,
    taskId: string,
    payload: AddRequest,
  ): Promise<TimeTrackingResponse> {
    const { data } = await axiosInstance.post<TimeTrackingResponse>(
      "/projects/" + projectId + "/tasks/" + taskId + "/time_trackings",
      payload,
    );
    return data;
  },
  async updateTaskTimeTracking(
    projectId: string,
    taskId: string,
    timeTrackingId: string,
    payload: UpdateRequest,
  ): Promise<TimeTrackingResponse> {
    const { data } = await axiosInstance.put<TimeTrackingResponse>(
      "/projects/" +
        projectId +
        "/tasks/" +
        taskId +
        "/time_trackings/" +
        timeTrackingId,
      payload,
    );
    return data;
  },
  async deleteTaskTimeTracking(
    projectId: string,
    taskId: string,
    timeTrackingId: string,
  ): Promise<void> {
    await axiosInstance.delete<void>(
      "/projects/" +
        projectId +
        "/tasks/" +
        taskId +
        "/time_trackings/" +
        timeTrackingId,
    );
  },
  async getTaskTimeTracking(
    projectId: string,
    taskId: string,
    timeTrackingId: string,
  ): Promise<TimeTrackingResponse> {
    const { data } = await axiosInstance.get<TimeTrackingResponse>(
      "/projects/" +
        projectId +
        "/tasks/" +
        taskId +
        "/time_trackings/" +
        timeTrackingId,
    );
    return data;
  },
  async getTaskTimeTrackings(
    projectId: string,
    taskId: string,
  ): Promise<SearchResponse> {
    const { data } = await axiosInstance.get<SearchResponse>(
      "/projects/" + projectId + "/tasks/" + taskId + "/time_trackings",
    );
    return data;
  },
};
