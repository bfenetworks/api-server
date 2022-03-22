# 配置文件说明
本文讲解如何配置 API Server，让 API Server可以正常运行。

## 概述

所有的配置文件都在配置文件根路径下。配置文件根路径在API Server启动时由启动命令的 -c 参数指定，默认为启动路径下的 `./conf`。

配置文件包括：

- api_server.toml: API Server 模块启动必须的参数
- nav_tree.toml: 导航栏相关配置
- i18n/: 国际化相关的配置，主要是错误提示信息等内容

下面逐个进行说明。

## api_server.toml 

该文件提供API Server的启动参数，分为下面几部分，分别说明。

### API Server Config

API Server的关键配置。

| 配置项            | 描述                                              |
| ----------------- | ------------------------------------------------- |
| ServerPort        | Int<br>API Server的服务端口                       |
| GracefulTimeOutMs | Int<br>API Server关闭时的优雅退出时间，单位为毫秒 |
| MonitorPort       | Int<br>监控端口                                   |

示例：

```
# API Server Config
[Server]
# server port
ServerPort          = 8183
# server graceful exit timeout
GracefulTimeOutMs   = 5000
# monitor port, don't start monitor server if less than 0
MonitorPort         = 8284
```
### Log Config

API Server的日志相关配置。这里的日志指API Server的access log（访问日志）和sql log（数据库日志）。

可以在Log Config分别配置上述两种日志的配置。其中sql log部分配置仅在Runtime Config中的RecordSQL为"true"时需要配置。

两种日志均包含如下各配置项：

| 配置项      | 描述                                                         |
| ----------- | ------------------------------------------------------------ |
| LogName     | String<br>日志名称                                           |
| LogLevel    | String<br>日志级别<br>合法值："DEBUG", "TRACE", "INFO", "WARNING", "ERROR", "CRITICAL" |
| RotateWhen  | String<br/>日志轮转时间<br/>合法值：包括"M","H","D","MIDNIGHT"<br>  "M" 表示每分钟轮转<br>   "H" 表示每小时轮转<br>   "D" 表示每天轮转<br>   "MIDNIGHT" 表示每天0点整轮转 |
| BackupCount | Int<br>日志删除前轮转次数，即最大的日志存储数量              |
| Format      | String<br>日志各字段记录和排列的格式。支持的字段包括：<br>  %T - Time (15:04:05 MST)<br>  %t - Time (15:04)<br>  %D - Date (2006/01/02)<br>  %d - Date (01/02/06)<br>  %L - Level<br>  %P - Pid of process<br>  %S - Source<br>  %M - Message |
| StdOut      | Bool<br>是否输出日志到StdOut                                 |

示例：

```
# Logger Config
# access log config
[Loggers.access]
LogName     = "access"
LogLevel    = "INFO"
RotateWhen  = "MIDNIGHT"
BackupCount = 1
Format      = "[%D %T] [%L] [%S] %M"
StdOut      = false

# sql log, you can skip this config if RunTime.RecordSQL is false
[Loggers.sql]
LogName     = "sql"
LogLevel    = "INFO"
RotateWhen  = "MIDNIGHT"
BackupCount = 1
Format      = "[%D %T] %M"
StdOut      = false
```

### Database Config

