/*
Use this resource to create dayu DDoS policy v2

Example Usage

```hcl
resource "tencentcloud_dayu_ddos_policy_v2" "asd" {
  resource_id = "bgpip-000004x0"
    white_ips   = null
    black_ips   = ["1.2.2.2", "1.2.3.5"]
    port_acls {
      action       = "drop"
      d_end_port   = 200
      d_start_port = 101
      priority     = 10
      protocol     = "all"
      s_end_port   = 100
      s_start_port = 1
    }
    port_acls {
      action       = "drop"
      d_end_port   = 300
      d_start_port = 201
      priority     = 10
      protocol     = "all"
      s_end_port   = 100
      s_start_port = 1
    }
    ddos_ai = "on"
    drop_options {
      drop_icmp        = false
      drop_other       = true
      drop_tcp         = false
      drop_udp         = false
      null_conn_enable = true
    }
    geo_ip_blocks {
      action      = "drop"
      area_list   = []
      region_type = "oversea"
    }
    packet_filters {
      action         = "transmit"
      d_end_port     = 12
      d_start_port   = 12
      depth          = 0
      is_include     = false
      match_begin    = "no_match"
      offset         = 0
      pkt_length_max = 1024
      pkt_length_min = 1024
      protocol       = "all"
      s_end_port     = 80
      s_start_port   = 80
    }
  speed_limits {
    dst_port_list = "1-80;8999"
    mode          = 1
    protocol_list = "SMP"
    speed_values {
      type  = 1
      value = 20
    }
    speed_values {
      type  = 2
      value = 10
    }

  }
}
```
*/
package tencentcloud

import (
	"context"
	"log"

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
				Optional:     true,
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
							ValidateFunc: validateAllowedStringValue(DDOS_PROTOCOL),
							Description:  "Protocol. Valid values are `tcp`, `udp`, `all`.",
						},
						"d_start_port": {
							Type:         schema.TypeInt,
							Optional:     true,
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
							ValidateFunc: validatePort,
							Description:  "Source start port. Valid value ranges: (0~65535).",
						},
						"s_end_port": {
							Type:         schema.TypeInt,
							Optional:     true,
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
							ValidateFunc: validateIntegerInRange(0, 1000),
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
				Optional: true,
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
			"connect_limits": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sd_new_limit": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The limit on the number of news per second based on source IP and destination IP.",
						},
						"sd_conn_limit": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Concurrent connection control based on source IP and destination IP.",
						},
						"dst_new_limit": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Limit on the number of news per second based on the destination IP.",
						},
						"dst_conn_limit": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Concurrent connection control based on destination IP and destination port.",
						},
						"bad_conn_threshold": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Abnormal connection detection conditions, empty connection guard switch, value are true and false.",
						},
						"syn_rate": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Anomalous connection detection condition, percentage of syn ack, value range [0,100].",
						},
						"syn_limit": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Anomaly connection detection condition, syn threshold, value range [0,100].",
						},
						"conn_timeout": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Abnormal connection detection condition, connection timeout, value range [0,65535].",
						},
						"null_conn_enable": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicate whether to drop ICMP protocol or not.",
						},
					},
				},
				Description: "Connection suppression configuration details.",
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
				Type:     schema.TypeList,
				Optional: true,
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

	resourceId := d.Get("resource_id").(string)

	// set IpBlackWhite
	var blackIps []interface{}
	var whiteIps []interface{}

	if v, ok := d.GetOk("black_ips"); ok {
		blackIps = v.(*schema.Set).List()
	}
	if v, ok := d.GetOk("white_ips"); ok {
		whiteIps = v.(*schema.Set).List()
	}

	antiddosService := AntiddosService{client: meta.(*TencentCloudClient).apiV3Conn}
	err := antiddosService.CreateIpBlackWhite(resourceId, blackIps, whiteIps)
	if err != nil {
		return err
	}

	// set PortAcls
	var portAclMapping []interface{}
	if v, ok := d.GetOk("port_acls"); ok {
		portAclMapping = v.([]interface{})
	}
	err = antiddosService.CreatePortAcl(ctx, resourceId, portAclMapping)
	if err != nil {
		return err
	}

	//set DropOption
	var dropMapping []interface{}
	if v, ok := d.GetOk("drop_options"); ok {
		dropMapping = v.([]interface{})
	}
	err = antiddosService.CreateDropOption(ctx, resourceId, dropMapping)
	if err != nil {
		return err
	}

	//set DDoSAI protection
	var ddosAi string
	if v, ok := d.GetOk("ddos_ai"); ok {
		ddosAi = v.(string)
	}
	err = antiddosService.CreateAIProtection(ctx, resourceId, ddosAi)
	if err != nil {
		return err
	}

	//set DDoSPolicyPacketFilter
	var packetFilterMapping []interface{}
	if v, ok := d.GetOk("packet_filters"); ok {
		packetFilterMapping = v.([]interface{})
	}
	err = antiddosService.CreatePacketFilter(ctx, resourceId, packetFilterMapping)
	if err != nil {
		return err
	}

	//set DDoSSpeedLimitConfig
	var speedLimitMapping []interface{}
	if v, ok := d.GetOk("speed_limits"); ok {
		speedLimitMapping = v.([]interface{})
	}
	err = antiddosService.CreateDDoSSpeedLimitConfig(ctx, resourceId, speedLimitMapping)
	if err != nil {
		return err
	}

	//set DDoSGeoIPBlockConfig
	var geoIpBlockMapping []interface{}
	if v, ok := d.GetOk("geo_ip_blocks"); ok {
		geoIpBlockMapping = v.([]interface{})
	}
	err = antiddosService.CreateDDoSGeoIPBlockConfig(ctx, resourceId, geoIpBlockMapping)
	if err != nil {
		return err
	}

	d.SetId(resourceId)

	return resourceTencentCloudDayuDdosPolicyReadV2(d, meta)
}

