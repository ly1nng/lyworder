# 打包依赖阶段使用golang作为基础镜像
From golang:1.19.13-bullseye AS builder

# 启用go module
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY . .

RUN make build

FROM golang:1.19.13-alpine

WORKDIR /app

COPY --from=builder /app/bin/* . 
COPY --from=builder /app/configs /app/configs
COPY --from=builder /app/static /app/static

EXPOSE 8080

ENTRYPOINT ["./lyworder","-port=8080","-config=configs/","-mode=release"] 

