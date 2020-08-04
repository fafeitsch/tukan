package params

import (
	"encoding/json"
	"github.com/goccy/go-yaml"
)

type Credentials struct {
	Login    string
	Password string
}

// Parameters describes all the settings of the VoIP phone. The phones have
// slightly different format for the parameter upload/download. While unmarshalling
// downloaded parameters, both formats can be unmarshalled.
// For marshalling, the upload format is used.
type Parameters struct {
	AcceptAllCertificates                    Setting `json:"AcceptAllCertificates,omitempty"`
	AcceptInvalidCSeq                        Setting `json:"AcceptInvalidCSeq,omitempty"`
	AccessCode                               Setting `json:"AccessCode,omitempty"`
	AccessCodeEnabled                        Setting `json:"AccessCodeEnabled,omitempty"`
	AccessCodeFor                            Setting `json:"AccessCodeFor,omitempty"`
	AccessCodeInternalNumberLength           Setting `json:"AccessCodeInternalNumberLength,omitempty"`
	ActiveRingbackDisable                    Setting `json:"ActiveRingbackDisable,omitempty"`
	AdaptiveJitterBufferInitialPrefetchValue Setting `json:"AdaptiveJitterBufferInitialPrefetchValue,omitempty"`
	AdaptiveJitterBufferMaximumDelay         Setting `json:"AdaptiveJitterBufferMaximumDelay,omitempty"`
	AdaptiveJitterBufferMinimumDelay         Setting `json:"AdaptiveJitterBufferMinimumDelay,omitempty"`
	AdminPassword                            Setting `json:"AdminPassword,omitempty"`
	AllowAccessWeb                           Setting `json:"AllowAccessWeb,omitempty"`
	AllowFragmentation                       Setting `json:"AllowFragmentation,omitempty"`
	AllowHttpOutgoing                        Setting `json:"AllowHttpOutgoing,omitempty"`
	AllowSpanning                            Setting `json:"AllowSpanning,omitempty"`
	AnonymousCallBlock                       Setting `json:"AnonymousCallBlock,omitempty"`
	AreaCodesCountry                         Setting `json:"AreaCodesCountry,omitempty"`
	AreaCodesIntCode                         Setting `json:"AreaCodesIntCode,omitempty"`
	AreaCodesIntPrefix                       Setting `json:"AreaCodesIntPrefix,omitempty"`
	AreaCodesLocalCode                       Setting `json:"AreaCodesLocalCode,omitempty"`
	AreaCodesLocalPrefix                     Setting `json:"AreaCodesLocalPrefix,omitempty"`
	AutoAdjustClockForDST                    Setting `json:"AutoAdjustClockForDST,omitempty"`
	AutoAdjustTime                           Setting `json:"AutoAdjustTime,omitempty"`
	AutoDetermineAddress                     Setting `json:"AutoDetermineAddress,omitempty"`
	AutomaticCheckForUpdates                 Setting `json:"AutomaticCheckForUpdates,omitempty"`
	AutomaticFKFilling                       Setting `json:"AutomaticFKFilling,omitempty"`
	AutomaticRebootEnabled                   Setting `json:"AutomaticRebootEnabled,omitempty"`
	AutomaticRebootTime                      Setting `json:"AutomaticRebootTime,omitempty"`
	AutomaticRebootWeekdays                  Setting `json:"AutomaticRebootWeekdays,omitempty"`
	// TODO support format: AvailableCodecs                          Setting              `json:"AvailableCodecs,omitempty"`
	BLFCallPickupCode                      Setting              `json:"BLFCallPickupCode,omitempty"`
	BLFURL                                 Setting              `json:"BLFURL,omitempty"`
	Backlight                              Setting              `json:"Backlight,omitempty"`
	BroadsoftACDEnabled                    Setting              `json:"BroadSoftACDEnabled,omitempty"`
	BroadsoftACDStatus                     Setting              `json:"BroadsoftACDStatus,omitempty"`
	BroadsoftLookupIncomingEnabled         Setting              `json:"BroadsoftLookupIncomingEnabled,omitempty"`
	BroadsoftLookupOutgoingEnabled         Setting              `json:"BroadsoftLookupOutgoingEnabled,omitempty"`
	BroadsoftRemoteOfficeVisible           Setting              `json:"BroadsoftRemoteOfficeVisible,omitempty"`
	CallDiverDisable                       Setting              `json:"CallDivertDisable,omitempty"`
	CallDivertAll                          []CallDivertAll      `json:"CallDivertAll,omitempty"`
	CallDivertBusy                         []CallDivertBusy     `json:"CallDivertBusy,omitempty"`
	CallDivertNoAnswer                     []CallDivertNoAnswer `json:"CallDivertNoAnswer,omitempty"`
	CallWaitingDisable                     Setting              `json:"CallWaitingDisable,omitempty"`
	CallsViaCallManager                    Setting              `json:"CallsViaCallManager,omitempty"`
	ClearSIPMessageWithBackKey             Setting              `json:"ClearSIPMessageWithBackKey,omitempty"`
	ColourSchemeBasic                      Setting              `json:"ColourSchemeBasic,omitempty"`
	ColourSchemeThress                     Setting              `json:"ColourSchemeThree,omitempty"`
	ConfigurationCode                      Setting              `json:"ConfigurationCode,omitempty"`
	ConfigurationWith                      Setting              `json:"ConfigurationWith,omitempty"`
	ConnectionEstablished                  Setting              `json:"ConnectionEstablished,omitempty"`
	ConnectionTerminated                   Setting              `json:"ConnectionTerminated,omitempty"`
	ContactsDownloadPath                   Setting              `json:"ContactsDownloadPath,omitempty"`
	Contrast                               Setting              `json:"Contrast,omitempty"`
	CoreDumpsEnabled                       Setting              `json:"CoreDumpsEnabled,omitempty"`
	DateOrder                              Setting              `json:"DateOrder,omitempty"`
	DebugLEvelMaskAUTOREBOOT               Setting              `json:"DebugLevelMaskAutoREBOOT,omitempty"`
	DebugLevelMaskAUDIO                    Setting              `json:"DebugLevelMaskAUDIO,omitempty"`
	DebugLevelMaskDisplayGUI               Setting              `json:"DebugLevelMaskDisplayGUI,omitempty"`
	DebugLevelMaskEXTENSIONBOARD           Setting              `json:"DebugLevelMaskEXTENSIONBOARD,omitempty"`
	DebugLevelMaskMTIMERS                  Setting              `json:"DebugLevelMaskMTIMERS,omitempty"`
	DebugLevelMaskNETWORK                  Setting              `json:"DebugLevelMaskNETWORK,omitempty"`
	DebugLevelMaskPCM                      Setting              `json:"DebugLevelMaskPCM,omitempty"`
	DebugLevelMaskSIP                      Setting              `json:"DebugLevelMaskSIP,omitempty"`
	DebugLevelMaskSYSCONF                  Setting              `json:"DebugLevelMaskSYSCONF,omitempty"`
	DebugLevelMaskWATCHDOG                 Setting              `json:"DebugLevelMaskWATCHDOG,omitempty"`
	DefaultAccount                         Setting              `json:"DefaultAccount,omitempty"`
	DefaultRingtone                        Setting              `json:"DefaultRingtone,omitempty"`
	DefaultURLForAdmin                     Setting              `json:"DefaultURLForAdmin,omitempty"`
	DefaultURLForUser                      Setting              `json:"DefaultURLForUser,omitempty"`
	DeriveTargetAddress                    Setting              `json:"DeriveTargetAddress,omitempty"`
	DeviceNameInNetwork                    Setting              `json:"DeviceNameInNetwork,omitempty"`
	DialingPlans                           []DialingPlan        `json:"DialingPlans,omitempty"`
	DisableWebUI                           Setting              `json:"DisableWebUI,omitempty"`
	DisplayDiversionInfo                   Setting              `json:"DisplayDiversionInfo,omitempty"`
	DistinctiveRingingEnabled              Setting              `json:"DistinctiveRingingEnabled,omitempty"`
	DnDListActive                          Setting              `json:"DnDListActive,omitempty"`
	Dnd                                    []Dnd                `json:"DnD,omitempty"`
	DoorStations                           []DoorStation        `json:"DoorStations,omitempty"`
	EnableAec                              Setting              `json:"EnableAec,omitempty"`
	EnablePortMirroring                    Setting              `json:"EnablePortMirroring,omitempty"`
	FirmwareDataServer                     Setting              `json:"FirmwareDataServer,omitempty"`
	FirmwareDownloadPath                   Setting              `json:"FirmwareDownloadPath,omitempty"`
	FlexibleSeatingEnabled                 Setting              `json:"FlexibleSeatingEnabled,omitempty"`
	FunctionKeys                           FunctionKeys         `json:"FunctionKeys,omitempty"`
	FunctionKeysIcons                      Setting              `json:"FunctionKeysIcons,omitempty"`
	HTTPAuthPassword                       Setting              `json:"HTTPAuthPassword,omitempty"`
	HTTPConnectionType                     Setting              `json:"HTTPConnectionType,omitempty"`
	HTTPPort                               Setting              `json:"HTTPPort,omitempty"`
	HTTPSPort                              Setting              `json:"HTTPSPort,omitempty"`
	HandSetMode                            Setting              `json:"HandsetMode,omitempty"`
	HandsFreeMode                          Setting              `json:"HandsFreeMode,omitempty"`
	HeadsetMode                            Setting              `json:"HeadsetMode,omitempty"`
	HoldOnTransferAttended                 Setting              `json:"HoldOnTransferAttended,omitempty"`
	HoldOnTransferUnattended               Setting              `json:"HoldOnTransferUnattended,omitempty"`
	HttpAuthUsername                       Setting              `json:"HTTPAuthUsername,omitempty"`
	IPAddressType                          Setting              `json:"IPAddressType,omitempty"`
	IPv4Address                            Setting              `json:"IPv4Address,omitempty"`
	IPv4AlternateDNSServer                 Setting              `json:"IPv4AlternateDNSServer,omitempty"`
	IPv4PreferredDNSServer                 Setting              `json:"IPv4PreferredDNSServer,omitempty"`
	IPv4StandardGateway                    Setting              `json:"IPv4StandardGateway,omitempty"`
	IPv4SubnetMask                         Setting              `json:"IPv4SubnetMask,omitempty"`
	IncCallsWithoutCallManager             Setting              `json:"IncCallsWithoutCallManager,omitempty"`
	IncomingCall                           Setting              `json:"IncomingCall,omitempty"`
	LANPort                                Setting              `json:"LANPort,omitempty"`
	LDAPAdditionalAttribute                Setting              `json:"LDAPAdditionalAttribute,omitempty"`
	LDAPAdditionalAttributeDialable        Setting              `json:"LDAPAdditionalAttributeDialable,omitempty"`
	LDAPBaseDN                             Setting              `json:"LDAPBaseDN,omitempty"`
	LDAPCity                               Setting              `json:"LDAPCity,omitempty"`
	LDAPCompany                            Setting              `json:"LDAPCompany,omitempty"`
	LDAPCountry                            Setting              `json:"LDAPCountry,omitempty"`
	LDAPDirectoryName                      Setting              `json:"LDAPDirectoryName,omitempty"`
	LDAPDisplayFormat                      Setting              `json:"LDAPDisplayFormat,omitempty"`
	LDAPEmail                              Setting              `json:"LDAPEmail,omitempty"`
	LDAPEnable                             Setting              `json:"LDAPEnable,omitempty"`
	LDAPFax                                Setting              `json:"LDAPFax,omitempty"`
	LDAPFirstName                          Setting              `json:"LDAPFirstName,omitempty"`
	LDAPLookup                             Setting              `json:"LDAPLookup,omitempty"`
	LDAPMaxNumberOfSearchResults           Setting              `json:"LDAPMaxNumberOfSearchResults,omitempty"`
	LDAPNameFilter                         Setting              `json:"LDAPNameFilter,omitempty"`
	LDAPNumberFilter                       Setting              `json:"LDAPNumberFilter,omitempty"`
	LDAPPAssword                           Setting              `json:"LDAPPassword,omitempty"`
	LDAPPhoneHome                          Setting              `json:"LDAPPhoneHome,omitempty"`
	LDAPPhoneMobile                        Setting              `json:"LDAPPhoneMobile,omitempty"`
	LDAPPhoneOffice                        Setting              `json:"LDAPPhoneOffice,omitempty"`
	LDAPResponseTimeout                    Setting              `json:"LDAPResponseTimeout,omitempty"`
	LDAPSecurity                           Setting              `json:"LDAPSecurity,omitempty"`
	LDAPServerAddress                      Setting              `json:"LDAPServerAddress,omitempty"`
	LDAPServerPort                         Setting              `json:"LDAPServerPort,omitempty"`
	LDAPStreet                             Setting              `json:"LDAPStreet,omitempty"`
	LDAPSurname                            Setting              `json:"LDAPSurname,omitempty"`
	LDAPUsername                           Setting              `json:"LDAPUsername,omitempty"`
	LDAPZIP                                Setting              `json:"LDAPZIP,omitempty"`
	LLDAPActive                            Setting              `json:"LLDAPActive,omitempty"`
	LLDAPPacketInterval                    Setting              `json:"LLDAPPacketInterval,omitempty"`
	LinkSpeedDuplexLanPort                 Setting              `json:"LinkSpeedDuplexLanPort,omitempty"`
	LinkedSpeedDuplexPcPort                Setting              `json:"LinkedSpeedDuplexPcPort,omitempty"`
	LogoutTimer                            Setting              `json:"LogoutTimer,omitempty"`
	LookupIncoming                         Setting              `json:"LookupIncomming,omitempty"` // sic!
	LookupOutgoing                         Setting              `json:"LookupOutgoing,omitempty"`
	MACAddress                             Setting              `json:"MACAddress,omitempty"`
	MainMenuContent                        Setting              `json:"MainMenuContent,omitempty"`
	MenuAdaptiveJitterBuffer               Setting              `json:"MenuAdaptiveJitterBuffer,omitempty"`
	MenuAdjustment                         Setting              `json:"MenuAdjustment,omitempty"`
	MenuAudo                               Setting              `json:"MenuAudio,omitempty"`
	MenuCLI                                Setting              `json:"MenuCLI,omitempty"`
	MenuCallDivert                         Setting              `json:"MenuCallDivert,omitempty"`
	MenuCallHistory                        Setting              `json:"MenuCallHistory,omitempty"`
	MenuCallSettings                       Setting              `json:"MenuCallSettings,omitempty"`
	MenuConnections                        Setting              `json:"MenuConnections,omitempty"`
	MenuCoreDumps                          Setting              `json:"MenuCoreDumps,omitempty"`
	MenuCorporate                          Setting              `json:"MenuCorporate,omitempty"`
	MenuDateAndTime                        Setting              `json:"MenuDateAndTime,omitempty"`
	MenuDebugging                          Setting              `json:"MenuDebugging,omitempty"`
	MenuDeskPhone                          Setting              `json:"MenuDeskPhone,omitempty"`
	MenuDevice                             Setting              `json:"MenuDevice,omitempty"`
	MenuDiallingPlans                      Setting              `json:"MenuDiallingPlans,omitempty"`
	MenuDisplay                            Setting              `json:"MenuDisplay,omitempty"`
	MenuDoNotDisturb                       Setting              `json:"MenuDoNotDisturb,omitempty"`
	MenuDoorStation                        Setting              `json:"MenuDoorStation,omitempty"`
	MenuEvents                             Setting              `json:"MenuEvents,omitempty"`
	MenuExpert                             Setting              `json:"menuExpert,omitempty"`
	MenuExpertAudio                        Setting              `json:"MenuExpertAudio,omitempty"`
	MenuExpertNetwork                      Setting              `json:"MenuExportNetwork,omitempty"`
	MenuExtensionModule1                   Setting              `json:"MenuExtensionModule1,omitempty"`
	MenuExtensionModule2                   Setting              `json:"MenuExtensionModule2,omitempty"`
	MenuExtensionModule3                   Setting              `json:"MenuExtensionModule3,omitempty"`
	MenuFirmwareUpdate                     Setting              `json:"MenuFirmwareUpdate,omitempty"`
	MenuIP                                 Setting              `json:"MenuIP,omitempty"`
	MenuKeysAndLEDs                        Setting              `json:"MenuKeysAndLEDs,omitempty"`
	MenuLan                                Setting              `json:"MenuLAN`
	MenuLocalPhonebook                     Setting              `json:"MenuLocalPhonebook,omitempty"`
	MenuMainMenu                           Setting              `json:"MenuMainMenu,omitempty"`
	MenuMessageNotification                Setting              `json:"MenuMessageNotification,omitempty"`
	MenuNetwork                            Setting              `json:"MenuNetwork,omitempty"`
	MenuOnlineDirectories                  Setting              `json:"MenuOnlineDirectories,omitempty"`
	MenuOnlineServices                     Setting              `json:"MenuOnlineServices',omitempty"`
	MenuPCAPLogging                        Setting              `json:"MenuPCAPLogging,omitempty"`
	MenuPhoneSystem                        Setting              `json:"MenuPhoneSystem,omitempty"`
	MenuPhoneWebServer                     Setting              `json:"MenuPhoneWebServer,omitempty"`
	MenuPictures                           Setting              `json:"MenuPictures,omitempty"`
	MenuProvisioningConfiguration          Setting              `json:"MenuProvisioningConfiguration,omitempty"`
	MenuPublic                             Setting              `json:"MenuPublic,omitempty"`
	MenuReboot                             Setting              `json:"MenuReboot,omitempty"`
	MenuRebootAndReset                     Setting              `json:"MenuRebootAndReset,omitempty"`
	MenuRingtones                          Setting              `json:"MenuRingtones,omitempty"`
	MenuSIPProtocol                        Setting              `json:"MenuSIPPRotocol,omitempty"`
	MenuSaveAndRestore                     Setting              `json:"MenuSaveAndRestore,omitempty"`
	MenuSecurity                           Setting              `json:"MenuSecurity,omitempty"`
	MenuSettings                           Setting              `json:"MenuSettings,omitempty"`
	MenuStatus                             Setting              `json:"MenuStatus,omitempty"`
	MenuStatusConnections                  Setting              `json:"MenuStatusConnections,omitempty"`
	MenuStorageAllocation                  Setting              `json:"MenuStorageAllocation,omitempty"`
	MenuSwitchSettings                     Setting              `json:"MenuSwitchSettings,omitempty"`
	MenuSystem                             Setting              `json:"MenuSystem,omitempty"`
	MenuSystemLogging                      Setting              `json:"MenuSystemLogging,omitempty"`
	MenuTelephony                          Setting              `json:"MenuTelephony,omitempty"`
	MenuVoIP                               Setting              `json:"MenuVoIP,omitempty"`
	MenuVoiceMail                          Setting              `json:"MenuVoiceMail,omitempty"`
	MenuWebConfigurator                    Setting              `json:"MenuWebConfigurator,omitempty"`
	MenuWebcam                             Setting              `json:"MenuWebcam,omitempty"`
	MenuXML                                Setting              `json:"MenuXML,omitempty"`
	MissedCallsNotificationActive          Setting              `json:"MissedCallsNotificationActive,omitempty"`
	NetworkType                            Setting              `json:"NetworkType,omitempty"`
	OffHook                                Setting              `json:"OffHook,omitempty"`
	OnHook                                 Setting              `json:"OnHook,omitempty"`
	OneMelodyRingtoneDoorStation           Setting              `json:"OneMelodyRingtoneDoorStation,omitempty"`
	OneMelodyRingtoneExternal              Setting              `json:"OneMelodyRingtoneExternal,omitempty"`
	OneMelodyRingtoneGroup                 Setting              `json:"OneMelodyRingtoneGroup,omitempty"`
	OneMelodyRingtoneInterval              Setting              `json:"OneMelodyRingtoneInterval,omitempty"`
	OneMelodyRingtoneOptional              Setting              `json:"OneMelodyRingtoneOptional,omitempty"`
	OutCallsViaFunctionKey                 Setting              `json:"OutCallsViaFunctionKey,omitempty"`
	OutgoingCall                           Setting              `json:"OutgoingCall,omitempty"`
	PCPort                                 Setting              `json:"PCPort,omitempty"`
	PIN                                    Setting              `json:"PIN,omitempty"`
	PacketTimeForRtp                       Setting              `json:"PacketTimeForRTP,omitempty"`
	PhoneLanguage                          Setting              `json:"PhoneLanguage,omitempty"`
	PhoneModel                             Setting              `json:"PhoneModel,omitempty"`
	PhoneName                              Setting              `json:"PhoneName,omitempty"`
	PhoneSystem                            Setting              `json:"PhoneSystem,omitempty"`
	ProgrammableKeys                       Setting              `json:"ProgrammableKeysDNDActionURLDisable,omitempty"`
	ProgrammableKeysConferenceActionURL    Setting              `json:"ProgrammableKeysConferenceActionURL,omitempty"`
	ProgrammableKeysConferenceDTMFCode     Setting              `json:"ProgrammableKeysDTMFCode,omitempty"`
	ProgrammableKeysConferenceFAC          Setting              `json:"ProgrammableKeysConferenceFAC,omitempty"`
	ProgrammableKeysConferenceType         Setting              `json:"ProgrammableKeysConferenceType,omitempty"`
	ProgrammableKeysDNDActionURLEnable     Setting              `json:"ProgrammableKeysDNDActionURLEnable,omitempty"`
	ProgrammableKeysDNDFACDisable          Setting              `json:"ProgrammableKeysDNDFACDisable,omitempty"`
	ProgrammableKeysDNDType                Setting              `json:"ProgrammableKeysDNDType,omitempty"`
	ProgrammableKeysDirectoryType          Setting              `json:"ProgrammableKeysDirectoryType,omitempty"`
	ProgrammableKeysHoldActionURL          Setting              `json:"ProgrammableKeysHoldActionURL,omitempty"`
	ProgrammableKeysHoldDTMFCode           Setting              `json:"ProgrammableKeysHoldDTMFCode,omitempty"`
	ProgrammableKeysHoldType               Setting              `json:"ProgrammableKeysHoldType,omitempty"`
	ProgrammableKeysMessagesActionURL      Setting              `json:"ProgrammableKeysMessagesActionURL,omitempty"`
	ProgrammableKeysMessagesFAC            Setting              `json:"ProgrammableKeysMessagesFAC,omitempty"`
	ProgrammableKeysMessagesType           Setting              `json:"ProgrammableKeysMessagesType,omitempty"`
	ProvisioningDirectLink                 Setting              `json:"ProvisioningDirectLink,omitempty"`
	ProvisioningServer                     Setting              `json:"ProvisioningServer,omitempty"`
	ProxyServerActive                      Setting              `json:"ProxyServerActive,omitempty"`
	ProxyServerAddress                     Setting              `json:"ProxyServerAddress,omitempty"`
	ProxyServerPort                        Setting              `json:"ProxyServerPort,omitempty"`
	QuickDialKeys                          []QuickDialKey       `json:"QuickDialKeys,omitempty"`
	RTPQoSDSCP                             Setting              `json:"RTPoSDSCP,omitempty"`
	RegistrationFailed                     Setting              `json:"RegistrationFailed,omitempty"`
	RegistrationSucceeded                  Setting              `json:"RegistrationSucceeded,omitempty"`
	RemoteControlAllow                     Setting              `json:"RemoteControlAllow,omitempty"`
	RemoteControlSource                    Setting              `json:"RemoteControlSource,omitempty"`
	SIPAccountFailover                     Setting              `json:"SIPAccountFailover,omitempty"`
	SIPG729AnnexB                          Setting              `json:"SIPG729AnnexB,omitempty"`
	SIPNoSrtpCalls                         Setting              `json:"SipNoSrtpCalls,omitempty"`
	SIPPort                                Setting              `json:"SIPPort,omitempty"`
	SIPPrack                               Setting              `json:"SIPPRack,omitempty"`
	SIPQoSDSCP                             Setting              `json:"SIPQoSDSCP,omitempty"`
	SIPRtpPort                             Setting              `json:"SIPRtpPort,omitempty"`
	SIPRtpRTCPXRServerAddress              Setting              `json:"SIPRtpRTCPXRServerAddress,omitempty"`
	SIPRtpRTCPXRServerPort                 Setting              `json:"SIPRtpRTCXRServerPort,omitempty"`
	SIPRtpRandomPort                       Setting              `json:"SipRtpRandomPort,omitempty"`
	SIPRtpSymmetricPort                    Setting              `json:"SIPRtpSymetricPort,omitempty"` // sic!
	SIPRtpUseRTCPXR                        Setting              `json:"SIPRtpUseRTCPXR,omitempty"`
	SIPRtprRtcp                            Setting              `json:"SIPRtprRtcp,omitempty"`
	SIPSCertificate                        Setting              `json:"SIPSCertificate,omitempty"`
	SIPSKeyPassword                        Setting              `json:"SIPSKeyPassword,omitempty"`
	SIPSPrivateKey                         Setting              `json:"SIPSPrivateKey,omitempty"`
	SIPSecurity                            Setting              `json:"SIPSecurity,omitempty"`
	SIPSecurityEnabled                     Setting              `json:"SIPSecurityEnabled,omitempty"`
	SIPSessionTimer                        Setting              `json:"SIPSessionTimer,omitempty"`
	SIPTimerT1                             Setting              `json:"SIPTimerT1,omitempty"`
	SIPTimersFailedRegistration            Setting              `json:"SIPTimersFailedRegistration,omitempty"`
	SIPTimersFailedSubscription            Setting              `json:"SIPTimersFailSubscription,omitempty"`
	SIPTimersSubscription                  Setting              `json:"SIPTimersSubscription,omitempty"`
	SIPTimersSubscriptionBLFFollowRegister Setting              `json:"SIPTimersSubscriptionBLFFollowRegister,omitempty"`
	SIPTransportProtocol                   Setting              `json:"SIPTransportProtocol,omitempty"`
	SIPrtp                                 Setting              `json:"SIPrtp,omitempty"`
	ScreenSaverTimeout                     Setting              `json:"ScreenSaverTimeout,omitempty"`
	Screensaver                            Setting              `json:"Screensaver,omitempty"`
	ScreensaverBacklight                   Setting              `json:"ScreensaverBacklight,omitempty"`
	ScreensaverHTTPSource                  Setting              `json:"ScreensaverHTTPSource,omitempty"`
	ScreensaverPictures                    Setting              `json:"ScreensaverPictures,omitempty"`
	SelectedCodecs                         Setting              `json:"SelectedCodecs,omitempty"`
	SelectedServesDisable                  Setting              `json:"SelectedServicesDisable,omitempty"`
	SemiAttendedTransferType               Setting              `json:"SemiAttendedTransferType,omitempty"`
	SettingsVersion                        Setting              `json:"SettingsVersion,omitempty"`
	ShowPIN                                Setting              `json:"ShowPIN,omitempty"`
	ShowPassword                           Setting              `json:"ShowPassword,omitempty"`
	ShowSIPMessagesOnDisplay               Setting              `json:"ShowSIPMessagesOnDisplay,omitempty"`
	Sip                                    []Sip                `json:"SIP,omitempty"`
	SoftReboots                            Setting              `json:"SoftReboots,omitempty"`
	SoftwareVariant                        Setting              `json:"SoftwareVariant,omitempty"`
	SoftwareVersion                        Setting              `json:"SoftwareVersion,omitempty"`
	StandbyBacklight                       Setting              `json:"StandbyBacklight,omitempty"`
	Startups                               Setting              `json:"Startups,omitempty"`
	SyslogEnabled                          Setting              `json:"SyslogEnabled,omitempty"`
	SyslogServer                           Setting              `json:"SyslogServer,omitempty"`
	SystemLocalPhonebookUpdateTime         Setting              `json:"SystemLocalPhonebookUpdateTime,omitempty"`
	SystemLocalPhonebookUrl                Setting              `json:"SystemLocalPhonebookUrl,omitempty"`
	TimeFormat                             Setting              `json:"TimeFormat,omitempty"`
	TimeServer                             Setting              `json:"TimeServer,omitempty"`
	TimeServerDHCP                         Setting              `json:"TimeServerDHCP,omitempty"`
	TimeServerProvisioning                 Setting              `json:"TimeServerProvisioning,omitempty"`
	TimeZone                               Setting              `json:"TimeZone,omitempty"`
	Timestamp                              Setting              `json:"Timestamp,omitempty"`
	ToneScheme                             Setting              `json:"TomeScheme,omitempty"`
	UserPassword                           Setting              `json:"UserPassword,omitempty"`
	VLANIdentifierLAN                      Setting              `json:"VLANIdentifierLAN,omitempty"`
	VLANIdentifierPC                       Setting              `json:"VLANIdentifierPC,omitempty"`
	VLANLocked                             Setting              `json:"VLANLocked,omitempty"`
	VLANPriorityLAN                        Setting              `json:"VLANPriorityLAN,omitempty"`
	VLANPriorityPC                         Setting              `json:"VLANPriorityDC,omitempty"`
	VLANTagging                            Setting              `json:"VLANTagging,omitempty"`
	Variant                                Setting              `json:"Variant,omitempty"`
	VoiceQuality                           Setting              `json:"VoiceQuality,omitempty"`
	VoicemailMessagesActive                Setting              `json:"VoicemailMessagesActive,omitempty"`
	WebUICallDivertDisable                 Setting              `json:"WebUICallDivertDisable,omitempty"`
	WebUICallWaitingDisable                Setting              `json:"WebUICallWaitingDisable,omitempty"`
	WebUILanguage                          Setting              `json:"WebUILanguage,omitempty"`
	WitholdNumberDisable                   Setting              `json:"WitholdNumberDisable,omitempty"`
	WorkgingCounter                        Setting              `json:"WorkingCounter,omitempty"`
	WorkingCounterSec                      Setting              `json:"WorkingCounterSec,omitempty"`
	XMLProviderName                        Setting              `json:"XMLProviderName,omitempty"`
	XSIAuthName                            Setting              `json:"XSIAuthName,omitempty"`
	XSIAuthPassword                        Setting              `json:"XSIAuthPassword,omitempty"`
	XSICallLogType                         Setting              `json:"XSICallLogType,omitempty"`
	XSIEnterpriseCommonDirectoryEnabled    Setting              `json:"XSIEnterpriseCommonDirectoryEnabled,omitempty"`
	XSIEnterpriseCommonDirectoryName       Setting              `json:"XSIEnterpriseCommonDirectoryName,omitempty"`
	XSIEnterpriseDirectoryEnabled          Setting              `json:"XSIEnterpriseDirectoryEnabled,omitempty"`
	XSIEnterpriseDirectoryName             Setting              `json:"XSIEnterpriseDirectoryName,omitempty"`
	XSIGroupCommonDirectoryEnabled         Setting              `json:"XSIGroupCommonDirectoryEnabled,omitempty"`
	XSIGroupCommonDirectoryName            Setting              `json:"XSIGroupCommonDirectoryName,omitempty"`
	XSIGroupDirectoryEnabled               Setting              `json:"XSIGroupDirectoryEnabled,omitempty"`
	XSIGroupDirectoryName                  Setting              `json:"XSIGroupDirectoryName,omitempty"`
	XSIPersonalDirectoryEnabled            Setting              `json:"XSIPersonalDirectoryEnabled,omitempty"`
	XSIPersonalDirectoryName               Setting              `json:"XSIPersonalDirectoryName,omitempty"`
	XSISIPAuthentication                   Setting              `json:"XSISIPAuthentication,omitempty"`
	XSISearchAnywhereInNameEnabled         Setting              `json:"XSISearchAnywhereInNameEnabled,omitempty"`
	XSIServer                              Setting              `json:"XSIServer,omitempty"`
	XmlEnablePrivateDirectory              Setting              `json:"XmlEnablePrivateDirectory,omitempty"`
	XmlEnableWhiteDirectory                Setting              `json:"XmlEnableWhiteDirectory,omitempty"`
	XmlEnableYellowDirectory               Setting              `json:"XmlEnableYellowDirectory,omitempty"`
	XmlNumberFilter                        Setting              `json:"XmlNumberFilter,omitempty"`
	XmlPassword                            Setting              `json:"XMLPassword,omitempty"`
	XmlPrivateDirectoryName                Setting              `json:"XMLPrivateDirectoryName,omitempty"`
	XmlServerAddress                       Setting              `json:"XmlServerAddress,omitempty"`
	XmlUsername                            Setting              `json:"XmlUsername,omitempty"`
	XmlWhiteDirectoryName                  Setting              `json:"XmlWhiteDirectoryName,omitempty"`
	XmlYellowDirectoryName                 Setting              `json:"XmlYellowDirectoryName,omitempty"`
}

