<template>
  <section class="Spa">
    <div class="left-container">
      <div class="left-top-container">
        <el-input v-model="receivedText" type="textarea" readonly />
        <div class="opeartion-container">
          <el-checkbox v-model="autoScroll" :disabled="!formState.status">自动滚动</el-checkbox>
          <el-checkbox v-model="timeStamp" :disabled="!formState.status">时间戳</el-checkbox>
        </div>
        <div class="opeartion-container">
          <el-select v-model="formState.receiveMode" placeholder="模式" :disabled="!formState.status">
            <el-option v-for="item in Modes" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
          <el-select v-model="formState.code" placeholder="编码" :disabled="!formState.status">
            <el-option v-for="item in Codes" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
          <el-button :disabled="receiveEmpty" @click="ClearReceiveText">清空</el-button>
        </div>
      </div>
      <div class="left-bottom-container">
        <div class="send-container">
          <el-input v-model="sendText" :autosize="{ minRows: 3, maxRows: 3 }" show-word-limit type="textarea" />
          <div class="opeartion-container">
            <el-button :disabled="sendEmpty || !formState?.status" type="primary" @click="handleSendText">发送</el-button>
            <el-button :disabled="sendEmpty" @click="ClearSendText">清空</el-button>
          </div>
        </div>
      </div>
    </div>
    <div class="right-container">
      <div class="right-top-container">
        <el-form label-width="auto" label-position="left" :model="formState">
          <el-form-item label="串口号">
            <el-select v-model="formState.port" placeholder="串口号" :disabled="connectDisabled">
              <el-option v-if="Coms?.length" v-for="com in Coms" :key="com?.value" :label="com?.label"
                :value="com?.value" />
            </el-select>
          </el-form-item>
          <el-form-item label="波特率">
            <el-select v-model="formState.baudRate" :disabled="connectDisabled" placeholder="波特率">
              <el-option v-for="item in BaudRates" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>
          <el-form-item label="数据位">
            <el-select v-model="formState.dataBits" placeholder="数据位">
              <el-option v-for="item in DataBites" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>
          <el-form-item label="停止位">
            <el-select v-model="formState.stopBits" placeholder="停止位">
              <el-option v-for="item in StopBites" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>
          <el-form-item label="校验位">
            <el-select v-model="formState.checkBits" placeholder="校验位">
              <el-option v-for="item in CheckBites" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>
          <el-form-item label="发送模式">
            <el-select v-model="formState.sendMode" placeholder="发送模式">
              <el-option v-for="item in Modes" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>
          <el-form-item label="自动发送">
            <el-switch v-model="formState.autoSend" :disabled="!formState.status" />
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
          <el-button :loading="loading" :disabled="StatusDisabled" style="width: 100%" v-if="!formState?.status"
            type="danger" @click="openSearialPort">连接串口</el-button>
          <el-button :loading="loading" :disabled="StatusDisabled" style="width: 100%" v-else type="success"
            @click="closeSearialPort">断开连接</el-button>
        </div>
      </div>
    </div>
  </section>
</template>
<script lang="ts" setup>
import { ElMessage } from 'element-plus'
import { debounce } from 'lodash-es'
import { GetSerialPorts, OpenSerial, CloseSerial, SendData } from '../../../bindings/changeme/serialportservice'
import { Events } from '@wailsio/runtime'
import { computed, onMounted, ref, watch } from 'vue'
interface OpeartionType {
  value: string
  label: string
}
const receivedText = ref('')
const sendText = ref('')
const loading = ref(false)
const timeStamp = ref(false)
const autoScroll = ref(false)
const CheckBites = ref([
  { value: 'None', label: 'None' },
  { value: 'Even', label: 'Even' },
  { value: 'Odd', label: 'Odd' }
])
const Coms = ref<OpeartionType[]>([])
const BaudRates = ref([
  { value: 9600, label: 9600 },
  { value: 19200, label: 19200 },
  { value: 18400, label: 18400 }
])
const Modes = ref([
  { label: 'Hex', value: 'Hex' },
  { label: 'Text', value: 'Text' }
])
const DataBites = ref([
  { label: 5, value: 5 },
  { label: 6, value: 6 },
  { label: 7, value: 7 },
  { label: 8, value: 8 }
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
  baudRate: 9600,
  dataBits: 5,
  autoSend: false,
  stopBits: 1,
  checkBits: 'None',
  frequency: 1000,
  receiveMode: 'Text',
  sendMode: 'Text',
  status: false,
  code: ''
})
const ClearReceiveText = () => {
  receivedText.value = ''
}
const ClearSendText = () => {
  sendText.value = ''
}
const StatusDisabled = computed(() => {
  return !formState.value.baudRate || !formState.value.port || !formState.value.dataBits || !formState.value.stopBits || !formState.value.stopBits || !formState.value.checkBits
})
const sendEmpty = computed(() => {
  return !sendText.value?.length
})
const connectDisabled = computed(() => {
  return formState.value.status
})
const receiveEmpty = computed(() => {
  return !receivedText.value?.length
})
const handleSendText = debounce(async () => {
  // Events.On("serial_data", (data) => {
  //   console.log("收到串口数据:", data);
  //   receivedText.value = data?.data;
  // });
  try {
    await SendData(receivedText.value)
    ElMessage.success('数据发送成功')
  } catch (error) {
    ElMessage.error('数据发送失败')
  }
}, 1000)
const openSearialPort = debounce(async () => {
  try {
    loading.value = true
    await OpenSerial(formState.value?.port, formState.value?.baudRate, formState.value?.dataBits, formState.value?.stopBits)
    formState.value.status = true
    ElMessage.success('串口连接成功')
    loading.value = false
  } catch (error) {
    ElMessage.error('串口连接失败')
    loading.value = false
  }
  finally {
    loading.value = false

  }
}, 1000)
const closeSearialPort = debounce(async () => {
  try {
    loading.value = true
    console.log('8888')
    await CloseSerial()
    console.log('999')
    formState.value.status = false
    ElMessage.success('串口断开成功')
    loading.value = false
  } catch (error) {
    ElMessage.error('串口断开失败')
    loading.value = false
  } finally {
    loading.value = false
  }
}, 1000)
const init = async () => {
  const data = await GetSerialPorts()
  Coms.value = data?.map((item) => ({ label: item, value: item })) ?? []
}
class TimeInterval {
  interval: number
  fn: Function
  lastTime: number
  timer: any
  constructor(fn: Function, interval: number = 500) {
    this.interval = interval
    this.fn = fn
    this.lastTime = 0
    this.loop(0)
  }
  loop(timestamp: number) {
    this.timer = requestAnimationFrame(TimeInterval.prototype.loop.bind(this))
    if (timestamp - this.lastTime > this.interval) {
      this.lastTime = timestamp
      typeof this.fn == 'function' && this.fn()
    }
  }
  clear() {
    cancelAnimationFrame(this.timer)
    this.timer = null
  }
}

watch(
  () => formState.value.autoSend,
  (newVal, oldVal) => {
    console.log("自动", newVal)
    const timeInterval =  new TimeInterval(() => {
        handleSendText
    }, formState.value.frequency)
    if (newVal) {
      timeInterval
    } else {
      timeInterval.clear()
    }
  }
)


onMounted(() => {
  init()
})
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
