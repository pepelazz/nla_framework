<template>
  <q-page padding>

    <comp-breadcrumb :list="[{label:'Тип задачи', to:'/taskType', docType: 'taskType'},
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
                    />
       </template>
      </div>

      <!--  кнопки   -->
      <comp-item-btn-save @save="save" @cancel="$router.push(docUrl)"/>

    </div>
  </q-page>
</template>

<script>
    export default {
        props: ['id'],
        computed: {
            docUrl: () => '/taskType',
            // локализация статусов
            selectOptions: function () {
                return ['client', 'deal'].map(v => {
                    return {value: v, label: this.$t(`taskType.table_name_${v}`)}
                })
            }
        },
        data() {
            return {
                item: null,
                flds: [
                    [
                        {name: 'title', type: 'string', label: 'Название'},
                        {name: 'table_name', type: 'select', selectOptions: () => this.selectOptions, label: 'К чему относится', readonly: () => this.item.id > 0},
                    ],
                ],
                // список полей для редактирования из options
                optionsFlds: [],
            }
        },
        methods: {
            resultModify(res) {
                if (res.table_name) res.table_name = {value: res.table_name, label: this.$t(`taskType.table_name_${res.table_name}`)}
                return res
            },
            save() {
                this.$utils.saveItem.call(this, {
                    method: 'task_type_update',
                    itemForSaveMod: {
                        table_name: this.item.table_name ? this.item.table_name.value : undefined,
                    },
                    resultModify: this.resultModify,
                })
            },
        },
        mounted() {
            let cb = (v) => {
                this.item = this.resultModify(v)
            }
            this.$utils.getDocItemById.call(this, {method: 'task_type_get_by_id', cb})
        }
    }
</script>
