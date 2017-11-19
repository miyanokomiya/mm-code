<template>
  <div id="app">
    <v-app id="inspire" v-resize="onResize">
      <v-navigation-drawer
        fixed
        v-model="drawer"
        app
      >
        <v-list dense>
          <v-list-tile @click="drawer = false">
            <v-list-tile-action>
              <v-icon>arrow_back</v-icon>
            </v-list-tile-action>
          </v-list-tile>
          <v-list-tile @click="">
            <v-list-tile-action>
              <v-icon>settings</v-icon>
            </v-list-tile-action>
            <v-list-tile-content>
              <v-list-tile-title>Settings</v-list-tile-title>
            </v-list-tile-content>
          </v-list-tile>
          <v-list-tile @click="">
            <v-list-tile-action>
            </v-list-tile-action>
            <v-list-tile-content>
              <v-list-tile-title>
                <v-checkbox
                  label="Chase File"
                  v-model="isChaseFile"
                />
              </v-list-tile-title>
            </v-list-tile-content>
          </v-list-tile>
          <v-list-tile @click="">
            <v-list-tile-action>
            </v-list-tile-action>
            <v-list-tile-content>
              <v-list-tile-title>
                <v-checkbox
                  label="Chase Cursor"
                  v-model="isChaseCursor"
                />
              </v-list-tile-title>
            </v-list-tile-content>
          </v-list-tile>
        </v-list>
      </v-navigation-drawer>
      <v-toolbar color="indigo" height="36px" dark fixed app>
        <v-toolbar-side-icon @click.stop="drawer = !drawer"></v-toolbar-side-icon>
        <v-toolbar-title>
        </v-toolbar-title>
        <div class="d-flex align-center" style="margin-left: auto">
          <v-btn icon @click="splitViewCount = (splitViewCount) % 4 + 1">
            <v-icon>chrome_reader_mode</v-icon>
          </v-btn>
        </div>
      </v-toolbar>
      <v-content>
        <v-layout row wrap>
          <v-flex
            v-for="(fileName, i) in currentViewFileList"
            :key="i"
            :xs6="currentViewFileList.length === 2 || (currentViewFileList.length === 3 && i !== 2) || currentViewFileList.length === 4"
            :xs12="currentViewFileList.length === 1 || (currentViewFileList.length === 3 && i === 2)"
          >
            <EditorContainer
              v-if="getViewFileForEditor(i)"
              class="editor-container"
              :class="{'editor-container-right': i % 2 === 1}"
              :editorHeight="adjustedEditorHeight"
              :editorKeyList="editorKeyList"
              :viewEditor="editors[getViewFileForEditor(i)]"
              :isChaseCursor="i === 0 ? isChaseCursor : false"
              @changeFile="val => setViewFileForEditor(i, val)"
            />
          </v-flex>
        </v-layout>
      </v-content>
    </v-app>
  </div>
</template>

<script>
import Vue from 'vue'
import EditorContainer from './components/EditorContainer'

let websocket = null
let connectionLoop = null

export default {
  components: {
    EditorContainer
  },
  data () {
    return {
      activeEditorName: null,
      editors: {},
      drawer: false,
      editorHeight: 300,
      settings: {
        isChaseFile: true,
        isChaseCursor: true,
        splitViewCount: 1
      },
      viewFileList: ['', '', '', '']
    }
  },
  computed: {
    activeEditor () {
      return this.editors[this.activeEditorName]
    },
    editorKeyList () {
      return Object.keys(this.editors)
    },
    currentViewFileList () {
      return this.viewFileList.slice(0, this.settings.splitViewCount)
    },
    adjustedEditorHeight () {
      if (this.splitViewCount < 3) {
        return this.editorHeight
      } else {
        return (this.editorHeight) / 2
      }
    },
    isChaseFile: {
      get () {
        return this.settings.isChaseFile
      },
      set (val) {
        this.settings.isChaseFile = val
      }
    },
    isChaseCursor: {
      get () {
        return this.settings.isChaseCursor
      },
      set (val) {
        this.settings.isChaseCursor = val
      }
    },
    splitViewCount: {
      get () {
        return this.settings.splitViewCount
      },
      set (val) {
        this.settings.splitViewCount = val
      }
    }
  },
  watch: {
    settings: {
      handler (to, from) {
        localStorage.setItem('settings', JSON.stringify(to))
      },
      deep: true
    },
    activeEditorName (to, from) {
      if (this.settings.isChaseFile) {
        this.setViewFileForEditor(0, to)
      }
    }
  },
  created () {
    this.initSocket()
    this.onResize()
    const settings = localStorage.getItem('settings')
    if (settings) {
      try {
        this.settings = Object.assign({}, this.settings, JSON.parse(settings))
      } catch (e) {
        console.log('Failed to load settings from local-strage.')
      }
    }
  },
  methods: {
    setViewFileForEditor (index, val) {
      Vue.set(this.viewFileList, index, val)
    },
    getViewFileForEditor (index) {
      const fileName = this.viewFileList[index]
      return fileName || this.activeEditorName
    },
    onResize () {
      const height = window.innerHeight - 36
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
        // ファイル取得
        if (jsonData.fileName) {
          this.activeEditorName = jsonData.fileName
          Vue.set(this.editors, jsonData.fileName, {
            fileName: jsonData.fileName,
            lines: jsonData.text.split('\n'),
            cursor: {
              row: 0,
              column: 0
            }
          })
        }
      } else if (type === 'line') {
        // 1行更新
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
}
.editor-container {
  border-left: 1px solid #fff;
  border-right: 1px solid #fff;
}
.editor-container.editor-container-right {
  border-left: none;
}

.tabs__item, .btn__content {
  text-transform: none;
}
</style>
