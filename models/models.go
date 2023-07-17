package main

import (
	"encoding/hex"
	"strconv"

	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"

	"crypto/md5"

	"encoding/json"

	"github.com/schollz/progressbar/v3"
)

type modelnfo struct {
	Name     string `json:"name"`
	FileName string `json:"filename"`
	Url      string `json:"url"`
	MD5Sum   string `json:"md5sum"`
	Filesize string `json:"filesize"`
}

func (m *modelnfo) String() string {
	return m.Name
}

type tee struct {
	io.Reader
	io.Closer
}

func Tee(r io.ReadCloser, w io.Writer) io.ReadCloser {
	return tee{io.TeeReader(r, w), r}
}

func download(url string, filesize int64) (io.ReadCloser, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	bar := progressbar.DefaultBytes(
		filesize,
		"downloading",
	)

	return Tee(response.Body, bar), nil
}

func getModels() ([]*modelnfo, error) {
	var jr io.Reader
	f, err := os.Open(path.Join(ModelsDir, "models.json"))
	if err != nil {
		r, err := download(Models, -1)
		if err != nil {
			return nil, err
		}

		f, err := os.Create(path.Join(ModelsDir, "models.json"))
		if err != nil {
			return nil, err
		}
		defer f.Close()
		jr = io.TeeReader(r, f)
	} else {
		defer f.Close()
		jr = f
	}

	var models []*modelnfo
	d := json.NewDecoder(jr)
	err = d.Decode(&models)
	if err != nil {
		os.Remove(path.Join(ModelsDir, "models.json"))
		return nil, err
	}

	downloadable := make([]*modelnfo, 0, len(models))
	for _, m := range models {

		if u, err := url.Parse(m.Url); err == nil && u.Scheme == "https" {
			downloadable = append(downloadable, m)
		}
	}

	return downloadable, nil
}

func changedModel(m *modelnfo) bool {
	f, err := os.Open(path.Join(ModelsDir, m.FileName))
	if err != nil {
		return true
	}
	defer f.Close()

	h := md5.New()

	io.Copy(h, f)

	hash := hex.EncodeToString(h.Sum(nil))

	return hash != m.MD5Sum
}

func downloadModel(info *modelnfo) error {
	if !changedModel(info) {
		return nil
	}

	size, err := strconv.ParseInt(info.Filesize, 10, 64)
	if err != nil {
		size = -1
	}
	r, err := download(info.Url, size)
	if err != nil {
		return err
	}
	defer r.Close()

	f, err := os.Create(path.Join(ModelsDir, info.FileName))
	if err != nil {
		return err
	}
	defer f.Close()

	h := md5.New()

	io.Copy(f, io.TeeReader(r, h))

	hash := hex.EncodeToString(h.Sum(nil))
	if hash != info.MD5Sum {
		return fmt.Errorf("Download failed with hash mismatch %s != %s", hash, info.MD5Sum)
	}

	fmt.Println("\n")
	fmt.Println("Name:", info.Name)
	fmt.Println("File:", info.FileName)

	return nil
}
