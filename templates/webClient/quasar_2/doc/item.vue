<template>
  <q-page padding>
    [[- if .Vue.Breadcrumb]]
    [[.Vue.Breadcrumb]]
    [[- else]]
    <comp-breadcrumb v-if="!isOpenInDialog" :list="[{label:'[[index .Vue.I18n "listTitle"]]', to:'/[[.Vue.RouteName]]',  docType: '[[.Name]]'}, [[if .IsRecursion]] parentProductBreadcrumb, [[end]] {label: item ? (item.title ? item.title : 'Редактирование') : '',  docType: 'edit'}]"/>
    [[- end]]


    <div v-if="item" class="q-mt-sm">
      <!--  поля формы    -->
      [[ define "vueItemRow" ]]
      <div class="[[.Class]]" [[- if .Fld.Vue.Vif]] v-if="[[.Fld.Vue.Vif]]" [[- end -]]>
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

      [[- if .IsRecursion -]]
      <div class="row q-col-gutter-md q-mb-sm q-mt-sm">
        <div class="col-md-8 col-xs-12" v-if="id !== 'new'">
          <comp-recursive-child-list :id='id' :readonly="[[if .Vue.Readonly -]][[.Vue.Readonly]][[else]]false[[- end]]" @update='save'/>
        </div>
      </div>
      [[end]]

      <!--  кнопки   -->
      <comp-item-btn-save v-if="!isOpenInDialog" @save="save" :readonly="[[.Vue.Readonly]]" @cancel="$router.push(docUrl)"/>
      <!--  при открытии в диалоге кнопку Отмена не показываем   -->
      <q-btn v-else color="secondary" :label="$t('message.save')" class="q-mr-sm" @click="save"/>

        [[range .Vue.Hooks.ItemHtml]]
            [[.]]
        [[- end]]

    </div>
  </q-page>
</template>

<script>
[[ .PrintVueImport "docItem" ]]
    import currentUserMixin from '../../../app/mixins/currentUser'
    export default {
        props: ['id', 'isOpenInDialog' [[- if .IsRecursion -]], 'parent_id'[[- end -]] ],
        components: {[[- .PrintComponents "docItem" -]]},
        mixins: [currentUserMixin, [[- .Vue.PrintMixins "docItem" -]] ],
        computed: {
            docUrl: function() {
              return [[if not .IsRecursion -]]'/[[.Vue.RouteName]]'[[else -]] this.parent_id ? `/[[.Vue.RouteName]]/${this.parent_id}` : '/[[.Vue.RouteName]]' [[- end]]
            },
        },
        data() {
            return {
                item: null,
                flds: [
                    [[- range .Flds]]
                        [[- if .Name]]
                        {name: '[[.Name]]', label: '[[.Vue.NameRu]]'[[if .Vue.IsRequired -]],  required: true[[- end]]},
                        [[- end]]
                    [[- end]]
                ],
                optionsFlds: [ [[- .PrintVueItemOptionsFld -]] ],
                [[if .IsRecursion -]]parentProductBreadcrumb: [], [[ end -]]
                [[ .PrintVueVars "docItem" ]]
            }
        },
        watch: {
          [[.PrintVueItemHookItemWatch]]
        },
        methods: {
          [[ .PrintVueMethods "docItem" ]]
            resultModify(res) {
                [[.PrintVueItemResultModify]]
                return res
            },
            save() {
                [[.PrintVueItemHookBeforeSave]]
                this.$utils.saveItem.call(this, {
                    method: '[[.PgName]]_update',
                    itemForSaveMod: {[[.PrintVueItemForSave]]},
                    resultModify: this.resultModify,
                })
            },
          reload() {
            let cb = (v) => {
              this.item = this.resultModify(v)
            }
            this.$utils.getDocItemById.call(this, {method: '[[.PgName]]_get_by_id', cb})
          }
        },
        mounted() {
           this.reload()
        }
    }
</script>
