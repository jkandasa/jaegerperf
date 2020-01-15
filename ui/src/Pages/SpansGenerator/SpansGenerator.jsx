import React from "react";
import { Row, Col, Input, Select, Button } from "antd";
import PageTitle from "../../Components/PageTitle/PageTitle";
import { triggerGenerateSpans } from "../../Services/Api";
import { redirect as r } from "../../Services/Routes";

const { TextArea } = Input;
const { Option } = Select;

const defaultCode = `target: "collector"
endpoint: http://jaegerqe-collector:14268/api/traces
serviceName: generated_span
numberOfDays: 10
spansPerDay: 10
spansPerSecond: 500 # maximum push span limit/sec
childDepth: 2
tags: 
  spans_generator: "jaegerperf"
  days: 10
startTime: 
`

class SpansGenerator extends React.Component {
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
    triggerGenerateSpans(this.state.codeString, this.state.language);
    r(this.props.history, "jobs");
  };

  render() {
    return (
      <React.Fragment>
        <PageTitle title={"Spans Generator"} />
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

export default SpansGenerator;
