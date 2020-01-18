import React from "react";
import { message } from "antd";
import * as API from "../../Services/Api";
import { redirect as r } from "../../Services/Routes";
import CodeSubmitForm from "../../Components/CodeSubmitForm/CodeSubmitForm"


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
`;

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
    API.triggerGenerateSpans(this.state.codeString, this.state.language)
      .then(res => {
        this.displayInfo(JSON.stringify(res.data));
        r(this.props.history, "jobs");
      })
      .catch(e => {
        this.displayError(e.message ? e.message : JSON.stringify(e));
      });
  };

  displayError = text => {
    message.error(text);
  };

  displayInfo = text => {
    message.info(text);
  };

  render() {
    return (
      <CodeSubmitForm 
        title="Spans Generator"
        language={this.state.language}
        onLanguageChange={this.onLanguageChange}
        codeString={this.state.codeString}
        onCodeChange={this.onChange}
        onSubmit={this.onSubmit}
      />
    );
  }
}

export default SpansGenerator;
