import {of} from 'rxjs'
import {take, map, catchError} from 'rxjs/operators'
import {ajax} from 'rxjs/ajax'
import {Notify} from 'quasar'
import config from './config'
import _ from 'lodash'
import fp from 'lodash/fp'
import moment from 'moment'

const postApiRequest = ({url, params, isShowError = true}) => {
  return ajax({
    url: `${config.apiUrl()}${url}`,
    method: 'POST',
    headers: getHttpHeaders(),
    body: {
      params,
    }
  }).pipe(
    take(1),
    map(processResponse()),
    catchError(processError(isShowError))
  )
}

const postCallPgMethod = ({method, params, isShowError = true, successMsg = null}) => {
  return ajax({
    url: `${config.apiUrl()}/api/call_pg_func`,
    method: 'POST',
    headers: getHttpHeaders(),
    body: {
      method,
      params,
    }
  }).pipe(
    take(1),
    map(processResponse({successMsg})),
    catchError(processError(isShowError))
  )
}

// метод когда результат вызова не нужен. Нужно только сообщение ок или ошибка
const pgCall = ({method, params, cb, successMsg}) => postCallPgMethod({
  method,
  params,
  successMsg
}).subscribe(v => cb ? cb(v.ok ? v.result : null) : null)

// метод для загрузки/создания документа по id
const getDocItemById = function ({method, cb}) {
  if (this.id === 'new') {
    this.item = {id: -1}
    // заполняем поля вновь созданного документа либо default значением, либо null
    _.flattenDeep(this.flds).map(v => {
      this.item[v.name] = v.default || null
    })
  } else {
    // дефолтный callback
    cb = cb || ((v) => {
      this.item = v
    })
    // добавляем в цепочку обработки функцию по извлечению полей из options
    const composeCb = (v) => fp.compose(cb, extractOptionsFlds(this))(v)
    pgCall({method, params: {id: this.id}, cb: composeCb})
  }
}

// функция сохранения/создани документа
const saveItem = function ({method, itemForSaveMod = {}, resultModify, cb}) {
  let item = Object.assign({}, this.item)
  const fldNames = _.flattenDeep(this.flds).map(({name}) => name)
  if (!item.options) item.options = {}
  // накатываем на копию item модификатор полей itemForSaveMod
  item = Object.assign(item, itemForSaveMod)
  // отдельно обрабатываем поля, которые сохраняются в колонку options. Если есть this.optionsFlds, то переносим поля с такими именем из item в item.options
  if (this.optionsFlds) this.optionsFlds.map(fldName => item.options[fldName] = item[fldName])
  const itemForSave = Object.assign({}, _.pick(item, ['id', ...fldNames, 'options']))
  let notFilledFlds = _.flattenDeep(this.flds).filter(v => v.required && !itemForSave[v.name])
  if (notFilledFlds.length > 0) {
    notFilledFlds.map(v => {
      this.$q.notify({
        color: 'negative',
        position: 'bottom',
        message: `не заполнено поле '${v.label}'`,
      })
    })
    return
  }
  postCallPgMethod({
    method,
    params: itemForSave,
    successMsg: 'изменения сохранены'
  }).subscribe(v => {
    if (v.ok) {
      // в случае создания нового документа, после сохранения и получения id из базы переходим по новому url
      if (this.item.id === -1) {
        // if (!this.docUrl) console.warn(`$utils.itemSave missed this.docUrl`)
        this.$router.push(`${this.docUrl}/${v.result.id}`)
      }
      resultModify = resultModify || (v => v)
      this.item = fp.compose(resultModify, extractOptionsFlds(this))(v.result)
      this.$emit('updated', this.item)
      if (cb) cb()
    }
  })
}

// shortcut для postCallPgMethod
const callPgMethod = (method, params, cb) => {
  postCallPgMethod({method, params}).subscribe(res => {
    if (res.ok && cb) cb(res.result)
  })
}

// функция для перекладывания поля из options в item, чтобы можно было их редактировать
const extractOptionsFlds = (that) => (v) => {
  if (that.optionsFlds) that.optionsFlds.map(fldName => v[fldName] = v.options[fldName])
  return v
}

const updateUrlQuery = (params = {}, isAdd = true) => {
  let searchParams = new URLSearchParams(window.location.search)
  Object.keys(params).forEach(function (key) {
    if (isAdd && params[key] !== null) {
      searchParams.set(key, params[key])
    } else {
      searchParams.delete(key)
    }
  })

  let newurl = window.location.protocol + '//' + window.location.host + window.location.pathname
  if (searchParams.toString().length > 0) {
    newurl = newurl + '?' + searchParams.toString()
  }
  window.history.pushState({path: newurl}, '', newurl)
}

const formatPgDateTime = (d) => {
  return d ? moment(d, 'YYYY-MM-DDTHH:mm:ss').format('DD-MM-YYYY HH:mm') : null
}

const formatPgDate = (d) => {
  return d ? moment(d, 'YYYY-MM-DDTHH:mm:ss').format('DD-MM-YYYY') : null
}

const notifySuccess = (msg) => {
  Notify.create({
    color: 'positive',
    position: 'bottom-right',
    message: msg
  })
}

const notifyError = (msg) => {
  Notify.create({
    color: 'negative',
    position: 'top-right',
    message: msg
  })
}

[[FunctionsList]]

export default {
  postApiRequest,
  postCallPgMethod,
  callPgMethod,
  pgCall,
  getDocItemById,
  saveItem,
  updateUrlQuery,
  formatPgDateTime,
  formatPgDate,
  notifySuccess,
  notifyError,
  _,
  [[ExportDefaultList]]
}

const getHttpHeaders = () => {
  let headers = {'Content-Type': 'application/json'}
  let authToken = localStorage.getItem(config.appName)
  if (authToken) headers['Auth-token'] = authToken
  return headers
}

const processResponse = ({successMsg} = {}) => (res) => {
  if (res.response && !res.response.ok) throw new Error(res.response.message)
  if (res.response && res.response.ok && successMsg) {
    Notify.create({
      color: 'positive',
      position: 'bottom-right',
      message: successMsg
    })
  }

  return res.response
}

const processError = (isShowError) => (err) => {
  const message = err.response ? err.response.message : err.message
  if (isShowError) {
    Notify.create({
      color: 'negative',
      position: 'top-right',
      message
    })
  }
  return of({ok: false, message})
}
