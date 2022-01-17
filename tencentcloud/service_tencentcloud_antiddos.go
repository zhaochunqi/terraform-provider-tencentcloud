package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	antiddos "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20200309"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type AntiddosService struct {
	client *connectivity.TencentCloudClient
}

func (me *AntiddosService) DescribeBlackWhiteIpList(instanceId string) (result []*antiddos.BlackWhiteIpRelation, err error) {
	request := antiddos.NewDescribeListBlackWhiteIpListRequest()
	offset := int64(0)
	request.Offset = &offset
	result = make([]*antiddos.BlackWhiteIpRelation, 0)
	limit := int64(DDOS_DESCRIBE_LIMIT)
	request.Limit = &limit
	var response *antiddos.DescribeListBlackWhiteIpListResponse
	for {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			response, err = me.client.UseAntiddosClient().DescribeListBlackWhiteIpList(request)
			if e, ok := err.(*errors.TencentCloudSDKError); ok {
				if e.GetCode() == "InternalError.ClusterNotFound" {
					return nil
				}
			}
			if err != nil {
				return resource.RetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL] read ddos blackwhile list failed, reason:%s\n", err.Error())
			return
		} else {
			result = append(result, response.Response.IpList...)
			if len(response.Response.IpList) < DDOS_DESCRIBE_LIMIT {
				break
			} else {
				offset = offset + limit
			}
		}
	}
	return
}

func (me *AntiddosService) CreateIpBlackWhite(instanceId string, blackIps []interface{}, whiteIps []interface{}) (err error) {

	blackIpsWithMask := make([]*antiddos.IpSegment, 0)
	whiteIpsWithMask := make([]*antiddos.IpSegment, 0)
	for _, blackIp := range blackIps {
		blackIpsWithMask = append(blackIpsWithMask, &antiddos.IpSegment{Ip: helper.String(blackIp.(string)), Mask: helper.IntUint64(0)})
	}
	for _, whiteIp := range whiteIps {
		whiteIpsWithMask = append(whiteIpsWithMask, &antiddos.IpSegment{Ip: helper.String(whiteIp.(string)), Mask: helper.IntUint64(0)})
	}
	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		request := antiddos.NewCreateDDoSBlackWhiteIpListRequest()
		ratelimit.Check(request.GetAction())
		request.InstanceId = helper.String(instanceId)
		request.IpList = blackIpsWithMask
		request.Type = helper.String(DDOS_BLACK_TYPE)
		_, err = me.client.UseAntiddosClient().CreateDDoSBlackWhiteIpList(request)
		if e, ok := err.(*errors.TencentCloudSDKError); ok {
			if e.GetCode() == "InternalError.ClusterNotFound" {
				return nil
			}
		}
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})

	if err != nil {
		return
	}
	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		request := antiddos.NewCreateDDoSBlackWhiteIpListRequest()
		ratelimit.Check(request.GetAction())
		request.InstanceId = helper.String(instanceId)
		request.IpList = whiteIpsWithMask
		request.Type = helper.String(DDOS_WHITE_TYPE)
		_, err = me.client.UseAntiddosClient().CreateDDoSBlackWhiteIpList(request)
		if e, ok := err.(*errors.TencentCloudSDKError); ok {
			if e.GetCode() == "InternalError.ClusterNotFound" {
				return nil
			}
		}
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		return
	}
	return
}

func (me *AntiddosService) CreatePortAcl(ctx context.Context, instanceId string, mapping []interface{}) (err error) {
	logId := getLogId(ctx)
	for _, vv := range mapping {
		v := vv.(map[string]interface{})
		dStartPort := v["d_start_port"].(uint64)
		dEndPort := v["d_end_port"].(uint64)
		sStartPort := v["s_start_port"].(uint64)
		sEndPort := v["s_end_port"].(uint64)
		protocol := v["kind"].(string)
		action := v["action"].(string)
		priority := v["priority"].(uint64)
		if dStartPort > dEndPort {
			err = fmt.Errorf("The `dStartPort` should not be greater than `dEndPort`.")
			return
		}
		if sStartPort > sEndPort {
			err = fmt.Errorf("The `sStartPort` should not be greater than `sEndPort`.")
			return
		}
		var request *antiddos.CreatePortAclConfigRequest
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			request = antiddos.NewCreatePortAclConfigRequest()
			ratelimit.Check(request.GetAction())
			request.InstanceId = common.StringPtr(instanceId)
			request.AclConfig = &antiddos.AclConfig{
				ForwardProtocol: common.StringPtr(protocol),
				DPortStart:      common.Uint64Ptr(dStartPort),
				DPortEnd:        common.Uint64Ptr(dEndPort),
				SPortStart:      common.Uint64Ptr(sStartPort),
				SPortEnd:        common.Uint64Ptr(sEndPort),
				Action:          common.StringPtr(action),
				Priority:        common.Uint64Ptr(priority),
			}
			_, err = me.client.UseAntiddosClient().CreatePortAclConfig(request)
			if e, ok := err.(*errors.TencentCloudSDKError); ok {
				if e.GetCode() == "InternalError.ClusterNotFound" {
					return nil
				}
			}
			if err != nil {
				return resource.RetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return
		}
	}
	return
}

