<template>
  <div class="main-container">
    <Sidebar active="gethelp" />
    <div class="get-help p-3">
      <h4 class="text-white fw-medium fs-5">Get Help</h4>
      <div class="Customer-help">
        <!--<p class="text-white-50 m-0 fw-normal">Customer Support</p>-->
        <div class="scroll-div">
          <!--<GetHelpItem-->
          <!--  :image="require('@/assets/help/faqs.svg')"-->
          <!--  title="FAQ's"-->
          <!--  subtitle="Frequently asked questions"-->
          <!--  btn-text="Visit Website"-->
          <!--  btn-action="() => console.log('Hello World')"-->
          <!--/>-->
          <GetHelpItem
            :image="require('@/assets/help/report-a-bug.svg')"
            title="Report a bug"
            subtitle="Send us a report technical issues"
            btn-text="Send Us"
            :btn-action="
              () => {
                reportBugDlgShown = true;
              }
            "
          />
          <!--<GetHelpItem-->
          <!--  :image="require('@/assets/help/feedback.svg')"-->
          <!--  title="Leave a feedback"-->
          <!--  subtitle="Submit feature requests and suggestions"-->
          <!--  btn-text="Rate Us"-->
          <!--  btn-action="() => console.log('Hello World')"-->
          <!--/>-->
          <GetHelpItem
            :image="require('@/assets/help/suggest-location.svg')"
            title="Suggest a Location"
            subtitle="Suggest your desired Location"
            btn-text="Suggest"
            :btn-action="
              () => {
                suggestLocationDlgShown = true;
              }
            "
          />
          <p class="text-white-50 m-0 fw-normal">General Info</p>
          <GetHelpItem
            :image="require('@/assets/help/privacy-policy.svg')"
            title="Privacy Policy"
            subtitle="Get to Known our terms of service"
            btn-text="Read More"
            :btn-action="visitPrivacyPolicy"
          />
          <GetHelpItem
            :image="require('@/assets/help/tos.svg')"
            title="Terms of Service"
            subtitle="Get to know our Terms of service"
            btn-text="Read More"
            :btn-action="visitTermsOfService"
          />
        </div>
      </div>
    </div>
    <div v-if="reportBugDlgShown" class="component-overlay">
      <ReportBugDlg :on-dismiss="() => (reportBugDlgShown = false)" />
    </div>
    <div v-if="suggestLocationDlgShown" class="component-overlay">
      <SuggestLocationDlg
        :on-dismiss="() => (suggestLocationDlgShown = false)"
      />
    </div>
  </div>
</template>
<style scoped>
.scroll-div {
  overflow-y: scroll !important;
  height: 400px;
}
.component-overlay {
  position: absolute;
  left: 0px;
  top: 0px;
  height: 100%;
  width: 100%;
  z-index: 2;
}
</style>
<script>
import Sidebar from "@/components/Component-Sidebar.vue";
import GetHelpItem from "@/components/help/GetHelpItem.vue";
import { PrivacyPolicyUrl, TermsOfServiceUrl } from "@/config";
import ReportBugDlg from "@/components/help/ReportBugDlg.vue";
import SuggestLocationDlg from "@/components/help/SuggestLocationDlg.vue";

const sender = window.ipcSender;
export default {
  components: {
    SuggestLocationDlg,
    ReportBugDlg,
    GetHelpItem,
    Sidebar,
  },
  data: () => {
    return {
      reportBugDlgShown: false,
      suggestLocationDlgShown: false,
    };
  },
  methods: {
    visitUrl(url) {
      sender.shellOpenExternal(url);
    },
    visitPrivacyPolicy() {
      this.visitUrl(PrivacyPolicyUrl);
    },
    visitTermsOfService() {
      this.visitUrl(TermsOfServiceUrl);
    },
  },
};
</script>
