<template>
  <div class="main-container">
    <div class="home-layout">
      <div>
        <Sidebar active="settings" />
        <div class="body-content px-4">
          <div>
            <h4 class="my-3">Settings</h4>
            <div class="row" style="flex-grow: 1">
              <div id="tabsTitle">
                <!--                <button
                  class="noBordersBtn tabTitleBtn"
                  v-on:click="onView('account')"
                  v-bind:class="{
                    activeBtn: view === 'account',
                  }"
                >
                  Account
                </button>-->

                <button
                  class="noBordersBtn tabTitleBtn"
                  v-on:click="onView('general')"
                  v-bind:class="{
                    activeBtn: view === 'general',
                  }"
                >
                  General
                </button>

                <button
                  class="noBordersBtn tabTitleBtn"
                  v-on:click="onView('connection')"
                  v-bind:class="{
                    activeBtn: view === 'connection',
                  }"
                >
                  Connection
                </button>
                <button
                  class="noBordersBtn tabTitleBtn"
                  v-on:click="onView('firewall')"
                  v-bind:class="{
                    activeBtn: view === 'firewall',
                  }"
                >
                  Kill Switch
                </button>
                <button
                  v-if="isSplitTunnelVisible"
                  class="noBordersBtn tabTitleBtn"
                  v-on:click="onView('splittunnel')"
                  v-bind:class="{
                    activeBtn: view === 'splittunnel',
                  }"
                >
                  Split Tunnel
                </button>
                <!-- <button
                  v-if="isLoggedIn"
                  class="noBordersBtn tabTitleBtn"
                  v-on:click="onView('antitracker')"
                  v-bind:class="{
                    activeBtn: view === 'antitracker',
                  }"
                >
                  Anti Tracker
                </button>
                <button
                  v-if="isLoggedIn"
                  class="noBordersBtn tabTitleBtn"
                  v-on:click="onView('dns')"
                  v-bind:class="{
                    activeBtn: view === 'dns',
                  }"
                >
                  DNS
                </button> -->
              </div>
            </div>
          </div>
          <div class="flexColumn" v-if="view === 'connection'">
            <connectionView />
          </div>
          <!--          <div class="flexColumn" v-else-if="view === 'account'">
            <accountView />
          </div>-->
          <div class="flexColumn" v-else-if="view === 'general'">
            <generalView />
          </div>
          <div class="flexColumn" v-else-if="view === 'firewall'">
            <firewallView />
          </div>
          <div class="flexColumn" v-else-if="view === 'splittunnel'">
            <splittunnelView />
          </div>
          <!-- <div class="flexColumn" v-else-if="view === 'networks'">
            <networksView />
          </div>
          <div class="flexColumn" v-else-if="view === 'antitracker'">
            <antitrackerView />
          </div>
          <div class="flexColumn" v-else-if="view === 'dns'">
            <dnsView
              :registerBeforeCloseHandler="doRegisterBeforeViewCloseHandler"
            />
          </div> -->
          <div class="flexColumn" v-else>
            <!-- no view defined -->
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
const sender = window.ipcSender;

import connectionView from "@/components/settings/settings-connection.vue";
// import accountView from "@/components/settings/settings-account.vue";
import generalView from "@/components/settings/settings-general.vue";
import firewallView from "@/components/settings/settings-firewall.vue";
import splittunnelView from "@/components/settings/settings-splittunnel.vue";
// import networksView from "@/components/settings/settings-networks.vue";
// import antitrackerView from "@/components/settings/settings-antitracker.vue";
// import dnsView from "@/components/settings/settings-dns.vue";
// import imgArrowLeft from "@/components/images/arrow-left.vue";
import Sidebar from "@/components/Component-Sidebar.vue";

