本文档描述 认证和授权 的相关接口。

认证包括两大类：
- 普通用户，使用用户名和密码登录
- token

API请求时需要将 会话密钥 放入Header 的 "Authorization"头 中，当前有如下形式：
- Session Key：形式为 "Authorization: Session {session_key}", 比如： "Authorization: Session sxddxefda8"
- Token：形式为 "Authorization: Token {token}", 比如： "Authorization: Token daalkfjdkx"

当前，所有的资源被划分为如下三个 scope：
- Product：产品线资源，比如产品线的转发规则
- Support：导出类资源， 用于BFE数据面模块从API-Server导出所需要的配置
- System： 全部的权限，包括全局配置（比如 BFECluster）、产品线资源和导出类资源

对于普通用户和token都会设定可访问资源的scope，只能访问 scope 内资源
- 如果设定的scope为Product，还需要进一步校验是否具有某个产品线的权限


# 1 用户

## 1.1 创建用户

### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 含义	| 创建用户 | |
| 端点 | /auth/users | |
| 版本 | v1 |  |
| method | POST | - |


### 输入参数
#### Body 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - |
| user_name | string | 用户名 |  Y |  |
| password | string | 用户密码 |  Y | 必填 |
| is_admin | bool | 是否是系统管理员 |  Y | 如果是，就是有 System 的权限，不然就是 Product的权限  |

#### HTTP BODY中参数示例
```
{
	"user_name": "user_demo",
	"password": "password@baidu.com",
	"is_admin": true
}
```


### 返回数据(Data内容)
无

## 1.2 删除用户

### 基本信息

| 项目  | 值  | 说明 | 
| - | - | - |
| 含义	| 删除用户 ||
| 端点 | /auth/users/{user_name} | |
| 版本	| v1 ||
| 动作	| DELETE | - |

### 输入参数
#### URL 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - |
| user_name | string | 待删除的用户名 |  Y | - |

### 返回数据(Data内容)
无


## 1.3 重置用户密码

### 基本信息

| 项目  | 值  | 说明 | 
| - | - | - |
| 含义	| 重置用户密码 ||
| 端点 | /auth/users/{user_name}/passwd | |
| 版本	| v1 ||
| 动作	| PATCH | - |

### 输入参数
#### URL 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - |
| user_name | string | 待修改密码的用户名 |  Y | - |

#### Body 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - |
| old_password | string | 旧的用户密码 |  N | 当被修改的用户为当前登录用户，需要填入旧密码 |
| password | string | 用户新密码 |  Y | - |

#### HTTP BODY中参数示例
```
{
	"old_password": "manager2123@$"
	"password": "manager2123@$"
}
```

### 返回数据(Data内容)
无


## 1.4 获取用户列表

### 基本信息

| 项目  | 值  | 说明 | 
| - | - | - |
| 含义	| 查看用户列表 ||
| 端点 | /auth/users | |
| 版本	| v1 ||
| 动作	| GET | - |

### 输入参数
无

### 返回数据(Data内容)
数组，每个元素为一个用户

#### 成功返回数据示例

```
[
    {
        "user_name": "user_demo1",
        "is_admin": true
    },
    {
        "user_name": "user_demo",
        "is_admin": false
    }
]
```


## 1.5 设置用户是否具有管理员权限

### 基本信息

| 项目  | 值  | 说明 | 
| - | - | - |
| 含义	| 设置用户是否有管理员权限 ||
| 端点 | /auth/users/{user_name}/is_admin | |
| 版本	| v1 ||
| 动作	| PATCH | - |

### 输入参数
#### URL 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - |
| user_name | string | 待修改权限的用户的用户名 |  Y | - |

#### Body 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - |
| is_admin | bool | 是否为系统管理员 |  Y | 系统管理员有System(所有)的权限 |

#### HTTP BODY中参数示例
```
{
	"is_admin": true
}
```

### 返回数据(Data内容)
无


## 1.6 为用户增加某个产品线的授权

### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 含义	| 为用户增加某个产品线的授权 | |
| 端点 | /auth/users/{user_name}/products/{product_name} | |
| 版本 | v1 |  |
| method | POST | - |

### 输入参数
#### URI 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - |
| product_name | string | 产品线名 |  Y | - |
| user_name | string | 用户名 |  Y | - |

### 返回数据(Data内容)
无

## 1.7 对用户取消某个产品线的授权

### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 含义	| 对用户取消某个产品线的授权 | |
| 端点 | /auth/users/{user_name}/products/{product_name} | |
| 版本 | v1 |  |
| method | DELETE | - |

### 输入参数
#### URI 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - |
| product_name | string | 产品线名 |  Y | - |
| user_name | string | 用户名 |  Y | - |

### 返回数据(Data内容)
无


## 1.8 获取对指定产品线有权限的用户列表

### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 含义	| 获取对指定产品线有权限的用户列表 | |
| 端点 | /auth/users/actions/search-by-product/{product_name} | |
| 版本 | v1 |  |
| 动作	| GET | - |

### 输入参数
#### URL 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - |
| product_name | string | 产品线名 |  Y | - |

### 返回数据(Data内容)

