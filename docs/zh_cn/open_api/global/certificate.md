# 证书

## 1 创建证书

### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 含义	| 创建证书 | |
| 端点 | /certificates | |
| 版本 | v1 |  |
| method | POST | - |


### 输入参数
#### Body 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - |
| cert_name | string | 证书名 |  Y | 必须唯一 |
| description | string | 证书描述 |  Y | |
| is_default | bool | 是否是默认证书 |  Y | 必须有且只有一个默认证书  |
| cert_file_name | string | 主证书文件名 | Y | |
| cert_file_content | string | 主证书文件内容 | Y | |
| key_file_name | string | 主证书Key名 | Y | |
| key_file_content | string | 主证书Key文件内容 | Y | |
| expired_date | string | 主证书过期时间 | Y | - |

#### HTTP BODY中参数示例
```
{
	"cert_name": "cert_demo",
	"description":"abc",
	"is_default": true,
	
	"cert_file_name": "demo_cert_file_name",
	"cert_file_content":"-----BEGIN ...-----END CERTIFICATE-----",
	"key_file_name": "demo_key_file_name",
	"key_file_content":"-----BEGIN RSA PRIVATE KEY-----...-----END RSA PRIVATE KEY-----",
	"expired_date":"2021-08-23 16:02:31"
}
```

### 返回数据(Data内容)
cert_file_content、key_file_content不返回，其他字段含义同输入参数。

#### 成功返回数据示例
```
{
	"cert_name": "cert_demo",
	"description":"abc",
	"is_default": true,
	
	"cert_file_name": "demo_cert_file_name",
	"key_file_name": "demo_key_file_name",
	"expired_date":"2021-08-23 16:02:31"
}
```


## 2 证书列表

### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 含义	| 获取全体证书信息列表 | |
| 端点	| /certificates | |
| 版本	| v1 | |
| 动作	| GET | - |


### 输入参数

### 返回数据(Data内容)
同创建接口

#### 成功返回数据示例
```
[ 
    {
    	"cert_name": "cert_demo",
    	"description":"abc",
    	"is_default": true,
    	"cert_file_name": "demo_cert_file_name",
    	"key_file_name": "demo_key_file_name",
    	"expired_date":"2021-08-23 16:02:31"
    }
]
```

## 3 证书详情

### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 含义	| 获取单个证书信息 |  |
| 端点	| /certificates/{cert_name} ||
| 版本	| v1 ||
| 动作	| GET | - |

### 输入参数

#### URI 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - |  
| cert_name | string | 证书名称 | Y | - |

### 返回数据(Data内容)
同创建接口


## 4 更新证书为默认证书

### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 含义	| 更新证书信息 ||
| 端点	| /certificates/{cert_name}/default ||
| 版本	| v1 ||
| 动作	| PATCH | - |

### 输入参数

#### URI 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - |  
| cert_name | string | 证书名称 | Y | - |


更新为默认证书时，旧的默认证书自动变为非默认证书

### 返回数据(Data内容)
同创建接口


## 5  删除证书

### 基本信息

| 项目  | 值  | 说明 | 
| - | - | - |
| 含义	| 删除证书 ||
| 端点	| /certificates/{cert_name} ||
| 版本	| v1 ||
| 动作	| DELETE | - |

### 输入参数

#### URI 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - |  
| cert_name | string | 证书名称 | Y | - |


- 默认证书不能被删除，全局必须有一个默认证书

### 返回数据(Data内容)
无
