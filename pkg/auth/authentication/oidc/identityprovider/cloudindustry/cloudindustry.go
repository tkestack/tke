package cloudindustry

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

	"github.com/caryxychen/cloudindustry-sdk-go/client/iam"
	iammodel "github.com/caryxychen/cloudindustry-sdk-go/model/iam"
	"github.com/dexidp/dex/connector"
	dexlog "github.com/dexidp/dex/pkg/log"
	dexserver "github.com/dexidp/dex/server"
	"github.com/pkg/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metainternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"tkestack.io/tke/api/auth"
	"tkestack.io/tke/pkg/apiserver/authentication/authenticator/oidc"
	"tkestack.io/tke/pkg/auth/authentication/oidc/identityprovider"
	"tkestack.io/tke/pkg/auth/util"
	"tkestack.io/tke/pkg/util/log"
)

const (
	ConnectorType = "cloudindustry"
)

var store idpStore

func init() {
	// create dex identity provider for cloudindustry connector.
	dexserver.ConnectorsConfig[ConnectorType] = func() dexserver.ConnectorConfig {
		return new(identityProvider)
	}
	store = idpStore{
		cache: map[string]*identityProvider{},
	}
}

// identityProvider is the third-party idp that support CloudIndustry.
type identityProvider struct {
	tenantID       string
	administrators []string
	config         *SDKConfig
	client         *iam.Client
}

type idpStore struct {
	sync.Mutex
	cache map[string]*identityProvider
}

type SDKConfig struct {
	SecretID       string `json:"secret_id"`
	SecretKey      string `json:"secret_key"`
	Endpoint       string `json:"endpoint"`
	IAMAPIEndpoint string `json:"iam_api_endpoint"`
	IAMAppEndpoint string `json:"iam_api_app_endpoint"`
	Region         string `json:"region"`
	MasterID       string `json:"master_id"`
}

var (
	_ identityprovider.IdentityProvider = &identityProvider{}
	_ identityprovider.UserGetter       = &identityProvider{}
	_ identityprovider.UserLister       = &identityProvider{}
	_ identityprovider.GroupGetter      = &identityProvider{}
	_ identityprovider.GroupLister      = &identityProvider{}
)

func NewCloudIndustryIdentityProvider(tenantID string, administrators []string, config *SDKConfig) (identityprovider.IdentityProvider, error) {
	bytes, err := json.Marshal(config)
	if err != nil {
		log.Warnf("illegal sdk config, err '%v'", err)
		return nil, err
	}

	log.Infof("NewCloudIndustryIdentityProvider, tenantID '%s', administrators '%v', config '%s'", tenantID, administrators, string(bytes))
	cpf := profile.NewClientProfile()
	if len(config.IAMAPIEndpoint) != 0 {
		cpf.HttpProfile.Endpoint = config.IAMAPIEndpoint
	} else {
		cpf.HttpProfile.Endpoint = fmt.Sprintf("iam-api.%s", config.Endpoint)
	}
	client, err := iam.NewClient(common.NewCredential(config.SecretID, config.SecretKey), config.Region, cpf)
	if err != nil {
		log.Warnf("init client failed, err: '%v'", err)
		return nil, err
	}

	store.Mutex.Lock()
	defer store.Mutex.Unlock()
	store.cache[tenantID] = &identityProvider{
		tenantID:       tenantID,
		administrators: append(administrators, tenantID),
		config:         config,
		client:         client,
	}
	return store.cache[tenantID], nil
}

func (i *identityProvider) Open(id string, logger dexlog.Logger) (connector.Connector, error) {
	return &cloudIndustryConnector{tenantID: id}, nil
}

func (i *identityProvider) Store() (*auth.IdentityProvider, error) {
	if i.tenantID == "" {
		return nil, fmt.Errorf("must specify tenantID")
	}
	bytes, err := json.Marshal(i.config)
	if err != nil {
		return nil, fmt.Errorf("mashal cloudindustry config failed: %+v", err)
	}
	return &auth.IdentityProvider{
		ObjectMeta: metav1.ObjectMeta{Name: i.tenantID},
		Spec: auth.IdentityProviderSpec{
			Name:           i.tenantID,
			Type:           ConnectorType,
			Administrators: i.administrators,
			Config:         string(bytes),
		},
	}, nil
}

type cloudIndustryConnector struct {
	tenantID string
}

func (c *cloudIndustryConnector) Prompt() string {
	return "Username"
}

