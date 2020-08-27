import React from "react"

import "./PageTitle.css"

const pageTitle = ({ title }) => {
  return (
    <div style={{ paddingTop: "5px", paddingBottom: "10px", width: "100%" }}>
      <div className="page-title">{title}</div>
    </div>
  )
}

export default pageTitle
