import React from "react";
import PageTitle from "../../Components/PageTitle/PageTitle";
import {
  Table,
  Button,
  Spin,
  Icon,
  Divider,
  Tooltip,
  Modal,
  message
} from "antd";
import uuid from "uuid/v4";
import * as API from "../../Services/Api";
import moment from "moment";
import { t } from "typy";
import Highlighter from "../../Components/SyntaxHighlighter/SyntaxHighlighter";

const { confirm } = Modal;

class Jobs extends React.Component {
  state = {
    data: [],
    loading: true
  };

  columns = [
    {
      title: "ID",
      dataIndex: "id"
    },
    {
      title: "Type",
      dataIndex: "type"
    },
    {
      title: "Status",
      dataIndex: "data.isRunning",
      render: v => (v ? "Running" : "Completed")
    },
    {
      title: "Modified Time",
      dataIndex: "modifiedTime",
      render: v => moment(v).format("DD-MMM-YYYY, HH:mm:ss"),
      defaultSortOrder: "descend",
      sorter: (a, b) => moment(a.modifiedTime) - moment(b.modifiedTime)
    },
    {
      title: "Actions",
      key: "actions",
      render: record => {
        return (
          <span>
            <Tooltip placement="left" title="Delete">
              <Icon
                type="delete"
                style={{ fontSize: "16px", color: "#08c" }}
                onClick={() => this.deleteConfirmation(record)}
              />
            </Tooltip>
            <Divider type="vertical" />
            <Tooltip placement="left" title="Download">
              <Icon
                type="download"
                style={{ fontSize: "16px", color: "#08c" }}
                onClick={() => this.saveRecord(record)}
              />
            </Tooltip>
          </span>
        );
      }
    }
  ];

  queryRunnerColumns = [
    {
      title: "Query",
      dataIndex: "name",
      defaultSortOrder: "ascend",
      sorter: (a, b) => {
        return a.name.localeCompare(b.name);
      }
    },
    {
      title: "Elapsed (ms)",
      dataIndex: "elapsed",
      render: v => parseInt(v / 1000)
    },
    {
      title: "Samples",
      dataIndex: "samples"
    },
    {
      title: "Errors Count",
      dataIndex: "errorsCount"
    },
    {
      title: "Error Percentage",
      dataIndex: "errorPercentage"
    },
    {
      title: "Content Length",
      dataIndex: "contentLength"
    }
  ];

  displayError = text => {
    message.error(text);
  };

  deleteRecord = r => {
    API.deleteJob(r.id)
      .then(() => {
        // deleted successfully
        this.fetchData();
      })
      .catch(e => {
        this.displayError(e.message ? e.message : JSON.stringify(e));
      });
  };

  deleteConfirmation = r => {
    const message = `ID\t: ${r.id}\nType\t: ${r.type}\nStatus\t: ${
      r.isRunning ? "Running" : "Completed"
    }`;
    confirm({
      width: 520,
      title: "Do you want to delete this item?",
      content: <pre>{message}</pre>,
      onOk: () => this.deleteRecord(r),
      onCancel() {}
    });
  };

  saveRecord = r => {
    const blob = new Blob([JSON.stringify(r, "", " ")], {
      type: "application/json"
    });
    const url = URL.createObjectURL(blob);

    const a = document.createElement("a");
    a.href = url;
    a.download = r.id + ".json";
    const click = new MouseEvent("click");

    // Push the download operation on the next tick
    requestAnimationFrame(() => {
      a.dispatchEvent(click);
    });
  };

  fetchData = () => {
    this.setState({ loading: true });
    API.jobs({})
      .then(res => {
        this.setState({ data: res.data, loading: false });
      })
      .catch(e => {
        this.displayError(e.message ? e.message : JSON.stringify(e));
        this.setState({ data: [], loading: false });
      });
  };

  componentDidMount() {
    this.fetchData();
  }

  expandRowFn = record => {
    let queryRunnerTable;
    if (record.type === "QueryRunner" && !record.isRunning) {
      if (t(record, "data.data.metrics.summary").isDefined) {
        const subData = t(record, "data.data.metrics.summary").safeArray;
        queryRunnerTable = (
          <Table
            style={{ marginBottom: "7px" }}
            columns={this.queryRunnerColumns}
            dataSource={subData}
            rowKey={uuid}
            pagination={false}
            bordered
            size="small"
          />
        );
      }
    }
    return (
      <React.Fragment>
        {queryRunnerTable}
        <Highlighter code={JSON.stringify(record, null, 2)} language="json" />
      </React.Fragment>
    );
  };

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
          icon="sync"
          shape="round"
        >
          Refresh
        </Button>
        <Spin
          size="large"
          tip="Loading..."
          spinning={this.state.loading}
          delay={300}
        >
          <Table
            columns={this.columns}
            dataSource={this.state.data}
            rowKey={uuid}
            expandedRowRender={this.expandRowFn}
            bordered
          />
        </Spin>
      </React.Fragment>
    );
  }
}

export default Jobs;
