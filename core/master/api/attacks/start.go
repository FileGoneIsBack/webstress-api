package attackapi

import (
	"api/core/models/apis"
	"api/core/models/floods"
	"api/core/models/servers"
	"fmt"
	"strings"
)


func SendAttack(conns int, atk *floods.Attack) (string, error) {
    srvCount, srvErr := servers.DistributeServers(atk, conns)
    apiCount, apiErr := apis.Send(atk, conns)

    partsOK, partsErr := []string{}, []string{}
    if srvErr == nil { partsOK = append(partsOK, fmt.Sprintf("servers:%d", srvCount)) }
    if apiErr == nil { partsOK = append(partsOK, fmt.Sprintf("APIs:%d", apiCount)) }
    if srvErr != nil { partsErr = append(partsErr, fmt.Sprintf("servers✗:%v", srvErr)) }
    if apiErr != nil { partsErr = append(partsErr, fmt.Sprintf("APIs✗:%v", apiErr)) }

    if len(partsOK) > 0 {
        return "Success: " + strings.Join(partsOK, " & "), nil
    }
    return "", fmt.Errorf("all failed: %s", strings.Join(partsErr, "; "))
}