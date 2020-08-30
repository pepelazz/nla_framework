<template>
  <div [[Vif]]>
    <q-btn color="primary" outline label="[[GetLabel]]" [[if GetIconSrc]]icon="[[GetIconSrc]]" [[end]] @click="open"/>
    <q-dialog v-model="isShowDialog" persistent>
      <q-card style="width: 700px; max-width: 80vw;">
        <q-card-section>
          <div class="row q-col-gutter-md q-mb-sm">
            <div class="text-h6">[[UpperCaseFirst GetLabel]]</div>
          </div>
          [[range GetUpdateFldsGrid]]
          [[- if .]]<div class="row q-col-gutter-md q-mb-sm">
            [[range .]]
            <div class='[[printf "%v" .Vue.ClassPrint]]'>
            [[PrintVueFldTemplate .]]
            </div>
            [[end]]
          </div>
          [[- end -]]
          [[end]]
        </q-card-section>
        <q-card-actions align="right" class="text-primary">
          <q-btn flat label="Отмена" v-close-popup />
          <q-btn flat label="Ок" @click="action"/>
        </q-card-actions>
      </q-card>
    </q-dialog>
  </div>
</template>

<script>
    import isRole from '../../../mixins/isRole'
    export default {
        props: ['item', 'currentUser'],
        mixins: [isRole],
        computed: {
          id: function () {
            return this.item?.id
          }
        },
        data() {
            return {
                isShowDialog: false,
                isReadonly: false, // заглушка для корректного отображения полей
            }
        },
        methods: {
            open() {
                this.isShowDialog = true
            },
            action() {
              [[range GetUpdateFlds -]]
                      [[- if .Vue.IsRequired]]
              if (!this.item.[[.Name]]) {
                this.$q.notify({
                  message: 'Не заполнено поле "[[.NameRu]]"',
                  type: 'negative',
                  position: 'top-right'
                })
                return
              }
                      [[- end -]]
              [[- end]]
                this.$utils.postCallPgMethod({method: '[[.Name]]_action', params: Object.assign(this.item, {action_name: '[[GetActionName]]', [[range GetUpdateFlds]] [[.Name]]: this.item.[[.Name]], [[end]]})}).subscribe(res => {
                    if (res.ok) {
                        this.isShowDialog = false
                        this.$emit('stateChanged')
                    }
                })
            }
        }
    }
</script>
