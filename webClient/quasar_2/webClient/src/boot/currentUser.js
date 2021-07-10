import CurrentUser from '../app/plugins/CurrentUser'

export default async ({app}) => {
  app.config.globalProperties.$currentUser = new CurrentUser()
}

