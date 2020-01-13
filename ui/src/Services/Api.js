import axios from "axios";
import { t } from "typy";
import qs from "qs";

export const HTTP_CODES = {
  OK: 200,
  ACCEPTED: 202,
  NO_CONTENT: 204,
  BAD_REQUEST: 400,
  UNAUTHORIZED: 401,
  NOT_FOUND: 404,
  REQUEST_FAILED: 422,
  INTERNAL_SERVER: 500,
  SERVICE_UNAVAILABLE: 503,
  GATEWAY_TIMEOUT: 504
};

export const HTTP_VERBS = {
  DELETE: "delete",
  GET: "get",
  PATCH: "patch",
  POST: "post",
  PUT: "put"
};

const myAxios = axios.create({
  paramsSerializer: params => qs.stringify(params, { arrayFormat: "repeat" })
});

// Add a request interceptor
myAxios.interceptors.request.use(
  request => {
    //console.log('Request:', request)
    return request;
  },
  error => {
    console.log("REQ-Error:", error);
    return Promise.reject(error);
  }
);

// Add a response interceptor
myAxios.interceptors.response.use(
  response => {
    //console.log('Response:', response)
    return response;
  },
  error => {
    // do some action
    return Promise.reject(error);
  }
);

const urls = {
  jobs: {
    default: "/jobs"
  },
  queryRunner: {
    default: "/queryRunner"
  },
  spansGenerator: {
    default: "/spansGenerator"
  },
  status: {
    default: "/status"
  }
};

const url = key => {
  return t(urls, key).safeString;
};

const getHeaders = () => {
  return {
    "X-Auth-Type-Browser-UI": "1"
  };
};

const newRequest = (method, url, queryParams, data, headers) =>
  myAxios.request({
    method: method,
    url: "/api" + url,
    //url: "http://localhost:8080/api" + url,
    data: data,
    headers: { ...getHeaders(), ...headers },
    params: queryParams
  });

export const triggerQueryRunner = (data, language) => {
  return newRequest(HTTP_VERBS.POST, url("queryRunner.default"), {}, data, {
    "Content-Type": "application/" + language
  });
};

export const triggerGenerateSpans = (data, language) => {
  return newRequest(HTTP_VERBS.POST, url("spansGenerator.default"), {}, data, {
    "Content-Type": "application/" + language
  });
};

export const status = () => {
  return newRequest(HTTP_VERBS.GET, url("status.default"), {}, {});
};

export const jobs = () => {
  return newRequest(HTTP_VERBS.GET, url("jobs.default"), {}, {});
};
