<template>
  <q-page padding>

    <comp-breadcrumb :list="[{label:'[[index .Vue.I18n "listTitle"]]', to:'/[[.Vue.RouteName]]'}, {label: item ? (item.title ? item.title : 'Редактирование') : ''}]"/>

    <div v-if="item" class="q-mt-sm">
      <!--  поля формы    -->
      [[ define "vueItemRow" ]]
      <div class="[[.Class]]">
        [[- if gt (len .Grid) 0]]
            [[- range .Grid -]]
                [[ template "vueItemRow" . ]]
            [[- end]]
        [[- else]]
          [[PrintVueFldTemplate .Fld]]
        [[- end]]
      </div>
      [[- end -]]
      [[range .Vue.Grid]]
        [[- template "vueItemRow" .]]
      [[end]]

      <!--  кнопки   -->
      <comp-item-btn-save @save="save" @cancel="$router.push(docUrl)"/>

    </div>
  </q-page>
</template>

<script>
[[ .PrintVueImport "docItem" ]]
    export default {
        props: ['id'],
        mixins: [ [[- .Vue.PrintMixins "docItem" -]] ],
        computed: {
            docUrl: () => '/[[.Vue.RouteName]]',
        },
        data() {
            return {
                item: null,
                flds: [
                    [[- range .Flds]]
                        {name: '[[.Name]]', label: '[[.Vue.NameRu]]'[[if .Vue.IsRequred -]],  required: true[[- end]]},
                    [[- end]]
                ],
                optionsFlds: [ [[- .PrintVueItemOptionsFld -]] ],
            }
        },
        methods: {
          [[ .PrintVueMethods "docItem" ]]
            resultModify(res) {
                [[.PrintVueItemResultModify]]
                return res
            },
            save() {
                this.$utils.saveItem.call(this, {
                    method: '[[.PgName]]_update',
                    itemForSaveMod: {[[.PrintVueItemForSave]]},
                    resultModify: this.resultModify,
                })
            },
        },
        mounted() {
            let cb = (v) => {
                this.item = this.resultModify(v)
            }
            this.$utils.getDocItemById.call(this, {method: '[[.PgName]]_get_by_id', cb})
        }
    }
</script>
