package configuration

type Main struct {
	Title		string
	UiType		string
}

type Network struct {
	Ip			string
	LoginPort	int
	PlusSalt	byte
	XorSalt		byte
	PreferIndex	uint16
}

type Features struct {
	UseUpdate	bool
	UseChisel	bool
	UseDefence	bool
}

type Config struct {
	Id			int
	Main		Main
	Network		Network
	Features	Features
}