func (p *Parameters) TransformFunctionKeyNames(original, replace string) (Parameters, []int) {
	keys := make([]FunctionKey, 0, len(p.FunctionKeys))
	changed := make([]int, 0, 0)
	for index, fnKey := range p.FunctionKeys {
		var key = FunctionKey{}
		if fnKey.DisplayName.String() == original {
			key = FunctionKey{DisplayName: Setting(replace)}
			changed = append(changed, index)
		}
		keys = append(keys, key)
	}
	return Parameters{FunctionKeys: keys}, changed
}

func (p *Parameters) ExtendedMarshalling(enable bool) {

}

type FunctionKeys []FunctionKey

type FunctionKey struct {
	AutomaticallyFilled Setting `json:"AutomaticallyFilled,omitempty"`
	CallDivertType      Setting `json:"CallDivertType,omitempty"`
	CallPickupCode      Setting `json:"CallPickupCode,omitempty"`
	Color               Setting `json:"Color,omitempty"`
	Connection          Setting `json:"Connection,omitempty"`
	DTMFCode            Setting `json:"DTMFCode,omitempty"`
	DisableCode         Setting `json:"DisableCode,omitempty"`
	DisplayName         Setting `json:"DisplayName,omitempty"`
	EnableCode          Setting `json:"EnableCode,omitempty"`
	LockProvisioning    Setting `json:"LockProvisioning,omitempty"`
	PhoneNumber         Setting `json:"PhoneNumber,omitempty"`
	Silent              Setting `json:"Silent,omitempty"`
	Type                Setting `json:"Type,omitempty"`
	Url                 Setting `json:"URL,omitempty"`
}

