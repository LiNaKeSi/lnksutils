package lnksutils

import "os"

type fileOpts struct {
	atomic   bool
	filemode os.FileMode
}

type FileOption func(*fileOpts) error

// 保存文件时，先写入与目标文件同目录下的一个临时文件，完成后再重命名为最终文件。
func WithAtomicSave(opt *fileOpts) error {
	opt.atomic = true
	return nil
}

// 保持文件时，如果目标文件不存在或使用了WithAtomicSave，则会使用m作为新文件的权限。
// NOTE: 若目标文件已经存在且没有使用WithAtomicSave则此操作无任何效果。
func WithFileMode(m os.FileMode) FileOption {
	return func(opt *fileOpts) error {
		opt.filemode = m
		return nil
	}
}

func initFileOpts(opts []FileOption) (*fileOpts, error) {
	ret := &fileOpts{
		filemode: os.FileMode(0600),
	}
	for _, opt := range opts {
		err := opt(ret)
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}
