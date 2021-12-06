本文讲解如何配置 API Server的配置，让 API Server可以正常运行。

# 配置文件说明
所有的配置文件都在 conf/ 文件夹下，包括：
- api_server.toml: API Server 模块启动必须的参数
- i18n/: 国际化相关的配置，主要是错误信息展示等内容
- nav_tree.toml: 导航栏相关配置

进行逐个讲解。

# api_server.toml 

```
# ---------------------------------
# API Server Config
[Server]
# server port
ServerPort          = 8183
# server graceful exit timeout
GracefulTimeOutMs   = 5000
# monitor port, don't start monitor server if less than 0
MonitorPort         = 8284


# ---------------------------------
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

# ---------------------------------
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


# ---------------------------------
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


# ---------------------------------
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

对于文件配置来讲，可能不清楚当前路径在哪，项目提供了环境变量来帮助配置文件路径：
- ${conf_dir}: 配置文件夹的根路径，是 启动命令的 -c 参数，默认为启动路径的 `./conf` 文件夹

你可以直接配置 `${conf_dir}/conf_file`，它将自动解析为一个绝对路径。

# i18n/ 文件夹
当需要返回错误时，会根据请求的 `Accept-Language` 来找到 i18n/{lang}.toml 的语言包返回。

默认语言是 英文。

# nav_tree.toml 
本配置用来控制前端的导航栏。如果你想关闭某些模块，修改该文件即可。