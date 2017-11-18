<template>
  <div id="app">
    <h3>{{fileName || 'NO FILE'}}</h3>
    <div ref="editorWrapper" class="editor-wrapper">
      <div class="line-count-wrapper">
        <ul class="line-count">
          <li v-for="(l, i) in lines" :key="i">
            <span ref="lineHighlight" v-if="i === cursor.row" class="line-highlight"></span>
            <span ref="cursor" class="cursor" v-if="i === cursor.row"></span>
            <span>{{i}}</span>
          </li>
        </ul>
      </div>
      <pre :class="`language-${languageType}`"><code v-html="prismHtml" :class="`language-${languageType}`"></code></pre>
    </div>
  </div>
</template>

<script>
import Vue from 'vue'
import Prism from 'prismjs'
import {languageType} from './commons/prismConfig'

let websocket = null
let connectionLoop = null

export default {
  name: 'this',
  data () {
    return {
      fileName: '',
      lines: [],
      cursor: {
        row: 0,
        column: 0
      }
    }
  },
  computed: {
    prismHtml () {
      const prismHtml = this.filterPrism(this.convertedText)
      return prismHtml
    },
    convertedText () {
      return this.lines.join('\n')
    },
    languageType () {
      return languageType(this.fileName)
    }
  },
  mounted () {
    this.initSocket()
  },
  updated () {
    if (this.$refs.editorWrapper) {
      const width = this.$refs.editorWrapper.scrollWidth
      if (this.$refs.lineHighlight.length > 0) {
        this.$refs.lineHighlight[0].style.width = `${width + 30}px`
      }
      if (this.$refs.cursor.length > 0) {
        // フォントサイズとパディングを考慮して絶妙な位置に調整
        this.$refs.cursor[0].style.left = `${13 + this.cursor.column * 9.6}px`
      }
    }
  },
  methods: {
    filterPrism (text) {
      let ret = text
      const type = Prism.languages[this.languageType]
      if (type) {
        ret = Prism.highlight(text, type)
      }
      return ret
    },
    initSocket () {
      if (websocket) {
        websocket.close()
      }
      websocket = new WebSocket('ws://localhost:8090/ws')
      websocket.onopen = (evt) => {
        this.onOpen(evt)
      }
      websocket.onclose = (evt) => {
        this.onClose(evt)
      }
      websocket.onmessage = (evt) => {
        this.onMessage(evt)
      }
      websocket.onerror = (evt) => {
        this.onError(evt)
      }
    },
    onOpen () {
      console.log('CONNECTED')
      this.doJoin()
    },
    onClose () {
      console.log('DISCONNECTED')
    },
    onMessage (evt) {
      const data = evt.data
      let jsonData = JSON.parse(data)
      const type = jsonData.type

      if (type === 'say') {
      } else if (type === 'text') {
        this.fileName = jsonData.fileName
        this.lines = jsonData.text.split('\n')
      } else if (type === 'line') {
        Vue.set(this.lines, parseInt(jsonData.r), jsonData.l)
        this.cursor.row = parseInt(jsonData.r)
        this.cursor.column = parseInt(jsonData.c)
      } else if (type === 'updates') {
        // 複数行更新
        jsonData.updates.forEach((update) => {
          let stackText = update.t

          // テキスト全体からみた変更箇所
          let start = 0
          let end = 0
          for (let i = 0; i < update.sr; i++) {
            start += this.lines[i].length + 1
            end = start
          }
          start += update.sc

          for (let i = update.sr; i < update.er; i++) {
            end += this.lines[i].length + 1
          }
          end += update.ec

          let text = this.lines.join('\n')
          let newText = text.slice(0, start) + stackText + text.slice(end)
          this.lines = newText.split('\n')
        })
      } else if (type === 'join') {
        // 誰かが入室
        console.log('new comer')
      }
    },
    onError (evt) {
      console.error('connection error. mm retry after 5sec.', evt)
      if (connectionLoop) {
        clearTimeout(connectionLoop)
      }
      connectionLoop = setTimeout(() => {
        this.initSocket()
      }, 5000)
    },
    doJoin () {
      websocket.send(
      'join ' +
        JSON.stringify({
          name: 'room'
        })
      )
    }
  }
}
</script>

<style lang="scss">
$back-color: #272822;

#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 20px;
}

.editor-wrapper {
  overflow: auto;
  background-color: $back-color;

  pre[class*="language-"] {
    overflow: initial;
    margin: 0;
    padding: 1em 24px 1.8em 44px;
  }
  code[class*="language-"] {
    padding-right: 24px;
  }
}

.line-count-wrapper {
  float: left;
  padding: 20px 0 0 0;
  color: #FFF;
  border-radius: 4px;
}

.line-count {
  list-style: none;
  padding: 0;
  margin: 0;
}

.line-count li {
  width: 30px;
  text-align: right;
  font-size: 0.8rem;
  line-height: 1.5;
  height: 1.5rem;
  position: relative;
}

.line-highlight {
  display: inline-block;
  position: absolute;
  height: 100%;
  width: 300px;
  top: -6px;
  left: 0;
  pointer-events: none;
  // background-color: rgba(0, 0, 255, 0.2);
  border-bottom: 2px solid rgba(255, 255, 255, 0.8);
}


.cursor {
  float: right;
  position: relative;
  border-left: 2px solid #fff;
  margin: 0 -2px 0 0;
  display: inline-block;
  height: 1rem;
  vertical-align: middle;

  -webkit-animation:blink 0.5s ease-in-out infinite alternate;
  -moz-animation:blink 0.5s ease-in-out infinite alternate;
  animation:blink 0.5s ease-in-out infinite alternate;
  @-webkit-keyframes blink{
    0% {opacity:1;}
    50% {opacity:1;}
    51% {opacity:0;}
    100% {opacity:0;}
  }
  @-moz-keyframes blink{
    0% {opacity:1;}
    50% {opacity:1;}
    51% {opacity:0;}
    100% {opacity:0;}
  }
  @keyframes blink{
    0% {opacity:1;}
    50% {opacity:1;}
    51% {opacity:0;}
    100% {opacity:0;}
  }
}
</style>
