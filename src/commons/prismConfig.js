import 'prismjs/themes/prism-okaidia.css'

import 'prismjs'
import 'prismjs/components/prism-scss'
import 'prismjs/components/prism-sass'
import 'prismjs/components/prism-go'
import 'prismjs/components/prism-json'
import 'prismjs/components/prism-ruby'
import 'prismjs/components/prism-yaml'
import 'prismjs/components/prism-coffeescript'
import 'prismjs/components/prism-typescript'
import 'prismjs/components/prism-bash'

export const languageType = fileName => {
  const split = fileName.split('.')
  const ex = split[split.length - 1].toLowerCase()
  switch (ex) {
    case 'js':
    case 'jsx':
      return 'javascript'
    case 'vue':
    case 'html':
      return 'markup'
    case 'md':
      return 'markdown'
    case 'rb':
      return 'ruby'
    case 'yml':
      return 'yaml'
    case 'css':
      return 'css'
    case 'scss':
      return 'scss'
    case 'go':
      return 'go'
    case 'json':
      return 'json'
    case 'ts':
      return 'typescript'
    default:
      if (ex) {
        return ex
      } else {
        return 'bash'
      }
  }
}

export const unescapeHTML = str => {
  const ret = str
    .replace(/&lt;/g, '<')
    .replace(/&gt;/g, '>')
    .replace(/&nbsp;/g, ' ')
    .replace(/&#13;/g, '\r')
    .replace(/&#10;/g, '\n')
  return ret
}
