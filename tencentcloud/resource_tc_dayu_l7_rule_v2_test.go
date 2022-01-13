package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testDayuL7RuleV2ResourceName = "tencentcloud_dayu_l7_rule_v2"
var testDayuL7RuleV2ResourceKey = testDayuL7RuleV2ResourceName + ".test_rule"

func TestAccTencentCloudDayuL7RuleV2Resource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDayuL7RuleV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDayuL7RuleV2, defaultDayuBgpIdV2, defaultDayuBgpIpV2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDayuL7RuleV2Exists(testDayuL7RuleV2ResourceKey),
					resource.TestCheckResourceAttrSet(testDayuL7RuleV2ResourceKey, "rule_id"),
					resource.TestCheckResourceAttrSet(testDayuL7RuleV2ResourceKey, "status"),
					resource.TestCheckResourceAttr(testDayuL7RuleV2ResourceKey, "resource_type", "bgpip"),
					resource.TestCheckResourceAttr(testDayuL7RuleV2ResourceKey, "name", "rule_test"),
					resource.TestCheckResourceAttr(testDayuL7RuleV2ResourceKey, "source_type", "2"),
					resource.TestCheckResourceAttr(testDayuL7RuleV2ResourceKey, "source_list.#", "2"),
					resource.TestCheckResourceAttr(testDayuL7RuleV2ResourceKey, "switch", "true"),
					resource.TestCheckResourceAttr(testDayuL7RuleV2ResourceKey, "protocol", "http"),
					resource.TestCheckResourceAttr(testDayuL7RuleV2ResourceKey, "health_check_code", "31"),
					resource.TestCheckResourceAttr(testDayuL7RuleV2ResourceKey, "health_check_switch", "true"),
					resource.TestCheckResourceAttr(testDayuL7RuleV2ResourceKey, "health_check_interval", "30"),
					resource.TestCheckResourceAttr(testDayuL7RuleV2ResourceKey, "health_check_path", "/"),
					resource.TestCheckResourceAttr(testDayuL7RuleV2ResourceKey, "health_check_method", "GET"),
					resource.TestCheckResourceAttr(testDayuL7RuleV2ResourceKey, "health_check_health_num", "5"),
					resource.TestCheckResourceAttr(testDayuL7RuleV2ResourceKey, "health_check_unhealth_num", "10"),
				),
			},
			{
				Config: fmt.Sprintf(testAccDayuL7RuleV2Update, defaultDayuBgpIdV2, defaultDayuBgpIpV2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDayuL7RuleV2Exists(testDayuL7RuleV2ResourceKey),
					testAccCheckDayuL7RuleV2Exists(testDayuL7RuleV2ResourceKey),
					resource.TestCheckResourceAttrSet(testDayuL7RuleV2ResourceKey, "rule_id"),
					resource.TestCheckResourceAttrSet(testDayuL7RuleV2ResourceKey, "status"),
					resource.TestCheckResourceAttr(testDayuL7RuleV2ResourceKey, "resource_type", "bgpip"),
					resource.TestCheckResourceAttr(testDayuL7RuleV2ResourceKey, "name", "rule_test"),
					resource.TestCheckResourceAttr(testDayuL7RuleV2ResourceKey, "source_type", "1"),
					resource.TestCheckResourceAttr(testDayuL7RuleV2ResourceKey, "source_list.#", "1"),
					resource.TestCheckResourceAttr(testDayuL7RuleV2ResourceKey, "switch", "false"),
					resource.TestCheckResourceAttr(testDayuL7RuleV2ResourceKey, "protocol", "http"),
				),
			},
		},
	})
}

func testAccCheckDayuL7RuleV2Destroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testDayuL7RuleResourceName {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		items := strings.Split(rs.Primary.ID, FILED_SP)
		if len(items) < 3 {
			return fmt.Errorf("broken ID of L7 rule")
		}
		resourceType := items[0]
		resourceId := items[1]
		ruleId := items[2]

		service := DayuService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, _, has, err := service.DescribeL7Rule(ctx, resourceType, resourceId, ruleId)
		if err != nil {
			_, _, has, err = service.DescribeL7Rule(ctx, resourceType, resourceId, ruleId)
		}
		if err != nil {
			return err
		}
		if !has {
			return nil
		} else {
			return fmt.Errorf("delete L7 rule %s fail, still on server", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckDayuL7RuleV2Exists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		items := strings.Split(rs.Primary.ID, FILED_SP)
		if len(items) < 3 {
			return fmt.Errorf("broken ID of L7 rule")
		}
		resourceType := items[0]
		resourceId := items[1]
		ruleId := items[2]

		service := DayuService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, _, has, err := service.DescribeL7Rule(ctx, resourceType, resourceId, ruleId)
		if err != nil {
			_, _, has, err = service.DescribeL7Rule(ctx, resourceType, resourceId, ruleId)
		}
		if err != nil {
			return err
		}
		if has {
			return nil
		} else {
			return fmt.Errorf("L7 rule %s not found on server", rs.Primary.ID)

		}
	}
}

const testAccDayuL7RuleV2 string = `
resource "tencentcloud_dayu_l7_rule_v2" "test_rule" {
  resource_type         = "bgpip"
  resource_id 			= "%s"
  resource_ip			= "%s"
  name					= "rule_test"
  domain				= "zhaoshaona.com"
  protocol				= "http"
  switch				= true
  source_type			= 2
  source_list 			= ["1.1.1.1:80","2.2.2.2"]
  health_check_switch	= true
  health_check_code		= 31
  health_check_interval = 30
  health_check_method	= "GET"
  health_check_path		= "/"
  health_check_health_num = 5
  health_check_unhealth_num = 10
}
`
const testAccDayuL7RuleV2Update string = `
resource "tencentcloud_dayu_l7_rule_v2" "test_rule" {
  resource_type         = "bgpip"
  resource_id 			= "%s"
  resource_ip			= "%s"
  name					= "rule_test"
  domain				= "zhaoshaona.com"
  protocol				= "http"
  switch				= false
  source_type			= 1
  source_list 			= ["zhaoshaona.com"]
}
`
