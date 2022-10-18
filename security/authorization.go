package security

// Scope type
type Scope string

// Scope values
var (
	ScopeGlobal = Scope("global")
	ScopeRealm  = Scope("realm")
	ScopeGroup  = Scope("group")
)

// Action type
type Action struct {
	Name  string
	Scope Scope
}

func (a Action) String() string {
	return a.Name
}

type Service int

const (
	BridgeService Service = iota
	EventService
	IDNowService
	PaperCardService
	SchedulerService
	SignatureService
	VoucherService
	AccreditationService
)

type API int

const (
	CommunicationAPI API = iota
	EventsAPI
	KycAPI
	ManagementAPI
	StatisticAPI
	TaskAPI
	IDNowAPI
	CardsAPI
	SchedulerAPI
	SignatureAPI
)

type ActionsIndex struct {
	index map[Service]map[API][]Action
}

func (a *ActionsIndex) addAction(service Service, api API, name string, scope Scope) Action {
	action := Action{Name: name, Scope: scope}

	if _, ok := a.index[service]; !ok {
		a.index[service] = map[API][]Action{}
	}

	a.index[service][api] = append(a.index[service][api], action)
	return action
}

func (a *ActionsIndex) GetActionsForAPIs(service Service, apis ...API) []Action {
	var actions []Action
	for _, api := range apis {
		actions = append(actions, a.index[service][api]...)
	}
	return actions
}

func (a *ActionsIndex) GetAllActions() []Action {
	var res []Action
	for _, apiActions := range a.index {
		for _, actions := range apiActions {
			res = append(res, actions...)
		}
	}
	return res
}

func (a *ActionsIndex) GetActionNamesForAPIs(service Service, apis ...API) []string {
	var names []string
	for _, api := range apis {
		for _, action := range a.index[service][api] {
			names = append(names, action.Name)
		}
	}
	return names
}

func (a *ActionsIndex) GetActionNamesForService(service Service) []string {
	var names []string
	for _, actions := range a.index[service] {
		for _, action := range actions {
			names = append(names, action.Name)
		}
	}
	return names
}

var Actions = ActionsIndex{index: map[Service]map[API][]Action{}}