func (f *FunctionKey) IsEmpty() bool {
	return (f.Type == "" && f.PhoneNumber == "" && f.DisplayName == "" && f.CallPickupCode == "") || f.Type == "-1"
}

func (f *FunctionKey) Merge(other FunctionKey) {
	if other.PhoneNumber != "" {
		f.PhoneNumber = other.PhoneNumber
	}
	if other.CallPickupCode != "" {
		f.CallPickupCode = other.PhoneNumber
	}
	if other.Type != "" {
		f.Type = other.Type
	}
	if other.DisplayName != "" {
		f.DisplayName = other.DisplayName
	}
}

type Setting string

func (s *Setting) String() string {
	return string(*s)
}

func (s *Setting) UnmarshalJSON(data []byte) error {
	setting := struct {
		Value json.RawMessage `json:"value,omitempty"`
	}{}
	err := json.Unmarshal(data, &setting)
	// data is not download-format, try the upload format:
	if err != nil {
		str := ""
		err = json.Unmarshal(data, &str)
		if err != nil {
			return err
		}
		*s = Setting(str)
		return nil
	}
	got := Setting(setting.Value)
	*s = got
	return nil
}

func (s *Setting) UnmarshalYAML(data []byte) error {
	setting := struct {
		Value string `json:"value,omitempty"`
	}{}
	err := yaml.Unmarshal(data, &setting)
	// data is not download-format, try the upload format:
	if err != nil {
		message := json.RawMessage{}
		err = yaml.Unmarshal(data, &message)
		if err != nil {
			return err
		}
		*s = Setting(message)
		return nil
	}
	got := Setting(setting.Value)
	*s = got
	return nil
}

