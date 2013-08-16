package summa

/*
#include <git2.h>
#include <git2/errors.h>
*/
import "C"

import (
	"runtime"
	"unsafe"
)

type IntBool bool
type UIntBool bool

type GitRepository struct {
	ptr *C.git_repository
}

type GitIndex struct {
	ptr *C.git_index
}

type GitOid struct {
	bytes [20]byte
}

type GitError struct {
	Message string
	Code    int
}

func init() {
	C.git_threads_init()
}

func (i IntBool) ToInt() C.int {
	if i {
		return C.int(1)
	}
	return C.int(0)
}

func (i UIntBool) ToInt() C.uint {
	if i {
		return C.uint(1)
	}
	return C.uint(0)
}

func (oid *GitOid) toC() *C.git_oid {
	return (*C.git_oid)(unsafe.Pointer(&oid.bytes))
}

func (e *GitError) Error() string {
	return e.Message
}

func GitErrorLast() error {
	err := C.giterr_last()
	if err == nil {
		return &GitError{"No message", 0}
	}
	return &GitError{C.GoString(err.message), int(err.klass)}
}

func GitRepositoryInit(path string, bare UIntBool) (*GitRepository, error) {
	repo := new(GitRepository)

	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	ret := C.git_repository_init(&repo.ptr, cpath, bare.ToInt())
	if ret < 0 {
		return nil, GitErrorLast()
	}

	runtime.SetFinalizer(repo, (*GitRepository).Free)
	return repo, nil
}

func GitRepositoryOpen(path string) (*GitRepository, error) {
	repo := new(GitRepository)

	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	ret := C.git_repository_open(&repo.ptr, cpath)
	if ret < 0 {
		return nil, GitErrorLast()
	}

	runtime.SetFinalizer(repo, (*GitRepository).Free)
	return repo, nil
}

func (r *GitRepository) Free() {
	runtime.SetFinalizer(r, nil)
	C.git_repository_free(r.ptr)
}

func (r *GitRepository) Commit(name, email string) error {
	var ret C.int

	var index *C.git_index
	ret = C.git_repository_index(&index, r.ptr)
	if ret < 0 {
		return GitErrorLast()
	}
	defer C.git_index_free(index)

	treeOid := new(C.git_oid)
	ret = C.git_index_write_tree(treeOid, index)
	if ret < 0 {
		return GitErrorLast()
	}

	tree := new(C.git_tree)
	ret = C.git_tree_lookup(&tree, r.ptr, treeOid)
	if ret < 0 {
		return GitErrorLast()
	}
	defer C.git_tree_free(tree)

	signature := new(C.git_signature)
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	cEmail := C.CString(email)
	defer C.free(unsafe.Pointer(cEmail))
	ret = C.git_signature_now(&signature, cName, cEmail)
	if ret < 0 {
		return GitErrorLast()
	}
	defer C.git_signature_free(signature)

	headOid := new(C.git_oid)
	cHead := C.CString("HEAD")
	defer C.free(unsafe.Pointer(cHead))
	ret = C.git_reference_name_to_id(headOid, r.ptr, cHead)

	commitOid := new(C.git_oid)
	cMessage := C.CString("")
	defer C.free(unsafe.Pointer(cMessage))

	if ret == 0 {
		head := new(C.git_commit)
		ret = C.git_commit_lookup(&head, r.ptr, headOid)
		if ret < 0 {
			return GitErrorLast()
		}
		defer C.git_commit_free(head)

		parents := make([]*C.git_commit, 1)
		parents[0] = head

		ret = C.git_commit_create(
			commitOid,
			r.ptr,
			cHead,
			signature,
			signature,
			nil,
			cMessage,
			tree,
			1,
			&parents[0],
		)
	} else {
		ret = C.git_commit_create(
			commitOid,
			r.ptr,
			cHead,
			signature,
			signature,
			nil,
			cMessage,
			tree,
			0,
			nil,
		)
	}

	if ret < 0 {
		return GitErrorLast()
	}

	ret = C.git_index_write(index)
	if ret < 0 {
		return GitErrorLast()
	}

	return nil
}

func (r *GitRepository) Index() (*GitIndex, error) {
	var ptr *C.git_index
	ret := C.git_repository_index(&ptr, r.ptr)
	if ret < 0 {
		return nil, GitErrorLast()
	}

	idx := &GitIndex{ptr}
	runtime.SetFinalizer(idx, (*GitIndex).Free)
	return idx, nil
}

func (i *GitIndex) Free() {
	runtime.SetFinalizer(i, nil)
	C.git_index_free(i.ptr)
}

func (i *GitIndex) Write() error {
	ret := C.git_index_write(i.ptr)
	if ret < 0 {
		return GitErrorLast()
	}
	return nil
}

func (i *GitIndex) WriteTree() (*GitOid, error) {
	oid := new(GitOid)
	ret := C.git_index_write_tree(oid.toC(), i.ptr)
	if ret < 0 {
		return nil, GitErrorLast()
	}

	return oid, nil
}

func (i *GitIndex) Add(file string) error {
	cfile := C.CString(file)
	defer C.free(unsafe.Pointer(cfile))

	ret := C.git_index_add_bypath(i.ptr, cfile)
	if ret < 0 {
		return GitErrorLast()
	}

	return i.Write()
}

func (i *GitIndex) Rm(file string) error {
	cfile := C.CString(file)
	defer C.free(unsafe.Pointer(cfile))

	ret := C.git_index_remove_bypath(i.ptr, cfile)
	if ret < 0 {
		return GitErrorLast()
	}

	return i.Write()
}
