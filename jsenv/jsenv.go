package jsenv

import (
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

type Config struct {
	// The name of the environment variable on the window.
	// Default: _jsenv
	WindowVarName string
	// Prefixes of environment variables to render into the window config.
	// Default: REACT_APP_, NODE_, HEROKU_
	EnvPrefixes []string
	// Indentation of each entry in the config dict, because beauty is important
	// and who knows how your index.html is intented.
	// Default: "  " (two spaces)
	Indent string
}

var DefaultConfig = Config{}

// InstallAt reads the file at indexPath, and overwrites it with an index.html with config.
// See README for more info.
func InstallAt(indexPath string, cfg Config) error {
	if cfg.WindowVarName == "" {
		cfg.WindowVarName = "_jsenv"
	}
	if len(cfg.EnvPrefixes) == 0 {
		cfg.EnvPrefixes = []string{"REACT_APP_", "NODE_", "HEROKU_"}
	}
	if cfg.Indent == "" {
		cfg.Indent = "  "
	}

	indexFile, err := os.Open(indexPath)
	if err != nil {
		return err
	}
	doc, err := html.Parse(indexFile)
	if err != nil {
		return err
	}

	if !installIntoTree(doc, cfg) {
		return errors.New("malformed html must have <head> element")
	}

	tf, err := ioutil.TempFile("", "jsenv.index.html")
	if err != nil {
		return err
	}
	if err := html.Render(tf, doc); err != nil {
		return err
	}
	if err := os.Rename(tf.Name(), indexPath); err != nil {
		return err
	}
	return nil
}

func installIntoTree(n *html.Node, cfg Config) bool {
	// We only care about the head node; if we get any other node (html root?),
	// keep walking. Or if we hit body, don't bother to keep walking because we know it's not there.
	isBody := n.Type == html.ElementNode && n.Data == "body"
	if isBody {
		return false
	}
	isHead := n.Type == html.ElementNode && n.Data == "head"
	if !isHead {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if installIntoTree(c, cfg) {
				return true
			}
		}
		return false
	}
	// If we're already installed into a head child, break
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if getAttr(c.Attr, "id") == "jsenv" {
			return true
		}
	}

	scriptNode := &html.Node{
		Type: html.ElementNode,
		Data: "script",
		Attr: []html.Attribute{{Key: "id", Val: "jsenv"}},
	}
	scriptNode.AppendChild(&html.Node{
		Type: html.TextNode,
		Data: formatEnv(cfg),
	})
	n.AppendChild(scriptNode)
	return true
}

func getAttr(attrs []html.Attribute, key string) string {
	for _, a := range attrs {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

func formatEnv(cfg Config) string {
	envvars := pickCopyableVars(cfg)
	content := &strings.Builder{}
	content.WriteString(fmt.Sprintf("\n%swindow.%s = {\n", cfg.Indent, cfg.WindowVarName))
	for k, v := range envvars {
		escaped := strings.ReplaceAll(v, `"`, `\"`)
		content.WriteString(fmt.Sprintf(`%s%s"%s": "%s",`, cfg.Indent, cfg.Indent, k, escaped))
		content.WriteByte('\n')
	}
	content.WriteString(fmt.Sprintf("%s};\n", cfg.Indent))
	return content.String()
}

func pickCopyableVars(cfg Config) map[string]string {
	envvars := make(map[string]string, 8)
	environ := os.Environ()
	sort.Strings(environ)
	for _, ev := range environ {
		varParts := strings.SplitN(ev, "=", 2)
		varKey := varParts[0]
		varVal := varParts[1]
		if shouldCopyEnv(varKey, cfg) {
			envvars[varKey] = varVal
		}
	}
	return envvars
}

func shouldCopyEnv(key string, cfg Config) bool {
	for _, prefix := range cfg.EnvPrefixes {
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}
	return false
}