func (me *AntiddosService) CreateDropOption(ctx context.Context, instanceId string, mapping []interface{}) (err error) {
	logId := getLogId(ctx)
	for _, vv := range mapping {
		v := vv.(map[string]interface{})
		dropTcp := v["drop_tcp"].(bool)
		dropUdp := v["drop_udp"].(bool)
		dropIcmp := v["drop_icmp"].(bool)
		dropOther := v["drop_other"].(bool)
		nullConnEnable := v["null_conn_enable"].(bool)
		var request *antiddos.CreateProtocolBlockConfigRequest
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			request = antiddos.NewCreateProtocolBlockConfigRequest()
			ratelimit.Check(request.GetAction())
			request.InstanceId = common.StringPtr(instanceId)
			request.ProtocolBlockConfig = &antiddos.ProtocolBlockConfig{
				DropTcp:                common.Int64Ptr(bool2int64(dropTcp)),
				DropUdp:                common.Int64Ptr(bool2int64(dropUdp)),
				DropIcmp:               common.Int64Ptr(bool2int64(dropIcmp)),
				DropOther:              common.Int64Ptr(bool2int64(dropOther)),
				CheckExceptNullConnect: common.Int64Ptr(bool2int64(nullConnEnable)),
			}
			_, err = me.client.UseAntiddosClient().CreateProtocolBlockConfig(request)
			if e, ok := err.(*errors.TencentCloudSDKError); ok {
				if e.GetCode() == "InternalError.ClusterNotFound" {
					return nil
				}
			}
			if err != nil {
				return resource.RetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return
		}
	}
	return
}

func (me *AntiddosService) CreatePacketFilter(ctx context.Context, instanceId string, mapping []interface{}) (err error) {
	logId := getLogId(ctx)
	for _, vv := range mapping {
		v := vv.(map[string]interface{})
		protocol := v["protocol"].(string)
		dStartPort := v["d_start_port"].(int64)
		dEndPort := v["d_end_port"].(int64)
		sStartPort := v["s_start_port"].(int64)
		sEndPort := v["s_end_port"].(int64)
		pktLengthMin := v["pkt_length_min"].(int64)
		pktLengthMax := v["pkt_length_max"].(int64)
		matchBegin := v["match_begin"].(string)
		matchType := v["match_type"].(string)
		matchStr := v["match_str"].(string)
		depth := v["depth"].(int64)
		offset := v["offset"].(int64)
		isInclude := v["is_include"].(bool)
		action := v["action"].(string)

		var request *antiddos.CreatePacketFilterConfigRequest
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			request := antiddos.NewCreatePacketFilterConfigRequest()
			ratelimit.Check(request.GetAction())
			request.InstanceId = common.StringPtr(instanceId)
			request.PacketFilterConfig = &antiddos.PacketFilterConfig{
				Protocol:   common.StringPtr(protocol),
				SportStart: common.Int64Ptr(sStartPort),
				SportEnd:   common.Int64Ptr(sEndPort),
				DportStart: common.Int64Ptr(dStartPort),
				DportEnd:   common.Int64Ptr(dEndPort),
				PktlenMin:  common.Int64Ptr(pktLengthMin),
				PktlenMax:  common.Int64Ptr(pktLengthMax),
				Action:     common.StringPtr(action),
				MatchBegin: common.StringPtr(matchBegin),
				MatchType:  common.StringPtr(matchType),
				Str:        common.StringPtr(matchStr),
				Depth:      common.Int64Ptr(depth),
				Offset:     common.Int64Ptr(offset),
				IsNot:      common.Int64Ptr(bool2int64(isInclude)),
			}

			_, err = me.client.UseAntiddosClient().CreatePacketFilterConfig(request)
			if e, ok := err.(*errors.TencentCloudSDKError); ok {
				if e.GetCode() == "InternalError.ClusterNotFound" {
					return nil
				}
			}
			if err != nil {
				return resource.RetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return
		}
	}
	return
}

