import {mapActions} from 'vuex';

export default {
  data() {
    return {
      form: {},
    };
  },
  methods: {
    ...mapActions(['userLogin']),
    async submit() {
      const {
        account,
        password,
      } = this.form;
      if (!account || !password) {
        this.$error('account and password can\'t be null');
        return;
      }
      const close = this.$loading();
      try {
        await this.userLogin({
          account,
          password,
        });
        this.$router.back();
      } catch (err) {
        this.$error(err);
      } finally {
        close();
      }
    },
  },
};