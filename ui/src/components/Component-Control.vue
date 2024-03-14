<template>
  <div>
    <div
      :class="[
        'body-content',
        'connection-layout',
        { 'w-60': this.serverSectionOpen },
      ]"
    >
      <div class="d-flex justify-content-end"></div>
      <div>
        <ConnectBlock
          :isChecked="isConnected"
          :isProgress="isInProgress"
          :onChecked="switchChecked"
          :onPauseResume="onPauseResume"
          :pauseState="this.$store.state.vpnState.pauseState"
        />

        <div>
          <SelectedServerBlock :onShowServersPressed="onShowServersPressed" />

          <ConnectionDetailsBlock
            :onShowPorts="onShowPorts"
            :onShowWifiConfig="onShowWifiConfig"
          />
        </div>
      </div>
    </div>
    <div v-if="this.serverSectionOpen" class="server-pane">
      <Servers
        :closeServerSection="closeServerSection"
        :onBack="backToMainView"
        :onFastestServer="onFastestServer"
        :onRandomServer="onRandomServer"
        :onServerChanged="onServerChanged"
      />
    </div>
    <div v-if="getPinPopupShown" class="pin-overlay">
      <GetPin
        :on-dismiss="() => (getPinPopupShown = false)"
        :connection-handler="connectionHandler"
      />
    </div>
  </div>
</template>

<script>
import Servers from "./Component-Servers.vue";
import ConnectBlock from "./blocks/block-connect.vue";
import ConnectionDetailsBlock from "./blocks/block-connection-details.vue";
import SelectedServerBlock from "@/components/blocks/block-selected-server.vue";
// import HopButtonsBlock from "./blocks/block-hop-buttons.vue";
const sender = window.ipcSender;
import { VpnStateEnum, VpnTypeEnum, PauseStateEnum } from "@/store/types";
import { isStrNullOrEmpty } from "@/helpers/helpers";
import GetPin from "@/components/GetPin.vue";

const viewTypeEnum = Object.freeze({
  default: "default",
  serversEntry: "serversEntry",
  serversExit: "serversExit",
});

async function connect(me, isConnect) {
  try {
    me.isConnectProgress = true;
    if (isConnect === true) await sender.Connect();
    else await sender.Disconnect();
  } catch (e) {
    console.error(e);
    sender.showMessageBoxSync({
      type: "error",
      buttons: ["OK"],
      message: `Failed to ${isConnect ? "connect" : "disconnect"}: ` + e,
    });
  } finally {
    me.isConnectProgress = false;
  }
}

function connected(me) {
  return me.$store.state.vpnState.connectionState === VpnStateEnum.CONNECTED;
}

