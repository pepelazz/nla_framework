import UserTasks from '../app/plugins/UserTasks'

// "async" is optional
export default ({Vue}) => {
  Vue.use({
    install: (Vue, options) => {
      Vue.prototype.$userTasks = new UserTasks()
    }
  }, {})
}
