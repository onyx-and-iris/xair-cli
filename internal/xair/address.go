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
	"mainmono": "/main/m",
	"matrix":   "/mtx/%02d",
	"strip":    "/ch/%02d",
	"bus":      "/bus/%02d",
	"headamp":  "/headamp/%03d",
	"snapshot": "/-snap",
}

func addressMapFromMixerKind(kind mixerKind) map[string]string {
	switch kind {
	case kindX32:
		return x32AddressMap
	default:
		return xairAddressMap
	}
}
