package splunk

import (
	"encoding/json"
	"strings"
)

type ACL struct {
	// The app context for the resource. Required for updating saved search ACL properties.
	// Allowed values are: The name of an app or system
	App string `schema:"app,omitempty" json:"app"`

	// Indicates if the active user can change permissions for this object.
	CanChangePerms bool `schema:"can_change_perms,omitempty" json:"can_change_perms"`

	// Indicates if the active user can change sharing to app level.
	CanShareApp bool `schema:"can_share_app,omitempty" json:"can_share_app"`

	// Indicates if the active user can change sharing to system level
	CanShareGlobal bool `schema:"can_share_global,omitempty" json:"can_share_global"`

	// Indicates if the active user can change sharing to user level.
	CanShareUser bool `schema:"can_share_user,omitempty" json:"can_share_user"`

	// Indicates if the active user can edit this object. Defaults to true.
	CanWrite bool `schema:"can_write,omitempty" json:"can_write"`

	// User name of resource owner. Defaults to the resource creator. Required for updating any knowledge object ACL properties.
	// nobody = All users may access the resource, but write access to the resource might be restricted.
	Owner string `schema:"owner,omitempty" json:"owner"`

	Perms struct {
		// Properties that indicate resource read permissions.
		Read []string `schema:"perms.read,omitempty" json:"read"`

		// Properties that indicate write permissions of the resource.
		Write []string `schema:"perms.write,omitempty" json:"write"`
	} `json:"perms"`

	// Indicates whether an admin or user with sufficient permissions can delete the entity.
	Removable bool `schema:"removable,omitempty" json:"removable"`

	// Indicates how the resource is shared. Required for updating any knowledge object ACL properties.
	// app: Shared within a specific app
	// global: (Default) Shared globally to all apps.
	// user: Private to a user
	Sharing string `schema:"sharing,omitempty" json:"sharing"`
}

type ACLFeed struct {
	Feed
	Entry []ACLEntry `schema:"-" json:"entry"`
}

type ACLEntry struct {
	Entry
	ACL ACL `json:"acl"`
}

func (c *Client) ACLPost(acl *ACL, path string) (f ACLFeed, e error) {
	params, e := encode(acl)
	if e != nil {
		return
	}

	// perm.read and perm.write needs to be a single comma delimited value
	if p, ok := params["perms.read"]; ok {
		params["perms.read"] = []string{strings.Join(p, ",")}
	}
	if p, ok := params["perms.write"]; ok {
		params["perms.write"] = []string{strings.Join(p, ",")}
	}

	b, e := c.Post(path, params)
	if e != nil {
		return
	}

	json.Unmarshal(b, &f)
	return
}