#### 智能串口调试助手

一个基于 Vue 3 和 Element Plus 构建的现代化串口通信调试工具。它专注于解决 Web 端串口通信中的数据编码与解码难题，提供了强大的 Hex（十六进制）数据处理能力，支持智能编码、解码以及严格的输入校验，旨在为开发者提供流畅的底层数据调试体验。

#### 核心特性

- **智能 Hex 输入框**：基于 Element Plus 深度封装，支持自动过滤非 Hex 字符、自动去除首尾及中间空格、自动转大写，完美兼容中文输入法。
- **多模式编码/解码**：
    - **Hex 模式**：支持发送和接收纯十六进制字节流（如 `AA 01`），自动处理字节对齐。
    - **ASCII 模式**：支持发送和接收人类可读的文本字符串（如 `Hello`），自动处理 UTF-8 编码。
- **数据智能转换**：内置智能算法，自动识别数字与文本，处理混合数据（如 `100px`）时自动切换编码策略，无需手动干预。
- **空格处理机制**：发送时自动剔除干扰空格，接收时自动格式化添加空格，完美平衡“机器解析”与“人类阅读”的需求。
- **用户体验优化**：解决了 `v-model.trim` 在组件库中的兼容性问题，采用 `@blur` 事件处理，确保输入流畅不卡顿。

#### 技术栈

- **框架**：Vue 3 (Composition API)、 Wails3
- **UI 库**：Element Plus
- **语言**：JavaScript / TypeScript、 Go
- **核心 API**：`TextEncoder`, `TextDecoder`, `DataView`, `Uint8Array`

#### 快速开始

1. **安装Node环境和Go环境**

2. **启动项目**

```
wails3 dev
```

3. **构建生产版本**

```
wails3 build
```

#### 功能模块详解

本项目的核心在于数据处理层，以下是主要模块的逻辑说明：

**1. 智能 Hex 输入组件**
为了解决 `v-model.trim` 在 Element Plus 中的兼容性问题，我们采用了 `@blur` 事件配合正则清洗的策略。

- **输入拦截**：监听 `keydown` 事件，直接阻止非 Hex 字符（非 0-9, A-F）的按键输入。
- **粘贴清洗**：监听 `paste` 事件，自动提取剪贴板中的 Hex 数据，去除空格和换行符。
- **IME 兼容**：通过 `compositionstart` 和 `compositionend` 监听中文输入法状态，避免在选词过程中误删字符。

**2. 串口编码工具类**
封装了 `SerialCodec` 对象，统一处理数据转换逻辑：

| 方法名 | 描述 | 示例输入 | 示例输出 |
| ------ |------ |------ |------ |
| `encodeHex` | 将数据转为字节流（Hex 模式） | `100` (Number) | `[0x64]` |
| `encodeAscii` | 将数据转为字节流（ASCII 模式） | `"100"` (String) | `[0x31, 0x30, 0x30]` |
| `decodeHex` | 将字节流转为 Hex 字符串 | `[0xAA, 0x01]` | `"AA 01"` |
| `decodeAscii` | 将字节流转为文本 | `[0x48, 0x69]` | `"Hi"` |

#### 项目结构

```
├── src
│   ├── components
│   │   ├── HexInput.vue       # 核心 Hex 输入组件
│   │   └── SerialMonitor.vue  # 串口监视器组件
│   ├── utils
│   │   └── codec.js           # 编解码核心逻辑
│   ├── App.vue
│   └── main.js
├── README.md
└── package.json
```

#### 使用示例

**发送数据**
在发送区输入 `FF 0A`（Hex 模式），点击发送，串口将实际输出两个字节：`0xFF` 和 `0x0A`。

**接收数据**
当串口收到字节 `[0x48, 0x65, 0x6C, 0x6C, 0x6F]`：

- 勾选 **Hex 显示**：界面显示 `48 65 6C 6C 6F`
- 勾选 **ASCII 显示**：界面显示 `Hello`

修改配置后需要运行以下命令

   wails3 task common:update:build-assets

   wails3 build

Mac平台打包

​	 wails3 package

Windows平台打包

​	   wails3 build GOOS=windows GOARCH=amd64

注意：

​	当前进行串口调试只支持Mac ARM平台，后续会增加Windows、Linux平台支持。



#### 许可证

MIT License
