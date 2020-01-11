import React from "react";
import PageTitle from "../../Components/PageTitle/PageTitle";

import MonacoEditor from "react-monaco-editor";

class SpansGenerator extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      code: "// type your code..."
    };
  }
  editorDidMount(editor, monaco) {
    console.log("editorDidMount", editor);
    editor.focus();
  }
  onChange(newValue, e) {
    console.log("onChange", newValue, e);
  }

  render() {
    const code = this.state.code;
    const options = {
      selectOnLineNumbers: true
    };
    return (
      <React.Fragment>
        <PageTitle title={"Spans Generator"} />
        <MonacoEditor
          width="800"
          height="600"
          language="javascript"
          theme="vs-dark"
          value={code}
          options={options}
          onChange={this.onChange}
          editorDidMount={this.editorDidMount}
        />
      </React.Fragment>
    );
  }
}

export default SpansGenerator;
