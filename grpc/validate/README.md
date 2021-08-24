# grpc 参数验证拦截

1. 通过 https://github.com/favadi/protoc-go-inject-tag 为 grpc request 结构体生成字段 [validator](https://github.com/go-playground/validator) tag；
2. 在 grpc 拦截器中校验请求参数；
3. 支持通过 context 传递 locale 信息；
