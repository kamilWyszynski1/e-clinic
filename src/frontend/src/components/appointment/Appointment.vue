<template>
  <div v-if="appointmentInfos.length">
    <div
      v-bind:key="a.id"
      v-for="(a, inx) in appointmentInfos"
      :class="[`appointment${inx % 2}`]"
    >
      <p>{{ a.state }}</p>
      <p>{{ a.comment }}</p>
      <p>{{ a.symptoms }}</p>
      <p>{{ a.scheduled_time }}</p>
      <p>{{ a.duration / 60 }}min</p>
    </div>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "Appointment",
  props: ["appointments"],
  data() {
    return {
      appointmentInfos: [],
    };
  },
  watch: {
    appointments() {
      for (const i in this.appointments) {
        console.log(this.appointments[i].id);
        axios
          .get(
            "http://localhost:8081/api/v1/Assistant/GetAppointment?aID=" +
              this.appointments[i].id
          )
          .then((res) => {
            console.log(res.data);
            const a = res.data.appointment;
            this.appointmentInfos.push({
              state: a.state,
              comment: res.data.form.comment,
              symptoms: res.data.form.symptoms,
              scheduled_time: a.scheduled_time,
              duration: a.duration,
            });
          });
      }
    },
  },
};
</script>

<style scoped>
.appointment0 {
  background-color: #b7b7a4;
}
.appointment1 {
  background-color: #ddbea9;
}
</style>