type Sip struct {
	AccountName                   Setting `json:"AccountName,omitempty"`
	Active                        Setting `json:"Active,omitempty"`
	AllowRouteHeaders             Setting `json:"AllowRouteHeaders,omitempty"`
	AuthenaticationName           Setting `json:"AuthenticationName,omitempty"`
	AuthenticationPassword        Setting `json:"AuthenticationPassword,omitempty"`
	AutoNegOfDTMFTransmission     Setting `json:"AutoNetOfDTMSTransmission,omitempty"`
	CLIPSource                    Setting `json:"CLIPSource,omitempty"`
	CLIR                          Setting `json:"CLIR,omitempty"`
	CallWaiting                   Setting `json:"CallWaiting,omitempty"`
	CallWaitingSignal             Setting `json:"CallWaitingSignal,omitempty"`
	CountMissedAcceptedCalls      Setting `json:"CountMissedAcceptedCalls,omitempty"`
	DNSQuery                      Setting `json:"DNSQuery,omitempty"`
	DTMFTransmission              Setting `json:"DTMFTransmission,omitempty"`
	DisplayName                   Setting `json:"DisplayName,omitempty"`
	Domain                        Setting `json:"Domain,omitempty"`
	FailoverServerAddress         Setting `json:"FailoverServerAddress,omitempty"`
	FailoverServerEnabled         Setting `json:"FailoverServerEnabled,omitempty"`
	FailoverServerPort            Setting `json:"FailoverServerPort,omitempty"`
	HeaderDoorstation             Setting `json:"HeaderDoorstation,omitempty"`
	HeaderExternal                Setting `json:"HeaderExternal,omitempty"`
	HeaderGroup                   Setting `json:"HeaderGroup,omitempty"`
	HeaderInternal                Setting `json:"HeaderInternal,omitempty"`
	HeaderOptional                Setting `json:"HeaderOptional,omitempty"`
	ICE                           Setting `json:"ICE,omitempty"`
	NATRefreshTime                Setting `json:"NATRefreshTime,omitempty"`
	OutboundProxyAddress          Setting `json:"OutboundProxyAddress,omitempty"`
	OutboundProxyMode             Setting `json:"OutboundProxyMode,omitempty"`
	OutboundProxyPort             Setting `json:"OutboundProxyPort,omitempty"`
	Provider                      Setting `json:"Provider,omitempty"`
	ProxyServerAddress            Setting `json:"ProxyServerAddress,omitempty"`
	ProxyServerPort               Setting `json:"ProxyServerPort,omitempty"`
	RegistrationServerAddress     Setting `json:"RegistrationServerAddress,omitempty"`
	RegistrationServerPort        Setting `json:"RegistrationSeverPort,omitempty"`
	RegistrationServerRefreshTiem Setting `json:"RegistrationServerRefreshTime,omitempty"`
	RequestCheckOptions           Setting `json:"RequestCheckOptions,omitempty"`
	ReregisterAlternative         Setting `json:"ReregisterAlternative,omitempty"`
	RingtoneDoorStation           Setting `json:"RingtoneDoorStation,omitempty"`
	RingtoneExternal              Setting `json:"RingtoneExternal,omitempty"`
	RingtoneGroup                 Setting `json:"RingtoneGroup,omitempty"`
	RingtoneInternal              Setting `json:"RingtoneInternal,omitempty"`
	RingtoneOptional              Setting `json:"RingtoneOptional,omitempty"`
	STUNEnabled                   Setting `json:"STUNEnabled,omitempty"`
	STUNRefreshTime               Setting `json:"STUNRefreshTime,omitempty"`
	STUNServerAddress             Setting `json:"STUNServerAddress,omitempty"`
	STUNServerPort                Setting `json:"STUNServerPort,omitempty"`
	Username                      Setting `json:"Username,omitempty"`
	VoiceMailActive               Setting `json:"VoiceMailActive,omitempty"`
	VoiceMailMailbox              Setting `json:"VoiceMailMailbox,omitempty"`
}

