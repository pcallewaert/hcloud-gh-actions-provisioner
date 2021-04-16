package main

import (
	"strings"

	"github.com/namsral/flag"
	"github.com/pcallewaert/hcloud-gh-actions-provisioner/pkg"
	"github.com/sirupsen/logrus"
)

var (
	loglevel           string
	nrOfBuilders       int
	hcloudToken        string
	githubPat          string
	githubOwner        string
	imageSnapshot      int
	hcloudFirewallName string
	hcloudLocation     string
	hcloudServerType   string
	namePrefix         string
	staticLabels       string
)

func main() {
	logrus.Info("Starting hcloud-gh-actions-provisioner...")
	parseConfig()
	labels := strings.Split(staticLabels, ",")
	sa := pkg.NewServerAdmin(githubPat, githubOwner, hcloudToken, hcloudFirewallName, hcloudLocation, hcloudServerType)
	if err := sa.ScaleTo(nrOfBuilders, imageSnapshot, namePrefix, labels); err != nil {
		logrus.Fatal(err)
	}
}

func parseConfig() {
	flag.StringVar(&loglevel, "loglevel", "INFO", "Log level")
	flag.IntVar(&nrOfBuilders, "number-of-builders", -1, "The number of builders that have to be scaled")
	flag.IntVar(&imageSnapshot, "image-snapshot", -1, "Image ID of the snapshot to use")
	flag.StringVar(&hcloudToken, "hcloud-token", "", "Hetzner Cloud API Token")
	flag.StringVar(&hcloudFirewallName, "hcloud-firewall-name", "", "Hetzner Firewall Name")
	flag.StringVar(&hcloudLocation, "hcloud-location", "fsn1", "Hetzner Location")
	flag.StringVar(&hcloudServerType, "hcloud-server-type", "cpx21", "Hetzner Server type")
	flag.StringVar(&githubPat, "github-pat", "", "Github Personal Access Token")
	flag.StringVar(&githubOwner, "github-owner", "", "Github Organisation owner")
	flag.StringVar(&namePrefix, "name-prefix", "hcloud-github-actions-", "Name prefix of the servers")
	flag.StringVar(&staticLabels, "static-labels", "", "Labels that are added to Github runner and hetzner server")

	flag.Parse()
	if lvl, err := logrus.ParseLevel(loglevel); err != nil {
		logrus.Warnf("Unable to parse %s as loglevel, falling back to INFO: %v", loglevel, err)
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(lvl)
	}
}
