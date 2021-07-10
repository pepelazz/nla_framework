<template>
  <div>
    <q-drawer :modelValue="rightSide" side="right" bordered @hide="$emit('hide')">
      <q-list separator>
        <q-item-label header>
          Задачи
        </q-item-label>
<!--        <div style="position: absolute; top: 10px; right: 10px">-->
<!--          <q-btn round flat color="secondary" icon="refresh" size="sm"/>-->
<!--        </div>-->
        <q-separator/>
        <component v-for="item in listForRender" :key="item.id" :is="item.template" :item="item"></component>
      </q-list>
      <q-separator/>
    </q-drawer>
    <comp-dialog-task-done ref="doneTaskDialog"/>

  </div>
</template>

<script>
    [[PrintImports]]
    export default {
        props: ['currentUser', 'rightSide'],
        components: {[[PrintComps]]},
        computed: {
            listForRender: function () {
                return this.list.filter(v => v.state === 'in_process').map(v => {
                    v.template = v.options && v.options.template ? v.options.template : 'defaultTmpl'
                    return v
                })
            },
            icon() {
                return function(tableName) {
                    return this.$config.breadcrumbIcons[tableName] || 'label'
                }
            }
        },
        data() {
            return {
                list: [],
            }
        },
        methods: {
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
