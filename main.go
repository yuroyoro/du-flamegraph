package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "du-flamegraph"
	app.Usage = "visualize disk usage as flamegraph"
	app.ArgsUsage = "[FILE]"
	app.HideHelp = true

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "width, w",
			Value: 1200,
			Usage: "width of image (default 1200)",
		},
		cli.IntFlag{
			Name:  "height, h",
			Value: 16,
			Usage: "height of each frame (default 16)",
		},
		cli.StringFlag{
			Name:  "flamegraph-script",
			Usage: "path of flamegraph.pl. if not given, find the script from $PATH",
		},
		cli.StringFlag{
			Name:  "out",
			Value: "./du-flamegraph.svg",
			Usage: "distination path of grenerated flamegraph. default is ./du-flamegraph.svg",
		},
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "show verbose log",
		},
	}

	app.Action = generate
	app.Run(os.Args)
}

var (
	verbose bool
	bases   []string
)

func generate(c *cli.Context) error {
	verbose = c.Bool("verbose")
	trace("Args : %+v", os.Args)

	bases = basePaths(c)

	m := map[string]float64{}

	for _, path := range bases {
		traverse(path, m)
	}

	trace("m : %+v", m)
	writeFlameGraph(c, m)

	return nil
}

func trace(format string, v ...interface{}) {
	if !verbose {
		return
	}

	log.Printf(format+"\n", v...)
}

func traverse(node string, m map[string]float64) {
	filepath.Walk(node, func(path string, info os.FileInfo, err error) error {
		trace("%s : %d : dir %v", path, info.Size(), info.IsDir())

		abspath, err := filepath.Abs(path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			if _, ok := m[abspath]; ok {
				trace("%s : already traversed", path)
				return filepath.SkipDir
			}

			m[abspath] = 0
			return nil
		}

		dir := filepath.Dir(abspath)
		m[dir] += float64(info.Size())

		return nil
	})
}

func basePaths(c *cli.Context) []string {
	res := make([]string, 0, len(c.Args()))
	for _, path := range c.Args() {
		abspath, err := filepath.Abs(path)
		if err != nil {
			log.Fatalf("basePaths: cound not resolve absolute path : %s : %s", path, err)
		}
		res = append(res, abspath)
	}

	return res
}

func relpath(path string, bases []string) (string, string) {
	for _, base := range bases {
		rel, err := filepath.Rel(base, path)
		if err == nil {
			if rel == "." {
				rel = base
			}
			return base, rel
		}
	}

	return path, path
}

func resolveScript(c *cli.Context) string {
	if script := c.String("flamegraph-script"); script != "" {
		trace("script: %s", script)
		if _, err := os.Stat(script); err == nil {
			return script
		}
	}

	if script, err := exec.LookPath("flamegraph.pl"); err == nil {
		return script
	}

	log.Fatalf("flamegraph.pl is not found in $PATH")

	return ""
}

func writeFlameGraph(c *cli.Context, m map[string]float64) {
	writeSVG(c, flamegraphData(m))
}

func flamegraphData(m map[string]float64) []byte {
	var b bytes.Buffer

	for path, size := range m {
		base, rel := relpath(path, bases)
		trace("rel : %s => %s", path, rel)
		if !filepath.IsAbs(rel) {

			arr := []string{base}
			arr = append(arr, strings.Split(rel, string(os.PathSeparator))...)
			rel = strings.Join(arr, ";")
		}

		if _, err := b.WriteString(fmt.Sprintf("%s %.2f\n", rel, size)); err != nil {
			log.Fatalf("flamegraphData: cound not write datat: %s", err)
		}
	}

	return b.Bytes()
}

func writeSVG(c *cli.Context, data []byte) {
	args := []string{
		"--title",
		"flame Graph of disk usage",
		"--width",
		fmt.Sprintf("%d", c.Int("width")),
		"--height",
		fmt.Sprintf("%d", c.Int("height")),
		"--countname",
		"bytes",
		"--nametype",
		"Path",
		"--colors",
		"aqua",
	}

	script := resolveScript(c)
	cmd := exec.Command(script, args...)
	cmd.Stdin = bytes.NewReader(data)
	cmd.Stderr = os.Stderr

	svg, err := cmd.Output()
	if err != nil {
		log.Fatalf("writeSVG: failed to run script %s : %s", script, err)
	}

	out := c.String("out")
	if err := ioutil.WriteFile(out, svg, 0644); err != nil {
		log.Fatalf("writeFlameGraph: cound not write to %s : %s", out, err)
	}
}
