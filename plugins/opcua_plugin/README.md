# OPCUA subscription plugin yaml file input config.

## For subscription implementation

Please use the below format to read from OPCUA.

```{"1": [{"node": "ns=2;s=Pressure", "name":"Pressure", "group": "D001", "db": "mssql", "historian": "influx", "sqlSp": "sp_sql_logging"}]}```

**node:** node id which you want to read.
**name:** which you want in the output
**group:** group name
**db:** db name 
**historian:** historian name
**sqlSp:** stored procedure name

**Please use the below format to subscribe to OPCUA**
```
input:
  opcua:
    endpoint: "opc.tcp://localhost:46010"
    nodeIDs:
      - '{"1": [{"node": "ns=2;s=Pressure", "name":"Pressure", "mtopic":"runtime","group": "D001", "db": "mssql", "historian": "influx", "sqlSp": "sp_sql_logging"}]}'
      - '{"2": [{"node": "ns=2;s=Temperature", "name":"Temperature", "mtopic":"runtime", "group": "D002", "db": "mssql", "historian": "influx", "sqlSp": "sp_sql_logging"}]}'
    subscribeEnabled: true
    insecure: true
```	

**below is an example of output.**
```
{
	"Pressure":"86.29516419786509",
	"db":"mssql",
	"group":"D001",
	"historian":"influx",
	"sqlSp":"sp_sql_logging",
	"timestamp_ms":1712419092292
}
```

## For set tBatchNodeIDs,
**Please use the below format to subscribe to tBatchNodeIDs.**
```
input:
  opcuatrigger:
    endpoint: "opc.tcp://localhost:46010"
    tNodeIDs:
      - '{"1": [{"node": "ns=2;s=Pressure", "group": "D001", "db": "mssql", "historian": "influx", "sqlSp": "sp_sql_logging"}]}'
      - '{"2": [{"node": "ns=2;s=Temperature", "group": "D002", "db": "mssql", "historian": "influx", "sqlSp": "sp_sql_logging"}]}'
    subscribeEnabled: true
    insecure: true
    tBatchNodeIDs:
      - '{"1": [{"node":"ns=2;s=Pressure", "name": "PressureNodeName"}, {"node": "ns=2;s=Humidity", "name": "HumidityNodeName"}]}'
      - '{"2": [{"node":"ns=2;s=Temperature","name": "TemperatureNodeName"}, {"node": "ns=2;s=Air Quality", "name": "Air Qualit Node Name"}]}'
```

**both tNodeIDs and tBatchNodeIDs lengths must be the same.**

So in tNodeIDs when the subscription value change then you will get tBatchNodeIDs value in output.  like when ```ns=2;s=Pressure``` value change ```ns=2;s=Pressure``` and ```ns=2;s=Humidity``` will be in output. 

**below is an example of output.**
```
{
	"HumidityNodeName":"39.837678163441815",
	"PressureNodeName":"82.27382485048597",
	"db":"mssql",
	"group":"D001",
	"historian":"influx",
	"sqlSp":"sp_sql_logging",
	"timestamp_ms":1712418932043,
	"trigger":"ns_2_s_Pressure",
	"value":82.27382485048597
}
```