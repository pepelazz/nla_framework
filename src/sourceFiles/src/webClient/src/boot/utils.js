import utils from '../app/plugins/utils'

export default async ({Vue, app}) => {
  Vue.use({
    install: (Vue, options) => {
      Vue.prototype.$utils = utils
    }
  }, {})
}
