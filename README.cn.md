# snmp-agent
本组件主要为采集prometheus的exporter数据并上报至snmpd而设计。 

使用的时候只需修改对应的snmp-agent.toml配置文件，在配置文件中加入所需要采集的exporter，以及相关采集配置即可。

## 使用说明

### address
组件地址端口

### timeout
连接超时，单位秒

### reconnect_interval
重连间隔，单位秒

### base_oid
注册至snmpd的基础oid

### scan_interval
扫描间隔，单位秒

### origin
数据来源，当前可选 exporter 或者 prometheus

### server
数据源地址

### oid
用于指定上报至snmp的oid

### name
用于指定需要采集的metrics (直接从exporter采集时候需要)

### labels
用于对上述metrics进行label过滤 (直接从exporter采集时候需要)

写入Labels配置时，以字符串数组形式写入，label和value之间```","```分隔，多个标签之间同样以```","```分隔
```toml
label_nodes = ["label1","value1","label2","value2"]
```

### value_label
从过滤后的metric中，选择标签作为其值（仅exporter采集的数据能使用）
不填或者空值，默认以value作为值

### query
用于输入PromQL对指定prometheus服务器进行查询获取数据
使用query后，无须使用name、labels进行过滤

### value_type
指定当前oid的数据类型
```toml
CEILINT 表示向上取整的int32类型
STRING 表示不做任何处理的string类型
```

### data_type
指定子数据类型，默认为无子数据
```toml
TABLE 表示是一个表格类型数据，根据后续table结构进行处理
LIST 表示一个列表类型数据，根据后续list结构进行处理
```

### table
当data_type为TABLE时，指定table对应的entry

table可以设置oid，label，value_type

table必须指定oid，不适用exporters.items.oid

可供选择的table.label默认有value，此外为mertic.label

可供选择的table.value_type同value_type

### list
exporter暂时不支持list

当data_type为LIST时，指定list对应的item

list可以设置oid，label，value_type

list必须指定oid，不适用exporters.items.oid

可供选择的list.labels同item.labels

可供选择的list.value_type同value_type

### 注意
本组件不支持采集prometheus中untyped类型指标

## 示例
```toml
[snmpd]
address = "127.0.0.1:161"
timeout = 10
reconnect_interval = 10
base_oid = "1.3.6.1.4.1.47032"
scan_interval = 10

[[exporters]]
origin = "exporter"
server = "http://127.0.0.1:9100/metrics"

    [[exporters.items]]
    oid = "1.3.6.1.4.1.47032.1.1"
    name = "metric_name_1"
    labels = []
    data_type = "TABLE"

    [[exporters.items.table]]
    oid = "1.3.6.1.4.1.47032.1.1.1.1"
    label = "label_name_1"
    value_type = "STRING"

    [[exporters.items.table]]
    oid = "1.3.6.1.4.1.47032.1.1.1.3"
    label = "label_name_2"
    value_type = "STRING"

    [[exporters.items]]
    oid = "1.3.6.1.4.1.47032.1.2"
    name = "metric_name_2"
    labels = []
    is_table = false
    value_type = "STRING"

    [[exporters.items]]
    oid = "1.3.6.1.4.1.47032.1.3"
    name = "metric_name_3"
    labels = []
    value_label = "label_name"
    value_type = "STRING"

[[exporters]]
origin = "exporter"
server = "http://127.0.0.1:9100/metrics"

    [[exporters.items]]
    oid = "1.3.6.1.4.1.47032.1.4"
    name = "metric_name_4"
    labels = ["label_name", "label_value"]
    value_type = "STRING"

[[exporters]]
origin = "prometheus"
server = "http://127.0.0.1:9090"

    [[exporters.items]]
    oid = "1.3.6.1.4.1.47032.1.5"
    query = "metric_name_5{label=\"label_value"}"
    data_type = "TABLE"

    [[exporters.items.table]]
    oid = "1.3.6.1.4.1.47032.1.5.1.1"
    label = "label_name"
    value_type = "STRING"

    [[exporters.items]]
    oid = "1.3.6.1.4.1.47032.1.6"
    query = "metric_name_6{label=\"label_value\"}"
    data_type = "LIST"

    [[exporters.items.table]]
    oid = "1.3.6.1.4.1.47032.1.6.1.1"
    labels = ["label_name_1", "label_value_1", "label_name_2", "label_value_2"]
    value_type = "CEILINT"
```

## pdu类型对应的数据类型如下
```go
pdu.VariableTypeInteger: int32(-123)

pdu.VariableTypeOctetString: "echo test"

pdu.VariableTypeNull: nil

pdu.VariableTypeObjectIdentifier: "1.3.6.1.4.1.47032.1.5"

pdu.VariableTypeIPAddress: net.IP{10, 10, 10, 10}

pdu.VariableTypeCounter32: uint32(123)

pdu.VariableTypeGauge32: uint32(123)

pdu.VariableTypeTimeTicks: 123 * time.Second

pdu.VariableTypeOpaque: []byte{1, 2, 3}

pdu.VariableTypeCounter64: uint64(12345678901234567890) 
```
