resource "exec" "foo" {
  command = "/bin/ls"
  only_if = "[\"$(/Users/gosuri/Dropbox/projects/go/src/github.com/gosuri/terraform-exec/tester)\" == \"1234\"]"
}
