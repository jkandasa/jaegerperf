import React from "react"
import { Input, Select, Row, Col, Button, Form, Space } from "antd"
import { redirect as r } from "../../Services/Routes"
import PageTitle from "../PageTitle/PageTitle"
import Editor from "@monaco-editor/react"
import { monaco } from "@monaco-editor/react"
import { CheckOutlined, CloseOutlined } from "@ant-design/icons"
import { withRouter } from "react-router"
import { infoMsg, errorMsg } from "../../Components/Message/Message"

monaco
  .init()
  .then((monaco) => {
    monaco.editor.defineTheme("console", {
      base: "vs-dark",
      inherit: true,
      rules: [
        { token: "number", foreground: "ace12e" },
        { token: "type", foreground: "73bcf7" },
        { token: "string", foreground: "f0ab00" },
        { token: "keyword", foreground: "cbc0ff" },
      ],
      colors: {
        "editor.background": "#151515",
        "editorGutter.background": "#292e34",
        "editorLineNumber.activeForeground": "#fff",
        "editorLineNumber.foreground": "#f0f0f0",
      },
    })
  })
  .catch((error) => console.error("An error occurred during initialization of Monaco: ", error))

const { Option } = Select

const options = {
  selectOnLineNumbers: true,
  scrollBeyondLastLine: false,
  contextmenu: true,
  autoIndent: "full",
  cursorBlinking: "phase",
  smoothScrolling: true,
  tabSize: 2,
  fontSize: 15,
}

class SubmitTemplate extends React.Component {
  state = {
    loading: true,
    templates: [],
    templateName: "",
    data: "",
    isEditorReady: false,
  }
  valueGetter = React.createRef()

  componentDidMount() {
    this.props
      .listTemplateFn()
      .then((res) => {
        const _templates = res.data.map((t) => t.name)
        this.setState({ templates: _templates, loading: false })
      })
      .catch((e) => {
        errorMsg(e.message ? e.message : JSON.stringify(e))
        this.setState({ loading: false })
      })
  }

  onChange = ({ target: { value } }) => {
    this.setState({ data: value })
  }

  onSubmit = (tName, execute) => {
    const code = this.valueGetter.current()
    this.props
      .saveTemplateFn({ name: tName, data: code })
      .then(() => {
        if (execute) {
          this.props
            .triggerFn(code, "yaml")
            .then((res1) => {
              infoMsg(JSON.stringify(res1.data), 10)
              r(this.props.history, "jobs")
            })
            .catch((er) => {
              console.log(er)
              errorMsg(er.message ? er.message : JSON.stringify(er))
            })
        }
      })
      .catch((e) => {
        console.log(e)
        errorMsg(e.message ? e.message : JSON.stringify(e))
      })
  }

  onTemplateChange = (templateName) => {
    this.props
      .getTemplateFn(templateName)
      .then((res) => {
        const f = res.data
        this.setState({ templateName: f.name, data: f.data, loading: false })
      })
      .catch((e) => {
        errorMsg(e.message ? e.message : JSON.stringify(e))
        this.setState({ loading: false })
      })
  }

  onTemplateNameChange = (name) => {
    this.setState({ templateName: name })
  }

  handleEditorDidMount = (_valueGetter) => {
    this.valueGetter.current = _valueGetter
    this.setState({ isEditorReady: true })
  }

  render() {
    const tOptions = this.state.templates.map((t) => (
      <Option key={t} value={t}>
        {t}
      </Option>
    ))

    return (
      <>
        <PageTitle title={this.props.title} />
        <Form labelCol={{ span: 3, offset: 0 }} wrapperCol={{ span: 21, offset: 0 }} size="middle">
          <Row>
            <Col span={24}>
              <Form.Item
                wrapperCol={{ span: 9, offset: 0 }}
                label="Source"
                colon={false}
                labelAlign="left"
                style={{ marginBottom: "5px" }}
              >
                <Select style={{ width: "100%" }} onChange={this.onTemplateChange}>
                  {tOptions}
                </Select>
              </Form.Item>
            </Col>
            <Col span={24}>
              <Form.Item label="Save As" colon={false} labelAlign="left" style={{ marginBottom: "5px" }}>
                <Input
                  placeholder="Template name"
                  value={this.state.templateName}
                  onChange={(name) => this.onTemplateNameChange(name.target.value)}
                />
              </Form.Item>
            </Col>
            <Col span={24} style={{ marginTop: "10px" }}>
              <Editor
                height="73vh"
                language="yaml"
                theme="console"
                value={this.state.data}
                options={options}
                editorDidMount={this.handleEditorDidMount}
              />
            </Col>
            <Col span={24}>
              <Space style={{ marginTop: "10px" }}>
                <Button size="default" shape="round">
                  <CloseOutlined />
                  Cancel
                </Button>
                <Button
                  size="default"
                  shape="round"
                  type="primary"
                  onClick={() => this.onSubmit(this.state.templateName, false)}
                  disabled={!this.state.isEditorReady || this.state.templateName === ""}
                >
                  <CheckOutlined />
                  Save
                </Button>
                <Button
                  size="default"
                  shape="round"
                  type="primary"
                  onClick={() => this.onSubmit(this.state.templateName, true)}
                  disabled={!this.state.isEditorReady || this.state.templateName === ""}
                >
                  <CheckOutlined />
                  Save & Execute
                </Button>
              </Space>
            </Col>
          </Row>
        </Form>
      </>
    )
  }
}

export default withRouter(SubmitTemplate)
