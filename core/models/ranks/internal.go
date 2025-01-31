package ranks

var Internal = map[string]*Rank{
	"admin": {
		Name:        "admin",
		Description: "network administrator/maintainer",
	},
	"basic": {
		Name:        "basic",
		Description: "basic User",
	},
	"vip": {
		Name:        "vip",
		Description: "VIP User",
	},
	"api": {
		Name:        "api",
		Description: "api access!",
	},
	"cnc": {
		Name:        "cnc",
		Description: "Access to the CnC!",
	},
	"member": {
		Name:        "member",
		Description: "member of our comminuty",
	},
}
