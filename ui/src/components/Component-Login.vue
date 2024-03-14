<template>
  <div class="main-container login-container">
    <div>
      <div class="card login-card">
        <form class="card-body" @submit.prevent="Login">
          <!-- Logo Start -->
          <div
            class="d-flex mt-2 mb-4 justify-content-center align-items-center"
          >
            <img src="@/assets/img/evolve-vpn-icon.png" alt="" />
            <span class="logo-text"><b>Golden Guard</b>VPN</span>
          </div>
          <!-- Logo End -->

          <div class="welcome-heading mb-4">
            Welcome back! <br />
            Glad to see you, Again!
          </div>

          <div class="mb-3">
            <input
              type="text"
              class="form-control form-control-lg"
              :class="{ 'is-invalid': usernameError }"
              ref="username"
              v-model="username"
              aria-describedby="helpId"
              placeholder="Enter your email"
              :disabled="isProcessing"
              v-on:keyup="keyup($event)"
            />

            <div class="invalid-feedback">{{ usernameErrorText }}</div>
          </div>
          <div class="mb-3">
            <input
              type="password"
              class="form-control form-control-lg"
              :class="{ 'is-invalid': passwordError }"
              v-model="password"
              ref="password"
              aria-describedby="helpId"
              placeholder="Enter your password"
              :disabled="isProcessing"
              v-on:keyup="keyup($event)"
            />
            <div class="invalid-feedback">{{ passwordErrorText }}</div>
          </div>
          <div class="text-end mb-2">
            <a href="#" class="float-right" @click.prevent="ForgetPassword"
              >Forget Password?</a
            >
          </div>
          <div class="d-grid gap-2">
            <button
              type="button"
              name=""
              ref="submitBtn"
              class="btn btn-primary btn-lg"
              :disabled="isProcessing"
              v-on:click="Login"
            >
              <span v-if="isProcessing">Authenticating...</span>
              <span v-else><i class="isax isax-login"></i> Login</span>
            </button>
          </div>
          <div class="mt-4 text-center fw-light">
            Don’t have an account?
            <a href="#" @click.prevent="CreateAccount">Register Now</a>
          </div>
          <div
            class="flexRow leftright_margins"
            style="margin: 20px 20px 0 20px"
          >
            <div
              class="flexRow flexRowRestSpace switcher_small_text"
              style="margin-right: 10px"
            >
              {{ firewallStatusText }}
            </div>

            <SwitchProgress
              :isChecked="this.$store.state.vpnState.firewallState.IsEnabled"
              :isProgress="firewallIsProgress"
              :onChecked="firewallOnChecked"
            />
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
// import spinner from "@/components/controls/control-spinner.vue";
import SwitchProgress from "@/components/controls/control-switch-small2.vue";

import Config from "@/config";

const sender = window.ipcSender;

function processError(e) {
  console.error(e);
  sender.showMessageBox({
    type: "error",
    buttons: ["OK"],
    message: e.toString(),
  });
}

