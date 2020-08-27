import React from "react"
import PageTitle from "../../Components/PageTitle/PageTitle"
import { Table, Button, Spin, Divider, Tooltip, Modal } from "antd"
import { api as API } from "../../Services/Api"
import moment from "moment"
import { t } from "typy"
import { DeleteOutlined, DownloadOutlined, SyncOutlined } from "@ant-design/icons"
import { errorMsg } from "../../Components/Message/Message"

const { confirm } = Modal

class Jobs extends React.Component {
  state = {
    data: [],
    loading: true,
  }

  columns = [
    {
      title: "ID",
      dataIndex: "id",
    },
    {
      title: "Type",
      dataIndex: "type",
    },
    {
      title: "Is Running",
      dataIndex: "isRunning",
      render: (v) => (v ? "Running" : "Completed"),
    },
    {
      title: "Modified Time",
      dataIndex: "modifiedTime",
      render: (v) => moment(v).format("DD-MMM-YYYY, HH:mm:ss"),
      defaultSortOrder: "descend",
      sorter: (a, b) => moment(a.modifiedTime) - moment(b.modifiedTime),
    },
    {
      title: "Actions",
      key: "actions",
      render: (record) => {
        return (
          <span>
            <Tooltip placement="left" title="Delete">
              <DeleteOutlined
                style={{ fontSize: "16px", color: "#08c" }}
                onClick={() => this.deleteConfirmation(record)}
              />
            </Tooltip>
            <Divider type="vertical" />
            <Tooltip placement="left" title="Download">
              <DownloadOutlined
                type="download"
                style={{ fontSize: "16px", color: "#08c" }}
                onClick={() => this.saveRecord(record)}
              />
            </Tooltip>
          </span>
        )
      },
    },
  ]

  queryRunnerColumns = [
    {
      title: "Query",
      dataIndex: "name",
      defaultSortOrder: "ascend",
      sorter: (a, b) => {
        return a.name.localeCompare(b.name)
      },
    },
    {
      title: "Elapsed (ms)",
      dataIndex: "elapsed",
      render: (v) => parseInt(v / 1000),
    },
    {
      title: "Samples",
      dataIndex: "samples",
    },
    {
      title: "Errors Count",
      dataIndex: "errorsCount",
    },
    {
      title: "Error Percentage",
      dataIndex: "errorPercentage",
    },
    {
      title: "Content Length",
      dataIndex: "contentLength",
    },
  ]

  deleteRecord = (r) => {
    API.jobs
      .delete(r.id)
      .then(() => {
        // deleted successfully
        this.fetchData()
      })
      .catch((e) => {
        errorMsg(e.message ? e.message : JSON.stringify(e))
      })
  }

  deleteConfirmation = (r) => {
    const message = `ID\t: ${r.id}\nType\t: ${r.type}\nStatus\t: ${r.isRunning ? "Running" : "Completed"}`
    confirm({
      width: 520,
      title: "Do you want to delete this item?",
      content: <pre>{message}</pre>,
      onOk: () => this.deleteRecord(r),
      onCancel() {},
    })
  }

  saveRecord = (r) => {
    const blob = new Blob([JSON.stringify(r, "", " ")], {
      type: "application/json",
    })
    const url = URL.createObjectURL(blob)

    const a = document.createElement("a")
    a.href = url
    a.download = r.id + ".json"
    const click = new MouseEvent("click")

    // Push the download operation on the next tick
    requestAnimationFrame(() => {
      a.dispatchEvent(click)
    })
  }

  fetchData = () => {
    this.setState({ loading: true })
    API.jobs
      .list()
      .then((res) => {
        this.setState({ data: res.data, loading: false })
      })
      .catch((e) => {
        errorMsg(e.message ? e.message : JSON.stringify(e))
        this.setState({ data: [], loading: false })
      })
  }

  componentDidMount() {
    this.fetchData()
  }

  expandRowFn = (record) => {
    let queryRunnerTable
    if (record.type === "query_runner" && !record.isRunning) {
      if (t(record, "data.report.summary").isDefined) {
        const subData = t(record, "data.report.summary").safeArray
        queryRunnerTable = (
          <Table
            style={{ marginBottom: "7px" }}
            columns={this.queryRunnerColumns}
            dataSource={subData}
            rowKey="name"
            pagination={false}
            bordered
            size="small"
          />
        )
      }
    }
    return (
      <React.Fragment>
        {queryRunnerTable}
        <pre style={{ margin: 0 }}>{JSON.stringify(record, null, 2)}</pre>
      </React.Fragment>
    )
  }

  render() {
    return (
      <React.Fragment>
        <PageTitle title={"Jobs"} />
        <Button
          disabled={this.state.loading}
          onClick={this.fetchData}
          type="primary"
          style={{ marginBottom: "7px" }}
          size="default"
          shape="round"
        >
          <SyncOutlined /> Refresh
        </Button>
        <Spin size="large" tip="Loading..." spinning={this.state.loading} delay={300}>
          <Table
            columns={this.columns}
            dataSource={this.state.data}
            rowKey="id"
            expandedRowRender={this.expandRowFn}
            bordered
            size="small"
          />
        </Spin>
      </React.Fragment>
    )
  }
}

export default Jobs
