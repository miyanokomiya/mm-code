<template>
  <div id="app">
    <h3>{{activeEditorName || 'NO FILE'}}</h3>
    <EditorBox
      class="editor-box"
      v-if="activeEditor"
      :fileName="activeEditor.fileName"
      :lines="activeEditor.lines"
      :cursor="activeEditor.cursor"
    />
  </div>
</template>

<script>
import Vue from 'vue'
import EditorBox from './components/EditorBox'

let websocket = null
let connectionLoop = null

export default {
  components: {
    EditorBox
  },
  data () {
    return {
      activeEditorName: null,
      editors: {}
    }
  },
  computed: {
    activeEditor () {
      return this.editors[this.activeEditorName]
    }
  },
  mounted () {
    this.initSocket()
  },
  methods: {
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
        this.activeEditorName = jsonData.fileName
        Vue.set(this.editors, jsonData.fileName, {
          fileName: jsonData.fileName,
          lines: jsonData.text.split('\n'),
          cursor: {
            row: 0,
            column: 0
          }
        })
      } else if (type === 'line') {
        const editor = this.activeEditor
        if (editor) {
          Vue.set(editor.lines, parseInt(jsonData.r), jsonData.l)
          editor.cursor.row = parseInt(jsonData.r)
          editor.cursor.column = parseInt(jsonData.c)
        }
      } else if (type === 'updates') {
        // 複数行更新
        const editor = this.activeEditor
        if (editor) {
          jsonData.updates.forEach((update) => {
            let stackText = update.t
            // テキスト全体からみた変更箇所
            let start = 0
            let end = 0
            for (let i = 0; i < update.sr; i++) {
              start += editor.lines[i].length + 1
              end = start
            }
            start += update.sc

            for (let i = update.sr; i < update.er; i++) {
              end += editor.lines[i].length + 1
            }
            end += update.ec

            let text = editor.lines.join('\n')
            let newText = text.slice(0, start) + stackText + text.slice(end)
            editor.lines = newText.split('\n')
          })
        }
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
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 20px;
}

.editor-box {
  height: 300px;
}
</style>
