package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr"
)

const defaultRegion = "cn-hangzhou"

type GetAuthTokenResponse struct {
	Data struct {
		AuthorizationToken string `json:"authorizationToken"`
		UserName           string `json:"tempUserName"`
	} `json:"data"`
}

func main() {
	var (
		repo      = getenv("PLUGIN_REPO")
		registry  = getenv("PLUGIN_REGISTRY")
		region    = getenv("PLUGIN_REGION")
		namespace = getenv("PLUGIN_NAMESPACE")
		key       = getenv("ACCESS_KEY", "PLUGIN_ACCESS_KEY")
		secret    = getenv("SECRET_ACCESS_KEY", "PLUGIN_SECRET_ACCESS_KEY")
		client    *cr.Client
		err       error
	)

	// set the region
	if region == "" {
		region = defaultRegion
	}

	if key != "" && secret != "" {
		client, err = cr.NewClientWithAccessKey(region, key, secret)
	} else {
		log.Fatalln("miss access key and access secret")
	}

	if err != nil {
		log.Fatal(err)
	}

	// default registry value
	if registry == "" {
		registry = fmt.Sprintf("registry.%s.aliyuncs.com", region)
	}

	fullReg := registry
	// split namespace and registry
	registry = strings.Split(registry, "/")[0]

	if namespace != "" && !strings.HasSuffix(fullReg, namespace) {
		fullReg = fmt.Sprintf("%s/%s", fullReg, namespace)
	}

	// must use the fully qualified repo name. If the
	// repo name does not have the registry prefix we
	// should prepend.
	if !strings.HasPrefix(repo, fullReg) {
		repo = fmt.Sprintf("%s/%s", fullReg, repo)
	}

	req := cr.CreateGetAuthorizationTokenRequest()
	req.Domain = fmt.Sprintf("cr.%s.aliyuncs.com", region)

	rawResp, err := client.GetAuthorizationToken(req)
	if err != nil {
		log.Fatal(err)
	}
	var resp GetAuthTokenResponse
	err = json.Unmarshal(rawResp.GetHttpContentBytes(), &resp)
	if err != nil {
		log.Fatal(err)
	}

	os.Setenv("PLUGIN_REPO", repo)
	os.Setenv("PLUGIN_REGISTRY", registry)
	os.Setenv("DOCKER_USERNAME", resp.Data.UserName)
	os.Setenv("DOCKER_PASSWORD", resp.Data.AuthorizationToken)

	// invoke the base docker plugin binary
	cmd := exec.Command("/kaniko/plugin.sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

}

func getenv(key ...string) (s string) {
	for _, k := range key {
		s = os.Getenv(k)
		if s != "" {
			return
		}
	}
	return
}
