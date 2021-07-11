package data

import "errors"

type Platform struct {
	Name string
}

var platforms = []Platform{
	{Name: "PC"}, {Name: "PS5"}, {Name: "SNES"},
}

func CreatePlatform(platform Platform) {
	platforms = append(platforms, platform)
}

func GetPlatforms() []Platform {
	return platforms
}

func GetPlatformByName(name string) (Platform, error) {
	var platform Platform
	for _, p := range platforms {
		if p.Name == name {
			platform = p
			break
		}
	}
	if platform == (Platform{}) {
		return Platform{}, errors.New("No platform with the specified name.")
	}
	return platform, nil
}
