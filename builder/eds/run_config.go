//go:generate packer-sdc struct-markdown
//go:generate packer-sdc mapstructure-to-hcl2 -type EdsUser,EdsOfficeSiteInternetAccess,EdsOfficeSiteCen,EdsOfficeSite,EdsImageFilter,EdsComputerTemplate,EdsPolicyGroup,EdsUserCommand,EdsArtifact,RunConfig

package eds

import (
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	"github.com/hashicorp/packer-plugin-sdk/uuid"
)

type EdsUser struct {
	Name         string `mapstructure:"name" required:"true"`
	Email        string `mapstructure:"email" required:"true"`
	Phone        string `mapstructure:"phone" required:"false"`
	OwnerType    string `mapstructure:"owner_type" required:"false"`
	OrgId        string `mapstructure:"org_id" required:"false"`
	Remark       string `mapstructure:"remark" required:"false"`
	RealNickName string `mapstructure:"real_nick_name" required:"false"`
}

type EdsOfficeSiteInternetAccess struct {
	Enabled   bool  `mapstructure:"enabled" required:"false"`
	Bandwidth int32 `mapstructure:"bandwidth" required:"false"`
}

type EdsOfficeSiteCen struct {
	Id         string `mapstructure:"id" required:"false"`
	OwnerId    int64  `mapstructure:"owner_id" required:"false"`
	VerifyCode string `mapstructure:"verify_code" required:"false"`
}

type EdsOfficeSite struct {
	Id             string                      `mapstructure:"id" required:"false"`
	Name           string                      `mapstructure:"name" required:"false"`
	CidrBlock      string                      `mapstructure:"cidr_block" required:"false"`
	InternetAccess EdsOfficeSiteInternetAccess `mapstructure:"internet_access" required:"false"`
	Cen            EdsOfficeSiteCen            `mapstructure:"cen" required:"false"`
}

type EdsImageFilter struct {
	// The ID of the image.
	ImageId string `mapstructure:"image_id" required:"false"`
	// The image name.
	ImageName string `mapstructure:"image_name" required:"false"`
	// The image version.
	ImageVersion string `mapstructure:"image_version" required:"false"`
	// The type of the image. Avaiable options: [SYSTEM, CUSTOM]
	ImageType string `mapstructure:"image_type" required:"false"`
	// The status of the image. Avaiable options: [Creating, Available,
	// CreateFailed]
	ImageStatus string `mapstructure:"image_status" required:"false"`
	// The protocol type. Available options: [HDX, ASP]
	ProtocolType string `mapstructure:"protocol_type" required:"false"`
	// The instance type of the cloud computer. You can call the
	// DescribeDesktopTypes operation to obtain the parameter value.
	InstanceType string `mapstructure:"desktop_instance_type" required:"false"`
	// The session type. Avaiable options: [SINGLE_SESSION, MULTIPLE_SESSION]
	SessionType string `mapstructure:"session_type" required:"false"`
	// Specifies whether the images are GPU-accelerated images. Available
	// options: [true, false]
	GpuCategory bool `mapstructure:"gpu_category" required:"false"`
	// The version of the GPU driver.
	GpuDriverVersion string `mapstructure:"gpu_driver_version" required:"false"`
	// The type of the operating system of the images.Available options:
	// [Linux, Windows]
	OsType string `mapstructure:"os_type" required:"false"`
	// The language of the OS. Available options: [en-US, zh-HK, zh-CN, ja-JP]
	OsLanguageType string `mapstructure:"language_type" required:"false"`
}

type EdsComputerTemplate struct {
	Id                       string         `mapstructure:"id" required:"false"`
	Name                     string         `mapstructure:"name" required:"false"`
	Description              string         `mapstructure:"description" required:"false"`
	SourceImageFilter        EdsImageFilter `mapstructure:"source_image_filter" required:"false"`
	InstanceType             string         `mapstructure:"instance_type" required:"true"`
	RootDiskSizeGib          int            `mapstructure:"root_disk_size_gib" required:"true"`
	RootDiskPerformanceLevel string         `mapstructure:"root_disk_performance_level" required:"false"`
	UserDiskSizeGib          []int32        `mapstructure:"user_disk_size_gib" required:"false"`
	UserDiskPerformanceLevel string         `mapstructure:"user_disk_performance_level" required:"false"`
	Language                 string         `mapstructure:"language" required:"false"`
}

