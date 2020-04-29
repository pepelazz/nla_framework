<template>
  <div class="text-grey-8 q-gutter-xs">
    <q-btn size="12px" flat dense round icon="more_vert">
      <q-menu>
        <q-list dense style="min-width: 100px">

          <!--  кнопка редактирования    -->
          <q-item v-if="isEdit" clickable v-close-popup @click="$emit('edit')">
            <q-item-section>редактировать</q-item-section>
          </q-item>

          <!--  кнопка удаления/восстановления        -->
          <q-item v-if="isDelete" clickable v-close-popup @click="confirmItemDeleteRecover(item)">
            <q-item-section v-if="!item.deleted">удалить</q-item-section>
            <q-item-section v-else>восстановить</q-item-section>
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
          title: 'Подтвердите',
            message: `${item.deleted ? 'восстановить' : 'удалить'} : <strong>${itemDetail}</strong>`,
            ok: 'ok',
            cancel: 'отмена',
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
