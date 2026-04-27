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
  },
  user: {
    search: () => axiosInstance.get("/users"),
  },
  project: {
    search: () => axiosInstance.get("/projects"),
  },
};

export { api };