export default {
  props: {
    onConnectionSettings: Function,
    onWifiSettings: Function,
    onDefaultView: Function,
  },

  components: {
    GetPin,
    // BlockSelectedServer,
    // HopButtonsBlock,
    Servers,
    ConnectBlock,
    SelectedServerBlock,
    ConnectionDetailsBlock,
  },
  mounted() {
    this.recalcScrollButtonVisiblity();

    // ResizeObserver sometimes is stopping to work for unknown reason. So, We do not use it for now
    // Instead, watchers are in use: isMinimizedUI, isMultiHop
    //const resizeObserver = new ResizeObserver(this.recalcScrollButtonVisiblity);
    //resizeObserver.observe(this.$refs.scrollArea);
  },
  data: function () {
    return {
      isShowScrollButton: false,
      isConnectProgress: false,
      uiView: viewTypeEnum.default,
      lastServersPingRequestTime: null,
      serverSectionOpen: false,
      getPinPopupShown: false,
    };
  },

  computed: {
    isConnected: function () {
      return connected(this);
    },
    isOpenVPN: function () {
      return this.$store.state.settings.vpnType === VpnTypeEnum.OpenVPN;
    },
    isMultihopAllowed: function () {
      return this.$store.getters["account/isMultihopAllowed"];
    },
    port: function () {
      // needed for watcher
      return this.$store.getters["settings/getPort"];
    },
    isInProgress: function () {
      if (this.isConnectProgress) return this.isConnectProgress;
      return (
        this.$store.state.vpnState.connectionState !== VpnStateEnum.CONNECTED &&
        this.$store.state.vpnState.connectionState !== VpnStateEnum.DISCONNECTED
      );
    },
    // needed for watcher
    conectionState: function () {
      return this.$store.state.vpnState.connectionState;
    },
    ActiveUntil: function () {
      return this.$store.state.account.accountStatus.ActiveUntil * 1000;
    },
    IsExpired: function () {
      return new Date().valueOf() > this.ActiveUntil;
    },
    IsActive: function () {
      return this.$store.state.account.accountStatus?.Active ?? false;
    },
  },

  watch: {
    conectionState(newValue, oldValue) {
      // show connection failure description:

      // only in case of changing to DISCONNECTED
      if (newValue !== VpnStateEnum.DISCONNECTED || newValue == oldValue)
        return;

      // if disconnection reason defined
      let failureInfo = this.$store.state.vpnState.disconnectedInfo;
      if (!failureInfo || isStrNullOrEmpty(failureInfo.ReasonDescription))
        return;

      sender.showMessageBoxSync({
        type: "error",
        buttons: ["OK"],
        message: `Failed to connect`,
        detail: failureInfo.ReasonDescription,
      });
    },
    port(newValue, oldValue) {
      if (!connected(this)) return;
      if (newValue == null || oldValue == null) return;
      if (newValue.port === oldValue.port && newValue.type === oldValue.type)
        return;
      connect(this, true);
    },
  },

  methods: {
    async switchChecked(isConnect) {
      console.log("Current Connection Status: ", isConnect);
      if (isConnect) {
        if (this.IsActive && !this.IsExpired) {
          connect(this, isConnect);
        } else {
          this.getPinPopupShown = true;
        }
      } else {
        connect(this, isConnect);
      }
    },
    connectionHandler() {
      connect(this, true);
    },
    async onPauseResume(seconds) {
      if (seconds == null || seconds == 0) {
        // RESUME
        if (this.$store.state.vpnState.pauseState != PauseStateEnum.Resumed)
          await sender.ResumeConnection();
      } else {
        // PAUSE
        await sender.PauseConnection(seconds);
      }
    },
    onShowServersPressed(isExitServers) {
      console.log(isExitServers);
      /*this.uiView = isExitServers
        ? viewTypeEnum.serversExit
        : viewTypeEnum.serversEntry;

      if (this.onDefaultView) this.onDefaultView(false);*/

      // request servers ping not more often than once per 30 seconds
      if (
        this.lastServersPingRequestTime == null ||
        (new Date().getTime() - this.lastServersPingRequestTime.getTime()) /
          1000 >
          30
      ) {
        sender.PingServers();
        this.lastServersPingRequestTime = new Date();
      }
      this.serverSectionOpen = true;
      // this.$router.push("/servers");
    },
    onShowPorts() {
      if (this.onConnectionSettings != null) this.onConnectionSettings();
    },
    onShowWifiConfig() {
      if (this.onWifiSettings != null) this.onWifiSettings();
    },
    backToMainView() {
      this.uiView = viewTypeEnum.default;
      if (this.onDefaultView) this.onDefaultView(true);

      setTimeout(this.recalcScrollButtonVisiblity, 1000);
    },
    onServerChanged(server, isExitServer) {
      if (server == null || isExitServer == null) return;

      let needReconnect = false;
      if (!isExitServer) {
        if (
          !this.$store.state.settings.serverEntry ||
          this.$store.state.settings.serverEntry.flag !== server.flag ||
          this.$store.state.settings.isRandomServer !== false
        ) {
          this.$store.dispatch("settings/isRandomServer", false);
          this.$store.dispatch("settings/serverEntry", server);
          needReconnect = true;
        }
      } else {
        if (
          !this.$store.state.settings.serverExit ||
          this.$store.state.settings.serverExit.flag !== server.flag ||
          this.$store.state.settings.isRandomExitServer !== false
        ) {
          this.$store.dispatch("settings/isRandomExitServer", false);
          this.$store.dispatch("settings/serverExit", server);
          needReconnect = true;
        }
      }
      if (this.$store.state.settings.isFastestServer !== false) {
        this.$store.dispatch("settings/isFastestServer", false);
        needReconnect = true;
      }
      if (needReconnect === true && connected(this)) connect(this, true);
    },
    onFastestServer() {
      this.$store.dispatch("settings/isFastestServer", true);
      if (connected(this)) connect(this, true);
    },
    closeServerSection() {
      this.serverSectionOpen = false;
    },
    onRandomServer(isExitServer) {
      if (isExitServer === true)
        this.$store.dispatch("settings/isRandomExitServer", true);
      else this.$store.dispatch("settings/isRandomServer", true);

      if (connected(this)) connect(this, true);
    },
    recalcScrollButtonVisiblity() {
      let sa = this.$refs.scrollArea;
      if (sa == null) {
        this.isShowScrollButton = false;
        return;
      }

      const isNeedToShow = function () {
        let pixelsToTheEndScroll =
          sa.scrollHeight - (sa.clientHeight + sa.scrollTop);
        // hide if the 'pixels to scroll' < 20
        if (pixelsToTheEndScroll < 20) return false;
        return true;
      };

      // hide - imediately; show - with 1sec delay
      if (!isNeedToShow()) this.isShowScrollButton = false;
      else {
        setTimeout(() => {
          this.isShowScrollButton = isNeedToShow();
        }, 1000);
      }
    },
    onScrollDown() {
      let sa = this.$refs.scrollArea;
      if (sa == null) return;
      sa.scrollTo({
        top: sa.scrollHeight,
        behavior: "smooth",
      });
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="scss" scoped>
@import "@/components/scss/constants";

.mainControl {
  height: 100vw;
  width: 100vh;
}
.w-60 {
  width: 55%;
}
.server-pane {
  width: 45%;
  position: absolute;
  top: 0;
  right: 0;
  border-left: 1px solid rgba(255, 255, 255, 0.15);
  height: 100%;
}
.body-content .connect-action button {
  background: none;
  border: none;
}

.pin-overlay {
  position: absolute;
  left: 0px;
  top: 0px;
  height: 100%;
  width: 100%;
  z-index: 2;
}
</style>
