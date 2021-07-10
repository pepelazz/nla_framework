<template>
  <div>
    <q-breadcrumbs gutter="sm">
      <q-breadcrumbs-el :label="$config.uiAppName" to="/"/>
      <!-- делаем плоский массив, потому что в нем может быть вложенный массив     -->
      <q-breadcrumbs-el v-for='item in [].concat(...list)' :key="item.label" :label="item.label" :to="item.to" :icon="icon(item)"/>
    </q-breadcrumbs>
    <q-separator class="q-mt-sm"/>
  </div>
</template>

<script>
    export default {
        props: ['list'],
        computed: {
            icon() {
                return function(item) {
                    if (item.docType === 'edit') return 'edit'
                    return item.icon || this.$config.breadcrumbIcons[item.docType] || 'label'
                }
            }
        }
    }
</script>
