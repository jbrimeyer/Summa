package summa

import (
	"os"
	"path"
)

// repoPath will return the absolute path to the repository
func repoPath(id string) string {
	return path.Join(config.GitRoot(), id[:2], id[2:])
}

func repoCreate(id string, user *summaUser, files snippetFiles) error {
	var err error
	absPath := repoPath(id)
	err = os.MkdirAll(absPath, 0755)
	if err != nil {
		return err
	}

	defer (func() {
		if err != nil {
			repoDelete(id)
		}
	})()

	for _, file := range files {
		var f *os.File
		filePath := path.Join(absPath, file.Filename)
		f, err = os.Create(filePath)
		if err != nil {
			return err
		}
		_, err = f.WriteString(file.Contents)
		if err != nil {
			return err
		}
	}

	// TODO: git init
	// TODO: git add
	// TODO: git commit

	return nil
}

func repoUpdate(id string, user *summaUser, files snippetFiles) error {
	return nil
}

func repoDelete(id string) error {
	return nil
}