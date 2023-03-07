package oktawave

import (
	"context"
	"github.com/oktawave-code/odk"
	oks "github.com/oktawave-code/oks-sdk"
)

type ClientConfig struct {
	odkAuth   *context.Context
	odkClient odk.APIClient
	oksAuth   *context.Context
	oksClient oks.APIClient
}

const ( // values not used in .tf files
	// Dictionary #1
	DICT_LANGUAGE_PL = 1
	DICT_LANGUAGE_EN = 2
	// Dictionary #40
	DICT_TICKET_NEW     = 134
	DICT_TICKET_RUNNING = 135
	DICT_TICKET_SUCCEED = 136
	DICT_TICKET_ERROR   = 137
)

const ( // <export> values exported to .tf files

	// Dictionary #15
	DICT_SERVICE_TYPE_HTTP  = 43
	DICT_SERVICE_TYPE_HTTPS = 44
	DICT_SERVICE_TYPE_SMTP  = 45
	DICT_SERVICE_TYPE_PORT  = 155
	DICT_SERVICE_TYPE_MYSQL = 287

	// Dictionary #16
	DICT_SESSION_HANDLING_SOURCE_IP = 46
	DICT_SESSION_HANDLING_NONE      = 47
	DICT_SESSION_HANDLING_COOKIE    = 280

	// Dictionary #17
	DICT_DISK_TIER_1 = 48
	DICT_DISK_TIER_2 = 49
	DICT_DISK_TIER_3 = 50
	DICT_DISK_TIER_4 = 895
	DICT_DISK_TIER_5 = 896

	// Dictionary #27
	DICT_INSTANCE_STATUS_ON           = 86
	DICT_INSTANCE_STATUS_OFF          = 87
	DICT_INSTANCE_STATUS_DELETED      = 126
	DICT_INSTANCE_STATUS_INITIALIZING = 1748

	// Dictionary #36
	DICT_IP_VERSION_IPV4 = 115
	DICT_IP_VERSION_IPV6 = 116
	DICT_IP_VERSION_BOTH = 565

	// Dictionary #55
	DICT_AUTOSCALING_OFF = 184
	DICT_AUTOSCALING_ON  = 185

	// Dictionary #77
	DICT_LB_ALGORITHM_LEAST_CONNECTION    = 281
	DICT_LB_ALGORITHM_LEAST_RESPONSE_TIME = 282
	DICT_LB_ALGORITHM_IP_HASH             = 288
	DICT_LB_ALGORITHM_ROUND_ROBIN         = 612

	// Dictionary #123
	DICT_IP_TYPE_STATIC    = 1106
	DICT_IP_TYPE_AUTOMATIC = 1107
	DICT_IP_TYPE_EXTERNAL  = 1108
	DICT_IP_TYPE_RESERVED  = 1109

	// Dictionary #140
	DICT_PUBLICATION_STATUS_PRIVATE               = 1245
	DICT_PUBLICATION_STATUS_UNDERGOING_ACCEPTANCE = 1246
	DICT_PUBLICATION_STATUS_PUBLIC                = 1247
	DICT_PUBLICATION_STATUS_REJECTED              = 1248

	// Dictionary #159
	DICT_LOGIN_TYPE_SSH_KEYS      = 1398
	DICT_LOGIN_TYPE_USER_AND_PASS = 1399

	// Dictionary #160
	DICT_AFFINITY_TYPE_NO_SEPARATION       = 1403
	DICT_AFFINITY_TYPE_MINIMIZE_SEPARATION = 1404
	DICT_AFFINITY_TYPE_MAXIMIZE_SEPARATION = 1405

	// Dictionary #162
	DICT_SHARE_TYPE_LINUX   = 1411
	DICT_SHARE_TYPE_WINDOWS = 1412

	// Dictionary #167
	DICT_ETHERNET_CONTROLLER_E1000   = 1442
	DICT_ETHERNET_CONTROLLER_VMXNET3 = 1443

	// Dictionary #301
	DICT_IP_MODE_NORMAL   = 1858
	DICT_IP_MODE_FLOATING = 1859
	DICT_IP_MODE_KAS      = 1860

	// Dictionary #302
	DICT_PROXY_PROTOCOL_NONE = 1861
	DICT_PROXY_PROTOCOL_V1   = 1862
	DICT_PROXY_PROTOCOL_V2   = 1863
) // </export>

type DCConfig struct {
	odkApiUrl string
	oksApiUrl string
}

var dcConfigs = map[string]DCConfig{
	"DC1_DEV": {
		odkApiUrl: "https://pl1-api.dev.oktawave.com/services",
		oksApiUrl: "https://k44s-api-devenv.i.k44sdev.oktawave.com",
	},
	"DC2_DEV": {
		odkApiUrl: "https://pl2-api.dev.oktawave.com/services",
		oksApiUrl: "https://k44s-api-devenv.i.k44sdev.oktawave.com",
	},
	"DC1": {
		odkApiUrl: "https://pl1-api.oktawave.com/services",
		oksApiUrl: "https://k44s-api.i.k44s.oktawave.com",
	},
	"DC2": {
		odkApiUrl: "https://pl2-api.oktawave.com/services",
		oksApiUrl: "https://k44s-api.i.k44s.oktawave.com",
	},
}
