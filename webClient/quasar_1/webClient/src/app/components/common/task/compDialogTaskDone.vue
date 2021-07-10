<template>
  <q-dialog v-model="isShowDialog" @hide="hide">
    <q-card style="width: 700px; max-width: 80vw;">
      <q-card-section>
        <div class="text-h6">{{task.task_type_title}}</div>
      </q-card-section>

      <q-card-section class="q-mt-md">
        <q-slider
          v-model="success_rate"
          :min="0"
          :max="10"
          :step="1"
          snap
          label-always
          :label-value="`Оценка успешности: ${success_rate}`"
          markers
          color="primary"
        />
      </q-card-section>

      <q-card-section>
        <q-input outlined label="Результат" v-model="taskResult" autogrow/>
      </q-card-section>

      <q-card-actions align="right">
        <q-btn flat label="Отмена" v-close-popup/>
        <q-btn flat label="OK" color="primary" @click="done"/>
      </q-card-actions>
    </q-card>

  </q-dialog>
</template>

<script>
    export default {
        data() {
            return {
                isShowDialog: false,
                task: {},
                success_rate: 0,
                taskResult: null,
            }
        },
        methods: {
            open(task) {
                this.task = task
                this.isShowDialog = true
            },
            hide() {
                this.success_rate = 0
                this.taskResult = null
            },
            done() {
                this.$utils.postCallPgMethod({method: 'task_action_to_finished', params: {id: this.task.id, result: this.taskResult, success_rate: this.success_rate}}).subscribe(res => {
                    if (res.ok) {
                      this.$emit('taskFinished', this.task.id)
                      this.isShowDialog = false
                    }
                })
            }
        }
    }
</script>
