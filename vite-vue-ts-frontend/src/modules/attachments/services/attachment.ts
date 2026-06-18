import { axiosInstance } from "../../../api/client";

import type { SearchResponse } from "../types/dto";

export const attachmentService = {
  async deleteProjectAttachment(
    projectId: string,
    attachmentId: string,
  ): Promise<void> {
    await axiosInstance.delete<void>(
      "/projects/" + projectId + "/attachments/" + attachmentId,
    );
  },
  async getProjectAttachments(projectId: string): Promise<SearchResponse> {
    const { data } = await axiosInstance.get<SearchResponse>(
      "/projects/" + projectId + "/attachments",
    );
    return data;
  },
  async deleteTaskAttachment(
    projectId: string,
    taskId: string,
    attachmentId: string,
  ): Promise<void> {
    await axiosInstance.delete<void>(
      "/projects/" +
        projectId +
        "/tasks/" +
        taskId +
        "/attachments/" +
        attachmentId,
    );
  },
  async getTaskAttachments(
    projectId: string,
    taskId: string,
  ): Promise<SearchResponse> {
    const { data } = await axiosInstance.get<SearchResponse>(
      "/projects/" + projectId + "/tasks/" + taskId + "/attachments",
    );
    return data;
  },
};
