# mykonf

Go 配置管理库，支持从多个来源加载和合并配置，具有清晰的优先级顺序。

## 功能特性

- YAML 配置文件加载
- 环境变量覆盖（支持自定义前缀）
- 嵌套结构体配置支持
- JSON 字符串解析（用于复杂的 map 和 struct 字段）
- 配置文件中的环境变量展开（`$VAR` 或 `${VAR}` 语法）
- 通过 struct tag 设置默认值
- 类型自动转换（duration、slice、bool 等）

## 安装

```bash
go get github.com/empirefox/mykonf
```

## 配置项优先级

配置加载遵循以下优先级顺序（从高到低）：

| 优先级 | 来源 | 说明 |
|--------|------|------|
| 1 (最高) | 环境变量 | 必须带有指定的 `envPrefix` 前缀 |
| 2 | YAML 配置文件 | 默认 `config.yaml` 或自定义路径 |
| 3 (最低) | 默认值 | 通过 `default:""` struct tag 指定 |

## 使用方法

### 基本用法

```go
package main

import (
    "log"
    "github.com/empirefox/mykonf"
)

type Config struct {
    Port     int    `yaml:"port" default:"8080"`
    Host     string `yaml:"host" default:"localhost"`
    Database struct {
        Host string `yaml:"host" default:"localhost"`
        Port int    `yaml:"port" default:"5432"`
        Name string `yaml:"name"`
    } `yaml:"database"`
}

func main() {
    conf := &Config{}
    err := mykonf.Load("APP_", conf)
    if err != nil {
        log.Fatal(err)
    }
    // conf 已加载完成
}
```

### YAML 配置文件

默认读取当前目录下的 `config.yaml`：

```yaml
port: 9090
host: "0.0.0.0"

database:
  host: "db.example.com"
  port: 3306
  name: "myapp"
```

### 自定义配置文件路径

通过环境变量指定配置文件路径：

```bash
export APP_SERVER_CONFIG=/etc/myapp/config.yaml
```

### 环境变量覆盖

环境变量命名规则：
- 前缀 + 字段路径（大写，用下划线连接）
- 嵌套字段使用下划线分隔

```bash
# 基本字段
export APP_PORT=9000
export APP_HOST="0.0.0.0"

# 嵌套字段
export APP_DATABASE_HOST="prod.db.example.com"
export APP_DATABASE_PORT=5432
```

### 配置文件中使用环境变量

配置文件中可以引用环境变量：

```yaml
database:
  password: $DB_PASSWORD
  connection_string: "${DB_USER}:${DB_PASSWORD}@${DB_HOST}"
```

### JSON 字符串解析

对于 map 或 struct 类型的字段，可以通过 JSON 字符串传递：

```bash
export APP_CALLBACK_MAP='{"push":"http://webhook/push","pull":"http://webhook/pull"}'
```

```go
type Config struct {
    CallbackMap map[string]string `yaml:"callback_map"`
}
```

### 切片字段

切片字段在环境变量中使用逗号分隔：

```bash
export APP_HOSTS="host1,host2,host3"
```

```go
type Config struct {
    Hosts []string `yaml:"hosts"`
}
```

### Duration 字段

支持 Go 标准的 duration 格式：

```yaml
timeout: 30s
interval: 5m
```

```bash
export APP_TIMEOUT="1m30s"
```

## API 参考

### Load

```go
func Load(envPrefix string, conf any) error
```

使用默认配置路径加载配置。配置文件路径通过 `{envPrefix}SERVER_CONFIG` 环境变量指定，默认为 `config.yaml`。

### LoadPath

```go
func LoadPath(envPrefix, path string, conf any) error
```

从指定路径加载配置文件。

参数：
- `envPrefix`: 环境变量前缀（如 `APP_`）
- `path`: 配置文件路径
- `conf`: 配置结构体指针

### ConfigPath

```go
func ConfigPath(envPrefix string) string
```

获取配置文件路径。检查 `{envPrefix}SERVER_CONFIG` 环境变量，默认返回 `config.yaml`。

## 完整示例

```go
package main

import (
    "log"
    "time"
    "github.com/empirefox/mykonf"
)

type Config struct {
    Listen      string            `yaml:"listen" default:":8080"`
    CallbackMap map[string]string `yaml:"callback_map"`
    Timeout     time.Duration     `yaml:"timeout" default:"30s"`

    Database struct {
        Host     string `yaml:"host" default:"localhost"`
        Port     int    `yaml:"port" default:"5432"`
        Username string `yaml:"username"`
        Password string `yaml:"password"`
    } `yaml:"database"`

    Features []string `yaml:"features"`
}

func main() {
    conf := &Config{}
    if err := mykonf.Load("MYAPP_", conf); err != nil {
        log.Fatal(err)
    }

    log.Printf("Listening on %s", conf.Listen)
    log.Printf("Database: %s:%d", conf.Database.Host, conf.Database.Port)
}
```

配置文件 `config.yaml`：

```yaml
listen: ":9090"
timeout: 1m

database:
  host: "db.example.com"
  port: 3306
  username: "admin"
  password: $DB_PASSWORD

features:
  - auth
  - logging
  - metrics
```

环境变量覆盖：

```bash
export MYAPP_LISTEN=":8080"
export MYAPP_DATABASE_HOST="prod.db.example.com"
export MYAPP_CALLBACK_MAP='{"event":"http://hooks/event"}'
export DB_PASSWORD="secret"
```

## 依赖

- [koanf](https://github.com/knadh/koanf) - 配置管理框架
- [defaults](https://github.com/creasty/defaults) - 默认值处理
- [mapstructure](https://github.com/go-viper/mapstructure) - 结构体解码

## License

MIT
