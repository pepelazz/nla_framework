import utils from '../app/plugins/utils'

export default async ({app}) => {
  app.config.globalProperties.$utils = utils
}

