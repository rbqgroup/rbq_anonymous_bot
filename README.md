# rbq_anonymous_bot

频道小编工具：向绒频道和群组匿名发表内容。

频道默认不开启匿名模式以满足部分交互需求，如果有编辑者仍想保持匿名请使用这个 bot 。

## 使用

使用 `/会话代号 [内容]` 向预设好的 会话ID 发送内容。

## 部署

开发时所用 golang 版本: `1.19.5`

```sh
go get
go build
```

然后创建配置文件 `config.json` 和编译出来的可执行文件放一起。

### 配置文件示例

```json
{
    "ver": 1,
    "proxy": "代理服务器地址(可选)",
    "apikey": "TG API KEY",
    "timeout": 600,
    "whitelist": [
        00000000
    ],
    "g": "G-0000000000000",
    "c2": "C-0000000000000",
    "c25": "C-0000000000000",
    "c3": "C-0000000000000",
    "g18": "G-0000000000000",
    "gy": "G-0000000000000"
}
```

- `whitelist` 是白名单，只允许这些 UID 使用这个 BOT 。
- `C` 开头的会话 ID 表示这是一个 频道 。
- `G` 开头的会话 ID 表示这是一个 群组 。

## 许可

Copyright (c) 2023 KagurazakaYashi rbq_anonymous_bot is licensed under Mulan PSL v2. You can use this software according to the terms and conditions of the Mulan PSL v2. You may obtain a copy of Mulan PSL v2 at: http://license.coscl.org.cn/MulanPSL2 THIS SOFTWARE IS PROVIDED ON AN “AS IS” BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE. See the Mulan PSL v2 for more details.
