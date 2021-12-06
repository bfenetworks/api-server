API Server 是 BFE 控制面核心模块，完成配置的录入、存储和导出。

# 控制面组件
![架构](/docs/zh_cn/assert/deploy_architecture.png)

图1：控制面组件

BFE控制面包含如下组件：
- API Server: 对外提供Open API接口，完成BFE配置的变更、存储和下发
- Dashboard: 管理控制台，用于BFE集群的可视化管理
    - 仓库地址在 [bfenetwork/dashboard](https://github.com/bfenetworks/dashboard)
- Conf Agent: 配置加载组件，完成最新配置的获取和 BFE 热加载的触发
    - 仓库地址在 [bfenetwork/conf-agent](https://github.com/bfenetworks/conf-agent)


# 快速开始
通过查看 [部署说明](/docs/zh_cn/deploy.md) 快速运行 API Server。

# 二次开发
API Sever 提供 OpenAPI 供第三方程序和 API Server 集成，接口定义详见 [API 文档](docs/zh_cn/open_api/SUMMARY.md)。


# 相关模块
- [BFE数据面：负载均衡器](https://github.com/bfenetworks/bfe)
- [BFE控制面：控制台](https://github.com/bfenetworks/dashboard)
- [BFE控制面：Conf Agent](https://github.com/bfenetworks/conf-agent)


## 关于BFE
- 官网：https://www.bfe-networks.net
- 书籍：[《深入理解BFE》](https://github.com/baidu/bfe-book) ：介绍网络接入的相关技术原理，说明BFE的设计思想，以及如何基于BFE搭建现代化的网络接入平台。现已开放全文阅读。
