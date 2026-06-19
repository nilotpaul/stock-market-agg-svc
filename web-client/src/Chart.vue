<template>
  <div class="container">
    <h2 class="title">Candlestick Viewer</h2>
    <p>NOTE: fields are pre-filled for ease of use.</p>

    <form class="form" @submit="fetchCandles">
      <input v-model="symbol" placeholder="Symbol (TCS)" />
      <select v-model="timeframe">
        <option value="1m">1m</option>
        <option value="5m">5m</option>
        <option value="15m">15m</option>
        <option value="1h">1h</option>
        <option value="1d">1d</option>
      </select>

      <input v-model="startDate" type="datetime-local" />
      <input v-model="endDate" type="datetime-local" />
      <input v-model="limit" type="number" />
      <button type="submit">Fetch</button>
    </form>

    <div ref="chartContainer" class="chart"></div>
    <div class="error">
      <p>{{ errMsg }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";

const symbol = ref("TCS");
const timeframe = ref("1m");
const startDate = ref("2026-01-01 09:16");
const endDate = ref("2026-01-01 15:29");
const limit = ref(1000);

const errMsg = ref("");

const chartContainer = ref(null);

let chart = null;
onMounted(() => {
  const Highcharts = window.Highcharts;

  chart = Highcharts.stockChart(chartContainer.value, {
    title: { text: "Candlestick Chart" },
    rangeSelector: { selected: 1 },
    series: [
      {
        type: "candlestick",
        name: symbol.value,
        data: [],
      },
    ],
  });
});

async function fetchCandles(e) {
  e.preventDefault();

  const params = new URLSearchParams({
    symbol: symbol.value,
    timeframe: timeframe.value,
    start_date: format(startDate.value),
    end_date: format(endDate.value),
    limit: limit.value,
  });

  const res = await fetch(`/api/v1/candles?${params}`);
  const json = await res.json();

  if (!res.ok) {
    errMsg.value = json?.message || "something went wrong";
    return;
  }

  errMsg.value = "";
  const series = json.candles.map((c) => [
    new Date(c.datetime).getTime(),
    Number(c.open),
    Number(c.high),
    Number(c.low),
    Number(c.close),
  ]);

  if (chart) {
    chart.series[0].update({
      name: symbol.value,
    });
    chart.series[0].setData(series, true, false, false);
    chart.redraw();
  }
}

const format = (v) => {
  if (!v) return "";
  return v.replace("T", " ") + ":00";
};
</script>

<style scoped>
.container {
  padding: 24px;
  max-width: 1100px;
  margin: auto;
  font-family: Arial, sans-serif;
}

.title {
  font-size: 20px;
  margin-bottom: 16px;
}

.form {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 10px;
  margin-bottom: 20px;
}

input,
select {
  padding: 8px;
  border: 1px solid #ccc;
  border-radius: 6px;
}

button {
  background: #2563eb;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
}

button:hover {
  background: #1d4ed8;
}

.chart {
  height: 500px;
  width: 100%;
}

.error {
  margin-top: 20px;
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: red;
  font-size: 1.25rem;
}
</style>
