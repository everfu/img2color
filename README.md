## img2color
使用Goland简单实现图片获取主题色

## 使用方法
1. 创建并登陆 vercel
2. 点击按钮快速导入
[![Deploy with Vercel](https://vercel.com/button)](https://vercel.com/new/clone?repository-url=https://github.com/wleelw/img2color&env=ALLOWED_REFERERS,KV_ENABLE&project-name=img2color&repo-name=img2color)
3. 设置环境变量
    
    ALLOWED_REFERERS: 请求refer（填写格式：'*,*.wzsco.top'，注意使用英文半角逗号分隔）

    KV_ENABLE: 是否启用 KV 存储进行缓存，如不需要请填写 false

## 增加缓存功能

1. 新建一个 vercel KV
   ![](https://github.com/wleelw/img2color/assets/74389842/fa435feb-c311-4044-8954-876729185027)
2. 将 KV 与项目绑定
   ![](https://github.com/wleelw/img2color/assets/74389842/c288a9ca-9416-46c5-af8a-9210e5f826c5)
3. 如果刚开始将 KV_ENABLE 设置为 false，需要到项目的环境变量中改为true
   ![](https://github.com/wleelw/img2color/assets/74389842/6985efab-4eb5-48c7-ba96-4675fc6e7ad7)
4. 选择任意 deploy 重新构建
   ![](https://github.com/wleelw/img2color/assets/74389842/4649aba6-c455-444b-bb50-5d6e8f7982a7)
5. 测试
   ![](https://github.com/wleelw/img2color/assets/74389842/f2e95b9b-4756-419b-aad6-2caab03e12ac)

## 贡献
[@Efu](https://github.com/efuo/)

## 仓库统计
![Alt](https://repobeats.axiom.co/api/embed/9a7ae5077e31ed3b650612589a712c220da1dd18.svg "Repobeats analytics image")
