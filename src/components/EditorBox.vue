<template>
  <div ref="editorWrapper" class="editor-wrapper">
    <!-- <div class="line-count-wrapper">
      <ul class="line-count">
        <li class="line" v-for="(l, i) in lines" :key="i">
          <span ref="lineHighlight" v-if="i === cursor.row" class="line-highlight"></span>
          <span ref="cursor" class="cursor" v-if="i === cursor.row"></span>
          <span>{{i + 1}}</span>
        </li>
      </ul>
    </div> -->
    <!-- <pre :class="`language-${languageType}`"><code v-html="prismHtml" :class="`language-${languageType}`"></code></pre> -->
    <!-- <editor :content="text" :sync="true" :height="`${editorHeight}px`" theme="monokai" :options="{fontSize: '14px', readOnly: true}"></editor> -->
    <div ref="editor" :style="{height: `${editorHeight}px`}"></div>
  </div>
</template>

<script>
import Prism from 'prismjs'
import {languageType} from '@/commons/prismConfig'
import ace from 'brace'
import 'brace/mode/javascript'
import 'brace/mode/html'
import 'brace/theme/monokai'

const SPACE_OF_TAB = '  '

export default {
  data () {
    return {
      editor: null
    }
  },
  props: {
    editorHeight: {
      type: Number,
      default: 300
    },
    fileName: {
      type: String,
      default: 'no name'
    },
    lines: {
      type: Array,
      default: () => []
    },
    cursor: {
      row: {
        type: Number,
        default: 0
      },
      column: {
        type: Number,
        default: 0
      }
    },
    isChaseCursor: {
      type: Boolean,
      default: true
    }
  },
  computed: {
    text: {
      get () {
        return this.lines.join('\n')
      },
      set (val) {
        console.log(val)
      }
    },
    prismHtml () {
      const prismHtml = this.filterPrism(this.convertedText)
      return prismHtml
    },
    convertedText () {
      const text = this.lines.join('\n')
      const ret = text.replace(/\t/g, SPACE_OF_TAB)
      return ret
    },
    languageType () {
      return languageType(this.fileName)
    },
    cursorColumnIndex () {
      const orgColumn = this.cursor.column
      // カーソル位置までの行内タブ文字数を数える
      const lineBefore = this.lines[this.cursor.row].slice(0, orgColumn)
      const count = (lineBefore.match(/\t/g) || []).length
      // タブをスペース換算した分を補う
      return orgColumn + count * (SPACE_OF_TAB.length - 1)
    }
  },
  mounted () {
    this.editor = ace.edit(this.$refs.editor)
    const options = {
      readOnly: true,
      fontSize: '14px'
    }
    this.editor.$blockScrolling = Infinity
    this.editor.getSession().setMode('ace/mode/' + this.languageType)
    this.editor.setTheme('ace/theme/monokai')
    this.editor.setValue(this.text, 1)
    this.editor.setOptions(options)
    this.editor.on('change', () => {
    })

    this.adjustAll()
  },
  watch: {
    text (to, from) {
      this.editor.setValue(to, 1)
      this.adjustAll()
    },
    cursor: {
      handler () {
        this.adjustAll()
      },
      deep: true
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
    adjustAll () {
      if (this.editor) {
        this.editor.gotoLine(this.cursor.row + 1, this.cursor.column, true)
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.editor-wrapper {
  text-align: left;
}
</style>