func (me *AntiddosService) DescribeListPacketFilterConfig(instanceId string) (result []*antiddos.PacketFilterRelation, err error) {
	request := antiddos.NewDescribeListPacketFilterConfigRequest()
	offset := int64(0)
	request.Offset = &offset
	result = make([]*antiddos.PacketFilterRelation, 0)
	limit := int64(DDOS_DESCRIBE_LIMIT)
	request.Limit = &limit
	var response *antiddos.DescribeListPacketFilterConfigResponse
	for {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			response, err = me.client.UseAntiddosClient().DescribeListPacketFilterConfig(request)
			if e, ok := err.(*errors.TencentCloudSDKError); ok {
				if e.GetCode() == "InternalError.ClusterNotFound" {
					return nil
				}
			}
			if err != nil {
				return resource.RetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL] read ddos blackwhile list failed, reason:%s\n", err.Error())
			return
		} else {
			result = append(result, response.Response.ConfigList...)
			if len(response.Response.ConfigList) < DDOS_DESCRIBE_LIMIT {
				break
			} else {
				offset = offset + limit
			}
		}
	}
	return
}

func (me *AntiddosService) CreateAIProtection(ctx context.Context, instanceId string, ddosAiSwitch string) (err error) {
	logId := getLogId(ctx)

	var request *antiddos.CreateDDoSAIRequest
	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		request = antiddos.NewCreateDDoSAIRequest()
		ratelimit.Check(request.GetAction())
		request.InstanceIdList = common.StringPtrs([]string{instanceId})
		request.DDoSAI = common.StringPtr(ddosAiSwitch)

		_, err = me.client.UseAntiddosClient().CreateDDoSAI(request)
		if e, ok := err.(*errors.TencentCloudSDKError); ok {
			if e.GetCode() == "InternalError.ClusterNotFound" {
				return nil
			}
		}
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *AntiddosService) DescribeListDDoSAI(instanceId string) (result []*antiddos.DDoSAIRelation, err error) {
	request := antiddos.NewDescribeListDDoSAIRequest()
	offset := int64(0)
	request.Offset = &offset
	result = make([]*antiddos.DDoSAIRelation, 0)
	limit := int64(DDOS_DESCRIBE_LIMIT)
	request.Limit = &limit
	var response *antiddos.DescribeListDDoSAIResponse
	for {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			response, err = me.client.UseAntiddosClient().DescribeListDDoSAI(request)
			if e, ok := err.(*errors.TencentCloudSDKError); ok {
				if e.GetCode() == "InternalError.ClusterNotFound" {
					return nil
				}
			}
			if err != nil {
				return resource.RetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL] read ddos blackwhile list failed, reason:%s\n", err.Error())
			return
		} else {
			result = append(result, response.Response.ConfigList...)
			if len(response.Response.ConfigList) < DDOS_DESCRIBE_LIMIT {
				break
			} else {
				offset = offset + limit
			}
		}
	}
	return
}

func (me *AntiddosService) CreateDDoSSpeedLimitConfig(ctx context.Context, instanceId string, mapping []interface{}) (err error) {
	logId := getLogId(ctx)
	for _, vv := range mapping {
		v := vv.(map[string]interface{})
		speedValues := v["speed_values"].([]interface{})
		speedValueList := make([]*antiddos.SpeedValue, 0)
		for _, speedValue := range speedValues {
			speedValueMap := speedValue.(map[string]interface{})
			speedValueType := speedValueMap["type"].(uint64)
			speedValueValue := speedValueMap["value"].(uint64)
			speedValueList = append(speedValueList, &antiddos.SpeedValue{Type: &speedValueType, Value: &speedValueValue})
		}
		mode := v["mode"].(uint64)
		protocolList := v["protocol_list"].(string)
		dstPortList := v["dst_port_list"].(string)

		var request *antiddos.CreateDDoSSpeedLimitConfigRequest
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			request = antiddos.NewCreateDDoSSpeedLimitConfigRequest()
			ratelimit.Check(request.GetAction())
			request.InstanceId = common.StringPtr(instanceId)
			request.DDoSSpeedLimitConfig = &antiddos.DDoSSpeedLimitConfig{
				Mode:         common.Uint64Ptr(mode),
				ProtocolList: common.StringPtr(protocolList),
				DstPortList:  common.StringPtr(dstPortList),
				SpeedValues:  speedValueList,
			}

			_, err = me.client.UseAntiddosClient().CreateDDoSSpeedLimitConfig(request)
			if e, ok := err.(*errors.TencentCloudSDKError); ok {
				if e.GetCode() == "InternalError.ClusterNotFound" {
					return nil
				}
			}
			if err != nil {
				return resource.RetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return
		}
	}
	return
}

