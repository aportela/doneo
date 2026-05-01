import { axiosInstance } from "./axios";

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
  user: {
    search: () => axiosInstance.get("/users"),
  },
  project: {
    search: () => axiosInstance.get("/projects"),
  },
  projectTypes: {
    search: () => axiosInstance.get("/project_types"),
  },
  projectStatuses: {
    search: () => axiosInstance.get("/project_statuses"),
  },
  projectPriorities: {
    search: () => axiosInstance.get("/project_priorities"),
  },
};

export { api };
