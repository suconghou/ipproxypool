## HTML抓取与下载任务

###  http/https 流量中继


直接将URL访问拼接到此域名后,将会通过此服务中转URL

默认不转发Referer,以跳过防盗链

默认超时时间600秒


###  静态资源服务器

参数`-d`指定一个目录,开启此服务

请求前缀`/public/` 将会映射到此目录

除此之外,当前工作目录下的`static`文件夹默认对外开放

### 资源下载进程组


API提交,提供批量下载任务

具体接口见 `/api/stask`

文件存储地址,由接口指定,请确保目标目录为非系统关键目录并且可写

`/api/task/add`
```
[{
	"url":"https://nodejs.org/dist/v16.17.0/node-v16.17.0.pkg",
	"mode":3
}]
```



### HTML解析


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

