<template>
  <div>
    <el-card>
      <div class="header">
        <div class="key">{{ redisKey }}</div>
        <div class="expire-time" v-if="expireTime">
          ExpireAt: {{ expireTime }}
        </div>
      </div>
      <div class="content" v-if="dataType === 'string'">
        <el-row>
          <el-col :span="4" class="label">Value:</el-col>
          <el-col :span="20" class="value">{{ data }}</el-col>
        </el-row>
      </div>
      <div class="content" v-else-if="dataType === 'hash'">
        <el-table :data="tableData" style="width: 100%">
          <el-table-column
            prop="field"
            label="Field"
            width="150"
          ></el-table-column>
          <el-table-column
            prop="value"
            label="Value"
            width="auto"
          ></el-table-column>
        </el-table>
      </div>
      <div class="content" v-else-if="dataType === 'set'">
        <el-table :data="tableData" style="width: 100%">
          <el-table-column
            prop="value"
            label="Member"
            width="auto"
          ></el-table-column>
        </el-table>
      </div>
      <div class="content" v-else-if="dataType === 'zset'">
        <el-table :data="tableData" style="width: 100%">
          <el-table-column
            prop="member"
            label="Member"
            width="150"
          ></el-table-column>
          <el-table-column
            prop="score"
            label="Score"
            width="auto"
          ></el-table-column>
        </el-table>
      </div>
      <div class="content" v-else-if="dataType === 'list'">
        <el-table :data="tableData" style="width: 100%">
          <el-table-column
            prop="index"
            label="#"
            width="auto"
          ></el-table-column>
          <el-table-column
            prop="value"
            label="Value"
            width="auto"
          ></el-table-column>
        </el-table>
      </div>
    </el-card>
  </div>
</template>

<script>
import axios from '@/components/axios/axios'

export default {
  name: 'KeyInfo',
  props: ['host', 'redisKey'],
  data () {
    return {
      dataType: '',
      data: '',
      tableData: [],
      setItems: [],
      listItems: [],
      expireTime: null
    }
  },
  watch: {
    redisKey () {
      if (this.redisKey) {
        this.loadKeyInfo()
      }
    }
  },
  async created () {
    await this.loadKeyInfo()
  },
  methods: {
    async loadKeyInfo () {
      const result = await axios.post('/api/rma/get_key_info', {
        host: this.host,
        key: this.redisKey
      })
      console.log(result)
      if (result.data.type === 'string') {
        this.dataType = 'string'
        this.data = result.data.value
      } else if (result.data.type === 'hash') {
        this.dataType = 'hash'
        this.tableData = Object.entries(
          result.data.value
        ).map(([field, value]) => ({field, value}))
        console.log(this.tableData)
      } else if (result.data.type === 'set') {
        this.dataType = 'set'
        this.tableData = result.data.value.map((value) => ({
          value
        }))
      } else if (result.data.type === 'zset') {
        this.dataType = 'zset'
        this.tableData = result.data.value.map(({member, score}) => ({
          member,
          score
        }))
      } else if (result.data.type === 'list') {
        this.dataType = 'list'
        this.tableData = result.data.value.map((value, index) => ({
          index: index + 1,
          value
        }))
      }
      if (result.data.ttl > 0) {
        const expireDate = new Date(Date.now() + result.data.ttl * 1000)
        this.expireTime = expireDate.toLocaleString()
      }
    }
  }
}
</script>

<style scoped>
.header {
  display: grid;
  grid-template-columns: 1fr auto;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid #ebeef5;
  padding-bottom: 10px;
  margin-bottom: 10px;
}

.key {
  justify-self: start;
  font-weight: bold;
}

.expire-time {
  justify-self: end;
}

.label {
  font-weight: bold;
}

.value {
  word-break: break-all;
}

.content {
  margin-bottom: 10px;
}
</style>
