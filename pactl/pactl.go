// Package pactl provides functionality to interact with the pactl comandline utility
package pactl

import (
	"encoding/json"
	"fmt"
	"maps"
	"os/exec"
	"slices"
)

const pactlExec = "pactl"

func (c *Card) print() {
	fmt.Println("Name: ", c.Name)
	fmt.Println("Driver: ", c.Driver)
	fmt.Println("Device Description: ", c.Properties.DeviceDescription)
	fmt.Println("Device Nick: ", c.Properties.DeviceNick)
	fmt.Println("Profiles: ")
	profileKeys := SortKeys(c.Profiles)
	for _, k := range profileKeys {
		fmt.Print("	", k, ": ")
		fmt.Println(c.Profiles[k].
			Description, "| available: ", c.Profiles[k].Available)
	}
	fmt.Println("Active profile: ", c.ActiveProfile)
	fmt.Println("---------------")
}

func FetchDeviceData() ([]Card, error) {
	jsonOutput, err := runCmd()
	if err != nil {
		return nil, err
	}
	cards, err := parseCards(jsonOutput)
	if err != nil {
		return nil, err
	}
	return cards, nil
}

func SortKeys[T any](unorderedMap map[string]T) []string {
	keys := slices.Sorted(maps.Keys(unorderedMap))
	return keys
}

func parseCards(jsonOutput []byte) ([]Card, error) {
	var cards []Card
	err := json.Unmarshal(jsonOutput, &cards)
	if err != nil {
		return nil, err
	}
	return cards, nil
}

func runCmd() ([]byte, error) {
	cmd := exec.Command(pactlExec, "-f", "json", "list", "cards")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return output, nil
}

func SetActiveProfile(card string, profile string) error {
	cmd := exec.Command(pactlExec, "set-card-profile", card, profile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, output)
	}
	return nil
}
