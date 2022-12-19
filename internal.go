package gomicro

import (
	"github.com/civet148/log"
	"strings"
)

func parseRegistry(strRegistry string) (typ RegistryType, addresses []string) {
	//--registry "etcd://127.0.0.1:2379,127.0.0.1:2380"
	var strAddress string
	ss := strings.Split(strRegistry, "://")
	count := len(ss)
	if count > 1 {
		strRegName := strings.ToLower(ss[0])
		strAddress = strings.ToLower(ss[1])
		switch strRegName {
		case "etcd":
			typ = RegistryType_ETCD
		default:
			log.Warnf("Unknown registry name [%s], use default MDNS", strRegName)
			typ = RegistryType_MDNS
		}
	} else {
		typ = RegistryType_MDNS
	}

	log.Infof("registry type [%s] address %+v", typ.String(), strAddress)
	return typ, strings.Split(strAddress, ",")
}