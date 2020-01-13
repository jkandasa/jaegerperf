import React from "react";
import PageTitle from "../../Components/PageTitle/PageTitle";
import { Table, Button } from "antd";
import uuid from "uuid/v4";
import { jobs } from "../../Services/Api";
import moment from 'moment'



const columns = [
  {
    title: "ID",
    dataIndex: "id",
  },
  {
    title: "Type",
    dataIndex: "type",
  },
  {
    title: "Status",
    dataIndex: "data.isRunning",
    render: v => (v ? "Running" : "Completed"),
  },
  {
    title: "Modified Time",
    dataIndex: "modifiedTime",
    render: v => moment(v).format('DD-MMM-YYYY, HH:mm:ss'),
    defaultSortOrder: "descend",
    sorter: (a, b) => moment(a.modifiedTime) -moment(b.modifiedTime)
  }
];

class Jobs extends React.Component {
  state = {
    data: []
  };

  fetchData = () => {
    jobs({})
      .then(res => {
        this.setState({ data: res.data });
      })
      .catch(e => {
        console.log(e);
      });
  };

  componentDidMount() {
    this.fetchData();
  }

  expandRowFn = record => {
    return <pre style={{ margin: 0 }}>{JSON.stringify(record, null, 2)}</pre>;
  };

  render() {
    return (
      <React.Fragment>
        <PageTitle title={"Jobs"} />
        <Button
          onClick={this.fetchData}
          type="primary"
          style={{ marginBottom: "7px" }}
          size="large"
        >
          Refresh
        </Button>

        <Table
          columns={columns}
          dataSource={this.state.data}
          rowKey={uuid}
          expandedRowRender={this.expandRowFn}
          bordered
        />
      </React.Fragment>
    );
  }
}

export default Jobs;
