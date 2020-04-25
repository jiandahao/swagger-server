FROM jiandahao/swagger-ui:v1.0.0 as swagger_ui

FROM golang:1.13
# Copy swagger ui dist package
COPY --from=swagger_ui /usr/share/nginx/html/ /swagger-server/dist/
COPY bin/swagger-server /swagger-server/
COPY template /swagger-server/template

WORKDIR /swagger-server/
EXPOSE 8088
CMD ["./swagger-server"]