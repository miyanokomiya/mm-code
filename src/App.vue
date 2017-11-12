<template>
  <div id="app">
    <h3>{{fileName}}</h3>
    <div>
      <div class="line-count-wrapper">
        <ul class="line-count">
          <li v-for="(l, i) in lines" :key="i">
            <span ref="lineHighlight" v-if="i === cursor.row" class="line-highlight"></span>
            <span>{{i}}</span>
          </li>
        </ul>
      </div>
      <div ref="codeWrapper" class="code-wrapper">
        <pre v-highlightjs="lines.join('\n')"><code :class="`language-${langueageType}`"></code></pre>
      </div>
    </div>
  </div>
</template>

<script>
import Vue from 'vue'

let websocket = null

export default {
  name: 'this',
  data () {
    return {
      fileName: 'no name',
      lines: [],
      cursor: {
        row: 0,
        column: 0
      }
    }
  },
  computed: {
    langueageType () {
      const split = this.fileName.split('.')
      const ex = split[split.length - 1].toLowerCase()
      switch (ex) {
        case 'js':
          return 'javascript'
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
          return ex
      }
    }
  },
  mounted () {
    websocket = new WebSocket('ws://localhost:8080/ws')
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
  updated () {
    const width = this.$refs.codeWrapper.clientWidth
    if (this.$refs.lineHighlight.length > 0) {
      this.$refs.lineHighlight[0].style.width = `${width + 30}px`
    }
  },
  methods: {
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
      console.error(evt)
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

.line-count-wrapper {
  float: left;
  padding: 20px 0 0 0;
  background-color: $back-color;
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

.code-wrapper {
  padding: 0 0 20px 0;
  background-color: $back-color;
  border-radius: 4px;
}
</style>
