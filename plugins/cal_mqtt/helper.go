package cal_mqtt

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"golang.org/x/exp/maps"
)

func parseMergeTopics(mergeTopicString []string) []mergeTopics {

	var topics = make([]mergeTopics, len(mergeTopicString))

	for i, v := range mergeTopicString {
		var mergeMap map[string][]map[string]string
		var mTopic mergeTopics

		err := json.Unmarshal([]byte(v), &mergeMap)
		if err != nil {
			log.Println(err)
			return topics
		}
		for key, values := range mergeMap {
			var topic mergeTopic
			var subTopics []mergeTopic
			for _, obj := range values {

				topic.topic = obj["topic"]
				topic.elements = strings.Split(obj["elements"], ",")
				subTopics = append(subTopics, topic)

			}
			mTopic.ID, _ = strconv.Atoi(key)
			mTopic.topics = subTopics

		}
		log.Println(mTopic)
		topics[i] = mTopic
	}
	return topics
}
func parseInputTopic(inputTopicString []string) inputTopic {
	var inputtopic inputTopic
	for _, v := range inputTopicString {
		var mergeMap map[string][]map[string]string

		err := json.Unmarshal([]byte(v), &mergeMap)
		if err != nil {
			log.Println(err)
			return inputtopic
		}
		for _, values := range mergeMap {

			for _, obj := range values {

				inputtopic.topic = obj["topic"]
				inputtopic.outputTopic = obj["outputtopic"]

			}
		}

	}
	return inputtopic
}
func (b *mqttReader) apply(opts *mqtt.ClientOptions) *mqtt.ClientOptions {
	opts = opts.SetAutoReconnect(false).
		SetClientID(b.clientId).
		SetConnectTimeout(b.timeout).
		SetKeepAlive(time.Duration(b.keepAlive) * time.Second)

	/*opts = b.will.apply(opts)

	 if b.tlsEnabled {
		opts = opts.SetTLSConfig(b.tlsConf)
	}

	opts = opts.SetUsername(b.username)
	opts = opts.SetPassword(b.password)

	for _, u := range b.urls {
		opts = opts.AddBroker(u.String())
	} */
	opts = opts.AddBroker(b.url)
	return opts
}
func getAllTopics(inputTopic inputTopic, mergeTopics []mergeTopics, qos uint8) map[string]byte {
	topics := make(map[string]byte) // add input topics to the list of available topics
	for _, v := range mergeTopics {
		for _, h := range v.topics {
			topics[h.topic] = qos // set filter bit to true for each topic in merge group
		}

	}
	topics[inputTopic.topic] = qos // set filter bit to true for each topic in merge group

	return topics
}
func getElementsParsed(topics []mergeTopics) map[string]any {
	d := make(map[string]any)
	for _, g := range topics {
		for _, h := range g.topics {
			if h.oldValue != nil {
				var dd map[string]any
				err := json.Unmarshal(h.oldValue, &dd)
				if err != nil {
					log.Panic(err)
				}
				for _, g := range h.elements {
					d[g] = dd[g]
				}
			}

		}
	}
	//log.Println("Girish getElementsParsed:", d)
	return d
}

func getServiceMessage(msg mqtt.Message, mergeTopicss []mergeTopics) ([]byte, error) {

	c := make(map[string]any)
	err := json.Unmarshal(msg.Payload(), &c)
	// panic on error
	if err != nil {
		log.Println(err)
		return nil, err
	}
	//log.Println("Girish getServiceMessage:", c)
	parsedFromCurrentMessage := getElementsParsed(mergeTopicss)
	maps.Copy(parsedFromCurrentMessage, c)
	//log.Println("Girish finalmap:", parsedFromCurrentMessage)
	data, err := json.Marshal(parsedFromCurrentMessage)
	if err != nil {
		return nil, err
	}
	return data, nil

}
func compareTopic(topic string, inpputTopic string) bool {
	var itopic, ctopic string
	iLastIndex := strings.LastIndexAny(inpputTopic, "/")
	if iLastIndex != -1 {
		itopic = inpputTopic[:iLastIndex]
	} else {
		itopic = inpputTopic
	}
	cLastIndex := strings.LastIndexAny(topic, "/")
	if cLastIndex != -1 {
		ctopic = topic[:cLastIndex]
	} else {
		ctopic = topic
	}
	return itopic == ctopic
}
func getOutputTopic(topic, outputTopic string) string {
	var itopic string
	iLastIndex := strings.LastIndexAny(topic, "/")
	if iLastIndex != -1 {
		itopic = topic[iLastIndex+1:]
		return fmt.Sprintf("%s/%s", itopic, outputTopic)
	}
	return outputTopic
}

/* func mergerMessage(currentMsg map[string]any, mergerTopic mergeTopics) map[string]any {
	for _, v := range mergerTopic.topics {
		if v.oldValue != nil {
			c := make(map[string]json.RawMessage)
			err := json.Unmarshal(v.oldValue, &c)
			if err != nil {
				log.Panic(err)
			}
			return getElementsFromMergeTag(v, c, currentMsg)

		}
	}
	return currentMsg

} */
