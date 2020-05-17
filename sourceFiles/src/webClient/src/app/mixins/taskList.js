import {Subject} from 'rxjs'
import {takeUntil} from 'rxjs/operators'

export default {
  computed: {
    taskListForRender: function () {
      return this.taskList.filter(v => v.state === 'in_process')
    }
  },
  data() {
    return {
      taskList: [],
      taskSubscribeDestroy$: new Subject(),
    }
  },
  mounted() {
    // подписываемся на получаение списка задач
    this.$userTasks.getList$().pipe(takeUntil(this.taskSubscribeDestroy$)).subscribe(v => {
      // из всего списка заданий выбираем только относящиеся к данному документу
      if (v.table_name === this.tableName && v.table_id === +this.id) {
        const i = this.taskList.findIndex(v1 => v1.id === v.id)
        i === -1 ? this.taskList.push(v) : this.taskList.splice(i, 1, v)
      }
    })
  },
  destroyed() {
    this.taskSubscribeDestroy$.next(true)
    this.taskSubscribeDestroy$.complete()
  }
}
