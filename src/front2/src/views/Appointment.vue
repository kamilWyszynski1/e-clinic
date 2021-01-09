<template>
  <div>
    <h1>Make an appointment</h1>
    <br />
    <form class="form" v-if="specialists" v-on:submit.prevent="submit">
      <div class="form-group row">
        <label for="comment" class="col-4 col-form-label"
          >Describe your case</label
        >
        <div class="col-4">
          <textarea
            id="comment"
            name="comment"
            cols="40"
            rows="5"
            required="required"
            class="form-control"
            v-model="comment"
          ></textarea>
        </div>
      </div>
      <div class="form-group row">
        <label for="symptoms" class="col-4 col-form-label"
          >Describe symptoms</label
        >
        <div class="col-4">
          <input
            id="symptoms"
            name="symptoms"
            placeholder="itching;pain;nausea"
            type="text"
            class="form-control"
            aria-describedby="symptomsHelpBlock"
            v-model="symptoms"
          />
          <span id="symptomsHelpBlock" class="form-text text-muted"
            >separate each symptom with ;</span
          >
        </div>
      </div>
      <div class="form-group row">
        <label for="select" class="col-4 col-form-label">Specialist</label>
        <div class="col-4">
          <select
            id="select"
            name="select"
            class="custom-select"
            required="required"
            v-model="feeID"
          >
            <option
              v-bind:value="spec.fee_id"
              v-bind:key="inx"
              v-for="(spec, inx) in specialists"
            >
              {{ spec.name }} {{ spec.surname }} -- {{ spec.speciality }}
            </option>
          </select>
        </div>
      </div>
      <div class="form-group row">
        <label for="example-datetime-local-input" class="col-4 col-form-label"
          >Date and time</label
        >
        <div class="col-4">
          <input
            class="form-control"
            type="datetime-local"
            value="2011-08-19T13:45:00"
            id="example-datetime-local-input"
            v-model="time"
          />
        </div>
      </div>
      <div class="form-group row">
        <label for="select" class="col-4 col-form-label">Duration</label>
        <div class="col-4">
          <select
            id="select"
            name="select"
            class="custom-select"
            required="required"
            v-model="duration"
          >
            <option
              v-bind:value="d.float"
              v-bind:key="inx"
              v-for="(d, inx) in durations"
            >
              {{ d.text }}
            </option>
          </select>
        </div>
      </div>
      <div class="form-group row">
        <div class="offset-4 col-4">
          <button name="submit" class="btn btn-primary">Submit</button>
        </div>
      </div>
    </form>

    <hr />
    <h1>Your appointments</h1>
    <table v-if="appointments" class="tg" id="customers">
      <thead>
        <tr>
          <th class="tg-0lax">Status</th>
          <th class="tg-0lax">Time</th>
          <th class="tg-0lax">Duration</th>
          <th class="tg-0lax">Moc</th>
          <th class="tg-0lax">Typ preparatu</th>
        </tr>
      </thead>
      <tbody>
        <!-- <tr v-bind:key="drug.id" v-for="drug in drugs">
              <td class="tg-0lax">{{ drug.name }}</td>
              <td class="tg-0lax">{{ drug.type_of_preparation }}</td>
              <td class="tg-0lax">{{ drug.common_name }}</td>
              <td class="tg-0lax">{{ drug.strength }}</td>
              <td class="tg-0lax">{{ drug.shape }}</td>
            </tr> -->
        <router-link
          v-bind:key="drug.id"
          v-for="drug in replacements"
          tag="tr"
          :to="{ name: 'drug', params: { id: drug.id } }"
        >
          <td class="tg-0lax">
            <b>{{ drug.name }}</b>
          </td>
          <td class="tg-0lax">{{ drug.type_of_preparation }}</td>
          <td class="tg-0lax">{{ drug.common_name }}</td>
          <td class="tg-0lax">{{ drug.strength }}</td>
          <td class="tg-0lax">{{ drug.shape }}</td>
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
      patientID: "147fc868-efa7-4f37-b85c-6a3f05e78b44",
      comment: "",
      symptoms: "",
      feeID: "",
      time: "",
      specialists: [],
      duration: "",
      appointments: [],
      durations: [
        {
          float: 0.5,
          text: "30min",
        },
        {
          float: 0.75,
          text: "45min",
        },
        {
          float: 1,
          text: "60min",
        },
        {
          float: 1.5,
          text: "90min",
        },
      ],
    };
  },
  methods: {
    clearForm() {
      this.feeID = "";
      this.time = "";
      this.duration = "";
      this.comment = "";
      this.symptoms = "";
    },
    getAppointments() {
      axios
        .post(`http://localhost:8081/api/v1/Assistant/GetAppointments`, {
          id: this.patientID,
          user_type: "patient",
          range: {
            from: "2020-10-13T06:01:07+02:00",
            to: "2021-01-04T13:01:18+01:00",
          },
        })
        .then((response) => {
          console.log(response.data);
          this.appointments = response.data.appointments;
        });
    },
    submit() {
      let specialist_id = "";
      let speciality = "";
      for (const i in this.specialists) {
        if (this.specialists[i].fee_id == this.feeID) {
          specialist_id = this.specialists[i].specialist_id;
          speciality = this.specialists[i].speciality;
          break;
        }
      }
      console.log(specialist_id);

      axios
        .post("http://localhost:8081/api/v1/Assistant/CreateAppointment", {
          patient_id: this.patientID,
          specialist_id: specialist_id,
          speciality: speciality,
          scheduled_at: new Date(this.time).toISOString(),
          duration: this.duration * 3600 * 1e9, // duration in nanoseconds
          patient_comment: this.comment,
          patient_symptoms: this.symptoms.split(";"),
        })
        .then((res) => {
          console.log(res);
        })
        .finally(this.clearForm());
    },
  },
  created() {
    axios
      .get(`http://localhost:8081/api/v1/Assistant/GetSpecialists`)
      .then((response) => {
        console.log(response.data);
        this.specialists = response.data.specialists;
      });
      this.getAppointments()
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