<template>
  <el-container>
    <el-aside>
      <el-button type="primary"
                 icon="el-icon-setting"
                 @click="displayNewInstanceDialog"> New Instance
      </el-button>

      <NewInstanceDialog
        ref="newInstanceDialog"
        @newInstanceFinished="newInstanceFinished">
      </NewInstanceDialog>

      <InstanceList ref="instanceList" v-on:click_redis_instance="openInstanceKeyTree"></InstanceList>
    </el-aside>

    <el-main>
      <KeyTree v-bind:bindInstance="openInstance"></KeyTree>
    </el-main>

  </el-container>
</template>

<script>
import InstanceList from '@/components/InstanceList'
import NewInstanceDialog from '@/components/NewInstanceDialog'
import KeyTree from '@/components/KeyTree'

export default {
  name: 'RMA',
  data () {
    return {
      openInstance: ''
    }
  },
  components: {NewInstanceDialog, InstanceList, KeyTree},
  methods: {
    displayNewInstanceDialog () {
      this.$refs.newInstanceDialog.dialogVisible = true
    },
    newInstanceFinished () {
      this.$refs.instanceList.update()
    },
    openInstanceKeyTree (clickHost) {
      this.openInstance = clickHost
    }
  }
}
</script>

<style scoped>

</style>
