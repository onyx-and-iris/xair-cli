package xair

type Option func(*engine)

func WithKind(kind string) Option {
	return func(e *engine) {
		e.Kind = MixerKind(kind)
		e.addressMap = addressMapForMixerKind(e.Kind)
	}
}
