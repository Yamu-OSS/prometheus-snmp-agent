# snmp-agent
This component is primarily designed to collect data from Prometheus exporters and report it to snmpd.

To use it, simply modify the corresponding snmp-agent.toml configuration file, adding the exporters you need to collect data from, along with their relevant collection configurations.

## Instructions

### address
Component address port.

### timeout
Connection timeout, in seconds.

### reconnect_interval
Reconnection interval, in seconds.

### base_oid
Base OID registered to snmpd.

### scan_interval
Scanning interval, in seconds.

### origin
Data source, currently selectable as either 'exporter' or 'prometheus'.

### server
Source data address.

### oid
Used to specify the OID reported to SNMP.

### name
Used to specify the metrics to be collected (required when collecting directly from exporters).

### labels
Used for labeling the above metrics (required when collecting directly from exporters).

When writing Labels configuration, it should be written in the form of a string array, with labels and values separated by ```","```, and multiple labels separated by ```","```.
```toml
label_nodes = ["label1","value1","label2","value2"]
```

### value_label
Selects a label from the filtered metric as its value (only applicable to data collected from exporters).
If left blank or empty, value is selected as default.

### query
Used to input PromQL to query specified Prometheus servers to obtain data.
When using query, there is no need to use name or labels for filtering.

### value_type
Specifies the data type of the current OID.
```toml
CEILINT represents int32 type rounded up
STRING represents string type with no processing
```

### data_type
Specifies the sub-data type, default is no sub-data.
```toml
TABLE represents table type data, processed according to subsequent table structure
LIST represents a list type data, processed according to subsequent list structure
```

### table
When data_type is TABLE, specify the entry corresponding to the table.

Table can set oid, label, value_type.

Table must specify oid, not applicable to exporters.items.oid.

Available table.labels default to value, and also include metric.labels.

Available table.value_type same as value_type.

### list
Exporters currently do not support list.

When data_type is LIST, specify the item corresponding to the list.

List can set oid, label, value_type.

List must specify oid, not applicable to exporters.items.oid.

Available list.labels same as item.labels.

Available list.value_type same as value_type.

### 注意
This component does not support collecting untyped metrics in Prometheus.

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

## PDU type corresponding data types are as follows
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
