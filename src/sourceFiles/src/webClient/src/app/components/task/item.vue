<template>
  <q-page padding>

    <comp-breadcrumb :list="[{label:'Задача', to:'/task', docType: 'task', icon: 'fas fa-tasks'},
    {label: item && item.title ? `${item.title}` : '',  docType: 'edit'}]"/>

    <div v-if="item" class="q-mt-sm">
      <!--  поля формы    -->
      <div class="row q-col-gutter-md q-mb-sm" v-for="fldRow in flds">
       <template v-if="Array.isArray(fldRow)">
          <comp-fld v-for="fld in fldRow" :key='fld.name'
                                :fld="item[fld.name]"
                                :type="fld.type"
                                @update="item[fld.name] = $event"
                                :label="fld.label"
                                :selectOptions="fld.selectOptions ? fld.selectOptions() : []"
                                :ajaxSelectTitle="item[fld.ajaxSelectTitle]"
                                :columnClass="fld.columnClass"
                                :compName="fld.compName"
                                :readonly="fld.readonly ? fld.readonly() : false"
                                :pgMethod="fld.pgMethod"
                    >
<!--            <template v-if="fld.name === 'table_id'">-->
<!--              <ref-table-id :fld="item[fld.name]" :label="fld.label" @update="item[fld.name] = $event"/>-->
<!--            </template>-->
          </comp-fld>
       </template>
        <template v-if="fldRow.name==='task_type_id'">
          <ref-table-id :id="id" :task_type_id="item.task_type_id" :table_id="item.table_id"
                        :task_type_title="item.task_type_title" :table_id_title="item.table_id_title"
                        :table_name="item.table_name"
                        @update="updateRefTable"/>
        </template>
      </div>

      <!--  кнопки   -->
      <comp-item-btn-save @save="save" @cancel="$router.push(docUrl)"/>

    </div>
  </q-page>
</template>

<script>
    import refTableId from './comp/refTableId'
    export default {
        components: {refTableId},
        props: ['id'],
        computed: {
            docUrl: () => '/task',
        },
        data() {
            return {
                item: null,
                flds: [
                    [
                        {name: 'executor_id', type: 'userId', label: 'Исполнитель', ajaxSelectTitle: 'executor_fullname', readonly: () => this.id > 0},
                    ],
                    {name: 'task_type_id'},
                    {name: 'table_id'},
                    [
                        {name: 'deadline', type: 'datetime', label: 'deadline', readonly: () => true}
                    ],
                    [
                        {name: 'content', type: 'string', label: 'Текст', columnClass: 'col-xs-12 col-sm-8 col-md-8'},
                        // {name: 'table_name', type: 'string', label: 'Название таблицы'},
                        // {name: 'table_id', type: 'number', label: 'id в таблице'},
                        // {name: 'executor_id', type: 'number', label: 'Исполнитель'},
                        // {name: 'manager_id', type: 'number', label: 'Постановщик'},
                        // {name: 'state', type: 'string', label: 'Статус'},
                        // {name: 'deadline', type: 'timestamp', label: 'Срок'},
                        // {name: 'date_completed', type: 'timestamp', label: 'Дата исполнения'},
                        // {name: 'result', type: 'string', label: 'Отчет об исполнении'},
                        // {name: 'success_rate', type: 'number', label: 'Оценка успешности (0-10)'},
                    ],
                ],
                // список полей для редактирования из options
                optionsFlds: [],
            }
        },
        methods: {
            updateRefTable({task_type_id, table_id}) {
              this.item.task_type_id = task_type_id
              this.item.table_id = table_id
            },
            resultModify(res) {
                return res
            },
            save() {
                this.$utils.saveItem.call(this, {
                    method: 'task_update',
                    itemForSaveMod: {},
                    resultModify: this.resultModify,
                })
            },
        },
        mounted() {
            let cb = (v) => {
                this.item = this.resultModify(v)
            }
            this.$utils.getDocItemById.call(this, {method: 'task_get_by_id', cb})
        }
    }
</script>
