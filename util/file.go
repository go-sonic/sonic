package util

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-sonic/sonic/util/xerr"
)

func ZipFile(dst string, srcs ...string) (err error) {
	// 创建准备写入的文件
	fw, err := os.Create(dst)
	if err != nil {
		return xerr.NoType.Wrap(err).WithMsg("create zip file err")
	}
	defer func() {
		if err = fw.Close(); err != nil {
			err = xerr.NoType.Wrap(err).WithMsg("close file")
		}
	}()
	// 通过 fw 来创建 zip.Write
	zw := zip.NewWriter(fw)
	defer func() {
		if err = zw.Close(); err != nil {
			err = xerr.NoType.Wrap(err).WithMsg("close zip file")
		}
	}()

	for _, src := range srcs {
		// 下面来将文件写入 zw ，因为有可能会有很多个目录及文件，所以递归处理
		err = filepath.Walk(src, func(path string, fi os.FileInfo, errBack error) (err error) {
			if errBack != nil {
				return errBack
			}

			// 通过文件信息，创建 zip 的文件信息
			fh, err := zip.FileInfoHeader(fi)
			if err != nil {
				return err
			}
			if path == src {
				fh.Name = filepath.Base(src)
			} else {
				fh.Name = filepath.Join(filepath.Base(src), strings.TrimPrefix(path, src))
			}
			// 替换文件信息中的文件名
			fh.Name = strings.TrimPrefix(fh.Name, string(filepath.Separator))

			// 这步开始没有加，会发现解压的时候说它不是个目录
			if fi.IsDir() {
				fh.Name += "/"
			}

			// 写入文件信息，并返回一个 Write 结构
			w, err := zw.CreateHeader(fh)
			if err != nil {
				return err
			}

			// 检测，如果不是标准文件就只写入头信息，不写入文件数据到 w
			// 如目录，也没有数据需要写
			if !fh.Mode().IsRegular() {
				return nil
			}

			// 打开要压缩的文件
			fr, err := os.Open(path)
			if err != nil {
				return err
			}
			defer fr.Close()

			// 将打开的文件 Copy 到 w
			_, err = io.Copy(w, fr)
			if err != nil {
				return err
			}
			// 输出压缩的内容

			return err
		})
		if err != nil {
			return xerr.NoType.Wrap(err).WithMsg("zip file err")
		}
	}
	return err
}

func Unzip(src string, dest string) ([]string, error) {
	r, err := zip.OpenReader(src)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	filenames := make([]string, 0, len(r.File))
	for _, f := range r.File {
		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			err := os.MkdirAll(fpath, os.ModePerm)
			if err != nil {
				return nil, xerr.WithStatus(err, xerr.StatusInternalServerError)
			}
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}

func CopyDir(srcPath, desPath string) error {
	if srcInfo, err := os.Stat(srcPath); err != nil {
		return err
	} else if !srcInfo.IsDir() {
		return xerr.WithMsg(nil, "src is not dir")
	}

	if err := MakeDir(desPath); err != nil {
		return err
	}
	if desInfo, err := os.Stat(desPath); err != nil {
		return err
	} else if !desInfo.IsDir() {
		return xerr.WithMsg(nil, "dest is not dir")
	}

	if strings.TrimSpace(srcPath) == strings.TrimSpace(desPath) {
		return xerr.WithMsg(nil, "srcPath=destPath")
	}

	err := filepath.Walk(srcPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}

		if path == srcPath {
			return nil
		}

		destNewPath := strings.ReplaceAll(path, srcPath, desPath)

		if !f.IsDir() {
			if _, err = CopyFile(path, destNewPath); err != nil {
				return err
			}
		} else if !FileIsExisted(destNewPath) {
			return MakeDir(destNewPath)
		}

		return nil
	})

	return err
}

func CopyFile(src, des string) (written int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer srcFile.Close()

	fi, _ := srcFile.Stat()
	perm := fi.Mode()

	desFile, err := os.OpenFile(des, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return 0, err
	}
	defer desFile.Close()

	return io.Copy(desFile, srcFile)
}

func FileIsExisted(filename string) bool {
	existed := true
	if _, err := os.Stat(filename); err != nil && os.IsNotExist(err) {
		existed = false
	}
	return existed
}

func MakeDir(dir string) error {
	if !FileIsExisted(dir) {
		if err := os.MkdirAll(dir, 0o777); err != nil { // os.ModePerm
			fmt.Println("MakeDir failed:", err)
			return err
		}
	}
	return nil
}
