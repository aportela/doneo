import { axiosInstance } from "./axios";
import { setupAxiosWithStore } from "./setupAxiosWithStore";
import { useSessionStore } from "../../stores/session";

setupAxiosWithStore(useSessionStore());

export { axiosInstance };
