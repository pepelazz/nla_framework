<template>
  <div>
    <q-drawer :value="rightSide" side="right" bordered @hide="$emit('hide')">
      <q-list separator>
        <q-item-label header>
          Задачи
        </q-item-label>
<!--        <div style="position: absolute; top: 10px; right: 10px">-->
<!--          <q-btn round flat color="secondary" icon="refresh" size="sm"/>-->
<!--        </div>-->
        <q-separator/>
        <q-item v-for="item in listForRender" :key="item.id">
          <q-item-section avatar top @click="$router.push(`/task/${item.id}`)">
            <q-avatar v-if="item.isDeadlinePass" icon="warning" color="orange" text-color="white"/>
            <q-avatar v-else icon="error_outline" color="info" text-color="white"/>
          </q-item-section>
          <q-item-section>
            <q-item-label>{{item.task_type_title}}</q-item-label>
            <q-item-label caption>{{formatDate(item.deadline)}}</q-item-label>
<!--            <q-item-label v-if="item.table_name === 'client'" caption @click="$router.push(`/client/${item.table_id}`)"><q-icon name="far fa-building"/> {{item.table_options.title}}</q-item-label>-->
<!--            <q-item-label v-if="item.table_name === 'deal'" caption @click="$router.push(`/client/${item.table_options.client_id}/deal/${item.table_id}`)"><q-icon name="opacity"/> {{item.table_options.client_title}} {{item.table_options.deal_title}}</q-item-label>-->
          </q-item-section>
          <q-item-section side>
            <div class="text-grey-8">
              <q-btn size="12px" flat dense round icon="done"/>
<!--              <q-btn size="12px" flat dense round icon="done" @click="$refs.doneTaskDialog.open(item)"/>-->
            </div>
          </q-item-section>
        </q-item>
      </q-list>
      <q-separator/>
    </q-drawer>
<!--    <done-task-dialog ref="doneTaskDialog"></done-task-dialog>-->

  </div>
</template>

<script>
    export default {
        props: ['currentUser', 'rightSide'],
        computed: {
            listForRender: function () {
                return this.list.filter(v => v.state === 'in_process')
            }
        },
        data() {
            return {
                list: [],
            }
        },
        methods: {
            formatDate(d) {
                return this.$utils.formatPgDate(d)
            },
        },
        mounted() {
            if (this.currentUser.id) {
                // первоначальная загрузка список задач
                this.$userTasks.loadList()
                // подписываемся на получаение обновлений
                this.$userTasks.getList$().subscribe(v => {
                    const i = this.list.findIndex(v1 => v1.id === v.id)
                    i === -1 ? this.list.push(v) : this.list.splice(i, 1, v)
                    this.$emit('updateCounter', this.listForRender.length)
                })
            }
        }
    }
</script>
