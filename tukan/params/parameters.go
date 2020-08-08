package params

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
)

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Parameters describes all the settings of the VoIP phone. The phones have
// slightly different format for the parameter upload/download. While unmarshalling
// downloaded parameters, both formats can be unmarshalled.
// For marshalling, the upload format is used.
type Parameters struct {
	AcceptAllCertificates                    string               `json:"AcceptAllCertificates,omitempty"`
	AcceptInvalidCSeq                        string               `json:"AcceptInvalidCSeq,omitempty"`
	AccessCode                               string               `json:"AccessCode,omitempty"`
	AccessCodeEnabled                        string               `json:"AccessCodeEnabled,omitempty"`
	AccessCodeFor                            string               `json:"AccessCodeFor,omitempty"`
	AccessCodeInternalNumberLength           int                  `json:"AccessCodeInternalNumberLength,omitempty"`
	ActiveRingbackDisable                    string               `json:"ActiveRingbackDisable,omitempty"`
	AdaptiveJitterBufferInitialPrefetchValue int                  `json:"AdaptiveJitterBufferInitialPrefetchValue,omitempty"`
	AdaptiveJitterBufferMaximumDelay         int                  `json:"AdaptiveJitterBufferMaximumDelay,omitempty"`
	AdaptiveJitterBufferMinimumDelay         int                  `json:"AdaptiveJitterBufferMinimumDelay,omitempty"`
	AdminPassword                            string               `json:"AdminPassword,omitempty"`
	AllowAccessWeb                           string               `json:"AllowAccessWeb,omitempty"`
	AllowFragmentation                       string               `json:"AllowFragmentation,omitempty"`
	AllowHttpOutgoing                        string               `json:"AllowHttpOutgoing,omitempty"`
	AllowSpanning                            string               `json:"AllowSpanning,omitempty"`
	AnonymousCallBlock                       string               `json:"AnonymousCallBlock,omitempty"`
	AreaCodesCountry                         string               `json:"AreaCodesCountry,omitempty"`
	AreaCodesIntCode                         string               `json:"AreaCodesIntCode,omitempty"`
	AreaCodesIntPrefix                       string               `json:"AreaCodesIntPrefix,omitempty"`
	AreaCodesLocalCode                       string               `json:"AreaCodesLocalCode,omitempty"`
	AreaCodesLocalPrefix                     string               `json:"AreaCodesLocalPrefix,omitempty"`
	AutoAdjustClockForDST                    string               `json:"AutoAdjustClockForDST,omitempty"`
	AutoAdjustTime                           string               `json:"AutoAdjustTime,omitempty"`
	AutoDetermineAddress                     string               `json:"AutoDetermineAddress,omitempty"`
	AutomaticCheckForUpdates                 string               `json:"AutomaticCheckForUpdates,omitempty"`
	AutomaticFKFilling                       string               `json:"AutomaticFKFilling,omitempty"`
	AutomaticRebootEnabled                   string               `json:"AutomaticRebootEnabled,omitempty"`
	AutomaticRebootTime                      string               `json:"AutomaticRebootTime,omitempty"`
	AutomaticRebootWeekdays                  int                  `json:"AutomaticRebootWeekdays,omitempty"`
	AvailableCodecs                          string               `json:"AvailableCodecs,omitempty"`
	BLFCallPickupCode                        string               `json:"BLFCallPickupCode,omitempty"`
	BLFURL                                   string               `json:"BLFURL,omitempty"`
	Backlight                                int                  `json:"Backlight,omitempty"`
	BroadsoftACDEnabled                      string               `json:"BroadSoftACDEnabled,omitempty"`
	BroadsoftACDStatus                       string               `json:"BroadsoftACDStatus,omitempty"`
	BroadsoftLookupIncomingEnabled           string               `json:"BroadsoftLookupIncomingEnabled,omitempty"`
	BroadsoftLookupOutgoingEnabled           string               `json:"BroadsoftLookupOutgoingEnabled,omitempty"`
	BroadsoftRemoteOfficeVisible             string               `json:"BroadsoftRemoteOfficeVisible,omitempty"`
	CallDiverDisable                         string               `json:"CallDivertDisable,omitempty"`
	CallDivertAll                            []CallDivertAll      `json:"CallDivertAll,omitempty"`
	CallDivertBusy                           []CallDivertBusy     `json:"CallDivertBusy,omitempty"`
	CallDivertNoAnswer                       []CallDivertNoAnswer `json:"CallDivertNoAnswer,omitempty"`
	CallWaitingDisable                       string               `json:"CallWaitingDisable,omitempty"`
	CallsViaCallManager                      string               `json:"CallsViaCallManager,omitempty"`
	ClearSIPMessageWithBackKey               string               `json:"ClearSIPMessageWithBackKey,omitempty"`
	ColourSchemeBasic                        string               `json:"ColourSchemeBasic,omitempty"`
	ColourSchemeThress                       string               `json:"ColourSchemeThree,omitempty"`
	ConfigurationCode                        string               `json:"ConfigurationCode,omitempty"`
	ConfigurationWith                        string               `json:"ConfigurationWith,omitempty"`
	ConnectionEstablished                    string               `json:"ConnectionEstablished,omitempty"`
	ConnectionTerminated                     string               `json:"ConnectionTerminated,omitempty"`
	ContactsDownloadPath                     string               `json:"ContactsDownloadPath,omitempty"`
	Contrast                                 string               `json:"Contrast,omitempty"`
	CoreDumpsEnabled                         string               `json:"CoreDumpsEnabled,omitempty"`
	DateOrder                                string               `json:"DateOrder,omitempty"`
	DebugLEvelMaskAUTOREBOOT                 string               `json:"DebugLevelMaskAutoREBOOT,omitempty"`
	DebugLevelMaskAUDIO                      string               `json:"DebugLevelMaskAUDIO,omitempty"`
	DebugLevelMaskDisplayGUI                 string               `json:"DebugLevelMaskDisplayGUI,omitempty"`
	DebugLevelMaskEXTENSIONBOARD             string               `json:"DebugLevelMaskEXTENSIONBOARD,omitempty"`
	DebugLevelMaskMTIMERS                    string               `json:"DebugLevelMaskMTIMERS,omitempty"`
	DebugLevelMaskNETWORK                    string               `json:"DebugLevelMaskNETWORK,omitempty"`
	DebugLevelMaskPCM                        string               `json:"DebugLevelMaskPCM,omitempty"`
	DebugLevelMaskSIP                        string               `json:"DebugLevelMaskSIP,omitempty"`
	DebugLevelMaskSYSCONF                    string               `json:"DebugLevelMaskSYSCONF,omitempty"`
	DebugLevelMaskWATCHDOG                   string               `json:"DebugLevelMaskWATCHDOG,omitempty"`
	DefaultAccount                           string               `json:"DefaultAccount,omitempty"`
	DefaultRingtone                          string               `json:"DefaultRingtone,omitempty"`
	DefaultURLForAdmin                       string               `json:"DefaultURLForAdmin,omitempty"`
	DefaultURLForUser                        string               `json:"DefaultURLForUser,omitempty"`
	DeriveTargetAddress                      string               `json:"DeriveTargetAddress,omitempty"`
	DeviceNameInNetwork                      string               `json:"DeviceNameInNetwork,omitempty"`
	DialingPlans                             []DialingPlan        `json:"DialingPlans,omitempty"`
	DisableWebUI                             string               `json:"DisableWebUI,omitempty"`
	DisplayDiversionInfo                     string               `json:"DisplayDiversionInfo,omitempty"`
	DistinctiveRingingEnabled                string               `json:"DistinctiveRingingEnabled,omitempty"`
	DnDListActive                            string               `json:"DnDListActive,omitempty"`
	Dnd                                      []Dnd                `json:"DnD,omitempty"`
	DoorStations                             []DoorStation        `json:"DoorStations,omitempty"`
	EnableAec                                string               `json:"EnableAec,omitempty"`
	EnablePortMirroring                      string               `json:"EnablePortMirroring,omitempty"`
	FirmwareDataServer                       string               `json:"FirmwareDataServer,omitempty"`
	FirmwareDownloadPath                     string               `json:"FirmwareDownloadPath,omitempty"`
	FlexibleSeatingEnabled                   string               `json:"FlexibleSeatingEnabled,omitempty"`
	FunctionKeys                             FunctionKeys         `json:"FunctionKeys,omitempty"`
	FunctionKeysIcons                        string               `json:"FunctionKeysIcons,omitempty"`
	HTTPAuthPassword                         string               `json:"HTTPAuthPassword,omitempty"`
	HTTPConnectionType                       string               `json:"HTTPConnectionType,omitempty"`
	HTTPPort                                 int                  `json:"HTTPPort,omitempty"`
	HTTPSPort                                int                  `json:"HTTPSPort,omitempty"`
	HandSetMode                              string               `json:"HandsetMode,omitempty"`
	HandsFreeMode                            string               `json:"HandsFreeMode,omitempty"`
	HeadsetMode                              string               `json:"HeadsetMode,omitempty"`
	HoldOnTransferAttended                   string               `json:"HoldOnTransferAttended,omitempty"`
	HoldOnTransferUnattended                 string               `json:"HoldOnTransferUnattended,omitempty"`
	HttpAuthUsername                         string               `json:"HTTPAuthUsername,omitempty"`
	IPAddressType                            string               `json:"IPAddressType,omitempty"`
	IPv4Address                              string               `json:"IPv4Address,omitempty"`
	IPv4AlternateDNSServer                   string               `json:"IPv4AlternateDNSServer,omitempty"`
	IPv4PreferredDNSServer                   string               `json:"IPv4PreferredDNSServer,omitempty"`
	IPv4StandardGateway                      string               `json:"IPv4StandardGateway,omitempty"`
	IPv4SubnetMask                           string               `json:"IPv4SubnetMask,omitempty"`
	IncCallsWithoutCallManager               string               `json:"IncCallsWithoutCallManager,omitempty"`
	IncomingCall                             string               `json:"IncomingCall,omitempty"`
	LANPort                                  string               `json:"LANPort,omitempty"`
	LDAPAdditionalAttribute                  string               `json:"LDAPAdditionalAttribute,omitempty"`
	LDAPAdditionalAttributeDialable          string               `json:"LDAPAdditionalAttributeDialable,omitempty"`
	LDAPBaseDN                               string               `json:"LDAPBaseDN,omitempty"`
	LDAPCity                                 string               `json:"LDAPCity,omitempty"`
	LDAPCompany                              string               `json:"LDAPCompany,omitempty"`
	LDAPCountry                              string               `json:"LDAPCountry,omitempty"`
	LDAPDirectoryName                        string               `json:"LDAPDirectoryName,omitempty"`
	LDAPDisplayFormat                        string               `json:"LDAPDisplayFormat,omitempty"`
	LDAPEmail                                string               `json:"LDAPEmail,omitempty"`
	LDAPEnable                               string               `json:"LDAPEnable,omitempty"`
	LDAPFax                                  string               `json:"LDAPFax,omitempty"`
	LDAPFirstName                            string               `json:"LDAPFirstName,omitempty"`
	LDAPLookup                               string               `json:"LDAPLookup,omitempty"`
	LDAPMaxNumberOfSearchResults             int                  `json:"LDAPMaxNumberOfSearchResults,omitempty"`
	LDAPNameFilter                           string               `json:"LDAPNameFilter,omitempty"`
	LDAPNumberFilter                         string               `json:"LDAPNumberFilter,omitempty"`
	LDAPPAssword                             string               `json:"LDAPPassword,omitempty"`
	LDAPPhoneHome                            string               `json:"LDAPPhoneHome,omitempty"`
	LDAPPhoneMobile                          string               `json:"LDAPPhoneMobile,omitempty"`
	LDAPPhoneOffice                          string               `json:"LDAPPhoneOffice,omitempty"`
	LDAPResponseTimeout                      int                  `json:"LDAPResponseTimeout,omitempty"`
	LDAPSecurity                             string               `json:"LDAPSecurity,omitempty"`
	LDAPServerAddress                        string               `json:"LDAPServerAddress,omitempty"`
	LDAPServerPort                           int                  `json:"LDAPServerPort,omitempty"`
	LDAPStreet                               string               `json:"LDAPStreet,omitempty"`
	LDAPSurname                              string               `json:"LDAPSurname,omitempty"`
	LDAPUsername                             string               `json:"LDAPUsername,omitempty"`
	LDAPZIP                                  string               `json:"LDAPZIP,omitempty"`
	LLDAPActive                              string               `json:"LLDAPActive,omitempty"`
	LLDAPPacketInterval                      string               `json:"LLDAPPacketInterval,omitempty"`
	LinkSpeedDuplexLanPort                   string               `json:"LinkSpeedDuplexLanPort,omitempty"`
	LinkedSpeedDuplexPcPort                  string               `json:"LinkedSpeedDuplexPcPort,omitempty"`
	LogoutTimer                              int                  `json:"LogoutTimer,omitempty"`
	LookupIncoming                           string               `json:"LookupIncomming,omitempty"` // sic!
	LookupOutgoing                           string               `json:"LookupOutgoing,omitempty"`
	MACAddress                               string               `json:"MACAddress,omitempty"`
	MainMenuContent                          string               `json:"MainMenuContent,omitempty"`
	MenuAdaptiveJitterBuffer                 string               `json:"MenuAdaptiveJitterBuffer,omitempty"`
	MenuAdjustment                           string               `json:"MenuAdjustment,omitempty"`
	MenuAudo                                 string               `json:"MenuAudio,omitempty"`
	MenuCLI                                  string               `json:"MenuCLI,omitempty"`
	MenuCallDivert                           string               `json:"MenuCallDivert,omitempty"`
	MenuCallHistory                          string               `json:"MenuCallHistory,omitempty"`
	MenuCallstrings                          string               `json:"MenuCallstrings,omitempty"`
	MenuConnections                          string               `json:"MenuConnections,omitempty"`
	MenuCoreDumps                            string               `json:"MenuCoreDumps,omitempty"`
	MenuCorporate                            string               `json:"MenuCorporate,omitempty"`
	MenuDateAndTime                          string               `json:"MenuDateAndTime,omitempty"`
	MenuDebugging                            string               `json:"MenuDebugging,omitempty"`
	MenuDeskPhone                            string               `json:"MenuDeskPhone,omitempty"`
	MenuDevice                               string               `json:"MenuDevice,omitempty"`
	MenuDiallingPlans                        string               `json:"MenuDiallingPlans,omitempty"`
	MenuDisplay                              string               `json:"MenuDisplay,omitempty"`
	MenuDoNotDisturb                         string               `json:"MenuDoNotDisturb,omitempty"`
	MenuDoorStation                          string               `json:"MenuDoorStation,omitempty"`
	MenuEvents                               string               `json:"MenuEvents,omitempty"`
	MenuExpert                               string               `json:"menuExpert,omitempty"`
	MenuExpertAudio                          string               `json:"MenuExpertAudio,omitempty"`
	MenuExpertNetwork                        string               `json:"MenuExportNetwork,omitempty"`
	MenuExtensionModule1                     string               `json:"MenuExtensionModule1,omitempty"`
	MenuExtensionModule2                     string               `json:"MenuExtensionModule2,omitempty"`
	MenuExtensionModule3                     string               `json:"MenuExtensionModule3,omitempty"`
	MenuFirmwareUpdate                       string               `json:"MenuFirmwareUpdate,omitempty"`
	MenuIP                                   string               `json:"MenuIP,omitempty"`
	MenuKeysAndLEDs                          string               `json:"MenuKeysAndLEDs,omitempty"`
	MenuLan                                  string               `json:"MenuLAN`
	MenuLocalPhonebook                       string               `json:"MenuLocalPhonebook,omitempty"`
	MenuMainMenu                             string               `json:"MenuMainMenu,omitempty"`
	MenuMessageNotification                  string               `json:"MenuMessageNotification,omitempty"`
	MenuNetwork                              string               `json:"MenuNetwork,omitempty"`
	MenuOnlineDirectories                    string               `json:"MenuOnlineDirectories,omitempty"`
	MenuOnlineServices                       string               `json:"MenuOnlineServices',omitempty"`
	MenuPCAPLogging                          string               `json:"MenuPCAPLogging,omitempty"`
	MenuPhoneSystem                          string               `json:"MenuPhoneSystem,omitempty"`
	MenuPhoneWebServer                       string               `json:"MenuPhoneWebServer,omitempty"`
	MenuPictures                             string               `json:"MenuPictures,omitempty"`
	MenuProvisioningConfiguration            string               `json:"MenuProvisioningConfiguration,omitempty"`
	MenuPublic                               string               `json:"MenuPublic,omitempty"`
	MenuReboot                               string               `json:"MenuReboot,omitempty"`
	MenuRebootAndReset                       string               `json:"MenuRebootAndReset,omitempty"`
	MenuRingtones                            string               `json:"MenuRingtones,omitempty"`
	MenuSIPProtocol                          string               `json:"MenuSIPPRotocol,omitempty"`
	MenuSaveAndRestore                       string               `json:"MenuSaveAndRestore,omitempty"`
	MenuSecurity                             string               `json:"MenuSecurity,omitempty"`
	Menustrings                              string               `json:"Menustrings,omitempty"`
	MenuStatus                               string               `json:"MenuStatus,omitempty"`
	MenuStatusConnections                    string               `json:"MenuStatusConnections,omitempty"`
	MenuStorageAllocation                    string               `json:"MenuStorageAllocation,omitempty"`
	MenuSwitchstrings                        string               `json:"MenuSwitchstrings,omitempty"`
	MenuSystem                               string               `json:"MenuSystem,omitempty"`
	MenuSystemLogging                        string               `json:"MenuSystemLogging,omitempty"`
	MenuTelephony                            string               `json:"MenuTelephony,omitempty"`
	MenuVoIP                                 string               `json:"MenuVoIP,omitempty"`
	MenuVoiceMail                            string               `json:"MenuVoiceMail,omitempty"`
	MenuWebConfigurator                      string               `json:"MenuWebConfigurator,omitempty"`
	MenuWebcam                               string               `json:"MenuWebcam,omitempty"`
	MenuXML                                  string               `json:"MenuXML,omitempty"`
	MissedCallsNotificationActive            string               `json:"MissedCallsNotificationActive,omitempty"`
	NetworkType                              string               `json:"NetworkType,omitempty"`
	OffHook                                  string               `json:"OffHook,omitempty"`
	OnHook                                   string               `json:"OnHook,omitempty"`
	OneMelodyRingtoneDoorStation             string               `json:"OneMelodyRingtoneDoorStation,omitempty"`
	OneMelodyRingtoneExternal                string               `json:"OneMelodyRingtoneExternal,omitempty"`
	OneMelodyRingtoneGroup                   string               `json:"OneMelodyRingtoneGroup,omitempty"`
	OneMelodyRingtoneInterval                string               `json:"OneMelodyRingtoneInterval,omitempty"`
	OneMelodyRingtoneOptional                string               `json:"OneMelodyRingtoneOptional,omitempty"`
	OutCallsViaFunctionKey                   string               `json:"OutCallsViaFunctionKey,omitempty"`
	OutgoingCall                             string               `json:"OutgoingCall,omitempty"`
	PCPort                                   int                  `json:"PCPort,omitempty"`
	PIN                                      string               `json:"PIN,omitempty"`
	PacketTimeForRtp                         string               `json:"PacketTimeForRTP,omitempty"`
	PhoneLanguage                            string               `json:"PhoneLanguage,omitempty"`
	PhoneModel                               string               `json:"PhoneModel,omitempty"`
	PhoneName                                string               `json:"PhoneName,omitempty"`
	PhoneSystem                              string               `json:"PhoneSystem,omitempty"`
	ProgrammableKeys                         string               `json:"ProgrammableKeysDNDActionURLDisable,omitempty"`
	ProgrammableKeysConferenceActionURL      string               `json:"ProgrammableKeysConferenceActionURL,omitempty"`
	ProgrammableKeysConferenceDTMFCode       string               `json:"ProgrammableKeysDTMFCode,omitempty"`
	ProgrammableKeysConferenceFAC            string               `json:"ProgrammableKeysConferenceFAC,omitempty"`
	ProgrammableKeysConferenceType           string               `json:"ProgrammableKeysConferenceType,omitempty"`
	ProgrammableKeysDNDActionURLEnable       string               `json:"ProgrammableKeysDNDActionURLEnable,omitempty"`
	ProgrammableKeysDNDFACDisable            string               `json:"ProgrammableKeysDNDFACDisable,omitempty"`
	ProgrammableKeysDNDType                  string               `json:"ProgrammableKeysDNDType,omitempty"`
	ProgrammableKeysDirectoryType            string               `json:"ProgrammableKeysDirectoryType,omitempty"`
	ProgrammableKeysHoldActionURL            string               `json:"ProgrammableKeysHoldActionURL,omitempty"`
	ProgrammableKeysHoldDTMFCode             string               `json:"ProgrammableKeysHoldDTMFCode,omitempty"`
	ProgrammableKeysHoldType                 string               `json:"ProgrammableKeysHoldType,omitempty"`
	ProgrammableKeysMessagesActionURL        string               `json:"ProgrammableKeysMessagesActionURL,omitempty"`
	ProgrammableKeysMessagesFAC              string               `json:"ProgrammableKeysMessagesFAC,omitempty"`
	ProgrammableKeysMessagesType             string               `json:"ProgrammableKeysMessagesType,omitempty"`
	ProvisioningDirectLink                   string               `json:"ProvisioningDirectLink,omitempty"`
	ProvisioningServer                       string               `json:"ProvisioningServer,omitempty"`
	ProxyServerActive                        string               `json:"ProxyServerActive,omitempty"`
	ProxyServerAddress                       string               `json:"ProxyServerAddress,omitempty"`
	ProxyServerPort                          int                  `json:"ProxyServerPort,omitempty"`
	QuickDialKeys                            []QuickDialKey       `json:"QuickDialKeys,omitempty"`
	RTPQoSDSCP                               int                  `json:"RTPoSDSCP,omitempty"`
	RegistrationFailed                       string               `json:"RegistrationFailed,omitempty"`
	RegistrationSucceeded                    string               `json:"RegistrationSucceeded,omitempty"`
	RemoteControlAllow                       string               `json:"RemoteControlAllow,omitempty"`
	RemoteControlSource                      string               `json:"RemoteControlSource,omitempty"`
	SIPAccountFailover                       string               `json:"SIPAccountFailover,omitempty"`
	SIPG729AnnexB                            string               `json:"SIPG729AnnexB,omitempty"`
	SIPNoSrtpCalls                           string               `json:"SipNoSrtpCalls,omitempty"`
	SIPPort                                  int                  `json:"SIPPort,omitempty"`
	SIPPrack                                 string               `json:"SIPPRack,omitempty"`
	SIPQoSDSCP                               int                  `json:"SIPQoSDSCP,omitempty"`
	SIPRtpPort                               int                  `json:"SIPRtpPort,omitempty"`
	SIPRtpRTCPXRServerAddress                string               `json:"SIPRtpRTCPXRServerAddress,omitempty"`
	SIPRtpRTCPXRServerPort                   int                  `json:"SIPRtpRTCXRServerPort,omitempty"`
	SIPRtpRandomPort                         string               `json:"SipRtpRandomPort,omitempty"`
	SIPRtpSymmetricPort                      string               `json:"SIPRtpSymetricPort,omitempty"` // sic!
	SIPRtpUseRTCPXR                          string               `json:"SIPRtpUseRTCPXR,omitempty"`
	SIPRtprRtcp                              string               `json:"SIPRtprRtcp,omitempty"`
	SIPSCertificate                          string               `json:"SIPSCertificate,omitempty"`
	SIPSKeyPassword                          string               `json:"SIPSKeyPassword,omitempty"`
	SIPSPrivateKey                           string               `json:"SIPSPrivateKey,omitempty"`
	SIPSecurity                              string               `json:"SIPSecurity,omitempty"`
	SIPSecurityEnabled                       string               `json:"SIPSecurityEnabled,omitempty"`
	SIPSessionTimer                          int                  `json:"SIPSessionTimer,omitempty"`
	SIPTimerT1                               int                  `json:"SIPTimerT1,omitempty"`
	SIPTimersFailedRegistration              int                  `json:"SIPTimersFailedRegistration,omitempty"`
	SIPTimersFailedSubscription              int                  `json:"SIPTimersFailSubscription,omitempty"`
	SIPTimersSubscription                    int                  `json:"SIPTimersSubscription,omitempty"`
	SIPTimersSubscriptionBLFFollowRegister   string               `json:"SIPTimersSubscriptionBLFFollowRegister,omitempty"`
	SIPTransportProtocol                     string               `json:"SIPTransportProtocol,omitempty"`
	SIPrtp                                   string               `json:"SIPrtp,omitempty"`
	ScreenSaverTimeout                       string               `json:"ScreenSaverTimeout,omitempty"`
	Screensaver                              string               `json:"Screensaver,omitempty"`
	ScreensaverBacklight                     int                  `json:"ScreensaverBacklight,omitempty"`
	ScreensaverHTTPSource                    string               `json:"ScreensaverHTTPSource,omitempty"`
	ScreensaverPictures                      string               `json:"ScreensaverPictures,omitempty"`
	SelectedCodecs                           []int                `json:"SelectedCodecs,omitempty"`
	SelectedServesDisable                    string               `json:"SelectedServicesDisable,omitempty"`
	SemiAttendedTransferType                 string               `json:"SemiAttendedTransferType,omitempty"`
	StringsVersion                           string               `json:"StringsVersion,omitempty"`
	ShowPIN                                  string               `json:"ShowPIN,omitempty"`
	ShowPassword                             string               `json:"ShowPassword,omitempty"`
	ShowSIPMessagesOnDisplay                 string               `json:"ShowSIPMessagesOnDisplay,omitempty"`
	Sip                                      []Sip                `json:"SIP,omitempty"`
	SoftReboots                              int                  `json:"SoftReboots,omitempty"`
	SoftwareVariant                          string               `json:"SoftwareVariant,omitempty"`
	SoftwareVersion                          string               `json:"SoftwareVersion,omitempty"`
	StandbyBacklight                         int                  `json:"StandbyBacklight,omitempty"`
	Startups                                 int                  `json:"Startups,omitempty"`
	SyslogEnabled                            string               `json:"SyslogEnabled,omitempty"`
	SyslogServer                             string               `json:"SyslogServer,omitempty"`
	SystemLocalPhonebookUpdateTime           string               `json:"SystemLocalPhonebookUpdateTime,omitempty"`
	SystemLocalPhonebookUrl                  string               `json:"SystemLocalPhonebookUrl,omitempty"`
	TimeFormat                               string               `json:"TimeFormat,omitempty"`
	TimeServer                               string               `json:"TimeServer,omitempty"`
	TimeServerDHCP                           string               `json:"TimeServerDHCP,omitempty"`
	TimeServerProvisioning                   string               `json:"TimeServerProvisioning,omitempty"`
	TimeZone                                 string               `json:"TimeZone,omitempty"`
	Timestamp                                int                  `json:"Timestamp,omitempty"`
	ToneScheme                               string               `json:"TomeScheme,omitempty"`
	UserPassword                             string               `json:"UserPassword,omitempty"`
	VLANIdentifierLAN                        int                  `json:"VLANIdentifierLAN,omitempty"`
	VLANIdentifierPC                         int                  `json:"VLANIdentifierPC,omitempty"`
	VLANLocked                               string               `json:"VLANLocked,omitempty"`
	VLANPriorityLAN                          string               `json:"VLANPriorityLAN,omitempty"`
	VLANPriorityPC                           string               `json:"VLANPriorityDC,omitempty"`
	VLANTagging                              string               `json:"VLANTagging,omitempty"`
	Variant                                  string               `json:"Variant,omitempty"`
	VoiceQuality                             string               `json:"VoiceQuality,omitempty"`
	VoicemailMessagesActive                  string               `json:"VoicemailMessagesActive,omitempty"`
	WebUICallDivertDisable                   string               `json:"WebUICallDivertDisable,omitempty"`
	WebUICallWaitingDisable                  string               `json:"WebUICallWaitingDisable,omitempty"`
	WebUILanguage                            string               `json:"WebUILanguage,omitempty"`
	WitholdNumberDisable                     string               `json:"WitholdNumberDisable,omitempty"`
	WorkingCounter                           int                  `json:"WorkingCounter,omitempty"`
	WorkingCounterSec                        int                  `json:"WorkingCounterSec,omitempty"`
	XMLProviderName                          string               `json:"XMLProviderName,omitempty"`
	XSIAuthName                              string               `json:"XSIAuthName,omitempty"`
	XSIAuthPassword                          string               `json:"XSIAuthPassword,omitempty"`
	XSICallLogType                           string               `json:"XSICallLogType,omitempty"`
	XSIEnterpriseCommonDirectoryEnabled      string               `json:"XSIEnterpriseCommonDirectoryEnabled,omitempty"`
	XSIEnterpriseCommonDirectoryName         string               `json:"XSIEnterpriseCommonDirectoryName,omitempty"`
	XSIEnterpriseDirectoryEnabled            string               `json:"XSIEnterpriseDirectoryEnabled,omitempty"`
	XSIEnterpriseDirectoryName               string               `json:"XSIEnterpriseDirectoryName,omitempty"`
	XSIGroupCommonDirectoryEnabled           string               `json:"XSIGroupCommonDirectoryEnabled,omitempty"`
	XSIGroupCommonDirectoryName              string               `json:"XSIGroupCommonDirectoryName,omitempty"`
	XSIGroupDirectoryEnabled                 string               `json:"XSIGroupDirectoryEnabled,omitempty"`
	XSIGroupDirectoryName                    string               `json:"XSIGroupDirectoryName,omitempty"`
	XSIPersonalDirectoryEnabled              string               `json:"XSIPersonalDirectoryEnabled,omitempty"`
	XSIPersonalDirectoryName                 string               `json:"XSIPersonalDirectoryName,omitempty"`
	XSISIPAuthentication                     string               `json:"XSISIPAuthentication,omitempty"`
	XSISearchAnywhereInNameEnabled           string               `json:"XSISearchAnywhereInNameEnabled,omitempty"`
	XSIServer                                string               `json:"XSIServer,omitempty"`
	XmlEnablePrivateDirectory                string               `json:"XmlEnablePrivateDirectory,omitempty"`
	XmlEnableWhiteDirectory                  string               `json:"XmlEnableWhiteDirectory,omitempty"`
	XmlEnableYellowDirectory                 string               `json:"XmlEnableYellowDirectory,omitempty"`
	XmlNumberFilter                          string               `json:"XmlNumberFilter,omitempty"`
	XmlPassword                              string               `json:"XMLPassword,omitempty"`
	XmlPrivateDirectoryName                  string               `json:"XMLPrivateDirectoryName,omitempty"`
	XmlServerAddress                         string               `json:"XmlServerAddress,omitempty"`
	XmlUsername                              string               `json:"XmlUsername,omitempty"`
	XmlWhiteDirectoryName                    string               `json:"XmlWhiteDirectoryName,omitempty"`
	XmlYellowDirectoryName                   string               `json:"XmlYellowDirectoryName,omitempty"`
}

