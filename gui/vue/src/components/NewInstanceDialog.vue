<template>
  <el-dialog class="new_instance_dialog" title="New Instance" :close-on-click-modal="false"
             :visible.sync="dialogVisible">
    <el-form>
      <el-row :gutter=15>
        <el-col :span=8>
          <el-form-item label="Host" required>
            <el-input v-model="instance.host" placeholder="127.0.0.1"></el-input>
          </el-form-item>
        </el-col>
        <el-col :span=8>
          <el-form-item label="Port" required>
            <el-input v-model.number="instance.port" type="number" placeholder="6379"></el-input>
          </el-form-item>
        </el-col>
        <el-col :span=8>
          <el-form-item label="Password">
            <el-input v-model="instance.password"></el-input>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter=15>
        <el-col :span=8>
          <el-form-item label="Scan Count">
            <el-input v-model.number="instance.count" type="number" placeholder="10000"></el-input>
          </el-form-item>
        </el-col>
        <el-col :span=8>
          <el-form-item label="Scan Total">
            <el-input v-model.number="instance.limit" type="number" placeholder="100000"></el-input>
          </el-form-item>
        </el-col>
        <el-col :span=8>
          <el-form-item label="Scan Match">
            <el-input v-model="instance.match" placeholder="*"></el-input>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter=15>
        <el-col :span=8>
          <el-form-item label="Redis Data Type">
            <el-input v-model="instance.types" placeholder="Separate by comma"></el-input>
          </el-form-item>
        </el-col>
        <el-col :span=8>
          <el-form-item label="Key Separator">
            <el-input v-model="instance.separators"></el-input>
          </el-form-item>
        </el-col>
        <el-col :span=8>
          <el-form-item label="Pause(ms)">
            <el-input v-model.number="instance.pause" type="number" placeholder="1000"></el-input>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter=15>
        <el-col :span=8>
          <el-form-item>
            <el-switch active-text="Cluster" inactive-text="Single" v-model="instance.cluster"></el-switch>
          </el-form-item>
        </el-col>
      </el-row>

    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button @click="dialogVisible = false">Cancel</el-button>
      <el-button type="primary" @click="confirm">Confirm</el-button>
    </div>
  </el-dialog>
</template>

<script>
import axios from '@/components/axios/axios'

export default {
  name: 'NewInstanceDialog',
  data () {
    return {
      instance: {
        host: '127.0.0.1',
        port: 6379,
        password: '',
        count: 10000,
        limit: 100000,
        match: '*',
        types: '',
        separators: ':',
        cluster: true,
        pause: 1000
      },
      dialogVisible: false
    }
  },
  props: [],
  methods: {
    confirm () {
      axios.post('/api/rma/start_analyze', this.instance)
        .then(response => {
          this.dialogVisible = false
          this.$emit('newInstanceFinished')
        })
        .catch(function (error) {
          console.log(error)
        })
    },
    show () {
      this.dialogVisible = true
    },
    hidden () {
      this.dialogVisible = false
    }
  }
}
</script>

<style scoped>
.new_instance_dialog {
  width: 60%;
  position: center;
}
</style>
