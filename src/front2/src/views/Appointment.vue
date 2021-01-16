<template>
  <div>
    <h2>Appointment</h2>
    <table v-if="appointment" class="tg" id="customers">
      <thead>
        <tr>
          <th class="tg-0lax">Status</th>
          <th class="tg-0lax">Time</th>
          <th class="tg-0lax">Duration(min)</th>
        </tr>
      </thead>
      <tbody>
        <td class="tg-0lax">
          <b>{{ appointment.appointment.state }}</b>
        </td>
        <td class="tg-0lax">{{ appointment.appointment.scheduled_time }}</td>
        <td class="tg-0lax">{{ appointment.appointment.duration / 30 }}</td>
      </tbody>
    </table>
    <br />
    <h2>Form</h2>
    <table v-if="appointment" class="tg" id="customers">
      <thead>
        <tr>
          <th class="tg-0lax">Form comment</th>
        </tr>
      </thead>
      <tbody>
        <td class="tg-0lax">
          {{ appointment.form.comment }}
        </td>
      </tbody>
    </table>
    <br />
    <h2>Prescription</h2>
    <table v-if="appointment" class="tg" id="customers">
      <thead>
        <tr>
          <th class="tg-0lax">Comment</th>
        </tr>
      </thead>
      <tbody>
        <td class="tg-0lax">{{ appointment.prescription.comment }}</td>
      </tbody>
    </table>

    <table v-if="appointment" class="tg" id="customers">
      <thead>
        <tr>
          <th class="tg-0lax">Drug</th>
          <th class="tg-0lax">Dosing</th>
        </tr>
      </thead>
      <tbody>
        <router-link
          v-bind:key="drug.id"
          v-for="drug in appointment.prescription.drugs"
          tag="tr"
          :to="{ name: 'drug', params: { id: drug.drug } }"
        >
          <td class="tg-0lax">
            <b>{{ drug.drug }}</b>
          </td>
          <td class="tg-0lax">{{ drug.dosing }}</td>
        </router-link>
      </tbody>
    </table>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "Appointment",
  data() {
    return {
      appointment: null,
    };
  },
  methods: {
    getAppointment(uid) {
      axios
        .get(`http://localhost:8081/api/v1/Assistant/GetAppointment?aID=${uid}`)
        .then((response) => {
          console.log(response.data);
          this.appointment = response.data;
        });
    },
  },
  created() {
    this.getAppointment(this.$route.params.id);
  },
  beforeRouteUpdate(to, from, next) {
    console.log(to, from, next);
    this.getContent(to.params.id);
    next();
  },
};
// #007bff
</script>

<style scoped>
label {
  font-family: Avenir, Helvetica, Arial, sans-serif !important;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}
</style>