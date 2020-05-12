# swagger-server
swagger server支持动态添加api文档，结合swagger ui，方便实现服务api文档的自动化管理。

- 主要实现思路：

用户通过http向服务端post一个api文档，服务端将完成文档的存储，并动态渲染swagger ui主页。

- 如何使用

  **启动服务**
  ```bash
  # 编译
  ./build.sh
  # 执行程序
  ./bin/swagger-server
  ```
  
  也可以通过docker完成服务部署：
  ```bash
  docker run -d -p 8088:8088 -v ${PWD}/docs:/swagger-server/docs jiandahao/swagger_server:v1.0.0
  ```
  
  `${PWD}/docs:/swagger-server/docs`挂载存放api文档的目录。

  **添加api文档**
  ```bash
  curl -X POST -H "Content-Type: multipart/form-data" -F "file=@${PWD}/_example/test_swagger.swagger.json" 127.0.0.1:8088/new
  ```

  访问文档： 在浏览器访问[127.0.0.1:8088](127.0.0.1:8088)

- 如何配合ci使用

  这边是以gitlab-ci为例了，使用其他ci工具道理是一样的。
  
  通常我们都是通过相关工具生成相应的swagger api文档，这里以protobuf为例，具体示例见[这里](_example/ci-demo),
  - swagger_deploy.sh: 完成根据proto文档自动生成相应的swagger api文档。
  - swagger2openapi.sh: 由于自动生成的swagger文档是2.0版本，swagger2openapi.sh完成将文档转换成openapi 3.0版本，这一步是可选的，是个人情况而定。

