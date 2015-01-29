package main

import (
	//"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	//"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

type ExecCmd struct {
	Cmd     string
	Timeout int
}

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
			},
			"output": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceExecCreate(d *schema.ResourceData, m interface{}) error {
	command := d.Get("command").(string)
	onlyIfCmd := d.Get("only_if").(string)

	id := fmt.Sprintf("%s: %s", time.Now().Format("20060102150405"), command)
	d.SetId(id)

	// run the only_if command and continue only on success
	if onlyIfCmd != "" {
		_, err := ExecuteCmd(&ExecCmd{Cmd: onlyIfCmd})
		if err != nil {
			// stop executing the command by returning nil
			return nil
		}
	}

	// Run the actual command
	cmd := &ExecCmd{Cmd: command}
	out, err := ExecuteCmd(cmd)
	if err != nil {
		d.Set("output", out)
		return err
	}
	log.Printf("[DEBUG] Exec output: %s", out)
	d.Set("output", out)
	return nil
}

func resourceExecRead(d *schema.ResourceData, m interface{}) error {
	// onlyIfCmd := d.Get("only_if").(string)
	// command := d.Get("command").(string)
	// id := fmt.Sprintf("%s: %s", time.Now().Format("20060102150405"), command)
	// d.SetId(id)
	// // run the only_if command
	// if onlyIfCmd != "" {
	// 	_, err := exec.Command(onlyIfCmd).Output()
	// 	if err != nil {
	// 		// stop executing the command by returning nil
	// 		return err
	// 		//return nil
	// 	}
	// }
	// d.SetId("")
	return nil
}

func resourceExecUpdate(d *schema.ResourceData, m interface{}) error {
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

	var out []byte
	log.Printf("[DEBUG] executing: %s", cmdWrapper.Name())
	out, err = exec.Command(cmdWrapper.Name()).Output()
	if err != nil {
		return output, err
	}
	output = string(out[:])
	return output, nil
}
