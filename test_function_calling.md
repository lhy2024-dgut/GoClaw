# 函数调用测试

## 功能说明
现在 GoClaw 支持 AI 自动执行本地操作，无需手动输入命令。

## 可用工具
1. **文件系统工具 (filesystem)**:
   - `create_dir`: 创建文件夹
   - `list_dir`: 列出目录内容
   - `read_file`: 读取文件
   - `write_file`: 写入文件
   - `delete`: 删除文件/文件夹

## 测试示例

### 1. 创建文件夹
用户: "请在当前目录创建一个名为 test_folder 的文件夹"
AI 会自动调用 `filesystem` 工具的 `create_dir` 操作

### 2. 列出目录
用户: "列出当前目录的内容"
AI 会自动调用 `filesystem` 工具的 `list_dir` 操作

### 3. 读取文件
用户: "读取 test_file.txt 的内容"
AI 会自动调用 `filesystem` 工具的 `read_file` 操作

### 4. 写入文件
用户: "创建一个名为 hello.txt 的文件，内容为 'Hello World'"
AI 会自动调用 `filesystem` 工具的 `write_file` 操作

### 5. 删除文件
用户: "删除 hello.txt 文件"
AI 会自动调用 `filesystem` 工具的 `delete` 操作

## 使用方法
1. 启动服务器: `./goclaw.exe serve`
2. 访问 `http://localhost:8080/chat`
3. 输入自然语言指令，AI 会自动执行相应操作