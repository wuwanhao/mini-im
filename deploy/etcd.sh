docker run -d \
-p 2379:2379 \
-p 2380:2380 \
--name etcd_3_5_13 \
quay.io/coreos/etcd:v3.5.13 \
/usr/local/bin/etcd --data-dir=/etcd-data --name node1 \
--listen-client-urls http://0.0.0.0:2379 \
--advertise-client-urls http://0.0.0.0:2379 \
--listen-peer-urls http://0.0.0.0:2380 \
--initial-advertise-peer-urls http://0.0.0.0:2380 \
--initial-cluster node1=http://0.0.0.0:2380 \
--log-level info --logger zap --log-outputs stderr