export default {
  props: {
    forceLoginAccount: {
      type: String,
      default: null,
    },
  },
  components: {
    // spinner,
    SwitchProgress,
  },
  data: function () {
    return {
      firewallIsProgress: false,
      username: "",
      password: "",
      isProcessing: false,
      usernameError: false,
      usernameErrorText: "",
      passwordError: false,
      passwordErrorText: "",
      rawResponse: null,
      apiResponseStatus: 0,

      isForceLogoutRequested: false,
    };
  },
  mounted() {
    if (this.$refs.username) this.$refs.username.focus();

    if (this.$route.params.forceLoginAccount == null) {
      if (this.$store.state.settings.isExpectedAccountToBeLoggedIn === true) {
        this.$store.dispatch("settings/isExpectedAccountToBeLoggedIn", false);
        setTimeout(() => {
          sender.showMessageBox({
            type: "info",
            buttons: ["OK"],
            message: `You are logged out.\n\nYou have been redirected to the login page to re-enter your credentials.`,
          });
        }, 0);
      }
    }
  },
  computed: {
    firewallStatusText: function () {
      if (this.$store.state.vpnState.firewallState.IsEnabled)
        return "Firewall enabled and blocking all traffic";
      return "Firewall disabled";
    },
  },
  methods: {
    async Login() {
      try {
        this.isProcessing = true;
        this.usernameError = false;
        this.passwordError = false;
        const resp = await sender.Login(this.username, this.password);
        this.apiResponseStatus = resp.APIStatus;
        let rawResponse = JSON.parse(resp.RawResponse);
        console.log(rawResponse);
        if (this.apiResponseStatus === 422) {
          let errors = rawResponse.errors;
          this.usernameError = "username" in errors || "email" in errors;
          if (this.usernameError) {
            if ("username" in errors) {
              this.usernameErrorText = errors["username"][0];
            } else if ("email" in errors) {
              this.usernameErrorText = errors["email"][0];
            }
          }
          this.passwordError = "password" in errors;
          if (this.passwordError) {
            this.passwordErrorText = errors["password"][0];
          }
        } else {
          let response = rawResponse.response;
          if (!response.auth) {
            throw new Error(`Invalid Credentials Provided`);
          } else if (!response.active) {
            throw new Error(
              `Account has been Deactivated by Provider. Please Contact Support`,
            );
          } else if (response.true) {
            throw new Error(
              `Account has been Expired. Please purchase a plan or renew previous.`,
            );
          } else if (!response.allowLogin) {
            await this.$router.push({
              name: "AccountLimit",
              params: {
                username: this.username,
                totalDevicesAllowed: response.totalSessionsAllowed,
                usedDevices: response.loggedInSessions,
                activeDevices: response.activeSessions,
              },
            });
          } else {
            try {
              await sender.ServerList();
            } catch (e) {
              console.error(e);
            }
          }
        }
      } catch (e) {
        console.error(e);
        sender.showMessageBoxSync({
          type: "error",
          buttons: ["OK"],
          message: "Failed to login",
          detail: `${e}`,
        });
      } finally {
        this.isProcessing = false;
      }
    },
    CreateAccount() {
      sender.shellOpenExternal(Config.RegisterUrl);
    },
    ForgetPassword() {
      sender.shellOpenExternal(Config.ForgetPasswordUrl);
    },
    Cancel() {
      this.rawResponse = null;
      this.apiResponseStatus = 0;
      this.captcha = "";
      this.confirmation2FA = "";
      this.isForceLogoutRequested = false;
    },
    keyup(event) {
      if (event.keyCode === 13) {
        // Cancel the default action, if needed
        event.preventDefault();
        this.Login();
      }
    },
    async firewallOnChecked(isEnabled) {
      this.firewallIsProgress = true;
      try {
        if (
          isEnabled === false &&
          this.$store.state.vpnState.firewallState.IsPersistent
        ) {
          let ret = await sender.showMessageBoxSync(
            {
              type: "question",
              message:
                "The always-on firewall is enabled. If you disable the firewall the 'always-on' feature will be disabled.",
              buttons: ["Disable Always-on firewall", "Cancel"],
            },
            true,
          );

          if (ret == 1) return; // cancel
          await sender.KillSwitchSetIsPersistent(false);
        }

        this.firewallIsProgress = true;
        await sender.EnableFirewall(isEnabled);
      } catch (e) {
        processError(e);
      } finally {
        this.firewallIsProgress = false;
      }
    },
  },
  watch: {},
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="scss" scoped>
.login-container .login-card {
  width: 350px;
  left: 410px;
  top: 41px;
  padding: 10px;
  background: rgba(255, 255, 255, 0.57);
  backdrop-filter: blur(10.1181px);
  border-radius: 8.06751px;
}

.login-container .login-card .logo-text {
  font-weight: 700;
  font-size: 22.7161px;
  line-height: 29px;
  margin-left: 10px;
}

.login-container .login-card .welcome-heading {
  font-weight: 700;
  font-size: 19.65px;
  line-height: 130%;
}

.login-container {
  background: url("~@/assets/img/login-bg-light.png");
}

/* Dark Styling */
.dark .login-container {
  background: url("~@/assets/img/login-bg-dark.png");
  color: #fff;
}

.dark .login-container .login-card {
  background: rgba(0, 0, 0, 0.67);
}
</style>
