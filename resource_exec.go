package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

// Default timeout for the command is 60s
const defaultTimeout = 60

// ExecCmd holds data necessary for a command to run
type ExecCmd struct {
	Cmd     string
	Timeout int
}

// Terraform schema for the 'exec' resource that is
// used in the provider configuration
func resourceExec() *schema.Resource {
	return &schema.Resource{
		Create: resourceExecCreate,
		Read:   resourceExecRead,
		Update: resourceExecUpdate,
		Delete: resourceExecDelete,

		Schema: map[string]*schema.Schema{
			"command": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"only_if": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				//ForceNew: true,
			},
			"timeout": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				//ForceNew: true,
			},
			"output": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceExecCreate(d *schema.ResourceData, m interface{}) error {
	timeout := d.Get("timeout").(int)

	cmd := &ExecCmd{
		Cmd:     d.Get("command").(string),
		Timeout: timeout,
	}

	onlyIf := &ExecCmd{
		Cmd:     d.Get("only_if").(string),
		Timeout: timeout,
	}

	// Set the id of the resource
	id := strings.Join(strings.Split(cmd.Cmd, " "), "-")
	d.SetId(id)

	// run the only_if command and continue only on success
	if onlyIf.Cmd != "" {
		_, err := ExecuteCmd(onlyIf)
		if err != nil {
			log.Printf("[DEBUG] Skipped execution (%s): `%s` exited with a failed state", cmd.Cmd, onlyIf.Cmd)
			// stop executing the command by returning nil
			return nil
		}
	}

	// run the actual command
	out, err := ExecuteCmd(cmd)
	if err != nil {
		d.Set("output", "")
	}
	log.Printf("[DEBUG] Command Output (%s): %s", cmd.Cmd, out)
	d.Set("output", out)
	return nil
}

func resourceExecUpdate(d *schema.ResourceData, m interface{}) error {
	timeout := d.Get("timeout").(int)

	cmd := &ExecCmd{
		Cmd:     d.Get("command").(string),
		Timeout: timeout,
	}

	onlyIf := &ExecCmd{
		Cmd:     d.Get("only_if").(string),
		Timeout: timeout,
	}

	// run the only_if command and continue only on success
	if onlyIf.Cmd != "" {
		_, err := ExecuteCmd(onlyIf)
		if err != nil {
			log.Printf("[DEBUG] Skipped execution (%s): `%s` exited with a failed state", cmd.Cmd, onlyIf.Cmd)
			// stop executing the command by returning nil
			return nil
		}
	}

	// run the acctual command
	out, err := ExecuteCmd(cmd)
	if err != nil {
		d.Set("output", "")
		return nil
	}

	log.Printf("[DEBUG] Command Output (%s): %s", cmd.Cmd, out)
	d.Set("output", out)
	return nil
}

func resourceExecRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceExecDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func ExecuteCmd(command *ExecCmd) (output string, err error) {
	// Wrap the command in a temp file
	var cmdWrapper *os.File
	cmdWrapper, err = ioutil.TempFile("", "exec")

	if err != nil {
		log.Fatal(fmt.Sprintf("Error while creating temp file: %s", err))
		return "", err
	}
	defer cmdWrapper.Close()

	if err = os.Chmod(cmdWrapper.Name(), 0755); err != nil {
		log.Fatal(fmt.Sprintf("Error while making the file executable: %s", err))
	}

	// Run the command in the current working directory
	var path string
	path, err = os.Getwd()
	if err != nil {
		log.Fatal(fmt.Sprintf("Error getting pwd: %s", err))
		return "", err
	}

	code := fmt.Sprintf("#!/usr/bin/env /bin/sh\ncd %s\n%s", path, command.Cmd)

	if err = ioutil.WriteFile(cmdWrapper.Name(), []byte(code), 0755); err != nil {
		log.Fatal(fmt.Sprintf("Error while writing to temp file: %s", err))
		return "", err
	}

	if command.Timeout == 0 {
		command.Timeout = defaultTimeout
	}

	// Run the command in a channel using select statement
	// with time.After for timingout calls that run too long
	var out []byte
	timeout := make(chan error)
	go func() {
		out, err = exec.Command(cmdWrapper.Name()).Output()
		timeout <- err
	}()

	select {
	case err := <-timeout:
		if err != nil {
			return "", err
		}
	case <-time.After(time.Duration(command.Timeout) * time.Second):
		log.Printf("[DEBUG] Execution (%s) timedout in %ds", command.Cmd, command.Timeout)
		return "", nil
	}

	return string(out[:]), nil
}
