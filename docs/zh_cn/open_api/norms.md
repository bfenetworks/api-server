# 整体说明

## API规范

一个典型API包括如下部分：
- 接口名：一句话描述接口名
- 基本信息
- 输入参数
- 返回参数

### 基本信息

**URL格式说明**
- API遵守一般RESTful风格，API的URL格式：
    - http://api_server:port/open-api/{ver}/{endpoint}?{arg=value}
    - 例子：http://127.0.0.1:8086/open-api/v1/products
- API URL各部分说明
    - api_server: 服务器地址，一般是域名或者IP地址
    - port：API服务的端口号
    - ver：当前API的版本
    - endpoint：REST风格的资源路径
    - arg：参数名
    - value：参数值

若无特殊说明，后续文档的具体API只描述Endopoint.

**method说明**

若无特殊说明，method 遵循如下约定：
- GET：读取
- POST: 创建
- PATCH：更新
- DELETE：删除

举例：

| 项目  | 值  | 说明 | 
| - | - | - |
| 含义	| 创建产品线 | |
| 端点 | /products | |
| 版本 | v1 |  |
| method | POST | - |

### 请求参数

请求参数分为：

- URI参数
- Query参数
- Body内容

举例：

| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - | 
| phone_list | []string | 联系人列表 |  Y | |
| HealthCheck | object | 健康检查配置 |  Y | |
| HealthCheck.Interval | int64 | 健康检查间隔 | Y | - |

HTTP BODY中参数示例
```
{ 
    "name": "bfe", 
    "description": "demo product", 
    "mail_list": ["op@bfenetwork.com"],
    "phone_list": ["13512341234", "13543214321"],
    "contact_person_list": ["manager@bfenetwork.com"]
}
```


### 返回数据

所有API的返回值格式为：

```
{
	"ErrNum": number,
	"ErrMsg": "string message",
	"Data": json_object
}
```

- ErrNum: 返回码
    - 200: 调用成功时
    - 调用失败时，
        - 402：没有调用权限造成的失败
        - 422：参数不合法造成的失败
        - 510：集群/分流规则创建时实例池未ready
        - 404：查询/修改/删除不存在的对象时
        - 555：创建重复对象时
        - 500：其他业务逻辑错误，一律返回500
- Data: 返回的数据结构
    - 调用成功时，返回json格式的数据
    - 调用失败时，返回null
- ErrMsg: 文本消息
    - 调用成功时，ErrMsg是success或空串
    - 调用失败时，ErrMsg是相关的错误信息


举例：
```
{
	"ErrNum": 200,
	"ErrMsg": "Succ",
	"Data": {
		"name": "bfe",
		"description": "demo product",
		"mail_list": ["op@bfenetwork.com"],
		"phone_list": ["13512341234", "13543214321"],
		"contact_person_list": ["manager@bfenetwork.com"]
	}
}
```

说明：API文档中中API的返回结果，仅给出Data部分。


## 鉴权机制
- API使用Token机制鉴权
- 访问时在HTTP Authorization HEADER中加入SessionKey/Token
- 鉴权详细机制见 [用户和鉴权](global/auth.md)
- Session Key的使用示例：

```
curl http://127.1:8086/open-api/v1/products/demo/clusters -H "Authorization: Session gc0JnZJpkMBmqJf1dbcV" 
```
