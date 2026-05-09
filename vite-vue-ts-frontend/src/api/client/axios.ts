import axios from "axios";
import {
  SERVER_API_BASE_PATH,
  SERVER_API_DEFAULT_TIMEOUT,
} from "../../constants";

export const axiosInstance = axios.create({
  baseURL: SERVER_API_BASE_PATH,
  timeout: SERVER_API_DEFAULT_TIMEOUT,
  withCredentials: true,
});
