# Live2D 桌面宠物 🤖

一个基于 Live2D 的桌面宠物应用程序，可以作为桌面置顶浮窗运行。

## 🎯 项目简介

本项目是一个**Live2D桌面宠物**，具有以下核心功能：

- 🖼️ **Live2D模型显示** - 支持在桌面显示动态 Live2D 角色
- 🎮 **动作控制** - 点击按钮触发不同的角色动作
- 🎵 **语音播放** - 输入文字后播放预设音频
- 🖥️ **桌面置顶** - 作为浮窗始终显示在桌面最顶层
- 🖱️ **拖拽移动** - 可自由拖拽调整位置

## 🚀 快速开始

### 运行网页版

```bash
cd live2d/Samples/TypeScript/Demo
npm install
npm run start
```

访问 `http://localhost:5000`

### 运行桌面版（Electron）

```bash
cd live2d-desktop
npm install
npm start
```

## 📖 目录结构

```
├── live2d/                    # Live2D SDK 源码
│   ├── Core/                  # Live2D 核心库
│   ├── Framework/             # 框架代码
│   └── Samples/               # 示例项目
│       └── TypeScript/Demo/   # 网页版Demo
│
├── live2d-desktop/            # Electron 桌面版
│   ├── main.js               # Electron 主进程
│   └── index.html            # 前端入口
│
└── README.md                  # 项目说明
```

## 🔧 功能扩展指南

如果你想完善以下功能，以下是开发思路：

### 1. 🎤 语音识别

**目标**：将用户的语音输入转换为文字

**方案**：
- **浏览器API**（推荐）：使用 `Web Speech API`
  ```javascript
  const recognition = new webkitSpeechRecognition();
  recognition.lang = 'zh-CN';
  recognition.onresult = (event) => {
    const text = event.results[0][0].transcript;
    console.log('识别结果:', text);
  };
  recognition.start();
  ```
- **第三方API**：如阿里云、腾讯云的语音识别服务

### 2. 😊 情绪化表达

**目标**：让角色根据文字内容展现不同的表情

**方案**：
- 准备多个表情文件（`.exp3.json`）
- 根据识别到的文字情感（正面/负面/中性）切换表情
- 使用 Live2D 的 Expression 功能实现表情过渡

### 3. 🎭 选择情绪

**目标**：允许用户选择角色的情绪状态

**方案**：
- 在界面添加情绪选择按钮（开心、悲伤、生气、惊讶等）
- 点击后调用 `model.setExpression('happy')` 切换表情

### 4. 🛠️ 调用工具

**目标**：根据用户输入执行相应操作（如查询天气、播放音乐等）

**方案**：
```
用户输入 → 语音识别 → 文字理解 → 判断意图 → 调用工具 → 反馈结果 → 语音合成/动作反应
```

示例工具调用逻辑：
```javascript
const tools = {
  '天气': () => fetchWeather(),
  '播放音乐': () => playMusic(),
  '打开浏览器': () => openBrowser()
};

function handleInput(text) {
  for (const [keyword, handler] of Object.entries(tools)) {
    if (text.includes(keyword)) {
      return handler();
    }
  }
}
```

## 🎨 自定义模型

1. 将你的 Live2D 模型文件（`.moc3`、`.model3.json`、纹理图片等）放入 `live2d/Samples/TypeScript/Demo/public/Resources/`
2. 在 `lappdefine.ts` 中修改 `ModelDir` 数组

## 📝 技术栈

- **前端**：TypeScript + Vite
- **3D引擎**：Live2D Cubism SDK
- **桌面**：Electron
- **语音识别**：Web Speech API / 第三方API

## 📄 许可证

本项目使用的 Live2D Cubism SDK 遵循 Live2D 开源许可协议。

详见 [LICENSE](live2d/LICENSE.md)
