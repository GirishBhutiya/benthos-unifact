package s7comm_plugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/robinson/gos7"
)

const addressRegexp = `^(?P<area>[A-Z]+)(?P<no>[0-9]+)\.(?P<type>[A-Z]+)(?P<start>[0-9]+)(?:\.(?P<extra>.*))?$`

var (
	regexAddr = regexp.MustCompile(addressRegexp)
	// Area mapping taken from https://github.com/robinson/gos7/blob/master/client.go
	areaMap = map[string]int{
		"PE": 0x81, // process inputs
		"PA": 0x82, // process outputs
		"MK": 0x83, // Merkers
		"DB": 0x84, // DB
		"C":  0x1C, // counters
		"T":  0x1D, // timers
	}
	// Word-length mapping taken from https://github.com/robinson/gos7/blob/master/client.go
	wordLenMap = map[string]int{
		"X":  0x01, // Bit
		"B":  0x02, // Byte (8 bit)
		"C":  0x03, // Char (8 bit)
		"S":  0x03, // String (8 bit)
		"W":  0x04, // Word (16 bit)
		"I":  0x05, // Integer (16 bit)
		"DW": 0x06, // Double Word (32 bit)
		"DI": 0x07, // Double integer (32 bit)
		"R":  0x08, // IEEE 754 real (32 bit)
		// see https://support.industry.siemens.com/cs/document/36479/date_and_time-format-for-s7-?dti=0&lc=en-DE
		"DT": 0x0F, // Date and time (7 byte)
	}
)

type S7DataItemWithAddressAndConverter struct {
	Address       string
	ConverterFunc converterFunc
	Item          gos7.S7DataItem
}
type subscriptionD struct {
	ID        int
	Address   string
	Name      string
	Group     string
	DB        string
	Historian string
	SqlSp     string
	DataType  string
	Value     any
}

func ParseSubscriptionDef(subscription []string, batchMaxSize int) ([]subscriptionD, [][]S7DataItemWithAddressAndConverter, error) {
	var parsedSubscription []subscriptionD
	addresses := make([]string, 0)
	for _, subscriptionElement := range subscription {
		var subscr map[string][]map[string]string
		var subsc subscriptionD
		err := json.Unmarshal([]byte(subscriptionElement), &subscr)
		if err != nil {
			return nil, nil, err
		}
		//log.Println("Girish ParseSubscription() json tagname: ", subscr)
		for key, values := range subscr {
			for _, obj := range values {

				subsc.ID, _ = strconv.Atoi(key)
				subsc.Address = obj["address"]
				subsc.Group = obj["group"]
				subsc.DB = obj["db"]
				subsc.Historian = obj["historian"]
				subsc.SqlSp = obj["sqlSp"]
				subsc.DataType = obj["datatype"]
				subsc.Name = obj["name"]

			}
		}
		addresses = append(addresses, subsc.Address)
		parsedSubscription = append(parsedSubscription, subsc)
	}
	batches, err := parseAddresses(addresses, batchMaxSize)
	if err != nil {
		return parsedSubscription, nil, err
	}
	return parsedSubscription, batches, nil

}

