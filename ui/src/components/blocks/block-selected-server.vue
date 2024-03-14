<template>
  <div>
    <div
      class="server-card d-flex align-items-center mx-5 justify-content-between"
    >
      <div class="flexRow">
        <serverNameControl
          class="serverName"
          style="max-width: 245px"
          :isLargeText="true"
          :server="this.server"
          :isCountryFirst="!(isFastestServer || isRandomServer)"
          :isFastestServer="isFastestServer"
          :isRandomServer="isRandomServer"
          :isShowPingPicture="!(isFastestServer || isRandomServer)"
        />
      </div>
      <serverPingInfoControl
        v-show="!(isFastestServer || isRandomServer)"
        :server="this.server"
        style="
          margin-left: 9px;
          margin-right: 8px;
          position: absolute;
          bottom: 12px;
          left: 30px;
        "
      />
      <div>
        <a href="#" class="mt-3" v-on:click.prevent="showServersList()">
          <button class="change-location">Change</button>
        </a>
      </div>
    </div>
  </div>
</template>

<script>
import serverNameControl from "@/components/controls/control-server-name.vue";
import serverPingInfoControl from "@/components/controls/control-server-ping.vue";
import { VpnStateEnum } from "@/store/types";

export default {
  props: ["onShowServersPressed", "isExitServer"],
  components: {
    serverNameControl,
    serverPingInfoControl,
  },
  computed: {
    server: function () {
      return this.$store.state.settings.serverEntry;
    },
    isConnected: function () {
      return (
        this.$store.state.vpnState.connectionState === VpnStateEnum.CONNECTED
      );
    },
    isConnecting: function () {
      return this.$store.getters["vpnState/isConnecting"];
    },
    isDisconnected: function () {
      return (
        this.$store.state.vpnState.connectionState === VpnStateEnum.DISCONNECTED
      );
    },
    isFastestServer: function () {
      if (
        (this.isDisconnected || this.$store.state.vpnState.isPingingServers) &&
        this.$store.getters["settings/isFastestServer"]
      )
        return true;
      return false;
    },
    isRandomServer: function () {
      if (!this.isDisconnected) return false;
      return this.isExitServer
        ? this.$store.getters["settings/isRandomExitServer"]
        : this.$store.getters["settings/isRandomServer"];
    },
  },
  methods: {
    showServersList() {
      if (this.onShowServersPressed != null)
        this.onShowServersPressed(this.isExitServer);
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
@import "@/components/scss/constants";

#main {
  @extend .left_panel_block;
  border: 1px solid rgba(139, 154, 171, 0.17);
  border-radius: 12px;
  margin: 10px;
}

.small_text {
  font-size: 14px;
  line-height: 17px;
  letter-spacing: -0.3px;
  color: var(--text-color-details);
}

img.pic {
  width: 35px;
  box-shadow: none !important;
}

.serverSelectBtn {
  padding: 16px;
  border: none;
  background-color: inherit;
  outline-width: 0;
  cursor: pointer;
  width: 100%;
}

.server-card {
  border: 1px solid rgba(139, 154, 171, 0.17);
  min-height: 72px;
  max-height: 72px;
  border-radius: 12px;
  padding: 0 15px;
  margin: 10px;
  position: relative;
}

.server-card .country-name {
  font-size: 16px;
  font-weight: 700;
}

.serverName {
  max-width: 270px;
}
.change-location {
  padding: 3px 10px;
  background-color: #ebc553;
  border: 1px solid #ebc553;
  font-weight: 500;
  color: #fff;
  border-radius: 50px;
}
.change-location:hover {
  background-color: #fff;
  color: #ebc553;
}
</style>