func (c *cloudIndustryConnector) Login(ctx context.Context, scopes connector.Scopes, key, accessToken string) (connector.Identity, bool, error) {
	ident := connector.Identity{}
	if key != "access_token" {
		return ident, false, nil
	}
	store.Mutex.Lock()
	provider, ok := store.cache[c.tenantID]
	if !ok {
		err := errors.Errorf("unexpected error, can't find config for tenant '%s'", c.tenantID)
		return ident, false, err
	}
	store.Mutex.Unlock()

	credential := common.NewCredential(provider.config.SecretID, provider.config.SecretKey)
	credential.Token = accessToken
	cpf := profile.NewClientProfile()
	if len(provider.config.IAMAppEndpoint) != 0 {
		cpf.HttpProfile.Endpoint = provider.config.IAMAppEndpoint
	} else {
		cpf.HttpProfile.Endpoint = fmt.Sprintf("iam-app-api.%s", provider.config.Endpoint)
	}
	cpf.HttpProfile.ReqTimeout = 30

	log.Infof("%v, %v", credential, cpf)
	client, err := iam.NewClient(credential, provider.config.Region, cpf)
	if err != nil {
		log.Warnf("failed to create cloudindustry client, err '%v'", err)
		return ident, false, err
	}

	account, err := client.DescribeAccount(iam.NewDescribeAccountRequest())
	if err != nil {
		log.Warnf("failed to describe account, err '%v'", err)
		return ident, false, err
	}
	if account.Response == nil {
		log.Warnf("got empty accounts")
		return ident, false, nil
	}
	userInfo := account.Response.Account

	ident.UserID = fmt.Sprintf("%d", userInfo.AccountId)
	ident.Username = userInfo.AccountName
	ident.PreferredUsername = userInfo.NickName
	ident.Email = userInfo.ContactMail
	ident.EmailVerified = userInfo.MailBindStatus == 1
	extra := map[string]string{
		oidc.TenantIDKey: c.tenantID,
	}
	ident.ConnectorData, _ = json.Marshal(extra)
	log.Infof("user '%s' login successful", ident.Username)
	return ident, true, nil
}

func (c *cloudIndustryConnector) Refresh(ctx context.Context, s connector.Scopes, identity connector.Identity) (connector.Identity, error) {
	return identity, nil
}

func (i *identityProvider) GetUser(ctx context.Context, name string, options *metav1.GetOptions) (*auth.User, error) {
	accountsRequest := iam.NewDescribeAccountsRequest()
	accountsRequest.AccountIds = []string{name}
	accounts, err := i.client.DescribeAccounts(accountsRequest)
	if err != nil {
		return nil, apierrors.NewInternalError(err)
	}
	if accounts.Response == nil || len(accounts.Response.Accounts) == 0 {
		return nil, apierrors.NewNotFound(auth.Resource("user"), name)
	}
	return &i.usersFromAccounts(accounts.Response.Accounts).Items[0], nil
}

func (i *identityProvider) ListUsers(ctx context.Context, options *metainternal.ListOptions) (*auth.UserList, error) {
	subIdsRequest := iam.NewDescribeSubIdsRequest()
	subIdsRequest.MasterId = i.config.MasterID

	ids, err := i.client.DescribeSubIds(subIdsRequest)
	if err != nil {
		log.Warnf("failed to describe subIds, err '%v'", err)
		return nil, apierrors.NewInternalError(err)
	}
	log.Infof("get subids '%s'", ids)

	if ids.Response == nil || len(ids.Response.SubIds) == 0 {
		return &auth.UserList{}, nil
	}

	accountsRequest := iam.NewDescribeAccountsRequest()
	accountsRequest.AccountIds = ids.Response.SubIds
	if options != nil {
		keyword, limit := util.ParseQueryKeywordAndLimit(options)
		accountsRequest.SearchKey = keyword
		accountsRequest.Limit = strconv.Itoa(limit)
	}
	bytes, _ := json.Marshal(accountsRequest)
	log.Debugf("request '%s'", string(bytes))

	accounts, err := i.client.DescribeAccounts(accountsRequest)
	if err != nil {
		log.Warnf("failed to describe accounts, err '%v'", err)
		return nil, apierrors.NewInternalError(err)
	}
	bytes, _ = json.Marshal(accounts)
	log.Debugf("accounts '%s'", string(bytes))
	if accounts.Response == nil || len(accounts.Response.Accounts) == 0 {
		return &auth.UserList{}, nil
	}
	return i.usersFromAccounts(accounts.Response.Accounts), nil
}

func (i *identityProvider) GetGroup(ctx context.Context, name string, options *metav1.GetOptions) (*auth.Group, error) {
	return nil, apierrors.NewNotFound(auth.Resource("group"), name)
}

func (i *identityProvider) ListGroups(ctx context.Context, options *metainternal.ListOptions) (*auth.GroupList, error) {
	return &auth.GroupList{}, nil
}

func (i *identityProvider) usersFromAccounts(accounts []iammodel.Account) *auth.UserList {
	result := &auth.UserList{}
	for _, account := range accounts {
		uid := fmt.Sprintf("%d", account.AccountId)
		extra := map[string]string{}
		if uid == i.config.MasterID {
			extra["administrator"] = "true"
		}
		result.Items = append(result.Items, auth.User{
			ObjectMeta: metav1.ObjectMeta{
				Name: uid,
			},
			Spec: auth.UserSpec{
				ID:          uid,
				Name:        account.AccountName,
				DisplayName: account.NickName,
				Email:       account.ContactMail,
				PhoneNumber: account.ContactMobile,
				TenantID:    i.tenantID,
				Extra:       extra,
			},
		})
	}
	return result
}
