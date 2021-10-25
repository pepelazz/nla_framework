<template>
  <div :class="columnClass">
    <div v-if="type=='separator' && vif" class="text-weight-bold form-separator-header">{{label}}</div>
    <!-- строка, число  -->
    <q-input v-if="type=='string' && vif"
             outlined :type='type' :label="label" :modelValue="fld"
             @update:model-value="v=>$emit('update', v)"
             :readonly="readonly" autogrow/>
    <!-- строка, число  -->
    <q-input v-if="type=='number' && vif"
             outlined :type='type' :label="label" :modelValue="fld"
             @update:model-value="v=>$emit('update', v)"
             :readonly="readonly"/>
    <q-input v-if="type=='phone' && vif"
             outlined :label="label" :modelValue="fld"
             @update:model-value="v=>$emit('update', v)"
             mask="+# (###) ### - ####"
             :readonly="readonly">
      <template v-slot:prepend><q-icon name="phone"/></template>
    </q-input>
    <!-- выбор варианта из списка   -->
    <q-select
      v-if="type=='select' && vif"
      outlined
      :label="label" :modelValue="fld"
      :options="selectOptions"
      @update:model-value="v=>$emit('update', v)"
      :readonly="readonly"
    >
      <template v-if="icon" v-slot:prepend>
        <q-icon :name="icon"/>
      </template>
    </q-select>
    <!-- выбор нескольких вариантов из списка   -->
    <q-select
      v-if="type=='selectMultiple' && vif"
      outlined
      :label="label" :modelValue="fld"
      multiple
      :options="selectOptions"
      @update:model-value="v=>$emit('update', v)"
      :readonly="readonly"
    />
    <!-- выбор пользователя   -->
    <comp-fld-user-search v-if="type=='userId' && vif" :label="label" :user="ajaxSelectTitle" :ext="ext"
                          @update="v=>$emit('update', v.id)"
                          :readonly="readonly"/>
    <!-- выбор ajax-селектора -->
    <comp-fld-ref-search v-if="type=='refSearch' && vif" :label="label" :pgMethod="pgMethod" :itemTitleFldName="itemTitleFldName"
                         @update="v=>$emit('update', v.id)" :readonly="readonly" :ext="ext"/>
    <!-- date   -->
    <comp-fld-date  v-if="type=='date' && vif" :label="label" :date-string="formatDateForSelector(fld)" @update="v=>$emit('update', v)" :readonly="readonly"/>
    <!-- datetime   -->
    <comp-fld-date-time  v-if="type=='datetime' && vif" :label="label" :date-string="formatDateTimeForSelector(fld)" @update="v=>$emit('update', v)" @clear="$emit('clear')" :readonly="readonly"/>

    <div class="q-gutter-sm" v-if="type=='checkbox' && vif">
      <q-checkbox v-model="localFld" :label="label" @update:model-value="v=>$emit('update', v)"/>
    </div>
    <!-- вариант кастомной директивы   -->
    <div v-if="compName && vif">
      <component :is="compName" :label="label" :fld="fld"
                 :type="type"
                 :item="item"
                 :ajaxSelectTitle = "ajaxSelectTitle"
                 :selectOptions="selectOptions"
                 :columnClass="columnClass"
                 :ext="ext"
                 @update="v=>$emit('update', v)"
                 :readonly="readonly"/>
    </div>

  </div>
</template>

<style lang="sass">
  .form-separator-header
    padding: 0.5rem 0
    border-bottom: 1px solid #ccc
    margin: 1rem 0 0
</style>

<script>
    import moment from 'moment'
    export default {
        props: {
            item: {},
            type: {},
            label: {},
            fld: {},
            selectOptions: null,
            ajaxSelectTitle: null,
            ext: null, // дополнительные параметры
            readonly: null,
            columnClass: {
                default: 'col-xs-12 col-sm-6 col-md-4'
            },
            icon: null,
            vif: {
                default: true,
            },
            compName: null,
            pgMethod: null,
            itemTitleFldName: null,
            href: null,
        },
        data() {
            return {
                localFld: null
            }
        },
        methods: {
            formatDateForSelector(d) {
                return d ? moment(d, 'YYYY-MM-DDTHH:mm:ss').format('DD-MM-YYYY') : null
            },
            formatDateTimeForSelector(d) {
                return d ? moment(d, 'YYYY-MM-DDTHH:mm:ss').format('DD-MM-YYYY HH:mm') : null
            },
        },
        mounted() {
            this.localFld = this.fld
        }
    }
</script>
