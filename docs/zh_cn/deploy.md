本文讲解如何部署 BFE 控制面组件。

# BFE控制面组件
BFE控制面包含如下组件：
- API Server: 对外提供Open API接口，完成BFE配置的变更、存储和下发
- Dashboard: 管理控制台，用于BFE集群的可视化管理
    - 仓库地址在 [bfenetwork/dashboard](https://github.com/bfenetworks/dashboard)
- Conf Agent: 配置加载组件，完成最新配置的获取和 BFE 热加载的触发
    - 仓库地址在 [bfenetwork/conf-agent](https://github.com/bfenetworks/conf-agent)


# 部署架构
![部署架构](./assert/deploy_architecture.png)

图1：部署架构

如图1所示：
- 配置变更：
    - 系统管理人员可以通过 BFE Dashboard 相关配置进行可视化管理
    - 也可以通过调用 BFE API Server 提供的 Open API进行管理
- 配置下发：
    - Conf Agent和BFE同机部署
    - Conf Agent和API Server通信，在发现有新的配置后，将新的配置读取到BFE所在的服务器，并触发BFE执行配置热加载
    - BFE被触发后，读取最新的配置并生效


# 部署步骤

部署步骤依次为：
1. API Server 部署
1. Dashboard 部署
1. Conf Agent 部署

## APIServer部署
1. 安装MySQL数据库：数据库版本5.6以上即可，具体安装过程本文不详细描述
1. 初始化数据库： 执行 `mysql -u{user} -p{password} < db_ddl.sql`
1. 获取API Server可执行程序
    - 方式一：通过源码编译：clone本仓库后进入项目根目录，执行 `make`，output文件夹包括了可执行文件和初始配置文件
    - 方式二：直接进入 [releases](https://github.com/bfenetworks/api-server/releases) 页面下载相应的编译产出
1. 修改初始配置文件，详见[配置文件说明](./config_param.md)
- 特别注意：绝大多数配置可以使用默认配置，最小修改集合为 **数据库用户名和密码**
1. 启动 API Server。执行`./api-server -c ./conf -sc api_server.toml -l ./log `。如果不需要指定启动参数，直接执行 `./api-server` 即可

## Dashboard部署
1. 获取 Dashboard 产出
    - 方式一：通过源码编译： clone [bfenetwork/dashboard](https://github.com/bfenetworks/dashboard) 仓库后进入项目根目录，执行 `sh build.sh`， output 文件夹就是静态配置文件
    - 方式二：直接进入 [dashboard/releases](https://github.com/bfenetworks/dashboard/releases) 页面下载相应的编译产出
1. 部署：将output文件夹的内容拷贝到 API Server 的 static 文件夹（默认在api-server可执行文件同级目录）中即可
1. 浏览器打开 http://host:{ServerPort} (ServerPort 为 API Server 部署时配置的端口号) 即可看到登录页面，测试账号和密码都是 `admin`。登陆后，请立刻修改您的admin的密码

可查看 [Dashboard 使用文档](https://github.com/bfenetworks/dashboard/blob/develop/docs/zh-cn/user-guide/SUMMARY.md) 了解BFE管理控制台的基本概念和使用流程。


## ConfAgent部署
Conf Agent和 BFE 转发引擎同机部署。

若未部署BFE转发引擎，可参考[BFE安装部署](https://www.bfe-networks.net/en_us/installation/install/)完成部署。

**若BFE已经上线运行，有历史配置数据，需要完善业务配置才能启动Conf Agent，本文不展开描述。**

1. 获取Conf Agent可执行程序
    - 方式一：通过源码编译：clone [bfenetwork/conf-agent](https://github.com/bfenetworks/conf-agent) 仓库后进入项目根目录，执行 `make`， output文件夹包括了可执行文件和初始配置文件
    - 方式二：直接进入 [releases](https://github.com/bfenetworks/conf-agent/releases) 页面下载相应的编译产出
1. [可选]微调部分当前不支持通过APIServer管理的BFE配置文件内容。详见 [可能需要手动维护的BFE配置文件](#keep)
1. 修改配置，详见 [配置文件说明](https://github.com/bfenetworks/conf-agent/blob/develop/docs/zh_cn/config.md)，使其能访问API Server导出最新的配置
1. 启动 Conf Agent。执行 `./conf-agent`

### <a id="keep">可能需要手动维护的BFE配置文件</a>

对于tls配置，现在APIServer不支持 配置 的文件有：
- client_ca
- client_crl
- session_ticket_key.data
- tls_rule_conf.data

能配置和下发的文件有：
- server_cert_conf.data  
 
在热加载tls配置时， tls_rule_conf.data 内容和 server_cert_conf.data 内容相关联。当两者不一致时就会出现关联检查失败而报错。
当前 默认的 tls_rule_conf.data 的配置依赖 `example.org` 证书配置 来指定 租户 `example_product` 的TLS协议配置。如果自行添加的证书中不存在 `example.org` 证书，并且继续使用默认的 tls_rule_conf.data 配置内容，ConfAgent 将在触发BFE热加载tls配置时报错。

当前这个问题的解决方案是自行维护 tls_rule_conf.data 文件（ConfAgent会直接使用修改后的文件内容），根据业务需求手动修改该文件的原始内容。

如果你不需要对租户进行自定义，建议修改方案为：

```
默认配置内容：
{
    "Version": "12",
    "DefaultNextProtos": ["http/1.1"],
    "Config": {
        "example_product": {
            "VipConf": [
                "10.199.4.14"
            ],
            "SniConf": ["example.org"],
            "CertName": "example.org",
            "NextProtos": [
                "h2;rate=100;isw=65535;mcs=200;level=0",
                "http/1.1"
            ],
            "Grade": "C",
            "ClientAuth": false,
            "ClientCAName": "example_ca"
        }
    }
}

修改后的简化配置：
{
    "Version": "12",
    "DefaultNextProtos": ["http/1.1"],
    "Config": {}
}

```
