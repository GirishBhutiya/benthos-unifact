create or replace procedure transfer(
   "HumidityNodeName":dec, 
   "PressureNodeName":dec, 
   "sqlSp": character varying, 
   "historian": character varying, 
   "db": character varying, 
   "group": character varying, 
   "trigger": character varying, 
   "timestamp_ms": int, 
   "value": dec
)
language plpgsql    
as $$
begin
     
    INSERT INTO data VALUES (
        "HumidityNodeName",
        "PressureNodeName",
        "sqlSp",
        "historian", 
        "db", 
        "group", 
        "trigger", 
        "timestamp_ms",
        "value"
    );

    commit;
end;$$;