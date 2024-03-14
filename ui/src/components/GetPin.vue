<template>
  <div class="w-100 h-100">
    <div
      class="overlay d-flex align-items-center justify-content-center w-100 h-100"
    >
      <div class="d-flex flex-column align-items-center h-100 w-100">
        <div class="logo w-100 d-flex mt-4 justify-content-between">
          <img src="@/assets/img/main-logo.svg" alt="no" />
          <button @click="onDismiss" class="btn p-0">
            <svg
              width="38"
              height="38"
              viewBox="0 0 38 38"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                d="M11.7716 29.6955H26.2164C28.8048 29.6955 30.1407 28.3595 30.1407 25.8069V11.2905C30.1407 8.73793 28.8048 7.40199 26.2164 7.40199H11.7716C9.19516 7.40199 7.84729 8.726 7.84729 11.2905V25.8069C7.84729 28.3595 9.19516 29.6955 11.7716 29.6955ZM14.9683 23.7314C14.3361 23.7314 13.8352 23.2305 13.8352 22.5864C13.8352 22.3001 13.9544 22.0138 14.1811 21.7991L17.4255 18.5428L14.1811 15.2983C13.9544 15.0836 13.8352 14.7974 13.8352 14.5111C13.8352 13.867 14.3361 13.3779 14.9683 13.3779C15.3023 13.3779 15.5647 13.4853 15.7794 13.7L19.0358 16.9444L22.304 13.6881C22.5426 13.4614 22.7931 13.3541 23.1151 13.3541C23.7473 13.3541 24.2483 13.855 24.2483 14.4872C24.2483 14.7854 24.129 15.0479 23.9024 15.2864L20.658 18.5428L23.9024 21.7872C24.1171 22.0138 24.2364 22.2882 24.2364 22.5864C24.2364 23.2305 23.7354 23.7314 23.0913 23.7314C22.7692 23.7314 22.4949 23.6122 22.2683 23.3975L19.0358 20.165L15.8033 23.3975C15.5886 23.6241 15.3023 23.7314 14.9683 23.7314Z"
                fill="white"
              />
            </svg>
          </button>
        </div>
        <div
          class="main-content d-flex align-items-center justify-content-center w-100 h-100"
        >
          <div class="d-flex justify-content-center align-items-center">
            <form
              @submit.prevent="submitPinForm"
              class="pin-form bg-dark d-flex flex-column justify-content-center text-center"
            >
              <h5 class="fw-bold fs-6">Enter Pin to Connect</h5>
              <p class="fw-normal w-75 mx-auto">
                Please enter the PIN you received from our
                <a
                  :href="pinUrl"
                  @click.prevent="visitWebsite(pinUrl)"
                  class="text-warning"
                  >website</a
                >.
              </p>
              <div class="mb-3">
                <input
                  type="password"
                  class="form-control form-control-lg"
                  :class="{ 'is-invalid': pinError }"
                  v-model="pin"
                  ref="password"
                  aria-describedby="helpId"
                  placeholder="Enter Pin"
                  :disabled="isProcessing"
                  v-on:keyup="keyup($event)"
                />
                <div class="invalid-feedback">{{ pinErrorText }}</div>
              </div>
              <button
                class="btn btn-primary my-2"
                type="submit"
                :disabled="isProcessing"
              >
                Connect to VPN
              </button>
              <button
                type="button"
                class="btn btn-outline-primary text-white"
                @click="visitWebsite(pinUrl)"
              >
                Get a new PIN
              </button>
            </form>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { NewPinUrl } from "@/config";
const sender = window.ipcSender;
export default {
  methods: {
    visitWebsite(url) {
      sender.shellOpenExternal(url);
    },
    async submitPinForm() {
      try {
        this.pinError = false;
        console.log(this.pin);
        if (this.pin.length <= 0) {
          this.pinError = true;
          this.pinErrorText = "Please enter a Pin.";
          return;
        }
        this.isProcessing = true;
        const resp = await sender.VerifyPin(this.pin);
        const currentTimeStamp = new Date().valueOf() / 1000;
        const timestamp = resp.Account.ActiveUntil * 1000;
        if (resp.Account.Active && currentTimeStamp < timestamp) {
          console.log("connect here");
          this.connectionHandler();
        } else {
          this.pinError = true;
          this.pinErrorText =
            currentTimeStamp > timestamp
              ? "The Pin is Expired."
              : "The Pin has already been used.";
        }
        console.log(resp);
      } catch (e) {
        console.log(e);
      } finally {
        this.isProcessing = false;
      }
    },
    keyup(event) {
      if (event.keyCode === 13) {
        // Cancel the default action, if needed
        event.preventDefault();
        this.submitPinForm();
      }
    },
  },
  data: function () {
    return {
      pin: "",
      pinError: false,
      pinErrorText: "",
      isProcessing: false,
    };
  },
  props: {
    onDismiss: Function,
    connectionHandler: Function,
  },
  components: {},
  computed: {
    pinUrl: function () {
      return NewPinUrl;
    },
  },
};
</script>
<style scoped>
.overlay {
  background-color: rgba(
    0,
    0,
    0,
    0.5
  ); /* Adjust the alpha value for transparency */
  backdrop-filter: blur(5px); /* Adjust the blur value as needed */
  border-radius: 10px; /* Adjust the border radius as needed */
  padding: 20px; /* Adjust the padding as needed */
  color: white;
}
.pin-form {
  border: 1px solid #dcdcdc;
  padding: 40px 30px;
  border-radius: 13px;
  background-color: #000000;
}
input[type="text"],
input[type="password"] {
  border-radius: 7px;
}
.btn {
  border-radius: 8px;
}
</style>
