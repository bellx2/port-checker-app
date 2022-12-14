<script setup>
import { reactive } from 'vue'
import { Greet, PortCheck, NetInfo, SpeedTest } from '../../wailsjs/go/main/App'

const data = reactive({
  result_json: {},
  result_net: null,
  result_speed: [],
  host: 'localhost',
  port: '58000',
  nums: '5',
  waiting_scan: false,
  waiting_speed: false
})

function init() {
  get_interface()
}

init()

function do_portscan() {
  data.result_json = {}
  data.waiting_scan = true
  PortCheck(data.host, Number(data.port), Number(data.nums)).then((result) => {
    data.result_json = JSON.parse(result)
    data.waiting_scan = false
  })
}

function get_interface() {
  data.result_net = null
  NetInfo().then((result) => {
    data.result_net = JSON.parse(result)
  })
}
function clear_speedtest() {
  data.result_speed = []
}

async function do_speedtest() {
  data.waiting_speed = true
  SpeedTest().then((result) => {
    data.result_speed.push(JSON.parse(result)[0])
    data.waiting_speed = false
  })
  // var result,
  //   d,
  //   result_json = null
  // result = await SpeedTest()
  // result_json = await JSON.parse(result)
  // var d = []
  // d.append(result_json[0])
  // data.result_speed = d
  // setTimeout(function () {
  //   console.log('I am the third log after 5 seconds')
  // }, 5000)
  // result = await SpeedTest()
  // result_json = await JSON.parse(result)
  // d = data.result_speed
  // d.append(result_json[0])
  // data.result_speed = d
  // data.waiting_speed = false
}
</script>

<template>
  <main>
    <nav class="navbar is-light" role="navigation" aria-label="main navigation">
      <div class="navbar-brand">
        <div class="navbar-item" href="https://bulma.io">
          <h1 class="title">Port Checker</h1>
        </div>
      </div>
    </nav>
    <section class="section">
      <div class="container">
        <div class="box">
          <div class="field">
            <label class="label">Interfaces</label>
            <span v-if="data.result_net">
              <span class="tag is-link">Global IP</span>
              {{ data.result_net['global_ip'] }} (ISP:
              {{ data.result_net['isp'] }})
              <br />
              <span v-for="i in data.result_net['interface']">
                <span class="tag is-link is-light">{{ i.name }}</span>
                {{ i.ip }}
              </span>
            </span>
          </div>
          <button class="button is-primary" @click="get_interface">
            Reload
          </button>
        </div>
        <div class="box">
          <div class="field">
            <label class="label">Speed Test (fast.com)</label>
            <div v-for="r in data.result_speed">
              Down :
              {{ Math.round(r['download'] * 100) / 100 }}
              Mbps / Up :
              {{ Math.round(r['upload'] * 100) / 100 }} Mbps / Ping :{{
                r['latency']
              }}
            </div>
            <span v-if="data.waiting_speed"> Speed Testing ...</span>
          </div>
          <button class="button m-1" @click="clear_speedtest">Clear</button>
          <button
            class="button is-primary m-1"
            @click="do_speedtest"
            :class="{ 'is-loading': data.waiting_speed }"
          >
            SpeedTest
          </button>
        </div>
        <div class="box">
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
                  placeholder="number of ports to do_portscan"
                />
              </div>
            </div>
          </div>
          <div class="field">
            <button class="button m-1" @click="data.result_json = {}">
              Clear
            </button>
            <button
              class="button m-1 is-primary"
              @click="do_portscan"
              :class="{ 'is-loading': data.waiting_scan }"
            >
              Start Scan
            </button>
          </div>
        </div>
      </div>
      <br />
      {{ data.waiting_scan ? 'Waiting Response ...' : '' }}
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
