package main

import (
	"compress/flate"
	"compress/gzip"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"SimpleIPLocation/internal/httpfs"
	"SimpleIPLocation/internal/utils"

	"github.com/andybalholm/brotli"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatal("no original dir path")
	}

	orgDir := args[1]
	pExist, err := utils.PathExist(orgDir)
	if err != nil {
		log.Fatal(err)
	}
	if !pExist {
		log.Fatal(orgDir, "was not exist")
	}
	staticDir := filepath.Dir(orgDir)

	cTypeList := []string{"gzip", "deflate", "br"}
	for _, compressType := range cTypeList {
		compressDir := filepath.Join(staticDir, compressType)
		err = os.MkdirAll(compressDir, utils.NormalDirPerm)
		if err != nil {
			log.Fatal("create dir ,err=", err)
		}
	}

	err = filepath.WalkDir(orgDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relPath, relErr := filepath.Rel(orgDir, path)
		if relErr != nil {
			return err
		}
		_ = relPath
		if d.IsDir() {
			// create sub dir
			for _, compressType := range cTypeList {
				compressDir := filepath.Join(staticDir, compressType)
				mkdirErr := os.MkdirAll(filepath.Join(compressDir, relPath), utils.NormalDirPerm)
				if mkdirErr != nil {
					return err
				}
			}
		} else {
			if !httpfs.IsNeedCompress(filepath.ToSlash(path)) {
				return nil
			}

			orgBytes, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			for _, compressType := range cTypeList {
				compressDir := filepath.Join(staticDir, compressType)
				const createFlag = os.O_RDWR | os.O_CREATE | os.O_TRUNC
				dstF, err := os.OpenFile(filepath.Join(compressDir, relPath), createFlag, utils.NormalFilePerm)
				if err != nil {
					return err
				}

				var cWriter io.WriteCloser
				switch compressType {
				case "gzip":
					gW, nWErr := gzip.NewWriterLevel(dstF, gzip.BestCompression)
					if nWErr != nil {
						return nWErr
					}
					cWriter = gW
				case "deflate":
					flateW, nWErr := flate.NewWriter(dstF, flate.BestCompression)
					if nWErr != nil {
						return nWErr
					}
					cWriter = flateW
				case "br":
					brW := brotli.NewWriterLevel(dstF, brotli.BestCompression)
					cWriter = brW
				}

				_, err = cWriter.Write(orgBytes)
				if err != nil {
					return err
				}

				err = cWriter.Close()
				if err != nil {
					return err
				}
				dstF.Close()
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
