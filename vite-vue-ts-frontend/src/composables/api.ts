import { axiosInstance } from "./axios";

import type { ProjectTypeInterface } from "../types/models/projectType";
const api = {
  auth: {
    signIn: (email: string, password: string) => {
      const params = {
        email: email,
        password: password,
      };
      return axiosInstance.post("/auth/signin", params);
    },
    signUp: (name: string, email: string, password: string) => {
      const params = {
        name: name,
        email: email,
        password: password,
      };
      return axiosInstance.post("/auth/signup", params);
    },
    signOut: () => axiosInstance.post("/auth/signout"),
    renewAccessToken: function () {
      return axiosInstance.post("/auth/renew_access_token");
    },
  },
  workspace: {
    search: () => axiosInstance.get("/workspaces"),
  },
  user: {
    search: () => axiosInstance.get("/users"),
  },
  project: {
    search: () => axiosInstance.get("/projects"),
  },
  projectTypes: {
    add: (workspaceId: string, projectType: ProjectTypeInterface) =>
      axiosInstance.post(
        "/workspaces/" + workspaceId + "/project-types",
        projectType,
      ),
    update: (workspaceId: string, projectType: ProjectTypeInterface) =>
      axiosInstance.put(
        "/workspaces/" + workspaceId + "/project-types/" + projectType.id,
        projectType,
      ),
    delete: (workspaceId: string, id: string) =>
      axiosInstance.delete(
        "/workspaces/" + workspaceId + "/project-types/" + id,
      ),
    get: (workspaceId: string, id: string) =>
      axiosInstance.get("/workspaces/" + workspaceId + "/project-types/" + id),
    search: (workspaceId: string) =>
      axiosInstance.get("/workspaces/" + workspaceId + "/project-types"),
  },
  projectStatuses: {
    search: (workspaceId: string) =>
      axiosInstance.get("/project_statuses?workspace_id=" + workspaceId),
  },
  projectPriorities: {
    search: (workspaceId: string) =>
      axiosInstance.get("/project_priorities?workspace_id=" + workspaceId),
  },
};

export { api };
