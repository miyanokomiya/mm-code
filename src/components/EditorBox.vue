<template>
  <div ref="editorWrapper" class="editor-wrapper">
    <div class="line-count-wrapper">
      <ul class="line-count">
        <li class="line" v-for="(l, i) in lines" :key="i">
          <span ref="lineHighlight" v-if="i === cursor.row" class="line-highlight"></span>
          <span ref="cursor" class="cursor" v-if="i === cursor.row"></span>
          <span>{{i + 1}}</span>
        </li>
      </ul>
    </div>
    <pre :class="`language-${languageType}`"><code v-html="prismHtml" :class="`language-${languageType}`"></code></pre>
  </div>
</template>

<script>
import Prism from 'prismjs'
import {languageType} from '@/commons/prismConfig'

const SPACE_OF_TAB = '  '

export default {
  props: {
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
    this.adjustAll()
  },
  updated () {
    this.adjustAll()
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
      if (this.$refs.editorWrapper) {
        const width = this.$refs.editorWrapper.scrollWidth
        if (this.$refs.lineHighlight.length > 0) {
          this.$refs.lineHighlight[0].style.width = `${width}px`
        }
        if (this.$refs.cursor.length > 0) {
          // フォントサイズとパディングを考慮して絶妙な位置に調整
          this.$refs.cursor[0].style.left = `${13 + this.cursorColumnIndex * 7.8}px`
        }
      }
      if (this.isChaseCursor) {
        this.adjustScroll()
      }
    },
    adjustScroll () {
      if (this.$refs.editorWrapper) {
        if (this.$refs.lineHighlight.length > 0) {
          const lineHighlight = this.$refs.lineHighlight[0]
          const parent = lineHighlight.parentElement
          const top = parent.offsetTop
          const current = this.$refs.editorWrapper.scrollTop
          const viewHeight = this.$refs.editorWrapper.offsetHeight
          const noMoveRange = viewHeight * 2 / 5
          const adjustRange = viewHeight * 3 / 4
          const nextTop = top - viewHeight * 2 / 5
          const dif = nextTop - current
          if (Math.abs(dif) < noMoveRange / 2) {
          } else if (Math.abs(nextTop - current) < adjustRange / 2) {
            this.$refs.editorWrapper.scrollTop += (dif / Math.abs(dif)) * 19
          } else {
            this.$refs.editorWrapper.scrollTop = nextTop
          }

          if (this.$refs.cursor.length > 0) {
            const cursor = this.$refs.cursor[0]
            const left = cursor.offsetLeft
            const current = this.$refs.editorWrapper.scrollLeft
            const viewWidth = this.$refs.editorWrapper.offsetWidth
            const noMoveRange = viewWidth * 2 / 5
            const adjustRange = viewWidth * 3 / 4
            const nextLeft = left - viewWidth / 2
            const dif = nextLeft - current
            if (Math.abs(dif) < noMoveRange / 2) {
            } else if (Math.abs(nextLeft - current) < adjustRange / 2) {
              // スクロール１文字分がぴったり合わないが許容範囲
              this.$refs.editorWrapper.scrollLeft += (dif / Math.abs(dif)) * 7.993
            } else {
              this.$refs.editorWrapper.scrollLeft = nextLeft
            }
          }
        }
      }
    }
  }
}
</script>

<style lang="scss" scoped>
$back-color: #272822;

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
    background-color: inherit;
    font-size: 13px;
    font-weight: 100;
    letter-spacing: 0;
    box-shadow: none;
  }
}

.line-count-wrapper {
  float: left;
  padding: 15px 0 0 0;
  color: #FFF;
  border-radius: 4px;
}

.line-count {
  list-style: none;
  padding: 0;
  margin: 0;
  font-size: 12px;
}

.line-count li.line {
  width: 30px;
  text-align: right;
  font-size: 0.9rem;
  line-height: 1.5;
  height: 19px;
  position: relative;
  margin-top: 0;
}

.line-highlight {
  display: inline-block;
  position: absolute;
  height: 100%;
  width: 300px;
  top: -2px;
  left: 0;
  pointer-events: none;
  // background-color: rgba(0, 0, 255, 0.2);
  border-bottom: 2px solid rgba(255, 255, 255, 0.8);
}

.cursor {
  float: right;
  position: relative;
  border-left: 2px solid #fff;
  margin: -2px -2px 0 0;
  display: inline-block;
  height: 100%;
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

code:after, code:before, kbd:after, kbd:before {
  content: "";
}
</style>
