<template>
  <div style="height:576px; min-height:576px; display: flex;">
    <div >
      <Sidebar v-if="showSidebar"/>
    </div>
    <div style="width:100%; margin: 20px 20px 20px 0;">
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
</template>

<script>
const sender = window.ipcSender;

import { DaemonConnectionType } from "@/store/types";
import { Platform, PlatformEnum, IsWindowHasFrame } from "@/platform/platform";
import Init from "@/components/Init.vue";
import Login from "@/components/Login.vue";
import Control from "@/components/Control.vue";
import Sidebar from "@/components/Sidebar.vue";
export default {
  components: {
    Init,
    Login,
    Control,
    Sidebar
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
      if (!this.isLoggedIn) return Login;
      return Control;
    },
    showSidebar: function () {
      const daemonConnection = this.$store.state.daemonConnectionState;
      if (
        daemonConnection == null ||
        daemonConnection === DaemonConnectionType.NotConnected ||
        daemonConnection === DaemonConnectionType.Connecting
      )
        return false;
      if (!this.isLoggedIn) return false;
      return true;
    },
    isMapBlured: function () {
      if (this.currentViewComponent !== Control) return "true";
      return "false";
    },
    isMinimizedButtonsVisible: function () {
      if (this.currentViewComponent !== Control) return false;
      if (this.isCanShowMinimizedButtons !== true) return false;
      return this.isMinimizedUI;
    },
    isMinimizedUI: function () {
      return this.$store.state.settings.minimizedUI;
    },
    minimizedButtonsTransition: function () {
      if (Platform() === PlatformEnum.Linux) return "smooth-display";
      return "fade";
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

#flexview {
  display: flex;
  flex-direction: row;
  height: 100%;
}

#left {
  width: 690px;
  min-width: 690px;
  max-width: 690px;
}

#right {
  width: 0%; // ???
  flex-grow: 1;
}

</style>
