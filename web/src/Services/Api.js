import axios from "axios"
import qs from "qs"

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
  GATEWAY_TIMEOUT: 504,
}

export const HTTP_VERBS = {
  DELETE: "delete",
  GET: "get",
  PATCH: "patch",
  POST: "post",
  PUT: "put",
}

const myAxios = axios.create({
  paramsSerializer: (params) => qs.stringify(params, { arrayFormat: "repeat" }),
})

// Add a request interceptor
myAxios.interceptors.request.use(
  (request) => {
    //console.log('Request:', request)
    return request
  },
  (error) => {
    console.log("REQ-Error:", error)
    return Promise.reject(error)
  }
)

// Add a response interceptor
myAxios.interceptors.response.use(
  (response) => {
    //console.log('Response:', response)
    return response
  },
  (error) => {
    // do some action
    return Promise.reject(error)
  }
)

const getHeaders = () => {
  return {
    "X-Auth-Type-Browser-UI": "1",
  }
}

const newRequest = (method, url, queryParams, data, headers) =>
  myAxios.request({
    method: method,
    url: "/api" + url,
    //url: "http://localhost:8080/api" + url,
    data: data,
    headers: { ...getHeaders(), ...headers },
    params: queryParams,
  })

export const api = {
  jobs: {
    list: () => newRequest(HTTP_VERBS.GET, "/jobs", {}, {}),
    delete: (jobId) => newRequest(HTTP_VERBS.DELETE, "/jobs/delete", { jobId }, {}),
  },
  query: {
    trigger: (data, language) =>
      newRequest(HTTP_VERBS.POST, "/query", {}, data, {
        "Content-Type": "application/" + language,
      }),
    listMetrics: (tags) => newRequest(HTTP_VERBS.GET, "/query/summary", { tags }, {}),
    listTags: () => newRequest(HTTP_VERBS.GET, "/query/tags", {}, {}),
    listTemplate: () => newRequest(HTTP_VERBS.GET, "/template/query", {}, {}),
    getTemplate: (filename) => newRequest(HTTP_VERBS.GET, "/template/query/" + filename, {}, {}),
    saveTemplate: (data) => newRequest(HTTP_VERBS.POST, "/template/query", {}, data),
  },
  generator: {
    trigger: (data, language) =>
      newRequest(HTTP_VERBS.POST, "/generator", {}, data, {
        "Content-Type": "application/" + language,
      }),
    listTemplate: () => newRequest(HTTP_VERBS.GET, "/template/generator", {}, {}),
    getTemplate: (filename) => newRequest(HTTP_VERBS.GET, "/template/generator/" + filename, {}, {}),
    saveTemplate: (data) => newRequest(HTTP_VERBS.POST, "/template/generator", {}, data),
  },
  status: {
    get: () => newRequest(HTTP_VERBS.GET, "/status", {}, {}),
  },
}