type EdsPolicyGroup struct {
	// Invitational Preview ONLY
	// AdminAccess                string `mapstructure:"admin_access" required:"false"`

	Id   string `mapstructure:"id" required:"false"`
	Name string `mapstructure:"name" required:"false"`
	// VisualQuality                 string `mapstructure:"visual_quality" required:"false"`
	// GpuAcceleration               string `mapstructure:"gpu_acceleration" required:"false"`
	// AppContentProtection          bool   `mapstructure:"app_content_protection" required:"false"`
	// InternetCommunicationProtocol string `mapstructure:"internet_communication_protocol" required:"false"`
	// MaxReconnectTime              int    `mapstructure:"max_reconnect_time" required:"false"`
	// WyAssistant                   bool   `mapstructure:"wy_assistant" required:"false"`

	// AuthorizeSecurityPolicyRule []struct {
	// 	Type        string `mapstructure:"type" required:"false"`
	// 	Policy      string `mapstructure:"policy" required:"false"`
	// 	PortRange   string `mapstructure:"port_range" required:"false"`
	// 	Description string `mapstructure:"description" required:"false"`
	// 	IpProtocol  string `mapstructure:"ip_protocol" required:"false"`
	// 	Priority    string `mapstructure:"priority" required:"false"`
	// 	CidrIp      string `mapstructure:"cidr_ip" required:"false"`
	// } `mapstructure:"authorize_security_policy_rule" required:"false"`

	// AuthorizeAccessPolicyRule []struct {
	// 	Description string `mapstructure:"description" required:"false"`
	// 	CidrIp      string `mapstructure:"cidr_ip" required:"false"`
	// } `mapstructure:"authorize_access_policy_rule" required:"false"`

	// ClientType []struct {
	// 	Enabled    bool   `mapstructure:"enabled" required:"false"`
	// 	ClientType string `mapstructure:"client_type" required:"false"`
	// } `mapstructure:"client_type" required:"false"`

	// Preemption struct {
	// 	PreemptLogin     bool     `mapstructure:"preempt_login" required:"false"`
	// 	PreemptLoginUser []string `mapstructure:"preempt_login_user" required:"false"`
	// } `mapstructure:"preemption" required:"false"`

	// Redirection struct {
	// 	Camera     bool   `mapstructure:"camera" required:"false"`
	// 	Clipboard  string `mapstructure:"clipboard" required:"false"`
	// 	LocalDrive string `mapstructure:"local_drive" required:"false"`
	// 	Network    bool   `mapstructure:"network" required:"false"`
	// 	Printer    bool   `mapstructure:"printer" required:"false"`
	// 	USB        bool   `mapstructure:"usb" required:"false"`
	// 	Video      bool   `mapstructure:"video" required:"false"`

	// 	Devices []struct {
	// 		DeviceType   string `mapstructure:"type" required:"false"`
	// 		RedirectType string `mapstructure:"redirect_type" required:"false"`
	// 	} `mapstructure:"devices" required:"false"`

	// 	DeviceRules []struct {
	// 		Type         string `mapstructure:"type" required:"false"`
	// 		Name         string `mapstructure:"name" required:"false"`
	// 		VendorId     string `mapstructure:"vendor_id" required:"false"`
	// 		ProductId    string `mapstructure:"product_id" required:"false"`
	// 		RedirectType string `mapstructure:"redirect_type" required:"false"`
	// 		OptCommand   string `mapstructure:"opt_command" required:"false"`
	// 		Platforms    string `mapstructure:"platforms" required:"false"`
	// 	} `mapstructure:"device_rules" required:"false"`

	// 	USBSupplyRules []struct {
	// 		VendorId       string `mapstructure:"vendor_id" required:"false"`
	// 		ProductId      string `mapstructure:"product_id" required:"false"`
	// 		Description    string `mapstructure:"description" required:"false"`
	// 		RedirectType   int    `mapstructure:"redirect_type" required:"false"`
	// 		DeviceClass    string `mapstructure:"device_class" required:"false"`
	// 		DeviceSubclass string `mapstructure:"device_subclass" required:"false"`
	// 		RuleType       int    `mapstructure:"rule_type" required:"false"`
	// 	} `mapstructure:"usb_supply_rules" required:"false"`
	// } `mapstructure:"redirection" required:"false"`

	// DomainList    string `mapstructure:"domain_list" required:"false"`
	// DomainResolve struct {
	// 	Enabled           bool `mapstructure:"enabled" required:"false"`
	// 	DomainResolveRule []struct {
	// 		Domain      string `mapstructure:"domain" required:"false"`
	// 		Policy      string `mapstructure:"policy" required:"false"`
	// 		Description string `mapstructure:"description" required:"false"`
	// 	} `mapstructure:"domain_resolve_rule" required:"false"`
	// } `mapstructure:"domain_resolve" required:"false"`

	// HelpDesk struct {
	// 	RemoteCoordinate            string `mapstructure:"remote_coordinate" required:"false"`
	// 	EndUserApplyAdminCoordinate string `mapstructure:"end_user_apply_admin_coordinate" required:"false"`
	// 	EndUserGroupCoordinate      string `mapstructure:"end_user_group_coordinate" required:"false"`
	// } `mapstructure:"help_desk" required:"false"`

	// HTML5 struct {
	// 	Enabled      bool   `mapstructure:"enabled" required:"false"`
	// 	FileTransfer string `mapstructure:"file_transfer" required:"false"`
	// } `mapstructure:"html5" required:"false"`

	// Recording struct {
	// 	Enabled           bool   `mapstructure:"enabled" required:"false"`
	// 	StartTime         string `mapstructure:"start_time" required:"false"`
	// 	EndTime           string `mapstructure:"end_time" required:"false"`
	// 	Fps               int    `mapstructure:"fps" required:"false"`
	// 	Audio             string `mapstructure:"audio" required:"false"`
	// 	Expires           int    `mapstructure:"expires" required:"false"`
	// 	Content           string `mapstructure:"content" required:"false"`
	// 	ContentExpires    int    `mapstructure:"content_expires" required:"false"`
	// 	Duration          int    `mapstructure:"duration" required:"false"`
	// 	UserNotify        string `mapstructure:"user_notify" required:"false"`
	// 	UserNotifyMessage string `mapstructure:"user_notify_message" required:"false"`
	// } `mapstructure:"recording" required:"false"`

	// Scope struct {
	// 	ScopeType  string   `mapstructure:"type" required:"false"`
	// 	ScopeValue []string `mapstructure:"value" required:"false"`
	// } `mapstructure:"scope" required:"false"`

	// Watermark struct {
	// 	Enabled           bool    `mapstructure:"enabled" required:"false"`
	// 	Type              string  `mapstructure:"type" required:"false"`
	// 	Transparency      string  `mapstructure:"transparency" required:"false"`
	// 	TransparencyValue int     `mapstructure:"transparency_value" required:"false"`
	// 	Color             int     `mapstructure:"color" required:"false"`
	// 	Degree            float64 `mapstructure:"degree" required:"false"`
	// 	FontSize          int     `mapstructure:"font_size" required:"false"`
	// 	FontStyle         string  `mapstructure:"font_style" required:"false"`
	// 	RowAmount         int     `mapstructure:"row_amount" required:"false"`
	// 	Security          string  `mapstructure:"security" required:"false"`
	// 	AntiCam           string  `mapstructure:"anti_cam" required:"false"`
	// 	Power             string  `mapstructure:"power" required:"false"`
	// } `mapstructure:"watermark" required:"false"`
}

