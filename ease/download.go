package ease

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

// URLGetter conveniently describes a function
// used for getting a HTTP response based on a provided URL string
type URLGetter func(string) (*http.Response, error)

// DownloadFile downloads a file determined by uri to the specified filepath
// Make sure that the uri or environment this code runs in is handling
// the authentication for downloading the requested file
func DownloadFile(uri, file string, g ...URLGetter) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", fmt.Errorf("Provided URL '%s' is invalid: %w", uri, err)
	}

	if len(g) == 0 || g[0] == nil {
		g = []URLGetter{http.Get}
	}

	var tempDirCleanUp DeferFunc = func() {}
	if file == "" {
		tmpDir, cleanUp, err := CreateTmpDir(path.Dir(u.Path))
		if err != nil {
			return "", fmt.Errorf("Temp dir creation for download of '%s' failed: %w",
				fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, u.Path), err)
		}
		file = tmpDir + "/" + path.Base(u.Path)
		tempDirCleanUp = cleanUp
		log.Debug().Msgf("File download temp destination '%s' initialized", file)
	}

	dst, err := os.Create(file)
	if err != nil {
		defer tempDirCleanUp()
		return "", err
	}
	defer dst.Close()

	// Get the data
	resp, err := g[0](u.String())
	if err != nil {
		defer tempDirCleanUp()
		return "", err
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(dst, resp.Body)
	if err != nil {
		defer tempDirCleanUp()
		return "", err
	}

	file, err = filepath.Abs(file)
	if err != nil {
		defer tempDirCleanUp()
		return "", fmt.Errorf("Actual path '%s' seems to not exist: %w", file, err)
	}

	return file, nil
}
