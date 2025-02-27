import Axios, { AxiosRequestConfig, AxiosResponse, AxiosInstance } from "axios";

const baseApiUrl = "localhost";
export const AXIOS_INSTANCE = Axios.create({ baseURL: baseApiUrl });

const handleAxiosRequest = <T>(
  axiosInstance: AxiosInstance,
  config: AxiosRequestConfig,
): Promise<T> => {
  const source = Axios.CancelToken.source();
  const promise = axiosInstance({
    ...config,
    cancelToken: source.token,
  }).then((res: AxiosResponse<T>) => {
    const { data } = res;
    return data;
  });

  // // eslint-disable-next-line
  // @ts-expect-error Required by Orval https://orval.dev/guides/custom-axios
  promise.cancel = () => {
    source.cancel("Query was cancelled by React Query");
  };
  return promise;
};

export const customInstance = <T>(config: AxiosRequestConfig): Promise<T> => {
  return handleAxiosRequest(AXIOS_INSTANCE, config);
};

interface AxiosError<T> {
  config: AxiosRequestConfig;
  code?: string;
  request?: unknown;
  response?: AxiosResponse<T>;
}

export type ErrorType<Error> = AxiosError<Error>;