func (p *Parameters) TransformFunctionKeyNames(original, replace string) (Parameters, []int) {
	keys := make([]FunctionKey, 0, len(p.FunctionKeys))
	changed := make([]int, 0, 0)
	for index, fnKey := range p.FunctionKeys {
		var key = FunctionKey{}
		if fnKey.DisplayName == original {
			key = FunctionKey{DisplayName: replace}
			changed = append(changed, index)
		}
		keys = append(keys, key)
	}
	return Parameters{FunctionKeys: keys}, changed
}

func (p *Parameters) UnmarshalJSON(data []byte) error {
	return unmarshalInternal(data, p)
}

type FunctionKeys []FunctionKey

type FunctionKey struct {
	AutomaticallyFilled string `json:"AutomaticallyFilled,omitempty"`
	CallDivertType      string `json:"CallDivertType,omitempty"`
	CallPickupCode      string `json:"CallPickupCode,omitempty"`
	Color               string `json:"Color,omitempty"`
	Connection          string `json:"Connection,omitempty"`
	DTMFCode            string `json:"DTMFCode,omitempty"`
	DisableCode         string `json:"DisableCode,omitempty"`
	DisplayName         string `json:"DisplayName,omitempty"`
	EnableCode          string `json:"EnableCode,omitempty"`
	LockProvisioning    string `json:"LockProvisioning,omitempty"`
	PhoneNumber         string `json:"PhoneNumber,omitempty"`
	Silent              string `json:"Silent,omitempty"`
	Type                string `json:"Type,omitempty"`
	Url                 string `json:"URL,omitempty"`
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

func (f *FunctionKey) UnmarshalJSON(data []byte) error {
	return unmarshalInternal(data, f)
}

type Sip struct {
	AccountName                   string `json:"AccountName,omitempty"`
	Active                        string `json:"Active,omitempty"`
	AllowRouteHeaders             string `json:"AllowRouteHeaders,omitempty"`
	AuthenaticationName           string `json:"AuthenticationName,omitempty"`
	AuthenticationPassword        string `json:"AuthenticationPassword,omitempty"`
	AutoNegOfDTMFTransmission     string `json:"AutoNetOfDTMSTransmission,omitempty"`
	CLIPSource                    string `json:"CLIPSource,omitempty"`
	CLIR                          string `json:"CLIR,omitempty"`
	CallWaiting                   string `json:"CallWaiting,omitempty"`
	CallWaitingSignal             string `json:"CallWaitingSignal,omitempty"`
	CountMissedAcceptedCalls      string `json:"CountMissedAcceptedCalls,omitempty"`
	DNSQuery                      string `json:"DNSQuery,omitempty"`
	DTMFTransmission              string `json:"DTMFTransmission,omitempty"`
	DisplayName                   string `json:"DisplayName,omitempty"`
	Domain                        string `json:"Domain,omitempty"`
	FailoverServerAddress         string `json:"FailoverServerAddress,omitempty"`
	FailoverServerEnabled         string `json:"FailoverServerEnabled,omitempty"`
	FailoverServerPort            int    `json:"FailoverServerPort,omitempty"`
	HeaderDoorstation             string `json:"HeaderDoorstation,omitempty"`
	HeaderExternal                string `json:"HeaderExternal,omitempty"`
	HeaderGroup                   string `json:"HeaderGroup,omitempty"`
	HeaderInternal                string `json:"HeaderInternal,omitempty"`
	HeaderOptional                string `json:"HeaderOptional,omitempty"`
	ICE                           string `json:"ICE,omitempty"`
	NATRefreshTime                int    `json:"NATRefreshTime,omitempty"`
	OutboundProxyAddress          string `json:"OutboundProxyAddress,omitempty"`
	OutboundProxyMode             string `json:"OutboundProxyMode,omitempty"`
	OutboundProxyPort             int    `json:"OutboundProxyPort,omitempty"`
	Provider                      string `json:"Provider,omitempty"`
	ProxyServerAddress            string `json:"ProxyServerAddress,omitempty"`
	ProxyServerPort               int    `json:"ProxyServerPort,omitempty"`
	RegistrationServerAddress     string `json:"RegistrationServerAddress,omitempty"`
	RegistrationServerPort        int    `json:"RegistrationSeverPort,omitempty"`
	RegistrationServerRefreshTiem int    `json:"RegistrationServerRefreshTime,omitempty"`
	RequestCheckOptions           int    `json:"RequestCheckOptions,omitempty"`
	ReregisterAlternative         string `json:"ReregisterAlternative,omitempty"`
	RingtoneDoorStation           string `json:"RingtoneDoorStation,omitempty"`
	RingtoneExternal              string `json:"RingtoneExternal,omitempty"`
	RingtoneGroup                 string `json:"RingtoneGroup,omitempty"`
	RingtoneInternal              string `json:"RingtoneInternal,omitempty"`
	RingtoneOptional              string `json:"RingtoneOptional,omitempty"`
	STUNEnabled                   string `json:"STUNEnabled,omitempty"`
	STUNRefreshTime               int    `json:"STUNRefreshTime,omitempty"`
	STUNServerAddress             string `json:"STUNServerAddress,omitempty"`
	STUNServerPort                int    `json:"STUNServerPort,omitempty"`
	Username                      string `json:"Username,omitempty"`
	VoiceMailActive               string `json:"VoiceMailActive,omitempty"`
	VoiceMailMailbox              string `json:"VoiceMailMailbox,omitempty"`
}

func (s *Sip) UnmarshalJSON(data []byte) error {
	return unmarshalInternal(data, s)
}

type Dnd struct {
	PhoneNumber string `json:"PhoneNumber,omitempty"`
	Name        string `json:"Name,omitempty"`
}

func (d *Dnd) UnmarshalJSON(data []byte) error {
	return unmarshalInternal(data, d)
}

type CallDivertAll struct {
	TargetMail string `json:"TargetMail,omitempty"`
	Target     string `json:"Target,omitempty"`
	VoiceMail  string `json:"VoiceMail,omitempty"`
	Active     string `json:"Active,omitempty"`
}

func (c *CallDivertAll) UnmarshalJSON(data []byte) error {
	return unmarshalInternal(data, c)
}

type QuickDialKey struct {
	Type      string `json:"Type,omitempty"`
	Number    string `json:"Number,omitempty"`
	FAC       string `json:"FAC,omitempty"`
	ActionURL string `json:"ActionURL,omitempty"`
}

func (q *QuickDialKey) UnmarshalJSON(data []byte) error {
	return unmarshalInternal(data, q)
}

type DoorStation struct {
	Password           string `json:"Password,omitempty"`
	Username           string `json:"Username,omitempty"`
	DTMFCode           string `json:"DTMFCode,omitempty"`
	CameraURL          string `json:"CameraURL,omitempty"`
	Name               string `json:"Name,omitempty"`
	PictureRefreshTime string `json:"PictureRefreshTime,omitempty"`
	SIPID              string `json:"SIPID,omitempty"`
}

func (d *DoorStation) UnmarshalJSON(data []byte) error {
	return unmarshalInternal(data, d)
}

type CallDivertNoAnswer struct {
	CallDivertBusy
	Delay string `json:"Delay,omitempty"`
}

func (c *CallDivertNoAnswer) UnmarshalJSON(data []byte) error {
	return unmarshalInternal(data, c)
}

type CallDivertBusy struct {
	VoiceMail  string `json:"VoiceMail,omitempty"`
	Active     string `json:"Active,omitempty"`
	TargetMail string `json:"TargetMail,omitempty"`
	Target     string `json:"Target,omitempty"`
}

func (c *CallDivertBusy) UnmarshalJSON(data []byte) error {
	return unmarshalInternal(data, c)
}

type DialingPlan struct {
	PhoneNumber string `json:"PhoneNumber,omitempty"`
	Comment     string `json:"Comment,omitempty"`
	Active      string `json:"Active,omitempty"`
	Connection  string `json:"Connection,omitempty"`
	UseAreaCode string `json:"UseAreaCode,omitempty"`
}

func (d *DialingPlan) UnmarshalJSON(data []byte) error {
	return unmarshalInternal(data, d)
}

func unmarshalInternal(data []byte, v interface{}) error {
	raw := make(map[string]interface{})
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	structValue := reflect.ValueOf(v)
	for index := 0; index < structValue.Elem().NumField(); index++ {
		field := structValue.Elem().Field(index)
		jsonName := jsonFieldName(structValue.Elem().Type().Field(index))
		if value, ok := raw[jsonName]; ok {
			if reflect.TypeOf(value).Kind() == reflect.Map {
				value = value.(map[string]interface{})["value"]
			}
			switch field.Kind() {
			case reflect.String:
				converted, ok := value.(string)
				if !ok {
					return fmt.Errorf("json property \"%s\" has wrong type, want \"string\", found \"%v\"", jsonName, reflect.TypeOf(value))
				}
				field.SetString(converted)
				break
			case reflect.Int:
				converted, ok := value.(float64)
				if !ok {
					return fmt.Errorf("json property \"%s\" has wrong type, want \"float64\", found \"%v\"", jsonName, reflect.TypeOf(value))
				}
				field.SetInt(int64(converted))
				break
			case reflect.Slice:
				marshal, _ := json.Marshal(value)
				target := reflect.New(field.Type()).Interface()
				err := json.Unmarshal(marshal, target)
				if err != nil {
					return fmt.Errorf("could not unmarshal field of type \"%v\" with name \"%s\": %v", field.Type(), jsonName, err)
				}
				field.Set(reflect.ValueOf(target).Elem())
				break
			default:
				log.Fatalf("type conversion for field \"%s\" and type \"%v\" not implemented", field.Type().Name(), field.Type())
			}
		}
	}
	return nil
}

func jsonFieldName(field reflect.StructField) string {
	lookup, ok := field.Tag.Lookup("json")
	if !ok {
		return field.Name
	}
	return strings.Split(lookup, ",")[0]
}
