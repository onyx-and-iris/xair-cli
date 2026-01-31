package xair

type MixerKind string

const (
	KindXAir MixerKind = "xair"
	KindX32  MixerKind = "x32"
)

func NewMixerKind(kind string) MixerKind {
	switch kind {
	case "xair":
		return KindXAir
	case "x32":
		return KindX32
	default:
		return KindXAir
	}
}
