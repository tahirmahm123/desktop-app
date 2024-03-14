<template>
  <div class="switch">
    <label v-bind:class="{ load: isProgress }" class="btn-connection">
      <input
        style="display: none"
        type="checkbox"
        :checked="isConnected"
        v-on:click="DoSwitch($event)"
      />
      <span v-if="isProgress"> Connection in progress </span>
      <span v-else-if="isConnected"> Click to Disconnect </span>
      <span v-else-if="isChecked">Is checked</span>
      <span v-else> Click to connect </span>

      <div :style="style"></div>
    </label>
  </div>
</template>

<script>
// import LottieAnimation from "lottie-vuejs/src/LottieAnimation.vue";

export default {
  props: ["onChecked", "isChecked", "isProgress", "checkedColor"],
  computed: {
    isConnected: function () {
      if (this.isProgress) return false;
      return this.isChecked === true;
    },
    style: function () {
      if (this.checkedColor == null || this.isConnected === false) return "";
      return `background: ${this.checkedColor}`;
    },
  },

  components: {
    // LottieAnimation,
  },
  methods: {
    DoSwitch(e) {
      if (!this.isConnected) e.preventDefault();
      if (this.onChecked) {
        if (this.isProgress) this.onChecked(false, e);
        else this.onChecked(!this.isConnected, e);
      }
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
.switch {
  width: fit-content;
  height: fit-content;
  margin: 0 auto;
}
.switch label {
  cursor: pointer;
  font-size: 16.5px;
  font-style: normal;
  font-weight: 500;
}

.switch input:checked {
  background-color: #ebc553;
}

.switch:hover {
  input {
    & + div {
      opacity: 0.7;
    }
  }
}
label.btn-connection {
  background: #fff;
  border-radius: 50px;
  position: relative;
  border: 3px solid #fff;
  color: black !important;
  cursor: pointer;
  z-index: 1;
  padding: 13.626px 27.253px;
}
.btn-connection:hover {
  border: 3px solid #ebc553;
  background: #fff;
}
</style>
