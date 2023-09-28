package ease_test

import (
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/itler/go-lib/ease"
	"github.com/stretchr/testify/assert"
)

func testGet(uri string) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader("test content")),
	}, nil
}

func TestDownload(t *testing.T) {
	testfile := "test-generated.txt"
	testUri := "https://raw.githubusercontent.com/itler/go-lib/main/README.md"
	defer func() {
		os.Remove(testfile)
	}()

	t.Run("happy_path", func(t *testing.T) {
		downFile, err := ease.DownloadFile(testUri, testfile)
		assert.NoError(t, err)
		currentDir, err := os.Getwd()
		assert.NoError(t, err, "the current ")
		assert.Equal(t, currentDir+"/"+testfile, downFile,
			"downloading a file to a provided plain name file should download the file "+
				"to the current working directory (of this test)")
	})

	t.Run("invalid_url", func(t *testing.T) {
		uri := "invalid_scheme://test.local"
		_, err := ease.DownloadFile(uri, testfile, testGet)
		assert.Error(t, err)
	})

	t.Run("implicit tmp location for downloaded file", func(t *testing.T) {
		down, err := ease.DownloadFile(testUri, "", testGet)
		assert.NoError(t, err)
		assert.Regexp(t, regexp.MustCompile("/tmp/.*itler-go-lib-main/README.md"), down)
	})

}
