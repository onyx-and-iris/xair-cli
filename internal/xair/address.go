package xair

var xairAddressMap = map[string]string{
	"bus": "/bus/%01d",
}

var x32AddressMap = map[string]string{
	"bus": "/bus/%02d",
}

func addressMapForMixerKind(kind MixerKind) map[string]string {
	switch kind {
	case KindX32:
		return x32AddressMap
	default:
		return xairAddressMap
	}
}
