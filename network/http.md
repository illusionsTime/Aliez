#### 无状态协议

HTTP自身不对请求和响应之间的通信状态进行保存，协议对于发送过的请求或响应都不做持久化处理。

#### 协议格式
### URI

HTTP使用统一资源标识符（URI）来传输数据和建立连接。URL（统一资源定位符）是一种特殊种类的URI，包含了用于查找的资源的足够的信息，我们一般常用的就是URL，而一个完整的URL包含下面几部分：
1. 协议部分
该URL的协议部分为http:，表示网页用的是HTTP协议，后面的//为分隔符
2. 域名部分
域名是www.fishbay.cn，发送请求时，需要向DNS服务器解析IP。如果为了优化请求，可以直接用IP作为域名部分使用
3. 端口部分
域名后面的80表示端口，和域名之间用:分隔，端口不是一个URL的必须的部分。如果端口是80，也可以省略不写
4. 虚拟目录部分
从域名的第一个/开始到最后一个/为止，是虚拟目录的部分。其中，虚拟目录也不是URL必须的部分，本例中的虚拟目录是/mix/
5. 文件名部分
从域名最后一个/开始到?为止，是文件名部分；如果没有?，则是从域名最后一个/开始到#为止，是文件名部分；如果没有?和#，那么就从域名的最后一个/从开始到结束，都是文件名部分。本例中的文件名是76.html，文件名也不是一个URL的必须部分，如果没有文件名，则使用默认文件名
6. 锚部分
从#开始到最后，都是锚部分。本部分的锚部分是first，锚也不是一个URL必须的部分
7. 参数部分
从?开始到#为止之间的部分是参数部分，又称为搜索部分、查询部分。本例中的参数是name=kelvin&password=123456，如果有多个参数，各个参数之间用&作为分隔符。

### req resp
HTTP报文分为HTTP请求报文和响应报文，请求报文由请求行（请求方法，请求资源的URL和HTTP的版本）、首部行和实体（通常不用）组成。响应报文由状态行（状态码，短语和HTTP版本）、首部行和实体（有些不用）组成。

#### 流程

①先检查输入的URL是否合法，然后查询浏览器的缓存，如果有则直接显示。

②通过DNS域名解析服务解析IP地址，先从浏览器缓存查询、然后是操作系统和hosts文件的缓存，如果没有查询本地服务器的缓存。

③通过TCP的三次握手机制建立连接，建立连接后向服务器发送HTTP请求，请求数据包。

④服务器收到浏览器的请求后，进行处理并响应。

⑤浏览器收到服务器数据后，如果可以就存入缓存。

⑥浏览器发送请求内嵌在HTML中的资源，例如css、js、图片和视频等，如果是未知类型会弹出对话框。

⑦浏览器渲染页面并呈现给用户。