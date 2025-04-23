package apis

import "api/core/models/floods"

type ApiMethod struct {
	Name string
}

func WireFloodMethods() {
    for _, api := range Apis {
        api.Methods = make(map[string]string, len(floods.Methods))
        for name, fm := range floods.Methods {
            api.Methods[name] = fm.DisplayName 
        }
    }
}