func resourceTencentCloudDayuDdosPolicyReadV2(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dayu_ddos_policy_v2.read")()

	logId := getLogId(contextNil)
	// ctx := context.WithValue(context.TODO(), logIdKey, logId)

	resourceId := "bgpip-000004x0"
	antiddosService := AntiddosService{client: meta.(*TencentCloudClient).apiV3Conn}
	ipList, err := antiddosService.DescribeBlackWhiteIpList(resourceId)
	if err != nil {
		log.Printf("[CRITAL]%s DescribeBlackWhiteIpList error, reason:%s\n", logId, err)
		return nil
	}
	blackIps := make([]string, 0)
	whiteIps := make([]string, 0)
	for _, ip := range ipList {
		if *ip.Type == DDOS_BLACK_TYPE {
			blackIps = append(blackIps, *ip.Ip)
		}
		if *ip.Type == DDOS_WHITE_TYPE {
			whiteIps = append(whiteIps, *ip.Ip)
		}
	}
	_ = d.Set("black_ips", blackIps)
	_ = d.Set("white_ips", whiteIps)
	aclList, err := antiddosService.DescribeListPortAclList(resourceId)
	if err != nil {
		log.Printf("[CRITAL]%s DescribeListPortAclList error, reason:%s\n", logId, err)
		return nil
	}
	portAcls := make([]map[string]interface{}, 0)
	for _, acl := range aclList {
		tmpPortAcl := make(map[string]interface{})
		tmpPortAcl["protocol"] = *acl.AclConfig.ForwardProtocol
		tmpPortAcl["d_start_port"] = *acl.AclConfig.DPortStart
		tmpPortAcl["d_end_port"] = *acl.AclConfig.DPortEnd
		tmpPortAcl["s_start_port"] = *acl.AclConfig.SPortStart
		tmpPortAcl["s_end_port"] = *acl.AclConfig.SPortEnd
		tmpPortAcl["action"] = *acl.AclConfig.Action
		tmpPortAcl["priority"] = *acl.AclConfig.Priority

		portAcls = append(portAcls, tmpPortAcl)
	}
	_ = d.Set("port_acls", portAcls)

	protocolBlockList, err := antiddosService.DescribeListProtocolBlockConfig(resourceId)
	if err != nil {
		log.Printf("[CRITAL]%s DescribeListProtocolBlockConfig error, reason:%s\n", logId, err)
		return nil
	}
	dropOptions := make([]map[string]interface{}, 0)
	for _, protocolBlock := range protocolBlockList {
		tmpDropOption := make(map[string]interface{})

		tmpDropOption["drop_tcp"] = *protocolBlock.ProtocolBlockConfig.DropTcp == 1
		tmpDropOption["drop_udp"] = *protocolBlock.ProtocolBlockConfig.DropUdp == 1
		tmpDropOption["drop_icmp"] = *protocolBlock.ProtocolBlockConfig.DropIcmp == 1
		tmpDropOption["drop_other"] = *protocolBlock.ProtocolBlockConfig.DropOther == 1
		tmpDropOption["null_conn_enable"] = *protocolBlock.ProtocolBlockConfig.CheckExceptNullConnect == 1
		dropOptions = append(dropOptions, tmpDropOption)
	}
	_ = d.Set("drop_options", dropOptions)

	aiConfigList, err := antiddosService.DescribeListDDoSAI(resourceId)
	if err != nil {
		log.Printf("[CRITAL]%s DescribeListDDoSAI error, reason:%s\n", logId, err)
		return nil
	}
	if len(aiConfigList) > 0 {
		ddosAi := *aiConfigList[0].DDoSAI
		_ = d.Set("ddos_ai", ddosAi)
	}

	geoIPBlockConfigList, err := antiddosService.DescribeListDDoSGeoIPBlockConfig(resourceId)
	if err != nil {
		log.Printf("[CRITAL]%s DescribeListDDoSGeoIPBlockConfig error, reason:%s\n", logId, err)
		return nil
	}
	geoIpBlocks := make([]map[string]interface{}, 0)
	for _, geoIPBlockConfig := range geoIPBlockConfigList {
		geoIpBlock := make(map[string]interface{})

		geoIpBlock["region_type"] = *geoIPBlockConfig.GeoIPBlockConfig.RegionType
		geoIpBlock["action"] = *geoIPBlockConfig.GeoIPBlockConfig.Action
		geoIpBlock["area_list"] = geoIPBlockConfig.GeoIPBlockConfig.AreaList
		geoIpBlocks = append(geoIpBlocks, geoIpBlock)
	}
	_ = d.Set("geo_ip_blocks", geoIpBlocks)

	speedLimitConfigList, err := antiddosService.DescribeListDDoSSpeedLimitConfig(resourceId)
	if err != nil {
		log.Printf("[CRITAL]%s speedLimitConfigList error, reason:%s\n", logId, err)
		return nil
	}
	speedLimitConfigs := make([]map[string]interface{}, 0)
	for _, speedLimitConfig := range speedLimitConfigList {
		tmpSpeedLimitConfig := make(map[string]interface{})
		speedValues := make([]map[string]interface{}, 0)
		for _, speedValue := range speedLimitConfig.SpeedLimitConfig.SpeedValues {
			tmpSpeedValue := make(map[string]interface{})
			tmpSpeedValue["type"] = speedValue.Type
			tmpSpeedValue["value"] = speedValue.Value
			speedValues = append(speedValues, tmpSpeedValue)
		}
		tmpSpeedLimitConfig["speed_values"] = speedValues
		tmpSpeedLimitConfig["mode"] = speedLimitConfig.SpeedLimitConfig.Mode
		tmpSpeedLimitConfig["protocol_list"] = speedLimitConfig.SpeedLimitConfig.ProtocolList
		tmpSpeedLimitConfig["dst_port_list"] = speedLimitConfig.SpeedLimitConfig.DstPortList
		speedLimitConfigs = append(speedLimitConfigs, tmpSpeedLimitConfig)
	}
	_ = d.Set("speed_limits", speedLimitConfigs)

	packetFilterConfigList, err := antiddosService.DescribeListPacketFilterConfig(resourceId)
	if err != nil {
		log.Printf("[CRITAL]%s DescribeListPacketFilterConfig error, reason:%s\n", logId, err)
		return nil
	}
	packetFilterConfigs := make([]map[string]interface{}, 0)
	for _, packetFilterConfig := range packetFilterConfigList {
		tmpPacketFilterConfig := make(map[string]interface{})

		tmpPacketFilterConfig["protocol"] = *packetFilterConfig.PacketFilterConfig.Protocol
		tmpPacketFilterConfig["d_start_port"] = *packetFilterConfig.PacketFilterConfig.DportStart
		tmpPacketFilterConfig["d_end_port"] = *packetFilterConfig.PacketFilterConfig.DportEnd
		tmpPacketFilterConfig["s_start_port"] = *packetFilterConfig.PacketFilterConfig.SportStart
		tmpPacketFilterConfig["s_end_port"] = *packetFilterConfig.PacketFilterConfig.SportEnd
		tmpPacketFilterConfig["pkt_length_min"] = *packetFilterConfig.PacketFilterConfig.PktlenMin
		tmpPacketFilterConfig["pkt_length_max"] = *packetFilterConfig.PacketFilterConfig.PktlenMax
		tmpPacketFilterConfig["match_begin"] = *packetFilterConfig.PacketFilterConfig.MatchBegin
		tmpPacketFilterConfig["match_type"] = *packetFilterConfig.PacketFilterConfig.MatchType
		tmpPacketFilterConfig["match_str"] = *packetFilterConfig.PacketFilterConfig.Str
		tmpPacketFilterConfig["depth"] = *packetFilterConfig.PacketFilterConfig.Depth
		tmpPacketFilterConfig["offset"] = *packetFilterConfig.PacketFilterConfig.Offset
		tmpPacketFilterConfig["is_include"] = *packetFilterConfig.PacketFilterConfig.IsNot == 1
		tmpPacketFilterConfig["action"] = *packetFilterConfig.PacketFilterConfig.Action
		packetFilterConfigs = append(packetFilterConfigs, tmpPacketFilterConfig)
	}
	_ = d.Set("packet_filters", packetFilterConfigs)

	connectLimitList, err := antiddosService.DescribeDDoSConnectLimitList(resourceId)
	if err != nil {
		log.Printf("[CRITAL]%s DescribeDDoSConnectLimitList error, reason:%s\n", logId, err)
		return nil
	}
	connectLimits := make([]map[string]interface{}, 0)
	for _, connectLimit := range connectLimitList {
		tmpconnectLimit := make(map[string]interface{})

		tmpconnectLimit["sd_new_limit"] = *connectLimit.ConnectLimitConfig.SdNewLimit
		tmpconnectLimit["sd_conn_limit"] = *connectLimit.ConnectLimitConfig.SdConnLimit
		tmpconnectLimit["dst_new_limit"] = *connectLimit.ConnectLimitConfig.DstNewLimit
		tmpconnectLimit["dst_conn_limit"] = *connectLimit.ConnectLimitConfig.DstConnLimit
		tmpconnectLimit["bad_conn_threshold"] = *connectLimit.ConnectLimitConfig.BadConnThreshold
		tmpconnectLimit["syn_rate"] = *connectLimit.ConnectLimitConfig.SynRate
		tmpconnectLimit["syn_limit"] = *connectLimit.ConnectLimitConfig.SynLimit
		tmpconnectLimit["conn_timeout"] = *connectLimit.ConnectLimitConfig.ConnTimeout
		tmpconnectLimit["null_conn_enable"] = *connectLimit.ConnectLimitConfig.NullConnEnable == 1

		connectLimits = append(connectLimits, tmpconnectLimit)
	}
	_ = d.Set("connect_limits", connectLimits)

	return nil
}

func resourceTencentCloudDayuDdosPolicyUpdateV2(d *schema.ResourceData, meta interface{}) error {
	return resourceTencentCloudDayuDdosPolicyRead(d, meta)
}

func resourceTencentCloudDayuDdosPolicyDeleteV2(d *schema.ResourceData, meta interface{}) error {
	return nil
}
