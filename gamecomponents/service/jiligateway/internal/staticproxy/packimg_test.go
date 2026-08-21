package staticproxy

import (
	"archive/tar"
	"bytes"
	"io/fs"
	"log"
	"os"
	"path"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestPackImg(t *testing.T) {
	wd := "/data/game/bin/cache/jilid.rslotszs001.com"
	os.Chdir(wd)

	fileSystem := os.DirFS(".")

	extmap := map[string]bool{}

	for _, e := range lo.Must(os.ReadDir("./")) {
		// fmt.Println(e.Name())
		var buf bytes.Buffer
		tw := tar.NewWriter(&buf)

		checkmap := map[string]bool{}
		fs.WalkDir(fileSystem, e.Name(), func(pth string, d fs.DirEntry, err error) error {
			if err != nil {
				log.Fatal(err)
			}
			if !d.Type().IsRegular() {
				return nil
			}

			ext := path.Ext(pth)
			extmap[ext] = true

			switch ext {
			case ".jpg" /*".plist",*/, ".png", ".webp":
				// fmt.Println(pth)
				basename := path.Base(pth)
				assert.False(t, checkmap[basename])
				checkmap[basename] = true

				content := lo.Must(os.ReadFile(pth))

				hdr := &tar.Header{
					Name: basename,
					Mode: 0644,
					Size: int64(len(content)),
				}
				if err := tw.WriteHeader(hdr); err != nil {
					log.Fatal(err)
				}
				if _, err := tw.Write(content); err != nil {
					log.Fatal(err)
				}
			}
			// if strings.HasSuffix(path, "")
			return nil
		})
		tw.Flush()
		os.WriteFile(e.Name()+".tar", buf.Bytes(), 0644)
	}

	// fmt.Println(extmap)
}

func TestPath(t *testing.T) {
	// fmt.Println(path.Clean("/a/b/c/"))
	assert.True(t, regIngame.MatchString("/en-US/ingame"))
	assert.True(t, regIngame.MatchString("/ingame"))
	assert.True(t, regIngame.MatchString("/en-US/intro"))
	assert.False(t, regIngame.MatchString("/en.US/intro"))
	assert.False(t, regIngame.MatchString("/en-US/x/ingame"))
	assert.False(t, regIngame.MatchString("/x/ingame"))
	assert.False(t, regIngame.MatchString("/x/ingame/xxx"))
	assert.False(t, regIngame.MatchString("/_nuxt/pages/ingame/action"))
	assert.True(t, regIngame.MatchString("/en-US/ingame/gamehistory/csh/xxx"))
}

func TestRmGTM(t *testing.T) {
	filepath := "/data/game/bin/cache/uat-history.jlfafafa3.com/en-US/intro-2.html"
	content, _ := os.ReadFile(filepath)

	// reg := regexp.MustCompile(`<script data-n-head="1" data-hid="gtm-script">.*</script>`)

	// // ret := reg.Find(content)
	// // os.Stdout.Write(ret)

	// newcontent := reg.ReplaceAll(content, []byte{})

	// reg2 := regexp.MustCompile(`<noscript data-n-head="1" data-hid="gtm-noscript" data-pbody="true">.*</noscript>`)

	// newcontent = reg2.ReplaceAll(newcontent, []byte{})

	newcontent := bytes.ReplaceAll(content, []byte("GTM-P8B36RZ"), []byte("GTM-1234567"))

	os.Stdout.Write(newcontent)
	os.WriteFile(filepath, newcontent, 0644)
}

func TestXxx(t *testing.T) {
	assert.True(t, regJILI.MatchString("/tks/JILI_NEW_1717655961.png"))
	assert.False(t, regJILI.MatchString("/tks/JILI1_NEW_1717655961.png"))
}
