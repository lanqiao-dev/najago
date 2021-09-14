# Golang 开发工具包

* validate, grpc 参数校验拦截器；
* i18n, 基于  go-i18n 的国际化相关的辅助函数；

## i18n 模块使用

目前 i18n 模块主要提供了相关工具函数，方便在基于 `context` 模块完成相关的 `Localizer` 设置和提取。`i18n/grpc.go` 文件主要提供了在 grpc 拦截器实现，便于在 grpc 服务中实现相关翻译功能。
### 翻译文件的相关处理

1. 通过 `go get -u github.com/nicksnyder/go-i18n/v2/goi18n` 安装 goi18n 命令工具；

2. 在源码中通过如下形式添加需要翻译的消息

```golang

package i18n

import (
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

var example = &goi18n.Message{ID: "invalid-captcha", Other: "Invalid captcha.", Description: "example message"}
```

3. 通过执行 `goi18n extract` 命令提取源码中的所有 `*i18n.Message` 消息，默认情况下输出 `active.en.toml` 文件，如果有必要可以通过 `-outdir` 指定输出目录。比如生成的 `active.en.toml` 内容如下：

```toml
[invalid-captcha]
description = "example message"
other = "Invalid captcha."
```

4. 创建想要翻译的目标语言文件，比如中文 `translate.zh.toml` 文件，然后执行 `goi18n merge active.en.toml translate.zh.toml`，此时 `translate.zh.toml` 文件内容如下。

```toml
[invalid-captcha]
description = "example message"
hash = "sha1-72a2384350d8939135f112a5c24f2da94875f1da"
other = "Invalid captcha."
```

5. 翻译 `translate.zh.toml` 内如，比如翻译后的内容如下

```toml
[invalid-captcha]
description = "example message"
hash = "sha1-72a2384350d8939135f112a5c24f2da94875f1da"
other = "无效的验证码。"
```

6. 重命名 `translate.zh.toml` 为 `active.zh.toml`，并将 `active.zh.toml` 通过 golang 代码加载；

7. 如果有新增消息时，执行  `goi18n extract` 提取新增的消息，并执行 `goi18n merge active.*.toml` 此时会输出所有新增的消息到 `translate.*.toml` 文件中，比如 `translate.zh.toml` 文件。待 `translate.zh.toml` 中新增的消息翻译完成后，再通过 `goi18n merge active.*.toml translate.*.toml` 进行合并。此时 `translate.*.toml` 文件中新增的消息会被添加到对应的 `active.*.toml` 文件中。新增消息翻译工作即完成。