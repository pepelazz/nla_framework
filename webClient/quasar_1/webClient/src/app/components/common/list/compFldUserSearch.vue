<template>
  <q-select
    outlined
    v-model="localUser"
    use-input
    hide-selected
    fill-input
    input-debounce="300"
    :label='label'
    :options='options'
    @filter="filterFn"
    :readonly="readonly"
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
            user: {
                type: String
            },
            label: {
                type: String
            },
            ext: {
                type: Object
            },
            readonly: {
                type: Boolean,
                default: false
            }
        },
        data() {
            return {
                localUser: {},
                options: []
            }
        },
        created() {
            const fullname = this.user ? this.user : null
            this.localUser = {label: fullname, value: fullname}
        },
        watch: {
            localUser: function (v) {
                if (v && v.user) this.$emit('update', {id: v.user.id, fullname: v.user.fullname})
            }
        },
        methods: {
            filterFn(val, update, abort) {
                update(() => {
                    this.$utils.postCallPgMethod({
                        method: 'user_list',
                        params: Object.assign({search_fullname: val, per_page: 20}, this.ext ? this.ext : {})
                    }).subscribe(res => {
                        if (res.ok) {
                            if (!res.result) res.result = []
                            this.options = res.result.map(v => {
                                return {
                                    label: v.fullname,
                                    value: v.fullname,
                                    user: v,
                                    avatar: v.avatar,
                                }
                            })
                        }
                    })
                })
            },
        },
    }
</script>
