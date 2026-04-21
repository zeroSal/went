package app

type BuildSpecs struct {
	version   string
	channel   string
	buildDate string
}

func NewBuildSpecs(
	version string,
	channel string,
	buildDate string,
) *BuildSpecs {
	return &BuildSpecs{
		version:   version,
		channel:   channel,
		buildDate: buildDate,
	}
}

func (specs *BuildSpecs) GetVersion() string {
	return specs.version
}

func (specs *BuildSpecs) GetChannel() string {
	return specs.channel
}

func (specs *BuildSpecs) GetBuildDate() string {
	return specs.buildDate
}
