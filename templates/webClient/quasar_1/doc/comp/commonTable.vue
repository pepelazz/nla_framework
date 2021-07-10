<template>
    <q-table
            [[if GetTableTitle]]title="[[GetTableTitle]]"[[end]]
            :data="list"
            :columns="columns"
            row-key="name"
            :pagination.sync="pagination"
            separator="[[GetSeparator]]"
    />
</template>
<script>
    export default {
        data() {
            return {
                pagination: {
                    rowsPerPage: [[GetRowsPerPage]]
                },
                columns: [
                    [[range GetColumns -]]
                    { name: '[[.Name]]', align: '[[.Align]]', label: '[[.Label]]', field: '[[.Field]]', sortable: [[.Sortable]] },
                    [[ end]]
                ],
                list: [],
            }
        },
        methods: {
            reload() {
                this.$utils.postCallPgMethod([[GetPgMethod]]).subscribe(res => {
                    if (res.ok) this.list = res.result
                })

            }
        },
        mounted() {
            this.reload()
        }
    }
</script>
