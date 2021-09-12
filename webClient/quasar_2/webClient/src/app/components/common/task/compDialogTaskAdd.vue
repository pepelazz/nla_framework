<template>
  <q-dialog v-model="isShowDialog" @hide="hide">
    <q-card style="width: 700px; max-width: 80vw;">
      <q-card-section>
        <div class="text-h6">Создать задачу</div>
        <div class="row q-col-gutter-md q-mb-sm">
          <div class="col-12">
            <comp-fld-ref-search label="Тип задачи" item="" pgMethod="task_type_list" @update="updateTaskType" :ext="{table_name: table_name}"/>
          </div>
        </div>
        <div class="row q-col-gutter-md q-mb-sm">
          <div class="col-xs-12 col-sm-6 col-md-6">
            <comp-fld-user-search label="Ответственный" user="" @update="v=> task.executor_id = v.id"/>
          </div>
          <div class="col-xs-12 col-sm-6 col-md-6">
            <comp-fld-date-time  label="Deadline" @update="v=> task.deadline = v"/>
          </div>
        </div>
        <div class="row q-col-gutter-md q-mb-sm">
          <div class="col-12">
            <q-input outlined label="Текст" v-model="task.content" autogrow/>
          </div>
        </div>
      </q-card-section>

      <q-card-actions align="right">
        <q-btn flat :label="$t('message.cancel')" v-close-popup/>
        <q-btn flat label="OK" color="primary" @click="done"/>
      </q-card-actions>
    </q-card>

  </q-dialog>
</template>

<script>
    export default {
        props: ['table_id', 'table_name'],
        data() {
            return {
                isShowDialog: false,
                task: {},
            }
        },
        methods: {
            open() {
                this.isShowDialog = true
            },
            hide() {
            },
            updateTaskType({id, item: taskType}) {
                this.task.task_type_id = id
                // если тип задачи связан с таблицей tableName, то проставляем связи
                const tableName = taskType.item.table_name
                if (tableName) {
                    this.task.table_name = tableName
                    this.task.table_id = +this.table_id
                } else {
                    this.task.table_name = null
                    this.task.table_id = null
                }
            },
            done() {
                let isError = false
                const requiredFlds = [{fld: 'task_type_id', title: 'Тип задачи'}, {fld: 'deadline', title: 'Deadline'}, {fld: 'executor_id', title: 'Ответственный'}]
                requiredFlds.map(v => {
                  if (!this.task[v.fld]) {
                      isError = true
                      this.$q.notify({message: `Поле "${v.title}" не заполнено`, type: 'negative'})
                  }
                })
                if (isError) return
                let params = Object.assign({id: -1}, this.task)
                this.$utils.postCallPgMethod({method: 'task_update', params}).subscribe(res => {
                    if (res.ok) {
                      this.isShowDialog = false
                      this.task = {}
                      this.$emit('updated')
                    }
                })
            }
        }
    }
</script>
