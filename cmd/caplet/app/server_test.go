package app

import (
	"testing"
)

func TestRetrieveServer(t *testing.T) {
	if _, err := retrieveServer("non-exist"); err == nil {
		t.Error("error must exist")
	}
}

func TestGrpc(t *testing.T) {
	//out := os.Stdout
	//req := &proto.PullImageRequest{}
	//req.Images = append(req.Images, &proto.Image{
	//	Host: "softleader",
	//	Repo: "helm",
	//})
	//
	//ep := &caplet.Endpoint{
	//	Target: "localhost",
	//	Port:   caplet.DefaultPort,
	//}

	//if err := ep.PullImage(out, req); err != nil {
	//	t.Error(err)
	//}

	//if err := caplet.PullImage(out, []*caplet.Endpoint{ep}, req, dur.DefaultDeadlineSecond); err != nil {
	//	t.Error(err)
	//}
}
