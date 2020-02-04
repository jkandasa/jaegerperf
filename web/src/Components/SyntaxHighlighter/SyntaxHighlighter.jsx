import React from "react"
import { Light as SyntaxHighlighter } from "react-syntax-highlighter";
import yaml from "react-syntax-highlighter/dist/esm/languages/hljs/yaml";
import json from "react-syntax-highlighter/dist/esm/languages/hljs/json";
import cStyle from "react-syntax-highlighter/dist/esm/styles/hljs/github";

SyntaxHighlighter.registerLanguage("yaml", yaml);
SyntaxHighlighter.registerLanguage("json", json);

const highlighter = ({ code, language }) => {
  return (
    <SyntaxHighlighter
      language={language}
      style={cStyle}
      showLineNumbers={true}
    >
      {code}
    </SyntaxHighlighter>
  );
};

export default highlighter;
