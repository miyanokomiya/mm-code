<template>
<div>
  <v-tabs v-model="viewEditorTabKey" :scrollable="false">
    <v-tabs-bar class="cyan" dark>
      <v-tabs-item
        v-for="key in editorKeyList"
        :key="key"
        :href="`#${uniqueID}_${key}`"
        ripple
      >
        {{sliceFileName(key)}}
      </v-tabs-item>
      <v-tabs-slider color="yellow"></v-tabs-slider>
    </v-tabs-bar>
  </v-tabs>
  <EditorBox
    v-if="viewEditor"
    :ref="viewEditor.fileName"
    :style="{height: `${editorHeight}px`}"
    :fileName="viewEditor.fileName"
    :lines="viewEditor.lines"
    :cursor="viewEditor.cursor"
    :isChaseCursor="isChaseCursor"
  />
</div>
</template>

<script>
import EditorBox from './EditorBox'

export default {
  components: {
    EditorBox
  },
  data () {
    return {
      uniqueID: Math.random()
    }
  },
  props: {
    editorHeight: {
      type: Number,
      default: 300
    },
    editorKeyList: {
      type: Array,
      default: () => []
    },
    viewEditor: {
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
      }
    },
    isChaseCursor: {
      type: Boolean,
      default: true
    }
  },
  computed: {
    viewEditorTabKey: {
      get () {
        return this.viewEditor.fileName ? `${this.uniqueID}_${this.viewEditor.fileName}` : null
      },
      set (val) {
        if (typeof val === 'string') {
          const index = val.indexOf('_')
          const fileName = val.slice(index + 1)
          this.$emit('changeFile', fileName)
        }
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
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
