package database

import (
	rranks "api/core/models/ranks"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
)

func (user *User) Sync() error {
	// Decode the base64 string into the ranks
	ranks, err := base64.RawStdEncoding.DecodeString(user.ranks)
	if err != nil {
		log.Println("failed to sync users ranks", err)
		return err
	}

	// Parse the decoded ranks into a list of permissions
	var permissions []*rranks.Rank
	err = json.Unmarshal(ranks, &permissions)
	if err != nil {
		log.Println("failed to unmarshal ranks", err)
		return err
	}

	// Initialize user.Ranks as a slice of *rranks.Rank
	var rankObjects []*rranks.Rank
	for _, permission := range permissions {
		if rank, ok := rranks.Internal[permission.Name]; ok {
			rankObjects = append(rankObjects, rank)
		}
	}

	user.Ranks = rankObjects
	return nil
}

func (user *User) UpdateRoles(name string, role *rranks.Rank, has bool) error {
	if user.HasPermission(name) && has {
		return nil
	}
	if user.HasPermission(name) && !has {
		for _, rank := range user.Ranks {
			if rank.Name == name {
				rank.Has = false
			}
		}
	} else if has {
		user.Ranks = append(user.Ranks, &rranks.Rank{
			Name:        name,
			Description: role.Description,
			Has:         has,
		})
	}
	r, err := json.Marshal(user.Ranks)
	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Println(string(r))
	user.ranks = base64.RawStdEncoding.EncodeToString(r)
	return Container.UpdateUser(user)
}

func (user *User) HasPermission(role string) bool {
	for _, rank := range user.Ranks {
		if rank.Name == role && rank.Has {
			return true
		}
	}
	return false
}
func (user *User) HasAdminRole() bool {
	return true
}

func (user *User) HasRoles(roles []string) bool {
	if len(roles) == 1 && roles[0] == "" || len(roles) == 0 {
		return true
	}
	var has bool
	for _, role := range roles {
		for _, rank := range user.Ranks {
			if rank.Name == role && rank.Has {
				has = true
			}
		}
	}
	return has
}

func (user *User) NewRoles() string {
	r, err := json.Marshal(user.Ranks)
	if err != nil {
		log.Println(err)
		return ""
	}
	fmt.Println(string(r))
	user.ranks = base64.RawStdEncoding.EncodeToString(r)
	return user.ranks
}
