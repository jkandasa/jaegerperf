import React from "react"
import PageTitle from "../../Components/PageTitle/PageTitle"
import {
  Table,
  Button,
  Spin,
  message,
  Checkbox,
  Tag,
  Alert,
  Select,
  Divider,
  Empty
} from "antd"
import uuid from "uuid/v4"
import * as API from "../../Services/Api"
import LineChart from "../../Components/LineChart/LineChart"
import CustomCard from "../../Components/CustomCard/CustomCard"

import "./QueryMetrics.css"

const { Option } = Select
const CheckboxGroup = Checkbox.Group

const chartOptions = ["bar", "line"]
const metricOptions = {
  elapsed: "Elapsed Time(ms)",
  errorsCount: "Errors Count",
  errorPercentage: "Error Percentage",
  contentLength: "Content Length"
}
class Jobs extends React.Component {
  state = {
    tags: [],
    selectedTags: [],
    indeterminate: false,
    loading: true,
    loadMetrics: false,
    metricsData: {},
    errorsCount: 0,
    minContentLength: -1,
    chartType: "bar",
    metricType: "elapsed"
  }

  displayError = text => {
    message.error(text)
  }

  fetchTags = () => {
    this.setState({ loading: true })
    API.tags({})
      .then(res => {
        this.setState({ tags: res.data.sort(), loading: false })
      })
      .catch(e => {
        console.log(e)
        this.displayError(e.message ? e.message : JSON.stringify(e))
        this.setState({ tags: [], loading: false })
      })
  }

  fetchMetrics = () => {
    this.setState({ loading: true, loadMetrics: true })
    API.listQueryMetrics(this.state.selectedTags)
      .then(res => {
        let errorsCount = 0
        let minContentLength = 1000000

        const updateWarnData = d => {
          errorsCount += d["errorsCount"]
          if (d["contentLength"] < minContentLength) {
            minContentLength = d["contentLength"]
          }
        }

        // convert array to map object
        // source: [{tags: ["tag-1"], data:[{name: "query-name", elapsed: 123}]},{...}]
        // target: {"query-name":{"tag-1":{elapsed: 123}, "tag-2":{elapsed: 123}}}
        const _objMap = {}
        this.state.selectedTags.forEach(t => {
          for (let index = 0; index < res.data.length; index++) {
            const q = res.data[index]
            if (Array.from(q.tags).includes(t)) {
              const tag = this.formatTag(t) // tag.1 => tag_1
              q.data.forEach(d => {
                if (_objMap[d.name] === undefined) {
                  _objMap[d.name] = {}
                }
                _objMap[d.name][tag] = d
                // if delete the name removed the name from original reference,
                // generates undefined names
                // delete _objMap[d.name][tag]["name"]
                updateWarnData(d)
              })
              break
            }
          }
        })

        this.setState({
          metricsData: _objMap,
          loading: false,
          errorsCount: errorsCount,
          minContentLength: minContentLength
        })
      })
      .catch(e => {
        console.log(e)
        this.displayError(e.message ? e.message : JSON.stringify(e))
        this.setState({ tags: [], loading: false })
      })
  }

  componentDidMount() {
    this.fetchTags()
  }

  onChange = checkedList => {
    this.setState({
      selectedTags: checkedList,
      indeterminate:
        !!checkedList.length && checkedList.length < this.state.tags.length,
      checkAll: checkedList.length === this.state.tags.length
    })
  }

  onCheckAllChange = e => {
    this.setState({
      selectedTags: e.target.checked ? this.state.tags : [],
      indeterminate: false,
      checkAll: e.target.checked
    })
  }

  formatTag(tag) {
    return tag.replace(".", "_")
  }

  getTableColumn = () => {
    const columns = [
      {
        title: "Query",
        dataIndex: "name",
        defaultSortOrder: "ascend",
        sorter: (a, b) => {
          return a.name.localeCompare(b.name)
        }
      }
    ]
    this.state.selectedTags.forEach(t => {
      const tag = this.formatTag(t)
      columns.push({
        title: t,
        dataIndex: tag + "." + this.state.metricType,
        render: v =>
          this.state.metricType === "elapsed" ? (v ? parseInt(v / 1000) : v) : v
      })
    })
    return columns
  }

  changeChartType = type => {
    this.setState({ chartType: type })
  }

  changeMetricType = type => {
    this.setState({ metricType: type })
  }

  warningMessage = () => {
    const content = []
    if (this.state.errorsCount > 0) {
      content.push(
        <p>
          Errors found on the query runs! Number of failed queries:{" "}
          {this.state.errorsCount}
        </p>
      )
    }
    if (this.state.minContentLength <= 100) {
      content.push(
        <p>
          Some of the query response content lengths are less than 100 bytes!
        </p>
      )
    }
    if (content.length > 0) {
      return (
        <Alert
          style={{ marginTop: "10px" }}
          message="Warning"
          description={content}
          type="warning"
          showIcon
        />
      )
    }
    return null
  }

