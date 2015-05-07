Dockerfile for [unBALANCE](https://github.com/jbrodriguez/unbalance)

To run

```shell
docker run -d --name unbalance \
-v /path/to/config/dir:/config \
-v /path/to/log/dir:/log \
-v /mnt:/mnt \
-v /root:/root \
-p 6237:6237 \
jbrodriguez/unbalance
```