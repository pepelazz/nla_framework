import {BehaviorSubject} from 'rxjs'
import config from './config'
import utils from './utils'
import {Loading} from 'quasar'

let user$ = new BehaviorSubject(null)
let isInLogingProcess$ = new BehaviorSubject(false) // флаг для стадии процесса логина

const CurrentUser = class {
  getUser$ = () => user$

  login = ({user, auth_token} = {}) => {
    // в случае если уже передали авторизованного пользователя (например, при первоначальной авторизации через соцсети)
    if (user) {
      user$.next(user)
      localStorage.setItem(config.appName, user.auth_token)
    } else {
      // вариант логина по токену, который ищем в localStorage
      // если auth_token не передан в параметрах, то ищем его в localStorage
      if (!auth_token) auth_token = localStorage.getItem(config.appName)
      // записываем токен в localstorage, чтобы в header запроса был подставлен новый токен
      localStorage.setItem(config.appName, auth_token)
      loginProcess(true)

      utils.postApiRequest({url: '/api/current_user', params: {auth_token}, isShowError: false}).subscribe(res => {
        if (res.ok) {
          user$.next(res.result)
          localStorage.setItem(config.appName, res.result.auth_token)
          loginProcess(false)
        } else {
          this.logout()
          loginProcess(false)
        }
      })
    }
  }

  logout = () => {
    user$.next(null)
    localStorage.removeItem(config.appName)
  }

  getIsInLogingProcess = () => isInLogingProcess$
}

const loginProcess = (isTrue) => {
  if (isTrue) {
    Loading.show({message: 'авторизация'})
    isInLogingProcess$.next(true)
  } else {
    Loading.hide()
    isInLogingProcess$.next(false)
  }
}

export default CurrentUser
