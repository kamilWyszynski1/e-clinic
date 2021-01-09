<template>
  <div>
    <div v-if="path == '/'">
      <h2>Your recent appointments</h2>
      {{ drugs }}
    </div>
    <div v-if="path == '/profile'">profile</div>
    <div v-if="path == '/appointments'">
      <div>
        <h2>Make an appointment</h2>
        <form
          action="http://localhost:8081/api/v1/Assistant/CreateAppointment"
          method="POST"
          enctype="application/x-www-form-urlencoded"
        >
          <vue-form-generator
            :schema="schema"
            :model="model"
            :options="formOptions"
          ></vue-form-generator>
        </form>
        <div></div>
      </div>
      <div>
        <h2>Recent appointments</h2>
        <Appointments v-bind:appointments="appointments"/>
      </div>
    </div>
    <div v-if="path == '/drugs'">
      <Drugs v-bind:drugs="drugs.data.drugs" />
    </div>
  </div>
</template>

<script>
import axios from "axios";
import Drugs from "../drug/Drugs";
import VueFormGenerator from "vue-form-generator";
import Appointments from "../appointment/Appointment.vue"

export default {
  name: "Patient",
  components: {
    Drugs,
    Appointments,
    "vue-form-generator": VueFormGenerator.component,
  },
  props: ["path", "id"],
  data() {
    return {
      drugs: null,
      appointments: [],
      specialists: [
        {
          id: "45607504-969d-447a-b94f-e33b59beaca0",
          name: "Doc#1",
          speciality: "Diagnostic Radiology",
        },
        {
          id: "c65512c4-9809-4f8e-9d6d-3e4681fc79d6",
          name: "Doc#2",
          speciality: "Dermatology",
        },
      ],
      model: {
        name: "John Doe",
        patient_comment: "",
        patient_symptoms: "",
        specialist_id: "",
        duration: "",
        date: "",
        time: "",
      },
      formOptions: {
        validateAfterLoad: true,
        validateAfterChanged: true,
        validateAsync: true,
      },
    };
  },
  created() {
    // axios
    //   .get("http://localhost:8081/api/v1/Assistant/GetSpecialists")
    //   .then((res) => {
    //     console.log(res.data.specialists);
    //     this.specialists = res.data.specialists;
    //   });
  },
  mounted() {
    if (this.path == "/drugs") {
      axios
        .get(
          "http://localhost:8081/api/v1/Assistant/GetDrugs?prefix=&offset=0&limit=20"
        )
        .then((response) => (this.drugs = response));
    } else if (this.path == "/appointments") {
    //   let d1 = new Date();
    //   let d2 = new Date();
      axios
        .post("http://localhost:8081/api/v1/Assistant/GetAppointments", {
          id: this.id,
          user_type: "patient",
          range: {
            from: "2020-10-12T02:30:54+02:00",
            to: "2021-01-03T09:30:41+01:00",
          },
        })
        .then((response) => {
            console.log(response.data.appointments)
            this.appointments = response.data.appointments}
            );
    }
  },
  computed: {
      specList() {
          let specList = [];
              for (const s in this.specialists) {
                specList.push({
                  id: this.specialists[s].id,
                  name: this.specialists[s].name + " - " + this.specialists[s].speciality,
                });
              }
              return specList;
      },
    schema() {
      var result = {
        fields: [
          {
            type: "input",
            inputType: "textarea",
            label: "Describe your illness",
            model: "patient_comment",
            validator: VueFormGenerator.validators.string,
          },

          {
            type: "input",
            inputType: "email",
            label: "Describe your symptoms",
            model: "patient_symptoms",
            placeholder: "Separted with ;",
          },
          {
            type: "select",
            label: "Pick specialist",
            model: "specialist_id",
            values: this.specList,
            
          },
          {
            type: "select",
            label: "Pick duration",
            model: "duration",
            values: ["30m", "45m", "60m", "90m"],
          },
          {
            type: "input",
            inputType: "datetime",
            label: "Pick time",
            model: "date",
          },
          {
            type: "submit",
            label: "",
            buttonText: "Submit",
            validateBeforeSubmit: true,
            onSubmit: () => {
              const m = this.model;
              m.patient_id = this.id;
              m.scheduled_at = new Date(m.date).toISOString();
              m.patient_symptoms = m.patient_symptoms.split(";");
              
              for (const i in this.specialists) {
                  if (this.specialists[i].id == m.specialist_id) {
                      m.speciality = this.specialists[i].speciality
                  }
              }
              axios
                .post(
                  "http://localhost:8081/api/v1/Assistant/CreateAppointment",
                  m
                )
                .then(function (response) {
                  console.log(response);
                })
                .catch(function (error) {
                  alert(error);
                  console.log(error);
                });
            },
          },
        ],
      };
      return result;
    },
  },
};
</script>