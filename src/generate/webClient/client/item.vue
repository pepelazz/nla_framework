<template>
  <q-page padding>

    <comp-breadcrumb :list="[{label:'Клиент', to:'client'}, {label:'Редактирование'}]"/>

    <div v-if="item" class="q-mt-sm">
      <!--  поля формы    -->
      

      <!--  кнопки   -->
      <comp-item-btn-save @save="save" @cancel="$router.push(docUrl)"/>

    </div>
  </q-page>
</template>

<script>
    export default {
        props: ['id'],
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
