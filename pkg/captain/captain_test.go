package captain

import (
	"testing"
)

func TestGenerateScript(t *testing.T) {
	//path, err := ioutil.TempDir(os.TempDir(), "test-generate-script-")
	//if err != nil {
	//	t.Error(err)
	//}
	//defer os.RemoveAll(path)
	//
	//cmd := exec.Command("helm", "create", "foo")
	//cmd.Dir = path
	//if err := cmd.Run(); err != nil {
	//	t.Error(err)
	//}
	//
	//cmd = exec.Command("helm", "package", "foo")
	//cmd.Dir = path
	//if err := cmd.Run(); err != nil {
	//	t.Error(err)
	//}
	//
	//chart, err := ioutil.ReadFile(filepath.Join(path, "foo-0.1.0.tgz"))
	//if err != nil {
	//	t.Error(err)
	//}
	//
	//log := logger.New(os.Stdout)
	//err = GenerateScript(log, "192.168.1.93:30051", &proto.GenerateScriptRequest{
	//	Chart: &proto.Chart{
	//		Content:  chart,
	//		FileName: "foo-0.1.0.tgz",
	//	},
	//	Pull: true,
	//	Retag: &proto.ReTag{
	//	},
	//	Save: true,
	//	Load: true,
	//}, dur.DefaultDeadlineSecond)
	//if err != nil {
	//	t.Error(err)
	//}
}
