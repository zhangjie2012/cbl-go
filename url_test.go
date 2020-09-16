package cbl

import "testing"

func TestUrlJoin(t *testing.T) {
	{
		host := "https://github.com"
		t.Log(UrlJoin(host, "zhangjie2020"))
	}

	{
		host := "https://github.com/"
		t.Log(UrlJoin(host, "/zhangjie2020"))
	}

	{
		host := "https://github.com/"
		t.Log(UrlJoin(host, "/zhangjie2020", "/cbl-go/"))
	}
}
