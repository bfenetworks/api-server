
# 转发规则 

## 1 转发规则的模型
BFE的转发规则分为两张表: 
- 基础转发表
- 高级转发表

转发规则的匹配顺序:
- 先匹配基础转发表 
- 再匹配高级转发表

说明:
- 如果基础转发表为空，那么直接进入高级转发表匹配
- 转发规则中使用到的集群，必须是已经就绪的状态。否则 API 服务器拒绝本次提交

### 基础转发表
基础转发表
- 指定域名和PATH，确定转发的目标集群
- 基础表可以满足简单场景的转发规则，转发规则匹配算法是O(LogN) 
- 基础表由于不是线性匹配，所以表中规则的顺序与匹配并无直接关系

示例:

| 条件 | 目标集群 | 说明 | 
| - | - | - |
|域名: www.abc.com <br> Path: 无 | Cluster1 | 请求域名是 www.abc.com，则转发到 Cluster1|
|域名: www.xyz.com <br> Path: /path1 | Cluster2 | 请求域名是 www.xyz.com，并且 PATH 是 /path1，则转发到 Cluster2|
|域 名 : www.test1.com, www.test2.com <br> Path: /path2  | Cluster3 | 请求域名是 www.test1.com 或 www.test2， 并且 PATH 是/path2，则转发到 Cluster3|
|域名: www.xyz.com <br> Path: /path1,/path2  | Cluster4 | 请求域名是 www.xyz.com，并且 PATH 是 /path1 或/path3，则转发到 Cluster4 |
|域名: a.com,b.com <br> Path: /path1,/path2  | Cluster5 | 请求域名是 a.com 或 b.com，并且 PATH 是 /path1 或/path2，则转发到 Cluster5 |
|域名: www.xyz.com <br> Path: /path1 | GO_TO_ADVANCED_RULES | 请求域名是 a.com 或 b.com，并且 PATH 是 /path1 或/path2，则转发到 Cluster5 <br> 请求域名是 www.xyz.com，并且 PATH 是 /path1，则直接跳转到高级转发表继续匹配 |

注:GO_TO_ADVANCED_RULES 表示如果条件匹配，那么还需要进入高级转发规则表，做进一 步的规则匹配，找到最合适的集群，下文还有详细的例子。

- 域名支持泛域名匹配
    - 例如:*.abc.com，可以匹配a.abc.com,b.abc.com等;但不支持a.b.abc.com
- PATH支持前缀匹配，例如:/path1*，可以匹配以下几种形式:
    - /path1
    - /path1/abc 
    - /path1/a/b/c


### 高级转发表
高级转发表
- 高级转发表是一个有序的列表，规则的先后顺序，与转发引擎执行的顺序相同
    - 规则匹配的时间复杂度是O(N)
- 表中每一条规则，都是用Condition元语编写的条件表达式，以及集群
- 高级转发表中，最后一条规则，必须是:default_t()，它表示缺省的转发规则
   - 当其他所有规则都不匹配时，无条件匹配default_t()的规则 

| 条件 | 目标集群 | 说明 | 
| - | - | - |
| req_host_in("www.xyz.com") &&req_path_in(“/path1”) &&req_cookie_value_in("key1", "value1", false) | Cluster1 | 请求域名是 www.xyz.com，请求 path 是 /path1 ， 并 且 带 有 cookie: key1=value1，则转发到 Cluster1 |
 | req_host_in("www.xyz.com") &&req_path_in(“/path1”) | Cluster2 | 请求域名是 www.xyz.com，请求 path 是 /path1，则转发到 Cluster2 |
 | default_t() | Cluster3 | 其他情况下，一律进入 Cluster3| 


### 组合使用
- 基础表可以处理大多数简单的转发规则需求
- 高级表可以非常自由的指定流量特征，实现更复杂的转发


典型场景:灰度发布
- 常规流量使用域名:domain1访问，进入当前服务集群clusterA
- 测试流量使用域名:domain1 加上一个 cookie 值访问，进入一个测试集群clusterB

以上需求的配置方法:
- 基础表中配置规则:
    - [条件:domain1，目标集群:GO_TO_ADVANCED_RULES]
- 高级表中配置规则:
    - [条件:req_host_in("www.xyz.com") && req_cookie_value_in("key1", "value1", false)，目标集群:clusterB]
    - [条件:req_host_in("www.xyz.com"), false)，目标集群:clusterA] 
    
配置解释:
- 基础规则表中，匹配domain1，GO_TO_ADVANCED_RULES会将流程直接转到高级转发 规则表，进一步匹配
- 高级转发规则表中，带有cookie匹配的规则在前面，所以如果是测试流量，会被命中， 进入测试集群 clusterB
- 如果请求是常规流量，不带有 cookie，则会匹配下一条转发规则，进入常规集群 clusterA


## 2 更新转发规则

### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 端点| /products/{product_name}/routes ||
| method | PATCH ||
| 含义 | 整体更新转发列表 |- |

### 输入参数

#### URI 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - | 
| product_name | string | 产品线名字 | Y | - |

#### Body参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - | 
| basic_forward_rules | []object | 基础规则列表 | N |  |
| basic_forward_rules[].host_names | []string | 域名列表 | N |  |
| basic_forward_rules[].paths | []string | 路径列表 | N |  |
| basic_forward_rules[].cluster_name | string | 目标集群 | Y | 如果是跳转到高级规则，使用 "GO_TO_ADVANCED_RULES" 关键字 |
| basic_forward_rules[].description | string | 描述 | N |  |
| forward_rules | []object | 高级规则列表 | N | 如果存在，最后一条规则的表达式必须是 default_t() |
| forward_rules[].name | string | 名字 | N |  |
| forward_rules[].description | string | 描述 | N |  |
| forward_rules[].expression | string | 表达式 | Y |  |
| forward_rules[].cluster_name | string | 目标集群 | Y |- |


#### Body 请求示例
```
{
    "basic_forward_rules": [
        {
            "host_names": ["a.com"],
            "paths": [
                "/aaa",
                "/abc"
            ],
            "cluster_name": "GO_TO_ADVANCED_RULES",
            "description": ""
        }
    ],
    "forward_rules": [
        {
            "name": "rule1",
            "description": "",
            "expression": "req_host_in(\"b.com\")",
            "cluster_name": "Cluster1"
        },
        {
            "name": "default",
            "description": "",
            "expression": "default_t()",
            "cluster_name": "Cluster2"
        }
    ]
}
```
   

### 返回数据(Data内容)	

#### 返回数据示例
```
{
    "basic_forward_rules": [
        {
            "host_names": ["a.com"],
            "paths": [
                "/aaa",
                "/abc"
            ],
            "cluster_name": "GO_TO_ADVANCED_RULES",
            "description": ""
        }
    ],
    "forward_rules": [
        {
            "name": "rule1",
            "description": "",
            "expression": "req_host_in(\"b.com\")",
            "cluster_name": "Cluster1"
        },
        {
            "name": "default",
            "description": "",
            "expression": "default_t()",
            "cluster_name": "Cluster2"
        }
    ]
}
```
 
## 3 获取转发规则列表

### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 端点  | /products/{product_name}/routes | |
| method | GET | |
| 含义 | 获取产品线的转发规则列表 ||

### 输入参数

#### URI 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - | 
| product_name | string | 产品线名字 | Y | - |
  

### 返回数据(Data内容)	
同更新转发规则