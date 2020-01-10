package auth_test

import (
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	authv1 "tkestack.io/tke/api/auth/v1"
	"tkestack.io/tke/pkg/util/names"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	tkeclientset "tkestack.io/tke/api/client/clientset/versioned"
	"tkestack.io/tke/pkg/auth/util"
	apiclient "tkestack.io/tke/test/util/tkeclient"
)

var _ = Describe("E2e", func() {
	var client tkeclientset.Interface
	var groupID string
	var userID string
	var policyID string
	var roleID string
	var err error
	BeforeEach(func() {
		client, err = apiclient.LoadOrSetupTKE()
		Expect(err).To(BeNil())
	})

	JustBeforeEach(func() {
		group := authv1.LocalGroup{
			Spec: authv1.LocalGroupSpec{
				DisplayName: "E2ETestLocalGroup",
				TenantID:    "default",
				Description: "Test",
			},
		}

		result, err := client.AuthV1().LocalGroups().Create(&group)
		Expect(err).To(BeNil())
		groupID = result.Name

		user := authv1.LocalIdentity{
			Spec: authv1.LocalIdentitySpec{
				DisplayName:    "E2ETestLocalIdentity",
				TenantID:       "default",
				Username:       names.Generator.GenerateName("e2etest"),
				HashedPassword: "MTIzNDU=",
			},
		}

		localidentity, err := client.AuthV1().LocalIdentities().Create(&user)
		Expect(err).To(BeNil())
		userID = localidentity.Name

		policy := authv1.Policy{Spec: authv1.PolicySpec{
			DisplayName: "E2ETestPolicy",
			TenantID:    "default",
			Category:    "common",
			Username:    "admin",
			Description: "Test",
			Statement: authv1.Statement{
				Actions:   []string{"list"},
				Effect:    authv1.Allow,
				Resources: []string{"deployment"},
			},
		}}

		pol, err := client.AuthV1().Policies().Create(&policy)
		Expect(err).To(BeNil())
		policyID = pol.Name

		role := authv1.Role{Spec: authv1.RoleSpec{
			DisplayName: "E2ETestRole",
			TenantID:    "default",
			Username:    "admin",
			Description: "Test",
			Policies:    []string{policyID},
		}}

		rol, err := client.AuthV1().Roles().Create(&role)
		Expect(err).To(BeNil())
		roleID = rol.Name
	})

	AfterEach(func() {
		_ = client.AuthV1().LocalIdentities().Delete(userID, &metav1.DeleteOptions{})

		err = wait.Poll(1*time.Second, 10*time.Second, func() (bool, error) {
			_, err = client.AuthV1().LocalIdentities().Get(userID, metav1.GetOptions{})
			if err != nil && errors.IsNotFound(err) {
				return true, nil
			}

			if err != nil {
				return false, err
			}
			return false, nil
		})
		Expect(err).To(BeNil())

		_ = client.AuthV1().LocalGroups().Delete(groupID, &metav1.DeleteOptions{})
		err = wait.Poll(1*time.Second, 10*time.Second, func() (bool, error) {
			_, err = client.AuthV1().LocalGroups().Get(groupID, metav1.GetOptions{})
			if err != nil && errors.IsNotFound(err) {
				return true, nil
			}

			if err != nil {
				return false, err
			}
			return false, nil
		})

		Expect(err).To(BeNil())

		_ = client.AuthV1().Policies().Delete(policyID, &metav1.DeleteOptions{})
		err = wait.Poll(1*time.Second, 10*time.Second, func() (bool, error) {
			_, err = client.AuthV1().Policies().Get(policyID, metav1.GetOptions{})
			if err != nil && errors.IsNotFound(err) {
				return true, nil
			}

			if err != nil {
				return false, err
			}
			return false, nil
		})

		Expect(err).To(BeNil())
	})

	It("test group binding", func() {

		localIdentity, err := client.AuthV1().LocalIdentities().Get(userID, metav1.GetOptions{})
		Expect(err).To(BeNil())

		_, err = client.AuthV1().LocalGroups().Get(groupID, metav1.GetOptions{})
		Expect(err).To(BeNil())

		var subjects []authv1.Subject
		subjects = append(subjects, authv1.Subject{ID: userID})
		binding := authv1.Binding{Users: subjects}
		localGroup := &authv1.LocalGroup{}
		err = client.AuthV1().RESTClient().Post().Resource("localgroups").SubResource("binding").Name(groupID).Body(&binding).Do().Into(localGroup)
		Expect(err).To(BeNil())

		found := false
		for _, sub := range localGroup.Status.Users {
			if sub.ID == userID && sub.Name == localIdentity.Spec.Username {
				found = true
			}
		}

		Expect(found).To(BeTrue())

		By(fmt.Sprintf("wait bind user %s to group %s", userID, groupID))
		found = false
		err = wait.Poll(1*time.Second, 30*time.Second, func() (bool, error) {
			var groupList = &authv1.GroupList{}
			err := client.AuthV1().RESTClient().Get().Resource("localidentities").SubResource("groups").Name(userID).Do().Into(groupList)
			if err != nil {
				return false, err
			}

			for _, grp := range groupList.Items {
				if grp.Name == groupID && grp.Spec.DisplayName == localGroup.Spec.DisplayName {
					found = true
					return true, nil
				}
			}
			return false, nil
		})
		Expect(err).To(BeNil())

		Expect(found).To(BeTrue())
	})

	It("test policy binding", func() {

		localIdentity, err := client.AuthV1().LocalIdentities().Get(userID, metav1.GetOptions{})
		Expect(err).To(BeNil())

		group, err := client.AuthV1().LocalGroups().Get(groupID, metav1.GetOptions{})
		Expect(err).To(BeNil())

		_, err = client.AuthV1().Policies().Get(policyID, metav1.GetOptions{})
		Expect(err).To(BeNil())

		binding := authv1.Binding{Users: []authv1.Subject{{ID: userID}}, Groups: []authv1.Subject{{ID: groupID}}}
		policy := &authv1.Policy{}
		err = client.AuthV1().RESTClient().Post().Resource("policies").SubResource("binding").Name(policyID).Body(&binding).Do().Into(policy)
		Expect(err).To(BeNil())

		found := false
		for _, sub := range policy.Status.Users {
			if sub.ID == userID && sub.Name == localIdentity.Spec.Username {
				found = true
			}
		}

		Expect(found).To(BeTrue())

		found = false
		for _, sub := range policy.Status.Groups {
			if sub.ID == groupID && sub.Name == group.Spec.DisplayName {
				found = true
			}
		}

		Expect(found).To(BeTrue())

		By(fmt.Sprintf("wait bind user %s to policy %s", userID, policyID))
		found = false
		err = wait.Poll(1*time.Second, 30*time.Second, func() (bool, error) {
			var policyList = &authv1.PolicyList{}
			err := client.AuthV1().RESTClient().Get().Resource("localidentities").SubResource("policies").Name(userID).Do().Into(policyList)
			if err != nil {
				return false, err
			}

			for _, pol := range policyList.Items {
				if pol.Name == policyID && pol.Spec.DisplayName == policy.Spec.DisplayName {
					found = true
					return true, nil
				}
			}
			return false, nil
		})
		Expect(err).To(BeNil())
		Expect(found).To(BeTrue())

		found = false
		err = wait.Poll(1*time.Second, 30*time.Second, func() (bool, error) {
			var policyList = &authv1.PolicyList{}
			err := client.AuthV1().RESTClient().Get().Resource("users").SubResource("policies").Name(util.CombineTenantAndName("default", userID)).Do().Into(policyList)
			if err != nil {
				return false, err
			}

			for _, pol := range policyList.Items {
				if pol.Name == policyID && pol.Spec.DisplayName == policy.Spec.DisplayName {
					found = true
					return true, nil
				}
			}
			return false, nil
		})
		Expect(err).To(BeNil())
		Expect(found).To(BeTrue())

		By(fmt.Sprintf("wait bind group %s to policy %s", groupID, policyID))
		found = false
		err = wait.Poll(1*time.Second, 30*time.Second, func() (bool, error) {
			var policyList = &authv1.PolicyList{}
			err := client.AuthV1().RESTClient().Get().Resource("localgroups").SubResource("policies").Name(groupID).Do().Into(policyList)
			if err != nil {
				return false, err
			}

			for _, pol := range policyList.Items {
				if pol.Name == policyID && pol.Spec.DisplayName == policy.Spec.DisplayName {
					found = true
					return true, nil
				}
			}
			return false, nil
		})
		Expect(err).To(BeNil())
		Expect(found).To(BeTrue())

		found = false
		err = wait.Poll(1*time.Second, 30*time.Second, func() (bool, error) {
			var policyList = &authv1.PolicyList{}
			err := client.AuthV1().RESTClient().Get().Resource("groups").SubResource("policies").Name(util.CombineTenantAndName("default", groupID)).Do().Into(policyList)
			if err != nil {
				return false, err
			}

			for _, pol := range policyList.Items {
				if pol.Name == policyID && pol.Spec.DisplayName == policy.Spec.DisplayName {
					found = true
					return true, nil
				}
			}
			return false, nil
		})
		Expect(err).To(BeNil())
		Expect(found).To(BeTrue())
	})

	It("test role binding", func() {

		localIdentity, err := client.AuthV1().LocalIdentities().Get(userID, metav1.GetOptions{})
		Expect(err).To(BeNil())

		group, err := client.AuthV1().LocalGroups().Get(groupID, metav1.GetOptions{})
		Expect(err).To(BeNil())

		_, err = client.AuthV1().Policies().Get(policyID, metav1.GetOptions{})
		Expect(err).To(BeNil())

		role, err := client.AuthV1().Roles().Get(roleID, metav1.GetOptions{})
		Expect(err).To(BeNil())

		binding := authv1.Binding{Users: []authv1.Subject{{ID: userID}}, Groups: []authv1.Subject{{ID: groupID}}}
		role = &authv1.Role{}
		err = client.AuthV1().RESTClient().Post().Resource("roles").SubResource("binding").Name(roleID).Body(&binding).Do().Into(role)
		Expect(err).To(BeNil())

		found := false
		for _, sub := range role.Status.Users {
			if sub.ID == userID && sub.Name == localIdentity.Spec.Username {
				found = true
			}
		}

		Expect(found).To(BeTrue())

		found = false
		for _, sub := range role.Status.Groups {
			if sub.ID == groupID && sub.Name == group.Spec.DisplayName {
				found = true
			}
		}

		Expect(found).To(BeTrue())

		By(fmt.Sprintf("wait bind user %s to role %s", userID, policyID))
		found = false
		err = wait.Poll(1*time.Second, 30*time.Second, func() (bool, error) {
			var roleList = &authv1.RoleList{}
			err := client.AuthV1().RESTClient().Get().Resource("localidentities").SubResource("roles").Name(userID).Do().Into(roleList)
			if err != nil {
				return false, err
			}

			for _, rol := range roleList.Items {
				if rol.Name == roleID && rol.Spec.DisplayName == role.Spec.DisplayName {
					found = true
					return true, nil
				}
			}
			return false, nil
		})
		Expect(err).To(BeNil())
		Expect(found).To(BeTrue())

		found = false
		err = wait.Poll(1*time.Second, 30*time.Second, func() (bool, error) {
			var roleList = &authv1.RoleList{}
			err := client.AuthV1().RESTClient().Get().Resource("users").SubResource("roles").Name(util.CombineTenantAndName("default", userID)).Do().Into(roleList)
			if err != nil {
				return false, err
			}

			for _, rol := range roleList.Items {
				if rol.Name == roleID && rol.Spec.DisplayName == role.Spec.DisplayName {
					found = true
					return true, nil
				}
			}
			return false, nil
		})
		Expect(err).To(BeNil())
		Expect(found).To(BeTrue())

		By(fmt.Sprintf("wait bind group %s to role %s", groupID, policyID))
		found = false
		err = wait.Poll(1*time.Second, 30*time.Second, func() (bool, error) {
			var roleList = &authv1.RoleList{}
			err := client.AuthV1().RESTClient().Get().Resource("localgroups").SubResource("roles").Name(groupID).Do().Into(roleList)
			if err != nil {
				return false, err
			}

			for _, rol := range roleList.Items {
				if rol.Name == roleID && rol.Spec.DisplayName == role.Spec.DisplayName {
					found = true
					return true, nil
				}
			}
			return false, nil
		})
		Expect(err).To(BeNil())
		Expect(found).To(BeTrue())

		found = false
		err = wait.Poll(1*time.Second, 30*time.Second, func() (bool, error) {
			var roleList = &authv1.RoleList{}
			err := client.AuthV1().RESTClient().Get().Resource("groups").SubResource("roles").Name(util.CombineTenantAndName("default", groupID)).Do().Into(roleList)
			if err != nil {
				return false, err
			}

			for _, rol := range roleList.Items {
				if rol.Name == roleID && rol.Spec.DisplayName == role.Spec.DisplayName {
					found = true
					return true, nil
				}
			}
			return false, nil
		})
		Expect(err).To(BeNil())
		Expect(found).To(BeTrue())
	})

	It("test role policybinding", func() {
		policy, err := client.AuthV1().Policies().Get(policyID, metav1.GetOptions{})
		Expect(err).To(BeNil())

		_, err = client.AuthV1().Roles().Get(roleID, metav1.GetOptions{})
		Expect(err).To(BeNil())

		binding := authv1.PolicyBinding{Policies: []string{policyID}}
		role := &authv1.Role{}
		err = client.AuthV1().RESTClient().Post().Resource("roles").SubResource("policybinding").Name(roleID).Body(&binding).Do().Into(role)
		Expect(err).To(BeNil())

		found := false
		for _, pol := range role.Spec.Policies {
			if pol == policyID {
				found = true
			}
		}

		Expect(found).To(BeTrue())

		By(fmt.Sprintf("wait bind policies %s to role %s", policyID, policyID))
		found = false
		err = wait.Poll(1*time.Second, 30*time.Second, func() (bool, error) {
			var policyList = &authv1.PolicyList{}
			err := client.AuthV1().RESTClient().Get().Resource("roles").SubResource("policies").Name(roleID).Do().Into(policyList)
			if err != nil {
				return false, err
			}

			for _, pol := range policyList.Items {
				if pol.Name == policyID && pol.Spec.DisplayName == policy.Spec.DisplayName {
					found = true
					return true, nil
				}
			}
			return false, nil
		})
		Expect(err).To(BeNil())
		Expect(found).To(BeTrue())

		role = &authv1.Role{}
		err = client.AuthV1().RESTClient().Post().Resource("roles").SubResource("policyunbinding").Name(roleID).Body(&binding).Do().Into(role)
		Expect(err).To(BeNil())

		found = false
		for _, pol := range role.Spec.Policies {
			if pol == policyID {
				found = true
			}
		}

		Expect(found).To(BeFalse())

		By(fmt.Sprintf("wait bind policies %s to role %s", policyID, policyID))
		found = false
		err = wait.Poll(1*time.Second, 30*time.Second, func() (bool, error) {
			var policyList = &authv1.PolicyList{}
			err := client.AuthV1().RESTClient().Get().Resource("roles").SubResource("policies").Name(roleID).Do().Into(policyList)
			if err != nil {
				return false, err
			}

			for _, pol := range policyList.Items {
				if pol.Name == policyID && pol.Spec.DisplayName == policy.Spec.DisplayName {
					found = true
					return false, nil
				}
			}

			found = false
			return true, nil
		})
		Expect(err).To(BeNil())
		Expect(found).To(BeFalse())
	})

})
