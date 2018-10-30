import { mapActions, mapState } from "vuex";

export default {
  name: "register",
  data() {
    return {
      account: "",
      password: "",
      email: "",
      pwd: "",
    };
  },
  computed: {
    ...mapState({
      userInfo: ({ user }) => user.info,
    })
  },
  methods: {
    ...mapActions([
      "userRegister",
    ]),
    async register() {
      const {
        account,
        password,
        email,
        pwd,
      } = this;
      if (!account || !password || !email) {
        this.xError(new Error("用户名、密码与邮箱不能为空"))
        return;
      }
      if (password != pwd) {
        this.xError(new Error("两次输入的密码不一致"));
        return;
      }
      const close = this.xLoading();
      try {
        await this.userRegister({
          account,
          password,
          email,
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
