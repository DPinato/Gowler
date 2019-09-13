package main

import (
	"reflect"
	"testing"
)

func TestReadInputArgs(t *testing.T) {
	assertCorrect := func(t *testing.T, got, want map[string]string) {
		t.Helper()
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Got %v, want %v", got, want)
		}
	}

	// test list
	t.Run("Correct input argument", func(t *testing.T) {
		got := readInputArgs([]string{"--csv", "/some/path/here.csv"})
		want := map[string]string{"--csv": "/some/path/here.csv"}
		assertCorrect(t, got, want)
	})

	t.Run("Multiple input arguments", func(t *testing.T) {
		got := readInputArgs([]string{"--csv", "/some/path/here.csv", "--num", "100"})
		want := map[string]string{"--csv": "/some/path/here.csv", "--num": "100"}
		assertCorrect(t, got, want)
	})

	t.Run("Bad input arguments 1", func(t *testing.T) {
		got := readInputArgs([]string{})
		want := map[string]string(nil)
		assertCorrect(t, got, want)
	})

	t.Run("Bad input arguments 2", func(t *testing.T) {
		got := readInputArgs([]string{"--csv"})
		want := map[string]string(nil)
		assertCorrect(t, got, want)
	})

	t.Run("Bad multiple input arguments", func(t *testing.T) {
		got := readInputArgs([]string{"--csv", "/some/path/here.csv", "-num", "100"})
		want := map[string]string(nil)
		assertCorrect(t, got, want)
	})

}

func TestReadCSVFile(t *testing.T) {
	assertCorrectLength := func(t *testing.T, got, want int) {
		t.Helper()
		if got != want {
			t.Errorf("Got %v, want %v", got, want)
		}
	}

	// test list
	t.Run("Parsing test.csv file all entries", func(t *testing.T) {
		got, _ := readCSVFile("test-files/test.csv")
		want := 14
		assertCorrectLength(t, len(got), want)
	})

	t.Run("Parsing test.csv file 5 entries", func(t *testing.T) {
		got, _ := readCSVFile("test-files/test.csv", 5)
		want := 5
		assertCorrectLength(t, len(got), want)
	})

}

func TestContains(t *testing.T) {
	assertCorrectMessage := func(t *testing.T, got, want bool) {
		t.Helper()
		if got != want {
			t.Errorf("Got %v, want %v", got, want)
		}
	}

	// test list
	t.Run("Find string in string array", func(t *testing.T) {
		got := Contains([]string{"ab", "cd", "ef"}, "ab")
		want := true
		assertCorrectMessage(t, got, want)
	})

	t.Run("Find non-existing string in string array", func(t *testing.T) {
		got := Contains([]string{"ab", "cd", "ef"}, "ae")
		want := false
		assertCorrectMessage(t, got, want)
	})

}

func TestContainsAt(t *testing.T) {
	assertCorrectMessage := func(t *testing.T, got, want int) {
		t.Helper()
		if got != want {
			t.Errorf("Got %v, want %v", got, want)
		}
	}

	// test list
	t.Run("Find string in string array", func(t *testing.T) {
		got := ContainsAt([]string{"ab", "cd", "ef"}, "ef")
		want := 2
		assertCorrectMessage(t, got, want)
	})

	t.Run("Find non-existing string in string array", func(t *testing.T) {
		got := ContainsAt([]string{"ab", "cd", "ef"}, "ae")
		want := -1
		assertCorrectMessage(t, got, want)
	})

}
