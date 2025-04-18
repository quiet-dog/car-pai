package guard

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/antchfx/xmlquery"
)

func GateawayControl(host, username, password string, data string) (*xmlquery.Node, error) {
	ctx := context.Background()
	r := NewDigestRequest(ctx, username, password) // username & password
	reqBody := bytes.NewBuffer([]byte(data))

	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s%s", host, "/ISAPI/Parking/channels/1/barrierGate"), reqBody)
	req.Header.Set("Content-Type", "application/xml")

	resp, _ := r.Do(req)
	defer resp.Body.Close()

	rsp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%s \n", rsp)
	return xmlquery.Parse(strings.NewReader(string(rsp)))
}
