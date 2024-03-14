<template>
  <div>
    <div class="text-secondary mb-3">Non-VPN traffic blocking:</div>
    <div class="setting-list-item">
      <div class="d-flex justify-content-between w-50">
        <div class="d-flex align-items-center">
          <label class="ms-3" for="onDemand">
            <h6 class="mb-0">On-demand</h6>
            <p
              class="mb-0 text-secondary text-italic"
              style="font-size: smaller"
            >
              The VPN Firewall can be activated manually or automatically with a
              VPN connection.
            </p>
          </label>
        </div>
        <div class="settingsRadioBtn d-flex align-items-center">
          <input
            ref="radioFWOnDemand"
            class="form-check-input me-1"
            type="radio"
            id="onDemand"
            name="firewall"
            value="false"
            v-on:click="onPersistentFWChange(false)"
          />
        </div>
      </div>
      <div class="d-flex justify-content-between w-50">
        <div class="d-flex align-items-center">
          <label class="ms-3" for="alwaysOn">
            <h6 class="mb-0">Always-on firewall</h6>
            <p
              class="mb-0 text-secondary text-italic"
              style="font-size: smaller"
            >
              The VPN Firewall starts at boot and runs without the VPN client.
            </p>
          </label>
        </div>
        <div class="settingsRadioBtn d-flex align-items-center ms-3">
          <input
            ref="radioFWPersistent"
            type="radio"
            class="form-check-input me-1"
            id="alwaysOn"
            name="firewall"
            value="true"
            v-on:click="onPersistentFWChange(true)"
          />
          <label class="defColor" for="alwaysOn"></label>
        </div>
      </div>
    </div>
    <div class="setting-list-item">
      <div>
        <h6 class="mb-0">
          Allow access to {{ AppName }} servers when Firewall is enabled
        </h6>
      </div>
      <div class="param form-check form-switch">
        <input
          type="checkbox"
          class="form-check-input me-1"
          id="firewallAllowApiServers"
          v-model="firewallAllowApiServers"
        />
      </div>
    </div>

    <!-- On-demand Firewall -->
    <div class="text-secondary mb-3">On-demand Firewall:</div>
    <div class="setting-list-item">
      <label class="defColor" for="firewallActivateOnConnect"
        >Activate {{ AppName }} Firewall on connect to VPN</label
      >

      <div class="param form-check form-switch">
        <input
          type="checkbox"
          class="form-check-input me-1"
          id="firewallActivateOnConnect"
          :disabled="IsPersistent === true"
          v-model="firewallActivateOnConnect"
        />
      </div>
    </div>
    <div class="setting-list-item">
      <label class="defColor" for="firewallActivateOnConnect"
        >Deactivate {{ AppName }} Firewall on disconnect from VPN</label
      >

      <div class="param form-check form-switch">
        <input
          type="checkbox"
          class="form-check-input me-1"
          id="firewallDeactivateOnDisconnect"
          :disabled="IsPersistent === true"
          v-model="firewallDeactivateOnDisconnect"
        />
      </div>
    </div>
    <div class="param form-check form-switch">
      <label class="defColor" for="firewallDeactivateOnDisconnect"></label>
    </div>

    <!-- LAN settings -->
    <!-- <div class="settingsBoldFont">LAN settings:</div>
    <div class="param form-check form-switch">
      <input
        type="checkbox"
        class="form-check-input me-1"
        id="firewallAllowLan"
        v-model="firewallAllowLan"
      />
      <label class="defColor" for="firewallAllowLan"
        >Allow LAN traffic when {{ AppName }} Firewall is enabled</label
      >
    </div>
    <div class="param form-check form-switch">
      <input
        type="checkbox"
        id="firewallAllowMulticast"
        class="form-check-input me-1"
        :disabled="firewallAllowLan === false"
        v-model="firewallAllowMulticast"
      />
      <label class="defColor" for="firewallAllowMulticast"
        >Allow Multicast when LAN traffic is allowed</label
      >
    </div> -->
  </div>
</template>

<script>
import config from "@/config";

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
  data: function () {
    return {
      AppName: config.AppName,
    };
  },
  mounted() {
    this.updatePersistentFwUiState();
  },

  methods: {
    updatePersistentFwUiState() {
      if (this.$store.state.vpnState.firewallState.IsPersistent) {
        this.$refs.radioFWPersistent.checked = true;
        this.$refs.radioFWOnDemand.checked = false;
      } else {
        this.$refs.radioFWPersistent.checked = false;
        this.$refs.radioFWOnDemand.checked = true;
      }
    },
    async onPersistentFWChange(value) {
      try {
        await sender.KillSwitchSetIsPersistent(value);
      } catch (e) {
        processError(e);
      }
      this.updatePersistentFwUiState();
    },
  },
  watch: {
    IsPersistent() {
      this.updatePersistentFwUiState();
    },
  },
  computed: {
    IsPersistent: function () {
      return this.$store.state.vpnState.firewallState.IsPersistent;
    },
    firewallAllowApiServers: {
      get() {
        return this.$store.state.vpnState.firewallState.IsAllowApiServers;
      },
      async set(value) {
        await sender.KillSwitchSetAllowApiServers(value);
      },
    },
    firewallAllowLan: {
      get() {
        return this.$store.state.vpnState.firewallState.IsAllowLAN;
      },
      async set(value) {
        await sender.KillSwitchSetAllowLAN(value);
      },
    },
    firewallAllowMulticast: {
      get() {
        return this.$store.state.vpnState.firewallState.IsAllowMulticast;
      },
      async set(value) {
        await sender.KillSwitchSetAllowLANMulticast(value);
      },
    },

    firewallActivateOnConnect: {
      get() {
        return this.$store.state.settings.firewallActivateOnConnect;
      },
      set(value) {
        this.$store.dispatch("settings/firewallActivateOnConnect", value);
      },
    },
    firewallDeactivateOnDisconnect: {
      get() {
        return this.$store.state.settings.firewallDeactivateOnDisconnect;
      },
      set(value) {
        this.$store.dispatch("settings/firewallDeactivateOnDisconnect", value);
      },
    },
  },
};
</script>

<style scoped lang="scss">
@import "@/components/scss/constants";

.defColor {
  @extend .settingsDefaultTextColor;
}

div.param {
  @extend .flexRow;
  margin-top: 3px;
}

div.fwDescription {
  @extend .settingsGrayLongDescriptionFont;
  margin-bottom: 17px;
  margin-left: 22px;
  max-width: 425px;
}

input:disabled {
  opacity: 0.5;
}

input:disabled + label {
  opacity: 0.5;
}

label {
  margin-left: 1px;
}
</style>
