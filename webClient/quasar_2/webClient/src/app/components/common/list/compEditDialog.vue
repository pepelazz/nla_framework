<template>
  <!-- ДИАЛОГ ДОБАВЛЕНИЯ/РЕДАКТИРОВАНИЯ  -->
  <q-dialog v-if="item" v-model="isShowAddDialog" :persistent="item.id !== -1">
    <q-card style="width: 700px; max-width: 80vw;">

      <!-- ЗАГОЛОВОК     -->
      <q-card-section v-if="item.id === -1 && labelNew">
        <div class="text-h6" >{{ labelNew }}</div>
      </q-card-section>
      <q-card-section v-if="item.id > 0 && labelEdit">
        <div class="text-h6" >{{ labelEdit }}</div>
      </q-card-section>

      <q-card-section class="q-pt-none">
        <div v-for="fldRow in flds" class="row q-col-gutter-md q-mb-sm">
          <div v-for="fld in fldRow" :class="fld.classCol || 'col-12'">
            <q-input v-if="fld.type === 'text' || fld.type === 'number'" :label="fld.label" v-model="item[fld.name]" :type="fld.type" outlined/>
            <q-select v-if="fld.type === 'select'" :label="fld.label" v-model="item[fld.name]" :options="fld.options" outlined/>
            <comp-fld-ref-search v-if="fld.type === 'ref'" outlined :pgMethod="fld.pgMethod" :ext="fld.ext || {}" :label="fld.label" :item='item[fld.name + "_title"]' :itemId='item[fld.name]'  @update="v=> item[fld.name] = v.id" @clear="item[fld.name] = null"/>
            <comp-fld-date v-if="fld.type === 'date'" outlined  :date-string="$utils.formatPgDate(item[fld.name])" @update="v=> item[fld.name] = v" />
          </div>
        </div>
      </q-card-section>

      <q-card-actions align="right" class="bg-white text-teal">
        <q-btn flat label="Отмена" @click="isShowAddDialog=false"/>
        <q-btn flat label="Сохранить" @click="save"/>
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script>
  import {ref} from 'vue'
  import $utils from 'src/app/plugins/utils'
  import _ from 'lodash'
  import { useQuasar } from 'quasar'

  export default {
    props: ['pgMethod', 'labelNew', 'labelEdit'],
    emits: ['update'],
    setup(props, {emit}) {
      const $q = useQuasar()
      const item = ref(null)
      const isShowAddDialog = ref(false)
      const flds = ref([])
      const show = (d) => {
        isShowAddDialog.value = true
        item.value = d.item
        flds.value = d.flds
      }
      const save = () => {
        const itemForSave = Object.assign({}, item.value)
        // проверка на required
        let isNotAllFilled = false
        _.flattenDeep(flds.value).filter(v => v.required).map(v => {
          if (!itemForSave[v.name]) {
            isNotAllFilled = true
            $q.notify({type: 'negative', message: `не заполнено поле: ${v.label}`})
          }
        })
        if (isNotAllFilled) return
        // для поля select преобразуем {label: '', value: ''} -> value
        _.flattenDeep(flds.value).filter(v => v.type === 'select').map(v => itemForSave[v.name] = itemForSave[v.name].value)
        // сохраняем в базу
        $utils.callPgMethod(props.pgMethod, itemForSave, (res) => {
          isShowAddDialog.value = false
          emit('update', res)
        })
      }

      return {
        item, isShowAddDialog, show, flds, save,
      }
    }
  }
</script>
