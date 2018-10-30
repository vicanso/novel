<template lang="pug">
#app
  Home.mainView
  transition
    router-view.childView
</template>
<style lang="sass" src="@/styles/app.sass"></style>

<script>
import Home from "@/views/Home";
import { mapActions } from "vuex";

export default {
  name: "app",
  components: {
    Home
  },
  data() {
    return {
      account: ""
    };
  },
  methods: {
    ...mapActions(["userGetInfo", "userGetSetting"])
  },
  async beforeMount() {
    const close = this.xLoading();
    try {
      await this.userGetInfo();
      await this.userGetSetting();
    } catch (err) {
      this.xError(err);
    } finally {
      close();
    }
  }
};
</script>
