<template>
  <div>
    <q-breadcrumbs gutter="sm">
      <q-breadcrumbs-el label="4PL" to="/"/>
      <q-breadcrumbs-el v-for='item in list' :key="item.label" :label="item.label" :to="item.to"/>
    </q-breadcrumbs>
    <q-separator class="q-mt-sm"/>
  </div>
</template>

<script>
  export default {
    props: ['list']
  }
</script>
