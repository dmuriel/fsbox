package fsbox

import (
	"embed"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

//go:embed testdata
var fsys embed.FS

func TestBox(t *testing.T) {
	b := New(fsys, "testdata")

	if b == nil {
		t.Errorf("should return a new box")
	}
}

func TestOpen(t *testing.T) {
	b := New(fsys, "testdata")

	t.Run("Opening a.txt", func(t *testing.T) {
		f, err := b.Open("a.txt")
		if err != nil {
			t.Errorf("should not be error reading a.txt, got %v", err)
			return
		}

		dat, err := io.ReadAll(f)
		if err != nil {
			t.Errorf("should not be error reading dat, got %v", err)
			return
		}

		if string(dat) != "file a" {
			t.Errorf("bad data")
			return
		}
	})

	t.Run("opening non-present file", func(t *testing.T) {
		_, err := b.Open("x.txt")
		if err == nil {
			t.Errorf("should return an error opening")
			return
		}
	})
}

func TestAddString(t *testing.T) {
	b := New(fsys, "testdata")

	err := b.AddString("a", "a")
	if err != ErrCannotAdd {
		t.Errorf("should return ErrCannotAdd but got %v", err)
		return
	}
}

func TestAddBytes(t *testing.T) {
	b := New(fsys, "testdata")

	err := b.AddBytes("a", []byte("a"))
	if err != ErrCannotAdd {
		t.Errorf("should return ErrCannotAdd but got %v", err)
		return
	}
}

func TestList(t *testing.T) {
	b := New(fsys, "testdata")

	received := b.List()
	expected := []string{
		"testdata/a.txt",
		"testdata/b.txt",
	}

	if strings.Join(received, "|") != strings.Join(expected, "|") {
		t.Errorf("should return the correct list of paths, got %v", received)
		return
	}
}

func TestFind(t *testing.T) {
	t.Run("File Exists", func(t *testing.T) {
		f, _ := fsys.Open("testdata/a.txt")
		xs, _ := ioutil.ReadAll(f)

		b := New(fsys, "testdata")
		bs, err := b.Find("a.txt")
		if err != nil {
			t.Errorf("should got no error with a.txt, got %v", err)
			return
		}

		if string(bs) != string(xs) {
			t.Errorf("contents seem invalid")
			return
		}
	})

	t.Run("File is not there", func(t *testing.T) {
		b := New(fsys, "testdata")
		_, err := b.Find("x.txt")
		if err == nil {
			t.Errorf("should got error with x.txt, got nil")
			return
		}
	})

}

func TestFindString(t *testing.T) {
	t.Run("File Exists", func(t *testing.T) {
		f, _ := fsys.Open("testdata/a.txt")
		xs, _ := ioutil.ReadAll(f)

		b := New(fsys, "testdata")
		bs, err := b.FindString("a.txt")
		if err != nil {
			t.Errorf("should got no error with a.txt, got %v", err)
			return
		}

		if bs != string(xs) {
			t.Errorf("contents seem invalid")
			return
		}
	})

	t.Run("File is not there", func(t *testing.T) {
		b := New(fsys, "testdata")
		_, err := b.FindString("x.txt")
		if err == nil {
			t.Errorf("should got error with x.txt, got nil")
			return
		}
	})

}

func TestHas(t *testing.T) {
	t.Run("File Exists", func(t *testing.T) {
		b := New(fsys, "testdata")
		has := b.Has("testdata/a.txt")
		if !has {
			t.Errorf("should return true for has(a.txt)")
		}
	})
}
