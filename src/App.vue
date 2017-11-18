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
                  label="Chase Scroll"
                  v-model="isChaseCursor"
                />
              </v-list-tile-title>
            </v-list-tile-content>
          </v-list-tile>
        </v-list>
      </v-navigation-drawer>
      <v-toolbar color="indigo" dark fixed app>
        <v-toolbar-side-icon @click.stop="drawer = !drawer"></v-toolbar-side-icon>
        <v-toolbar-title>
          <!-- {{sliceFileName(activeEditorName)}} -->
          <v-menu offset-y>
            <v-btn color="primary" dark slot="activator">
              {{sliceFileName(viewEditorName)}}
              <v-icon>arrow_drop_down</v-icon>
            </v-btn>
            <v-list>
              <v-list-tile v-for="key in editorKeyList.concat().sort()" :key="key" @click="viewEditorName = key" v-if="key !== viewEditorName">
                <v-list-tile-title>{{sliceFileName(key)}}</v-list-tile-title>
              </v-list-tile>
            </v-list>
          </v-menu>
        </v-toolbar-title>
      </v-toolbar>
      <v-content>
        <v-tabs v-model="viewEditorTabKey" :scrollable="false">
          <v-tabs-bar class="cyan" dark>
            <v-tabs-item
              v-for="key in editorKeyList"
              :key="key"
              :href="'#editor_' + key"
              ripple
            >
              {{sliceFileName(key)}}
            </v-tabs-item>
            <v-tabs-slider color="yellow"></v-tabs-slider>
          </v-tabs-bar>
          <EditorBox
            v-if="viewEditor"
            :ref="viewEditor.fileName"
            class="editor-box"
            :style="{height: `${editorHeight}px`}"
            :fileName="viewEditor.fileName"
            :lines="viewEditor.lines"
            :cursor="viewEditor.cursor"
            :autoChase="isChaseCursor"
          />
        </v-tabs>
      </v-content>
    </v-app>
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
      editors: {},
      drawer: false,
      viewEditorName: null,
      editorHeight: 300,
      settings: {
        isChaseFile: true,
        isChaseCursor: true
      }
    }
  },
  computed: {
    activeEditor () {
      return this.editors[this.activeEditorName]
    },
    viewEditor () {
      return this.editors[this.viewEditorName]
    },
    editorKeyList () {
      return Object.keys(this.editors)
    },
    viewEditorTabKey: {
      get () {
        return this.viewEditorName ? `editor_${this.viewEditorName}` : null
      },
      set (val) {
        if (typeof val === 'string') {
          this.viewEditorName = val ? val.slice(7) : null
        }
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
    }
  },
  watch: {
    settings: {
      handler (to, from) {
        localStorage.setItem('settings', JSON.stringify(to))
      },
      deep: true
    }
  },
  mounted () {
    this.initSocket()
    const settings = localStorage.getItem('settings')
    if (settings) {
      try {
        this.settings = JSON.parse(settings)
      } catch (e) {
        console.log('Failed to load settings from local-strage.')
      }
    }
  },
  methods: {
    sliceFileName (name) {
      if (name) {
        return name.slice(name.lastIndexOf('/') + 1)
      } else {
        return null
      }
    },
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
        // ファイル取得
        if (jsonData.fileName) {
          this.activeEditorName = jsonData.fileName
          if (!this.viewEditorName || this.settings.isChaseFile) {
            this.viewEditorName = this.activeEditorName
          }
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
        if (this.settings.isChaseFile) {
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
        if (this.settings.isChaseFile) {
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

.tabs__item, .btn__content {
  text-transform: none;
}
</style>
