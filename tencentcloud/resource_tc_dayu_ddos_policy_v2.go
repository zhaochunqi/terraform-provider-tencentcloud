/*
Use this resource to create dayu DDoS policy

Example Usage

```hcl
resource "tencentcloud_dayu_ddos_policy" "test_policy" {
  resource_type = "bgpip"
  name          = "tf_test_policy"
  black_ips     = ["1.1.1.1"]
  white_ips     = ["2.2.2.2"]

  drop_options {
    drop_tcp           = true
    drop_udp           = true
    drop_icmp          = true
    drop_other         = true
    drop_abroad        = true
    check_sync_conn    = true
    s_new_limit        = 100
    d_new_limit        = 100
    s_conn_limit       = 100
    d_conn_limit       = 100
    tcp_mbps_limit     = 100
    udp_mbps_limit     = 100
    icmp_mbps_limit    = 100
    other_mbps_limit   = 100
    bad_conn_threshold = 100
    null_conn_enable   = true
    conn_timeout       = 500
    syn_rate           = 50
    syn_limit          = 100
  }

  port_limits {
    start_port = "2000"
    end_port   = "2500"
    protocol   = "all"
    action     = "drop"
    kind       = 1
  }

  packet_filters {
    protocol       = "tcp"
    action         = "drop"
    d_start_port   = 1000
    d_end_port     = 1500
    s_start_port   = 2000
    s_end_port     = 2500
    pkt_length_max = 1400
    pkt_length_min = 1000
    is_include     = true
    match_begin    = "begin_l5"
    match_type     = "pcre"
    depth          = 1000
    offset         = 500
  }

  watermark_filters {
    tcp_port_list = ["2000-3000", "3500-4000"]
    udp_port_list = ["5000-6000"]
    offset        = 50
    auto_remove   = true
    open_switch   = true
  }
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceTencentCloudDayuDdosPolicyV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDayuDdosPolicyCreateV2,
		Read:   resourceTencentCloudDayuDdosPolicyReadV2,
		Update: resourceTencentCloudDayuDdosPolicyUpdateV2,
		Delete: resourceTencentCloudDayuDdosPolicyDeleteV2,

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Resource id of the DDoS.",
			},
			"ddos_ai": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue(DDOS_AI_SWITCH),
				Description:  "AI protection switch,  Valid values are `on`, `off`.",
			},
			"port_acls": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateAllowedStringValue(DAYU_PROTOCOL),
							Description:  "Protocol. Valid values are `tcp`, `udp`, `all`.",
						},
						"d_start_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validatePort,
							Description:  "Destination start port. Valid value ranges: (0~65535).",
						},
						"d_end_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      65535,
							ValidateFunc: validatePort,
							Description:  "Destination end port. Valid value ranges: (0~65535). It must be greater than `d_start_port`.",
						},
						"s_start_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validatePort,
							Description:  "Source start port. Valid value ranges: (0~65535).",
						},
						"s_end_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      65535,
							ValidateFunc: validatePort,
							Description:  "Source end port. Valid value ranges: (0~65535). It must be greater than `s_start_port`.",
						},
						"action": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateAllowedStringValue(DDOS_PORT_ACTION),
							Description:  "Action of port to take. Valid values: `drop`, `transmit`, `forward`.",
						},
						"priority": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateStringLengthInRange(0, 1000),
							Description:  "The priority of port acl rule. Valid value ranges: (0~1000). The smaller the number, the higher the priority, and the default priority is 10",
						},
					},
				},
				Description: "Port acls.",
			},
			"black_ips": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateIp,
				},
				Optional:    true,
				Description: "Black IP list.",
			},
			"white_ips": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateIp,
				},
				Optional:    true,
				Description: "White IP list.",
			},

			"packet_filters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateAllowedStringValue(DDOS_PROTOCOL),
							Description:  "Protocol. Valid values: `tcp`, `udp`, `icmp`, `all`.",
						},
						"d_start_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validatePort,
							Description:  "Start port of the destination. Valid value ranges: (0~65535).",
						},
						"d_end_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validatePort,
							Description:  "End port of the destination. Valid value ranges: (0~65535). It must be greater than `d_start_port`.",
						},
						"s_start_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validatePort,
							Description:  "Start port of the source. Valid value ranges: (0~65535).",
						},
						"s_end_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validatePort,
							Description:  "End port of the source. Valid value ranges: (0~65535). It must be greater than `s_start_port`.",
						},
						"pkt_length_min": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateIntegerInRange(0, 1500),
							Description:  "The minimum length of the packet. Valid value ranges: (0~1500)(Mbps).",
						},
						"pkt_length_max": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateIntegerInRange(0, 1500),
							Description:  "The max length of the packet. Valid value ranges: (0~1500)(Mbps). It must be greater than `pkt_length_min`.",
						},
						"match_begin": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateAllowedStringValue(DDOS_MATCH_SWITCH),
							Description:  "Indicate whether to check load or not.",
						},
						"match_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateAllowedStringValue(DDOS_MATCH_TYPE),
							Description:  "Match type. Valid values: `sunday` and `pcre`. `sunday` means key word match while `pcre` means regular match.",
						},
						"match_str": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The key word or regular expression.",
						},
						"depth": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateIntegerInRange(0, 1500),
							Description:  "The depth of match. Valid value ranges: (0~1500).",
						},
						"offset": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateIntegerInRange(0, 1500),
							Description:  "The offset of match. Valid value ranges: (0~1500).",
						},
						"is_include": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Indicate whether to include the key word/regular expression or not.",
						},
						"action": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateAllowedStringValue(DDOS_PACKET_ACTION),
							Description:  "Action of port to take. Valid values: `drop`, `drop_black`,`drop_rst`,`drop_black_rst`,`transmit`.`drop`(drop the packet), `drop_black`(drop the packet and black the ip),`drop_rst`(drop the packet and disconnect),`drop_black_rst`(drop the packet, black the ip and disconnect),`transmit`(transmit the packet).",
						},
					},
				},
				Description: "Message filter options list.",
			},
			"drop_options": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"drop_tcp": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicate whether to drop TCP protocol or not.",
						},
						"drop_udp": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicate to drop UDP protocol or not.",
						},
						"drop_icmp": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicate whether to drop ICMP protocol or not.",
						},
						"drop_other": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicate whether to drop other protocols(exclude TCP/UDP/ICMP) or not.",
						},
						"null_conn_enable": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicate to enable null connection or not.",
						},
					},
				},
				Description: "Option list of abnormal check of the DDos policy, should set at least one policy.",
			},
			"speed_limits": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"speed_values": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validateIntegerInRange(1, 2),
										Description:  "Type of rate. Valid value: `1`(pps)、`2`(bps).",
									},
									"value": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Value of rate.",
									},
								},
							},
							Description: "Speed limit values, up to one per type of speed limit value, this field array has at least one speed limit value.",
						},
						"mode": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validateIntegerInRange(1, 2),
							Description:  "The mode of speed limit. Valid value ranges: (1-2).",
						},
						"protocol_list": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Protocol list of speed limit, Valid value: `ALL`、`TCP`、`UDP`、`SMP` and custom protocol number. When customizing the protocol number range, you can only fill in the protocol number, multiple ranges; Separation; no other agreements or protocol numbers can be filled in when all is filled out.",
						},
						"dst_port_list": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "List of port ranges, up to 8, multiple; separated, the range is represented with -.",
						},
					},
				},
				Description: "Watermark policy options, and only support one watermark policy at most.",
			},
			"geo_ip_blocks": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The zone type.",
						},
						"action": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Block the action.",
						},
						"area_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type:        schema.TypeInt,
								Required:    true,
								Description: "Area number list.",
							},
							Description: "Area LIst. When the RegionType is customized, the AreaList must be filled in, and a maximum of 128 must be filled in.",
						},
					},
				},
				Description: "Option list of abnormal check of the DDos policy, should set at least one policy.",
			},
		},
	}
}

func resourceTencentCloudDayuDdosPolicyCreateV2(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dayu_ddos_policy.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	// resourceType := d.Get("resource_type").(string)
	resourceId := d.Get("resource_id").(string)

	// set IpBlackWhite
	blackIps := d.Get("black_ips").(*schema.Set).List()
	whiteIps := d.Get("white_ips").(*schema.Set).List()
	antiddosService := AntiddosService{client: meta.(*TencentCloudClient).apiV3Conn}
	err := antiddosService.CreateIpBlackWhite(resourceId, blackIps, whiteIps)
	if err != nil {
		return err
	}

	// set PortAcls
	portAclMapping := d.Get("port_filters").([]interface{})
	err = antiddosService.CreatePortAcl(ctx, resourceId, portAclMapping)
	if err != nil {
		return err
	}

	//set DropOption
	dropMapping := d.Get("drop_options").([]interface{})
	err = antiddosService.CreateDropOption(ctx, resourceId, dropMapping)
	if err != nil {
		return err
	}

	//set DDoSPolicyPacketFilter
	packetFilterMapping := d.Get("packet_filters").([]interface{})
	err = antiddosService.CreatePacketFilter(ctx, resourceId, packetFilterMapping)
	if err != nil {
		return err
	}

	//set DDoSSpeedLimitConfig
	speedLimitMapping := d.Get("speed_limits").([]interface{})
	err = antiddosService.CreateDDoSSpeedLimitConfig(ctx, resourceId, speedLimitMapping)
	if err != nil {
		return err
	}

	//set DDoSAI protection
	ddosAi := d.Get("ddos_ai").(string)
	err = antiddosService.CreateAIProtection(ctx, resourceId, ddosAi)
	if err != nil {
		return err
	}

	//set DDoSGeoIPBlockConfig
	geoIpBlockMapping := d.Get("geo_ip_blocks").(*schema.Set).List()
	err = antiddosService.CreateDDoSGeoIPBlockConfig(ctx, resourceId, geoIpBlockMapping)
	if err != nil {
		return err
	}

	return resourceTencentCloudDayuDdosPolicyReadV2(d, meta)
}

func resourceTencentCloudDayuDdosPolicyReadV2(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dayu_ddos_policy_v2.read")()

	// logId := getLogId(contextNil)
	// ctx := context.WithValue(context.TODO(), logIdKey, logId)

	// resourceId := d.Get("resource_id").(string)
	// antiddosService := AntiddosService{client: meta.(*TencentCloudClient).apiV3Conn}

	return nil
}

func resourceTencentCloudDayuDdosPolicyUpdateV2(d *schema.ResourceData, meta interface{}) error {
	return resourceTencentCloudDayuDdosPolicyRead(d, meta)
}

func resourceTencentCloudDayuDdosPolicyDeleteV2(d *schema.ResourceData, meta interface{}) error {
	return nil
}
