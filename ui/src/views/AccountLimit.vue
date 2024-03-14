<template>
  <div id="main" class="row">
    <div id="leftPanel">
      <div style="margin: 20px; width: 300px">
        <div class="large_text" style="margin-top: 20px">
          <button
            v-if="showContinueBtn"
            class="btn p-0 mr-5"
            v-on:click="onTryAgain"
          >
            Back
          </button>
          All Connected Devices
        </div>
        <div style="height: 16px"></div>
        <div class="device-list">
          <div
            class="device-list-item"
            v-for="device in activeDevices"
            v-bind:key="device.tokenId"
          >
            <div class="d-flex align-items-center">
              <img v-bind:src="deviceIcon(device.details.type)" alt="" />
              <span class="device-details">
                <div class="device-name">{{ device.details.name }}</div>
                <div class="device-status">
                  <span v-if="device.currentSession" class="text-red">
                    Current Device
                  </span>
                  <span v-else> Last seen {{ device._last_used_at }} </span>
                </div>
              </span>
            </div>
            <button
              v-if="!device.currentSession"
              class="btn p-0"
              v-on:click="logoutOne(device.tokenId)"
            >
              <img
                src="@/assets/minus.svg"
                class="minus-btn"
                style="cursor: pointer"
                alt=""
              />
            </button>
          </div>
        </div>
        <button
          v-if="showContinueBtn"
          class="btn btn-primary btn-lg w-100"
          v-on:click="continueToHome"
        >
          Continue
        </button>
        <button
          v-bind:class="[
            `btn ${
              showContinueBtn ? 'btn-secondary mt-2' : 'btn-primary'
            } btn-lg w-100`,
          ]"
          v-on:click="logoutAll"
        >
          Log out from all devices
        </button>

        <div v-if="!showContinueBtn" style="height: 16px"></div>
        <div class="centered">
          <button
            v-if="!showContinueBtn"
            class="btn btn-primary-outline"
            v-on:click="onTryAgain"
          >
            Go back
          </button>
        </div>
      </div>
    </div>

    <div id="rightPanel">
      <div class="text-center">
        <img src="@/assets/devices-limit.svg" />
        <div class="large_text">Devices limit reached</div>
        <div style="height: 22px"></div>
        <div class="small_text">
          According to your subscription plan you can use your VPN account only
          on {{ totalDevicesAllowed }} devices.
        </div>
        <!--<div class="elementFooter d-none">
          <div class="small_text">Do you think there is some issue?</div>
          <div style="height: 2px"></div>
          <button class="link linkFont" v-on:click="onContactSupport">
            Contact Support Team
          </button>
        </div>-->
      </div>
    </div>
  </div>
</template>

<script>
// import { isValidURL } from "@/helpers/helpers";
import Windows from "@/assets/windows-dark.svg";
import Linux from "@/assets/linux-dark.svg";
import iOS from "@/assets/ios-dark.svg";
import Android from "@/assets/android-dark.svg";
import MacOS from "@/assets/mac-dark.svg";

const sender = window.ipcSender;

export default {
  mounted() {
    // this.username = this.$route.params.username;
    // this.totalDevicesAllowed = this.$store.state.account.totalDevicesAllowed;
    // this.usedDevices = this.$store.state.account.usedDevices;
    // this.activeDevices = this.$store.state.account.activeDevices;
  },
  data: function () {
    return {
      // username: null,
      // totalDevicesAllowed: this.$store.state.account.totalDevicesAllowed,
      // activeDevices: this.$store.state.account.activeDevice,
    };
  },
  computed: {
    totalDevicesAllowed: function () {
      return this.$store.state.account.totalDevicesAllowed;
    },
    // usedDevices: function() {
    //   return this.$store.state.account.usedDevices;
    // },
    activeDevices: function () {
      return this.$store.state.account.activeDevice;
    },
    showContinueBtn: function () {
      return (
        this.activeDevices.filter((value) => value.currentSession).length > 0
      );
    },
  },
  methods: {
    onTryAgain: function () {
      this.$router.push("/");
    },
    logoutAll: async function () {
      await sender.LogoutAll();
    },
    logoutOne: async function (id) {
      await sender.LogoutDevice(id);
    },
    continueToHome: function () {
      var session = this.$store.state.account.session;
      sender.UpdateSession({ ...session, UserLoggedIn: true });
      this.$router.push({
        name: "Main",
      });
    },
    deviceIcon: function (type) {
      type = type.toLowerCase();
      if (type.startsWith("win")) {
        return Windows;
      } else if (type.startsWith("ios")) {
        return iOS;
      } else if (type.startsWith("mac")) {
        return MacOS;
      } else if (type.startsWith("linux")) {
        return Linux;
      } else if (type.startsWith("android")) {
        return Android;
      }
    },
  },
};
</script>

<style scoped lang="scss">
@import "@/components/scss/constants";

#main {
  height: 100%;
  display: flex;
  flex-direction: row;
}

.dark #main {
  background-color: #000;
}

#leftPanel {
  min-width: 364px;
  max-width: 364px;
  color: #000;
  flex-direction: column;
  display: flex;
  align-items: center;
}

#rightPanel {
  flex-direction: row;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  width: 460px;
  background: #9490b0;
}

.large_text {
  font-weight: 600;
  font-size: 18px;
  line-height: 120%;

  text-align: center;
}

.device-list {
  color: #000;
  height: 365px;
  overflow-y: auto;
  margin-bottom: 20px;
  padding-left: 10px;
}

.device-list .device-list-item {
  background: #f7f7f7;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-radius: 8px;
  padding: 15px 15px;
  margin-bottom: 10px;
}

.device-list .device-list-item:hover {
  background: #eeeeee;
}

.dark .device-list .device-list-item {
  background: #3f3f3f;
}

.dark .device-list .device-list-item:hover {
  background: #313131;
}

.device-list-item img.minus-btn:hover {
  opacity: 0.6;
}

.device-list-item img.minus-btn:focus {
  opacity: 1;
}

.device-list-item .device-details {
  margin-left: 10px;
}

.device-list-item .device-details .device-name {
  font-family: "Outfit";
  font-style: normal;
  font-weight: 600;
  font-size: 14px;
  margin-bottom: 5px;
  line-height: 15px;
}

.device-list-item .device-details .device-status {
  font-family: "Outfit";
  font-style: normal;
  font-weight: 400;
  font-size: 12px;
  line-height: 13px;
  opacity: 0.5;
}

.small_text {
  font-size: 15px;
  line-height: 18px;
  text-align: center;
  letter-spacing: -0.3px;
  opacity: 0.7;
}

.small_text2 {
  font-size: 14px;
  line-height: 17px;
  text-align: center;
  letter-spacing: -0.3px;
}

.verticalSpace {
  margin-top: auto;
  margin-right: 0;
}

.linkFont {
  font-size: 12px;
  line-height: 18px;
  text-align: center;
  letter-spacing: -0.4px;
}

.centered {
  flex-direction: column;
  display: flex;
  justify-content: center;
  align-items: center;
}

.elementFooter {
  @extend .centered;
  margin-top: 36px;
}
</style>
