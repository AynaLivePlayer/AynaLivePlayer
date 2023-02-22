package model

import "fmt"

type VersionInfo struct {
	Version Version
	Info    string
}

type Version uint32

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", (v>>16)&0xff, (v>>8)&0xff, v&0xff)
}

func (v Version) Major() uint8 {
	return uint8((v >> 16) & 0xff)
}

func (v Version) Minor() uint8 {
	return uint8((v >> 8) & 0xff)
}

func (v Version) Patch() uint8 {
	return uint8(v & 0xff)
}

func VersionFromString(s string) Version {
	var major, minor, patch uint8
	_, err := fmt.Sscanf(s, "%d.%d.%d", &major, &minor, &patch)
	if err != nil {
		return 0
	}
	return Version(major)<<16 | Version(minor)<<8 | Version(patch)
}
