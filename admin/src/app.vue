<template lang="pug">
  #app
    main-nav
    main-header
    transition
      .contentWrapper: router-view
</template>

<script>
import {mapActions} from 'vuex';

import MainNav from './components/main-nav';
import MainHeader from './components/main-header';
export default {
  components: {
    MainNav,
    MainHeader,
  },
  methods: {
    ...mapActions(['userGetInfo']),
  },
  async beforeMount() {
    const close = this.$loading();
    try {
      await this.userGetInfo();
    } catch (err) {
      this.$error(err);
    } finally {
      close();
    }
  },
};
</script>

<style src="./app.sass" lang="sass"></style>
