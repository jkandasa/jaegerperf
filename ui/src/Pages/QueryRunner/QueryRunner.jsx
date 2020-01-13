import React from "react";
import { Row, Col, Input, Select, Button } from "antd";
import PageTitle from "../../Components/PageTitle/PageTitle";
import { triggerQueryRunner } from "../../Services/Api";
import { redirect as r } from "../../Services/Routes";

const { TextArea } = Input;
const { Option } = Select;

const defaultCode = `hostUrl: http://jaegerqe-query:16686
tests:
  - name: with_limit_1
    type: search
    iteration: 5
    query:
      service: generated_span
      lookback: 7d
      limit: 1

  - name: with_limit_100
    type: search
    iteration: 5
    query:
      service: generated_span
      lookback: 7d
      limit: 100
  
  - name: with_limit_500
    type: search
    iteration: 5
    query:
      service: generated_span
      lookback: 7d
      limit: 500

  - name: with_limit_1000
    type: search
    iteration: 5
    query:
      service: generated_span
      lookback: 7d
      limit: 1000

  - name: with_limit_1500
    type: search
    iteration: 5
    query:
      service: generated_span
      lookback: 7d
      limit: 1500

  - name: with_limit_2000
    type: search
    iteration: 5
    query:
      service: generated_span
      lookback: 7d
      limit: 2000

  - name: services
    type: services
    iteration: 5
`

class QueryRunner extends React.Component {
  state = {
    codeString: defaultCode,
    language: "yaml"
  };

  onChange = ({ target: { value } }) => {
    this.setState({ codeString: value });
  };
  onLanguageChange = value => {
    this.setState({ language: value });
  };
  onSubmit = () => {
    triggerQueryRunner(this.state.codeString, this.state.language);
    r(this.props.history, "jobs");
  };

  render() {
    return (
      <React.Fragment>
        <PageTitle title={"Query Runner"} />
        <Row gutter={["10", "10"]}>
          <Col>
            <span style={{ fontWeight: "600" }}>Language Selection: </span>
            <Select
              style={{ width: 200 }}
              value={this.state.language}
              onChange={this.onLanguageChange}
            >
              <Option value="yaml">YAML</Option>
              <Option value="json">JSON</Option>
            </Select>
          </Col>
          <Col>
            <TextArea
              style={{ minHeight: "50vh" }}
              value={this.state.codeString}
              onChange={this.onChange}
            />
          </Col>
          <Col>
            <Button size="large" type="primary" onClick={this.onSubmit}>
              Submit
            </Button>
            <Button size="large" style={{ marginLeft: "7px" }}>
              Cancel
            </Button>
          </Col>
        </Row>
      </React.Fragment>
    );
  }
}

export default QueryRunner;
