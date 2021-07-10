import config from '../app/plugins/config'

// "async" is optional
export default async ({app}) => {
  app.config.globalProperties.$config = config
}
