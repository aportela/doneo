import { axiosInstance } from "./axios";

import type { UserInterface } from "../types/models/user";
import type { ProjectTypeInterface } from "../types/models/projectType";
import type { ProjectStatusInterface } from "../types/models/projectStatus";
import type { ProjectPriorityInterface } from "../types/models/projectPriority";

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
    add: (user: UserInterface) => axiosInstance.post("/users", user),
    update: (user: UserInterface) =>
      axiosInstance.put("/users/" + user.id, user),
    delete: (id: string) => axiosInstance.delete("/users/" + id),
    unDelete: (id: string) =>
      axiosInstance.patch("/users/" + id, { deletedAt: null }),
    get: (id: string) => axiosInstance.get("/users/" + id),
    search: () => axiosInstance.get("/users"),
  },
  project: {
    search: () => axiosInstance.get("/projects"),
  },
  projectType: {
    add: (projectType: ProjectTypeInterface) =>
      axiosInstance.post("/project-types", projectType),
    update: (projectType: ProjectTypeInterface) =>
      axiosInstance.put("/project-types/" + projectType.id, projectType),
    delete: (id: string) => axiosInstance.delete("/project-types/" + id),
    get: (id: string) => axiosInstance.get("/project-types/" + id),
    search: () => axiosInstance.get("/project-types"),
  },
  projectStatus: {
    add: (projectStatus: ProjectStatusInterface) =>
      axiosInstance.post("/project-statuses", projectStatus),
    update: (projectStatus: ProjectStatusInterface) =>
      axiosInstance.put("/project-statuses/" + projectStatus.id, projectStatus),
    patchIndex: (id: string, newIndex: number) => {
      const params = {
        index: newIndex,
      };
      return axiosInstance.patch("/project-statuses/" + id, params);
    },
    delete: (id: string) => axiosInstance.delete("/project-statuses/" + id),
    get: (id: string) => axiosInstance.get("/project-statuses/" + id),
    search: () => axiosInstance.get("/project-statuses"),
  },
  projectPriority: {
    add: (projectPriority: ProjectPriorityInterface) =>
      axiosInstance.post("/project-priorities", projectPriority),
    update: (projectPriority: ProjectPriorityInterface) =>
      axiosInstance.put(
        "/project-priorities/" + projectPriority.id,
        projectPriority,
      ),
    patchIndex: (id: string, newIndex: number) => {
      const params = {
        index: newIndex,
      };
      return axiosInstance.patch("/project-priorities/" + id, params);
    },
    delete: (id: string) => axiosInstance.delete("/project-priorities/" + id),
    get: (id: string) => axiosInstance.get("/project-priorities/" + id),
    search: () => axiosInstance.get("/project-priorities"),
  },
};

export { api };
