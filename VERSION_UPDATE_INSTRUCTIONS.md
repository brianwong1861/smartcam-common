# SmartCam Common v1.3.0 - 时间格式更新指南

## 📋 更新内容

已成功修改 `shared-logging/logger.go` 文件，实现了以下更改：

### ✅ 完成的修改

1. **添加自定义时间编码器**
   ```go
   // CustomTimeEncoder encodes time in YYYYMMDD:HHMMSS format
   func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
       enc.AppendString(t.Format("20060102:150405"))
   }
   ```

2. **更新时间格式配置**
   ```go
   TimeKey:    "time",           // 从 "timestamp" 改为 "time"
   EncodeTime: CustomTimeEncoder, // 从 ISO8601 改为自定义格式
   ```

3. **时间格式示例**
   - **旧格式**: `"ts":1758847903.5403278`
   - **新格式**: `"time":"20250926:143022"`

## 🚀 版本发布步骤

请按以下步骤手动完成版本更新：

### 第1步：进入smartcam-common目录
```bash
cd /Users/brianwong/Metair/Repo/SmartCam/backend/smartcam-common
```

### 第2步：提交更改
```bash
# 添加修改的文件到Git
git add shared-logging/logger.go

# 提交更改
git commit -m "feat: update logging time format to YYYYMMDD:HHMMSS

- Add CustomTimeEncoder for better readability
- Change timestamp format from ISO8601 to YYYYMMDD:HHMMSS
- Update time key from 'timestamp' to 'time'

🤖 Generated with Claude Code"
```

### 第3步：创建版本标签
```bash
# 创建新版本标签
git tag -a v1.3.0 -m "Version v1.3.0 - Custom time format for logging

Features:
- Custom time encoder with YYYYMMDD:HHMMSS format
- Improved log readability for users
- Maintains all existing logging functionality

🤖 Generated with Claude Code"
```

### 第4步：推送到远程仓库
```bash
# 推送代码和标签
git push origin main
git push origin v1.3.0
```

## 📦 更新微服务依赖

版本发布后，更新所有微服务的go.mod文件：

### 更新命令
```bash
# 在每个微服务目录中运行
go get github.com/brianwong1861/smartcam-common@v1.3.0
go mod tidy
```

### 需要更新的服务
- iam-service
- api-gateway-service
- dealer-management-service
- device-management-service
- financial-management-service
- notification-management-service
- video-stream-service
- realtime-communication-service

## 🔄 重启服务以应用更改

```bash
# 停止所有服务
docker-compose down

# 重新构建并启动服务
docker-compose build --no-cache
docker-compose up -d
```

## ✅ 验证新的日志格式

更新后，日志应该显示如下格式：
```json
{
  "level": "info",
  "time": "20250926:143022",
  "caller": "main.go:42",
  "message": "Service started successfully",
  "service_name": "iam-service"
}
```

---

**注意**: 所有使用 `shared-logging` 的服务都会自动采用新的时间格式，无需修改各个服务的代码。