import axios from 'axios'

export default async ({app}) => {
  app.config.globalProperties.$axios = axios
}

