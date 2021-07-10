import config from '../app/plugins/config'

// "async" is optional
export default async ({Vue}) => {
  Vue.use({
    install: (Vue, options) => {
      Vue.prototype.$config = config
    }
  }, {})
}
