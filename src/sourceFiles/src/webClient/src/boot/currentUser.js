import CurrentUser from '../app/plugins/CurrentUser'

// "async" is optional
export default ({Vue}) => {
  Vue.use({
    install: (Vue, options) => {
      Vue.prototype.$currentUser = new CurrentUser()
    }
  }, {})
}
