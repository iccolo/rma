import axios from 'axios'

axios.defaults.timeout = 10000
axios.defaults.responseType = 'json'

const instance = axios.create({
  baseURL: ''
})

export default instance
