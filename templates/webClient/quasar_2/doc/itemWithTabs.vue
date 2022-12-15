<template>
    <q-page padding>

        <comp-breadcrumb v-if="!isOpenInDialog" :list="[{label:'[[index .Vue.I18n "listTitle"]]', to:'/[[.Vue.RouteName]]',  docType: '[[.Name]]'}, [[if .IsRecursion]] parentProductBreadcrumb, [[end]] {label: item ? (item.title ? item.title : 'Редактирование') : '',  docType: 'edit'}]"/>

        <div v-if="item" class="q-mt-sm">
            <q-tabs
                    v-model="tab"
                    dense
                    class="text-grey"
                    active-color="primary"
                    indicator-color="primary"
                    align="left"
                    narrow-indicator
            >
                [[.PrintVueItemTabs]]
            </q-tabs>

            <q-separator />

            <q-tab-panels v-model="tab">
                [[.PrintVueItemTabPanels]]
            </q-tab-panels>

        </div>
    </q-page>
</template>

<script>
    [[ .PrintVueImport "docItemWithTabs" ]]
    import queryString from 'query-string'

    export default {
        props: ['id', 'isOpenInDialog' [[- if .IsRecursion -]], 'parent_id'[[- end -]]],
        components: {[[- .PrintComponents "docItemWithTabs" -]]},
        mixins: [ [[- .Vue.PrintMixins "docItemWithTabs" -]] ],
        computed: {
            docUrl: function() {
                return [[if not .IsRecursion -]]'/[[.Vue.RouteName]]'[[else -]] this.parent_id ? `/[[.Vue.RouteName]]/${this.parent_id}` : '/[[.Vue.RouteName]]' [[- end]]
            },
        },
        data() {
            return {
                tableName: '[[.Name]]',
                tab: null,
                item: null,
                [[if .IsRecursion -]]parentProductBreadcrumb: [], [[ end -]]
            }
        },
        watch: {
            // смена название таба в url при переключении
            tab(v) {
                // если открыли в даилоге, то название табов в url не меняем
                if (!this.isOpenInDialog) this.$utils.updateUrlQuery({tab: v})
            }
        },
        methods: {
            resultModify(res) {
                [[.PrintVueItemResultModify]]
                return res
            },
        },
        mounted() {
            let cb = (v) => {
                this.item = this.resultModify(v)
            }
            this.$utils.getDocItemById.call(this, {method: '[[.PgName]]_get_by_id', cb})
            // извлекаем название таба только в случае открытия не в диалоге
            if (!this.isOpenInDialog) {
                const parsedQuery = queryString.parse(location.search)
                this.tab = parsedQuery.tab || 'info'
            } else {
                this.tab = 'info'
            }
        },
    }
</script>
