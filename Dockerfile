#FROM mermade/swagger2openapi:latest
#COPY pkg/swagger /swagger/
#ADD swagger2openapi.sh /usr/local/bin/
#WORKDIR /swagger/
#RUN /usr/local/bin/swagger2openapi.sh
# swagger-ui
# TODO： 更改镜像
FROM gitlab.xunlei.cn/xbase/swagger-ui:v0.9.3
COPY bin/swagger-server /usr/local/bin

CMD ["./swagger-server"]