func parseAddresses(addresses []string, batchMaxSize int) ([][]S7DataItemWithAddressAndConverter, error) {
	parsedAddresses := make([]S7DataItemWithAddressAndConverter, 0, len(addresses))

	for _, address := range addresses {
		item, converterFunc, err := handleFieldAddress(address)
		if err != nil {
			return nil, fmt.Errorf("address %q: %w", address, err)
		}

		newS7DataItemWithAddressAndConverter := S7DataItemWithAddressAndConverter{
			Address:       address,
			ConverterFunc: converterFunc,
			Item:          *item,
		}

		parsedAddresses = append(parsedAddresses, newS7DataItemWithAddressAndConverter)
	}

	// check for duplicates

	for i, a := range parsedAddresses {
		for j, b := range parsedAddresses {
			if i == j {
				continue
			}
			if a.Item.Area == b.Item.Area && a.Item.DBNumber == b.Item.DBNumber && a.Item.Start == b.Item.Start {
				return nil, fmt.Errorf("duplicate address %v", a)
			}
		}
	}

	// Now split the addresses into batches based on the batchMaxSize
	batches := make([][]S7DataItemWithAddressAndConverter, 0)
	for i := 0; i < len(parsedAddresses); i += batchMaxSize {
		end := i + batchMaxSize
		if end > len(parsedAddresses) {
			end = len(parsedAddresses)
		}
		batches = append(batches, parsedAddresses[i:end])
	}

	return batches, nil
}
func handleFieldAddress(address string) (*gos7.S7DataItem, converterFunc, error) {
	// Parse the address into the different parts
	if !regexAddr.MatchString(address) {
		return nil, nil, fmt.Errorf("invalid address %q", address)
	}
	names := regexAddr.SubexpNames()[1:]
	parts := regexAddr.FindStringSubmatch(address)[1:]
	if len(names) != len(parts) {
		return nil, nil, fmt.Errorf("names %v do not match parts %v", names, parts)
	}
	groups := make(map[string]string, len(names))
	for i, n := range names {
		groups[n] = parts[i]
	}

	// Check that we do have the required entries in the address
	if _, found := groups["area"]; !found {
		return nil, nil, errors.New("area is missing from address")
	}

	if _, found := groups["no"]; !found {
		return nil, nil, errors.New("area index is missing from address")
	}
	if _, found := groups["type"]; !found {
		return nil, nil, errors.New("type is missing from address")
	}
	if _, found := groups["start"]; !found {
		return nil, nil, errors.New("start address is missing from address")
	}
	dtype := groups["type"]

	// Lookup the item values from names and check the params
	area, found := areaMap[groups["area"]]
	if !found {
		return nil, nil, errors.New("invalid area")
	}
	wordlen, found := wordLenMap[dtype]
	if !found {
		return nil, nil, errors.New("unknown data type")
	}
	areaidx, err := strconv.Atoi(groups["no"])
	if err != nil {
		return nil, nil, fmt.Errorf("invalid area index: %w", err)
	}
	start, err := strconv.Atoi(groups["start"])
	if err != nil {
		return nil, nil, fmt.Errorf("invalid start address: %w", err)
	}

	// Check the amount parameter if any
	var extra, bit int
	switch dtype {
	case "S":
		// We require an extra parameter
		x := groups["extra"]
		if x == "" {
			return nil, nil, errors.New("extra parameter required")
		}

		extra, err = strconv.Atoi(x)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid extra parameter: %w", err)
		}
		if extra < 1 {
			return nil, nil, fmt.Errorf("invalid extra parameter %d", extra)
		}
	case "X":
		// We require an extra parameter
		x := groups["extra"]
		if x == "" {
			return nil, nil, errors.New("extra parameter required")
		}

		bit, err = strconv.Atoi(x)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid extra parameter: %w", err)
		}
		if bit < 0 || bit > 7 {
			// Ensure bit address is valid
			return nil, nil, fmt.Errorf("invalid extra parameter: bit address %d out of range", bit)
		}
	default:
		if groups["extra"] != "" {
			return nil, nil, errors.New("extra parameter specified but not used")
		}
	}

	// Get the required buffer size
	amount := 1
	var buflen int
	switch dtype {
	case "X", "B", "C": // 8-bit types
		buflen = 1
	case "W", "I": // 16-bit types
		buflen = 2
	case "DW", "DI", "R": // 32-bit types
		buflen = 4
	case "DT": // 7-byte
		buflen = 7
	case "S":
		amount = extra
		// Extra bytes as the first byte is the max-length of the string and
		// the second byte is the actual length of the string.
		buflen = extra + 2
	default:
		return nil, nil, errors.New("invalid data type")
	}

	// Setup the data item
	item := &gos7.S7DataItem{
		Area:     area,
		WordLen:  wordlen,
		Bit:      bit,
		DBNumber: areaidx,
		Start:    start,
		Amount:   amount,
		Data:     make([]byte, buflen),
	}

	// Determine the type converter function
	f := determineConversion(dtype)
	return item, f, nil
}