  render() {
    const data = []
    if (this.state.loadMetrics) {
      if (!this.state.loading) {
        data.push(
          <CustomCard
            key="overview"
            title="Metrics"
            body={
              <React.Fragment>
                <div>
                  Tags:{" "}
                  {this.state.selectedTags.map(t => {
                    return (
                      <Tag key={t} color="blue">
                        {t}
                      </Tag>
                    )
                  })}
                </div>
                {this.warningMessage()}
              </React.Fragment>
            }
            extra={
              <React.Fragment>
                <span style={{ marginRight: "7px" }}>Metric Type:</span>
                <Select
                  size="small"
                  style={{ width: 210 }}
                  defaultValue={this.state.metricType}
                  onChange={this.changeMetricType}
                >
                  {Object.keys(metricOptions).map(k => {
                    return (
                      <Option key={k} value={k}>
                        {metricOptions[k]}
                      </Option>
                    )
                  })}
                </Select>
              </React.Fragment>
            }
          />
        )

        // convert to data to chart format
        // source: {"query-name":{"tag-1":{elapsed: 123}, "tag-2":{elapsed: 123}}}
        // target: {labels: ["query-name", ...], datasets:[{label: "tag-1", data:[{x: "query-name", y: 123}, {...}]]}
        const category = [...Object.keys(this.state.metricsData).sort()]

        // temporary conversion
        // source: {"query-name":{"tag-1":{elapsed: 123}, "tag-2":{elapsed: 123}}}
        // target: {"tag-1": [{x: "query-name"}, {...}]}
        const _tDS = {}
        category.forEach(qn => {
          const ts = this.state.metricsData[qn]
          Object.keys(ts).forEach(t => {
            const _ts = ts[t]
            if (_tDS[t] === undefined) {
              _tDS[t] = []
            }
            _tDS[t].push({
              x: qn,
              y:
                this.state.metricType === "elapsed"
                  ? parseInt(_ts["elapsed"] / 1000)
                  : _ts[this.state.metricType]
            })
          })
        })

        // final dataset conversion
        // source: {"tag-1": [{x: "query-name"}, {...}]}
        // target: [{label: "tag-1", data:[{x: "query-name", y: 123}, {...}]]
        const _ds = []
        Object.keys(_tDS).forEach(t => {
          _ds.push({
            label: t,
            data: _tDS[t],
            type: this.state.chartType,
            fill: false,
            // yAxisID: 'left-y-axis',
            borderWidth: 2
          })
        })

        const cData = {
          labels: category,
          datasets: _ds
        }

        //console.log(JSON.stringify(cData, null, 2))
        //console.log(JSON.stringify(this.state.metricsData, null, 2))

        data.push(
          <CustomCard
            key="chart"
            title={"Chart for '" + metricOptions[this.state.metricType] + "'"}
            body={
              cData.datasets.length > 0 ? (
                <LineChart
                  data={cData}
                  yAxesLabel={metricOptions[this.state.metricType]}
                />
              ) : (
                <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} />
              )
            }
            extra={
              <React.Fragment>
                <Select
                  size="small"
                  style={{ width: 90 }}
                  defaultValue={this.state.chartType}
                  onChange={this.changeChartType}
                >
                  {chartOptions.map(v => {
                    return (
                      <Option key={v} value={v}>
                        {v.toUpperCase()}
                      </Option>
                    )
                  })}
                </Select>
              </React.Fragment>
            }
          />
        )

        // convert to table format
        // source: {"query-name":{"tag-1":{elapsed: 123}, "tag-2":{elapsed: 123}}}
        // target: [{name: "query-name", "tag-1": {"elapsed": 123}, "tag-2": {"elapsed": 123}}, {...}]
        const _tData = []
        Object.keys(this.state.metricsData).forEach(qk => {
          _tData.push({
            name: qk,
            ...this.state.metricsData[qk]
          })
        })

        //console.log(JSON.stringify(_tData, null, 2))

        data.push(
          <CustomCard
            key="table"
            size="small"
            title={"Table for '" + metricOptions[this.state.metricType] + "'"}
            style={{ marginBottom: "10px" }}
            body={
              <Table
                size="small"
                pagination={{ defaultPageSize: 15, hideOnSinglePage: true }}
                columns={this.getTableColumn()}
                dataSource={_tData}
                rowKey={uuid}
                bordered
                className="qm-table"
              />
            }
          />
        )
      }
    } else {
      data.push(
        <React.Fragment key="tags">
          <CustomCard
            size="small"
            title="Select tags to continue"
            body={
              <React.Fragment>
                <Checkbox
                  indeterminate={this.state.indeterminate}
                  onChange={this.onCheckAllChange}
                  checked={this.state.checkAll}
                >
                  Select all
                </Checkbox>
                <Divider type="vertical" />
                <CheckboxGroup
                  options={this.state.tags}
                  value={this.state.selectedTags}
                  onChange={this.onChange}
                />
              </React.Fragment>
            }
          />
          <Button
            size="default"
            shape="round"
            icon="check"
            type="primary"
            onClick={this.fetchMetrics}
            disabled={
              this.state.loading || this.state.selectedTags.length === 0
            }
          >
            Continue
          </Button>
        </React.Fragment>
      )
    }
    return (
      <React.Fragment>
        <PageTitle title={"Query Metrics"} />
        <Spin
          size="large"
          tip="Loading..."
          spinning={this.state.loading}
          delay={300}
        >
          {data}
        </Spin>
      </React.Fragment>
    )
  }
}

export default Jobs
