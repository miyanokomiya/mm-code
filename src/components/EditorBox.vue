<template>
  <div ref="editorWrapper" class="editor-wrapper">
    <div class="line-count-wrapper">
      <ul class="line-count">
        <li v-for="(l, i) in lines" :key="i">
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

export default {
  name: 'this',
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
    autoScroll: {
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
      return this.lines.join('\n')
    },
    languageType () {
      return languageType(this.fileName)
    }
  },
  updated () {
    if (this.$refs.editorWrapper) {
      const width = this.$refs.editorWrapper.scrollWidth
      if (this.$refs.lineHighlight.length > 0) {
        this.$refs.lineHighlight[0].style.width = `${width}px`
      }
      if (this.$refs.cursor.length > 0) {
        // フォントサイズとパディングを考慮して絶妙な位置に調整
        this.$refs.cursor[0].style.left = `${13 + this.cursor.column * 9.6}px`
      }
    }
    if (this.autoScroll) {
      this.adjustScroll()
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
    adjustScroll () {
      if (this.$refs.editorWrapper) {
        if (this.$refs.lineHighlight.length > 0) {
          const lineHighlight = this.$refs.lineHighlight[0]
          const parent = lineHighlight.parentElement
          const top = parent.offsetTop
          this.$refs.editorWrapper.scrollTop = top - 200
        }
      }
    }
  }
}
</script>

<style lang="scss">
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
