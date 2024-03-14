<template>
  <div>
    <div class="text-muted mb-3">WIREGUARD SETTINGS</div>
    <!-- Wireguard -->
    <div>
      <div class="flexRow paramBlock">
        <div class="defColor paramName">Preferred port:</div>
        <select v-model="port" class="form-select form-select-sm">
          <option
            v-for="item in prefferedPorts"
            :value="item.port"
            :key="item.text"
          >
            {{ item.text }}
          </option>
        </select>
      </div>

      <div class="flexRow paramBlock">
        <div class="defColor paramName">Rotate key every:</div>
        <select
          v-model="wgKeyRegenerationInterval"
          class="form-select form-select-sm"
        >
          <option
            v-for="item in wgRegenerationIntervals"
            :value="item.seconds"
            :key="item.seconds"
          >
            {{ item.text }}
          </option>
        </select>
      </div>

      <div v-if="IsAccountActive">
        <div class="settingsBoldFont text-muted">
          Wireguard key information:
        </div>

        <spinner :loading="isProcessing" />
        <div class="flexRow paramBlock">
          <div class="defColor paramName">Local IP Address:</div>
          <div class="detailedParamValue">
            {{ this.$store.state.account.session.WgLocalIP }}
          </div>
        </div>
        <div class="flexRow paramBlockDetailedConfig">
          <div class="defColor paramName">Public key:</div>
          <div class="detailedParamValue">
            {{ this.$store.state.account.session.WgPublicKey }}
          </div>
        </div>
        <div class="flexRow paramBlockDetailedConfig">
          <div class="defColor paramName">Generated:</div>
          <div class="detailedParamValue">
            {{ wgKeysGeneratedDateStr }}
          </div>
        </div>
        <div class="flexRow paramBlockDetailedConfig">
          <div class="defColor paramName">Expiration date:</div>
          <div class="detailedParamValue">
            {{ wgKeysExpirationDateStr }}
          </div>
        </div>
        <div class="flexRow paramBlockDetailedConfig">
          <div class="defColor paramName">Will be automatically rotated:</div>
          <div class="detailedParamValue">
            {{ wgKeysWillBeRegeneratedStr }}
          </div>
        </div>

        <button class="btn btn-primary py-1" v-on:click="onWgKeyRegenerate">
          Regenerate
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import spinner from "@/components/controls/control-spinner.vue";
import {
  VpnStateEnum,
  PortTypeEnum,
  V2RayObfuscationEnum,
} from "@/store/types";

import { enumValueName, dateDefaultFormat } from "@/helpers/helpers";
// import { SetInputFilterNumbers } from "@/helpers/renderer";

const sender = window.ipcSender;

