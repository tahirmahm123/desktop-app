<template>
  <div id="flexview">
    <div class="flexColumn">
      <div class="leftPanelTopSpace">
        <transition name="smooth-display">
          <div
            v-if="isMinimizedButtonsVisible"
            class="minimizedButtonsPanel leftPanelTopMinimizedButtonsPanel"
            v-bind:class="{
              minimizedButtonsPanelRightElements: isWindowHasFrame,
            }"
          >
            <button v-on:click="onAccountSettings()" title="Account settings">
              <img src="@/assets/user.svg" />
            </button>

            <button v-on:click="onSettings()" title="Settings">
              <img src="@/assets/settings.svg" />
            </button>
          </div>
        </transition>
      </div>
      <div class="flexColumn" style="min-height: 0px">
        <transition name="fade" mode="out-in">
          <component
            v-bind:is="currentViewComponent"
            :onConnectionSettings="onConnectionSettings"
            :onWifiSettings="onWifiSettings"
            :onFirewallSettings="onFirewallSettings"
            :onAntiTrackerSettings="onAntitrackerSettings"
            :onDefaultView="onDefaultLeftView"
            id="left"
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
import Login from "@/components/Component-Login.vue";
import Control from "@/components/Component-Control.vue";
import TheMap from "@/components/Component-Map.vue";
import ParanoidModePassword from "@/components/ParanoidModePassword.vue";
import Sidebar from "@/components/Component-Sidebar.vue";
export default {
  components: {
    Init,
    Login,
    Control,
    TheMap,
    Sidebar,
    ParanoidModePassword,
  },
  data: function () {
    return {
      isCanShowMinimizedButtons: true,
    };
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
    isMapBlured: function () {
      if (this.currentViewComponent !== Control) return "true";
      return "false";
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
    onFirewallSettings: function () {
      sender.ShowFirewallSettings();
    },
    onAntitrackerSettings: function () {
      sender.ShowAntitrackerSettings();
    },
  },
};
</script>

<style scoped lang="scss">
@import "@/components/scss/constants";
</style>
