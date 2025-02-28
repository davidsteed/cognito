import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from "axios";

export const handleAxiosRequest = <T>(
  axiosInstance: AxiosInstance,
  config: AxiosRequestConfig,
): Promise<T> => {
  const source = axios.CancelToken.source();
  const promise = axiosInstance({
    ...config,
    cancelToken: source.token,
  }).then((res: AxiosResponse<T>) => {
    const { data } = res;
    return data;
  });

  // eslint-disable-next-line
  // @ts-expect-error Required by Orval https://orval.dev/guides/custom-axios
  promise.cancel = () => {
    source.cancel("Query was cancelled by React Query");
  };

  return promise;
};

export interface AxiosError<T> {
  config: AxiosRequestConfig;
  code?: string;
  request?: unknown;
  response?: AxiosResponse<T>;
}
