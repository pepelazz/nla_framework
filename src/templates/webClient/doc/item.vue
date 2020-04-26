<template>
  <q-page padding>

    <comp-breadcrumb :list="[{label:'[[.NameRu | UpperCaseFirst]]', to:'[[.Vue.RouteName]]'}, {label:'Редактирование'}]"/>

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
[[ .Vue.PrintImport "docItem" ]]
    export default {
        props: ['id'],
        mixins: [ [[- .Vue.PrintMixins "docItem" -]] ],
        computed: {
            docUrl: () => '/city',
        },
        data() {
            return {
                item: null,
                flds: [
                    [
                        {name: 'title', type: 'string', label: 'Название', required: true},

                    ],
                ],
            }
        },
        methods: {
            save() {
                this.$utils.saveItem.call(this, {
                    method: 'city_update',
                    itemForSaveMod: {},
                    errMsgModify(msg) {
                        // локализация ошибки
                        // if (msg.includes('')) return ''
                        return msg
                    }
                })
            },
        },
        mounted() {
            let cb = (v) => {
                this.item = v
            }
            this.$utils.getDocItemById.call(this, {method: 'city_get_by_id', cb})
        }
    }
</script>
