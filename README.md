# rbq_anonymous_bot

频道小编工具：向绒频道和群组匿名发表内容。

频道默认不开启匿名模式以满足部分交互需求，如果有编辑者仍想保持匿名请使用这个 bot 。

## 功能

- 匿名发送内容到各个频道或群组。
- 解析推特推文链接，自动移除跟踪代码，拉取并发表推文里面包含的图文。
- 自动为视频信息加 `#视频` 标签。
- 一起当 _~更新姬~_ 。

## 使用

使用 `/会话代号 [内容]` 向预设好的 会话ID 发送内容。

- 当传送的内容（除了命令外）只有文本且只是一条 Twitter 链接，则解析该推文（自动取出 评论/转推/喜欢/回复 的数量、推文作者名和昵称、正文、图片、视频）然后发送到目标，并移除跟踪代码，而无需自己上传图文。可以在配置文件中禁用和指定拉取服务器。
  - 注意：这项功能只能用于**文字和图片**推文，视频等其他类型推文不支持。如果开启此功能的情况下发送包含不支持附件的单推文链接将会导致失败。
- 当传送的内容为视频时，自动在正文前面添加 `#视频` 标签。可以在配置文件中禁用或自定义此功能。
- 输入 `/chatid` 命令时，将返回当前会话和用户的名称和 ID 。该命令仅在配置文件中启用调试模式时有效，并且无视使用者白名单限制。可用于使用者加入白名单之前获取自己的 用户/群组/频道 的 ID 。

## 部署

开发时所用 golang 版本: `1.19.5`

```sh
go get
go build
```

然后创建配置文件 `config.json` ，和编译出来的可执行文件 `rbq_anonymous_bot` 放一起。

### 使用 Docker 部署

1. 使用 `build_linux.bat` 或参考里面的操作生成可执行文件压缩档 `bin/rbq_anonymous_bot.xz` 。
2. 修改 `./docker.sh` 为需要的 Docker 操作。
3. 将 `bin/rbq_anonymous_bot.xz` + `config.json` + `Dockerfile` + `docker.sh` 复制到服务器中的同一个文件夹中。
4. 进入服务器中的该文件夹，执行 `chmod +x docker.sh` 和 `./docker.sh` 即可运行
5. 让 bot 转发一条消息，等待大约一分钟，该 Docker 容器状态会显示为 `healthy` 。
6. 如果没有出现停止问题，可以将 RESTART POLICIES 设置为 `Always` 。

### 配置文件示例

```json
{
    "ver": 1,
    "debug": true,
    "proxy": "http://127.0.0.1:8080",
    "apikey": "xxxxxxxxxx:*******-***********-***************",
    "healthcheck": "healthcheck.lock",
    "timeout": 600,
    "whitelist": [
        00000000
    ],
    "to": {
        "d": "C-0000000000000",
        "g": "G-0000000000000",
        "c2": "C-0000000000000",
        "c25": "C-0000000000000",
        "c3": "C-0000000000000",
        "g18": "G-0000000000000",
        "gy": "G-0000000000000"
    },
    "nitterHost": "nitter.net",
    "headVideo": "#视频 ",
    "headAnimation": "",
    "headPhoto": "",
    "headText": ""
}
```

- `var` 配置文件版本号（填 `1` ）。
- `debug` 调试模式。显示所有通信日志，并将无命令的内容直接返回。
- `proxy` 是代理服务器，支持 `http` 和 `socks5`，不需要时留空字符串。
- `apikey` Telegram 的会话令牌（去问 [BotFather](https://t.me/BotFather) 要）。
- `healthcheck` Docker 健康检查用会话文件名，需要和 `Dockerfile` 中的 `HEALTHCHECK` 想对应。
- `whitelist` 是白名单，只允许这些 UID 使用这个 BOT 。
- `to` 是会话代号（预定义的发送目标）。
  - key 是命令，例如 `"c2"` 表示 `/c2` 命令。
  - `C` 开头的会话 ID 表示这是一个 **频道** 。
  - `G` 开头的会话 ID 表示这是一个 **群组** 。
  - `P` 开头的会话 ID 表示这是一个 **私聊** 。
- `nitterHost` 当传送的内容（除了命令外）只有文本且只是一条 Twitter 链接，则解析该推文（自动取出 评论/转推/喜欢/回复 的数量、推文作者名和昵称、正文、图片、视频）然后发送到目标，而无需自己上传图文。由于 Twitter 对 API 的限制，因此借由 Nitter 来达到此目标，你需要在此指定一个 Nitter 服务器域名，建议使用自建服务器。此项留空字符串时则禁用这个功能。
- `head*` 是自动在发送的不同类型消息中添加的前缀。
  - 视频优先级高于其他。

## 许可

Copyright (c) 2023 KagurazakaYashi rbq_anonymous_bot is licensed under Mulan PSL v2. You can use this software according to the terms and conditions of the Mulan PSL v2. You may obtain a copy of Mulan PSL v2 at: http://license.coscl.org.cn/MulanPSL2 THIS SOFTWARE IS PROVIDED ON AN “AS IS” BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE. See the Mulan PSL v2 for more details.
