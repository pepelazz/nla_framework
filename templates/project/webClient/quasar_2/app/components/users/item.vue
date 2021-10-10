<template>
  <q-page padding>

    <comp-breadcrumb :list="[{label:'Пользователи', to:'/users', docType: 'users'},
    {label: item ? `${item.last_name} ${item.first_name}` : '',  docType: 'edit'}]"/>

    <div v-if="item" class="q-mt-sm">
      <!--  поля формы    -->
      <div class="row q-col-gutter-md q-mb-sm">
        <div class="col-md-4 col-xs-6">
          <q-input outlined type='text' v-model="item.last_name" :label="$t('user.last_name')" class='q-mb-sm col-md-4 col-xs-6'/>
        </div>
        <div class="col-md-4 col-xs-6">
          <q-input outlined type='text' v-model="item.first_name" :label="$t('user.first_name')" class='q-mb-sm col-md-4 col-xs-6'/>
        </div>
      </div>
      <div class="row q-col-gutter-md q-mb-sm">
        <div class="col-md-4 col-xs-6">
          <q-select outlined multiple v-model="item.role" :label="$t('user.roles')" :options="options"/>
        </div>
        <div class="col-md-4 col-xs-6">
          <q-select outlined v-model="item.state" :label="$t('user.state')"  :options="stateOptions"/>
        </div>
      </div>
      <div class="row q-col-gutter-md q-mb-sm">
        <div class="col-md-4 col-xs-6">
          <q-input outlined type='text' v-model="item.grade" :label="$t('user.grade')" class='q-mb-sm col-md-4 col-xs-6'/>
        </div>
        <div class="col-md-4 col-xs-6">
          <q-input outlined mask="+# (###) ### - ####" v-model="item.phone" :label="$t('user.phone')" :readonly='false'  class='q-mb-sm col-md-4 col-xs-6' ><template v-slot:prepend><q-icon name="phone"/></template></q-input>
        </div>
      </div>
      <!-- аватарка     -->
      <div class="row q-col-gutter-md">
        <div class="col-xs-12 col-sm-6 col-md-4">
          <comp-fld-img :label="$t('user.photo')" :fld="item.avatar" :ext="{fldName: 'avatar', uploadUrl: 'upload_profile_image', methodUpdate: 'user_update', tableId: item.id}"/>
        </div>
        <div class="col-md-4 col-xs-6">
          <q-input outlined  v-model="item.email" :label="$t('user.email')" :readonly='true'  class='q-mb-sm col-md-4 col-xs-6' ><template v-slot:prepend><q-icon name="email"/></template></q-input>
        </div>
      </div>

      <!--  кнопки   -->
      <comp-item-btn-save @save="save" @cancel="$router.push(docUrl)"/>

    </div>
  </q-page>
</template>

<script>
  import _ from 'lodash'
  import roles from './roles'

  // const i18nState = {
  //   waiting_auth: 'ожидает авторизации',
  //   working: 'работает',
  //   fired: 'уволен',
  // }
  export default {
    props: ['id'],
    computed: {
      docUrl: () => '/users',
      // локализация статусов
      stateOptions: function () {
        return ['waiting_auth', 'working', 'fired'].map(v => {
          return {value: v, label:  this.$t(`user.state_${v}`)}
        })
      }
    },
    data() {
      return {
        item: null,
        flds: [
          {name: 'first_name',  required: true},
          {name: 'last_name', required: true},
          {name: 'role'},
          {name: 'state'},
          {name: 'grade'},
          {name: 'phone'},

        ],
        // flds: [
        //     [
        //     ],
        //     [
        //         {
        //             name: 'role',
        //             type: 'selectMultiple',
        //             label: 'роли',
        //             selectOptions: () => this.options
        //         },
        //         {
        //             name: 'state',
        //             type: 'select',
        //             label: 'Текущий статус',
        //             selectOptions: () => this.stateOptions
        //         },
        //     ],
        //     [
        //         {name: 'grade', type: 'string', label: 'Должность'},
        //         {name: 'phone', type: 'string', label: 'Телефон', icon: 'phone', mask: "+# (###) ### - ####"},
        //     ],
        //     // [
        //     //   {name: 'avatar', compName: 'comp-fld-img', label: 'Фото', ext: {fldName: 'avatar', uploadUrl: 'upload_profile_image', methodUpdate: 'current_user_update'}},
        //     // ],
        //     // for codeGenerate #flds_slot
        // ],
        options: roles,
        // for codeGenerate #ptionsFlds: ['state'] - не менять последовательность, это ключ для поиска и добавления новых полей
        optionsFlds: ['state'],
      }
    },
    methods: {
      resultModify(res) {
        res.role = res.role.map(roleName => _.find(this.options, {value: roleName})).filter(v => v)
        if (res.state) res.state = {value: res.state, label: this.$t(`user.state_${res.state}`)}
        return res
      },
      save() {
        this.$utils.saveItem.call(this, {
          method: 'user_update',
          itemForSaveMod: {
            role: this.item.role.map(({value}) => value).filter(v => v),
            state: this.item.state ? this.item.state.value : undefined,
          },
          resultModify: this.resultModify,
        })
      },
    },
    mounted() {
      let cb = (v) => {
        this.item = this.resultModify(v)
      }
      this.$utils.getDocItemById.call(this, {method: 'user_get_by_id', cb})
    }
  }
</script>