type Dnd struct {
	PhoneNumber Setting `json:"PhoneNumber,omitempty"`
	Name        Setting `json:"Name,omitempty"`
}

type CallDivertAll struct {
	TargetMail Setting `json:"TargetMail,omitempty"`
	Target     Setting `json:"Target,omitempty"`
	VoiceMail  Setting `json:"VoiceMail,omitempty"`
	Active     Setting `json:"Active,omitempty"`
}

type QuickDialKey struct {
	Type      Setting `json:"Type,omitempty"`
	Number    Setting `json:"Number,omitempty"`
	FAC       Setting `json:"FAC,omitempty"`
	ActionURL Setting `json:"ActionURL,omitempty"`
}

type DoorStation struct {
	Password           Setting `json:"Password,omitempty"`
	Username           Setting `json:"Username,omitempty"`
	DTMFCode           Setting `json:"DTMFCode,omitempty"`
	CameraURL          Setting `json:"CameraURL,omitempty"`
	Name               Setting `json:"Name,omitempty"`
	PictureRefreshTime Setting `json:"PictureRefreshTime,omitempty"`
	SIPID              Setting `json:"SIPID,omitempty"`
}

type CallDivertNoAnswer struct {
	CallDivertBusy
	Delay Setting `json:"Delay,omitempty"`
}

type CallDivertBusy struct {
	VoiceMail  Setting `json:"VoiceMail,omitempty"`
	Active     Setting `json:"Active,omitempty"`
	TargetMail Setting `json:"TargetMail,omitempty"`
	Target     Setting `json:"Target,omitempty"`
}

type DialingPlan struct {
	PhoneNumber Setting `json:"PhoneNumber,omitempty"`
	Comment     Setting `json:"Comment,omitempty"`
	Active      Setting `json:"Active,omitempty"`
	Connection  Setting `json:"Connection,omitempty"`
	UseAreaCode Setting `json:"UseAreaCode,omitempty"`
}
