<template>
  <section class="Spa">
    <div class="left-container">
      <div class="left-top-container">
        <el-input v-model="receivedText" type="textarea" disabled />
        <div class="opeartion-container">
          <el-select v-model="formState.receiveMode" placeholder="模式" clearable>
            <el-option v-for="item in Modes" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
          <el-select v-model="formState.code" placeholder="编码">
            <el-option v-for="item in Codes" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
          <el-button :disabled="receiveEmpty" @click="ClearReceiveText">清空</el-button>
        </div>
      </div>
      <div class="left-bottom-container">
        <div class="send-container">
          <el-input v-model="sendText" :autosize="{ minRows: 3, maxRows: 3 }" show-word-limit type="textarea" />
          <div class="opeartion-container">
            <el-button :disabled="sendEmpty" type="primary" @click="handleSendText">发送</el-button>
            <el-button :disabled="sendEmpty" @click="ClearSendText">清空</el-button>
          </div>
        </div>
      </div>
    </div>
    <div class="right-container">
      <div class="right-top-container">
        <el-form label-width="auto" label-position="left" :model="formState">
          <el-form-item label="串口号">
            <el-select v-model="formState.port" placeholder="串口号">
              <el-option v-for="item in Coms" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>
          <el-form-item label="波特率">
            <el-select v-model="formState.baudRate" placeholder="波特率">
              <el-option v-for="item in BaudRates" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>
          <el-form-item label="数据位">
            <el-select v-model="formState.dataBite" placeholder="数据位">
              <el-option v-for="item in DataBites" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>
          <el-form-item label="停止位">
            <el-select v-model="formState.stopBite" placeholder="停止位">
              <el-option v-for="item in StopBites" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>
          <el-form-item label="校验位">
            <el-select v-model="formState.checkBite" placeholder="校验位">
              <el-option v-for="item in CheckBites" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>
          <el-form-item label="发送模式">
            <el-select v-model="formState.sendMode" placeholder="发送模式" clearable>
              <el-option v-for="item in Modes" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>
          <el-form-item label="自动发送">
            <el-switch v-model="formState.autoSend" />
          </el-form-item>
          <el-form-item v-if="formState?.autoSend" label="发送周期">
            <el-input-number style="width: 100%" v-model="formState.frequency" :min="1">
              <template #suffix>
                <span>毫秒</span>
              </template>
            </el-input-number>
          </el-form-item>
        </el-form>
      </div>
      <div class="right-top-container">
        <div class="button">
          <el-button :disabled="StatusDisabled" style="width: 100%" v-if="!formState?.status" type="danger" @click="handlePortStatue(true)">连接串口</el-button>
          <el-button :disabled="StatusDisabled" style="width: 100%" v-else type="success" @click="handlePortStatue(false)">断开串口</el-button>
        </div>
      </div>
    </div>
  </section>
</template>
<script lang="ts" setup>
import { computed, ref } from 'vue'
const receivedText = ref('')
const sendText = ref('')
const CheckBites = ref([
  { value: 'None', label: 'None' },
  { value: 'Even', label: 'Even' },
  { value: 'Odd', label: 'Odd' }
])
const Coms = ref([
  { value: 'COM1', label: 'COM1' },
  { value: 'COM1', label: 'COM1' },
  { value: 'COM1', label: 'COM1' }
])
const BaudRates = ref([
  { value: 2400, label: 2400 },
  { value: 4800, label: 4800 },
  { value: 9600, label: 9600 },
  { value: 19200, label: 19200 },
  { value: 18400, label: 18400 }
])
const Modes = ref([{ label: 'HEX', value: 'HEX' }])
const DataBites = ref([
  { label: '5', value: '5' },
  { label: '6', value: '6' },
  { label: '7', value: '7' },
  { label: '8', value: '8' }
])
const StopBites = ref([
  { label: 1, value: 1 },
  { label: 2, value: 2 }
])
const Codes = ref([
  { label: 'ASCII', value: 'ASCII' },
  { label: 'UTF-8', value: 'UTF-8' },
  { label: 'UTF-16', value: 'UTF-16' }
])
const formState = ref({
  port: '',
  baudRate: null,
  dataBite: null,
  autoSend: false,
  stopBite: null,
  checkBite: null,
  frequency: 1000,
  receiveMode: '',
  sendMode: '',
  status: false,
  code: ''
})
const handlePortStatue = (status: boolean) => {
  formState.value.status = status
}
const ClearReceiveText = () => {
  receivedText.value = ''
}
const ClearSendText = () => {
  sendText.value = ''
}
const StatusDisabled = computed(() => {
  return !formState.value.baudRate || !formState.value.port || !formState.value.dataBite || !formState.value.stopBite || !formState.value.stopBite || !formState.value.checkBite
})
const sendEmpty = computed(() => {
  return !sendText.value?.length
})
const receiveEmpty = computed(() => {
  return !receivedText.value?.length
})
const handleSendText = () => {}
</script>
<style lang="scss" scoped>
.Spa {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: space-between;
  .left-container {
    flex: 1;
    display: flex;
    row-gap: 14px;
    flex-direction: column;
    .left-top-container {
      flex: 1;
      display: flex;
      flex-direction: column;
      row-gap: 10px;
      padding: 14px;
      border-bottom: 1px solid rgb(198, 198, 198);
      .el-textarea {
        flex: 1;
        ::v-deep(.el-textarea__inner) {
          height: 100%;
        }
      }

      .opeartion-container {
        display: flex;
        column-gap: 10px;

        button {
          margin: 0;
        }
      }
    }

    .left-bottom-container {
      display: flex;
      justify-content: center;
      flex-direction: column;
      padding: 10px;
      .send-container {
        display: flex;
        justify-content: space-between;
        align-items: center;
        column-gap: 10px;
        .opeartion-container {
          display: flex;
          flex-direction: column;
          justify-content: space-between;
          row-gap: 8px;
          button {
            margin: 0;
          }
        }
      }
    }
  }

  .right-container {
    border-left: 1px solid rgb(207, 207, 207);
    width: 240px;
    padding: 12px;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    align-content: center;
    .button {
      flex: 1;
      padding: 6px 10px;
      button {
        margin: 0;
      }
    }
  }
}
</style>
