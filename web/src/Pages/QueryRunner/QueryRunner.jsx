import React from "react"
import { api as API } from "../../Services/Api"
import SubmitTemplate from "../../Components/SubmitTemplate/SubmitTemplate"

const QueryRunner = () => {
  return (
    <SubmitTemplate
      title="Query Runner"
      listTemplateFn={API.query.listTemplate}
      saveTemplateFn={API.query.saveTemplate}
      getTemplateFn={API.query.getTemplate}
      triggerFn={API.query.trigger}
    />
  )
}

export default QueryRunner
