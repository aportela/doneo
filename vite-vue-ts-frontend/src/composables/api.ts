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
      return axiosInstance.post("/auth/renew-access-token");
    },
  },
  user: {
    search: () => axiosInstance.get("/users"),
  },
  project: {
    search: () => axiosInstance.get("/projects"),
  },
  projectTypes: {
    add: (projectType: ProjectTypeInterface) =>
      axiosInstance.post("/project-types", projectType),
    update: (projectType: ProjectTypeInterface) =>
      axiosInstance.put("/project-types/" + projectType.id, projectType),
    delete: (id: string) => axiosInstance.delete("/project-types/" + id),
    get: (id: string) => axiosInstance.get("/project-types/" + id),
    search: () => axiosInstance.get("/project-types"),
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
