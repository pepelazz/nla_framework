import axios from 'axios'

export default async ({Vue}) => {
  Vue.use({
    install: (Vue, options) => {
      Vue.prototype.$axios = axios
    }
  }, {})
}
