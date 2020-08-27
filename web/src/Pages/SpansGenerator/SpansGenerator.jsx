import React from "react"
import { api as API } from "../../Services/Api"
import SubmitTemplate from "../../Components/SubmitTemplate/SubmitTemplate"

const SpansGenerator = () => {
  return (
    <SubmitTemplate
      title="Spans Generator"
      listTemplateFn={API.generator.listTemplate}
      saveTemplateFn={API.generator.saveTemplate}
      getTemplateFn={API.generator.getTemplate}
      triggerFn={API.generator.trigger}
    />
  )
}

export default SpansGenerator
