// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import Prism from 'prismjs'
import 'prismjs/themes/prism-okaidia.css'

import 'prismjs/components/prism-scss'
import 'prismjs/components/prism-sass'
import 'prismjs/components/prism-go'
import 'prismjs/components/prism-json'
import 'prismjs/components/prism-ruby'
import 'prismjs/components/prism-yaml'
import 'prismjs/components/prism-coffeescript'
import 'prismjs/components/prism-typescript'

Vue.config.productionTip = false

Vue.directive('highlightjs', {
  deep: true,
  bind: function bind (el, binding) {
    // on first bind, highlight all targets
    let targets = el.querySelectorAll('code')
    let target
    let i

    for (i = 0; i < targets.length; i += 1) {
      target = targets[i]

      if (typeof binding.value === 'string') {
        // if a value is directly assigned to the directive, use this
        // instead of the element content.
        target.textContent = binding.value
      }

      Prism.highlightElement(target)
    }
  },
  componentUpdated: function componentUpdated (el, binding) {
    // after an update, re-fill the content and then highlight
    let targets = el.querySelectorAll('code')
    let target
    let i

    for (i = 0; i < targets.length; i += 1) {
      target = targets[i]
      if (typeof binding.value === 'string') {
        target.textContent = binding.value
        Prism.highlightElement(target)
      }
    }
  }
})

/* eslint-disable no-new */
new Vue({
  el: '#app',
  template: '<App/>',
  components: { App }
})
