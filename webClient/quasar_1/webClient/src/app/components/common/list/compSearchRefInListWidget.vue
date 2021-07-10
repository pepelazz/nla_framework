<template>
  <q-select
    outlined
    v-model="item"
    use-input
    hide-selected
    fill-input
    input-debounce="300"
    :label='label'
    :options='options'
    @filter="filterFn"
  >
    <template v-slot:no-option>
      <q-item>
        <q-item-section class="text-grey">
          Ничего не найдено
        </q-item-section>
      </q-item>
    </template>
  </q-select>
</template>

<script>
  import _ from 'lodash'

  export default {
    props: {
      label: {},
      pgMethod: {},
      title: {},
      ext: {},
      existIds: {}, // список уже выбранных, чтобы отфильтровать из списка загруженных
    },
    data() {
      return {
        item: null,
        options: []
      }
    },
    watch: {
      item(v) {
        if (v) this.$emit('update', v.value)
      }
    },
    methods: {
      filterFn(val, update, abort) {
        update(() => {
          this.loadData(val)
        })
      },
      loadData(search_text = '') {
        let params = {search_text, per_page: 20}
        if (this.ext) params = Object.assign({search_text, per_page: 20}, this.ext)
        this.$utils.postCallPgMethod({
          method: this.pgMethod,
          params
        }).subscribe(res => {
          if (res.ok) {
            if (!res.result) res.result = []
            // отфильтровываем уже выбранные элементы
            if (this.existIds && Array.isArray(this.existIds)) res.result = res.result.filter(v => !this.existIds.includes(v.id))
            this.options = res.result.map(v => {
              let label = v[this.title ? this.title : 'title']
              if (_.isFunction(this.title)) label = this.title(v)
              if (!label) label = 'неверно указан параметр для label'
              return {
                label,
                value: v,
              }
            })
          }
        })
      }
    },
    mounted() {
      this.loadData()
    }
  }
</script>
