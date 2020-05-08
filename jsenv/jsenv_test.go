package jsenv_test

import (
	"github.com/lithictech/runtime-js-env/jsenv"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rgalanakis/golangal"
	"io/ioutil"
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

	addEnvVar := golangal.EnvVars()

	Describe("Install", func() {
		It("installs config at window._jsenv", func() {
			index := mustWrite(`<html><head></head><body /></html>`)
			addEnvVar("REACT_APP_YO", "ma")
			addEnvVar("NODE_BLAH", "staging")
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
			addEnvVar("HEROKU_APP_NAME", "great")
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
			addEnvVar("REACT_APP_MYVAR", "test1=yeah")
			Expect(jsenv.InstallAt(index, jsenv.DefaultConfig)).To(Succeed())
			Expect(mustRead(index)).To(HavePrefix(`<html><head><script id="jsenv">
  window._jsenv = {
    "REACT_APP_MYVAR": "test1=yeah",
  };
</script></head><body></body></html>`))
		})
		It("escapes var values with quote", func() {
			index := mustWrite(`<html><head></head><body /></html>`)
			addEnvVar("REACT_APP_MYVAR", `x'"quoted"'y`)
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
