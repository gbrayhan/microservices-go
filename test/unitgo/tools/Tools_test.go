package tools

import (
  "testing"

  "github.com/gbrayhan/microservices-go/tools"
)

func TestGetMD5Hash(t *testing.T) {
  if tools.GetMD5Hash("") != "d41d8cd98f00b204e9800998ecf8427e" {
    t.Error("Error in test tools.GetMD5Hash ")
  }

  if tools.GetMD5Hash("a") != "0cc175b9c0f1b6a831c399e269772661" {
    t.Error("Error in test tools.GetMD5Hash ")
  }

  if tools.GetMD5Hash("abc") != "900150983cd24fb0d6963f7d28e17f72" {
    t.Error("Error in test tools.GetMD5Hash ")
  }

  if tools.GetMD5Hash("message digest") != "f96b697d7cb7938d525a2f31aaf161d0" {
    t.Error("Error in test tools.GetMD5Hash ")
  }
  if tools.GetMD5Hash("abcdefghijklmnopqrstuvwxyz") != "c3fcd3d76192e4007dfb496cca67e13b" {
    t.Error("Error in test tools.GetMD5Hash ")
  }
}
