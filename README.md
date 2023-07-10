![icon](macOS/rbqanonymousbot/Assets.xcassets/AppIcon.appiconset/rbq%205.png)

# [RBQ Anonymous Bot v1.2.2](https://github.com/rbqgroup/rbq_anonymous_bot)

频道小编工具：向绒频道和群组匿名发表内容。

频道默认不开启匿名模式以满足部分交互需求，如果有编辑者仍想保持匿名请使用这个 bot 。

## 功能

- 匿名发送内容到各个频道或群组。
- 接收投稿到指定频道或群组。
- 一键转发指定消息到指定频道或群组。
- 解析推特推文链接，自动移除跟踪代码，拉取并发表推文里面包含的图文。
  - 由于推特的一系列操作，该功能可能不再可用。
  - 目前实现方式为从多个 Nitter 实例中间接拉取推文，支持自动禁用获取失败的实例。
- 自动为视频信息加 `#视频` 标签。
- 一起当 _更新姬_ 。

## 安装

从 [Release](releases) 下载相应系统的可执行文件即可，无需安装。

| Release 文件（压缩包）   | 系统    | 最低版 | 位  | 体系结构            |
| ------------------------ | ------- | ------ | --- | ------------------- |
| `bin/*_Linux32.zip`      | Linux   | 2.6    | 32  | i386 (x86)          |
| `bin/*_Linux64.zip`      | Linux   | 2.6    | 64  | amd64(x86-64)       |
| `bin/*_macOSI64.dmg`     | macOS   | 10.13  | 64  | amd64(x86-64)       |
| `bin/*_macOSM64.dmg`     | macOS   | 11     | 64  | arm64(AppleSilicon) |
| `bin/*_Windows32.cab`    | Windows | 7      | 32  | i386 (x86)          |
| `bin/*_Windows64.cab`    | Windows | 7      | 64  | amd64(x86-64)       |
| `bin/*_WindowsARM64.cab` | Windows | 10     | 64  | arm64(aarch64)      |

## 使用

**本 bot 仅支持 supergroup 群组！**

### 观众投稿收件箱

直接私聊发送图片或视频，并附带来源链接，将会自动将这些内容转发到设置好的投稿接收频道/群组中，等待编辑进行审核处理。

**在投稿收件箱群组中，本 bot 必须为管理员！**

### 管理员投稿到频道

使用 `/会话代号 [内容]` 向预设好的 会话 ID 发送内容。

- 当传送的内容（除了命令外）只有文本且只是一条 Twitter 链接，则解析该推文（自动取出 评论/转推/喜欢/回复 的数量、推文作者名和昵称、正文、图片、视频）然后发送到目标，并移除跟踪代码，而无需自己上传图文。可以在配置文件中禁用和指定拉取服务器。
  - 注意：这项功能只能用于**文字和图片**推文，视频等其他类型推文不支持。如果开启此功能的情况下发送包含不支持附件的单推文链接将会导致失败。
- 当使用命令对一个消息进行回复时，则将这条消息转发到目标。
- 当传送的内容为视频时，自动在正文前面添加 `#视频` 标签。可以在配置文件中禁用或自定义此功能。
- 输入 `/<bot用户名>` 或 `/<bot用户名>@<bot用户名>` 命令时，将回复关于信息和当前会话信息。bot 名称自动取当前名称。

## 编译

开发时所用 golang 版本: `1.19.5`

```sh
go get
go build
```

### 跨平台编译

在 Windows x64 中也可以通过批处理一键生成全平台二进制文件：

```bat
build.bat
```

批处理脚本最后会调用 `MAKECAB` 和 `7z` 命令进行压缩。

## 部署

1. 创建配置文件 `config.json` ，和编译出来的可执行文件放一起。
2. Linux 或 macOS 需要使用 `chmod +x [可执行文件名]` 给予权限。
3. 运行可执行文件。

### 使用 Docker 部署

