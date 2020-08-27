import { message } from "antd"

export const errorMsg = (text, duration = 5) => {
  message.error(text, duration)
}

export const infoMsg = (text, duration = 5) => {
  message.info(text, duration)
}
