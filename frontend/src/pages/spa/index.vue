<template>
  <section class="Spa">
    <div class="left-container">
      <div class="left-top-container">
        <el-input ref="textareaRef" v-model="receivedText" type="textarea" readonly />
        <div class="opeartion-container">
          <el-checkbox v-model="timeStamp">显示时间戳</el-checkbox>
          <el-checkbox v-model="autoScroll">自动滚动</el-checkbox>
           <el-select v-model="formState.receiveMode" placeholder="编码模式">
            <el-option v-for="item in Codes" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
          <el-button :disabled="receiveEmpty" @click="ClearReceiveText">清空</el-button>
        </div>
      </div>
      <div class="left-bottom-container">
        <div class="send-container">
          <el-input v-model.trim="sendText" :autosize="{ minRows: 3, maxRows: 3 }" show-word-limit type="textarea" @change="(value) => { isHex && writeHex(value) }"
/>
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
            <el-select v-model="formState.dataBits" placeholder="数据位" :disabled="connectDisabled">
              <el-option v-for="item in DataBites" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>
          <el-form-item label="停止位">
            <el-select v-model="formState.stopBits" placeholder="停止位" :disabled="connectDisabled">
              <el-option v-for="item in StopBites" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>
          <el-form-item label="校验位">
            <el-select v-model="formState.parityMode" placeholder="校验位" :disabled="connectDisabled">
              <el-option v-for="item in ParityModes" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>
          <el-form-item label="编码模式">
            <el-select v-model="formState.sendMode" placeholder="编码模式" :disabled="connectDisabled" @change="handleStringHex">
              <el-option v-for="item in Codes" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>
          <el-form-item label="自动发送">
            <el-switch v-model="formState.autoSend" :disabled="!formState.status" />
          </el-form-item>
          <el-form-item v-if="formState?.autoSend" label="发送周期">
            <el-input-number style="width: 100%" v-model="formState.frequency" :min="100">
              <template #suffix>
                <span>ms</span>
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
import dayjs from 'dayjs'
import { ElMessage } from 'element-plus'
import { debounce, values } from 'lodash-es'
import { GetSerialPorts, OpenSerial, CloseSerial, SendData } from '../../../bindings/changeme/serialportservice'
import { Events } from '@wailsio/runtime'
import { computed, nextTick, onMounted, ref, watch } from 'vue'
import { baudRates, codes, createPersistentTask, dataBites, hexToString, parityModes, stopBites, stringToHex } from './datas'
interface OpeartionType {
  value: string
  label: string
}
interface DataType {
  timeStamp: string
  text: string
}
let timeInterval: any = null
const sendText = ref('')
const loading = ref(false)
const timeStamp = ref(false)
const autoScroll = ref(false)
const textareaRef = ref()
const ParityModes = ref(parityModes)
const Coms = ref<OpeartionType[]>([])
const BaudRates = ref(baudRates)
const receiveArrayText = ref<DataType[]>([])
const receivedText = computed(() => {
  if (timeStamp.value) return receiveArrayText.value?.map((cur) => {
    return `【${cur?.timeStamp}】${cur?.text}`
   }).join('\n')
  return receiveArrayText.value.map((cur) => {
    return `${cur?.text}`
   }).join('\n')
})
setInterval(() => {
  receiveArrayText.value.push({timeStamp: dayjs(new Date()).format('YYYY-MM-DD HH:mm:ss.SSS'), text: "888"})
}, 10000)
const DataBites = ref(dataBites)
const StopBites = ref(stopBites)
const Codes = ref(codes)
const formState = ref({
  port: '',
  baudRate: 9600,
  dataBits: 8,
  autoSend: false,
  stopBits: 1,
  parityMode: 0,
  frequency: 1000,
  receiveMode: 'UTF-8',
  sendMode: 'UTF-8',
  status: false,
})
const ClearReceiveText = () => {
  receiveArrayText.value = []
}
const ClearSendText = () => {
  sendText.value = ""
}
const StatusDisabled = computed(() => {
  return !formState.value.baudRate || !formState.value.port || !formState.value.dataBits || !formState.value.stopBits || !formState.value.stopBits
})
const sendEmpty = computed(() => {
  return !sendText.value?.length
})
const connectDisabled = computed(() => {
  return formState.value.status
})
const receiveEmpty = computed(() => {
  return !receiveArrayText.value?.length
})
const handleSendText = debounce(async () => {
  // Events.On("serial_data", (data) => {
  //   console.log("收到串口数据:", data);
  //   receivedText.value = data?.data;
  // });
  console.log("执行时间戳", new Date().getTime())
  try {
    await SendData(sendText.value)
    ElMessage.success('数据发送成功')
  } catch (error) {
    ElMessage.error('数据发送失败')
  }
}, 1000)
const openSearialPort = debounce(async () => {
  try {
    loading.value = true
    await OpenSerial(formState.value?.port, formState.value?.baudRate, formState.value?.dataBits, formState.value?.stopBits, formState.value?.parityMode)
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

watch(() => formState.value.autoSend, (newVal, _) => {
   if (newVal) {
    timeInterval = createPersistentTask(handleSendText, formState.value.frequency, { useBackgroundTime: true });
    console.log("开始", new Date().getTime())
  }
  if (!newVal && timeInterval) {
    timeInterval()
    console.log("关闭")
  }
});
const isHex = computed(() => {
  return formState.value.sendMode === "HEX"
})
const isUtf8 = computed(() => {
  return formState.value.sendMode === "UTF-8"
})
const writeString = (value) => {
   sendText.value = value
}
const writeHex = (value) => {
  sendText.value = stringToHex(value)
}
const handleStringHex = () => {
  if (isHex.value) {
    writeHex(sendText.value)
  }
  if (isUtf8.value) {
    writeString(hexToString(sendText.value))
  }
  console.log('999',formState.value.sendMode, sendText.value)
}

const scrollBottom = () => {
  nextTick(() => {
    requestAnimationFrame(() => {
      const textarea = textareaRef.value?.$el.querySelector('textarea')
      textarea.scrollTop = textarea.scrollHeight;
    })
  })
}

watch(() => receivedText.value, (val) => {
  if (autoScroll.value) {
    scrollBottom()
  }
})
watch(() => formState.value.frequency, (val) => {
  timeInterval()
  timeInterval = createPersistentTask(handleSendText, val, { useBackgroundTime: true });
});

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