数组，每个元素为一个用户

#### 成功返回数据示例

```
[
    {
        "user_name": "user_demo",
        "is_admin": false
    }
]
```


# 2 session key

## 2.1 使用账号名密码创建session key

### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 含义	| 使用账号密码得到session key（可用来登录） | |
| 端点 | /auth/session-keys | |
| 版本 | v1 |  |
| method | POST | - |


### 输入参数
#### Body 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - |
| user_name | string | 用户名 |  Y |  |
| password | string | 用户密码 |  Y | - |

#### HTTP BODY中参数示例
```
{
	"user_name": "manager2",
	"password": "manager2123@$"
}
```

### 返回数据(Data内容)
| 参数名 | 类型 |参数含义 | 补充描述 |
| - | -  | - | - | 
| session_key | string | 会话密钥 | 在后续请求中需要在Header中带上该值，格式为 "Authorization: Session iMQW0z5ZwK_6FnPPT7Xj" |
| user_name | string | 用户名 |  |
| is_admin | bool | 是否是系统管理员 |   如果是，就是有 System 的权限  |

#### 成功返回数据示例
```
{
    "user_name": "user_demo",
    "session_key": "iMQW0z5ZwK_6FnPPT7Xj",
    "is_admin": false
}
```

## 2.2 删除 session key

### 基本信息

| 项目  | 值  | 说明 | 
| - | - | - |
| 含义	| 删除 session key ||
| 端点 | /auth/session-keys/{session_key} | |
| 版本	| v1 ||
| 动作	| DELETE | - |

### 输入参数
#### URL 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - |
| session_key | string | 待删除的session key |  Y | - |

### 返回数据(Data内容)
无

# 3 Token


## 3.1 创建Token

### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 含义	| 创建Token（同时完成产品线绑定） | |
| 端点 | /auth/tokens | |
| 版本 | v1 |  |
| method | POST | - |

### 输入参数
#### Body 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - |
| name | string | token名字 |  Y |  name必须全局唯一 |
| scope | string | scope |  Y | 只能指定一个scope  |
| product_name | string | 产品线名 |  Y | 如果scope 为 Product，必须且只能绑定一个产品线 |

#### HTTP BODY中参数示例
```
{
	"name": "token_demo",
	"scope": "Product",
	"product_name": "product_demo"
}
```

### 返回数据(Data内容)
| 参数名 | 类型 |参数含义 | 补充描述 |
| - | -  | - | - | 
| token | string |  | 在后续请求中需要在Header中带上该值，格式为 "Authorization: Token Px2szn6R1HQo-WRSIJyt" |

#### 成功返回数据示例
```
{
    "token": "Px2szn6R1HQo-WRSIJyt"
}
```


## 3.2 删除Token

### 基本信息

| 项目  | 值  | 说明 | 
| - | - | - |
| 含义	| 删除token ||
| 端点 | /auth/tokens/{token_name} | |
| 版本	| v1 ||
| 动作	| DELETE | - |

### 输入参数
#### URL 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - |
| token_name | string | 待删除的token name |  Y | - |

### 返回数据(Data内容)
无


## 3.3 查看Token详情

### 基本信息

| 项目  | 值  | 说明 | 
| - | - | - |
| 含义	| 查看Token详情 ||
| 端点 | /auth/tokens/{token_name} | |
| 版本	| v1 ||
| 动作	| GET | - |

### 输入参数
#### URL 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - |
| token_name | string | token name |  Y | - |

### 返回数据(Data内容)
| 参数名 | 类型 |参数含义 |  补充描述 |
| - | -  | - | - | - |
| name | string | token名字 | |
| product_name | string | 产品线名 |  |
| token | string | token的值 | |
| scope | string | scope | - |

#### 成功返回数据示例

```
{
    "name": "token_demo",
    "product_name": "product_demo",
    "token": "Xim4h3tR_Gp7o4h",
    "scope": "Product"
}
```

## 3.4 查看Token列表

### 基本信息

| 项目  | 值  | 说明 | 
| - | - | - |
| 含义	| 查看Token列表 ||
| 端点 | /auth/tokens | |
| 版本	| v1 ||
| 动作	| GET | - |

### 输入参数
无

### 返回数据(Data内容)
数组，每个元素为Token (详见“查看Token详情”)

#### 成功返回数据示例

```
[
    {
        "name": "token_demo",
        "product_name": "product_demo",
        "token": "Xim4h3tR_Gp7o4h",
        "scope": "Product"
    }
]
```

## 3.5 获取对指定产品线有权限的Token列表

### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 含义	| 获取对指定产品线有权限的Token列表 | |
| 端点 | /auth/tokens/actions/search-by-product/{product_name} | |
| 版本 | v1 |  |
| 动作	| GET | - |

### 输入参数
#### URL 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - |
| product_name | string | 产品线名 |  Y | - |


### 返回数据(Data内容)

数组，每个元素为一个token对象(详见“查看Token详情”)

#### 成功返回数据示例

```
[
    {
        "name": "token_demo",
        "token": "Xim4h3tR_Gp7o4h",
        "scope": "Proudct"
    }
]
```