type EdsUserCommand struct {
	Type      string `mapstructure:"type" required:"true"`
	Content   string `mapstructure:"content" required:"true"`
	Encoding  string `mapstructure:"encoding" required:"true"`
	EndUserId string `mapstructure:"end_user_id" required:"false"`
	Role      string `mapstructure:"role" required:"false"`
	Timeout   uint64 `mapstructure:"timeout" required:"false"`
}
type EdsUserCommands []EdsUserCommand

type EdsArtifact struct {
	ImageName        string `mapstructure:"image_name" required:"true"`
	ImageDescription string `mapstructure:"description" required:"false"`
}

type RunConfig struct {
	Comm communicator.Config `mapstructure:",squash"`

	ResourceGroupId         string              `mapstructure:"resource_group_id" required:"false"`
	ComputerPoolId          string              `mapstructure:"computer_pool_id" required:"false"`
	DesktopName             string              `mapstructure:"desktop_name" required:"false"`
	DesktopNameSuffix       bool                `mapstructure:"desktop_name_suffix" required:"false"`
	Hostname                string              `mapstructure:"hostname" required:"false"`
	DesktopIp               string              `mapstructure:"desktop_ip" required:"false"`
	VolumeEncryptionEnabled bool                `mapstructure:"volume_encryption_enabled" required:"false"`
	VolumeEncryptionKey     string              `mapstructure:"volume_encryption_key" required:"false"`
	EndUser                 EdsUser             `mapstructure:"end_user" required:"false"`
	OfficeSite              EdsOfficeSite       `mapstructure:"office_site" required:"false"`
	ComputerTemplate        EdsComputerTemplate `mapstructure:"computer_template" required:"false"`
	PolicyGroup             EdsPolicyGroup      `mapstructure:"policy_group" required:"false"`
	UserCommands            EdsUserCommands     `mapstructure:"user_commands" required:"false"`
	Artifact                EdsArtifact         `mapstructure:"artifact" required:"true"`
}

