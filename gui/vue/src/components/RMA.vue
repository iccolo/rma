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
      <el-container>
        <KeyTree v-bind:bindInstance="openInstance" v-on:click_key="showKeyInfo"></KeyTree>
      </el-container>
      <el-container>
        <el-drawer
          :visible="clickRedisKey !== ''"
          :direction="'btt'"
          :size="'30%'"
          :with-header="false"
          :title="'Key Info'"
          :before-close="handleDrawerClose"
          :show-close="true"
          :modal="false">
          <KeyInfo v-bind:host="openInstance" :redisKey="clickRedisKey"></KeyInfo>
        </el-drawer>
      </el-container>
    </el-main>

  </el-container>
</template>

<script>
import InstanceList from '@/components/InstanceList'
import NewInstanceDialog from '@/components/NewInstanceDialog'
import KeyTree from '@/components/KeyTree'
import KeyInfo from '@/components/KeyInfo'

export default {
  name: 'RMA',
  data () {
    return {
      openInstance: '',
      clickRedisKey: ''
    }
  },
  components: {NewInstanceDialog, InstanceList, KeyTree, KeyInfo},
  methods: {
    displayNewInstanceDialog () {
      this.$refs.newInstanceDialog.dialogVisible = true
    },
    newInstanceFinished () {
      this.$refs.instanceList.update()
    },
    openInstanceKeyTree (clickHost) {
      this.openInstance = clickHost
    },
    showKeyInfo (key) {
      this.clickRedisKey = key
    },
    handleDrawerClose () {
      this.clickRedisKey = ''
    }
  }
}
</script>

<style scoped>

</style>
