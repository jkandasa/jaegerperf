import JobsPage from "../Pages/Jobs/Jobs"
import SpansGeneratorPage from "../Pages/SpansGenerator/SpansGenerator"
import QueryRunnerPage from "../Pages/QueryRunner/QueryRunner"
import QueryMetricsPage from "../Pages/QueryMetrics/QueryMetrics"
import { t } from "typy"

const routes = [
  {
    id: "sgPage",
    title: "Spans Generator",
    to: "/spans-generator",
    component: SpansGeneratorPage,
  },
  {
    id: "qrPage",
    title: "Query Runner",
    to: "/query-runner",
    component: QueryRunnerPage,
  },
  {
    id: "jPage",
    title: "Jobs",
    to: "/jobs",
    component: JobsPage,
  },
  {
    id: "queryMetrics",
    title: "Query Metrics",
    to: "/query-metrics",
    component: QueryMetricsPage,
  },
]

const hiddenRoutes = [
  {
    to: "/",
    component: "",
  },
]

const routeMap = {
  home: "/",
  dashboard: "/",
  spansGenerator: "/spans-generator",
  queryRunner: "/query-runner",
  jobs: "/jobs",
}

const redirect = (history, name = "home", urlParams = {}) => {
  const to = t(routeMap, name).safeString
  if (to) {
    let finalPath = to
    Object.keys(urlParams).forEach((key) => {
      finalPath = finalPath.replace(":" + key, urlParams[key])
    })
    history.push(finalPath)
  }
}

export { routes, hiddenRoutes, routeMap, redirect }
