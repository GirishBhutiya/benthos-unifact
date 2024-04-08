# Modbus subscription plugin yaml file input config.

## For subscription implementation

Please use the below format to read from Modbus.

```{"1": [{"address": "1", "addresstype":"inputregister","name":"Pressure", "group": "D001", "db": "mssql", "historian": "influx", "sqlSp": "sp_sql_logging"}]}```

**address:** Modbus address which you want to read.<br />
**addresstype:** Modbus address type like inputregister, coils, discrete, and holding<br />
**name:** which you want in the output<br />
**group:** group name<br />
**db:** db name <br />
**historian:** historian name<br />
**sqlSp:** stored procedure name<br />

```
modbus:
    endpoint: "localhost:10502"
    slaveid: 1
    timeout: 10
    subscriptions: 
      - '{"1": [{"address": "1", "addresstype":"inputregister","name":"Pressure", "group": "D001", "db": "mssql", "historian": "influx", "sqlSp": "sp_sql_logging"}]}'
      - '{"2": [{"address": "1", "addresstype":"coils","name":"Temprature", "group": "D001", "db": "mssql", "historian": "influx", "sqlSp": "sp_sql_logging"}]}'
      - '{"3": [{"address": "1", "addresstype":"discrete","name":"Discrete", "group": "D001", "db": "mssql", "historian": "influx", "sqlSp": "sp_sql_logging"}]}'
      - '{"3": [{"address": "1", "addresstype":"holding","name":"Holding", "group": "D001", "db": "mssql", "historian": "influx", "sqlSp": "sp_sql_logging"}]}'
    subscribeEnabled: true
```

**below is an example of output.**
```
{
	"db":"mssql",
	"group":"D001",
	"historian":"influx",
	"name":"Holding",
	"sqlSp":"sp_sql_logging",
	"timestamp_ms":1712332006930,
	"value":100
}
```

## For set tSubscription,
Please use the below format to subscribe to tSubscriptions.
```
modbustrigger:
    endpoint: "localhost:10502"
    slaveid: 1
    timeout: 10
    subscriptions: 
      - '{"2": [{"address": "1", "addresstype":"coils","name":"Temprature", "group": "D001", "db": "mssql", "historian": "influx", "sqlSp": "sp_sql_logging"}]}'
      - '{"3": [{"address": "1", "addresstype":"holding","name":"Holding", "group": "D001", "db": "mssql", "historian": "influx", "sqlSp": "sp_sql_logging"}]}'
    tsubscriptions:
      - '{"1": [{"address":"1","addresstype":"holding", "name": "PressureNodeName"}, {"address": "3", "addresstype":"coils","name": "HumidityNodeName"}]}'
      - '{"2": [{"address":"4","addresstype":"discrete","name": "TemperatureNodeName"}, {"address": "5", "addresstype":"holding","name": "Air Qualit Node Name"}]}'
```	  

**both subscription and tsubscription length must be the same.**
So in tsubsription when the subscription value changes then you will get tsubscription value in the output.  like when the coils value changes holding address 1 and coils address 3 will be in output. 

**below is an example of output.**
```
{
	"Air_Qualit_Node_Name":100,
	"TemperatureNodeName":0,
	"db":"mssql",
	"group":"D001",
	"historian":"influx",
	"sqlSp":"sp_sql_logging",
	"timestamp_ms":1712332088851,
	"trigger":"Holding",
	"value":100
}
```