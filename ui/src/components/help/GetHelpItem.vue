<script>
export default {
  props: {
    image: String,
    title: String,
    subtitle: String,
    btnText: String,
    btnAction: Function,
  },
  computed: {},
};
</script>

<template>
  <div
    class="d-flex align-items-center justify-content-between px-3 py-2 rounded-3 my-2 custom-border"
  >
    <div class="d-flex align-items-center gap-3">
      <img :src="image" alt="no" class="item-icon" />
      <div class="icon-side-text">
        <h6 class="m-0 text-white">{{ title }}</h6>
        <p class="m-0 text-white-50">{{ subtitle }}s</p>
      </div>
    </div>
    <div class="">
      <button
        class="btn btn-primary text-white py-1 px-3 rounded-4"
        @click="btnAction"
      >
        {{ btnText }}
      </button>
    </div>
  </div>
</template>

<style scoped>
.custom-border {
  border: 1px solid rgba(255, 255, 255, 0.15);
}
.item-icon {
  height: 36px;
  width: 36px;
}
</style>
