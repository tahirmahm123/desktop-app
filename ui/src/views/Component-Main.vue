<template>
  <div class="main-container">
    <div class="home-layout">
      <div>
        <Sidebar v-if="showSidebar" active="connect" />
        <transition mode="out-in" name="fade">
          <component
            v-bind:is="currentViewComponent"
            :onConnectionSettings="onConnectionSettings"
            :onDefaultView="onDefaultLeftView"
            :onWifiSettings="onWifiSettings"
          ></component>
        </transition>
      </div>
    </div>
  </div>
</template>

<script>
const sender = window.ipcSender;

import { DaemonConnectionType } from "@/store/types";
import { IsWindowHasFrame } from "@/platform/platform";
import Init from "@/components/Component-Init.vue";
import Control from "@/components/Component-Control.vue";
import Sidebar from "@/components/Component-Sidebar.vue";

export default {
  components: {
    Init,
    Control,
    Sidebar,
  },
  computed: {
    isWindowHasFrame: function () {
      return IsWindowHasFrame();
    },
    isLoggedIn: function () {
      return this.$store.getters["account/isLoggedIn"];
    },
    currentViewComponent: function () {
      const daemonConnection = this.$store.state.daemonConnectionState;
      if (
        daemonConnection == null ||
        daemonConnection === DaemonConnectionType.NotConnected ||
        daemonConnection === DaemonConnectionType.Connecting
      )
        return Init;
      return Control;
    },
    showSidebar: function () {
      const daemonConnection = this.$store.state.daemonConnectionState;
      return !(
        daemonConnection == null ||
        daemonConnection === DaemonConnectionType.NotConnected ||
        daemonConnection === DaemonConnectionType.Connecting
      );
    },
  },
  watch: {
    isMinimizedUI() {
      this.updateUIState();
    },
  },
  methods: {
    onAccountSettings: function () {
      //if (this.$store.state.settings.minimizedUI)
      sender.ShowAccountSettings();
      //else this.$router.push({ name: "settings", params: { view: "account" } });
    },
    onSettings: function () {
      sender.ShowSettings();
    },
    onConnectionSettings: function () {
      sender.ShowConnectionSettings();
    },
    onWifiSettings: function () {
      sender.ShowWifiSettings();
    },
    onDefaultLeftView: function (isDefaultView) {
      this.isCanShowMinimizedButtons = isDefaultView;
    },
    onMaximize: function (isMaximize) {
      this.$store.dispatch("settings/minimizedUI", !isMaximize);
      this.updateUIState();
    },
    updateUIState: function () {
      sender.uiMinimize(this.isMinimizedUI);
    },
  },
};
</script>

<style lang="scss" scoped>
@import "@/components/scss/constants";
</style>
