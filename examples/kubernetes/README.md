# Kubernetes 部署示例

本目录提供基于 Kustomize 的“一键部署”示例，包含 MySQL 数据库初始化和 API Server 所需的基础资源，适用于本地开发/测试/演示环境。

重要：请从本目录运行部署命令，示例使用的命名空间为 `bfe-system`（由 `kustomization.yaml` 指定）。

## 目录文件

- `namespace.yaml` — 创建命名空间 `bfe-system`。
- `mysql.yaml` — MySQL 组件：`Secret`、`Service`、`Deployment`、`ConfigMap`（内嵌 `db_ddl.sql`）和初始化 `Job`（`mysql-init-job`）。
- `configmap.yaml` — API Server 配置（ConfigMap）。
- `deployment.yaml` — API Server 的 `Deployment`。
- `service.yaml` — API Server 的 `Service`。
- `kustomization.yaml` — Kustomize 清单（一键组合上述资源）。
- `README.md` — 本说明文档。

## 一键部署

在 `examples/kubernetes` 目录下执行：

```sh
kubectl apply -k .
```

## 验证部署是否成功

- 检查命名空间和资源是否创建：

```sh
kubectl get ns bfe-system
kubectl get all -n bfe-system
```

- 检查 Pod 状态（等待所有 Pod 状态为 `Running`/`Completed`）：

```sh
kubectl get pods -n bfe-system
```

- 检查初始化 Job（`mysql-init-job`）状态并查看日志：

```sh
kubectl get job mysql-init-job -n bfe-system
# 获取 job 关联的 Pod 名称
JOB_POD=$(kubectl get pods -n bfe-system -l job-name=mysql-init-job -o jsonpath='{.items[0].metadata.name}')
kubectl logs -n bfe-system $JOB_POD
```

- 验证 MySQL 中的表：

```sh
# 获取 mysql Pod
MYSQL_POD=$(kubectl get pods -n bfe-system -l app=mysql -o jsonpath='{.items[0].metadata.name}')

# 从 Secret 读取 root 密码并执行 SHOW TABLES
MYSQL_PW=$(kubectl get secret mysql-secret -n bfe-system -o jsonpath='{.data.mysql-root-password}' | base64 --decode)

kubectl exec -n bfe-system -it $MYSQL_POD -- \
  mysql -uroot -p"$MYSQL_PW" -e "SHOW TABLES IN open_bfe;"
```

期望结果：应包含以下表（示例最终版）：

- `bfe_clusters`, `products`, `domains`, `clusters`, `lb_matrices`, `sub_clusters`, `pools`,
- `route_basic_rules`, `route_advance_rules`, `route_cases`, `certificates`, `extra_files`,
- `config_versions`, `users`, `user_products`

如果这些表都存在，则初始化成功。

## 访问控制台

部署成功后，通过浏览器访问 Node 的 `30083` 端口并登录控制台:

- URL：http://{NodeIP}:30083/
- 初始账户密码：admin/admin

说明：端口 `30083` 在 `service.yaml` 中以 `nodePort: 30083` 的形式配置（即 NodePort）。

## 常见故障排查

1. 初始化 Job 报错（SQL 语法或执行中断）

  - 查看 Job Pod 日志（如上）。常见错误包括 SQL 语法错误（如 `ERROR 1064`），或脚本在中间被中断。
  - 如果 Job 日志显示 SQL 语法错误，定位到错误行并修复 `mysql.yaml` 中的 `ConfigMap`（`db_ddl.sql`）内容，然后删除并重建 Job：

```sh
# 编辑并保存修改后的 mysql.yaml（或通过 CI/本地脚本生成新的 ConfigMap）
kubectl delete job mysql-init-job -n bfe-system
kubectl apply -k .
```

  - 注意：本示例使用 `emptyDir` 为 MySQL 存储（无持久化）。如果 Job 已成功写入数据但需要重新运行，先清空 MySQL 数据目录或删除 Pod（会丢失数据）；在生产场景请改用 `PersistentVolume`。

2. 找不到 MySQL Pod 或 Pod CrashLoop

```sh
kubectl describe pod <pod-name> -n bfe-system
kubectl logs <pod-name> -n bfe-system
```

检查是否为镜像拉取失败、权限、环境变量或资源限制问题。

3. 无法通过 API Server 访问数据库

- 检查 `Service` 是否存在并指向正确的 Pod 标签：

```sh
kubectl get svc -n bfe-system
kubectl describe svc mysql -n bfe-system
```

- 在 API Server Pod 内测试连接：

```sh
kubectl exec -n bfe-system -it deploy/api-server -- sh -c 'nc -vz mysql.bfe-system.svc.cluster.local 3306'
```

## 生产环境建议

- 使用 `StatefulSet` 与 `PersistentVolume` 替代 `Deployment` + `emptyDir`，以保证数据持久化。
- 不要在 ConfigMap 中保存敏感信息（使用 `Secret` 管理密码）。
- 使用 CI 将根目录的 `db_ddl.sql` 自动注入到 ConfigMap 的生成流程，避免手动同步差异。