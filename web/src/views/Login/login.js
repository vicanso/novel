import { mapActions, mapState } from "vuex";

export default {
  name: "login",
  data() {
    return {
      account: "",
      password: "",
    };
  },
  computed: {
    ...mapState({
      userInfo: ({ user }) => user.info,
    })
  },
  methods: {
    ...mapActions([
      "userLogin",
    ]),
    async login() {
      const {
        account,
        password,
      } = this;
      if (!account || !password) {
        this.xError(new Error("用户名或账号不能为空"))
        return;
      }
      const close = this.xLoading();
      try {
        await this.userLogin({
          account,
          password,
        });
        this.$router.back();
      } catch (err) {
        this.xError(err);
      } finally {
        close();
      }
    },
  }
};
