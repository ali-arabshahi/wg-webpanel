<template>
  <div id="app">
    <router-view></router-view>
  </div>
</template>

<script>
import Vue from 'vue';
export default {
  name: 'app',
  components: {},
  data: () => ({}),
  methods: {},
  created() {
    Vue.prototype.$vs = this.$vs;
  }
};
</script>

<style>
  @font-face {
    font-family: "verdana";
    src: url("./assets/fonts/verdana.ttf") format("ttf");
  }

  body {
    font-family: "verdana";
  }
</style>
