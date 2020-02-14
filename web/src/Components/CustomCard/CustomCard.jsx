import React from "react"
import { Card } from "antd"

import "./CustomCard.css"

const customCard = ({ title, extra, body }) => {
  return (
    <Card
      className="custom-card"
      size="small"
      title={title ? title : false}
      extra={extra}
      bordered
      type="inner"
      style={{ marginBottom: "10px" }}
    >
      {body}
    </Card>
  )
}

export default customCard
