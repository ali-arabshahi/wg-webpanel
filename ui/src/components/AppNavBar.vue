<template>
  <header>
    <vs-navbar :collapse="false"  id="appnavID" v-model="activeItem" class="nabarx">
      <div slot="title">
         <vs-button type="line" icon="density_small" @click="toggle"></vs-button>
        <!-- <vs-navbar-title> Hello world </vs-navbar-title> -->
      </div>

      <!-- <vs-navbar-item index="0">
        <a href="#">Home</a>
      </vs-navbar-item>
      <vs-navbar-item index="1">
        <a href="#">News</a>
      </vs-navbar-item>
      <vs-navbar-item index="2">
        <a href="#">Update</a>
      </vs-navbar-item> -->
    </vs-navbar>
  </header>
</template>

<script>
export default {
  name: "AppNavBar",
  data: () => ({
    activeItem: 0,
  }),
   methods: {
    toggle() {
      this.$store.dispatch("actionSidebare");
    },
  },
};
</script>
<style scoped>
/* #appnavID {
  position: fixed;
} */
</style>
