package ui

import "github.com/ngundlach/alsa-profile-switcher/pactl"

type dataFetchedMsg struct {
	cards []pactl.Card
}

type errorMsg struct {
	err error
}

type profileChangedMsg struct{}
