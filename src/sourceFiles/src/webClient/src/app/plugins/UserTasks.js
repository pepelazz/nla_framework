import {ReplaySubject} from 'rxjs'
// import config from './config'
import utils from './utils'
//
let taskList$ = new ReplaySubject()

const UserTasks = class {
  getList$ = () => taskList$

  loadList = () => {
    utils.postCallPgMethod({method: 'task_list_for_user', params: {state: 'in_process'}}).subscribe(res => {
      if (res.ok) {
        res.result.map(v => {
          taskList$.next(modifyTask(v))
        })
      }
    })
  }

  addTask = (task) => {
    taskList$.next(modifyTask(task))
  }
}

const modifyTask = (v) => {
  if (new Date(v.deadline) < new Date()) {
    v.isDeadlinePass = true
  }
  return v
}

export default UserTasks
