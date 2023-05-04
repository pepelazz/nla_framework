<template>
  <div class="text-grey-8 q-gutter-xs">
    <q-btn size="12px" flat dense round icon="more_vert">
      <q-menu>
        <q-list dense style="min-width: 100px">

          <!--  кнопка редактирования    -->
          <q-item v-if="isEdit" clickable v-close-popup @click="$emit('edit')">
            <q-item-section>{{$t('message.edit')}}</q-item-section>
          </q-item>

          <!--  кнопка удаления/восстановления        -->
          <q-item v-if="isDelete" clickable v-close-popup @click="confirmItemDeleteRecover(item)">
            <q-item-section v-if="!item.deleted">{{$t('message.delete')}}</q-item-section>
            <q-item-section v-else>{{$t('message.restore')}}</q-item-section>
          </q-item>

        </q-list>
      </q-menu>
    </q-btn>
  </div>
</template>

<script>
  export default {
    props: ['item', 'itemProp', 'isDelete', 'isEdit', 'pgMethod'],
    methods: {
      confirmItemDeleteRecover(item) {
        // let that = this
        const itemDetail = item[this.itemProp] ? item[this.itemProp] : ''
        this.$q.dialog({
            title: this.$t('message.dialog_delete_title'),
            message: `${item.deleted ? this.$t('message.restore') : this.$t('message.delete') } : <strong>${itemDetail}</strong>`,
            ok: this.$t('message.ok'),
            cancel: this.$t('message.cancel'),
            html: true,
        }).onOk(() => {
          this.$utils.postCallPgMethod({
            method: this.pgMethod,
            params: {id: item.id, deleted: !item.deleted}
          }).subscribe(v => {
            this.$emit('reload-list')
          })
        })
      },
    }
  }
</script>
