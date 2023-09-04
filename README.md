# 高性能高可用的分布式抖音后端应用
## 项目使用grpc+mysql+redis开发，使用gin实现请求访问
## 使用Kubernetes进行容器化部署，实现高可用部署、服务注册与发现
## 使用istio实现流量分发，与稳定性治理（限流与熔断）

数据库设计文件在：\kubernetes-manifests\db_deploymnet\mysql
涉及到部署的代码部分：如服务地址端口，以实际部署为准