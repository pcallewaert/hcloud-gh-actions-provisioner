package main

import (
	"github.com/namsral/flag"
	"github.com/pcallewaert/hcloud-gh-actions-provisioner/pkg"
	"github.com/sirupsen/logrus"
)

var (
	loglevel      string
	nrOfBuilders  int
	hcloudToken   string
	githubPat     string
	githubOwner   string
	imageSnapshot int
)

func main() {
	logrus.Info("Starting hcloud-gh-actions-provisioner...")
	parseConfig()
	sa := pkg.NewServerAdmin(githubPat, githubOwner, hcloudToken)
	if err := sa.ScaleTo(nrOfBuilders, imageSnapshot); err != nil {
		logrus.Fatal(err)
	}
}

func parseConfig() {
	flag.StringVar(&loglevel, "loglevel", "INFO", "Log level")
	flag.IntVar(&nrOfBuilders, "number-of-builders", -1, "The number of builders that have to be scaled")
	flag.IntVar(&imageSnapshot, "image-snapshot", -1, "Image ID of the snapshot to use")
	flag.StringVar(&hcloudToken, "hcloud-token", "", "Hetzner Cloud API Token")
	flag.StringVar(&githubPat, "github-pat", "", "Github Personal Access Token")
	flag.StringVar(&githubOwner, "github-owner", "", "Github Organisation owner")

	flag.Parse()
	if lvl, err := logrus.ParseLevel(loglevel); err != nil {
		logrus.Warnf("Unable to parse %s as loglevel, falling back to INFO: %v", loglevel, err)
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(lvl)
	}
}