var (
	COMSendEmail = Actions.addAction(BridgeService, CommunicationAPI, "COM_SendEmail", ScopeRealm)
	COMSendSMS   = Actions.addAction(BridgeService, CommunicationAPI, "COM_SendSMS", ScopeRealm)

	EVGetActions       = Actions.addAction(BridgeService, EventsAPI, "EV_GetActions", ScopeGlobal)
	EVGetEvents        = Actions.addAction(BridgeService, EventsAPI, "EV_GetEvents", ScopeRealm)
	EVGetEventsSummary = Actions.addAction(BridgeService, EventsAPI, "EV_GetEventsSummary", ScopeRealm)
	EVGetUserEvents    = Actions.addAction(BridgeService, EventsAPI, "EV_GetUserEvents", ScopeGroup)

	KYCGetActions                      = Actions.addAction(BridgeService, KycAPI, "KYC_GetActions", ScopeGlobal)
	KYCGetUserInSocialRealm            = Actions.addAction(BridgeService, KycAPI, "KYC_GetUserInSocialRealm", ScopeRealm)
	KYCGetUser                         = Actions.addAction(BridgeService, KycAPI, "KYC_GetUser", ScopeGroup)
	KYCGetUserByUsernameInSocialRealm  = Actions.addAction(BridgeService, KycAPI, "KYC_GetUserByUsernameInSocialRealm", ScopeRealm)
	KYCGetUserByUsername               = Actions.addAction(BridgeService, KycAPI, "KYC_GetUserByUsername", ScopeGroup)
	KYCValidateUserInSocialRealm       = Actions.addAction(BridgeService, KycAPI, "KYC_ValidateUserInSocialRealm", ScopeRealm)
	KYCValidateUser                    = Actions.addAction(BridgeService, KycAPI, "KYC_ValidateUser", ScopeGroup)
	KYCSendSmsConsentCodeInSocialRealm = Actions.addAction(BridgeService, KycAPI, "KYC_SendSmsConsentCodeInSocialRealm", ScopeRealm)
	KYCSendSmsConsentCode              = Actions.addAction(BridgeService, KycAPI, "KYC_SendSmsConsentCode", ScopeGroup)
	KYCSendSmsCodeInSocialRealm        = Actions.addAction(BridgeService, KycAPI, "KYC_SendSmsCodeInSocialRealm", ScopeRealm)
	KYCSendSmsCode                     = Actions.addAction(BridgeService, KycAPI, "KYC_SendSmsCode", ScopeGroup)
	KYCValidateUserBasicID             = Actions.addAction(BridgeService, KycAPI, "KYC_ValidateUserBasicID", ScopeRealm) /***TO BE REMOVED WHEN MULTI-ACCREDITATION WILL BE IMPLEMENTED***/ /***TO BE REMOVED WHEN MULTI-ACCREDITATION WILL BE IMPLEMENTED***/

	MGMTGetActions                          = Actions.addAction(BridgeService, ManagementAPI, "MGMT_GetActions", ScopeGlobal)
	MGMTGetRealms                           = Actions.addAction(BridgeService, ManagementAPI, "MGMT_GetRealms", ScopeGlobal)
	MGMTGetRealm                            = Actions.addAction(BridgeService, ManagementAPI, "MGMT_GetRealm", ScopeRealm)
	MGMTGetClient                           = Actions.addAction(BridgeService, ManagementAPI, "MGMT_GetClient", ScopeRealm)
	MGMTGetClients                          = Actions.addAction(BridgeService, ManagementAPI, "MGMT_GetClients", ScopeRealm)
	MGMTGetRequiredActions                  = Actions.addAction(BridgeService, ManagementAPI, "MGMT_GetRequiredActions", ScopeRealm)
	MGMTDeleteUser                          = Actions.addAction(BridgeService, ManagementAPI, "MGMT_DeleteUser", ScopeGroup)
	MGMTGetUser                             = Actions.addAction(BridgeService, ManagementAPI, "MGMT_GetUser", ScopeGroup)
	MGMTUpdateUser                          = Actions.addAction(BridgeService, ManagementAPI, "MGMT_UpdateUser", ScopeGroup)
	MGMTLockUser                            = Actions.addAction(BridgeService, ManagementAPI, "MGMT_LockUser", ScopeGroup)
	MGMTUnlockUser                          = Actions.addAction(BridgeService, ManagementAPI, "MGMT_UnlockUser", ScopeGroup)
	MGMTGetUsers                            = Actions.addAction(BridgeService, ManagementAPI, "MGMT_GetUsers", ScopeGroup)
	MGMTCreateUser                          = Actions.addAction(BridgeService, ManagementAPI, "MGMT_CreateUser", ScopeGroup)
	MGMTCreateUserInSocialRealm             = Actions.addAction(BridgeService, ManagementAPI, "MGMT_CreateUserInSocialRealm", ScopeRealm)
	MGMTGetUserChecks                       = Actions.addAction(BridgeService, ManagementAPI, "MGMT_GetUserChecks", ScopeGroup)
	MGMTGetUserAccountStatus                = Actions.addAction(BridgeService, ManagementAPI, "MGMT_GetUserAccountStatus", ScopeGroup)
	MGMTGetUserAccountStatusByEmail         = Actions.addAction(BridgeService, ManagementAPI, "MGMT_GetUserAccountStatusByEmail", ScopeRealm)
	MGMTGetRolesOfUser                      = Actions.addAction(BridgeService, ManagementAPI, "MGMT_GetRolesOfUser", ScopeGroup)
	MGMTAddRoleToUser                       = Actions.addAction(BridgeService, ManagementAPI, "MGMT_AddRoleToUser", ScopeGroup)
	MGMTDeleteRoleForUser                   = Actions.addAction(BridgeService, ManagementAPI, "MGMT_DeleteRoleForUser", ScopeGroup)
	MGMTGetGroupsOfUser                     = Actions.addAction(BridgeService, ManagementAPI, "MGMT_GetGroupsOfUser", ScopeGroup)
	MGMTSetGroupsToUser                     = Actions.addAction(BridgeService, ManagementAPI, "MGMT_SetGroupsToUser", ScopeGroup)
	MGMTAssignableGroupsToUser              = Actions.addAction(BridgeService, ManagementAPI, "MGMT_AssignableGroupsToUser", ScopeGroup)
	MGMTGetAvailableTrustIDGroups           = Actions.addAction(BridgeService, ManagementAPI, "MGMT_GetAvailableTrustIDGroups", ScopeRealm)
	MGMTGetTrustIDGroups                    = Actions.addAction(BridgeService, ManagementAPI, "MGMT_GetTrustIDGroups", ScopeGroup)
	MGMTSetTrustIDGroups                    = Actions.addAction(BridgeService, ManagementAPI, "MGMT_SetTrustIDGroups", ScopeGroup)
	MGMTGetClientRolesForUser               = Actions.addAction(BridgeService, ManagementAPI, "MGMT_GetClientRolesForUser", ScopeGroup)
	MGMTAddClientRolesToUser                = Actions.addAction(BridgeService, ManagementAPI, "MGMT_AddClientRolesToUser", ScopeGroup)
	MGMTDeleteClientRolesFromUser           = Actions.addAction(BridgeService, ManagementAPI, "MGMT_DeleteClientRolesFromUser", ScopeGroup)
	MGMTResetPassword                       = Actions.addAction(BridgeService, ManagementAPI, "MGMT_ResetPassword", ScopeGroup)
	MGMTExecuteActionsEmail                 = Actions.addAction(BridgeService, ManagementAPI, "MGMT_ExecuteActionsEmail", ScopeGroup)
	MGMTRevokeAccreditations                = Actions.addAction(BridgeService, ManagementAPI, "ACCR_RevokeAccreditations", ScopeGroup)
	MGMTSendSmsCode                         = Actions.addAction(BridgeService, ManagementAPI, "MGMT_SendSmsCode", ScopeGroup)
	MGMTSendOnboardingEmail                 = Actions.addAction(BridgeService, ManagementAPI, "MGMT_SendOnboardingEmail", ScopeGroup)
	MGMTSendOnboardingEmailInSocialRealm    = Actions.addAction(BridgeService, ManagementAPI, "MGMT_SendOnboardingEmailInSocialRealm", ScopeRealm)
	MGMTSendReminderEmail                   = Actions.addAction(BridgeService, ManagementAPI, "MGMT_SendReminderEmail", ScopeGroup)
	MGMTResetSmsCounter                     = Actions.addAction(BridgeService, ManagementAPI, "MGMT_ResetSmsCounter", ScopeGroup)
	MGMTCreateRecoveryCode                  = Actions.addAction(BridgeService, ManagementAPI, "MGMT_CreateRecoveryCode", ScopeGroup)
	MGMTCreateActivationCode                = Actions.addAction(BridgeService, ManagementAPI, "MGMT_CreateActivationCode", ScopeGroup)
	MGMTGetCredentialsForUser               = Actions.addAction(BridgeService, ManagementAPI, "MGMT_GetCredentialsForUser", ScopeGroup)
	MGMTDeleteCredentialsForUser            = Actions.addAction(BridgeService, ManagementAPI, "MGMT_DeleteCredentialsForUser", ScopeGroup)
	MGMTResetCredentialFailuresForUser      = Actions.addAction(BridgeService, ManagementAPI, "MGMT_ResetCredentialFailuresForUser", ScopeGroup)
	MGMTClearUserLoginFailures              = Actions.addAction(BridgeService, ManagementAPI, "MGMT_ClearUserLoginFailures", ScopeGroup)
	MGMTGetAttackDetectionStatus            = Actions.addAction(BridgeService, ManagementAPI, "MGMT_GetAttackDetectionStatus", ScopeGroup)
	MGMTGetRoles                            = Actions.addAction(BridgeService, ManagementAPI, "MGMT_GetRoles", ScopeRealm)
	MGMTGetRole                             = Actions.addAction(BridgeService, ManagementAPI, "MGMT_GetRole", ScopeRealm)
	MGMTCreateRole                          = Actions.addAction(BridgeService, ManagementAPI, "MGMT_CreateRole", ScopeRealm)
	MGMTUpdateRole                          = Actions.addAction(BridgeService, ManagementAPI, "MGMT_UpdateRole", ScopeRealm)
	MGMTDeleteRole                          = Actions.addAction(BridgeService, ManagementAPI, "MGMT_DeleteRole", ScopeRealm)
	MGMTGetGroups                           = Actions.addAction(BridgeService, ManagementAPI, "MGMT_GetGroups", ScopeRealm)
	MGMTIncludedInGetGroups                 = Actions.addAction(BridgeService, ManagementAPI, "MGMT_IncludedInGetGroups", ScopeGroup)
	MGMTCreateGroup                         = Actions.addAction(BridgeService, ManagementAPI, "MGMT_CreateGroup", ScopeRealm)
	MGMTDeleteGroup                         = Actions.addAction(BridgeService, ManagementAPI, "MGMT_DeleteGroup", ScopeGroup)
	MGMTGetAuthorizations                   = Actions.addAction(BridgeService, ManagementAPI, "MGMT_GetAuthorizations", ScopeGroup)
	MGMTUpdateAuthorizations                = Actions.addAction(BridgeService, ManagementAPI, "MGMT_UpdateAuthorizations", ScopeGroup)
	MGMTAddAuthorization                    = Actions.addAction(BridgeService, ManagementAPI, "MGMT_AddAuthorization", ScopeGroup)
	MGMTGetAuthorization                    = Actions.addAction(BridgeService, ManagementAPI, "MGMT_GetAuthorization", ScopeGroup)
	MGMTDeleteAuthorization                 = Actions.addAction(BridgeService, ManagementAPI, "MGMT_DeleteAuthorization", ScopeGroup)
	MGMTGetClientRoles                      = Actions.addAction(BridgeService, ManagementAPI, "MGMT_GetClientRoles", ScopeRealm)
	MGMTCreateClientRole                    = Actions.addAction(BridgeService, ManagementAPI, "MGMT_CreateClientRole", ScopeRealm)
	MGMTDeleteClientRole                    = Actions.addAction(BridgeService, ManagementAPI, "MGMT_DeleteClientRole", ScopeRealm)
	MGMTGetRealmCustomConfiguration         = Actions.addAction(BridgeService, ManagementAPI, "MGMT_GetRealmCustomConfiguration", ScopeRealm)
	MGMTUpdateRealmCustomConfiguration      = Actions.addAction(BridgeService, ManagementAPI, "MGMT_UpdateRealmCustomConfiguration", ScopeRealm)
	MGMTGetRealmAdminConfiguration          = Actions.addAction(BridgeService, ManagementAPI, "MGMT_GetRealmAdminConfiguration", ScopeRealm)
	MGMTUpdateRealmAdminConfiguration       = Actions.addAction(BridgeService, ManagementAPI, "MGMT_UpdateRealmAdminConfiguration", ScopeRealm)
	MGMTGetRealmBackOfficeConfiguration     = Actions.addAction(BridgeService, ManagementAPI, "MGMT_GetRealmBackOfficeConfiguration", ScopeGroup)
	MGMTUpdateRealmBackOfficeConfiguration  = Actions.addAction(BridgeService, ManagementAPI, "MGMT_UpdateRealmBackOfficeConfiguration", ScopeGroup)
	MGMTGetUserRealmBackOfficeConfiguration = Actions.addAction(BridgeService, ManagementAPI, "MGMT_GetUserRealmBackOfficeConfiguration", ScopeRealm)
	MGMTGetFederatedIdentities              = Actions.addAction(BridgeService, ManagementAPI, "MGMT_GetFederatedIdentities", ScopeGroup)
	MGMTLinkShadowUser                      = Actions.addAction(BridgeService, ManagementAPI, "MGMT_LinkShadowUser", ScopeGroup)
	MGMTGetIdentityProviders                = Actions.addAction(BridgeService, ManagementAPI, "MGMT_GetIdentityProviders", ScopeRealm)

	STGetActions                      = Actions.addAction(BridgeService, StatisticAPI, "ST_GetActions", ScopeGlobal)
	STGetStatistics                   = Actions.addAction(BridgeService, StatisticAPI, "ST_GetStatistics", ScopeRealm)
	STGetStatisticsIdentifications    = Actions.addAction(BridgeService, StatisticAPI, "ST_GetStatisticsIdentifications", ScopeRealm)
	STGetStatisticsUsers              = Actions.addAction(BridgeService, StatisticAPI, "ST_GetStatisticsUsers", ScopeRealm)
	STGetStatisticsAuthenticators     = Actions.addAction(BridgeService, StatisticAPI, "ST_GetStatisticsAuthenticators", ScopeRealm)
	STGetStatisticsAuthentications    = Actions.addAction(BridgeService, StatisticAPI, "ST_GetStatisticsAuthentications", ScopeRealm)
	STGetStatisticsAuthenticationsLog = Actions.addAction(BridgeService, StatisticAPI, "ST_GetStatisticsAuthenticationsLog", ScopeRealm)
	STGetMigrationReport              = Actions.addAction(BridgeService, StatisticAPI, "ST_GetMigrationReport", ScopeRealm)

	TSKDeleteDeniedToUUsers = Actions.addAction(BridgeService, TaskAPI, "TSK_DeleteDeniedToUUsers", ScopeGlobal)

	IDNGetActions     = Actions.addAction(IDNowService, IDNowAPI, "IDN_GetActions", ScopeGlobal)
	IDNVideoIdentInit = Actions.addAction(IDNowService, IDNowAPI, "IDN_Init", ScopeGroup)
	IDNAutoIdentInit  = Actions.addAction(IDNowService, IDNowAPI, "IDN_AutoIdentInit", ScopeGroup)

	PCGetActions            = Actions.addAction(PaperCardService, CardsAPI, "PC_GetActions", ScopeGlobal)
	PCGetConfigurationRealm = Actions.addAction(PaperCardService, CardsAPI, "PC_GetConfigurationRealm", ScopeRealm)
	PCSetConfigurationRealm = Actions.addAction(PaperCardService, CardsAPI, "PC_SetConfigurationRealm", ScopeRealm)
	PCGetConfigurationSelf  = Actions.addAction(PaperCardService, CardsAPI, "PC_GetConfigurationSelf", ScopeRealm)
	PCSetConfigurationSelf  = Actions.addAction(PaperCardService, CardsAPI, "PC_SetConfigurationSelf", ScopeRealm)
	PCGetConfigurationBatch = Actions.addAction(PaperCardService, CardsAPI, "PC_GetConfigurationBatch", ScopeRealm)
	PCSetConfigurationBatch = Actions.addAction(PaperCardService, CardsAPI, "PC_SetConfigurationBatch", ScopeRealm)
	PCPreview               = Actions.addAction(PaperCardService, CardsAPI, "PC_Preview", ScopeRealm)
	PCCreateBatch           = Actions.addAction(PaperCardService, CardsAPI, "PC_CreateBatch", ScopeRealm)
	PCGetBatches            = Actions.addAction(PaperCardService, CardsAPI, "PC_GetBatches", ScopeRealm)
	PCGetBatch              = Actions.addAction(PaperCardService, CardsAPI, "PC_GetBatch", ScopeRealm)
	PCDeleteBatch           = Actions.addAction(PaperCardService, CardsAPI, "PC_DeleteBatch", ScopeRealm)
	PCActivateBatch         = Actions.addAction(PaperCardService, CardsAPI, "PC_ActivateBatch", ScopeRealm)
	PCBlockBatch            = Actions.addAction(PaperCardService, CardsAPI, "PC_BlockBatch", ScopeRealm)
	PCDownloadBatch         = Actions.addAction(PaperCardService, CardsAPI, "PC_DownloadBatch", ScopeRealm)

	SDLRGetActions = Actions.addAction(SchedulerService, SchedulerAPI, "SDLR_GetActions", ScopeGlobal)
	SDLRGetTasks   = Actions.addAction(SchedulerService, SchedulerAPI, "SDLR_GetTasks", ScopeGlobal)
	SDLRAddTasks   = Actions.addAction(SchedulerService, SchedulerAPI, "SDLR_AddTasks", ScopeGlobal)
	SDLRDeleteTask = Actions.addAction(SchedulerService, SchedulerAPI, "SDLR_DeleteTask", ScopeGlobal)

	SIGGetActions = Actions.addAction(SignatureService, SignatureAPI, "SIG_GetActions", ScopeGlobal)
	// CH/AES was the default value: keep the same name
	SIGSignDocuments_CH_AES = Actions.addAction(SignatureService, SignatureAPI, "SIG_SignDocuments", ScopeRealm)
	SIGSignDocuments_CH_QES = Actions.addAction(SignatureService, SignatureAPI, "SIG_SignDocuments_CH_QES", ScopeRealm)
	SIGSignDocuments_EU_AES = Actions.addAction(SignatureService, SignatureAPI, "SIG_SignDocuments_EU_AES", ScopeRealm)
	SIGSignDocuments_EU_QES = Actions.addAction(SignatureService, SignatureAPI, "SIG_SignDocuments_EU_QES", ScopeRealm)

	VOUGetActions          = Actions.addAction(VoucherService, ManagementAPI, "VOU_GetActions", ScopeGlobal)
	VOUGetBatches          = Actions.addAction(VoucherService, ManagementAPI, "VOU_GetBatches", ScopeRealm)
	VOUCreateBatch         = Actions.addAction(VoucherService, ManagementAPI, "VOU_CreateBatch", ScopeRealm)
	VOUGetBatch            = Actions.addAction(VoucherService, ManagementAPI, "VOU_GetBatch", ScopeRealm)
	VOURevokeBatch         = Actions.addAction(VoucherService, ManagementAPI, "VOU_RevokeBatch", ScopeRealm)
	VOUDownloadBatch       = Actions.addAction(VoucherService, ManagementAPI, "VOU_DownloadBatch", ScopeRealm)
	VOUGetVoucher          = Actions.addAction(VoucherService, ManagementAPI, "VOU_GetVoucher", ScopeRealm)
	VOUGetConfiguration    = Actions.addAction(VoucherService, ManagementAPI, "VOU_GetConfiguration", ScopeRealm)
	VOUUpdateConfiguration = Actions.addAction(VoucherService, ManagementAPI, "VOU_UpdateConfiguration", ScopeRealm)
	VOUGetAbuseCounter     = Actions.addAction(VoucherService, ManagementAPI, "VOU_GetAbuseCounter", ScopeGroup)
	VOUResetAbuseCounter   = Actions.addAction(VoucherService, ManagementAPI, "VOU_ResetAbuseCounter", ScopeGroup)

	ACCRGetActions                   = Actions.addAction(AccreditationService, ManagementAPI, "ACCR_GetActions", ScopeGlobal)
	ACCRGetAllAccreditations         = Actions.addAction(AccreditationService, ManagementAPI, "ACCR_GetAllAccreditations", ScopeGlobal)
	ACCRGetEnabledAccreditations     = Actions.addAction(AccreditationService, ManagementAPI, "ACCR_GetEnabledAccreditations", ScopeRealm)
	ACCRGetAccreditation             = Actions.addAction(AccreditationService, ManagementAPI, "ACCR_GetAccreditation", ScopeRealm)
	ACCREnableAccreditation          = Actions.addAction(AccreditationService, ManagementAPI, "ACCR_EnableAccreditation", ScopeRealm)
	ACCRDisableAccreditation         = Actions.addAction(AccreditationService, ManagementAPI, "ACCR_DisableAccreditation", ScopeRealm)
	ACCRConfigureAccreditation       = Actions.addAction(AccreditationService, ManagementAPI, "ACCR_ConfigureAccreditation", ScopeRealm)
	ACCRGetAccreditationGroups       = Actions.addAction(AccreditationService, ManagementAPI, "ACCR_GetAccreditationGroups", ScopeGroup)
	ACCREnableAccreditationForGroup  = Actions.addAction(AccreditationService, ManagementAPI, "ACCR_EnableAccreditationForGroup", ScopeGroup)
	ACCRDisableAccreditationForGroup = Actions.addAction(AccreditationService, ManagementAPI, "ACCR_DisableAccreditationForGroup", ScopeGroup)
)
