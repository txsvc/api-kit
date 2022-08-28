package config

import (
	"fmt"
)

func NewAppInfo(name, shortName, copyright, about string, major, minor, fix int) Info {
	return Info{
		name:         name,
		shortName:    shortName,
		copyright:    copyright,
		about:        about,
		majorVersion: major,
		minorVersion: minor,
		fixVersion:   fix,
	}
}

func (i *Info) Name() string {
	return i.name
}

func (i *Info) ShortName() string {
	return i.shortName
}

func (i *Info) Copyright() string {
	return i.copyright
}

func (i *Info) About() string {
	return i.about
}

func (i *Info) MajorVersion() int {
	return i.majorVersion
}

func (i *Info) MinorVersion() int {
	return i.minorVersion
}

func (i *Info) FixVersion() int {
	return i.fixVersion
}

func (i *Info) VersionString() string {
	return fmt.Sprintf("%d.%d.%d", i.majorVersion, i.minorVersion, i.fixVersion)
}

func (i *Info) UserAgentString() string {
	return fmt.Sprintf("%s %d.%d.%d", i.shortName, i.majorVersion, i.minorVersion, i.fixVersion)
}

func (i *Info) ServerString() string {
	return fmt.Sprintf("%s %d.%d.%d", i.shortName, i.majorVersion, i.minorVersion, i.fixVersion)
}
