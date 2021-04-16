package pkg

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v33/github"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/sirupsen/logrus"
)

const userdata = `
#cloud-config
disable_root: 1
ssh_pwauth: 0
write_files:
  - content: |
      #!/bin/bash
      token=$1
      labels=$2
      version=$(curl https://api.github.com/repos/actions/runner/releases/latest 2> /dev/null | jq -r .tag_name)
      version_stripped=${version#"v"}
      cd $HOME
      mkdir -p actions-runner
      cd actions-runner
      wget https://github.com/actions/runner/releases/download/${version}/actions-runner-linux-x64-${version_stripped}.tar.gz
      tar xvfz actions-runner-linux-x64-${version_stripped}.tar.gz
      ./config.sh --url https://github.com/%s --token ${token} --labels ${labels} --unattended
      echo "ImageOS=ubuntu20" >> .env
      sudo ./svc.sh install
      sudo ./svc.sh start
    owner: runner:runner
    permissions: '755'
    path: /home/runner/install-gh-runner.sh
runcmd:
  - sudo -H -u runner /home/runner/install-gh-runner.sh %s %s
`

type ServerAdmin struct {
	githubOwner        string
	githubPat          string
	githubClient       *github.Client
	hcloudClient       *hcloud.Client
	hcloudFirewallName string
	hcloudLocation     string
	hcloudServerType   string
}

func NewServerAdmin(githubPat, githubOwner, hcloudToken, hcloudFirewallName, hcloudLocation, hcloudServerType string) *ServerAdmin {
	githubClient := setupGithubClient(githubPat)
	hcloudClient := setupHcloudClient(hcloudToken)
	return &ServerAdmin{
		githubClient:       githubClient,
		hcloudClient:       hcloudClient,
		githubOwner:        githubOwner,
		githubPat:          githubPat,
		hcloudFirewallName: hcloudFirewallName,
		hcloudServerType:   hcloudServerType,
		hcloudLocation:     hcloudLocation,
	}
}

func (sa *ServerAdmin) ScaleTo(number, imageSnapshot int, namePrefix string, staticLabels []string) error {
	listCurrentlyRunning, err := sa.listRunners(namePrefix)
	if err != nil {
		return err
	}
	logrus.Debugf("We have %d running servers, our target is %d", len(listCurrentlyRunning), number)
	switch delta := number - len(listCurrentlyRunning); {
	case delta < 0:
		logrus.Debug("We have to scale down")
		for i := 0; i < -delta; i++ {
			server, err := sa.findServerToRemove(listCurrentlyRunning, namePrefix)
			if err != nil {
				return err
			}
			if server != nil {
				logrus.Debugf("Removing server %s", server.GetName())
				if err := sa.removeServer(server.GetID(), server.GetName()); err != nil {
					logrus.Warnf("Unable to remove the server %s: %v", server.GetName(), err)
				}
			}
		}
		return nil
	case delta > 0:
		logrus.Debug("We have to scale up")
		for i := 0; i < delta; i++ {
			serverName := fmt.Sprintf("%s%d", namePrefix, time.Now().UnixNano())
			logrus.Debugf("Spinning up server %s", serverName)
			sa.spinUpServer(serverName, userdata, imageSnapshot, staticLabels)
		}
		return nil
	default:
		logrus.Debug("Already on target number of servers")
		return nil
	}
}

// listRunners only returns the runners we created
func (sa *ServerAdmin) listRunners(namePrefix string) (result []*github.Runner, err error) {
	runners, _, err := sa.githubClient.Actions.ListOrganizationRunners(context.Background(), sa.githubOwner, nil)
	if err != nil {
		return
	}
	for _, x := range runners.Runners {
		if strings.HasPrefix(x.GetName(), namePrefix) {
			result = append(result, x)
		}
	}
	return
}

func (sa *ServerAdmin) findServerToRemove(runners []*github.Runner, namePrefix string) (*github.Runner, error) {
	var oldestServer *github.Runner
	var oldestDate time.Time
	for _, x := range runners {
		if x.GetBusy() {
			logrus.Debugf("%s is not idle, so we can't delete it", x.GetName())
			continue
		}
		date := strings.TrimPrefix(x.GetName(), namePrefix)
		i, err := strconv.ParseInt(date, 10, 64)
		if err != nil {
			return nil, err
		}
		tm := time.Unix(0, i)
		if oldestDate.IsZero() || tm.Before(oldestDate) {
			oldestServer = x
			oldestDate = tm
		}
	}
	return oldestServer, nil
}

func (sa *ServerAdmin) removeServer(runnerId int64, serverName string) error {
	_, err := sa.githubClient.Actions.RemoveOrganizationRunner(context.Background(), sa.githubOwner, runnerId)
	if err != nil {
		return err
	}
	server, _, err := sa.hcloudClient.Server.GetByName(context.Background(), serverName)
	if err != nil {
		return err
	}
	_, err = sa.hcloudClient.Server.Delete(context.Background(), server)
	if err != nil {
		return err
	}
	return nil
}

func (sa *ServerAdmin) spinUpServer(serverName, userdata string, imageSnapshot int, staticLabels []string) {
	image, _, err := sa.hcloudClient.Image.GetByID(context.Background(), imageSnapshot)
	if err != nil {
		logrus.Fatalf("Error retrieving image: %v", err)
	}
	serverType, _, err := sa.hcloudClient.ServerType.GetByName(context.Background(), sa.hcloudServerType)
	if err != nil {
		logrus.Fatalf("Error retrieving servertype: %v", err)
	}
	location, _, err := sa.hcloudClient.Location.GetByName(context.Background(), sa.hcloudLocation)
	if err != nil {
		logrus.Fatalf("Error retrieving location: %v", err)
	}
	var fw *hcloud.Firewall
	if sa.hcloudFirewallName != "" {
		fw, _, err = sa.hcloudClient.Firewall.GetByName(context.Background(), sa.hcloudFirewallName)
		if err != nil {
			logrus.Fatalf("Error retrieving firewall: %v", err)
		}
	}
	token, _, err := sa.githubClient.Actions.CreateOrganizationRegistrationToken(context.Background(), sa.githubOwner)
	if err != nil {
		logrus.Errorf("Error retrieving github repos: %v", err)
		os.Exit(1)
	}
	logrus.Debugf("Organization Registration Token: %s", token.GetToken())
	labels := []string{"hcloud-runner"}
	labels = append(labels, staticLabels...)
	labelsJoined := strings.Join(labels, ",")
	formattedUserData := fmt.Sprintf(userdata, sa.githubOwner, token.GetToken(), labelsJoined)
	logrus.Debugf("cloud-init config: %s:", formattedUserData)
	// Setup hcloud server
	hetznerLabels := map[string]string{"runner": "automated"}
	for _, x := range staticLabels {
		hetznerLabels[x] = ""
	}
	opts := hcloud.ServerCreateOpts{
		Image:      image,
		ServerType: serverType,
		Location:   location,
		Name:       serverName,
		Labels:     hetznerLabels,
		Firewalls:  []*hcloud.ServerCreateFirewall{{Firewall: *fw}},
		UserData:   formattedUserData,
	}
	_, _, err = sa.hcloudClient.Server.Create(context.Background(), opts)
	if err != nil {
		logrus.Errorf("Failed to create server: %v", err)
	}
	logrus.Info("server created")
}
