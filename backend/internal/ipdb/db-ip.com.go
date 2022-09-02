package ipdb

import (
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"SimpleIPLocation/internal/utils"

	"github.com/fujiwara/shapeio"
)

// downloadDBIP 從 db-ip.com 下載 ip 資料庫
func downloadDBIP(dstPath string, bytesPerSec float64) error {
	// 建立要寫入的檔案

	dstFile, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, utils.NormalFilePerm)
	if err != nil {
		return errors.New("create new ipdb file failed, err=" + err.Error())
	}
	defer dstFile.Close()

	// 生成下載 URL

	nowTime := time.Now()
	donwloadURL := fmt.Sprintf("https://download.db-ip.com/free/dbip-city-lite-%d-%02d.mmdb.gz", nowTime.Year(), nowTime.Month())

	req, _ := http.NewRequestWithContext(context.Background(), "GET", donwloadURL, http.NoBody)
	req.Header.Add("user-agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.New("client.Do failed, err=" + err.Error())
	}
	defer resp.Body.Close()

	// 限制下載速率

	var reader io.Reader
	if bytesPerSec > 0 {
		sioR := shapeio.NewReader(resp.Body)
		sioR.SetRateLimit(bytesPerSec)
		reader = sioR
	} else {
		reader = resp.Body
	}

	// 進行 gzip 解壓縮

	gzReader, err := gzip.NewReader(reader)
	if err != nil {
		return errors.New("gzip.NewReader failed, err=" + err.Error())
	}
	defer gzReader.Close()

	// 開始寫入檔案

	const copyBufSize int64 = 32 * 1024
	for {
		_, err = io.CopyN(dstFile, gzReader, copyBufSize)
		if err != nil {
			if err == io.EOF {
				break
			}
			return errors.New("write ipdb failed, err" + err.Error())
		}
	}

	return nil
}