export default {
  components: {
    connectionView,
    // accountView,
    generalView,
    firewallView,
    splittunnelView,
    // networksView,
    // antitrackerView,
    // dnsView,
    // imgArrowLeft,
    Sidebar,
  },
  mounted() {
    this.onBeforeViewCloseHandler = null;
    if (this.$route.params.view != null) this.view = this.$route.params.view;
    this.$store.dispatch("uiState/currentSettingsViewName", this.view);
  },
  data: function () {
    return {
      view: "general",
      // Handler which will be called before closing current view (null - in case if no handler registered for current view).
      // Handler MUST be 'async' function and MUST return 'true' to allow to switch current view
      onBeforeViewCloseHandler: Function,
    };
  },
  computed: {
    isSplitTunnelVisible() {
      return this.$store.getters["isSplitTunnelEnabled"];
    },
    versionSingle: function () {
      if (this.versionDaemon === this.versionUI) return this.versionDaemon;
      return null;
    },
    versionDaemon: function () {
      let v = this.$store.state.daemonVersion;
      if (!v) return "version unknown";
      return `v${v}`;
    },
    versionUI: function () {
      let v = sender.appGetVersion();
      if (!v) return "version unknown";
      return `v${v}`;
    },
  },
  methods: {
    goBack: async function () {
      if (this.$store.state.settings.minimizedUI) {
        sender.closeCurrentWindow();
      } else {
        // Call async 'BeforeViewCloseHandler' for current view (if exists). Block view change if handler return != true
        if (this.onBeforeViewCloseHandler != null) {
          if ((await this.onBeforeViewCloseHandler()) !== true) return;
        }

        this.$router.push("/");
      }

      this.onBeforeViewCloseHandler = null; // forget 'onBeforeViewCloseHandler' for current view
      this.$store.dispatch("uiState/currentSettingsViewName", null);
    },
    onView: async function (viewName) {
      // Call async 'BeforeViewCloseHandler' for current view (if exists). Block view change if handler return != true
      if (this.onBeforeViewCloseHandler != null) {
        if ((await this.onBeforeViewCloseHandler()) !== true) return;
      }

      this.onBeforeViewCloseHandler = null; // forget 'onBeforeViewCloseHandler' for current view
      this.view = viewName;
      this.$store.dispatch("uiState/currentSettingsViewName", this.view);
    },
    doRegisterBeforeViewCloseHandler: function (handler) {
      // Register handler which will be called before closing current view
      // Handler MUST be 'async' function and MUST return 'true' to allow to switch current view
      this.onBeforeViewCloseHandler = handler;
    },
  },
};
</script>

<style scoped lang="scss">
@import "@/components/scss/constants";

$back-btn-width: 50px;
$min-title-height: 26px;

div.row {
  display: flex;
  flex-direction: row;
  width: 100%;
}

#main {
  height: 100%;

  font-size: 13px;
  line-height: 16px;
  letter-spacing: -0.58px;
}

#leftPanel {
  padding-top: 50px;
  background: var(--background-color-alternate);
  min-width: 232px;
  max-width: 232px;
  height: 100vh;
}

#leftPanelHeader {
  padding-bottom: 38px;
}

#tabsTitle {
  align-items: center;
  flex-direction: row;
  margin-bottom: 10px;
  display: flex;
  overflow: auto;
}

.rightPanel {
  margin: 20px;
  width: 100vw;
}

.rightPanel * {
}

#backBtn {
  min-width: $back-btn-width;
  max-width: $back-btn-width;

  display: flex;
  justify-content: center;
  align-items: center;
}

.Header {
  font-style: normal;
  font-weight: 800;
  font-size: 24px;
  line-height: 29px;
  letter-spacing: -0.3px;
  text-transform: capitalize;
}

button.noBordersBtn {
  border: none;
  background-color: inherit;
  outline-width: 0;
  cursor: pointer;
  width: auto;
}

button.tabTitleBtn {
  display: flex;
  padding: 6px 10px;
  font-size: 14px;
  line-height: 17px;

  color: var(--text-color-settings-menu);
}

button.activeBtn {
  font-weight: 500;
  color: #ebc553;
}

div.version {
  color: gray;
}

.form-check-input:checked {
  background-color: #ebc553 !important;
  border-color: #ebc553 !important;
}
</style>
