<template>
  <div class="switch">
    <label  v-bind:class="{ load: isProgress }">
      <input
        type="checkbox"
        :checked="isConnected"
        v-on:click="DoSwitch($event)"
      />
      <span v-if="isProgress">
        <LottieAnimation
            :width="150"
            :height="150"
            path="./lottie/connecting.json"
        />
      </span>
      <span v-else-if="isConnected">
      <LottieAnimation
            :width="150"
            :height="150"
            :loop="false"
            path="./lottie/connected.json"
        />
      </span>
      <span v-else-if="isChecked">Is checked</span>
      <span v-else>
        <LottieAnimation
      :width="150"
      :height="150"
      :loop="false"
      path="./lottie/disconnectingnew.json"
  />
      </span>
        <!-- <LottieAnimation
        v-if="isConnected"
            :width="150"
            :height="150"
            :loop="false"
            path="./lottie/connected.json"
        />
             -->
          <!--  -->
        <!--  -->
      <div :style="style"></div>
    </label>



  </div>
</template>

<script>
import LottieAnimation from "lottie-vuejs/src/LottieAnimation.vue"
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

components:{
  LottieAnimation
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

.switch{
  width: fit-content;
  height: fit-content;
}
  .switch:hover {
    input {
      & + div {
        opacity: 0.7;
      }
    }
  }
</style>
