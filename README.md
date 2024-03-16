# Mikanani 番剧更新推送 Bot

![](https://img.shields.io/badge/license-MIT-blue)
![](https://img.shields.io/badge/GO-1.22-blue)
![](https://img.shields.io/badge/PRs-welcome-green)

> 因为 Mikanani 是一个简体中文网站，所以这个专案不会提供英文版和繁体中文版

这个 Bot 会自动获取你添加的番剧的更新剧集，并给你推送更新剧集的种子

## 可用命令

- `/add <番剧名称>` 把番剧添加到自动更新列表
- `/del <番剧名称>` 把番剧移出自动更新列表
    - 这个专案没有检查番剧是否完结的功能，所以请在番剧完结时手动移出自动更新列表
- `/get` 获取所有标记为未下载的项目
- `/help` 获取可用命令

## 演示

- 使用命令添加番剧

![image](https://github.com/ArsFy/mikanani-push-bot/assets/93700457/af566bfe-f860-4429-8982-9d44a8776fb8)

- 主动推送

![image](https://github.com/ArsFy/mikanani-push-bot/assets/93700457/f4057030-259d-407e-abdc-22a88e4074b2)

- 获取所有标记为未下载的项目

![image](https://github.com/ArsFy/mikanani-push-bot/assets/93700457/1d08d713-f1ca-4f6c-b403-17871a62ef3b)

## 部署
### 1. 下载
你可以在 [Releases](https://github.com/ArsFy/mikanani-push-bot/releases) 中下载预设的压缩档，或者使用 Golang 编译二进制可执行档

下载完成后的目录结构应该如下:

```
- mikanani-push-bot
  - config.example.json
  - mikanani-push-bot
```

### 2. 配置

把 `config.example.json` 重命名为 `config.json`，并按下面的对照表修改配置档

```js
{
    // Bot Token, @BotFather
    "token": "",
    // 你的 TG 数字 ID，可以透过 @getmyid_bot 或者类似的 Bot 获取
    "admin": -1,
    // 来源字幕组，暂时只支援设定单个来源，对应的番剧不存在该来源就会添加失败
    "source": "ANi"
}
```

### 3. 运行

在 `config.json` 存在的目录中运行 `mikanani-push-bot`

```bash
chmod 775 ./mikanani-push-bot
./mikanani-push-bot
```

## 需要了解的事项
1. 这个专案使用了 JSON 档案数据库，如有必要可以二次开发使它支援 Database 以使其得到更高的并发性能
2. 这个专案不支援多使用者
3. 这个专案不会下载番剧本体，使用者仍需要手动下载种子

-----

### 使用 Linux Service

[<img src="https://opengraph.githubassets.com/0ce367d2a8cee652c1242cb4a99af11939ad2161e47eac849791a8695027a549/ArsFy/add_service" width="35%" style="border-radius: 5px" />](https://github.com/ArsFy/add_service)
