<template>
  <div class="drugs">
    <div v-bind:key="drug.id" v-for="drug in drugs">
      <div class="drug">
        <h3>
          <a v-bind:href="'/drugs/' + drug.id">{{ drug.name }}</a>
        </h3>
        <table class="tg">
          <thead>
            <tr>
              <th class="tg-0lax">Rodzaj preparatu</th>
              <th class="tg-0lax">Nazwa powszechna</th>
              <th class="tg-0lax">Moc</th>
              <th class="tg-0lax">Typ preparatu</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td class="tg-0lax">{{ drug.type_of_preparation }}</td>
              <td class="tg-0lax">{{ drug.common_name }}</td>
              <td class="tg-0lax">{{ drug.strength }}</td>
              <td class="tg-0lax">{{ drug.shape }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "Drugs",
  props: ["drugs"],
};
</script>

<style scoped>
.tg {
  border-collapse: collapse;
  border-spacing: 0;
}
.tg td {
  border-color: black;
  border-style: solid;
  border-width: 1px;
  font-family: Arial, sans-serif;
  font-size: 14px;
  overflow: hidden;
  padding: 10px 5px;
  word-break: normal;
}
.tg th {
  border-color: black;
  border-style: solid;
  border-width: 1px;
  font-family: Arial, sans-serif;
  font-size: 14px;
  font-weight: normal;
  overflow: hidden;
  padding: 10px 5px;
  word-break: normal;
}
.tg .tg-0lax {
  text-align: left;
  vertical-align: top;
}
</style>