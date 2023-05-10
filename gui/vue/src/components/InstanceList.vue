<template>
  <el-container>
    <el-menu :default-active=null class="redis-instance-list" mode="vertical" @select="handleSelect">
      <el-menu-item v-for="item in instance_list" :index="item.host" :key="item.host"> {{ item.host }}</el-menu-item>
    </el-menu>
  </el-container>
</template>

<script>
import axios from '@/components/axios/axios'

export default {
  name: 'InstanceList',
  data () {
    return {
      instance_list: [
        {
          host: '127.0.0.1',
          analyze_start_time: '2022-01-02 12:13:16',
          analyze_end_time: '2022-01-02 12:13:16',
          is_finish: false
        }
      ]
    }
  },
  created () {
    this.update()
  },
  methods: {
    update () {
      axios.post('/api/get_instance_list')
        .then(response => {
          if (response.data) {
            console.log(response.data)
          } else {
            console.log('no instance')
          }
          this.instance_list = response.data
        })
        .catch(function (error) {
          console.log(error)
        })
      console.log('instance list update')
      console.log(this.instance_list)
    },
    handleSelect (index) {
      const host = this.instance_list.find(item => item.host === index).host
      this.$emit('click_redis_instance', host)
    }
  }
}
</script>

<style scoped>
.redis-instance-list {
  width: 200px;
  height: 100%;
  border-right: 1px solid #eee;
}
</style>
