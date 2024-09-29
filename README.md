## DEV env 部署步骤
1. 运行 etcd.sh 脚本，启动 etcd
2. docker-compose up -d，启动 mysql 和 redis
3. 执行 deploy/sql 目录下的 sql 脚本，构建基础表结构