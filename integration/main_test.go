package integration_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/onsi/gomega/gbytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"github.com/micahyoung/docker-layer-extract/integration/helpers"
)

var cliBin string
var imagePath string
var expectedLayerID string

var _ = Describe("Main", func() {
	BeforeSuite(func() {
		var err error
		var session *gexec.Session
		dockerfilePath := filepath.Join("fixtures", fmt.Sprintf("Dockerfile.%s", runtime.GOOS))
		imageTag := "docker-layer-extract-ci"

		cliBin, err = gexec.Build("github.com/micahyoung/docker-layer-extract")
		Expect(err).ToNot(HaveOccurred())

		buildCmd := helpers.DockerBuildCommand(dockerfilePath, imageTag)
		session, err = gexec.Start(buildCmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(session, 1*time.Minute).Should(gexec.Exit(0))

		inspectCmd := helpers.DockerInspectCommand(imageTag)
		inspectBuffer := new(bytes.Buffer)
		session, err = gexec.Start(inspectCmd, inspectBuffer, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(session, 1*time.Minute).Should(gexec.Exit(0))

		expectedLayerID, err = helpers.ParseInspectLatestLayerID(inspectBuffer.String())
		Expect(err).ToNot(HaveOccurred())

		imageDir, _ := ioutil.TempDir("", "image-tempdir")
		imagePath = filepath.Join(imageDir, "image.tar")
		saveCmd := helpers.DockerImageSaveCommand(imageTag, imagePath)
		session, err = gexec.Start(saveCmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(session, 1*time.Minute).Should(gexec.Exit(0))

		Expect(err).ToNot(HaveOccurred())
	})

	AfterSuite(func() {
		os.Remove(imagePath)
		gexec.CleanupBuildArtifacts()
	})

	It("lists contents of saved image file", func() {
		command := exec.Command(cliBin, "list", "-i", imagePath)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)

		Eventually(session, 1*time.Minute).Should(gexec.Exit(0))
		Eventually(session).Should(gbytes.Say(expectedLayerID))

		Expect(err).ToNot(HaveOccurred())
	})
})
