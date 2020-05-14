export default {
    data () {
        return {
            [[GetFld]]FilterOptions: [],
            [[GetFld]]Options: [],
        }
    },
    methods: {
        [[GetFld]]CreateValue (val, done) {
            if (val.length > 0) {
                if (!this.[[GetFld]]Options.includes(val)) {
                    this.[[GetFld]]Options.push(val)
                }
                done(val, 'toggle')
            }
        },
        [[GetFld]]FilterFn (val, update) {
            update(() => {
                if (val === '') {
                    this.[[GetFld]]FilterOptions = this.[[GetFld]]Options
                }
                else {
                    const needle = val.toLowerCase()
                    this.[[GetFld]]FilterOptions = this.[[GetFld]]Options.filter(
                        v => v.toLowerCase().indexOf(needle) > -1
                    )
                }
            })
        }
    },
    mounted() {
        this.$utils.postCallPgMethod({method: '[[GetDoc]]_[[GetFld]]_list', params: {}}).subscribe(res => {
            if (res.ok) {
                this.[[GetFld]]Options = res.result
            }
        })
    },
}
