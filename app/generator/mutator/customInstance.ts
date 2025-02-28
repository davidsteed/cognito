import axios, { AxiosRequestConfig, AxiosResponse } from "axios";
import { fetchAuthSession } from "@aws-amplify/core";
import { AxiosError } from "@/app/utils";

const customInstance = axios.create({
  baseURL: "https://api.testawsreact.com/demo-api",
});

// Add auth token to every request
customInstance.interceptors.request.use(async (config) => {
  try {
    const session = await fetchAuthSession();
    const token = session.tokens?.idToken?.toString();

    if (token) {
      config.headers = config.headers || {};
      config.headers.Authorization = `Bearer ${token}`;
    }
  } catch (error) {
    console.warn("No auth token available", error);
  }

  return config;
});

export default customInstance;

export type ErrorType<E = unknown> = AxiosError<E>;

// ✅ Ensure `customMutator` is strictly typed
export const customMutator: <T>(
  config: AxiosRequestConfig,
  options?: Record<string, unknown>,
) => Promise<T> = async <T>(
  config: AxiosRequestConfig,
  options?: Record<string, unknown>,
): Promise<T> => {
  const response: AxiosResponse<T> = await customInstance.request<T>({
    ...config,
    ...(options as object), // Prevents `any` inference
  });

  return response.data;
};
