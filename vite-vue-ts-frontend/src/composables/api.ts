import { axiosInstance } from "./axios";

const api = {
  user: {
    search: () => axiosInstance.get("/users"),
  },
  project: {
    search: () => axiosInstance.get("/projects"),
  },
};

export { api };
