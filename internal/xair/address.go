package xair

var xairAddressMap = map[string]string{
	"main":     "/lr",
	"strip":    "/ch/%02d",
	"bus":      "/bus/%01d",
	"headamp":  "/headamp/%02d",
	"snapshot": "/-snap",
}

var x32AddressMap = map[string]string{
	"main":     "/main/st",
	"mainmono": "/main/mono",
	"strip":    "/ch/%02d",
	"bus":      "/bus/%02d",
	"headamp":  "/headamp/%02d",
	"snapshot": "/-snap",
}

func addressMapForMixerKind(kind MixerKind) map[string]string {
	switch kind {
	case KindX32:
		return x32AddressMap
	default:
		return xairAddressMap
	}
}
