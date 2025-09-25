# Shared Logging Package

统一的日志包，为所有SmartCam微服务提供标准化的日志功能。

## 特性

- 标准化的日志格式（JSON格式，使用ISO8601时间戳）
- 统一的HTTP请求日志中间件
- 标准化的GORM数据库日志
- 请求ID和关联ID支持
- 结构化的错误处理

## 基本使用

### 1. 创建Logger

```go
import "shared-logging"

// 生产环境
logger, err := logging.NewProductionLogger("my-service")
if err != nil {
    log.Fatal("Failed to initialize logger:", err)
}

// 开发环境
logger, err := logging.NewDevelopmentLogger("my-service")
if err != nil {
    log.Fatal("Failed to initialize logger:", err)
}

// 自定义配置
config := &logging.LoggerConfig{
    Level:       "info",
    Format:      "json",
    Output:      "stdout",
    ServiceName: "my-service",
}
logger, err := logging.NewLogger(config)
```

### 2. HTTP中间件

```go
import (
    "github.com/gin-gonic/gin"
    "shared-logging"
)

router := gin.New()

// 基础中间件
router.Use(logging.RequestID())           // 生成请求ID
router.Use(logging.Recovery(logger))      // 统一的panic恢复
router.Use(logging.HTTPLogger(logger))    // HTTP请求日志
router.Use(logging.CORS())               // CORS支持
```

### 3. GORM集成

```go
import (
    "gorm.io/gorm"
    "shared-logging"
)

// 使用默认配置
gormLogger := logging.NewGormLogger(logger, nil)

// 或自定义配置
config := &logging.GormLoggerConfig{
    LogLevel:                  gormlogger.Info,
    SlowThreshold:             200 * time.Millisecond,
    IgnoreRecordNotFoundError: true,
}
gormLogger := logging.NewGormLogger(logger, config)

// 在GORM中使用
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    Logger: gormLogger,
})
```

### 4. 结构化日志字段

```go
// 使用预定义字段
logger.Info("User logged in",
    logging.Field.UserID(123),
    logging.Field.RequestID("req-123"),
    logging.Field.HTTPMethod("POST"),
)

// 基本字段类型
logger.Info("Operation completed",
    logging.Field.String("operation", "create_user"),
    logging.Field.Int("count", 10),
    logging.Field.Bool("success", true),
    logging.Field.Error(err),
)
```

## 日志格式标准

### HTTP请求日志格式
```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "level": "info",
  "message": "HTTP request completed",
  "service_name": "api-gateway",
  "request_id": "req-12345",
  "http_method": "POST",
  "http_path": "/api/users",
  "http_status": 201,
  "client_ip": "192.168.1.1",
  "latency": "150ms",
  "user_agent": "Mozilla/5.0...",
  "body_size": 1024
}
```

### 数据库查询日志格式
```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "level": "debug",
  "message": "Database query executed",
  "service_name": "user-service",
  "request_id": "req-12345",
  "elapsed": "50ms",
  "rows_affected": 1,
  "sql": "SELECT * FROM users WHERE id = $1"
}
```

### 错误日志格式
```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "level": "error",
  "message": "Database connection failed",
  "service_name": "user-service",
  "error": "connection timeout",
  "caller": "service/user.go:45"
}
```

## 迁移指南

### 从现有日志系统迁移

1. 替换logger初始化：
```go
// 原来
logger, _ := zap.NewProduction()

// 现在
logger, _ := logging.NewProductionLogger("service-name")
```

2. 替换HTTP中间件：
```go
// 原来
router.Use(gin.Logger())

// 现在
router.Use(logging.HTTPLogger(logger))
```

3. 替换GORM logger：
```go
// 原来
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

// 现在
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    Logger: logging.NewGormLogger(logger, nil),
})
```

## 最佳实践

1. **服务名称一致性**：确保服务名称与Docker容器名称一致
2. **请求ID传播**：在微服务间传递请求ID
3. **结构化字段**：使用预定义字段而不是字符串拼接
4. **错误处理**：总是记录错误的上下文信息
5. **性能考虑**：在生产环境使用INFO级别，开发环境使用DEBUG级别