API Server需要将配置存放在MySQL数据库中，并使用[go-sql-driver](https://github.com/go-sql-driver/mysql)来访问数据库。

Database Config用于指定go-sql-drive所使用的配置，配置参数的具体说明可以参考[go-sql-driver官方文档](https://pkg.go.dev/github.com/go-sql-driver/mysql#section-readme)。

主要配置项如下：

| 配置项               | 描述                                                         |
| -------------------- | ------------------------------------------------------------ |
| DBName               | String<br>数据库名                                           |
| Addr                 | String<br>数据库地址<br>Net参数设置为"tcp"时，格式为"IP:Port"，如"127.0.0.1:3306" |
| Net                  | String<br>网络类型<br>合法值：<br>  "tcp" - 使用TCP连接数据库<br>  "unix" - 使用Unix Domain Socket连接数据库 |
| User                 | String<br>数据库用户名                                       |
| Passwd               | String<br>数据库密码                                         |
| MultiStatements      | Bool<br>是否允许一个SQL Query中包含多个Statements            |
| MaxAllowedPacket     | Int<br>Mysql 服务器端允许的最大数据包大小                    |
| ParseTime            | Bool<br>是否自动将日期和时间的值解析为Golang的时间对象 time.Time |
| AllowNativePasswords | Bool<br>是否允许使用MySQL native password authentication method |
| Driver               | String<br>数据库驱动类型<br>当前支持："mysql"                |
| MaxOpenConns         | Int<br>最大活跃连接数                                        |
| MaxIdleConns         | Int<br>最大空闲连接数                                        |
| ConnMaxIdleTimeMs    | Int<br>连接最大空闲时间，单位为毫秒                          |
| ConnMaxLifetimeMs    | Int<br>连接最大生命期，单位为毫秒                            |

示例：

```
# Database Config
# see https://github.com/go-sql-driver/mysql/blob/master/dsn.go#L37
[Databases.bfe_db]
DBName              = "open_bfe"
Addr                = "127.0.0.1:3306"
Net                 = "tcp"
User                = "{user}"
Passwd              = "{password}"
MultiStatements     = true
MaxAllowedPacket    = 67108864
ParseTime           = true
AllowNativePasswords= true

Driver              = "mysql"
MaxOpenConns        = 100
MaxIdleConns        = 100
ConnMaxIdleTimeMs   = 50000
ConnMaxLifetimeMs   = 50000
```

### Dependence Config

Dependence Config 指定API Server的部分依赖文件的路径。

| 配置项      | 描述                                                         |
| ----------- | ------------------------------------------------------------ |
| NavTreeFile | String<br/>导航栏相关的配置文件路径                          |
| I18nDir     | String<br/>国际化相关的配置文件路径                          |
| UIIcon      | String<br/>自定义Dashboard上Icon的文件路径<br>Dashboard页面默认使用BFE Icon。需要时可以在此处指定自定义的Icon文件路径 |
| UILogo      | String<br>自定义Dashboard上Logo的文件路径<br>Dashboard页面默认使用BFE Logo。需要时可以在此处指定自定义的Logo文件路径 |

注：文件路径中，可以用 ${conf_dir} 表示配置文件根路径。

示例：

```
# Dependence Config
[Depends]
# NavTreeFile path
NavTreeFile = "${conf_dir}/nav_tree.toml"
# i18n conf dir path
I18nDir     = "${conf_dir}/i18n"
# dashboard icon
UIIcon      = "https://raw.githubusercontent.com/bfenetworks/bfe/develop/docs/images/logo/icon/color/bfe-icon-color.svg"
# dashboard logo
UILogo      = "https://raw.githubusercontent.com/bfenetworks/bfe/develop/docs/images/logo/horizontal/color/bfe-horizontal-color.png"
```

### Runtime Config

运行时配置。

| 配置项            | 描述                                                         |
| ----------------- | ------------------------------------------------------------ |
| SkipTokenValidate | Bool<br>是否跳过Token验证<br>建议设为False<br>若设为True，可以使用"Skip {role_name}"作为authorization header来调用API，例如：Headers[Authorization] = "Skip System"<br> |
| RecordSQL         | Bool<br>是否保存数据库操作日志                               |
| SessionExpireDay  | Int<br>会话过期时间，单位为天                                |
| StaticFilePath    | String<br>静态文件路径。对API请求进行动态路由失败时，若该路径下有静态文件，则返回静态文件 |
| Debug             | Bool<br>是否在API的响应中包含Debug信息                       |

示例：

```
# Runtime Config
[RunTime]
# you can use "Skip {role_name} as authorization header to access api server if open this optional
# eg: Headers[Authorization] = "Skip System"
# don't open it on production environment
SkipTokenValidate   = false
# sql will be record to log file when this option be opend
RecordSQL           = false
# how long use must login again
SessionExpireDay    = 10
# static file path, when dynamic router not be matched, static file will be return if found
StaticFilePath      = "./static"
# debug info will be add to response when this option be opend
Debug               = false

```

## nav_tree.toml 

该文件用来控制Dashboard的导航栏。

如果想在导航栏中关闭某些模块，在该文件中将模块对应的内容注释掉即可。

## i18n/ 文件夹

该文件夹下保存国际化相关的配置，主要是错误提示信息等内容。

当需要返回错误时，会根据请求的 `Accept-Language` 来找到 i18n/{lang}.toml 的语言包返回。

默认语言是英文。

