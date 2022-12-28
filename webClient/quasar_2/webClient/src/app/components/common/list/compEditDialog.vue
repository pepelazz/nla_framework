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
            <q-input v-if="(fld.type === 'text' || fld.type === 'number') && fld.vif(item)" :label="fld.label" v-model="item[fld.name]" :type="fld.type" outlined :readonly="fld.readonly"/>
            <q-select v-if="fld.type === 'select' && fld.vif(item)" :label="fld.label" v-model="item[fld.name]" :options="fld.options" outlined :readonly="fld.readonly"/>
            <comp-fld-ref-search v-if="fld.type === 'ref' && fld.vif(item)" outlined :pgMethod="fld.pgMethod" :ext="fld.ext || {}" :label="fld.label" :item='item[fld.name + "_title"]' :itemId='item[fld.name]'  @update="v=> fld.updateFn ? fld.updateFn(item, v) :item[fld.name] = v.id" @clear="item[fld.name] = null" :readonly="fld.readonly"/>
            <comp-fld-date v-if="fld.type === 'date' && fld.vif(item)" :label="fld.label" outlined  :date-string="$utils.formatPgDate(item[fld.name])" @update="v=> item[fld.name] = v" :readonly="fld.readonly"/>
            <q-checkbox v-if="fld.type === 'checkbox' && fld.vif(item)" :label="fld.label" v-model="item[fld.name]" />
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
    let beforeSaveCb
    const show = (d) => {
      isShowAddDialog.value = true
      item.value = d.item
      flds.value = d.flds.map(v => v.map(fld => {
        if (!fld.vif) fld.vif = () => true
        return fld
      }))
      beforeSaveCb = d.beforeSaveCb
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
      // если указан модификатор beforeSaveCb, то выполняем вызов функции
      if (beforeSaveCb) beforeSaveCb(itemForSave)
      // сохраняем в базу
      if (props.pgMethod) {
        $utils.callPgMethod(props.pgMethod, itemForSave, (res) => {
          isShowAddDialog.value = false
          emit('update', res)
        })
      } else {
        isShowAddDialog.value = false
        emit('update', itemForSave)
      }
    }

    return {
      item, isShowAddDialog, show, flds, save,
    }
  }
}
</script>