func (c *RunConfig) Prepare(ctx *interpolate.Context) []error {
	// If we are not given an explicit ssh_keypair_name or
	// ssh_private_key_file, then create a temporary one, but only if the
	// temporary_key_pair_name has not been provided and we are not using
	// ssh_password.
	if c.Comm.SSHKeyPairName == "" && c.Comm.SSHTemporaryKeyPairName == "" &&
		c.Comm.SSHPrivateKeyFile == "" && c.Comm.SSHPassword == "" {

		c.Comm.SSHTemporaryKeyPairName = fmt.Sprintf("packer-%s", uuid.TimeOrderedUUID())
	}

	if c.Comm.SSHUsername == "" {
		// AliCloud default username is root for Linux image.
		c.Comm.SSHUsername = "root"
	}

	// IMPORTANT: Trigger communicator preparation to ensure that the
	// communicator is valid
	errs := c.Comm.Prepare(ctx)

	//
	if c.Artifact.ImageName == "" {
		errs = append(errs, fmt.Errorf("artifact.image_name must be specified"))
	}
	//
	if c.ComputerTemplate.SourceImageFilter.ImageId == "" &&
		c.ComputerTemplate.SourceImageFilter.ImageName == "" {
		errs = append(errs, fmt.Errorf("source_image_filter.image_id or source_image_filter.image_name must be specified"))
	}
	//
	if c.ComputerTemplate.RootDiskSizeGib < 40 {
		errs = append(errs, fmt.Errorf("root_disk_size_gib must be at least 40"))
	}
	//
	if len(c.ComputerTemplate.UserDiskSizeGib) <= 0 {
		errs = append(errs, fmt.Errorf("user_disk_size_gib must be specified"))
	}
	for i, size := range c.ComputerTemplate.UserDiskSizeGib {
		if size < 40 {
			errs = append(errs, fmt.Errorf("user_disk_size_gib[%d] must be at least 40", i))
		}
	}
	//
	if c.VolumeEncryptionEnabled && c.VolumeEncryptionKey == "" {
		errs = append(errs, fmt.Errorf("artifact.volume_encryption_key must be specified when volume_encryption_enabled is true"))
	}

	return errs
}
