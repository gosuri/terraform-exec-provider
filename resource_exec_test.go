package main

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"exec": testAccProvider,
	}
}

func TestResourceExecCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccResourceExecDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccResourceExecConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("exec.foo", "output", "success\n"),
				),
			},
		},
	})
}

func TestResourceExecUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccResourceExecDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccResourceExecConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("exec.foo", "output", "success\n"),
				),
			},
			resource.TestStep{
				Config: testAccResourceExecConfig_basic_2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("exec.foo", "output", "success2\n"),
				),
			},
		},
	})
}

func TestResourceExecCreateTestFail(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccResourceExecDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccResourceExecConfig_test_fail,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExecResourceIsNil("exec.foo"),
				),
			},
			resource.TestStep{
				Config: testAccResourceExecConfig_test_pass,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("exec.foo", "output", "success\n"),
				),
			},
		},
	})
}

func testAccCheckExecResourceIsNil(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[r]
		if ok {
			return fmt.Errorf("Resource exists: %s", r)
		}
		return nil
	}
}

func testAccResourceExecDestroy(s *terraform.State) error {
	return nil
}

const testAccResourceExecConfig_basic = `
resource "exec" "foo" {
	command = "echo 'success'"
}
`
const testAccResourceExecConfig_basic_2 = `
resource "exec" "foo" {
	command = "echo 'success2'"
}
`
const testAccResourceExecConfig_test_pass = `
resource "exec" "foo" {
	command = "echo 'success'"
	only_if = "true"
}
`
const testAccResourceExecConfig_test_fail = `
resource "exec" "foo" {
	command = "echo 'success'"
	only_if = "false"
}
`
const testAccResourceExecConfig_timeout = `
resource "exec" "foo" {
	command = "sleep 2 && echo 'success'"
	timeout = 1
}
`
const testAccResourceExecConfig_fail = `
resource "exec" "foo" {
	command = "echo 'failure' >&2 && exit 1"
}
`
