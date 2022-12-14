<script setup>
import { reactive } from 'vue'
import { Greet, PortCheck, NetInfo, SpeedTest } from '../../wailsjs/go/main/App'

const data = reactive({
  result_json: {},
  result_net: null,
  result_speed: null,
  host: 'localhost',
  port: '58000',
  nums: '5',
  waiting: false,
  waiting_speed: false
})

function check() {
  data.result_json = {}
  data.waiting = true
  PortCheck(data.host, Number(data.port), Number(data.nums)).then((result) => {
    data.result_json = JSON.parse(result)
    data.waiting = false
  })
}

function getif() {
  data.result_net = null
  NetInfo().then((result) => {
    data.result_net = JSON.parse(result)
  })
}

function getspeed() {
  data.result_speed = null
  data.waiting_speed = true
  SpeedTest().then((result) => {
    data.result_speed = JSON.parse(result)
    data.waiting_speed = false
  })
}
</script>

<template>
  <main>
    <section class="section">
      <div class="container">
        <h1 class="title">Port Checker</h1>
        <div class="field box">
          <span v-if="data.result_net">
            {{ data.result_net['global_ip'] }} / {{ data.result_net['isp'] }}
            <span v-for="i in data.result_net['interface']">
              {{ i.name }} : {{ i.ip }} /
            </span>
          </span>
          <button class="button is-primary" @click="getif">Network Info</button>
        </div>
        <div class="field box">
          <span v-if="data.waiting_speed"> Speed Testing ...</span>
          <span v-if="data.result_speed">
            Down :
            {{ Math.round(data.result_speed[0]['download'] * 100) / 100 }} Mbps
            / Up :
            {{ Math.round(data.result_speed[0]['upload'] * 100) / 100 }} Mbps /
            Ping :{{ data.result_speed[0]['latency'] }}
          </span>
          <button class="button is-primary" @click="getspeed">SpeedTest</button>
        </div>
        <div class="field">
          <label class="label">Host or IP : Start Port : x Nums</label>
          <div class="field is-grouped">
            <div class="control">
              <input
                v-model="data.host"
                autocomplete="off"
                class="input is-primary is-small"
                type="text"
                placeholder="host name or ip address"
              />
            </div>
            <div class="control">
              <input
                v-model="data.port"
                autocomplete="off"
                class="input is-primary is-small"
                type="text"
                placeholder="start port number"
              />
            </div>
            <div class="control">
              <input
                v-model="data.nums"
                autocomplete="off"
                class="input is-primary is-small"
                type="text"
                placeholder="number of ports to check"
              />
            </div>
          </div>
        </div>
        <div class="field">
          <button class="button m-1" @click="data.result_json = {}">
            Clear
          </button>
          <button class="button m-1 is-primary" @click="check">
            Start Scan
          </button>
        </div>
      </div>
      <br />
      {{ data.waiting ? 'Waiting Response ...' : '' }}
      <table class="table">
        <tbody>
          <tr v-for="p in data.result_json['ports']">
            <td>
              <b>{{ p.port }}</b>
            </td>
            <td v-if="p.open">
              <button class="button is-small is-success">OPEN</button>
            </td>
            <td v-if="!p.open">
              <button class="button is-small is-danger">CLOSED</button>
            </td>
            <td>{{ p.status }}</td>
          </tr>
        </tbody>
      </table>
    </section>
  </main>
</template>
