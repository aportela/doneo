import { axiosInstance } from "../../../api/client";

import type {
  AddRequest,
  UpdateRequest,
  NoteResponse,
  SearchResponse,
} from "../types/dto";

export const noteService = {
  async addProjectNote(
    projectId: string,
    payload: AddRequest,
  ): Promise<NoteResponse> {
    const { data } = await axiosInstance.post<NoteResponse>(
      "/projects/" + projectId + "/notes",
      payload,
    );
    return data;
  },
  async updateProjectNote(
    projectId: string,
    noteId: string,
    payload: UpdateRequest,
  ): Promise<NoteResponse> {
    const { data } = await axiosInstance.put<NoteResponse>(
      "/projects/" + projectId + "/notes/" + noteId,
      payload,
    );
    return data;
  },
  async deleteProjectNote(projectId: string, noteId: string): Promise<void> {
    await axiosInstance.delete<void>(
      "/projects/" + projectId + "/notes/" + noteId,
    );
  },
  async getProjectNotes(projectId: string): Promise<SearchResponse> {
    const { data } = await axiosInstance.get<SearchResponse>(
      "/projects/" + projectId + "/notes",
    );
    return data;
  },
  async addTaskNote(
    projectId: string,
    taskId: string,
    payload: AddRequest,
  ): Promise<NoteResponse> {
    const { data } = await axiosInstance.post<NoteResponse>(
      "/projects/" + projectId + "/tasks/" + taskId + "/notes",
      payload,
    );
    return data;
  },
  async updateTaskote(
    projectId: string,
    taskId: string,
    noteId: string,
    payload: UpdateRequest,
  ): Promise<NoteResponse> {
    const { data } = await axiosInstance.put<NoteResponse>(
      "/projects/" + projectId + "/tasks/" + taskId + "/notes/" + noteId,
      payload,
    );
    return data;
  },
  async deleteTaskNote(
    projectId: string,
    taskId: string,
    noteId: string,
  ): Promise<void> {
    await axiosInstance.delete<void>(
      "/projects/" + projectId + "/tasks/" + taskId + "/notes/" + noteId,
    );
  },
  async getTaskNotes(
    projectId: string,
    taskId: string,
  ): Promise<SearchResponse> {
    const { data } = await axiosInstance.get<SearchResponse>(
      "/projects/" + projectId + "/tasks/" + taskId + "/notes",
    );
    return data;
  },
};
