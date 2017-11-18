<template>
  <v-app id="app" v-resize="onResize">
    <v-navigation-drawer
      fixed
      v-model="drawer"
      app
    >
      <v-list dense>
        <v-list-tile @click="">
          <v-list-tile-action>
            <v-icon>home</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title>Home</v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>
        <v-list-tile @click="">
          <v-list-tile-action>
            <v-icon>contact_mail</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title>Contact</v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>
      </v-list>
    </v-navigation-drawer>
    <v-toolbar color="indigo" dark fixed app>
      <v-toolbar-side-icon @click.stop="drawer = !drawer"></v-toolbar-side-icon>
      <v-toolbar-title>Application</v-toolbar-title>
    </v-toolbar>
    <v-content>
      <v-tabs grow v-model="viewEditorTabKey">
        <v-tabs-bar class="cyan" dark>
          <v-tabs-item
            v-for="key in editorKeyList"
            :key="key"
            :href="'#editor_' + key"
            ripple
          >
            {{key.slice(key.lastIndexOf('/') + 1)}}
          </v-tabs-item>
          <v-tabs-slider color="yellow"></v-tabs-slider>
        </v-tabs-bar>
        <v-tabs-items>
          <v-tabs-content
            v-for="key in editorKeyList"
            :key="key"
            :id="`editor_${key}`"
          >
            <EditorBox
              :ref="editors[key].fileName"
              class="editor-box"
              :style="{height: `${editorHeight}px`}"
              :fileName="editors[key].fileName"
              :lines="editors[key].lines"
              :cursor="editors[key].cursor"
              :autoChase="autoChase"
            />
          </v-tabs-content>
        </v-tabs-items>
      </v-tabs>
    </v-content>
    <!-- <v-footer color="indigo" app>
      <span class="white--text">&copy; 2017</span>
    </v-footer> -->
  </v-app>
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
      editors: {},
      drawer: false,
      viewEditorName: null,
      editorHeight: 300,
      autoChase: true
    }
  },
  computed: {
    activeEditor () {
      return this.editors[this.activeEditorName]
    },
    editorKeyList () {
      return Object.keys(this.editors)
    },
    viewEditorTabKey: {
      get () {
        return this.viewEditorName ? `editor_${this.viewEditorName}` : null
      },
      set (val) {
        this.viewEditorName = val.slice(7)
      }
    }
  },
  mounted () {
    this.initSocket()
  },
  methods: {
    onResize () {
      const height = window.innerHeight - (56 + 48)
      this.editorHeight = height
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
        if (!this.viewEditorName || this.autoChase) {
          this.viewEditorName = this.activeEditorName
        }
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
        // 1行更新
        if (this.autoChase) {
          this.viewEditorName = this.activeEditorName
        }
        const editor = this.activeEditor
        if (editor) {
          Vue.set(editor.lines, parseInt(jsonData.r), jsonData.l)
          editor.cursor.row = parseInt(jsonData.r)
          editor.cursor.column = parseInt(jsonData.c)
        }
      } else if (type === 'updates') {
        // 複数行更新
        if (this.autoChase) {
          this.viewEditorName = this.activeEditorName
        }
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
}

.editor-box {
}
</style>
