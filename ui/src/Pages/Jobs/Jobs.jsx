import React from "react";
import PageTitle from "../../Components/PageTitle/PageTitle";
import { Table, Button, Spin } from "antd";
import uuid from "uuid/v4";
import { jobs } from "../../Services/Api";
import moment from "moment";
import { t } from "typy";

const columns = [
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
  }
];

const queryRunnerColumns = [
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
  }
];

class Jobs extends React.Component {
  state = {
    data: [],
    loading: true
  };

  fetchData = () => {
    this.setState({ loading: true });
    jobs({})
      .then(res => {
        this.setState({ data: res.data, loading: false });
      })
      .catch(e => {
        console.log(e);
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
            columns={queryRunnerColumns}
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
        <pre style={{ margin: 0 }}>{JSON.stringify(record, null, 2)}</pre>
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
          size="large"
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
            columns={columns}
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
