package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudGaapSecurityRule_basic(t *testing.T) {
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapSecurityRuleDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapSecurityRuleBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapSecurityRuleExists("tencentcloud_gaap_security_rule.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "cidr_ip", "1.1.1.1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "action", "ACCEPT"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "port", "80"),
				),
			},
			{
				ResourceName:      "tencentcloud_gaap_security_rule.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudGaapSecurityRule_drop(t *testing.T) {
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapSecurityRuleDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapSecurityRuleDrop,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapSecurityRuleExists("tencentcloud_gaap_security_rule.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "cidr_ip", "1.1.1.1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "action", "DROP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "port", "80"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapSecurityRule_name(t *testing.T) {
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapSecurityRuleDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapSecurityRuleWithName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapSecurityRuleExists("tencentcloud_gaap_security_rule.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "cidr_ip", "1.1.1.1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "action", "ACCEPT"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "port", "80"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "name", "ci-test-gaap-sr"),
				),
			},
			{
				Config: testAccGaapSecurityRuleUpdateName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapSecurityRuleExists("tencentcloud_gaap_security_rule.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "name", "ci-test-gaap-sr-new"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapSecurityRule_ipSubnet(t *testing.T) {
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapSecurityRuleDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapSecurityRuleIpSubnet,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapSecurityRuleExists("tencentcloud_gaap_security_rule.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "cidr_ip", "192.168.1.0/24"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "action", "ACCEPT"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "port", "80"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapSecurityRule_allProtocols(t *testing.T) {
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapSecurityRuleDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapSecurityRuleAllProtocols,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapSecurityRuleExists("tencentcloud_gaap_security_rule.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "cidr_ip", "1.1.1.1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "action", "ACCEPT"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "protocol", "ALL"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "port", "ALL"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapSecurityRule_AllPorts(t *testing.T) {
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapSecurityRuleDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapSecurityRuleAllPorts,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapSecurityRuleExists("tencentcloud_gaap_security_rule.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "cidr_ip", "1.1.1.1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "action", "ACCEPT"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "port", "ALL"),
				),
			},
		},
	})
}

func testAccCheckGaapSecurityRuleExists(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no security rule ID is set")
		}

		service := GaapService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		rule, err := service.DescribeSecurityRule(context.TODO(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if rule == nil {
			return errors.New("no security rule not exist")
		}

		*id = rs.Primary.ID

		return nil
	}
}

func testAccCheckGaapSecurityRuleDestroy(id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*TencentCloudClient).apiV3Conn
		service := GaapService{client: client}

		rule, err := service.DescribeSecurityRule(context.TODO(), *id)
		if err != nil {
			return err
		}

		if rule != nil {
			return errors.New("security rule still exists")
		}

		return nil
	}
}

var testAccGaapSecurityRuleBasic = fmt.Sprintf(`
resource tencentcloud_gaap_security_policy "foo" {
  proxy_id = "%s"
  action   = "ACCEPT"
}

resource tencentcloud_gaap_security_rule "foo" {
  policy_id = tencentcloud_gaap_security_policy.foo.id
  cidr_ip   = "1.1.1.1"
  action    = "ACCEPT"
  protocol  = "TCP"
  port      = "80"
}
`, defaultGaapProxyId)

var testAccGaapSecurityRuleWithName = fmt.Sprintf(`
resource tencentcloud_gaap_security_policy "foo" {
  proxy_id = "%s"
  action   = "ACCEPT"
}

resource tencentcloud_gaap_security_rule "foo" {
  name      = "ci-test-gaap-sr"
  policy_id = tencentcloud_gaap_security_policy.foo.id
  cidr_ip   = "1.1.1.1"
  action    = "ACCEPT"
  protocol  = "TCP"
  port      = "80"
}
`, defaultGaapProxyId)

var testAccGaapSecurityRuleDrop = fmt.Sprintf(`
resource tencentcloud_gaap_security_policy "foo" {
  proxy_id = "%s"
  action   = "ACCEPT"
}

resource tencentcloud_gaap_security_rule "foo" {
  policy_id = tencentcloud_gaap_security_policy.foo.id
  cidr_ip   = "1.1.1.1"
  action    = "DROP"
  protocol  = "TCP"
  port      = "80"
}
`, defaultGaapProxyId)

var testAccGaapSecurityRuleUpdateName = fmt.Sprintf(`
resource tencentcloud_gaap_security_policy "foo" {
  proxy_id = "%s"
  action   = "ACCEPT"
}

resource tencentcloud_gaap_security_rule "foo" {
  name      = "ci-test-gaap-sr-new"
  policy_id = tencentcloud_gaap_security_policy.foo.id
  cidr_ip   = "1.1.1.1"
  action    = "ACCEPT"
  protocol  = "TCP"
  port      = "80"
}
`, defaultGaapProxyId)

var testAccGaapSecurityRuleIpSubnet = fmt.Sprintf(`
resource tencentcloud_gaap_security_policy "foo" {
  proxy_id = "%s"
  action   = "ACCEPT"
}

resource tencentcloud_gaap_security_rule "foo" {
  policy_id = tencentcloud_gaap_security_policy.foo.id
  cidr_ip   = "192.168.1.0/24"
  action    = "ACCEPT"
  protocol  = "TCP"
  port      = "80"
}
`, defaultGaapProxyId)

var testAccGaapSecurityRuleAllProtocols = fmt.Sprintf(`
resource tencentcloud_gaap_security_policy "foo" {
  proxy_id = "%s"
  action   = "ACCEPT"
}

resource tencentcloud_gaap_security_rule "foo" {
  policy_id = tencentcloud_gaap_security_policy.foo.id
  cidr_ip   = "1.1.1.1"
  action    = "ACCEPT"
}
`, defaultGaapProxyId)

var testAccGaapSecurityRuleAllPorts = fmt.Sprintf(`
resource tencentcloud_gaap_security_policy "foo" {
  proxy_id = "%s"
  action   = "ACCEPT"
}

resource tencentcloud_gaap_security_rule "foo" {
  policy_id = tencentcloud_gaap_security_policy.foo.id
  cidr_ip   = "1.1.1.1"
  action    = "ACCEPT"
  protocol  = "TCP"
}
`, defaultGaapProxyId)
