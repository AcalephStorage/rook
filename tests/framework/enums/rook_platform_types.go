/*
Copyright 2016 The Rook Authors. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package enums

import (
	"fmt"
	"strings"
)

//RookPlatformType is enum if rook platforms
type RookPlatformType int

const (
	//Kubernetes plafrom is default
	Kubernetes RookPlatformType = iota + 1
	//BareMetal plafrom
	BareMetal
	//StandAlone plafrom is default
	StandAlone
	//None platform
	None
)

var platforms = [...]string{
	"Kubernetes",
	"BareMetal",
	"StandAlone",
	"None",
}

func (platform RookPlatformType) String() string {
	return platforms[platform-1]
}

//GetRookPlatFormTypeFromString returns rook platform type
func GetRookPlatFormTypeFromString(name string) (RookPlatformType, error) {
	switch {
	case strings.EqualFold(name, Kubernetes.String()):
		return Kubernetes, nil
	case strings.EqualFold(name, BareMetal.String()):
		return BareMetal, nil
	case strings.EqualFold(name, StandAlone.String()):
		return StandAlone, nil
	case strings.EqualFold(name, None.String()):
		return None, nil
	default:
		return None, fmt.Errorf("Unsupported Rook Platform Type: " + name)
	}
}
