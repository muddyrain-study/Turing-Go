package main

import (
	"os"
	"testing"
)

func TestCreate(t *testing.T) {
	f, err := os.Create("test.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
}

func TestMkdir(t *testing.T) {
	err := os.Mkdir("test", os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func TestMkdirAll(t *testing.T) {
	err := os.MkdirAll("test/test/test", os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func TestRemove(t *testing.T) {
	err := os.Remove("test.txt")
	if err != nil {
		panic(err)
	}
}

func TestRemoveAll(t *testing.T) {
	err := os.RemoveAll("test")
	if err != nil {
		panic(err)
	}
}

func TestGetWd(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	t.Log(dir)
}
func TestChdir(t *testing.T) {
	err := os.Chdir("test")
	if err != nil {
		panic(err)
	}
}

func TestTempDir(t *testing.T) {
	dir := os.TempDir()
	t.Log(dir)
}

func TestRename(t *testing.T) {
	err := os.Rename("test", "test1")
	if err != nil {
		panic(err)
	}
}

func TestChmod(t *testing.T) {
	err := os.Chmod("test1", 0777)
	if err != nil {
		panic(err)
	}
}
func TestChown(t *testing.T) {
	err := os.Chown("test1", os.Getuid(), os.Getgid())
	if err != nil {
		panic(err)
	}
}
