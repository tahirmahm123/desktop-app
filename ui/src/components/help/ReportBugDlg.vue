<script>
export default {
  props: {
    onDismiss: Function,
  },
  computed: {},
};
</script>

<template>
  <div class="overlay d-flex justify-content-center align-items-center">
    <div class="confirm-bug px-4 py-3 rounded-3">
      <div class="d-flex justify-content-between align-items-center">
        <h5 class="m-0 text-white fw-bold fs-6">Report a Bug</h5>
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
      <form action="#">
        <label for="email-bug" class="text-white mt-2">Enter Email</label>
        <input
          type="email"
          id="email-bug"
          class="form-control bg-transparent text-white"
          placeholder="admin@gmail.com"
        />
        <label for="text-area" class="text-white mt-2"
          >Enter your message here</label
        >
        <textarea
          id="text-area"
          rows="5"
          class="form-control rounded bg-transparent"
          placeholder="Problem description"
        ></textarea>
        <button class="btn btn-primary w-100 text-white mt-2">
          <img src="" alt="" />
          Send
        </button>
        <div class="check-box d-flex align-items-start mt-2 gap-2">
          <input type="checkbox" class="" />
          <div>
            <h6 class="text-white m-0">Send detailed system info</h6>
            <p class="check-box-p text-white m-0">
              Attach detailed system info on running processes, installed
              application, memory, CPU , Battery state , and app log.
            </p>
          </div>
        </div>
      </form>
    </div>
  </div>
</template>

<style scoped>
.custom-border {
  border: 1px solid rgba(255, 255, 255, 0.15);
}
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
  width: 100%;
  height: 100%;
}
.item-icon {
  height: 36px;
  width: 36px;
}
.confirm-bug {
  width: 400px;
  background-color: #000;
  border: 1px solid rgba(255, 255, 255, 0.2);
}
</style>
