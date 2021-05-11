<template>
  <q-layout view="lHh Lpr lFf">
    <q-header elevated> </q-header>
    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script>
export default {
  name: "MainLayout",
  components: {},
  data() {
    return {};
  }
};
</script>
