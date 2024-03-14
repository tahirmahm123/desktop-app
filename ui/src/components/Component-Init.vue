<template>
  <div class="main-container login-container">
    <div>
      <div class="card login-card">
        <div class="card-body">
          <spinner :loading="isProcessing" />

          <div v-if="isInitialization" class="main small_text"></div>
          <div class="main" v-else-if="isDaemonInstalling">
            Installing {{ AppName }} Daemon ...
            <div class="small_text" style="margin-top: 10px">
              Please follow the instructions in the dialog
            </div>
          </div>
          <div v-else-if="isConnecting" class="main small_text">
            Connecting ...
          </div>
          <div v-else class="flexColumn">
            <div>
              <div class="large_text">
                Error connecting to {{ AppName }} daemon
              </div>
              <div v-if="daemonIsOldVersionError">
                <div class="small_text">
                  Unsupported {{ AppName }} daemon version v{{
                    currDaemonVer
                  }}
                  (minimum required v{{ minRequiredVer }}).
                </div>
                <div class="small_text">
                  Please update daemon by downloading latest version from
                  <button
                    class="noBordersTextBtn settingsLinkText"
                    v-on:click="visitWebsite"
                  >
                    {{ AppName }} website
                  </button>
                  .
                </div>
              </div>
              <div v-else>
                <div class="small_text">
                  Not connected to daemon. Please, ensure {{ AppName }} daemon
                  is running and try to reconnect.
                </div>
                <div class="small_text">
                  The latest daemon version can be downloaded from
                  <button
                    class="noBordersTextBtn settingsLinkText"
                    v-on:click="visitWebsite"
                  >
                    {{ AppName }} website
                  </button>
                  .
                </div>
              </div>
            </div>
            <div class="text-center">
              <button class="btn btn-primary" v-on:click="ConnectToDaemon">
                Retry ...
              </button>
              <br />
              <button
                class="noBordersTextBtn settingsLinkText"
                v-on:click="visitWebsite"
              >
                {{ AppUrl }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import spinner from "@/components/controls/control-spinner.vue";
import { DaemonConnectionType } from "@/store/types";

const sender = window.ipcSender;
import config from "@/config";

export default {
  components: {
    spinner,
  },
  data: function () {
    return {
      AppName: config.AppName,
      AppUrl: config.AppUrl,
      isProcessing: false,
      isDelayElapsedAfterMount: false,
    };
  },
  mounted() {
    // In order to avoid text blinking, we are showing blank view first few seconds
    // untill 'daemonConnectionState' will not be initialised.
    // The blank view also will be visible first few seconds even after 'daemonConnectionState' was intialized by 'Connecting'
    setTimeout(() => {
      this.isDelayElapsedAfterMount = true;
    }, 3000);
  },
  methods: {
    async ConnectToDaemon() {
      try {
        await sender.ConnectToDaemon();
      } catch (e) {
        console.error(e);
      }
    },
    visitWebsite() {
      sender.shellOpenExternal(this.AppUrl);
    },
  },
  computed: {
    isDaemonInstalling: function () {
      return this.$store.state.daemonIsInstalling;
    },
    isInitialization: function () {
      return (
        (this.$store.state.daemonConnectionState == null &&
          !this.isDaemonInstalling &&
          this.isDelayElapsedAfterMount == false) ||
        (this.isConnecting && this.isDelayElapsedAfterMount == false)
      );
    },
    isConnecting: function () {
      return (
        this.$store.state.daemonConnectionState ===
        DaemonConnectionType.Connecting
      );
    },
    minRequiredVer: function () {
      return config.MinRequiredDaemonVer;
    },
    currDaemonVer: function () {
      return this.$store.state.daemonVersion;
    },
    daemonIsOldVersionError: function () {
      return this.$store.state.daemonIsOldVersionError;
    },
  },
  watch: {
    isConnecting() {
      this.isProcessing = this.isConnecting;
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
@import "@/components/scss/constants";

.login-container .login-card {
  width: 350px;
  left: 410px;
  top: 41px;
  height: 520px;
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

.main {
  padding: 15px;
  margin-top: -100px;
  height: 100%;

  display: flex;
  flex-flow: column;
  justify-content: center;
  align-items: center;
}

.large_text {
  margin: 12px;
  font-weight: 600;
  font-size: 18px;
  line-height: 120%;

  color: #2a394b;
}

.small_text {
  margin: 12px;
  margin-top: 0px;

  font-size: 13px;
  line-height: 17px;
  letter-spacing: -0.208px;

  color: #98a5b3;
}

//.btn {
//  margin: 30px 0 0 0;
//  width: 90%;
//  height: 28px;
//  background: #ffffff;
//  border-radius: 10px;
//  border: 1px solid #7d91a5;
//
//  font-size: 15px;
//  line-height: 20px;
//  text-align: center;
//  letter-spacing: -0.4px;
//  color: #6d849a;
//
//  cursor: pointer;
//}
</style>
