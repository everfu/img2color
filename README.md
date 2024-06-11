<div align="center">
   <h1>Img2Color</h1>

使用 Go 语言编写的一个简单的 API，用于获取图片的主色，支持使用 Vercel KV 存储桶。

</div>

## 返回格式
```json
{
   "RGB": "#756576"
}
```

## 使用方法

1. 创建并登录 Vercel 账号。
2. 点击右侧按钮。 [![Deploy with Vercel](https://vercel.com/button)](https://vercel.com/new/clone?repository-url=https://github.com/everfu/img2color&env=ALLOWED_REFERERS,KV_ENABLE&project-name=img2color&repo-name=img2color)
3. 配置环境变量。
    * ALLOWED_REFERERS: 设置可请求的域名（示例: '*, *.efu.me'）。如果不设置，将允许所有域名请求。
    * KV_ENABLE: 是否启用 Vercel KV 存储（true/false）。
4. 部署项目。

## KV 存储
创建一个 KV 桶，链接刚刚的项目。

如果在刚开始的时候没有启用 KV 存储，可以在项目设置中启用（设置：`KV_ENABLE: true`）。

重新生成项目即可。

测试：`https://examp.com/api?img=https://example.com/image.jpg`

## 版权

[MIT](LICENSE) License © 2024-至今 [EverFu](https://github.com/everfu)