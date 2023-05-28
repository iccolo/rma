<template>
  <el-tree v-if="bindInstance"
           ref="keyTree"
           :props="props"
           :load="loadNode"
           :render-content="renderContent"
           :key="bindInstance"
           @node-click="handleNodeClick"
           lazy>
  </el-tree>
</template>

<script>
import axios from '@/components/axios/axios'

export default {
  name: 'KeyTree',
  data () {
    return {
      props: {
        label: 'label',
        children: 'children',
        isLeaf: 'isLeaf'
      },
      count: 1
    }
  },
  props: ['bindInstance'],
  methods: {
    loadNode (node, resolve) {
      console.log('loadNode', node)
      if (node.level === 0) {
        if (!this.bindInstance) {
          return
        }
        axios.post('/api/rma/get_key_type', {host: this.bindInstance})
          .then(response => {
            console.log('get_key_type response.data', response.data)
            let leafs = []
            for (let item of response.data) {
              leafs.push({label: item})
            }
            console.log('tree data', leafs)
            resolve(leafs)
          }).catch(error => {
            console.log(error)
          })
      } else {
        if (node.isLeaf) {
          resolve([])
          return
        }
        // 获取keyPrefix和keyType
        const {keyPrefix, keyType} = this.backtrace(node)
        // 拉下一层节点
        axios.post('/api/rma/expand', {
          host: this.bindInstance,
          key_type: keyType,
          key_prefix: keyPrefix,
          num_limit: 200,
          sort_var: 1
        })
          .then(response => {
            console.log('/api/expand response.data:', response.data)
            let leafs = []
            for (let item of response.data) {
              leafs.push({label: item.segment, isLeaf: item.child_num === 0, info: item})
            }
            console.log('leafs:', leafs)
            resolve(leafs)
          })
          .catch(error => {
            console.log(error)
          })
      }
    },
    backtrace (node) {
      let cursor = node
      let keyPrefix = ''
      let keyType = ''
      while (cursor.level > 0) {
        if (cursor.level === 1) {
          keyType = cursor.label
        } else {
          keyPrefix = cursor.label + keyPrefix
        }
        cursor = cursor.parent
      }
      return {keyPrefix: keyPrefix, keyType: keyType}
    },
    renderContent (h, {node, data, store}) {
      // 构建节点的内容
      let content = [h('span', {}, node.label)]
      if (node.data.info) {
        content.push(h('span', {
          style: {
            fontSize: '80%',
            color: '#424040',
            marginLeft: '20px',
            span: 6
          }
        }, [this.formatBytes(node.data.info.total_size)]))

        if (node.data.info.child_num > 0) {
          content.push(h('span', {
            style: {
              fontSize: '80%',
              color: '#424040',
              marginLeft: '20px',
              span: 8
            }
          }, ['子key数量:', node.data.info.child_num]))
        }
      }

      return h(
        'span',
        {
          class: 'custom-tree-node'
          // style: {
          //   display: 'inline-block',
          //   width: '100%'
          // }
        },
        content
      )
    },
    formatBytes (bytes) {
      if (bytes === 0) return '0 B'
      const sizes = ['B', 'KB', 'MB', 'GB']
      const i = Math.floor(Math.log(bytes) / Math.log(1024))
      const value = parseFloat((bytes / Math.pow(1024, i)).toFixed(2))
      return value + ' ' + sizes[i]
    },
    handleNodeClick (nodeData) {
      if (nodeData.isLeaf) {
        this.$emit('click_key', nodeData.label)
      }
    }
  }
}
</script>

<style scoped>
.custom-tree-node {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 14px;
  padding-right: 8px;
}
</style>
