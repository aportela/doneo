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
    add: (projectType: ProjectTypeInterface) =>
      axiosInstance.post("/project_types", projectType),
    update: (projectType: ProjectTypeInterface) =>
      axiosInstance.put("/project_types/" + projectType.id, projectType),
    delete: (id: string) => axiosInstance.delete("/project_types/" + id),
    get: (id: string) => axiosInstance.get("/project_types/" + id),
    search: (workspaceId: string) =>
      axiosInstance.get("/project_types?workspaceid=" + workspaceId),
  },
  projectStatuses: {
    search: (workspaceId: string) =>
      axiosInstance.get("/project_statuses?workspaceid=" + workspaceId),
  },
  projectPriorities: {
    search: (workspaceId: string) =>
      axiosInstance.get("/project_priorities?workspaceid=" + workspaceId),
  },
};

export { api };
