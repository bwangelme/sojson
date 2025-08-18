# SoJSON - JSON 工具

一个使用 Go + Gin 后端和纯 JavaScript 前端开发的 JSON 去除转义和格式化工具。

## 功能特性

- **JSON 去除转义**：将转义的 JSON 字符串转换为正常的 JSON
- **JSON 格式化**：美化 JSON 显示，支持可配置的缩进（2空格、4空格、压缩）
- **JSON 验证**：检查 JSON 格式是否正确
- **组合处理**：一键去除转义并格式化
- **实时处理**：输入即时显示结果
- **错误提示**：详细的 JSON 格式错误信息
- **现代界面**：响应式设计，支持移动端
- **键盘快捷键**：提高使用效率
- **文件下载**：支持将结果下载为文件

## 技术栈

- **后端**: Go 1.21 + Gin Web Framework
- **前端**: 纯 JavaScript (ES6+) + HTML5 + CSS3
- **跨域**: CORS 支持
- **部署**: 可独立部署，无需额外依赖

## 安装运行

### 环境要求
- Go 1.21 或更高版本

### 快速开始

1. **克隆项目**
```bash
git clone <项目地址>
cd sojson
```

2. **安装依赖**
```bash
go mod tidy
```

3. **运行项目**
```bash
go run main.go
```

4. **访问应用**
打开浏览器访问: `http://localhost:8080`

### 生产部署

```bash
# 编译为可执行文件
go build -o sojson main.go

# 运行
./sojson
```

## API 接口

### 基础信息
- **Base URL**: `http://localhost:8080/api`
- **Content-Type**: `application/json`

### 接口列表

#### 1. 去除转义
```http
POST /api/unescape
Content-Type: application/json

{
    "text": "转义的JSON字符串"
}
```

#### 2. 格式化 JSON
```http
POST /api/format
Content-Type: application/json

{
    "text": "JSON字符串",
    "indent": 2  // 可选，默认为2
}
```

#### 3. 完整处理（去除转义+格式化）
```http
POST /api/process
Content-Type: application/json

{
    "text": "转义的JSON字符串",
    "indent": 2  // 可选，默认为2
}
```

#### 4. 验证 JSON
```http
POST /api/validate
Content-Type: application/json

{
    "text": "要验证的JSON字符串"
}
```

### 响应格式

#### 成功响应
```json
{
    "result": "处理后的JSON字符串",
    "success": true
}
```

#### 验证响应
```json
{
    "valid": true,
    "message": "JSON格式正确"
}
```

#### 错误响应
```json
{
    "error": "错误信息",
    "success": false
}
```

## 使用说明

### 功能操作

1. **选择功能**：点击顶部的功能按钮选择要执行的操作
2. **输入数据**：在左侧输入框中粘贴或输入 JSON 数据
3. **设置缩进**：在设置面板中选择缩进方式
4. **处理数据**：点击"处理"按钮或使用快捷键
5. **查看结果**：在右侧输出框中查看处理结果
6. **复制下载**：使用操作按钮复制或下载结果

### 键盘快捷键

- `Ctrl/Cmd + Enter`：执行处理
- `Ctrl/Cmd + L`：清空输入
- `Ctrl/Cmd + V`：粘贴到输入框
- `Ctrl/Cmd + C`：复制输出结果
- `Esc`：隐藏提示消息

### 示例数据

应用内置了示例数据，点击设置面板中的示例按钮即可快速体验各种功能。

## 项目结构

```
sojson/
├── main.go              # 主程序文件
├── go.mod              # Go 模块文件
├── templates/          # HTML 模板
│   └── index.html     # 主页模板
├── static/            # 静态文件
│   ├── css/
│   │   └── style.css  # 样式文件
│   └── js/
│       └── app.js     # 前端逻辑
└── README.md          # 项目说明
```

## 开发相关

### 本地开发
```bash
# 启用开发模式（自动重载）
go run main.go

# 或使用 air（需要安装）
air
```

### 代码规范
- 后端遵循 Go 标准规范
- 前端使用 ES6+ 标准
- 注释使用中文，提高可读性

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request 来改进这个项目。
