import axios from "axios";
import { fetchAuthSession } from "@aws-amplify/core";

const customAxiosInstance = axios.create();

customAxiosInstance.interceptors.request.use(async (config) => {
  try {
    const session = await fetchAuthSession();
    const token = session.tokens?.accessToken?.toString();

    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
  } catch (error) {
    console.warn("No auth token available", error);
  }

  return config;
});

export default customAxiosInstance;
