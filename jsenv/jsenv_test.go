package jsenv_test

import (
	"github.com/lithictech/runtime-js-env/jsenv"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"testing"
)

func TestJsenv(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "jsenv package Suite")
}

var _ = Describe("jsenv", func() {
	mustWrite := func(contents string) string {
		tf, err := ioutil.TempFile("", "jsenv")
		Expect(err).ToNot(HaveOccurred())
		_, err = tf.WriteString(contents)
		Expect(err).ToNot(HaveOccurred())
		return tf.Name()
	}

	mustRead := func(path string) string {
		b, err := ioutil.ReadFile(path)
		Expect(err).ToNot(HaveOccurred())
		return string(b)
	}

	var addedEnvs []string
	BeforeEach(func() {
		addedEnvs = []string{}
	})
	AfterEach(func() {
		for _, e := range addedEnvs {
			Expect(os.Unsetenv(e)).To(Succeed())
		}
	})
	mustSetenv := func(k, v string) {
		addedEnvs = append(addedEnvs, k)
		Expect(os.Setenv(k, v)).To(Succeed())
	}

	Describe("Install", func() {
		It("installs config at window._jsenv", func() {
			index := mustWrite(`<html><head></head><body /></html>`)
			mustSetenv("REACT_APP_YO", "ma")
			mustSetenv("NODE_BLAH", "staging")
			Expect(jsenv.InstallAt(index, jsenv.DefaultConfig)).To(Succeed())
			Expect(mustRead(index)).To(HavePrefix(`<html><head><script id="jsenv">
  window._jsenv = {
    "NODE_BLAH": "staging",
    "REACT_APP_YO": "ma",
  };
</script></head><body></body></html>`))
		})
		It("noops on reinstall", func() {
			index := mustWrite(`<html><head></head><body /></html>`)
			mustSetenv("HEROKU_APP_NAME", "great")
			Expect(jsenv.InstallAt(index, jsenv.DefaultConfig)).To(Succeed())
			Expect(jsenv.InstallAt(index, jsenv.DefaultConfig)).To(Succeed())
			Expect(mustRead(index)).To(HavePrefix(`<html><head><script id="jsenv">
  window._jsenv = {
    "HEROKU_APP_NAME": "great",
  };
</script></head><body></body></html>`))
		})
		It("handles env vars with multiple =", func() {
			index := mustWrite(`<html><head></head><body /></html>`)
			mustSetenv("REACT_APP_MYVAR", "test1=yeah")
			Expect(jsenv.InstallAt(index, jsenv.DefaultConfig)).To(Succeed())
			Expect(mustRead(index)).To(HavePrefix(`<html><head><script id="jsenv">
  window._jsenv = {
    "REACT_APP_MYVAR": "test1=yeah",
  };
</script></head><body></body></html>`))
		})
		It("escapes var values with quote", func() {
			index := mustWrite(`<html><head></head><body /></html>`)
			mustSetenv("REACT_APP_MYVAR", `x'"quoted"'y`)
			Expect(jsenv.InstallAt(index, jsenv.DefaultConfig)).To(Succeed())
			Expect(mustRead(index)).To(HavePrefix(`<html><head><script id="jsenv">
  window._jsenv = {
    "REACT_APP_MYVAR": "x'\"quoted\"'y",
  };
</script></head><body></body></html>`))
		})
		It("handles non-html string (is valid html5)", func() {
			index := mustWrite(`this is not html`)
			Expect(jsenv.InstallAt(index, jsenv.DefaultConfig)).To(Succeed())
			Expect(mustRead(index)).To(HavePrefix(`<html><head><script id="jsenv">
  window._jsenv = {
  };
</script></head><body>this is not html</body></html>`))
		})
		It("adds head if missing", func() {
			index := mustWrite(`<html><body>`)
			Expect(jsenv.InstallAt(index, jsenv.DefaultConfig)).To(Succeed())
			Expect(mustRead(index)).To(HavePrefix(`<html><head><script id="jsenv">
  window._jsenv = {
  };
</script></head><body></body></html>`))
		})
	})
})
