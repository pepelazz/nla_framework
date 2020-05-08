<template>
    <q-page padding>

        <comp-breadcrumb :list="[{label:'[[index .Vue.I18n "listTitle"]]', to:'/[[.Vue.RouteName]]',  docType: '[[.Name]]'}, {label: item ? (item.title ? item.title : 'Редактирование') : '',  docType: 'edit'}]"/>

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
                [[.Vue.PrintItemTabs]]
            </q-tabs>

            <q-separator />

            <q-tab-panels v-model="tab" :keep-alive="true">
                [[.Vue.PrintItemTabPanels]]
            </q-tab-panels>

        </div>
    </q-page>
</template>

<script>
    [[ .PrintVueImport "docItemWithTabs" ]]
    import queryString from 'query-string'

    export default {
        props: ['id'],
        components: {[[- .Vue.PrintComponents "docItemWithTabs" -]]},
        mixins: [ [[- .Vue.PrintMixins "docItemWithTabs" -]] ],
        computed: {
            docUrl: () => '/[[.Vue.RouteName]]',
        },
        data() {
            return {
                tableName: '[[.Name]]',
                tab: null,
                item: null,
            }
        },
        watch: {
            // смена название таба в url при переключении
            tab(v) {
                this.$utils.updateUrlQuery({tab: v})
            }
        },
        mounted() {
            let cb = (v) => {
                this.item = v
            }
            this.$utils.getDocItemById.call(this, {method: '[[.PgName]]_get_by_id', cb})
            // извлекаем название таба
            const parsedQuery = queryString.parse(location.search)
            this.tab = parsedQuery.tab || 'info'
        },
    }
</script>