export default {
  components: {
    spinner,
  },
  data: function () {
    return {
      isPortModified: false,
      isProcessing: false,
      openvpnManualConfig: false,
    };
  },
  mounted() {},
  watch: {
    // If port was changed in conneted state - reconnect
    async port(newValue, oldValue) {
      if (this.isPortModified === false) return;
      if (newValue == null || oldValue == null) return;
      if (newValue.port === oldValue.port && newValue.type === oldValue.type)
        return;
      await this.reconnect();
    },
  },

  methods: {
    async reconnect() {
      if (
        !this.$store.getters["vpnState/isConnected"] &&
        !this.$store.getters["vpnState/isConnecting"]
      )
        return; // not connected. Reconnection not required

      // Re-connect
      try {
        await sender.Connect();
      } catch (e) {
        console.error(e);
        sender.showMessageBoxSync({
          type: "error",
          buttons: ["OK"],
          message: `Failed to connect: ` + e,
        });
      }
    },

    onVPNConfigFileLocation: function () {
      const file = this.userDefinedOvpnFile;
      if (file) sender.shellShowItemInFolder(file);
    },
    onWgKeyRegenerate: async function () {
      try {
        this.isProcessing = true;
        await sender.WgRegenerateKeys();
      } catch (e) {
        console.log(`ERROR: ${e}`);
        sender.showMessageBoxSync({
          type: "error",
          buttons: ["OK"],
          message: "Error generating WireGuard keys",
          detail: e,
        });
      } finally {
        this.isProcessing = false;
      }
    },
    getWgKeysGenerated: function () {
      if (
        this.$store.state.account == null ||
        this.$store.state.account.session == null ||
        this.$store.state.account.session.WgKeyGenerated == null
      )
        return null;
      return new Date(this.$store.state.account.session.WgKeyGenerated);
    },
    formatDate: function (d) {
      if (d == null) return null;
      return dateDefaultFormat(d);
    },
    onShowHelpObfsproxy: function () {
      try {
        this.$refs.helpDialogObfsproxy.showModal();
      } catch (e) {
        console.error(e);
      }
    },
  },
  computed: {
    tooltipOptionEditRequiresDisconnection: function () {
      if (this.isDisconnected) return null;
      return "Disconnect VPN to edit this option";
    },
    isDisconnected: function () {
      return (
        this.$store.state.vpnState.connectionState === VpnStateEnum.DISCONNECTED
      );
    },
    isCanUseIPv6InTunnel: function () {
      return this.$store.getters["isCanUseIPv6InTunnel"];
    },
    IsAccountActive: function () {
      // if no info about account status - let's believe that account is active
      if (
        !this.$store.state.account ||
        !this.$store.state.account.accountStatus
      )
        return true;
      return this.$store.state.account?.accountStatus?.Active === true;
    },
    mtu: {
      get() {
        return this.$store.state.settings.mtu;
      },
      set(value) {
        this.$store.dispatch("settings/mtu", value);
      },
    },
    isMtuBadValue: function () {
      if (
        this.mtu != null &&
        this.mtu != "" &&
        this.mtu != 0 &&
        (this.mtu < 1280 || this.mtu > 65535)
      ) {
        return true;
      }
      return false;
    },
    userDefinedOvpnFile: function () {
      if (!this.$store.state.settings.daemonSettings) return null;
      return this.$store.state.settings.daemonSettings.UserDefinedOvpnFile;
    },
    wgKeyRegenerationInterval: {
      get() {
        return this.$store.state.account.session.WgKeysRegenIntervalSec;
      },
      set(value) {
        // daemon will send back a Hello response with updated 'session.WgKeysRegenIntervalSec'
        sender.WgSetKeysRotationInterval(value);
      },
    },
    wgKeysGeneratedDateStr: function () {
      return this.formatDate(this.getWgKeysGenerated());
    },
    wgKeysWillBeRegeneratedStr: function () {
      let t = this.getWgKeysGenerated();
      if (t == null) return null;

      t.setSeconds(
        t.getSeconds() +
          this.$store.state.account.session.WgKeysRegenIntervalSec,
      );

      let now = new Date();
      if (t < now) {
        // Do not show planned regeneration date in the past (it can happen after the computer wake up from a long sleep)
        // Show 'today' as planned date to regenerate keys in this case.
        // (the max interval to check if regeneration required is defined on daemon side, it is less than 24 hours)
        t = now;
      }

      return this.formatDate(t);
    },
    wgKeysExpirationDateStr: function () {
      let t = this.getWgKeysGenerated();
      if (t == null) return null;
      t.setSeconds(t.getSeconds() + 40 * 24 * 60 * 60); // 40 days
      return this.formatDate(t);
    },
    wgRegenerationIntervals: function () {
      let ret = [{ text: "1 day", seconds: 24 * 60 * 60 }];
      for (let i = 2; i <= 30; i++) {
        ret.push({ text: `${i} days`, seconds: i * 24 * 60 * 60 });
      }
      return ret;
    },
    wgQuantumResistanceStr: function () {
      if (this.$store.state.account.session.WgUsePresharedKey === true)
        return "Enabled";
      return "Disabled";
    },
    isShowAddPortOption: function () {
      const V2RayType = this.$store.getters["settings/getV2RayConfig"];
      const isV2Ray =
        V2RayType === V2RayObfuscationEnum.QUIC ||
        V2RayType === V2RayObfuscationEnum.TCP;

      if (
        (this.$store.state.settings.isMultiHop === true && isV2Ray === false) ||
        this.$store.getters["settings/isConnectionUseObfsproxy"]
      )
        return false;

      const ranges = this.$store.getters["vpnState/portRanges"];
      if (!ranges || ranges.length <= 0) return false;

      return true;
    },
    port: {
      get() {
        return this.$store.getters["settings/getPort"];
      },
      set(value) {
        this.isPortModified = true;

        if (value == "valueAddCustomPort") {
          // we need it just to update UI to show current port (except 'Add custom port...')
          this.$store.dispatch("settings/setPort", this.port);

          try {
            this.$refs.addCustomPortDlg.showModal();
          } catch (e) {
            console.error(e);
          }
          return;
        }

        this.$store.dispatch("settings/setPort", value);
      },
    },

    // Return suitable ports for current connection type.
    // If Obfuscation is enabled - the ports can differ from the default ports:
    // - Obfsproxy uses only TCP ports
    // - V2Ray (TCP) uses only TCP ports
    // - V2Ray (QUIC) uses only UDP ports
    prefferedPorts: function () {
      let ret = [];
      let ports = this.$store.getters["vpnState/connectionPorts"];

      // create UI items for ports
      ports.forEach((p) =>
        ret.push({
          text: `${enumValueName(PortTypeEnum, p.type)} ${p.port}`,
          key: `${enumValueName(PortTypeEnum, p.type)} ${p.port}`,
          port: p,
        }),
      );
      return ret;
    },
  },
};
</script>

<style lang="scss" scoped>
@import "@/components/scss/constants";
@import "@/components/scss/platform/base";

.defColor {
  @extend .settingsDefaultTextColor;
}

div.param {
  @extend .flexRow;
  margin-top: 3px;
}

input:disabled {
  opacity: 0.5;
}

input:disabled + label {
  opacity: 0.5;
}

div.paramBlock {
  @extend .flexRow;
  margin-top: 10px;
}

div.paramBlockDetailedConfig {
  @extend .flexRow;
  margin-top: 5px;
}

div.detailedConfigBlock {
  margin-left: 22px;
  max-width: 325px;
}

div.detailedConfigBlock input {
  width: 100%;
}

div.detailedConfigBlock select {
  width: 100%;
}

div.detailedConfigParamBlock {
  @extend .flexRow;
  margin-top: 10px;
  width: 100%;
}

div.detailedParamValue {
  opacity: 0.7;

  overflow-wrap: break-word;
  -webkit-user-select: text;
  letter-spacing: 0.1px;
}

div.paramName {
  min-width: 161px;
  max-width: 161px;
}

div.settingsRadioBtnProxy {
  @extend .settingsRadioBtn;
  padding-right: 20px;
}

div.paramBlockText {
  margin-top: 16px;
  margin-right: 21px;
}

select {
  background: linear-gradient(180deg, #ffffff 0%, #ffffff 100%);
  border: 0.5px solid rgba(0, 0, 0, 0.2);
  border-radius: 3.5px;
  width: 186px;
}

div.description {
  @extend .settingsGrayLongDescriptionFont;
  margin-left: 22px;
}

input.proxyParam {
  width: 100px;
}

div.disabled {
  pointer-events: none;
  opacity: 0.5;
}
</style>
