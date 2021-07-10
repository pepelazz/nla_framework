<template>
  <q-input :dense='dense' outlined :label="label" v-model="date" :readonly="readonly" mask="##-##-####">
    <template v-slot:append v-if="is_remove">
      <q-btn round flat size="sm" icon="clear" @click="clear"/>
    </template>
    <template v-slot:prepend v-if="!readonly">
      <q-icon name="event" class="cursor-pointer">
        <q-popup-proxy ref="qDateProxy" transition-show="scale" transition-hide="scale">
          <q-date v-model="date" @input="() => $refs.qDateProxy.hide()" mask="DD-MM-YYYY" :readonly="readonly"/>
        </q-popup-proxy>
      </q-icon>
    </template>
  </q-input>
</template>

<script>
  import {date as qDate} from 'quasar'
  // import moment from 'moment'
  export default {
    props: ['dateString', 'label', 'dense', 'is_remove', 'readonly'],
    data() {
      return {
        // date: this.dateString || qDate.formatDate(new Date(), 'DD-MM-YYYY'),
        date: null,
      }
    },
    watch: {
      dateString: function (newVal, oldVal) {
        // обрабатываем только случай когда значение поменялось с null на новое
        // if (newVal && !oldVal && !this.date) this.date = newVal
        if (newVal) this.date = newVal
      },
      date: function (newVal, oldVal) {
        // console.log('newVal:', newVal, 'oldVal:', oldVal, 'newVal !== oldVal:', newVal !== oldVal)
        if (newVal && newVal !== oldVal) {
          let d = qDate.extractDate(newVal, 'DD-MM-YYYY')
          // обрабатываю дату только если между 1900 и 2050 годами
          if (qDate.isBetweenDates(d, new Date(1900, 1, 1), new Date(2050, 1, 1))) {
            let res = qDate.formatDate(d, 'YYYY-MM-DDTHH:mm:ss')
            // let res = moment(newVal, 'DD-MM-YYYY')
            this.$emit('update', res)
          }
        }
      },
    },
    methods: {
      clear() {
        this.$emit('clear')
        this.date = null
      }
    },
    mounted() {
      if (this.dateString) this.date = this.dateString
    }
  }
</script>
