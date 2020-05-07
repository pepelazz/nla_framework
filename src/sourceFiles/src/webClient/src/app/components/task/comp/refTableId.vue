<template>
    <div class="col-12">
    <div class="row q-col-gutter-md q-mb-sm">
      <div class="col-xs-12 col-sm-6 col-md-4">
        <comp-fld-ref-search label="Тип задачи" :item="task_type_title" pgMethod="task_type_list"
                             @update="updateTaskType" :readonly="id>0"/>
      </div>
      <div v-if="tableName" class="col-xs-12 col-sm-6 col-md-4">
        <comp-fld-ref-search ref="tableIdSearch" :label="$config.tablesForTask[tableName]" :item="table_id_title" :pgMethod="tableName + '_list'"
                             @update="updateTableId" :readonly="id>0"/>
      </div>
    </div>
    </div>
</template>

<script>
    export default {
        props: ['id', 'task_type_id', 'table_id', 'task_type_title', 'table_name', 'table_id_title'],
        data() {
            return {
                localTaskTypeId: null,
                localTableId: null,
                tableName: null,
            }
        },
        methods: {
            updateTaskType({id, item: taskType}) {
                this.tableName = taskType.item.table_name
                this.localTaskTypeId = id
                // сбрасываем table_id
                this.localTableId = null
                // сбрасываем надпись в поле select
                this.$nextTick(() => {
                  if (this.$refs.tableIdSearch) this.$refs.tableIdSearch.clear()
                })
                // если тип задачи не связан ни с какой таблицей, то отправляем update
                if (!this.tableName) this.update()
            },
            updateTableId({id}) {
                this.localTableId = id
                this.update()
            },
            update() {
                this.$emit('update', {task_type_id: this.localTaskTypeId, table_id: this.localTableId})
            }
        },
        mounted() {
            this.localTaskTypeId = this.task_type_id
            this.localTableId = this.table_id
            this.tableName = this.table_name
        }
    }
</script>
