const MonacoWebpackPlugin = require("monaco-editor-webpack-plugin")

module.exports = function override(config, env) {
  config.plugins.push(
    new MonacoWebpackPlugin({
      languages: ["yaml"],
      features: [
        "!accessibilityHelp",
        "!bracketMatching",
        "!caretOperations",
        "clipboard",
        "!codeAction",
        "!codelens",
        "!colorDetector",
        "!comment",
        "!contextmenu",
        "!coreCommands",
      ],
    })
  )
  return config
}
