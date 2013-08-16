package summa

import (
	"os"
	"path"
)

// repoPath will return the absolute path to the repository
func repoPath(id string) string {
	return path.Join(config.GitRoot(), id[:2], id[2:])
}

// repoCreate will create a new repository in the filesystem
func repoCreate(id string, u *User, files snippetFiles) error {
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

	repo, err := GitRepositoryInit(absPath, false)
	if err != nil {
		return err
	}

	index, err := repo.Index()
	if err != nil {
		return err
	}

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

		err = index.Add(file.Filename)
		if err != nil {
			return err
		}
	}

	return repo.Commit(u.DisplayName, u.Email)
}

func repoUpdate(id string, u *User, oldFiles, newFiles snippetFiles) error {
	absPath := repoPath(id)

	repo, err := GitRepositoryOpen(absPath)
	if err != nil {
		return err
	}

	index, err := repo.Index()
	if err != nil {
		return err
	}

	for _, file := range oldFiles {
		filePath := path.Join(absPath, file.Filename)
		err = os.Remove(filePath)
		if err != nil {
			return err
		}

		infoLog.Printf("Removing %s", file.Filename)
		err = index.Rm(file.Filename)
		if err != nil {
			return err
		}
	}

	for _, file := range newFiles {
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

		infoLog.Printf("Adding %s", file.Filename)
		err = index.Add(file.Filename)
		if err != nil {
			return err
		}
	}

	infoLog.Printf("Committing")

	return repo.Commit(u.DisplayName, u.Email)
}

// repoDelete will permanently delete the repository from the filesystem
func repoDelete(id string) error {
	return os.RemoveAll(repoPath(id))
}
