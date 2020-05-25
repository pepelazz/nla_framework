<template>
  <div>
    <div v-if="item">
      <comp-dadata-address-dialog ref="addressDialog" @update="updateAddress"/>
      <div class="row q-col-gutter-md q-mb-sm" >
        <div class="col-12">
          <q-input outlined label="Адрес" v-model="item.full">
            <template v-slot:prepend v-if="item.geo_lat">
              <q-icon name="place" @click="openOnMap"><q-tooltip>показать на карте</q-tooltip></q-icon>
            </template>
            <template v-slot:append>
              <q-icon name="search" @click="$refs.addressDialog.open()" />
            </template>
          </q-input>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
    export default {
        props: {
            fld: {},
            label: {},
            readonly: null,
        },
        data() {
            return {
                item: null,
            }
        },
        watch: {
            item: {
                handler(v) {
                    this.$emit('update', v)
                },
                deep: true,
            }
        },
        methods: {
            updateAddress(addr) {
                let addrValue, addrData
                // в ранних версиях возвращался массив, в новых уже объект. Поэтому проверяем переданный аргумент
                if (Array.isArray(addr)) {
                    if (!addr[0]) return
                    addrValue = addr[0].value
                    addrData = addr[0].data
                } else {
                    addrValue = addr.value
                    addrData = addr.data
                }
                const address_full = addrValue
                const {city, region, area, street, postal_code, house, place, geo_lat, geo_lon} = addrData
                this.item.full = address_full
                this.item.geo_lat = geo_lat
                this.item.geo_lon = geo_lon
                this.item.source = {city, region, area, street, postal_code, house, place}
            },
            openOnMap() {
                let win = window.open(`http://maps.yandex.ru/?text=${this.item.geo_lat},${this.item.geo_lon}`, '_blank')
                win.focus()
            }
        },
        mounted() {
            this.item = this.fld || {full: null}
        }
    }
</script>
