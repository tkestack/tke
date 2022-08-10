package policyrolecache

import (
	"fmt"
	"k8s.io/apimachinery/pkg/util/sets"
	"reflect"
	"sync"
	apiauthzv1 "tkestack.io/tke/api/authz/v1"
)

type PolicyRoleCache interface {
	GetRolesByPolicy(policyName string) sets.String
	PutByRole(role *apiauthzv1.Role)
	DeleteRole(role *apiauthzv1.Role)
}

type policyRoleCache struct {
	rw sync.RWMutex
	// key: policyName
	// value: roleName set
	store map[string]sets.String
}

var Cache = &policyRoleCache{store: map[string]sets.String{}}

func (c *policyRoleCache) GetRolesByPolicy(policyName string) sets.String {
	c.rw.RLocker().Lock()
	defer c.rw.RLocker().Unlock()
	return c.store[policyName]
}

func (c *policyRoleCache) UpdateByRole(old, new *apiauthzv1.Role) {
	if reflect.DeepEqual(old.Policies, new.Policies) {
		return
	}
	c.rw.Lock()
	defer c.rw.Unlock()
	roleName := fmt.Sprintf("%s/%s", old.Namespace, old.Name)
	for _, oldPolicy := range old.Policies {
		delete(c.store[oldPolicy], roleName)
	}
	for _, newPolicy := range new.Policies {
		set := c.store[newPolicy]
		if set == nil {
			set = sets.String{}
		}
		set[roleName] = sets.Empty{}
		c.store[newPolicy] = set
	}
}

func (c *policyRoleCache) PutByRole(role *apiauthzv1.Role) {
	c.rw.Lock()
	defer c.rw.Unlock()
	roleName := fmt.Sprintf("%s/%s", role.Namespace, role.Name)
	policies := role.Policies
	for _, policy := range policies {
		set := c.store[policy]
		if set == nil {
			set = sets.String{}
		}
		set[roleName] = sets.Empty{}
		c.store[policy] = set
	}
}

func (c *policyRoleCache) DeleteRole(role *apiauthzv1.Role) {
	c.rw.Lock()
	defer c.rw.Unlock()
	roleName := fmt.Sprintf("%s/%s", role.Namespace, role.Name)
	for _, policy := range role.Policies {
		delete(c.store[policy], roleName)
	}
}

func (c *policyRoleCache) DeletePolicy(policyName string) {
	c.rw.Lock()
	defer c.rw.Unlock()
	delete(c.store, policyName)
}
