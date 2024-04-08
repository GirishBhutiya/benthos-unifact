# CSV subscription plugin yaml file input config.

Please use the below format to read from the CSV
```
input:
  csv:
    timeout: 15 
    ftphost: "217.21.84.185"
    ftpusername: "u974007693"
    ftppassword: "Password99!@3"
    ftpport: 21
    path: ./csv
    sourcetype: "networkpath"
    networkpath: "./csv"
    file: "/HMITT.csv"
    node: '{"1": [{"name":"Pressure", "group": "D001", "db": "mssql", "historian": "influx", "sqlSp": "sp_sql_logging"}]}'
```	
**sourcetype:** can be networkpath or ftp. if set ftp then ftphost,ftpusername, ftppassword and ftpport is required. if networkpath is set then a file path is required.<br />
**file:** filename with extension. <br />
**timeout:** time to set timeout time<br />
**ftphost:** ftp host URL<br />
**ftpusername:** username of ftphost<br />
**ftppassword:** password of ftphost<br />
**ftpport:** ftp port 21 or 22<br />
**node:** parameter which you want in the output.<br />

**The plugin will read the last row of the CSV file and send it to output with the node. Once there is a new row and the last row value changes it will send that row to output.**