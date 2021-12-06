# 产品线

## 1 创建产品线

### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 含义	| 创建产品线 | |
| 端点 | /products | |
| 版本 | v1 |  |
| method | POST | - |


### 输入参数
#### Body 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - | 
| name | string | 产品线名 |  Y | 必须全局唯一 |
| description | string | 产品线描述 |  N | |
| mail_list | []string | 产品线团队的邮箱的列表 |  N | |
| phone_list | []string | 值班电话的列表 | N | |
| contact_person_list | []string | 接口人邮箱的列表 | N | - |

#### HTTP BODY中参数示例
```
{ 
    "name": "product-demo", 
    "description": "demo product", 
    "mail_list": ["op@bfenetwork.com"],
    "phone_list": ["13512341234", "13543214321"],
    "contact_person_list": ["manager@bfenetwork.com"]
}
```

### 返回数据(Data内容)
字段含义同输入参数。

#### 成功返回数据示例
```
{ 
    "name": "product-demo",
    "description": "demo product", 
    "mail_list": ["op@bfenetwork.com"],
    "phone_list": ["13512341234", "13543214321"],
    "contact_person_list": ["manager@bfenetwork.com"]
}
```


## 2 产品线列表

### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 含义	| 获取全体产品线信息列表 | |
| 端点	| /products | |
| 版本	| v1 | |
| 动作	| GET | - |


### 输入参数
#### Query 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - | 
| domain | string | 域名 |  N | 产品线下的域名 |
| cluster | string | 集群名 |  N | 产品线下的集群，当domain查询参数不存在时，按照cluster进行过滤 |

### 返回数据(Data内容)
产品线列表

#### 成功返回数据示例
```
[ 
    { 
        "name": "product-demo",
        "description": "demo product description", 
        "mail_list": [“op@bfenetwork.com”], 
        "phone_list":  [“13512341234”, “13543214321”],
        "contact_person_list": ["manager@bfenetwork.com"]
    },
    { 
        "name": "test-product", 
        "description": "test product description", 
        "mail_list": [“op@bfenetwork.com”], 
        "phone_list": ["13512341234"], 
        "contact_person_list": [manager2@bfenetwork.com]
    }
]
```

## 3 产品线详情

### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 含义	| 获取单个产品线信息 |  |
| 端点	| /products/{product_name} ||
| 版本	| v1 ||
| 动作	| GET | - |

### 输入参数

#### URI 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - |  
| product_name | string | 产品线名称 | Y | - |

### 返回数据(Data内容)
同创建接口



## 4 更新产品线

### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 含义	| 更新产品线信息 ||
| 端点	| /products/{product_name} ||
| 版本	| v1 ||
| 动作	| PATCH | - |

### 输入参数

#### URI 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - |  
| product_name | string | 产品线名称 | Y | - |

HTTP Body 参数：
```
{
    "description": "demo product"
    "mail_list": ["op@bfenetwork.com"], 
    "phone_list": ["13512341234"], 
    "contact_person_list": ["manager@bfenetwork.com"]
}
```

- 必须携带完整的产品线信息，不支持仅提交部分数据项
- 产品线名字无法更新

### 返回数据(Data内容)
同创建接口


## 5 删除产品线

### 基本信息

| 项目  | 值  | 说明 | 
| - | - | - |
| 含义	| 删除产品线 ||
| 端点	| /products/{product_name} ||
| 版本	| v1 ||
| 动作	| DELETE | - |

### 输入参数

#### URI 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - |  
| product_name | string | 产品线名称 | Y | - |

### 返回数据(Data内容)
同创建接口

