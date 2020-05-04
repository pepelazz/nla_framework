<template>
  <q-page padding>

    <comp-breadcrumb :list="[{label:'Пользователи', to:'/users'}, {label:'Редактирование'}]"/>

    <div v-if="item" class="q-mt-sm">
      <!--  поля формы    -->
      <div class="row q-col-gutter-md q-mb-sm" v-for="fldRow in flds">
        <comp-fld v-for="fld in fldRow" :key='fld.name'
                  :fld="item[fld.name]"
                  :type="fld.type"
                  @update="item[fld.name] = $event"
                  :label="fld.label"
                  :selectOptions="fld.selectOptions ? fld.selectOptions() : []"
        />
      </div>

      <!--  кнопки   -->
      <comp-item-btn-save @save="save" @cancel="$router.push(docUrl)"/>

    </div>
  </q-page>
</template>

<script>
    import _ from 'lodash'
    import roles from './roles'

    export default {
        props: ['id'],
        computed: {
            docUrl: () => '/users',
        },
        data() {
            return {
                item: null,
                flds: [
                    [
                        {name: 'first_name', type: 'string', label: 'Имя', required: true},
                        {name: 'last_name', type: 'string', label: 'Фамилия', required: true},
                    ],
                    [{
                        name: 'role',
                        type: 'selectMultiple',
                        label: 'Роли',
                        selectOptions: () => this.options
                    }],
                ],
                options: roles,
            }
        },
        methods: {
            save() {
                this.$utils.saveItem.call(this, {
                    method: 'user_update',
                    itemForSaveMod: {role: this.item.role.map(({value}) => value).filter(v => v)},
                    resultModify: (res) => {
                        res.role = res.role.map(roleName => _.find(this.options, {value: roleName})).filter(v => v)
                        return res
                    }
                })
            },
        },
        mounted() {
            let cb = (v) => {
                this.item = v
                // преобразуем роли из строк в объекты
                this.item.role = this.item.role.map(roleName => _.find(this.options, {value: roleName})).filter(v => v)
            }
            this.$utils.getDocItemById.call(this, {method: 'user_get_by_id', cb})
        }
    }
</script>
