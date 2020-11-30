FROM golang:alpine AS builder

#为我们镜像设置必要的环境变量
ENV GO111MODEULE=on \
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64 \
	GOPROXY=https://goproxy.cn,direct

#移动到工作目录: /build
WORKDIR /build

#将复制项目中的go.mod 和go.sum文件并下载依赖信息
COPY go.mod ./
COPY go.sum ./
RUN  go mod download

#代码复制到容器中
COPY . .

#将我们代码编译成二进制可以执行文件 app
RUN go build -o bluebell .

###########################
#  接下来创建一个小镜像       #
###########################
FROM centos
#复制配置文件
COPY ./wait-for.sh ./app/
COPY ./templates ./app/templates
COPY ./static ./app/static
#COPY ./conf ./app/conf

#builer镜像中把/build/bluebell_app 拷贝到/app
COPY --from=builder /build/bluebell /app/
#指定工作目录
WORKDIR /app

#映射端口
EXPOSE 9000

#RUN set -eux; \
#	yum  install -y netcat && \
#        chmod 755 wait-for.sh

#运行需要启动的命令
ENTRYPOINT ["./bluebell","./conf/config.yaml"]

