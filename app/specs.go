package app

type Specs struct {
	version   string
	channel   string
	buildDate string
}

func NewSpecs(
	version string,
	channel string,
	buildDate string,
) *Specs {
	return &Specs{
		version:   version,
		channel:   channel,
		buildDate: buildDate,
	}
}

func (specs *Specs) GetVersion() string {
	return specs.version
}

func (specs *Specs) GetChannel() string {
	return specs.channel
}

func (specs *Specs) GetBuildDate() string {
	return specs.buildDate
}
