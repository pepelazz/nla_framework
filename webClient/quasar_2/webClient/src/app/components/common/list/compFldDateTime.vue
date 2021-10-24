<template>
  <q-input outlined :label="label" v-model="date" :readonly="readonly" mask="##-##-#### ##:##">
    <template v-slot:prepend  v-if="!readonly">
      <q-icon name="event" class="cursor-pointer">
        <q-popup-proxy ref="qDateProxy" transition-show="scale" transition-hide="scale">
          <q-date v-model="date" mask="DD-MM-YYYY HH:mm" :readonly="readonly">
            <div class="row items-center justify-end">
              <q-btn v-close-popup label="Ok" color="primary" flat />
            </div>
          </q-date>
        </q-popup-proxy>
      </q-icon>
    </template>
    <template v-slot:append v-if="!readonly">
      <q-icon name="access_time" class="cursor-pointer">
        <q-popup-proxy transition-show="scale" transition-hide="scale">
          <q-time v-model="date" mask="DD-MM-YYYY HH:mm" format24h >
            <div class="row items-center justify-end">
              <q-btn v-close-popup label="Ok" color="primary" flat />
            </div>
          </q-time>
        </q-popup-proxy>
      </q-icon>
    </template>
  </q-input>
</template>

<script>
import {date as qDate} from 'quasar'
export default {
  props: ['dateString', 'label', 'readonly'],
  data() {
    return {
      date: null,
    }
  },
  watch: {
    date: function (newVal, oldVal) {
      newVal = newVal?.trim()
      oldVal = oldVal?.trim()
      // console.log('newVal:', newVal, 'oldVal:', oldVal, 'newVal !== oldVal:', newVal.trim() !== oldVal, 'length', newVal.length)
      if (newVal !== oldVal) {
        let d = qDate.extractDate(newVal, 'DD-MM-YYYY HH:mm')
        // обрабатываю дату только если между 1900 и 2050 годами
        if (qDate.isBetweenDates(d, new Date(1900, 1, 1), new Date(2050, 1, 1))) {
          let res = qDate.formatDate(d, 'YYYY-MM-DDTHH:mm:ss')
          this.$emit('update', res)
        }
      }
    },
  },
  mounted() {
    if (this.dateString) this.date = this.dateString
  }
}
</script>