func (me *AntiddosService) DescribeListDDoSSpeedLimitConfig(instanceId string) (result []*antiddos.DDoSSpeedLimitConfigRelation, err error) {
	request := antiddos.NewDescribeListDDoSSpeedLimitConfigRequest()
	offset := uint64(0)
	request.Offset = &offset
	result = make([]*antiddos.DDoSSpeedLimitConfigRelation, 0)
	limit := uint64(DDOS_DESCRIBE_LIMIT)
	request.Limit = &limit
	var response *antiddos.DescribeListDDoSSpeedLimitConfigResponse
	for {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			response, err = me.client.UseAntiddosClient().DescribeListDDoSSpeedLimitConfig(request)
			if e, ok := err.(*errors.TencentCloudSDKError); ok {
				if e.GetCode() == "InternalError.ClusterNotFound" {
					return nil
				}
			}
			if err != nil {
				return resource.RetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL] read ddos blackwhile list failed, reason:%s\n", err.Error())
			return
		} else {
			result = append(result, response.Response.ConfigList...)
			if len(response.Response.ConfigList) < DDOS_DESCRIBE_LIMIT {
				break
			} else {
				offset = offset + limit
			}
		}
	}
	return
}

func (me *AntiddosService) CreateDDoSGeoIPBlockConfig(ctx context.Context, instanceId string, mapping []interface{}) (err error) {
	logId := getLogId(ctx)
	for _, vv := range mapping {
		v := vv.(map[string]interface{})
		regionType := v["RegionType"].(string)
		action := v["Action"].(string)
		areaList := v["AreaList"].([]int64)
		var request *antiddos.CreateDDoSGeoIPBlockConfigRequest
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			request := antiddos.NewCreateDDoSGeoIPBlockConfigRequest()
			ratelimit.Check(request.GetAction())
			request.InstanceId = common.StringPtr(instanceId)
			request.DDoSGeoIPBlockConfig = &antiddos.DDoSGeoIPBlockConfig{
				RegionType: common.StringPtr(regionType),
				Action:     common.StringPtr(action),
			}
			if regionType == "customized" {
				if len(areaList) == 0 {
					err := fmt.Errorf("When regionType is `customized`, must set area_list.")
					return retryError(err)
				}
				request.DDoSGeoIPBlockConfig.AreaList = common.Int64Ptrs(areaList)
			}

			_, err = me.client.UseAntiddosClient().CreateDDoSGeoIPBlockConfig(request)
			if e, ok := err.(*errors.TencentCloudSDKError); ok {
				if e.GetCode() == "InternalError.ClusterNotFound" {
					return nil
				}
			}
			if err != nil {
				return resource.RetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return
		}
	}
	return
}

func (me *AntiddosService) DescribeListDDoSGeoIPBlockConfig(instanceId string) (result []*antiddos.DDoSGeoIPBlockConfigRelation, err error) {
	request := antiddos.NewDescribeListDDoSGeoIPBlockConfigRequest()
	offset := uint64(0)
	request.Offset = &offset
	result = make([]*antiddos.DDoSGeoIPBlockConfigRelation, 0)
	limit := uint64(DDOS_DESCRIBE_LIMIT)
	request.Limit = &limit
	var response *antiddos.DescribeListDDoSGeoIPBlockConfigResponse
	for {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			response, err = me.client.UseAntiddosClient().DescribeListDDoSGeoIPBlockConfig(request)
			if e, ok := err.(*errors.TencentCloudSDKError); ok {
				if e.GetCode() == "InternalError.ClusterNotFound" {
					return nil
				}
			}
			if err != nil {
				return resource.RetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL] read ddos blackwhile list failed, reason:%s\n", err.Error())
			return
		} else {
			result = append(result, response.Response.ConfigList...)
			if len(response.Response.ConfigList) < DDOS_DESCRIBE_LIMIT {
				break
			} else {
				offset = offset + limit
			}
		}
	}
	return
}

func (me *AntiddosService) DescribeListProtocolBlockConfig(instanceId string) (result []*antiddos.ProtocolBlockRelation, err error) {
	request := antiddos.NewDescribeListProtocolBlockConfigRequest()
	offset := int64(0)
	request.Offset = &offset
	result = make([]*antiddos.ProtocolBlockRelation, 0)
	limit := int64(DDOS_DESCRIBE_LIMIT)
	request.Limit = &limit
	var response *antiddos.DescribeListProtocolBlockConfigResponse
	for {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			response, err = me.client.UseAntiddosClient().DescribeListProtocolBlockConfig(request)
			if e, ok := err.(*errors.TencentCloudSDKError); ok {
				if e.GetCode() == "InternalError.ClusterNotFound" {
					return nil
				}
			}
			if err != nil {
				return resource.RetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL] read ddos blackwhile list failed, reason:%s\n", err.Error())
			return
		} else {
			result = append(result, response.Response.ConfigList...)
			if len(response.Response.ConfigList) < DDOS_DESCRIBE_LIMIT {
				break
			} else {
				offset = offset + limit
			}
		}
	}
	return
}
