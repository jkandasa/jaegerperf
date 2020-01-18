import React from "react";
import {message } from "antd";
import * as API from "../../Services/Api";
import { redirect as r } from "../../Services/Routes";
import CodeSubmitForm from "../../Components/CodeSubmitForm/CodeSubmitForm"

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
`;

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

  displayError = text => {
    message.error(text);
  };

  displayInfo = text => {
    message.info(text);
  };

  onSubmit = () => {
    API.triggerQueryRunner(this.state.codeString, this.state.language)
      .then(res => {
        this.displayInfo(JSON.stringify(res.data));
        r(this.props.history, "jobs");
      })
      .catch(e => {
        this.displayError(e.message ? e.message : JSON.stringify(e));
      });
  };

  render() {
    return (
      <CodeSubmitForm 
      title="Query Runner"
      language={this.state.language}
      onLanguageChange={this.onLanguageChange}
      codeString={this.state.codeString}
      onCodeChange={this.onChange}
      onSubmit={this.onSubmit}
    />
    );
  }
}

export default QueryRunner;
