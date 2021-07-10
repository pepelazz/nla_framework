<template>
  <q-select
    outlined
    v-model="localUser"
    use-input
    fill-input
    emit-value
    map-options
    multiple
    use-chips
    :label='label'
    :options='options'
  >
    <template v-slot:no-option>
      <q-item>
        <q-item-section class="text-grey">
          No results
        </q-item-section>
      </q-item>
    </template>
  </q-select>
</template>

<script>
    export default {
        props: {
            userIds: {
                type: Array,
            },
            label: {
                type: String
            },
        },
        data() {
            return {
                localUser: [],
                options: [],
                filterOptions: [],
            }
        },
        watch: {
            localUser: function (v) {
                if (v) this.$emit('update', v)
            }
        },
        methods: {
            filterFn (val, update) {
                update(() => {
                    if (val === '') {
                        this.filterOptions = this.options
                    } else {
                        const needle = val.toLowerCase()
                        this.filterOptions = this.options.filter(v => v.label.toLowerCase().indexOf(needle) > -1)
                    }
                })
            }
        },
        mounted() {
            this.$utils.postCallPgMethod({
                method: 'user_list',
                params: Object.assign({per_page: 1000, order_by: 'fullname asc'})
            }).subscribe(res => {
                if (res.ok) {
                    if (!res.result) res.result = []
                    this.options = res.result.map(v => {
                        return {
                            label: `${v.fullname}`,
                            value: v.id,
                        }
                    })
                    if (this.userIds) this.localUser = this.userIds
                    this.$forceUpdate()
                }
            })
        }
    }
</script>
