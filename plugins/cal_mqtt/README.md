# MQTT subscription plugin yaml file input config.
## For the set subscription topic,

Please use the below format to subscribe and publish to MQTT.

```{"1": [{"topic":"pressure", "outputtopic":"caltemprature"}]}```

**topic:** This topic will be subscribed to and checked for value changes. If there is any change in this topic, then it will be published to MQTT with 'outputtopic'.
***outputtopic:*** An MQTT message will be published on this topic if there are any changes in the subscribed MQTT topic (which is provided in the 'topic').

**For Subtopics**
We can also subscribe to multiple subtopics with #, and below is an example of that.

```{"1": [{"topic":"ia/raw/opcua/#", "outputtopic":"caltemprature"}]}```

**topic:** In this format, we will subscribe to all subtopics that are under 'ia/raw/opcua/' like 'ia/raw/opcua/temparature', 'ia/raw/opcua/pressure', etc.
**outputtopic:** In this subtopic format, the output topic will be "topic/outputtopic", so in the above example, the output topic will be 'temparature/caltemprature', 'pressure/caltemprature'.


## For set tSubscription,
Please use the below format to subscribe to tSubscriptions.

```'{"1": [{"topic":"caljson", "elements": "element1,element2"},{"topic":"caljson1", "elements": "element3,element4"}]}'```

**topic:** provided topics will be subscribed like caljson, caljson1, and the above format is dynamic; you can add more topics in the above JSON array.
**elements:** a comma-separated list of all elements that you want in the output from that topic, like element1, element2 will be from caljson topic, and element3 and element4 will be from caljson1 topic.

**So the output of both of these is:**
All elements from subscription topic (pressure) + element1, element2 from caljson + element3, element4 from caljson1

**below is the example benthos input for MQTT**
```
input:
  mqtttrigger:
    url: localhost:1883
    client_id: benthos-mqtt
    subscriptions:
      - '{"1": [{"topic":"ia/raw/opcua/#", "outputtopic":"caltemprature"}]}'
    tsubscriptions:
      - '{"1": [{"topic":"caljson", "elements": "Cal_Runtime,Cal_Qty,Cal_scrap,Cal_OEE,Cal_Availability,averageJogVelocity"},{"topic":"caljson1", "elements": "job_product,job_shoporder,job_machineRate"}]}'
    qos: 1
    keepalive: 60
    timeout: 10s 
```

**Below is the example JSON.**
```
{
	"temparature": "123456",
	"db": "mysql",
	"element1": "1234",
	"element2": "string data",
	"element3": "another data",
	"element4": "567890",
	"timestamp_ms":1711631870112
}
```
