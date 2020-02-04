import React from "react";
import { Row, Col, Input, Select, Button } from "antd";
import PageTitle from "../PageTitle/PageTitle";

const { TextArea } = Input;
const { Option } = Select;

const codeSubmitForm = ({title, language, onLanguageChange, codeString, onCodeChange, onSubmit}) => {
  return (
    <React.Fragment>
      <PageTitle title={title} />
      <Row gutter={["10", "10"]}>
        <Col>
          <span style={{ fontWeight: "600" }}>Language Selection: </span>
          <Select
            style={{ width: 200 }}
            value={language}
            onChange={onLanguageChange}
          >
            <Option value="yaml">YAML</Option>
            <Option value="json">JSON</Option>
          </Select>
        </Col>
        <Col>
          <TextArea
            style={{ minHeight: "50vh" }}
            value={codeString}
            onChange={onCodeChange}
          />
        </Col>
        <Col>
          <Button
            size="default"
            shape="round"
            icon="check"
            type="primary"
            onClick={onSubmit}
          >
            Submit
          </Button>
          <Button
            size="default"
            shape="round"
            icon="close"
            style={{ marginLeft: "7px" }}
          >
            Cancel
          </Button>
        </Col>
      </Row>
    </React.Fragment>
  );
};

export default codeSubmitForm;
