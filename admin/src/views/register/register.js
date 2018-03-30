import {mapActions} from 'vuex';

export default {
  data() {
    return {
      form: {},
    };
  },
  methods: {
    ...mapActions(['userRegister']),
    async submit() {
      const {
        account,
        password,
        email,
        confirmPassword,
      } = this.form;
      if (!account || !password || !email) {
        this.$error('account, email and password can\'t be null');
        return;
      }
      if (password !== confirmPassword) {
        this.$error('the password isn\' the same');
        return;
      }
      const close = this.$loading();
      try {
        await this.userRegister({
          account,
          password,
          email,
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
