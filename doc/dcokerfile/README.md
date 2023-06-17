# docker 打包

打包需要的文件如下
```
├── Dockerfile
├── dst-admin-go
├── run.sh
```

**打包**
```
docker build -t dst-admin-go .
```

**运行**
```
dcoker run -d -p8082:8082 dst-admin-go
```
等看到 docker logs  说明启动成功了

```text
[GIN-debug] GET    /favicon.ico              --> github.com/gin-gonic/gin.(*RouterGroup).StaticFile.func1 (6 handlers)
[GIN-debug] HEAD   /favicon.ico              --> github.com/gin-gonic/gin.(*RouterGroup).StaticFile.func1 (6 handlers)
[GIN-debug] GET    /asset-manifest.json      --> github.com/gin-gonic/gin.(*RouterGroup).StaticFile.func1 (6 handlers)
[GIN-debug] HEAD   /asset-manifest.json      --> github.com/gin-gonic/gin.(*RouterGroup).StaticFile.func1 (6 handlers)
[GIN-debug] GET    /                         --> github.com/gin-gonic/gin.(*RouterGroup).StaticFile.func1 (6 handlers)
[GIN-debug] HEAD   /                         --> github.com/gin-gonic/gin.(*RouterGroup).StaticFile.func1 (6 handlers)
[GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.
[GIN-debug] Listening and serving HTTP on :8082
```