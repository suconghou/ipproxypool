## 模块化的爬虫代理组件

### 模块一: 代理IP池

参数 `-proxyfetch` 开启IP Proxy 爬虫

将会收集可用的代理IP,存储在内存

API `/api/proxy/info` 可获取此相关信息

每隔一分钟dump为json存储,以备后续启动时直接载入

### 模块二: http/https/socks代理

参数 `-proxylisten` 开启此功能

将会监听指定端口,提供代理

示例: `-proxylisten 0.0.0.0:6677`

将开启代理服务器监听端口6677,提供http/https/socks5代理


### 模块三: http/https 流量中继

此模块默认开启

直接将URL访问拼接到此域名后,将会通过此服务中转URL

默认不转发Referer,以跳过防盗链

默认超时时间600秒


### 模块四: 静态资源服务器

参数`-d`指定一个目录,开启此服务

请求前缀`/public/` 将会映射到此目录

除此之外,当前工作目录下的`static`文件夹默认对外开放

### 模块五: 资源下载进程组

此模块默认开启

API提交,提供批量下载任务

具体接口见 `/api/stask`

文件存储地址,由接口指定,请确保目标目录为非系统关键目录并且可写



### 模块六: 爬虫解析

此模块默认开启

提供网页批量抓取与分析

具体接口见 `/api/fetch`


### 数据预处理

`​/api/fetch/(json|html|text)`


1. **返回原始JSON**

  `/api/fetch/json`
  
2. **返回html全文拼接**

   `/api/fetch/html`

3. **返回text全文拼接**

   `/api/fetch/text`

4. **返回指定选择器属性值**
    `/api/fetch/query`
    

```
{
	"urls":
	[
		{
			"url":"https://v2ex.com/?tab=hot"
			
		},
		{
			"url":"https://v2ex.com/"
		}
	],
	"query":
	{
		"#TopicsHot .item_hot_topic_title a":{
			"attrs":{
				"href":"link"
			},
			"methods":{
				"text":"title"
			}
		},
		".box .item":{
			"query":{
				".item_title a":{
					"attrs":{"href":"href"},
					"methods":{"text":"title"}
				},
				".count_livid":{
					"methods":{
						"text":"count"
					}
				},
				".node":{
					"attrs":{
						"href":"topic_link"
					},
					"methods":{
						"text":"topic"
					}
				}
				
			}
		}
	}
}
```


## 应用示例

