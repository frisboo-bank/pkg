package simpledi

import "maps"

type tracked map[string]int

func (s tracked) add(info dependencyInfo) tracked {
	newTracked := make(tracked, len(s)+1)

	maps.Copy(newTracked, s)

	newTracked[info.key] = len(s)

	return newTracked
}

func (s tracked) ordered() []string {
	keys := make([]string, len(s))

	for key, idx := range s {
		keys[idx] = key
	}

	return keys
}
