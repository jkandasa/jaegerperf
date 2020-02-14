import React from "react"
import { message } from "antd"
import * as API from "../../Services/Api"
import { redirect as r } from "../../Services/Routes"
import CodeSubmitForm from "../../Components/CodeSubmitForm/CodeSubmitForm"

const defaultCode = `hostUrl: http://jaegerqe-query:16686
tags:
  - test-base-line
  - version 1.x
tests:
  - name: 1.last 12 hours limit 100
    type: search
    iteration: 10
    statusCode: 200
    queryParams:
      service: jaegerperf_generator
      limit: 100
      lookback: custom
      start: -12h
      end: 0h

  - name: 1.last 12 hours limit 1000
    type: search
    iteration: 10
    statusCode: 200
    queryParams:
      service: jaegerperf_generator
      limit: 1000
      lookback: custom
      start: -12h
      end: 0h

  - name: 1.last 12 hours limit 2000
    type: search
    iteration: 10
    statusCode: 200
    queryParams:
      service: jaegerperf_generator
      limit: 2000
      lookback: custom
      start: -12h
      end: 0h

  - name: 2.last 24 hours limit 100
    type: search
    iteration: 10
    statusCode: 200
    queryParams:
      service: jaegerperf_generator
      limit: 100
      lookback: custom
      start: -24h
      end: 0h

  - name: 2.last 24 hours limit 1000
    type: search
    iteration: 10
    statusCode: 200
    queryParams:
      service: jaegerperf_generator
      limit: 1000
      lookback: custom
      start: -24h
      end: 0h

  - name: 2.last 24 hours limit 2000
    type: search
    iteration: 10
    statusCode: 200
    queryParams:
      service: jaegerperf_generator
      limit: 2000
      lookback: custom
      start: -24h
      end: 0h

  - name: 3.last 2 days limit 100
    type: search
    iteration: 10
    statusCode: 200
    queryParams:
      service: jaegerperf_generator
      limit: 100
      lookback: custom
      start: -48h
      end: 0h

  - name: 3.last 2 days limit 1000
    type: search
    iteration: 10
    statusCode: 200
    queryParams:
      service: jaegerperf_generator
      limit: 1000
      lookback: custom
      start: -48h
      end: 0h

  - name: 3.last 2 days limit 2000
    type: search
    iteration: 10
    statusCode: 200
    queryParams:
      service: jaegerperf_generator
      limit: 2000
      lookback: custom
      start: -48h
      end: 0h

  - name: 4.last 7 days limit 100
    type: search
    iteration: 10
    statusCode: 200
    queryParams:
      service: jaegerperf_generator
      limit: 100
      lookback: custom
      start: -168h
      end: 0h

  - name: 4.last 7 days limit 1000
    type: search
    iteration: 10
    statusCode: 200
    queryParams:
      service: jaegerperf_generator
      limit: 1000
      lookback: custom
      start: -168h
      end: 0h

  - name: 4.last 7 days limit 2000
    type: search
    iteration: 10
    statusCode: 200
    queryParams:
      service: jaegerperf_generator
      limit: 2000
      lookback: custom
      start: -168h
      end: 0h
`

class QueryRunner extends React.Component {
  state = {
    codeString: defaultCode,
    language: "yaml"
  }

  onChange = ({ target: { value } }) => {
    this.setState({ codeString: value })
  }

  onLanguageChange = value => {
    this.setState({ language: value })
  }

  displayError = text => {
    message.error(text)
  }

  displayInfo = text => {
    message.info(text)
  }

  onSubmit = () => {
    API.triggerQueryRunner(this.state.codeString, this.state.language)
      .then(res => {
        this.displayInfo(JSON.stringify(res.data))
        r(this.props.history, "jobs")
      })
      .catch(e => {
        this.displayError(e.message ? e.message : JSON.stringify(e))
      })
  }

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
    )
  }
}

export default QueryRunner
