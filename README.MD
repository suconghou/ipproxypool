

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

  `/api/fetch/json/map`
  
2. **返回html全文拼接**

   `/api/fetch/html/html`

3. **返回text全文拼接**

   `/api/fetch/text/text`

4. **返回指定选择器属性值**

    * 返回页面所有a标签地址
      * `/api/fetch/a/href`
    * 返回页面所有img标签src地址
      * `/api/fetch/img/src`
    * 返回页面.item元素的data-id属性和text内容
      * `/api/fetch/.item/data-id|text`

