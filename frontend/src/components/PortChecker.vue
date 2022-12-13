<script setup>
import { reactive } from 'vue'
import { Greet, PortCheck } from '../../wailsjs/go/main/App'

const data = reactive({
  result_json: {},
  host: 'localhost',
  port: '58000',
  nums: '5',
  waiting: false
})

function check() {
  data.result_json = {}
  data.waiting = true
  PortCheck(data.host, Number(data.port), Number(data.nums)).then((result) => {
    data.result_json = JSON.parse(result)
    data.waiting = false
  })
}
</script>

<template>
  <main>
    <section class="section">
      <div class="container">
        <h1 class="title">Port Checker</h1>
        <div class="field">
          <label class="label">ホスト名 or IP : 開始ポート : x 個数</label>
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
          <button class="button is-small" @click="data.result_json = {}">
            結果消去
          </button>
          <button class="button is-primary is-small" @click="check">
            ポートチェック
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