1. 创建配置文件 `config.json` 。
2. 使用 `build_linux.bat` 或参考里面的操作生成可执行文件压缩档 `bin/rbq_anonymous_bot.xz` 。
3. 修改 `./docker.sh` 为需要的 Docker 操作。
4. 将 `bin/rbq_anonymous_bot.xz` + `config.json` + `Dockerfile` + `docker.sh` 复制到服务器中的同一个文件夹中。
5. 进入服务器中的该文件夹，执行 `chmod +x docker.sh` 和 `./docker.sh` 即可运行
6. 让 bot 转发一条消息或与 bot 私聊输入 `/<bot用户名>`，等待大约一分钟，该 Docker 容器状态会显示为 `healthy` 。
7. 如果没有出现停止问题，可以将 RESTART POLICIES 设置为 `Always` 。

### macOS 系统中添加启动参数

1. 打开 Release 中的相应平台的 `.dmg` 文件，找到里面的 `.app` 文件，将其复制到 `应用程序` 文件夹.
2. 右键点击改 `.app` 文件，选择 `显示包内容` 。
3. 编辑 `Contents/Resources/run.sh` 脚本文件，在里面注释位置处添加参数。

### 配置文件示例

```json
{
  "ver": 1,
  "debug": -1,
  "proxy": "http://127.0.0.1:8080",
  "apikey": "xxxxxxxxxx:*******-***********-***************",
  "healthcheck": "healthcheck.lock",
  "timezone": 8,
  "timeout": 600,
  "whitelist": [00000000],
  "defto": -000000000,
  "to": {
    "t": "G-0000000000000",
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

- `var`(`int8`): 配置文件版本号（填 `1` ）。
- `debug`(`int64`): 调试模式。显示所有通信日志，并将无命令的内容直接返回给指定 ID。填写 `-1` 为关闭，填写 `0` 为仅开启，填写 用户ID 为开启并将这个 用户ID 在任意处发的内容都转发给这个 用户ID 。
- `proxy`(`string`) 是代理服务器，支持 `http` 和 `socks5`，不需要时留空字符串。
- `apikey`(`string`) Telegram 的会话令牌（去问 [BotFather](https://t.me/BotFather) 要）。
- `healthcheck`(`string`): Docker 健康检查用会话文件名，需要和 `Dockerfile` 中的 `HEALTHCHECK` 相对应。
- `timezone`(`int8`): GMT 时间偏移量，用于显示时间时所用的时区。
- `timeout`(`int`): 与 Telegram 服务器的连接超时时间。
- `whitelist`(`[]int64`): 是白名单，只允许这些 UID 使用这个 BOT 。
- `defto`(`int64`) 是默认收件人（投稿收件箱），直接私聊发送图片或视频，并附带来源链接，将会自动将这些内容转发到这里设置好的投稿接收频道/群组中，等待编辑进行审核处理。
- `to`(`map[string]string`): 是会话代号（预定义的发送目标）。
  - key 是命令，例如 `"c2"` 表示 `/c2` 命令。
  - `C` 开头的会话 ID 表示这是一个 **频道** 。
  - `G` 开头的会话 ID 表示这是一个 **群组** 。
  - `P` 开头的会话 ID 表示这是一个 **私聊** 。
- `nitterHost`(`[]string`): 当传送的内容（除了命令外）只有文本且只是一条 Twitter 链接，则解析该推文（自动取出 评论/转推/喜欢/回复 的数量、推文作者名和昵称、正文、图片、视频）然后发送到目标，而无需自己上传图文。由于 Twitter 对 API 的限制，因此借由 Nitter 来达到此目标，你需要在此指定一个或多个 Nitter 服务器域名（程序将自动按顺序取用并屏蔽失败的域名），此项留空字符串时则禁用这个功能。
- `head*`(`string`): 是自动在发送的不同类型消息中添加的前缀。
  - 视频优先级高于其他。

**所有数组类型条目数不要超过 100 条。**

## 许可

Copyright (c) 2023 KagurazakaYashi rbq_anonymous_bot is licensed under Mulan PSL v2. You can use this software according to the terms and conditions of the Mulan PSL v2. You may obtain a copy of Mulan PSL v2 at: http://license.coscl.org.cn/MulanPSL2 THIS SOFTWARE IS PROVIDED ON AN “AS IS” BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE. See the Mulan PSL v2 for more details.
