<template>
  <div>
    <el-tree :data="data" :render-content="renderContent"></el-tree>
  </div>
</template>

<script>
export default {
  data () {
    return {
      data: [
        {
          id: 1,
          label: 'Node 1',
          children: [
            {
              id: 2,
              label: 'Node 1-1',
              comment: 'a'
            },
            {
              id: 3,
              label: 'Node 1-2',
              comment: 'a'
            }
          ]
        }
      ]
    }
  },
  methods: {
    renderContent (h, {node, data, store}) {
      // 自定义节点显示方式
      return <span>{data.label} (ID: {data.comment})</span>
    }
  }
}
</script>
