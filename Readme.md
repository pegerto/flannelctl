### Flannelctl view utility

Help debuging fetching flannel information form etcd

- View flannel configuration

```
flannelctl config
+--------------+---------------------+
|              |               Value |
+==============+=====================+
|      Network |      192.168.0.0/17 |
+--------------+---------------------+
|    SubnetLen |                  26 |
+--------------+---------------------+
|      Backend |    {"Type":"vxlan"} |
+--------------+---------------------+
```


- View active subnets 

```
flannelctl subnet
+----------------------+------------------+----------------+
|               Subnet |         PublicIP |    BackendType |
+======================+==================+================+
|    192.168.10.128/26 |    10.10.100.101 |          vxlan |
+----------------------+------------------+----------------+
|     192.168.5.128/26 |    10.10.100.102 |          vxlan |
+----------------------+------------------+----------------+
|     192.168.3.128/26 |    10.10.100.103 |          vxlan |
+----------------------+------------------+----------------+
|     192.168.2.192/26 |    10.10.100.104 |          vxlan |
+----------------------+------------------+----------------+

```