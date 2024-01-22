package recycle

type Geocode struct {
	CZ int
	PK int
	KV int
	VY int
}

func getGeocodes() Geocode {
	return Geocode{
		CZ: 3600051684,
		PK: 3600442466,
		KV: 3600442314,
		VY: 3600442453,
	}

}
