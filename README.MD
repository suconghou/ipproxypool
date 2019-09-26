

## 进程持续100个代理

提供api,也提供接口代理

基于 channel



http://tonybai.com/2015/11/17/tcp-programming-in-golang/


内存缓存最优数据


1. api 提交代理

代理类型
http
https
socket5

2. 显示代理接口


3. 并发抓取数据

header配置

失败重试



### 数据预处理

​/api/fetch/(json|html|text|任意选择器)/(任意属性|html|text)



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


## 环境变量

`PROXY_LISTEN` 配置此环境变量开启代理服务器

示例: `PROXY_LISTEN=0.0.0.0:6677 ./main`

将开启代理服务器监听端口6677,提供http/https/socks5代理
