package main

import "testing"

func TestDigest(t *testing.T) {
	contents := map[string]string{
		"abc":                                    "a9993e364706816aba3e25717850c26c9cd0d89d",
		`fn main() { println!("hello world"); }`: "b99f10377346277fd92eadbd29bd18999ee6f8b0",
		"你好，世界":                                  "3becb03b015ed48050611c8d7afe4b88f70d5a20",
	}

	for in, want := range contents {
		if got := digest([]byte(in)); got != (want) {
			t.Errorf("got: [%s], want: [%s]", got, want)
		}
	}
}

func TestUnhexMixin(t *testing.T) {
	digests42 := map[string]string{
		"3becb03b015ed48050611c8d7afe4b88f70d5a20": "ADvssDsBXtSAUGEcjXr-S4j3DVog",
		"a9993e364706816aba3e25717850c26c9cd0d89d": "AKmZPjZHBoFquj4lcXhQwmyc0Nid",
		"b99f10377346277fd92eadbd29bd18999ee6f8b0": "ALmfEDdzRid_2S6tvSm9GJme5viw",
	}

	digests6 := map[string]string{
		"3becb03b015ed48050611c8d7afe4b88f70d5a20": "DVog",
		"a9993e364706816aba3e25717850c26c9cd0d89d": "0Nid",
		"b99f10377346277fd92eadbd29bd18999ee6f8b0": "5viw",
	}
	for in, want := range digests42 {
		if got, _ := unhexMixin(in, 42); got != want {
			t.Errorf("got: [%s], want: [%s]", got, want)
		}
	}

	for in, want := range digests6 {
		if got, _ := unhexMixin(in, 6); got != want {
			t.Errorf("got: [%s], want: [%s]", got, want)
		}
	}
}
